"""Measure CPU image generation through the pinned ComfyUI launcher."""

from __future__ import annotations

import argparse
import contextlib
import hashlib
import json
import os
import pathlib
import socket
import subprocess
import time
import typing
import urllib.error
import urllib.request
import uuid

import PIL.Image
from python.runfiles import runfiles


def _mapping(value: object) -> dict[str, object]:
    """Validate a JSON object at the HTTP boundary."""
    if not isinstance(value, dict) or not all(
        isinstance(key, str) for key in value
    ):
        raise ValueError("Expected a JSON object")
    return dict(value)


def _request(
    url: str, data: typing.Optional[dict[str, object]] = None
) -> dict[str, object]:
    """Exchange a bounded JSON request with the local inference server."""
    payload: typing.Optional[bytes] = (
        json.dumps(data).encode() if data is not None else None
    )
    request: urllib.request.Request = urllib.request.Request(
        url, data=payload, headers={"Content-Type": "application/json"}
    )
    opener: urllib.request.OpenerDirector = urllib.request.build_opener(
        urllib.request.ProxyHandler({})
    )
    with opener.open(request, timeout=5) as response:
        return _mapping(json.loads(response.read(8 * 1024 * 1024)))


def _resolve(value: str) -> pathlib.Path:
    """Resolve a declared Bazel runfile without relying on the caller's cwd."""
    resolver: typing.Optional[runfiles.Runfiles] = runfiles.Create()
    if resolver is None:
        raise RuntimeError("Bazel runfiles are unavailable")
    location: typing.Optional[str] = resolver.Rlocation(value)
    if location is None:
        raise FileNotFoundError(value)
    return pathlib.Path(location).resolve(strict=True)


def _workflow(
    prompt: str, seed: int, steps: int, size: int
) -> dict[str, object]:
    """Construct the SD-Turbo API workflow at the requested size."""
    return {
        "1": {
            "class_type": "CheckpointLoaderSimple",
            "inputs": {"ckpt_name": "sd_turbo.safetensors"},
        },
        "2": {
            "class_type": "CLIPTextEncode",
            "inputs": {"text": prompt, "clip": ["1", 1]},
        },
        "3": {
            "class_type": "CLIPTextEncode",
            "inputs": {"text": "", "clip": ["1", 1]},
        },
        "4": {
            "class_type": "EmptyLatentImage",
            "inputs": {"width": size, "height": size, "batch_size": 1},
        },
        "5": {
            "class_type": "KSampler",
            "inputs": {
                "model": ["1", 0],
                "positive": ["2", 0],
                "negative": ["3", 0],
                "latent_image": ["4", 0],
                "seed": seed,
                "steps": steps,
                "cfg": 1.0,
                "sampler_name": "euler",
                "scheduler": "normal",
                "denoise": 1.0,
            },
        },
        "6": {
            "class_type": "VAEDecode",
            "inputs": {"samples": ["5", 0], "vae": ["1", 2]},
        },
        "7": {
            "class_type": "SaveImage",
            "inputs": {"images": ["6", 0], "filename_prefix": "benchmark"},
        },
    }


def _z_image_workflow(
    prompt: str, seed: int, steps: int, size: int
) -> dict[str, object]:
    """Construct the upstream Z-Image-Turbo workflow using built-in nodes."""
    return {
        "1": {
            "class_type": "UNETLoader",
            "inputs": {
                "unet_name": "z_image_turbo_bf16.safetensors",
                "weight_dtype": "default",
            },
        },
        "2": {
            "class_type": "CLIPLoader",
            "inputs": {
                "clip_name": "qwen_3_4b.safetensors",
                "type": "lumina2",
                "device": "cpu",
            },
        },
        "3": {
            "class_type": "VAELoader",
            "inputs": {"vae_name": "ae.safetensors"},
        },
        "4": {
            "class_type": "CLIPTextEncode",
            "inputs": {"clip": ["2", 0], "text": prompt},
        },
        "5": {
            "class_type": "ConditioningZeroOut",
            "inputs": {"conditioning": ["4", 0]},
        },
        "6": {
            "class_type": "EmptySD3LatentImage",
            "inputs": {"width": size, "height": size, "batch_size": 1},
        },
        "7": {
            "class_type": "ModelSamplingAuraFlow",
            "inputs": {"model": ["1", 0], "shift": 3.0},
        },
        "8": {
            "class_type": "KSampler",
            "inputs": {
                "model": ["7", 0],
                "positive": ["4", 0],
                "negative": ["5", 0],
                "latent_image": ["6", 0],
                "seed": seed,
                "steps": steps,
                "cfg": 1.0,
                "sampler_name": "res_multistep",
                "scheduler": "simple",
                "denoise": 1.0,
            },
        },
        "9": {
            "class_type": "VAEDecode",
            "inputs": {"samples": ["8", 0], "vae": ["3", 0]},
        },
        "10": {
            "class_type": "SaveImage",
            "inputs": {"images": ["9", 0], "filename_prefix": "benchmark"},
        },
    }


