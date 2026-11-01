"""Prepare local CPU inference state and execute the pinned ComfyUI server."""

from __future__ import annotations

import argparse
import os
import pathlib
import typing

import comfyui_frontend_package
from python.runfiles import runfiles


def _resolve(value: str) -> pathlib.Path:
    """Resolve a declared Bazel runfile without relying on the caller's cwd."""
    resolver: typing.Optional[runfiles.Runfiles] = runfiles.Create()
    if resolver is None:
        raise RuntimeError("Bazel runfiles are unavailable")
    location: typing.Optional[str] = resolver.Rlocation(value)
    if location is None:
        raise FileNotFoundError(value)
    return pathlib.Path(location).resolve(strict=True)


def _require_memory_limit(limit_gib: int) -> None:
    """Require an effective cgroup v2 RAM ceiling and disabled process swap."""
    entries: list[str] = (
        pathlib.Path("/proc/self/cgroup").read_text().splitlines()
    )
    relative_name: str = next(
        (entry[3:] for entry in entries if entry.startswith("0::")), ""
    )
    relative: pathlib.PurePosixPath = pathlib.PurePosixPath(
        relative_name.lstrip("/")
    )
    if not relative_name or ".." in relative.parts:
        raise RuntimeError("Inference requires a cgroup v2 memory scope")
    mount: pathlib.Path = pathlib.Path("/sys/fs/cgroup")
    current: pathlib.Path = mount / relative
    memory_bounded: bool = False
    swap_disabled: bool = False
    directory: pathlib.Path
    for directory in (current, *current.parents):
        if directory != mount and mount not in directory.parents:
            break
        memory_file: pathlib.Path = directory / "memory.max"
        swap_file: pathlib.Path = directory / "memory.swap.max"
        if memory_file.exists():
            memory_value: str = memory_file.read_text().strip()
            if (
                memory_value != "max"
                and int(memory_value) <= limit_gib * 1024**3
            ):
                memory_bounded = True
        if swap_file.exists() and swap_file.read_text().strip() == "0":
            swap_disabled = True
    if not memory_bounded or not swap_disabled:
        raise RuntimeError(
            f"Inference requires a cgroup RAM ceiling <= {limit_gib} GiB and zero swap; "
            "use the bounded systemd-run command in tools/comfyui/README.md"
        )


def _main() -> None:
    """Configure writable state and replace this process with upstream ComfyUI."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(__doc__)
    parser.add_argument("--server", required=True)
    parser.add_argument(
        "--model", action="append", required=True, metavar="MODEL_PATH=RUNFILE"
    )
    parser.add_argument("--directory", default="out/comfyui/session")
    parser.add_argument("--port", type=int, default=8188)
    parser.add_argument("--threads", type=int)
    parser.add_argument("--memory-limit-gib", type=int, default=0)
    args: argparse.Namespace = parser.parse_args()
    if args.threads is not None and args.threads < 1:
        parser.error("threads must be positive")
    if args.memory_limit_gib < 0:
        parser.error("memory limit must not be negative")
    if args.memory_limit_gib:
        _require_memory_limit(args.memory_limit_gib)
    frontend_module: typing.Optional[str] = comfyui_frontend_package.__file__
    if frontend_module is None:
        raise RuntimeError("The bundled frontend has no filesystem location")
    frontend_root: pathlib.Path = (
        pathlib.Path(frontend_module).resolve(strict=True).parent / "static"
    )
    workspace: pathlib.Path = pathlib.Path(
        os.environ.get("BUILD_WORKSPACE_DIRECTORY", os.getcwd())
    )
    directory: pathlib.Path = (workspace / args.directory).resolve()
    directory.mkdir(parents=True, exist_ok=True)
    cache_directory: pathlib.Path = directory / "cache"
    temporary_directory: pathlib.Path = directory / "temp"
    cache_directory.mkdir(exist_ok=True)
    temporary_directory.mkdir(exist_ok=True)
    specification: str
    for specification in args.model:
        relative_name: str
        runfile: str
        separator: str
        relative_name, separator, runfile = specification.partition("=")
        relative: pathlib.PurePosixPath = pathlib.PurePosixPath(relative_name)
        if (
            not separator
            or relative.is_absolute()
            or ".." in relative.parts
            or len(relative.parts) < 2
        ):
            parser.error(
                "model must be a relative category/file path and a declared runfile"
            )
        destination: pathlib.Path = directory / "models" / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        model: pathlib.Path = _resolve(runfile)
        if destination.is_symlink():
            if destination.resolve() != model:
                raise RuntimeError("Existing model links to another input")
        elif destination.exists():
            raise RuntimeError("Refusing to replace an existing model")
        else:
            destination.symlink_to(model)
    port: int = args.port
    environment: dict[str, str] = dict(os.environ)
    environment.update(
        {
            "HF_HUB_OFFLINE": "1",
            "TRANSFORMERS_OFFLINE": "1",
            "HF_HUB_DISABLE_TELEMETRY": "1",
            "HF_HOME": str(cache_directory / "huggingface"),
            "TORCH_HOME": str(cache_directory / "torch"),
            "XDG_CACHE_HOME": str(cache_directory),
            "TMPDIR": str(temporary_directory),
            "TMP": str(temporary_directory),
            "TEMP": str(temporary_directory),
            "PYTHONDONTWRITEBYTECODE": "1",
        }
    )
    if args.threads is not None:
        environment.update(
            {
                "OMP_NUM_THREADS": str(args.threads),
                "MKL_NUM_THREADS": str(args.threads),
                "OPENBLAS_NUM_THREADS": str(args.threads),
            }
        )
    command: list[str] = [
        str(_resolve(args.server)),
        "--cpu",
        "--force-fp32",
        "--offline",
        "--disable-all-custom-nodes",
        "--disable-auto-launch",
        "--front-end-root",
        str(frontend_root),
        "--listen",
        "127.0.0.1",
        "--port",
        str(port),
        "--base-directory",
        str(directory),
    ]
    print("ComfyUI state: " + str(directory), flush=True)
    os.execve(command[0], command, environment)


if __name__ == "__main__":
    _main()
