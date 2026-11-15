"""Generate reference-anchored still sequences through native ComfyUI operations.

Reference sequences edit the preceding image; composer images are independent.
Approved keyframes can also be assembled without another sampling pass.
Model loading and identity conditioning belong to the enclosing workflow. Native
ComfyUI owns sampling, cancellation and immediate per-image PNG serialization.
"""

import io
import math
import os
import pathlib
import re
import shutil
import signal
import subprocess
import sys
import tempfile
import time
import typing
import uuid

import PIL.Image
import PIL.ImageSequence
import torch


def _actions(text: str) -> list[tuple[int, int, str]]:
    """Parse nonempty action lines with explicit, nonempty scene boundaries."""
    result: list[tuple[int, int, str]] = []
    scene: int = 1
    position: int = 0
    for line in text.splitlines():
        value: str = line.strip()
        if not value:
            continue
        if value == "---":
            if position == 0:
                raise ValueError(
                    "Each scene must contain at least one action."
                )
            scene += 1
            position = 0
        else:
            position += 1
            result.append((scene, position, value))
    if not result or position == 0 or len(result) > 128:
        raise ValueError(
            "Enter 1 to 128 actions; use --- only between nonempty scenes."
        )
    return result


def _image(value: torch.Tensor) -> None:
    """Require a single RGB-compatible image before any native model operation."""
    if (
        value.ndim != 4
        or value.shape[0] != 1
        or value.shape[-1] not in (3, 4)
        or min(value.shape[1:3]) < 1
    ):
        raise ValueError("Supply one RGB or RGBA still image per input.")