def _wait_ready(
    server: subprocess.Popen[bytes],
    endpoint: str,
    directory: pathlib.Path,
    timeout: float,
) -> None:
    """Wait for the owned server, failing on exit or the readiness deadline."""
    deadline: float = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if server.poll() is not None:
            raise RuntimeError("ComfyUI exited; inspect server.log")
        try:
            stats: dict[str, object] = _request(endpoint + "/system_stats")
            system: dict[str, object] = _mapping(stats.get("system"))
            arguments: object = system.get("argv")
            if (
                not isinstance(arguments, list)
                or str(directory) not in arguments
            ):
                raise RuntimeError(
                    "Port belongs to a different ComfyUI server"
                )
            return
        except (urllib.error.URLError, TimeoutError):
            time.sleep(0.25)
    raise TimeoutError("ComfyUI readiness deadline exceeded")


def _generate(
    endpoint: str,
    workflow: dict[str, object],
    server: subprocess.Popen[bytes],
    timeout: float,
) -> dict[str, object]:
    """Submit a workflow and require successful execution before returning."""
    response: dict[str, object] = _request(
        endpoint + "/prompt", {"prompt": workflow}
    )
    identifier: object = response.get("prompt_id")
    if not isinstance(identifier, str):
        raise RuntimeError("ComfyUI did not accept the workflow")
    deadline: float = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if server.poll() is not None:
            raise RuntimeError("ComfyUI exited during generation")
        history: dict[str, object] = _request(
            endpoint + "/history/" + identifier
        )
        if identifier in history:
            result: dict[str, object] = _mapping(history[identifier])
            status: dict[str, object] = _mapping(result.get("status"))
            if status.get("status_str") != "success" or not status.get(
                "completed"
            ):
                raise RuntimeError("Workflow failed: " + json.dumps(status))
            return result
        time.sleep(0.25)
    raise TimeoutError("Image-generation deadline exceeded")


def _stop(server: subprocess.Popen[bytes]) -> None:
    """Terminate only the owned inference process and wait for cleanup."""
    if server.poll() is None:
        server.terminate()
        try:
            server.wait(timeout=15)
        except subprocess.TimeoutExpired:
            server.kill()
            server.wait(timeout=5)


