"""Chain native img2img latent samples and save a looping GIF with installed Pillow.

ComfyUI owns model loading, conditioning, sampling and interruption. This output
node only sequences those operations and serializes the generated RGB frames.
"""

import io
import pathlib
import re
import sys
import typing

import numpy
import PIL.GifImagePlugin
import PIL.Image
import torch


class _VAE(typing.Protocol):
    """The native decoder interface consumed without owning the model."""

    def decode(self, samples: torch.Tensor) -> torch.Tensor:
        """Return an NHWC image batch from a native latent tensor."""
        ...


class _UIResult(typing.TypedDict):
    """Native saved-image preview metadata, including animation playback."""

    images: list[dict[str, str]]
    """Saved GIF descriptors consumed by the image viewer and history."""

    animated: tuple[bool]
    """Request native animated preview rather than a static canvas thumbnail."""


class _Result(typing.TypedDict):
    """Native output contract for the GIF preview and decoded frame batch."""

    ui: _UIResult
    """Output-file descriptors consumed by ComfyUI's image viewer/history."""

    result: tuple[torch.Tensor]
    """All generated frames as one native IMAGE batch."""


class ChainedImageToImageGIF:
    """Sample ordered frame poses with shared model and noise/reference controls."""

    INPUT_IS_LIST: bool = True
    """Consume the complete conditioning sequence instead of independent mapping."""

    RETURN_TYPES: tuple[str] = ("IMAGE",)
    """Allow subsequent native consumers to use the generated frame batch."""

    FUNCTION: str = "generate"
    """Framework execution entry point."""

    OUTPUT_NODE: bool = True
    """Write and preview the GIF when the graph is queued."""

    CATEGORY: str = "imagegen/workflow"
    """Menu shared with the first-party Prompt Sequence node."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Use installed native sampler choices and ordinary model/latent inputs."""
        import comfy.samplers

        return {
            "required": {
                "model": ("MODEL",),
                "positive": ("CONDITIONING",),
                "negative": ("CONDITIONING",),
                "latent_image": ("LATENT",),
                "vae": ("VAE",),
                "seed": (
                    "INT",
                    {
                        "default": 0,
                        "min": 0,
                        "max": 0xFFFFFFFFFFFFFFFF,
                        "control_after_generate": True,
                    },
                ),
                "steps": ("INT", {"default": 40, "min": 1, "max": 10000}),
                "cfg": ("FLOAT", {"default": 3.0, "min": 0.0, "max": 100.0}),
                "sampler_name": (comfy.samplers.KSampler.SAMPLERS,),
                "scheduler": (comfy.samplers.KSampler.SCHEDULERS,),
                "denoise": (
                    "FLOAT",
                    {"default": 0.45, "min": 0.0, "max": 1.0},
                ),
                "fps": ("FLOAT", {"default": 6.0, "min": 0.1, "max": 60.0}),
                "filename_prefix": ("STRING", {"default": "animation"}),
            },
            "optional": {
                "frame_source": (
                    ["previous", "original"],
                    {"default": "previous"},
                ),
                "frame_seed_mode": (
                    ["increment", "fixed"],
                    {"default": "increment"},
                ),
            },
            "hidden": {"unique_id": "UNIQUE_ID"},
        }

    def generate(
        self,
        model: list[object],
        positive: list[object],
        negative: list[object],
        latent_image: list[dict[str, torch.Tensor]],
        vae: list[_VAE],
        seed: list[int],
        steps: list[int],
        cfg: list[float],
        sampler_name: list[str],
        scheduler: list[str],
        denoise: list[float],
        fps: list[float],
        filename_prefix: list[str],
        unique_id: list[str] | None = None,
        frame_source: list[str] | None = None,
        frame_seed_mode: list[str] | None = None,
        _frame_inputs: typing.Callable[
            [int, torch.Tensor | None],
            tuple[object, object, dict[str, torch.Tensor]],
        ]
        | None = None,
    ) -> _Result:
        """Sample each pose from the selected latent and publish one complete GIF."""
        if not positive or len(negative) not in (1, len(positive)):
            raise ValueError(
                "Supply frame prompts and either one shared negative or one per frame."
            )
        for values in (
            model,
            latent_image,
            vae,
            seed,
            steps,
            cfg,
            sampler_name,
            scheduler,
            denoise,
            fps,
        ):
            if len(values) != 1:
                raise ValueError(
                    "Model, latent and sampling settings must each have one shared value."
                )
        if not filename_prefix or not re.fullmatch(
            r"[A-Za-z0-9][A-Za-z0-9_-]{0,127}", filename_prefix[0]
        ):
            raise ValueError(
                "Use a filename prefix containing only letters, digits, underscores and hyphens."
            )
        if not 0 < fps[0] <= 60 or not 0 <= denoise[0] <= 1:
            raise ValueError(
                "Use a positive frame rate up to 60 and denoise between zero and one."
            )
        source: str = (
            "previous"
            if frame_source is None
            else (frame_source[0] if len(frame_source) == 1 else "")
        )
        seed_mode: str = (
            "increment"
            if frame_seed_mode is None
            else (frame_seed_mode[0] if len(frame_seed_mode) == 1 else "")
        )
        if source not in ("previous", "original") or seed_mode not in (
            "increment",
            "fixed",
        ):
            raise ValueError(
                "Select one valid frame source and frame seed mode."
            )
        current: dict[str, torch.Tensor] = latent_image[0]
        if current["samples"].shape[0] != 1:
            raise ValueError(
                "Use a single initial image, with latent batch size one."
            )

        import folder_paths
        import nodes

        frames: list[torch.Tensor] = []
        for index, conditioning in enumerate(positive):
            if source == "original":
                current = {
                    name: value.clone() if torch.is_tensor(value) else value
                    for name, value in latent_image[0].items()
                }
            if unique_id and "server" in sys.modules:
                import server

                server.PromptServer.instance.send_progress_text(
                    f"Frame {index + 1} of {len(positive)}",
                    unique_id[0],
                    server.PromptServer.instance.client_id,
                )
            negative_conditioning: object = negative[
                0 if len(negative) == 1 else index
            ]
            if _frame_inputs is not None:
                conditioning, negative_conditioning, current = _frame_inputs(
                    index, frames[-1] if frames else None
                )
            elif source == "previous" and frames:
                current = dict(latent_image[0])
                current.update(nodes.VAEEncode().encode(vae[0], frames[-1])[0])
            current = nodes.common_ksampler(
                model[0],
                (seed[0] + (index if seed_mode == "increment" else 0))
                % (2**64),
                steps[0],
                cfg[0],
                sampler_name[0],
                scheduler[0],
                conditioning,
                negative_conditioning,
                current,
                denoise=denoise[0],
            )[0]
            decoded: torch.Tensor = (
                nodes.VAEDecode()
                .decode(vae[0], current)[0]
                .detach()
                .cpu()
                .float()
            )
            if decoded.shape[0] != 1 or (
                frames and decoded.shape != frames[0].shape
            ):
                raise ValueError(
                    "Every generated frame must have the same dimensions and batch size one."
                )
            frames.append(decoded)

        images: torch.Tensor = torch.cat(frames, dim=0)
        pictures: list[PIL.Image.Image] = [
            PIL.Image.fromarray(
                numpy.clip(frame[0].numpy() * 255, 0, 255).astype(numpy.uint8)
            )
            for frame in frames
        ]
        duration: int = max(10, round(100 / fps[0]) * 10)
        encoded: io.BytesIO = io.BytesIO()
        # Encode every frame explicitly: Image.save(save_all=True) merges equal
        # frames, which would lose the one-prompt/one-frame media contract.
        palette: list[PIL.Image.Image] = [
            picture.convert(
                "P", palette=PIL.Image.Palette.ADAPTIVE, colors=256
            )
            for picture in pictures
        ]
        header: list[bytes]
        header, _ = PIL.GifImagePlugin.getheader(palette[0], info={"loop": 0})
        for block in header:
            encoded.write(block)
        for picture in palette:
            for block in PIL.GifImagePlugin.getdata(
                picture,
                duration=duration,
                disposal=2,
                include_color_table=True,
            ):
                encoded.write(block)
        encoded.write(b";")

        root: pathlib.Path = pathlib.Path(
            folder_paths.get_output_directory()
        ).resolve()
        directory: str
        basename: str
        counter: int
        subfolder: str
        directory, basename, counter, subfolder, _ = (
            folder_paths.get_save_image_path(
                filename_prefix[0], str(root), images.shape[2], images.shape[1]
            )
        )
        if not pathlib.Path(directory).resolve().is_relative_to(root):
            raise ValueError(
                "GIF output must stay inside ComfyUI's output directory."
            )
        while True:
            filename: str = f"{basename}_{counter:05d}.gif"
            path: pathlib.Path = pathlib.Path(directory, filename)
            try:
                with path.open("xb") as output:
                    try:
                        output.write(encoded.getvalue())
                    except BaseException:
                        path.unlink(missing_ok=True)
                        raise
                break
            except FileExistsError:
                counter += 1
        return {
            "ui": {
                "images": [
                    {
                        "filename": filename,
                        "subfolder": subfolder,
                        "type": "output",
                    }
                ],
                "animated": (True,),
            },
            "result": (images,),
        }


NODE_CLASS_MAPPINGS: dict[str, type[ChainedImageToImageGIF]] = {
    "ImagegenChainedImageToImageGIF": ChainedImageToImageGIF,
}
"""Native registry for the first-party chained-frame output node."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenChainedImageToImageGIF": "Chained Image to Image GIF",
}
"""Readable workflow editor title."""