def _run_interpolation(
    arguments: list[str], environment: dict[str, str]
) -> None:
    """Run an owned process group while honoring native interruption and deadline."""
    import comfy.model_management

    comfy.model_management.throw_exception_if_processing_interrupted()
    with subprocess.Popen(
        arguments,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        start_new_session=True,
        env=environment,
    ) as process:
        completed: bool = False
        deadline: float = time.monotonic() + 300
        output: bytes
        errors: bytes
        try:
            while True:
                comfy.model_management.throw_exception_if_processing_interrupted()
                remaining: float = deadline - time.monotonic()
                if remaining <= 0:
                    raise subprocess.TimeoutExpired(arguments, 300)
                try:
                    output, errors = process.communicate(
                        timeout=min(0.25, remaining)
                    )
                except subprocess.TimeoutExpired:
                    continue
                comfy.model_management.throw_exception_if_processing_interrupted()
                if process.returncode:
                    raise subprocess.CalledProcessError(
                        process.returncode, arguments, output, errors
                    )
                completed = True
                return
        finally:
            if not completed:
                try:
                    os.killpg(process.pid, signal.SIGTERM)
                except ProcessLookupError:
                    pass
                try:
                    process.communicate(timeout=1)
                except subprocess.TimeoutExpired:
                    pass
                # Descendants can close their pipes before ignoring termination.
                try:
                    os.killpg(process.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
                process.communicate()


def _gif_frames(
    images: list[dict[str, str]],
    duration: int,
    interpolation_multiplier: int = 1,
) -> tuple[list[PIL.Image.Image], list[int]]:
    """Prepare one scene's playback, synthesizing motion only inside that scene."""
    import folder_paths

    root: pathlib.Path = pathlib.Path(
        folder_paths.get_output_directory()
    ).resolve()
    pictures: list[PIL.Image.Image] = []
    durations: list[int] = []
    for image in images:
        source: pathlib.Path = (
            root / image["subfolder"] / image["filename"]
        ).resolve()
        if not source.is_relative_to(root):
            raise ValueError(
                "Sequence images must remain inside the output directory."
            )
        with PIL.Image.open(source) as picture:
            for frame in PIL.ImageSequence.Iterator(picture):
                pictures.append(frame.convert("RGB"))
                durations.append(frame.info.get("duration", duration))
    if interpolation_multiplier > 1 and len(pictures) > 1:
        with tempfile.TemporaryDirectory(
            prefix=".image-sequence-", dir=root
        ) as temporary:
            working: pathlib.Path = pathlib.Path(temporary)
            source_gif: pathlib.Path = working / "keyframes.gif"
            target_gif: pathlib.Path = working / "interpolated.gif"
            pictures[0].save(
                source_gif,
                format="GIF",
                save_all=True,
                append_images=pictures[1:],
                duration=durations,
                loop=0,
                disposal=2,
                optimize=False,
            )
            with PIL.Image.open(source_gif) as source:
                durations = [
                    frame.info["duration"]
                    for frame in PIL.ImageSequence.Iterator(source)
                ]
                if len(durations) == 1:
                    pictures = [source.convert("RGB")]
            if len(durations) > 1:
                engine: str = (
                    shutil.which("video2x") or "/opt/imagegen/Video2X.AppImage"
                )
                environment: dict[str, str] = dict(os.environ)
                environment.update(TMPDIR=temporary, XDG_CACHE_HOME=temporary)
                try:
                    _run_interpolation(
                        [
                            engine,
                            "--appimage-extract-and-run",
                            "-i",
                            str(source_gif),
                            "-o",
                            str(target_gif),
                            "-p",
                            "rife",
                            "--rife-model",
                            "rife-v4.25",
                            "-m",
                            str(interpolation_multiplier),
                            "-t",
                            "100",
                            "-d",
                            "0",
                            "-a",
                            "none",
                            "-c",
                            "gif",
                            "--pix-fmt",
                            "rgb8",
                            "--no-copy-streams",
                            "--no-progress",
                        ],
                        environment,
                    )
                except subprocess.CalledProcessError as error:
                    raise RuntimeError(
                        "RIFE interpolation failed: "
                        + error.stderr.decode(errors="replace")[-2000:]
                    ) from error
                with PIL.Image.open(target_gif) as interpolated:
                    pictures = [
                        frame.convert("RGB")
                        for frame in PIL.ImageSequence.Iterator(interpolated)
                    ]
                if len(pictures) != 1 + interpolation_multiplier * (
                    len(durations) - 1
                ):
                    raise ValueError(
                        "RIFE returned an unexpected number of frames."
                    )
                holds: list[int] = []
                for hold in durations[:-1]:
                    ticks, remainder = divmod(
                        hold // 10, interpolation_multiplier
                    )
                    holds.extend(
                        10 * (ticks + (index < remainder))
                        for index in range(interpolation_multiplier)
                    )
                durations = holds + durations[-1:]
    return pictures, durations


def _save_gif(
    images: list[dict[str, str]],
    prefix: str,
    duration: int,
    interpolation_multiplier: int = 1,
    scene_lengths: list[int] | None = None,
) -> dict[str, str]:
    """Publish one sequence GIF, retaining hard cuts between scene groups."""
    import comfy.model_management
    import folder_paths

    lengths: list[int] = (
        scene_lengths if scene_lengths is not None else [len(images)]
    )
    if any(length < 1 for length in lengths) or sum(lengths) != len(images):
        raise ValueError(
            "GIF scene groups must partition the generated images."
        )
    pictures: list[PIL.Image.Image] = []
    durations: list[int] = []
    start: int = 0
    for length in lengths:
        comfy.model_management.throw_exception_if_processing_interrupted()
        frames: list[PIL.Image.Image]
        holds: list[int]
        frames, holds = _gif_frames(
            images[start : start + length], duration, interpolation_multiplier
        )
        pictures.extend(frames)
        durations.extend(holds)
        start += length
    comfy.model_management.throw_exception_if_processing_interrupted()
    root: pathlib.Path = pathlib.Path(
        folder_paths.get_output_directory()
    ).resolve()
    filename: str = prefix + ".gif"
    path: pathlib.Path = root / filename
    with io.BytesIO() as encoded:
        pictures[0].save(
            encoded,
            format="GIF",
            save_all=True,
            append_images=pictures[1:],
            duration=durations,
            loop=0,
            disposal=2,
            optimize=False,
        )
        with path.open("xb") as output:
            try:
                output.write(encoded.getvalue())
            except BaseException:
                path.unlink(missing_ok=True)
                raise
    return {"filename": filename, "subfolder": "", "type": "output"}


class ReferenceImages:
    """Select up to three persistent views and an independent continuation image."""

    RETURN_TYPES: tuple[str, str] = ("IMAGE", "IMAGE")
    """Prepared square reference batch and optional original continuation image."""
    RETURN_NAMES: tuple[str, str] = ("references", "continuation_image")
    """Native graph labels for persistent anchors and editable starting pixels."""
    FUNCTION: str = "prepare"
    """Native entry point for loading the chosen image files."""
    CATEGORY: str = "imagegen/workflow"
    """Menu shared with the sequence sampler."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Reuse native input-file discovery and upload widgets for every image."""
        import nodes

        choices: list[str] = list(
            nodes.LoadImage.INPUT_TYPES()["required"]["image"][0]
        )
        return {
            "required": {
                "reference_1": (
                    choices,
                    {
                        "image_upload": True,
                        "tooltip": "Primary character/appearance reference; remains active across every scene.",
                    },
                ),
                "reference_2": (
                    ["(none)", *choices],
                    {
                        "image_upload": True,
                        "tooltip": "Optional additional view. Select (none) to omit it.",
                    },
                ),
                "reference_3": (
                    ["(none)", *choices],
                    {
                        "image_upload": True,
                        "tooltip": "Optional additional view or style reference; all selected references influence the whole image.",
                    },
                ),
                "continuation": (
                    ["(none)", *choices],
                    {
                        "image_upload": True,
                        "tooltip": "Accepted preceding image when resuming inside a scene. Does not replace the persistent references.",
                    },
                ),
            }
        }

    @classmethod
    def IS_CHANGED(cls, **inputs: str) -> str:
        """Include native content hashes so replacing an uploaded file invalidates cache."""
        import nodes

        return ":".join(
            str(nodes.LoadImage.IS_CHANGED(name))
            for name in inputs.values()
            if name != "(none)"
        )

    @classmethod
    def VALIDATE_INPUTS(cls, **inputs: str) -> bool | str:
        """Check uploaded-file existence using the native loader's validation contract."""
        import nodes

        for name in inputs.values():
            if name != "(none)":
                result: bool | str = nodes.LoadImage.VALIDATE_INPUTS(name)
                if result is not True:
                    return result
        return True

    def prepare(
        self,
        reference_1: str,
        reference_2: str = "(none)",
        reference_3: str = "(none)",
        continuation: str = "(none)",
    ) -> tuple[torch.Tensor, torch.Tensor | None]:
        """Pad reference views without cropping; retain continuation at full resolution."""
        import nodes

        views: list[torch.Tensor] = []
        for name in (reference_1, reference_2, reference_3):
            if name == "(none)":
                continue
            image: torch.Tensor = nodes.LoadImage().load_image(name)[0]
            _image(image)
            height: int = image.shape[1]
            width: int = image.shape[2]
            scale: float = 512 / max(height, width)
            resized: torch.Tensor = nodes.ImageScale().upscale(
                image[..., :3],
                "lanczos",
                max(1, round(width * scale)),
                max(1, round(height * scale)),
                "disabled",
            )[0]
            vertical: int = 512 - resized.shape[1]
            horizontal: int = 512 - resized.shape[2]
            views.append(
                torch.nn.functional.pad(
                    resized.movedim(-1, 1),
                    (
                        horizontal // 2,
                        horizontal - horizontal // 2,
                        vertical // 2,
                        vertical - vertical // 2,
                    ),
                    mode="replicate",
                ).movedim(1, -1)
            )
        if not views:
            raise ValueError("Select at least one appearance reference.")
        previous: torch.Tensor | None = None
        if continuation != "(none)":
            previous = nodes.LoadImage().load_image(continuation)[0]
            _image(previous)
        return torch.cat(views, dim=0), previous


class ScenePrompts:
    """Compose explicit fixed appearance fields with scene-local settings/actions."""

    RETURN_TYPES: tuple[str, str] = ("STRING", "STRING")
    """Shared appearance and the scene-aware action list for the sampler."""
    RETURN_NAMES: tuple[str, str] = ("shared_prompt", "actions")
    """Graph labels describing the two composed text outputs."""
    FUNCTION: str = "compose"
    """Native entry point for deterministic prompt composition."""
    CATEGORY: str = "imagegen/workflow"
    """Menu shared with the native sequence output node."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Expose character, clothing, style, scene settings and actions separately."""
        return {
            "required": {
                "character_prompt": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "1girl, adult woman, black hair, green eyes",
                    },
                ),
                "clothing_prompt": (
                    "STRING",
                    {"multiline": True, "default": "red coat, black boots"},
                ),
                "style_prompt": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "masterpiece, best quality, newest, anime illustration",
                    },
                ),
                "scene_settings": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "snowy park\nwarm library",
                        "tooltip": "One setting per nonempty line, matching the scenes. A single line applies to every scene.",
                    },
                ),
                "actions": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "standing\nwaving\n---\nsitting\nreading",
                        "tooltip": "One desired image per line; --- starts the next scene and setting.",
                    },
                ),
            }
        }

    def compose(
        self,
        character_prompt: str,
        clothing_prompt: str,
        style_prompt: str,
        scene_settings: str,
        actions: str,
    ) -> tuple[str, str]:
        """Resolve fixed fields once, rejecting ambiguous scene-setting counts."""
        instructions: list[tuple[int, int, str]] = _actions(actions)
        if not character_prompt.strip():
            raise ValueError("Enter a shared character description.")
        settings: list[str] = [
            line.strip()
            for line in scene_settings.splitlines()
            if line.strip()
        ]
        if len(settings) not in (1, instructions[-1][0]):
            raise ValueError(
                "Enter one shared setting or exactly one setting per scene."
            )
        shared: str = ", ".join(
            part.strip()
            for part in (style_prompt, character_prompt, clothing_prompt)
            if part.strip()
        )
        composed: list[str] = []
        for scene, position, action in instructions:
            if scene > 1 and position == 1:
                composed.append("---")
            setting: str = settings[0 if len(settings) == 1 else scene - 1]
            composed.append(f"{setting}, {action}")
        return shared, "\n".join(composed)