def _main() -> None:
    """Prepare writable state, launch CPU inference, and preserve evidence."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(__doc__)
    parser.add_argument("--launcher", required=True)
    parser.add_argument("--server", required=True)
    parser.add_argument(
        "--model", action="append", required=True, metavar="MODEL_PATH=RUNFILE"
    )
    parser.add_argument(
        "--profile", choices=["sd-turbo", "z-image-turbo"], default="sd-turbo"
    )
    parser.add_argument("--size", type=int, default=512)
    parser.add_argument("--directory", default="out/comfyui/runs")
    parser.add_argument("--threads", type=int, default=8)
    parser.add_argument("--memory-limit-gib", type=int, default=0)
    parser.add_argument("--steps", type=int)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--timeout", type=float, default=900)
    parser.add_argument(
        "--prompt",
        default="A red ceramic teapot on a wooden table beside a window, "
        "soft morning light, watercolor painting",
    )
    args: argparse.Namespace = parser.parse_args()
    if args.steps is None:
        args.steps = 8 if args.profile == "z-image-turbo" else 1
    if args.steps < 1 or args.steps > (
        50 if args.profile == "z-image-turbo" else 4
    ):
        parser.error("steps are outside the selected model's supported range")
    if args.size < 64 or args.size > 2048 or args.size % 64:
        parser.error("size must be a multiple of 64 between 64 and 2048")
    if args.threads < 1 or args.timeout <= 0:
        parser.error("threads and timeout must be positive")
    workspace: pathlib.Path = pathlib.Path(
        os.environ.get("BUILD_WORKSPACE_DIRECTORY", os.getcwd())
    )
    directory: pathlib.Path = (workspace / args.directory).resolve()
    directory = directory / (
        time.strftime("%Y%m%dT%H%M%S") + "-" + uuid.uuid4().hex[:8]
    )
    directory.mkdir(parents=True, exist_ok=True)
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        port: int = listener.getsockname()[1]
    command: list[str] = [
        str(_resolve(args.launcher)),
        "--server",
        args.server,
        "--directory",
        str(directory),
        "--port",
        str(port),
        "--threads",
        str(args.threads),
    ]
    if args.memory_limit_gib:
        command.extend(["--memory-limit-gib", str(args.memory_limit_gib)])
    specification: str
    for specification in args.model:
        command.extend(["--model", specification])
    print("ComfyUI state: " + str(directory), flush=True)
    started: float = time.monotonic()
    with contextlib.ExitStack() as resources:
        log: typing.BinaryIO = resources.enter_context(
            (directory / "server.log").open("wb")
        )
        server: subprocess.Popen[bytes] = subprocess.Popen(
            command, stdout=log, stderr=subprocess.STDOUT
        )
        resources.callback(_stop, server)
        endpoint: str = "http://127.0.0.1:" + str(port)
        _wait_ready(server, endpoint, directory, min(args.timeout, 180))
        ready: float = time.monotonic()
        workflow: dict[str, object] = (
            _z_image_workflow(args.prompt, args.seed, args.steps, args.size)
            if args.profile == "z-image-turbo"
            else _workflow(args.prompt, args.seed, args.steps, args.size)
        )
        (directory / "workflow.json").write_text(
            json.dumps(workflow, indent=2) + "\n"
        )
        history: dict[str, object] = _generate(
            endpoint, workflow, server, args.timeout
        )
        generated: float = time.monotonic()
        (directory / "history.json").write_text(
            json.dumps(history, indent=2) + "\n"
        )
        images: list[pathlib.Path] = list(
            (directory / "output").glob("benchmark_*.png")
        )
        if len(images) != 1:
            raise RuntimeError("Expected exactly one generated PNG")
        with PIL.Image.open(images[0]) as result:
            result.load()
            if result.format != "PNG" or result.size != (args.size, args.size):
                raise RuntimeError("Invalid generated image")
        model_hashes: dict[str, str] = {}
        for specification in args.model:
            relative_name: str
            runfile: str
            relative_name, _, runfile = specification.partition("=")
            with _resolve(runfile).open("rb") as model_file:
                model_hashes[relative_name] = hashlib.file_digest(
                    model_file, "sha256"
                ).hexdigest()
        report: dict[str, object] = {
            "prompt": args.prompt,
            "seed": args.seed,
            "steps": args.steps,
            "threads": args.threads,
            "backend": "cpu",
            "profile": args.profile,
            "width": args.size,
            "height": args.size,
            "model_sha256": next(iter(model_hashes.values())),
            "models_sha256": model_hashes,
            "startup_seconds": ready - started,
            "generation_seconds": generated - ready,
            "image": str(images[0]),
            "system_stats": _request(endpoint + "/system_stats"),
        }
        (directory / "report.json").write_text(
            json.dumps(report, indent=2) + "\n"
        )
        print(json.dumps(report, indent=2), flush=True)


if __name__ == "__main__":
    try:
        _main()
    except KeyboardInterrupt:
        raise SystemExit(130) from None
