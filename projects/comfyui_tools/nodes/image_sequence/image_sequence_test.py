"""Exercise scene resets, native sampling feedback and repeatable still outputs.

The fixture replaces GPU model boundaries, not the sequence implementation.
Saved PNGs and a JSON receipt expose order, seeds, feedback and native metadata.
"""

import json
import os
import pathlib
import signal
import sys
import tempfile
import time
import types
import typing
import unittest
from unittest import mock

import numpy
import PIL.Image
import PIL.ImageSequence
import PIL.PngImagePlugin
import torch

from projects.comfyui_tools.nodes.image_sequence import image_sequence


class _NativeFixture:
    """Provide deterministic GPU boundaries and real still-image serialization."""

    def __init__(self, directory: pathlib.Path) -> None:
        """Bind the output directory and expose every sampled transition."""
        self.directory: pathlib.Path = directory
        """Fixture-owned output directory."""
        self.calls: list[dict[str, object]] = []
        """Ordered observations at the native sampler boundary."""
        self.interrupt_after: int = 0
        """Completed image count after which cancellation is requested."""
        self.guides: list[dict[str, object]] = []
        """ControlNet input and application observations for this run."""

    def control(
        self,
        positive: object,
        negative: object,
        control_net: object,
        image: torch.Tensor,
        strength: float,
        start_percent: float,
        end_percent: float,
        vae: object = None,
    ) -> tuple[object, object]:
        """Record guide selection while exposing conditioning to the sampler."""
        self.guides.append(
            {
                "source": image.mean().item(),
                "strength": strength,
                "start": start_percent,
                "end": end_percent,
                "model": control_net,
            }
        )
        return ({"prompt": positive, "guide": len(self.guides)}, negative)

    def sample(
        self, *args: object, **kwargs: object
    ) -> tuple[dict[str, torch.Tensor]]:
        """Record actual feedback and produce a different deterministic latent."""
        latent: dict[str, torch.Tensor] = typing.cast(
            dict[str, torch.Tensor], args[8]
        )
        self.calls.append(
            {
                "seed": args[1],
                "source": latent["samples"].mean().item(),
                "denoise": kwargs["denoise"],
                "positive": args[6],
                "model": args[0],
            }
        )
        return (
            {
                "samples": torch.full_like(
                    latent["samples"], len(self.calls) / 10
                )
            },
        )

    def save_images(
        self,
        images: torch.Tensor,
        filename_prefix: str,
        prompt: object = None,
        extra_pnginfo: dict[str, object] | None = None,
    ) -> dict[str, object]:
        """Serialize pixels and native metadata so downstream artifacts are inspectable."""
        filename: str = filename_prefix + ".png"
        metadata: PIL.PngImagePlugin.PngInfo = PIL.PngImagePlugin.PngInfo()
        for key, value in (extra_pnginfo or {}).items():
            metadata.add_text(key, json.dumps(value))
        metadata.add_text("prompt", json.dumps(prompt))
        PIL.Image.fromarray(
            (images[0].numpy() * 255).astype(numpy.uint8)
        ).save(self.directory / filename, pnginfo=metadata)
        return {
            "ui": {
                "images": [
                    {"filename": filename, "subfolder": "", "type": "output"}
                ]
            }
        }

    def interrupt(self) -> None:
        """Model the native cancellation boundary between completed images."""
        if (
            self.interrupt_after
            and len(list(self.directory.glob("*.png"))) >= self.interrupt_after
        ):
            raise InterruptedError("Fixture cancellation")

    def modules(self) -> dict[str, object]:
        """Supply native node interfaces while retaining real tensors and files."""
        management: types.SimpleNamespace = types.SimpleNamespace(
            throw_exception_if_processing_interrupted=self.interrupt
        )
        return {
            "folder_paths": types.SimpleNamespace(
                get_output_directory=lambda: str(self.directory)
            ),
            "nodes": types.SimpleNamespace(
                common_ksampler=self.sample,
                CLIPTextEncode=lambda: types.SimpleNamespace(
                    encode=lambda clip, text: (text,)
                ),
                EmptyLatentImage=lambda: types.SimpleNamespace(
                    generate=lambda width, height, batch_size: (
                        {
                            "samples": torch.zeros(
                                batch_size, 4, height // 8, width // 8
                            )
                        },
                    )
                ),
                VAEEncode=lambda: types.SimpleNamespace(
                    encode=lambda vae, pixels: (
                        {
                            "samples": torch.full(
                                (
                                    1,
                                    4,
                                    pixels.shape[1] // 8,
                                    pixels.shape[2] // 8,
                                ),
                                pixels.mean().item(),
                            )
                        },
                    )
                ),
                VAEDecode=lambda: types.SimpleNamespace(
                    decode=lambda vae, latent: (
                        torch.full(
                            (
                                1,
                                latent["samples"].shape[2] * 8,
                                latent["samples"].shape[3] * 8,
                                3,
                            ),
                            latent["samples"].mean().item(),
                        ),
                    )
                ),
                ImageScale=lambda: types.SimpleNamespace(
                    upscale=lambda image, method, width, height, crop: (
                        torch.nn.functional.interpolate(
                            image.movedim(-1, 1),
                            size=(height, width),
                            mode="nearest",
                        ).movedim(1, -1),
                    )
                ),
                SaveImage=lambda: types.SimpleNamespace(
                    save_images=self.save_images
                ),
                ControlNetApplyAdvanced=lambda: types.SimpleNamespace(
                    apply_controlnet=self.control
                ),
                interrupt_processing=self.interrupt,
            ),
            "comfy_extras": types.SimpleNamespace(
                nodes_canny=types.SimpleNamespace(
                    Canny=types.SimpleNamespace(
                        execute=lambda image, low_threshold, high_threshold: (
                            image,
                        )
                    )
                )
            ),
            "comfy_extras.nodes_canny": types.SimpleNamespace(
                Canny=types.SimpleNamespace(
                    execute=lambda image, low_threshold, high_threshold: (
                        image,
                    )
                )
            ),
            "comfy": types.SimpleNamespace(model_management=management),
            "comfy.model_management": management,
        }