class OptionalLora:
    """Select an already-installed LoRA through ComfyUI's native model/CLIP loader."""

    RETURN_TYPES: tuple[str, str] = ("MODEL", "CLIP")
    """Native diffusion model and text encoder after optional adaptation."""
    FUNCTION: str = "apply"
    """Native entry point for applying the selected local weights."""
    CATEGORY: str = "imagegen/workflow"
    """Menu shared with the composer and sequence sampler."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """List local weights with a no-adaptation default; never download assets."""
        import folder_paths

        return {
            "required": {
                "model": ("MODEL",),
                "clip": ("CLIP",),
                "lora_name": (
                    ["(none)", *folder_paths.get_filename_list("loras")],
                ),
                "strength_model": (
                    "FLOAT",
                    {"default": 1.0, "min": -2.0, "max": 2.0, "step": 0.05},
                ),
                "strength_clip": (
                    "FLOAT",
                    {"default": 1.0, "min": -2.0, "max": 2.0, "step": 0.05},
                ),
            }
        }

    def apply(
        self,
        model: object,
        clip: object,
        lora_name: str,
        strength_model: float,
        strength_clip: float,
    ) -> tuple[object, object]:
        """Return the original pair when disabled, otherwise use installed weights."""
        import folder_paths
        import nodes

        if not math.isfinite(strength_model) or not math.isfinite(
            strength_clip
        ):
            raise ValueError("Use finite LoRA strengths.")
        if lora_name == "(none)":
            return model, clip
        if lora_name not in folder_paths.get_filename_list("loras"):
            raise ValueError("Select an installed LoRA or (none).")
        return nodes.LoraLoader().load_lora(
            model, clip, lora_name, strength_model, strength_clip
        )


class ImageSequence:
    """Sample scene-aware actions and save each completed still immediately."""

    RETURN_TYPES: tuple[str] = ("IMAGE",)
    """Completed images as a CPU batch for optional downstream consumers."""
    FUNCTION: str = "generate"
    """Native execution entry point."""
    OUTPUT_NODE: bool = True
    """Enable execution, PNG output and native App gallery rendering."""
    CATEGORY: str = "imagegen/workflow"
    """Menu shared with the reference selector."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Expose the sequence editor and routine controls without graph rewiring."""
        import comfy.samplers

        return {
            "required": {
                "model": ("MODEL",),
                "clip": ("CLIP",),
                "vae": ("VAE",),
                "shared_prompt": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "masterpiece, best quality, newest, highres, illustration, blue robot, red chest plate",
                        "tooltip": "Appearance and style tags used for every image. Put scene-specific locations in the actions.",
                    },
                ),
                "actions": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "standing, snowy park, beside a bench\nsitting, snowy park, on a bench\n---\nstanding, library, beside a bookshelf\nreading, library, holding an open book",
                        "tooltip": "One complete desired image per nonempty line. A line containing --- starts a scene. Generation mode selects independent images or previous-image edits.",
                    },
                ),
                "negative_prompt": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "worst quality, low quality, lowres, blurry, signature, username, logo, text, watermark",
                    },
                ),
                "change_strength": (
                    "FLOAT",
                    {
                        "default": 0.45,
                        "min": 0.01,
                        "max": 1.0,
                        "step": 0.01,
                        "tooltip": "Previous-image edits only: lower preserves more; higher allows more movement. Independent images and scene starts use fresh noise.",
                    },
                ),
                "width": (
                    "INT",
                    {"default": 1024, "min": 64, "max": 4096, "step": 8},
                ),
                "height": (
                    "INT",
                    {"default": 1024, "min": 64, "max": 4096, "step": 8},
                ),
                "seed": (
                    "INT",
                    {
                        "default": 42,
                        "min": 0,
                        "max": 0xFFFFFFFFFFFFFFFF,
                        "control_after_generate": True,
                    },
                ),
                "steps": ("INT", {"default": 30, "min": 1, "max": 1000}),
                "cfg": ("FLOAT", {"default": 5.5, "min": 0.0, "max": 100.0}),
                "sampler_name": (comfy.samplers.KSampler.SAMPLERS,),
                "scheduler": (comfy.samplers.KSampler.SCHEDULERS,),
                "start_image": (
                    "INT",
                    {
                        "default": 1,
                        "min": 1,
                        "max": 128,
                        "tooltip": "Absolute action number across scenes. Previous-image mode requires continuation inside a scene; independent mode does not.",
                    },
                ),
                "image_count": (
                    "INT",
                    {
                        "default": 0,
                        "min": 0,
                        "max": 128,
                        "tooltip": "0 generates all remaining actions; 1 previews or regenerates one action.",
                    },
                ),
                "filename_prefix": ("STRING", {"default": "sequence"}),
            },
            "optional": {
                "continuation_image": ("IMAGE",),
                "gif_frame_duration_ms": (
                    "INT",
                    {
                        "default": 500,
                        "min": 20,
                        "max": 10000,
                        "step": 10,
                        "tooltip": "Playback hold per generated image in scene and combined GIFs; no intermediate motion is generated.",
                    },
                ),
                "control_net": ("CONTROL_NET",),
                "control_strength": (
                    "FLOAT",
                    {
                        "default": 0.35,
                        "min": 0.0,
                        "max": 2.0,
                        "step": 0.05,
                        "tooltip": "Previous-image Canny guidance within a scene. Zero disables it; lower it for larger pose changes.",
                    },
                ),
                "control_end": (
                    "FLOAT",
                    {"default": 0.8, "min": 0.01, "max": 1.0, "step": 0.01},
                ),
                "interpolation_multiplier": (
                    [1, 2, 4, 8],
                    {
                        "default": 1,
                        "tooltip": "RIFE frames per transition. One disables smoothing; scene cuts and original PNGs are preserved.",
                    },
                ),
                "generation_mode": (
                    ["continue_previous", "independent"],
                    {
                        "default": "continue_previous",
                        "tooltip": "Continue previous pixels within a scene, or generate each image independently.",
                    },
                ),
                "seed_mode": (
                    ["increment", "fixed"],
                    {
                        "default": "increment",
                        "tooltip": "Increment by absolute action index, or reuse the same base seed for every image.",
                    },
                ),
                "refiner_steps": (
                    "INT",
                    {
                        "default": 0,
                        "min": 0,
                        "max": 1000,
                        "tooltip": "Optional second pass on each image's latent. Zero disables it.",
                    },
                ),
                "refiner_denoise": (
                    "FLOAT",
                    {"default": 0.4, "min": 0.01, "max": 1.0, "step": 0.01},
                ),
            },
            "hidden": {
                "prompt": "PROMPT",
                "extra_pnginfo": "EXTRA_PNGINFO",
                "unique_id": "UNIQUE_ID",
            },
        }

    @classmethod
    def IS_CHANGED(cls, **inputs: object) -> float:
        """Permit repeat runs with unchanged controls while retaining deterministic seeds."""
        return float("nan")

    def generate(
        self,
        model: object,
        clip: object,
        vae: object,
        shared_prompt: str,
        actions: str,
        negative_prompt: str,
        change_strength: float,
        width: int,
        height: int,
        seed: int,
        steps: int,
        cfg: float,
        sampler_name: str,
        scheduler: str,
        start_image: int,
        image_count: int,
        filename_prefix: str,
        continuation_image: torch.Tensor | None = None,
        prompt: object = None,
        extra_pnginfo: dict[str, object] | None = None,
        unique_id: str | None = None,
        gif_frame_duration_ms: int = 500,
        control_net: object = None,
        control_strength: float = 0.35,
        control_end: float = 0.8,
        interpolation_multiplier: int = 1,
        generation_mode: str = "continue_previous",
        seed_mode: str = "increment",
        refiner_steps: int = 0,
        refiner_denoise: float = 0.4,
    ) -> dict[str, object]:
        """Generate the selected action range, preserving native conditioning and outputs."""
        instructions: list[tuple[int, int, str]] = _actions(actions)
        if generation_mode not in (
            "continue_previous",
            "independent",
        ) or seed_mode not in ("increment", "fixed"):
            raise ValueError("Select a supported generation and seed mode.")
        if (
            not 0 <= refiner_steps <= 1000
            or not math.isfinite(refiner_denoise)
            or not 0 < refiner_denoise <= 1
        ):
            raise ValueError(
                "Use 0 to 1000 refinement steps and finite denoise above zero up to one."
            )
        if (
            interpolation_multiplier not in (1, 2, 4, 8)
            or gif_frame_duration_ms < interpolation_multiplier * 10
        ):
            raise ValueError(
                "Use interpolation 1, 2, 4 or 8 and at least 10 milliseconds per interpolated frame."
            )
        if (
            interpolation_multiplier > 1
            and not shutil.which("video2x")
            and not pathlib.Path("/opt/imagegen/Video2X.AppImage").is_file()
        ):
            raise ValueError(
                "Install the configured Video2X/RIFE engine or set interpolation to one."
            )
        if (
            not math.isfinite(control_strength)
            or not 0 <= control_strength <= 2
            or not math.isfinite(control_end)
            or not 0 < control_end <= 1
        ):
            raise ValueError(
                "Use finite control strength from 0 to 2 and end percentage above zero up to one."
            )
        if (
            not 20 <= gif_frame_duration_ms <= 10000
            or gif_frame_duration_ms % 10
        ):
            raise ValueError(
                "Use a GIF frame duration from 20 to 10000 milliseconds in multiples of 10."
            )
        if (
            not 1 <= start_image <= len(instructions)
            or not 0 <= image_count <= 128
        ):
            raise ValueError(
                "Select an existing start image and a count between 0 and 128."
            )
        if (
            not math.isfinite(change_strength)
            or not 0 < change_strength <= 1
            or not math.isfinite(cfg)
            or not 0 <= cfg <= 100
        ):
            raise ValueError(
                "Use finite change strength above zero up to one and CFG between 0 and 100."
            )
        if (
            any(
                size < 64 or size > 4096 or size % 8
                for size in (width, height)
            )
            or not 1 <= steps <= 1000
            or not 0 <= seed < 2**64
        ):
            raise ValueError(
                "Use dimensions divisible by 8 from 64 to 4096, 1 to 1000 steps, and a valid unsigned seed."
            )
        if not re.fullmatch(
            r"[A-Za-z0-9][A-Za-z0-9_-]{0,63}", filename_prefix
        ):
            raise ValueError(
                "Use a filename prefix of 1 to 64 letters, digits, underscores or hyphens."
            )
        if continuation_image is not None:
            _image(continuation_image)
        if (
            instructions[start_image - 1][1] != 1
            and continuation_image is None
            and generation_mode == "continue_previous"
        ):
            raise ValueError(
                "Starting inside a scene requires the accepted preceding image as continuation."
            )

        import comfy.model_management
        import nodes

        previous: torch.Tensor | None = continuation_image
        stop: int = (
            min(len(instructions), start_image - 1 + image_count)
            if image_count
            else len(instructions)
        )
        negative: object = nodes.CLIPTextEncode().encode(
            clip, negative_prompt
        )[0]
        run: str = uuid.uuid4().hex
        frames: list[torch.Tensor] = []
        saved: list[dict[str, str]] = []
        scene_lengths: list[int] = []
        scene_start: int = 0
        for index in range(start_image - 1, stop):
            comfy.model_management.throw_exception_if_processing_interrupted()
            scene: int
            position: int
            action: str
            scene, position, action = instructions[index]
            if unique_id and "server" in sys.modules:
                import server

                server.PromptServer.instance.send_progress_text(
                    f"Scene {scene}, image {index + 1} of {len(instructions)}",
                    unique_id,
                    server.PromptServer.instance.client_id,
                )
            positive_text: str = ", ".join(
                part for part in (shared_prompt.strip(), action) if part
            )
            positive: object = nodes.CLIPTextEncode().encode(
                clip, positive_text
            )[0]
            fresh: bool = position == 1 or generation_mode == "independent"
            denoise: float = 1.0 if fresh else change_strength
            controlled_negative: object = negative
            effective_control: float = 0.0
            latent: dict[str, torch.Tensor]
            if fresh:
                latent = nodes.EmptyLatentImage().generate(width, height, 1)[0]
            else:
                if previous is None:
                    raise ValueError(
                        "The preceding image is unavailable for this action."
                    )
                resized: torch.Tensor = nodes.ImageScale().upscale(
                    previous[..., :3], "lanczos", width, height, "center"
                )[0]
                latent = nodes.VAEEncode().encode(vae, resized)[0]
                if control_net is not None and control_strength:
                    import comfy_extras.nodes_canny

                    edges: torch.Tensor = (
                        comfy_extras.nodes_canny.Canny.execute(
                            resized, 0.1, 0.3
                        )[0]
                    )
                    positive, controlled_negative = (
                        nodes.ControlNetApplyAdvanced().apply_controlnet(
                            positive,
                            negative,
                            control_net,
                            edges,
                            control_strength,
                            0.0,
                            control_end,
                            vae=vae,
                        )
                    )
                    effective_control = control_strength
            image_seed: int = (
                seed + (index if seed_mode == "increment" else 0)
            ) % (2**64)
            sampled: dict[str, torch.Tensor] = nodes.common_ksampler(
                model,
                image_seed,
                steps,
                cfg,
                sampler_name,
                scheduler,
                positive,
                controlled_negative,
                latent,
                denoise=denoise,
            )[0]
            if refiner_steps:
                comfy.model_management.throw_exception_if_processing_interrupted()
                sampled = nodes.common_ksampler(
                    model,
                    image_seed,
                    refiner_steps,
                    cfg,
                    sampler_name,
                    scheduler,
                    positive,
                    controlled_negative,
                    sampled,
                    denoise=refiner_denoise,
                )[0]
            previous = (
                nodes.VAEDecode()
                .decode(vae, sampled)[0]
                .detach()
                .cpu()
                .float()
            )
            _image(previous)
            metadata: dict[str, object] = dict(extra_pnginfo or {})
            metadata["image_sequence"] = {
                "run": run,
                "scene": scene,
                "scene_image": position,
                "action": index + 1,
                "seed": image_seed,
                "positive": positive_text,
                "negative": negative_prompt,
                "denoise": denoise,
                "steps": steps,
                "cfg": cfg,
                "sampler": sampler_name,
                "scheduler": scheduler,
                "width": width,
                "height": height,
                "gif_frame_duration_ms": gif_frame_duration_ms,
                "control_strength": effective_control,
                "control_end": control_end,
                "interpolation_multiplier": interpolation_multiplier,
                "generation_mode": generation_mode,
                "seed_mode": seed_mode,
                "refiner_steps": refiner_steps,
                "refiner_denoise": refiner_denoise,
            }
            output: dict[str, object] = nodes.SaveImage().save_images(
                previous,
                f"{filename_prefix}_{run}_scene{scene:03d}_image{index + 1:03d}",
                prompt=prompt,
                extra_pnginfo=metadata,
            )
            descriptors: list[dict[str, str]] = typing.cast(
                dict[str, list[dict[str, str]]], output["ui"]
            )["images"]
            saved.extend(descriptors)
            frames.append(previous)
            if index + 1 == stop or instructions[index + 1][0] != scene:
                scene_lengths.append(len(saved) - scene_start)
                scene_start = len(saved)
        animation: dict[str, str] = _save_gif(
            saved,
            f"{filename_prefix}_{run}_sequence",
            gif_frame_duration_ms,
            interpolation_multiplier,
            scene_lengths,
        )
        return {
            "ui": {"images": [*saved, animation], "animated": (True,)},
            "result": (torch.cat(frames, dim=0),),
        }


class KeyframeGIF:
    """Publish approved image batches with optional reversible RIFE playback."""

    RETURN_TYPES: tuple[str] = ("IMAGE",)
    """Retain the approved images for downstream native consumers."""

    FUNCTION: str = "save"
    """Native framework entry point for model-free GIF assembly."""

    OUTPUT_NODE: bool = True
    """Publish the animation through ComfyUI output metadata."""

    CATEGORY: str = "imagegen/workflow"
    """Shared first-party workflow menu."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Expose playback controls independently of image generation."""
        return {
            "required": {
                "images": ("IMAGE",),
                "frame_duration_ms": (
                    "INT",
                    {"default": 320, "min": 20, "max": 10000, "step": 10},
                ),
                "interpolation_multiplier": ([1, 2, 4, 8], {"default": 4}),
                "ping_pong": ("BOOLEAN", {"default": True}),
                "filename_prefix": ("STRING", {"default": "keyframes"}),
            },
            "hidden": {"prompt": "PROMPT", "extra_pnginfo": "EXTRA_PNGINFO"},
        }

    def save(
        self,
        images: torch.Tensor,
        frame_duration_ms: int,
        interpolation_multiplier: int,
        ping_pong: bool,
        filename_prefix: str,
        prompt: object = None,
        extra_pnginfo: dict[str, object] | None = None,
    ) -> dict[str, object]:
        """Validate before writing, retain source PNGs and publish one closed loop."""
        if images.ndim != 4 or not 2 <= images.shape[0] <= 128:
            raise ValueError(
                "Supply 2 to 128 approved keyframes in one image batch."
            )
        _image(images[:1])
        if not re.fullmatch(
            r"[A-Za-z0-9][A-Za-z0-9_-]{0,127}", filename_prefix
        ):
            raise ValueError(
                "Use a filename prefix containing letters, digits, underscores and hyphens."
            )
        if (
            interpolation_multiplier not in (1, 2, 4, 8)
            or not 20 <= frame_duration_ms <= 10000
            or frame_duration_ms % 10
            or frame_duration_ms < interpolation_multiplier * 10
        ):
            raise ValueError(
                "Use interpolation 1, 2, 4 or 8 and a multiple of 10 milliseconds with at least 10 milliseconds per interpolated frame."
            )
        if (
            interpolation_multiplier > 1
            and not shutil.which("video2x")
            and not pathlib.Path("/opt/imagegen/Video2X.AppImage").is_file()
        ):
            raise ValueError(
                "Install the configured Video2X/RIFE engine or set interpolation to one."
            )
        import comfy.model_management
        import nodes

        run: str = uuid.uuid4().hex[:12]
        saved: list[dict[str, str]] = []
        for index in range(images.shape[0]):
            comfy.model_management.throw_exception_if_processing_interrupted()
            output: dict[str, object] = nodes.SaveImage().save_images(
                images[index : index + 1].detach().cpu().float(),
                f"{filename_prefix}_{run}_keyframe{index + 1:03d}",
                prompt=prompt,
                extra_pnginfo=extra_pnginfo,
            )
            descriptors: list[dict[str, str]] = typing.cast(
                dict[str, list[dict[str, str]]], output["ui"]
            )["images"]
            saved.extend(descriptors)
        playback: list[dict[str, str]] = (
            saved + list(reversed(saved[:-1])) if ping_pong else saved
        )
        animation: dict[str, str] = _save_gif(
            playback,
            f"{filename_prefix}_{run}_sequence",
            frame_duration_ms,
            interpolation_multiplier,
        )
        return {
            "ui": {"images": [*saved, animation], "animated": (True,)},
            "result": (images,),
        }


NODE_CLASS_MAPPINGS: dict[
    str,
    type[ReferenceImages]
    | type[ImageSequence]
    | type[ScenePrompts]
    | type[OptionalLora]
    | type[KeyframeGIF],
] = {
    "ImagegenReferenceImages": ReferenceImages,
    "ImagegenImageSequence": ImageSequence,
    "ImagegenScenePrompts": ScenePrompts,
    "ImagegenOptionalLora": OptionalLora,
    "ImagegenKeyframeGIF": KeyframeGIF,
}
"""Native registry for references, prompt composition, LoRA loading and sampling."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenReferenceImages": "Sequence References",
    "ImagegenImageSequence": "Image Sequence",
    "ImagegenScenePrompts": "Scene Prompts",
    "ImagegenOptionalLora": "Optional LoRA",
    "ImagegenKeyframeGIF": "Approved Keyframes to GIF",
}
"""Readable native node names in the graph and App builder."""