class _ImageSequenceTest(unittest.TestCase):
    """Test observable sequence behavior before any artistic-quality claim."""

    def test_composer_two_pass_outputs_and_installed_lora(self) -> None:
        """Compose two scenes through native boundaries and decode real outputs."""
        with tempfile.TemporaryDirectory() as temporary:
            directory: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
            root: pathlib.Path = (
                pathlib.Path(directory) / "composer"
                if directory
                else pathlib.Path(temporary)
            )
            root.mkdir(parents=True, exist_ok=True)
            fixture: _NativeFixture = _NativeFixture(root)
            modules: dict[str, object] = fixture.modules()
            loads: list[tuple[str, float, float]] = []

            def load_lora(
                model: object,
                clip: object,
                name: str,
                model_strength: float,
                clip_strength: float,
            ) -> tuple[object, object]:
                """Observe the native model/encoder boundary without model weights."""
                loads.append((name, model_strength, clip_strength))
                return "adapted-model", "adapted-clip"

            typing.cast(types.SimpleNamespace, modules["nodes"]).LoraLoader = (
                lambda: types.SimpleNamespace(load_lora=load_lora)
            )
            typing.cast(
                types.SimpleNamespace, modules["folder_paths"]
            ).get_filename_list = lambda kind: ["character.safetensors"]
            with mock.patch.dict("sys.modules", modules):
                composer: image_sequence.ScenePrompts = (
                    image_sequence.ScenePrompts()
                )
                shared, actions = composer.compose(
                    character_prompt="adult woman, black hair, green eyes",
                    clothing_prompt="red coat, black boots",
                    style_prompt="best quality, anime illustration",
                    scene_settings="snowy park\nwarm library",
                    actions="standing\nwaving\n---\nsitting\nreading",
                )
                original: object = object()
                self.assertEqual(
                    image_sequence.OptionalLora().apply(
                        original, original, "(none)", 1, 1
                    ),
                    (original, original),
                )
                model, clip = image_sequence.OptionalLora().apply(
                    original, original, "character.safetensors", 0.7, 0.8
                )
                result: dict[str, object] = (
                    image_sequence.ImageSequence().generate(
                        model=model,
                        clip=clip,
                        vae=object(),
                        shared_prompt=shared,
                        actions=actions,
                        negative_prompt="low quality",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=42,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="composer",
                        generation_mode="independent",
                        seed_mode="fixed",
                        refiner_steps=3,
                        refiner_denoise=0.4,
                        control_net=object(),
                        control_strength=1,
                    )
                )
            self.assertEqual(loads, [("character.safetensors", 0.7, 0.8)])
            self.assertEqual(len(fixture.calls), 8)
            self.assertEqual(
                [call["source"] for call in fixture.calls[::2]], [0, 0, 0, 0]
            )
            self.assertEqual(
                [
                    round(typing.cast(float, call["source"]), 1)
                    for call in fixture.calls[1::2]
                ],
                [0.1, 0.3, 0.5, 0.7],
            )
            self.assertEqual(
                [call["seed"] for call in fixture.calls], [42] * 8
            )
            self.assertEqual(
                [call["denoise"] for call in fixture.calls], [1, 0.4] * 4
            )
            self.assertEqual(
                [call["model"] for call in fixture.calls],
                ["adapted-model"] * 8,
            )
            self.assertEqual(fixture.guides, [])
            entries: list[dict[str, str]] = typing.cast(
                dict[str, list[dict[str, str]]], result["ui"]
            )["images"]
            stills: list[dict[str, str]] = [
                entry
                for entry in entries
                if entry["filename"].endswith(".png")
            ]
            self.assertEqual(len(stills), 4)
            observations: list[dict[str, object]] = []
            for index, entry in enumerate(stills):
                with PIL.Image.open(root / entry["filename"]) as picture:
                    metadata: dict[str, object] = json.loads(
                        picture.info["image_sequence"]
                    )
                    self.assertEqual(metadata["seed"], 42)
                    self.assertEqual(
                        metadata["generation_mode"], "independent"
                    )
                    self.assertEqual(metadata["refiner_steps"], 3)
                    self.assertEqual(metadata["refiner_denoise"], 0.4)
                    self.assertIn(shared, metadata["positive"])
                    self.assertIn(
                        "snowy park" if index < 2 else "warm library",
                        metadata["positive"],
                    )
                    observations.append(metadata)
            animations: list[dict[str, object]] = []
            for entry in entries:
                if not entry["filename"].endswith(".gif"):
                    continue
                with PIL.Image.open(root / entry["filename"]) as picture:
                    animations.append(
                        {
                            "frames": picture.n_frames,
                            "duration": sum(
                                frame.info["duration"]
                                for frame in PIL.ImageSequence.Iterator(
                                    picture
                                )
                            ),
                            "loop": picture.info["loop"],
                        }
                    )
            self.assertEqual(
                animations, [{"frames": 4, "duration": 2000, "loop": 0}]
            )
            (root / "composer-verification.json").write_text(
                json.dumps(
                    {
                        "images": observations,
                        "animations": animations,
                        "native_calls": fixture.calls,
                    },
                    indent=2,
                )
                + "\n"
            )

    def test_composer_preflight_rejects_invalid_inputs(self) -> None:
        """Reject context, LoRA and sampling errors before native execution."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            modules: dict[str, object] = fixture.modules()
            typing.cast(
                types.SimpleNamespace, modules["folder_paths"]
            ).get_filename_list = lambda kind: ["character.safetensors"]
            with mock.patch.dict("sys.modules", modules):
                with self.assertRaises(ValueError):
                    image_sequence.ScenePrompts().compose(
                        "character",
                        "coat",
                        "illustration",
                        "park\nlibrary\nbeach",
                        "standing\n---\nsitting",
                    )
                with self.assertRaises(ValueError):
                    image_sequence.ScenePrompts().compose(
                        "", "coat", "illustration", "park", "standing"
                    )
                with self.assertRaises(ValueError):
                    image_sequence.OptionalLora().apply(
                        object(), object(), "missing.safetensors", 1, 1
                    )
                with self.assertRaises(ValueError):
                    image_sequence.OptionalLora().apply(
                        object(), object(), "(none)", float("nan"), 1
                    )
                for mode, seed_mode, refiner_steps, refiner_denoise in [
                    ("invalid", "fixed", 2, 0.4),
                    ("independent", "invalid", 2, 0.4),
                    ("independent", "fixed", -1, 0.4),
                    ("independent", "fixed", 2, float("nan")),
                ]:
                    with self.assertRaises(ValueError):
                        image_sequence.ImageSequence().generate(
                            model=object(),
                            clip=object(),
                            vae=object(),
                            shared_prompt="character",
                            actions="standing",
                            negative_prompt="",
                            change_strength=0.45,
                            width=64,
                            height=64,
                            seed=42,
                            steps=2,
                            cfg=5.5,
                            sampler_name="euler",
                            scheduler="normal",
                            start_image=1,
                            image_count=0,
                            filename_prefix="invalid",
                            generation_mode=mode,
                            seed_mode=seed_mode,
                            refiner_steps=refiner_steps,
                            refiner_denoise=refiner_denoise,
                        )
            self.assertEqual(fixture.calls, [])
            self.assertEqual(list(root.iterdir()), [])

    def test_composer_partial_run_without_continuation(self) -> None:
        """Regenerate an independent action with stable context and fixed seed."""
        with tempfile.TemporaryDirectory() as temporary:
            fixture: _NativeFixture = _NativeFixture(pathlib.Path(temporary))
            with mock.patch.dict("sys.modules", fixture.modules()):
                shared, actions = image_sequence.ScenePrompts().compose(
                    "character",
                    "coat",
                    "illustration",
                    "park",
                    "standing\nwaving\n---\nsitting\nreading",
                )
                image_sequence.ImageSequence().generate(
                    model=object(),
                    clip=object(),
                    vae=object(),
                    shared_prompt=shared,
                    actions=actions,
                    negative_prompt="",
                    change_strength=0.45,
                    width=64,
                    height=64,
                    seed=42,
                    steps=2,
                    cfg=5.5,
                    sampler_name="euler",
                    scheduler="normal",
                    start_image=4,
                    image_count=1,
                    filename_prefix="partial-composer",
                    generation_mode="independent",
                    seed_mode="fixed",
                )
            self.assertEqual(len(fixture.calls), 1)
            self.assertEqual(fixture.calls[0]["seed"], 42)
            self.assertEqual(fixture.calls[0]["source"], 0)
            self.assertIn("park, reading", fixture.calls[0]["positive"])

    def test_cancel_between_refinement_passes_preserves_completed_images(
        self,
    ) -> None:
        """Native cancellation prevents refinement and retains earlier PNGs."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            modules: dict[str, object] = fixture.modules()

            def cancel_after_first_pass() -> None:
                """Interrupt the second image before its refinement starts."""
                if len(fixture.calls) == 3:
                    raise InterruptedError("Interrupted before refinement")

            typing.cast(
                types.SimpleNamespace, modules["comfy.model_management"]
            ).throw_exception_if_processing_interrupted = (
                cancel_after_first_pass
            )
            with mock.patch.dict("sys.modules", modules):
                with self.assertRaises(InterruptedError):
                    image_sequence.ImageSequence().generate(
                        model=object(),
                        clip=object(),
                        vae=object(),
                        shared_prompt="character",
                        actions="standing\nwaving",
                        negative_prompt="",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=42,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="cancel-refiner",
                        generation_mode="independent",
                        seed_mode="fixed",
                        refiner_steps=2,
                        refiner_denoise=0.4,
                    )
            self.assertEqual(len(fixture.calls), 3)
            self.assertEqual(len(list(root.glob("*.png"))), 1)
            self.assertEqual(list(root.glob("*.gif")), [])

    def test_native_cancel_terminates_interpolation_process_tree(self) -> None:
        """Cancel real descendants while preserving completed keyframes and cleanup."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            marker: pathlib.Path = root / "processes.json"
            ready: pathlib.Path = root / "descendant-ready"
            engine: pathlib.Path = root / "interpolation-probe.py"
            interpreter: pathlib.Path = root / "python"
            interpreter.symlink_to(sys.executable)
            engine.write_text(
                f"#!{interpreter}\n"
                '"""Run a bounded process tree to exercise native cancellation."""\n'
                "import json, os, pathlib, signal, subprocess, sys, time\n"
                f"ready = pathlib.Path({str(ready)!r})\n"
                "child = subprocess.Popen([sys.executable, '-c', "
                f"{('import pathlib, signal, time; signal.signal(signal.SIGTERM, signal.SIG_IGN); pathlib.Path(' + repr(str(ready)) + ').touch(); time.sleep(5)')!r}], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)\n"
                "while not ready.exists(): time.sleep(0.01)\n"
                f"pathlib.Path({str(marker)!r}).write_text(json.dumps([os.getpid(), child.pid]))\n"
                "time.sleep(1.5)\n"
                "sys.exit(7)\n"
            )
            engine.chmod(0o700)
            fixture: _NativeFixture = _NativeFixture(root)
            modules: dict[str, object] = fixture.modules()

            def cancel_when_running() -> None:
                """Request native cancellation only after both processes are alive."""
                if marker.exists():
                    raise InterruptedError("Native cancellation during RIFE")

            typing.cast(
                types.SimpleNamespace, modules["comfy.model_management"]
            ).throw_exception_if_processing_interrupted = cancel_when_running
            started: float = time.monotonic()
            try:
                with (
                    mock.patch.dict("sys.modules", modules),
                    mock.patch("shutil.which", return_value=str(engine)),
                ):
                    with self.assertRaises(InterruptedError):
                        image_sequence.ImageSequence().generate(
                            model="anchored",
                            clip=object(),
                            vae=object(),
                            actions="standing\nsitting\n---\nreading",
                            shared_prompt="robot",
                            negative_prompt="",
                            change_strength=0.45,
                            width=64,
                            height=64,
                            seed=42,
                            steps=2,
                            cfg=5.5,
                            sampler_name="euler",
                            scheduler="normal",
                            start_image=1,
                            image_count=0,
                            filename_prefix="cancelled",
                            interpolation_multiplier=2,
                        )
                self.assertLess(time.monotonic() - started, 3)
                self.assertEqual(len(fixture.calls), 3)
                self.assertEqual(len(list(root.glob("*.png"))), 3)
                self.assertEqual(list(root.glob("*.gif")), [])
                self.assertEqual(list(root.glob(".image-sequence-*")), [])
                for pid in json.loads(marker.read_text()):
                    status: pathlib.Path = pathlib.Path(f"/proc/{pid}/stat")
                    deadline: float = time.monotonic() + 0.5

                    def running(status: pathlib.Path = status) -> bool:
                        """Handle process exit during a procfs status read."""
                        try:
                            return status.read_text().split()[2] not in (
                                "Z",
                                "X",
                            )
                        except (FileNotFoundError, ProcessLookupError):
                            return False

                    while running() and time.monotonic() < deadline:
                        time.sleep(0.01)
                    self.assertFalse(running())
            finally:
                if marker.exists():
                    for pid in json.loads(marker.read_text()):
                        try:
                            os.kill(pid, signal.SIGKILL)
                        except ProcessLookupError:
                            pass

    def test_interpolation_keeps_scene_cuts_and_playback_time(self) -> None:
        """Decode interpolated GIFs and ensure no synthesized cross-scene frames."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            calls: list[list[str]] = []

            def interpolate(
                arguments: list[str], environment: dict[str, str]
            ) -> None:
                """Substitute the RIFE boundary with observable intermediate pixels."""
                calls.append(arguments)
                source: pathlib.Path = pathlib.Path(
                    arguments[arguments.index("-i") + 1]
                )
                target: pathlib.Path = pathlib.Path(
                    arguments[arguments.index("-o") + 1]
                )
                with PIL.Image.open(source) as picture:
                    colors: list[int] = [
                        frame.convert("RGB").getpixel((0, 0))[0]
                        for frame in PIL.ImageSequence.Iterator(picture)
                    ]
                self.assertIn(colors, ([25, 51], [76, 102]))
                pictures: list[PIL.Image.Image] = [
                    PIL.Image.new("RGB", (64, 64), (value, value, value))
                    for value in [
                        colors[0],
                        (colors[0] + colors[1]) // 2,
                        colors[1],
                    ]
                ]
                pictures[0].save(
                    target,
                    save_all=True,
                    append_images=pictures[1:],
                    duration=250,
                    loop=0,
                )

            with (
                mock.patch.dict("sys.modules", fixture.modules()),
                mock.patch("shutil.which", return_value="/fixture/video2x"),
                mock.patch(
                    "projects.comfyui_tools.nodes.image_sequence.image_sequence._run_interpolation",
                    side_effect=interpolate,
                ),
            ):
                image_sequence.ImageSequence().generate(
                    model="anchored",
                    clip=object(),
                    vae=object(),
                    actions="standing\nraising hand\n---\nsitting\nreading",
                    shared_prompt="robot",
                    negative_prompt="blur",
                    change_strength=0.45,
                    width=64,
                    height=64,
                    seed=42,
                    steps=2,
                    cfg=5.5,
                    sampler_name="euler",
                    scheduler="normal",
                    start_image=1,
                    image_count=0,
                    filename_prefix="interpolated",
                    interpolation_multiplier=2,
                )
            self.assertEqual(len(calls), 2)
            self.assertEqual(len(list(root.glob("*.gif"))), 1)
            for file in root.glob("*.gif"):
                with PIL.Image.open(file) as picture:
                    durations: list[int] = [
                        frame.info["duration"]
                        for frame in PIL.ImageSequence.Iterator(picture)
                    ]
                    self.assertEqual(
                        durations,
                        [250, 250, 500] * 2,
                    )
                    self.assertEqual(picture.info["loop"], 0)

    def test_canny_feedback_resets_and_can_be_disabled(self) -> None:
        """Guide only same-scene edits, including supplied partial continuation."""
        for strength, continuation, actions in (
            (0.35, False, "standing\nraising hand\n---\nsitting\nreading"),
            (0.0, False, "standing\nraising hand\n---\nsitting\nreading"),
            (0.35, True, "standing\nraising hand"),
        ):
            with self.subTest(strength=strength, continuation=continuation):
                with tempfile.TemporaryDirectory() as temporary:
                    fixture: _NativeFixture = _NativeFixture(
                        pathlib.Path(temporary)
                    )
                    with mock.patch.dict("sys.modules", fixture.modules()):
                        image_sequence.ImageSequence().generate(
                            model="anchored",
                            clip=object(),
                            vae=object(),
                            actions=actions,
                            shared_prompt="robot",
                            negative_prompt="blur",
                            change_strength=0.45,
                            width=64,
                            height=64,
                            seed=42,
                            steps=2,
                            cfg=5.5,
                            sampler_name="euler",
                            scheduler="normal",
                            start_image=2 if continuation else 1,
                            image_count=0,
                            filename_prefix="controlled",
                            control_net="canny",
                            control_strength=strength,
                            control_end=0.8,
                            continuation_image=torch.full((1, 64, 64, 3), 0.6)
                            if continuation
                            else None,
                        )
                    expected: list[float] = (
                        [0.6]
                        if continuation
                        else [0.1, 0.3]
                        if strength
                        else []
                    )
                    self.assertEqual(len(fixture.guides), len(expected))
                    for guide, source in zip(
                        fixture.guides, expected, strict=True
                    ):
                        self.assertAlmostEqual(
                            guide["source"], source, places=5
                        )
                        self.assertEqual(guide["strength"], 0.35)
                        self.assertEqual(guide["end"], 0.8)
                    for number, call in enumerate(fixture.calls):
                        self.assertEqual(
                            isinstance(call["positive"], dict),
                            bool(strength and (continuation or number % 2)),
                        )
                    with PIL.Image.open(
                        sorted(pathlib.Path(temporary).glob("*.png"))[-1]
                    ) as picture:
                        metadata: dict[str, object] = json.loads(
                            picture.info["image_sequence"]
                        )
                        self.assertEqual(
                            metadata["control_strength"], strength
                        )
        for patch in (
            {"control_strength": -1},
            {"control_strength": float("nan")},
            {"control_end": 0},
            {"control_end": 1.1},
        ):
            with self.subTest(patch=patch):
                with self.assertRaises(ValueError):
                    image_sequence.ImageSequence().generate(
                        model=object(),
                        clip=object(),
                        vae=object(),
                        actions="standing",
                        shared_prompt="robot",
                        negative_prompt="",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=1,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="invalid",
                        **patch,
                    )

    def test_scene_feedback_outputs_and_continuation(self) -> None:
        """Generate two scenes, inspect PNGs, and continue with the absolute seed."""
        directory: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
        with tempfile.TemporaryDirectory(dir=directory or None) as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            with mock.patch.dict("sys.modules", fixture.modules()):
                result: dict[str, object] = (
                    image_sequence.ImageSequence().generate(
                        model="anchored",
                        clip=object(),
                        vae=object(),
                        actions="standing, park\nsitting, park\n---\nstanding, library\nreading, library",
                        shared_prompt="blue robot",
                        negative_prompt="blur",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=42,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler_ancestral",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="sequence",
                        prompt={"workflow": "fixture"},
                        extra_pnginfo={"workflow": {"version": 0.4}},
                    )
                )
            self.assertEqual(
                [call["seed"] for call in fixture.calls], [42, 43, 44, 45]
            )
            self.assertEqual(
                [call["denoise"] for call in fixture.calls],
                [1.0, 0.45, 1.0, 0.45],
            )
            for actual, expected in zip(
                [call["source"] for call in fixture.calls],
                [0.0, 0.1, 0.0, 0.3],
                strict=True,
            ):
                self.assertAlmostEqual(actual, expected, places=5)
            self.assertTrue(
                all(call["model"] == "anchored" for call in fixture.calls)
            )
            files: list[pathlib.Path] = sorted(root.glob("*.png"))
            self.assertEqual(len(files), 4)
            animations: list[pathlib.Path] = sorted(root.glob("*.gif"))
            self.assertEqual(len(animations), 1)
            gif_receipt: list[dict[str, object]] = []
            for animation in animations:
                with PIL.Image.open(animation) as playback:
                    colors: list[int] = []
                    durations: list[int] = []
                    for frame in PIL.ImageSequence.Iterator(playback):
                        self.assertEqual(frame.size, (64, 64))
                        colors.append(frame.convert("RGB").getpixel((0, 0))[0])
                        durations.append(frame.info["duration"])
                    self.assertEqual(playback.info["loop"], 0)
                expected: list[int] = [25, 51, 76, 102]
                self.assertEqual(colors, expected)
                self.assertEqual(durations, [500] * len(expected))
                gif_receipt.append(
                    {
                        "file": animation.name,
                        "colors": colors,
                        "durations_ms": durations,
                    }
                )
            output_images: list[dict[str, str]] = typing.cast(
                dict[str, list[dict[str, str]]], result["ui"]
            )["images"]
            self.assertEqual(len(output_images), 5)
            with PIL.Image.open(files[1]) as picture:
                metadata: dict[str, object] = json.loads(
                    picture.info["image_sequence"]
                )
                self.assertEqual(metadata["seed"], 43)
                self.assertEqual(metadata["scene"], 1)
                self.assertEqual(metadata["action"], 2)
                self.assertIn("blue robot", metadata["positive"])
                self.assertEqual(
                    json.loads(picture.info["workflow"]), {"version": 0.4}
                )
            continuation: torch.Tensor = torch.full((1, 64, 64, 3), 0.2)
            with mock.patch.dict("sys.modules", fixture.modules()):
                image_sequence.ImageSequence().generate(
                    model="anchored",
                    clip=object(),
                    vae=object(),
                    actions="standing\nsitting\nreading",
                    shared_prompt="blue robot",
                    negative_prompt="blur",
                    change_strength=0.45,
                    width=64,
                    height=64,
                    seed=42,
                    steps=2,
                    cfg=5.5,
                    sampler_name="euler_ancestral",
                    scheduler="normal",
                    start_image=3,
                    image_count=1,
                    filename_prefix="continued",
                    continuation_image=continuation,
                )
            self.assertEqual(fixture.calls[-1]["seed"], 44)
            self.assertAlmostEqual(fixture.calls[-1]["source"], 0.2, places=5)
            self.assertAlmostEqual(continuation.mean().item(), 0.2, places=5)
            for animation in root.glob("continued*.gif"):
                with PIL.Image.open(animation) as playback:
                    self.assertEqual(playback.n_frames, 1)
                    self.assertEqual(playback.info["duration"], 500)
            if directory:
                receipt: pathlib.Path = pathlib.Path(
                    directory, "image-sequence-result.json"
                )
                receipt.write_text(
                    json.dumps(
                        {
                            "calls": fixture.calls,
                            "outputs": result["ui"],
                            "native_gpu_inference": False,
                            "gif_playback": gif_receipt,
                        },
                        indent=2,
                    )
                )
                for file in [*root.glob("*.png"), *root.glob("*.gif")]:
                    pathlib.Path(directory, file.name).write_bytes(
                        file.read_bytes()
                    )

    def test_invalid_input_before_sampling(self) -> None:
        """Reject missing continuity, empty scenes and invalid controls before work."""
        base: dict[str, object] = dict(
            model=object(),
            clip=object(),
            vae=object(),
            actions="standing\nsitting",
            shared_prompt="robot",
            negative_prompt="",
            change_strength=0.45,
            width=64,
            height=64,
            seed=1,
            steps=2,
            cfg=5.5,
            sampler_name="euler",
            scheduler="normal",
            start_image=1,
            image_count=0,
            filename_prefix="sequence",
        )
        for patch in [
            {"actions": ""},
            {"actions": "---\nfirst"},
            {"actions": "first\n---"},
            {"actions": "first\n---\n---\nlast"},
            {"start_image": 2},
            {"start_image": 3},
            {"start_image": 0},
            {"image_count": -1},
            {"change_strength": 0.0},
            {"change_strength": float("nan")},
            {"width": 65},
            {"filename_prefix": "../unsafe"},
            {"continuation_image": torch.zeros(2, 64, 64, 3)},
            {"gif_frame_duration_ms": 0},
            {"gif_frame_duration_ms": 25},
            {"gif_frame_duration_ms": 10001},
        ]:
            with self.subTest(patch=patch):
                with self.assertRaises(ValueError):
                    image_sequence.ImageSequence().generate(**(base | patch))

    def test_interruption_keeps_completed_images(self) -> None:
        """Cancel before further sampling while preserving already saved PNGs."""
        with tempfile.TemporaryDirectory() as temporary:
            fixture: _NativeFixture = _NativeFixture(pathlib.Path(temporary))
            fixture.interrupt_after = 1
            with mock.patch.dict("sys.modules", fixture.modules()):
                with self.assertRaises(InterruptedError):
                    image_sequence.ImageSequence().generate(
                        model=object(),
                        clip=object(),
                        vae=object(),
                        actions="standing\nsitting\nreading",
                        shared_prompt="robot",
                        negative_prompt="",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=1,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="sequence",
                    )
            self.assertEqual(len(fixture.calls), 1)
            self.assertEqual(
                len(list(pathlib.Path(temporary).glob("*.png"))), 1
            )

    def test_interrupted_sequence_preserves_pngs_without_gif(self) -> None:
        """Retain completed PNGs without publishing a partial sequence GIF."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            fixture.interrupt_after = 2
            with mock.patch.dict("sys.modules", fixture.modules()):
                with self.assertRaises(InterruptedError):
                    image_sequence.ImageSequence().generate(
                        model=object(),
                        clip=object(),
                        vae=object(),
                        actions="standing\nsitting\n---\nreading",
                        shared_prompt="robot",
                        negative_prompt="",
                        change_strength=0.45,
                        width=64,
                        height=64,
                        seed=1,
                        steps=2,
                        cfg=5.5,
                        sampler_name="euler",
                        scheduler="normal",
                        start_image=1,
                        image_count=0,
                        filename_prefix="sequence",
                        gif_frame_duration_ms=200,
                    )
            self.assertEqual(len(fixture.calls), 2)
            self.assertEqual(len(list(root.glob("*.png"))), 2)
            animations: list[pathlib.Path] = list(root.glob("*.gif"))
            self.assertEqual(animations, [])
            self.assertEqual(len(list(root.glob("*.png"))), 2)

    def test_repeated_images_preserve_playback_time(self) -> None:
        """Allow GIF frame coalescing while retaining each generated image's hold."""
        with tempfile.TemporaryDirectory() as temporary:
            root: pathlib.Path = pathlib.Path(temporary)
            fixture: _NativeFixture = _NativeFixture(root)
            modules: dict[str, object] = fixture.modules()
            native: types.SimpleNamespace = typing.cast(
                types.SimpleNamespace, modules["nodes"]
            )
            native.common_ksampler = lambda *args, **kwargs: (
                {"samples": torch.full((1, 4, 8, 8), 0.25)},
            )
            with mock.patch.dict("sys.modules", modules):
                image_sequence.ImageSequence().generate(
                    model=object(),
                    clip=object(),
                    vae=object(),
                    actions="standing\nstanding",
                    shared_prompt="robot",
                    negative_prompt="",
                    change_strength=0.45,
                    width=64,
                    height=64,
                    seed=1,
                    steps=2,
                    cfg=5.5,
                    sampler_name="euler",
                    scheduler="normal",
                    start_image=1,
                    image_count=0,
                    filename_prefix="repeated",
                    gif_frame_duration_ms=300,
                )
            for animation in root.glob("*.gif"):
                with PIL.Image.open(animation) as playback:
                    self.assertEqual(
                        sum(
                            frame.info["duration"]
                            for frame in PIL.ImageSequence.Iterator(playback)
                        ),
                        600,
                    )

    def test_reference_views_remain_uncropped(self) -> None:
        """Prepare only selected views, retaining edge details and optional continuation."""
        with tempfile.TemporaryDirectory() as temporary:
            fixture: _NativeFixture = _NativeFixture(pathlib.Path(temporary))
            portrait: torch.Tensor = torch.ones(1, 64, 32, 3)
            portrait[:, :8] = 0.25
            portrait[:, -8:] = 0.75
            landscape: torch.Tensor = torch.full((1, 32, 64, 3), 0.5)
            sources: dict[str, torch.Tensor] = {
                "portrait.png": portrait,
                "landscape.png": landscape,
            }
            modules: dict[str, object] = fixture.modules()
            native: types.SimpleNamespace = typing.cast(
                types.SimpleNamespace, modules["nodes"]
            )
            native.LoadImage = lambda: types.SimpleNamespace(
                load_image=lambda name: (sources[name], None)
            )
            with mock.patch.dict("sys.modules", modules):
                references, continuation = (
                    image_sequence.ReferenceImages().prepare(
                        "portrait.png",
                        "landscape.png",
                        "(none)",
                        "portrait.png",
                    )
                )
            self.assertEqual(tuple(references.shape), (2, 512, 512, 3))
            self.assertAlmostEqual(
                references[0, 0, 256].mean().item(), 0.25, places=5
            )
            self.assertAlmostEqual(
                references[0, -1, 256].mean().item(), 0.75, places=5
            )
            self.assertTrue(torch.equal(continuation, portrait))
            self.assertEqual(tuple(portrait.shape), (1, 64, 32, 3))


if __name__ == "__main__":
    unittest.main()
