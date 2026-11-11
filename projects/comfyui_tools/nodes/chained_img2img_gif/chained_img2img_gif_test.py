"""Exercise ordered native-sampler feedback and actual GIF output before GPU acceptance."""

import os
import pathlib
import shutil
import tempfile
import types
import typing
import unittest
from unittest import mock

import PIL.Image
import torch

from projects.comfyui_tools.nodes.chained_img2img_gif import (
    chained_img2img_gif,
)
from projects.comfyui_tools.nodes.qwen_chained_img2img_gif import (
    qwen_chained_img2img_gif,
)


class _Decoder:
    """Decode deterministic latent colors through the native VAE interface."""

    def __init__(self, temporal: bool = False) -> None:
        """Select the single-frame temporal layout returned by Qwen's VAE."""
        self.temporal = temporal

    def decode(self, samples: torch.Tensor) -> torch.Tensor:
        """Return distinguishable RGB frames in image or temporal layout."""
        images: torch.Tensor = samples[:, :3].movedim(1, -1) / 4
        return images.unsqueeze(1).bfloat16() if self.temporal else images


class _NativeDecoder:
    """Supply ComfyUI's documented decoded-image normalization boundary."""

    def decode(
        self, vae: _Decoder, samples: dict[str, torch.Tensor]
    ) -> tuple[torch.Tensor]:
        """Normalize temporal VAE output to the native NHWC image contract."""
        images: torch.Tensor = vae.decode(samples["samples"])
        if images.ndim == 5:
            images = images.reshape(-1, *images.shape[-3:])
        return (images,)


class _ChainedGIFTest(unittest.TestCase):
    """Check media properties, feedback, validation and interrupted writes."""

    def test_original_reference_and_fixed_noise(self) -> None:
        """Preserve the original masked latent and reuse noise for every pose."""
        for temporal in (False, True):
            with self.subTest(temporal=temporal):
                self._check_original_reference(temporal)

    def _check_original_reference(self, temporal: bool) -> None:
        """Publish a real GIF for both spatial and single-frame temporal VAEs."""
        with tempfile.TemporaryDirectory() as directory:
            seeds: list[int] = []
            originals: list[float] = []
            initial: torch.Tensor = torch.zeros(1, 4, 8, 8)
            mask: torch.Tensor = torch.ones(1, 8, 8)

            def sample(
                model: object,
                seed: int,
                steps: int,
                cfg: float,
                sampler_name: str,
                scheduler: str,
                positive: object,
                negative: object,
                latent: dict[str, torch.Tensor],
                denoise: float,
            ) -> tuple[dict[str, torch.Tensor]]:
                """Mutate the supplied copy to expose accumulated reference drift."""
                seeds.append(seed)
                originals.append(latent["samples"].mean().item())
                self.assertTrue(torch.equal(latent["noise_mask"], mask))
                latent["samples"].add_(1)
                return (latent,)

            folders: types.SimpleNamespace = types.SimpleNamespace(
                get_output_directory=lambda: directory,
                get_save_image_path=lambda prefix, root, width, height: (
                    root,
                    prefix,
                    1,
                    "",
                    prefix,
                ),
            )
            with mock.patch.dict(
                "sys.modules",
                {
                    "nodes": types.SimpleNamespace(
                        common_ksampler=sample, VAEDecode=_NativeDecoder
                    ),
                    "folder_paths": folders,
                },
            ):
                chained_img2img_gif.ChainedImageToImageGIF().generate(
                    model=[object()],
                    positive=["rest", "press", "rest"],
                    negative=["shared"],
                    latent_image=[{"samples": initial, "noise_mask": mask}],
                    vae=[_Decoder(temporal)],
                    seed=[42],
                    steps=[2],
                    cfg=[3.0],
                    sampler_name=["euler"],
                    scheduler=["beta"],
                    denoise=[0.6],
                    fps=[6.0],
                    filename_prefix=["anchored"],
                    frame_source=["original"],
                    frame_seed_mode=["fixed"],
                )
            self.assertEqual(seeds, [42, 42, 42])
            self.assertEqual(originals, [0.0, 0.0, 0.0])
            self.assertEqual(initial.mean().item(), 0)
            with PIL.Image.open(
                next(pathlib.Path(directory).glob("*.gif"))
            ) as gif:
                self.assertEqual(gif.n_frames, 3)
            artifact: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
            if artifact:
                shutil.copyfile(
                    next(pathlib.Path(directory).glob("*.gif")),
                    pathlib.Path(
                        artifact, f"anchored-fixed-noise-{temporal}.gif"
                    ),
                )

    def test_qwen_previous_image_conditioning(self) -> None:
        """Re-encode each decoded frame as the next edit reference at denoise one."""
        references: list[float] = []
        sampled_references: list[float] = []

        class EditEncoder:
            """Observe the native Qwen edit encoder boundary."""

            @classmethod
            def execute(
                cls,
                clip: object,
                prompt: str,
                vae: object,
                image1: torch.Tensor,
            ) -> types.SimpleNamespace:
                """Retain the actual image seen by positive and negative encoding."""
                reference: float = image1.mean().item()
                references.append(reference)
                return types.SimpleNamespace(
                    result=({"prompt": prompt, "reference": reference},)
                )

        class ImageEncoder:
            """Encode the same pixel reference supplied to Qwen conditioning."""

            def encode(
                self, vae: object, pixels: torch.Tensor
            ) -> tuple[dict[str, torch.Tensor]]:
                """Use invertible uniform colors for the native image boundary."""
                return (
                    {
                        "samples": pixels.movedim(-1, 1)
                        .mean(dim=1, keepdim=True)
                        .repeat(1, 4, 1, 1)
                        * 4
                    },
                )

        def sample(
            *args: object, **kwargs: object
        ) -> tuple[dict[str, torch.Tensor]]:
            """Require both conditioning branches and sampler image to advance."""
            positive = typing.cast(dict[str, object], args[6])
            negative = typing.cast(dict[str, object], args[7])
            latent = typing.cast(dict[str, torch.Tensor], args[8])
            expected: float = len(sampled_references) / 4
            self.assertEqual(positive["reference"], expected)
            self.assertEqual(negative["reference"], expected)
            self.assertEqual(latent["samples"].mean().item(), expected * 4)
            self.assertEqual(kwargs["denoise"], 1.0)
            sampled_references.append(expected)
            return ({"samples": latent["samples"] + 1},)

        with tempfile.TemporaryDirectory() as directory:
            folders = types.SimpleNamespace(
                get_output_directory=lambda: directory,
                get_save_image_path=lambda prefix, root, width, height: (
                    root,
                    prefix,
                    1,
                    "",
                    prefix,
                ),
            )
            with mock.patch.dict(
                "sys.modules",
                {
                    "nodes": types.SimpleNamespace(
                        common_ksampler=sample,
                        VAEDecode=_NativeDecoder,
                        VAEEncode=ImageEncoder,
                        NODE_CLASS_MAPPINGS=chained_img2img_gif.NODE_CLASS_MAPPINGS,
                    ),
                    "folder_paths": folders,
                    "comfy_extras": types.ModuleType("comfy_extras"),
                    "comfy_extras.nodes_qwen": types.SimpleNamespace(
                        TextEncodeQwenImageEditPlus=EditEncoder
                    ),
                    "node_helpers": types.SimpleNamespace(
                        conditioning_set_values=lambda conditioning, values: (
                            conditioning
                        )
                    ),
                },
            ):
                result = qwen_chained_img2img_gif.QwenChainedImageToImageGIF().generate(
                    model=[object()],
                    prompts=["first", "second", "third"],
                    negative_prompt=["shared"],
                    latent_image=[{"samples": torch.zeros(1, 4, 8, 8)}],
                    vae=[_Decoder(True)],
                    clip=[object()],
                    reference_image=[torch.zeros(1, 8, 8, 3)],
                    seed=[42],
                    steps=[2],
                    cfg=[4.0],
                    sampler_name=["euler"],
                    scheduler=["simple"],
                    denoise=[1.0],
                    fps=[6.0],
                    filename_prefix=["qwen_chain"],
                    frame_source=["previous"],
                    frame_seed_mode=["fixed"],
                )
            self.assertEqual(references, [0.0, 0.0, 0.25, 0.25, 0.5, 0.5])
            self.assertEqual(result["result"][0].shape, (3, 8, 8, 3))
            path = next(pathlib.Path(directory).glob("*.gif"))
            with PIL.Image.open(path) as gif:
                self.assertEqual(gif.n_frames, 3)
            artifact = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
            if artifact:
                shutil.copyfile(
                    path, pathlib.Path(artifact, "qwen-chained-reference.gif")
                )

    def test_chained_feedback_and_gif(self) -> None:
        """Use decoded pixels after a nonidentity VAE round trip, retaining masks."""
        with tempfile.TemporaryDirectory() as directory:
            calls: list[tuple[object, float, object, int]] = []
            mask: torch.Tensor = torch.ones(1, 8, 8)

            class ImageEncoder:
                """Expose latent reuse by changing the decoded-image round trip."""

                def encode(
                    self, vae: object, pixels: torch.Tensor
                ) -> tuple[dict[str, torch.Tensor]]:
                    """Produce a deliberately different latent from decoded pixels."""
                    return (
                        {
                            "samples": pixels.movedim(-1, 1)
                            .mean(dim=1, keepdim=True)
                            .repeat(1, 4, 1, 1)
                            * 2
                        },
                    )

            def sample(
                model: object,
                seed: int,
                steps: int,
                cfg: float,
                sampler_name: str,
                scheduler: str,
                positive: object,
                negative: object,
                latent: dict[str, torch.Tensor],
                denoise: float,
            ) -> tuple[dict[str, torch.Tensor]]:
                """Record native call inputs and advance a deterministic latent."""
                self.assertTrue(torch.equal(latent["noise_mask"], mask))
                calls.append(
                    (positive, latent["samples"].mean().item(), negative, seed)
                )
                return ({"samples": latent["samples"] + 1},)

            folders: types.SimpleNamespace = types.SimpleNamespace(
                get_output_directory=lambda: directory,
                get_save_image_path=lambda prefix, root, width, height: (
                    root,
                    prefix,
                    1,
                    "",
                    prefix,
                ),
            )
            with mock.patch.dict(
                "sys.modules",
                {
                    "nodes": types.SimpleNamespace(
                        common_ksampler=sample,
                        VAEDecode=_NativeDecoder,
                        VAEEncode=ImageEncoder,
                        NODE_CLASS_MAPPINGS=chained_img2img_gif.NODE_CLASS_MAPPINGS,
                    ),
                    "folder_paths": folders,
                },
            ):
                result: chained_img2img_gif._Result = (
                    chained_img2img_gif.ChainedImageToImageGIF().generate(
                        model=[object()],
                        positive=["first", "second", "third"],
                        negative=["shared"],
                        latent_image=[
                            {
                                "samples": torch.zeros(1, 4, 8, 8),
                                "noise_mask": mask,
                            }
                        ],
                        vae=[_Decoder()],
                        seed=[10],
                        steps=[2],
                        cfg=[3.0],
                        sampler_name=["euler"],
                        scheduler=["beta"],
                        denoise=[0.45],
                        fps=[6.0],
                        filename_prefix=[
                            "comic_001",
                            "comic_002",
                            "comic_003",
                        ],
                    )
                )
            self.assertEqual(
                calls,
                [
                    ("first", 0.0, "shared", 10),
                    ("second", 0.5, "shared", 11),
                    ("third", 0.75, "shared", 12),
                ],
            )
            output: pathlib.Path = next(pathlib.Path(directory).glob("*.gif"))
            with PIL.Image.open(output) as gif:
                self.assertEqual(gif.n_frames, 3)
                self.assertEqual(gif.info["loop"], 0)
                self.assertEqual(gif.size, (8, 8))
                for index in range(3):
                    gif.seek(index)
                    self.assertEqual(gif.info["duration"], 170)
            self.assertEqual(result["result"][0].shape, (3, 8, 8, 3))
            self.assertEqual(result["ui"]["animated"], (True,))
            artifact: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
            if artifact:
                shutil.copyfile(
                    output, pathlib.Path(artifact, "chained-feedback.gif")
                )

    def test_invalid_input_before_sampling(self) -> None:
        """Reject mismatched batches, frame inputs and output paths early."""
        inputs: dict[str, object] = {
            "model": [object()],
            "positive": ["first", "second"],
            "negative": ["shared"],
            "latent_image": [{"samples": torch.zeros(1, 4, 8, 8)}],
            "vae": [_Decoder()],
            "seed": [10],
            "steps": [2],
            "cfg": [3.0],
            "sampler_name": ["euler"],
            "scheduler": ["beta"],
            "denoise": [0.45],
            "fps": [6.0],
            "filename_prefix": ["comic"],
        }
        for patch in [
            {"positive": []},
            {"negative": [1, 2, 3]},
            {"model": [1, 2]},
            {"fps": [0.0]},
            {"filename_prefix": ["../outside"]},
            {"latent_image": [{"samples": torch.zeros(2, 4, 8, 8)}]},
            {"frame_source": ["invalid"]},
            {"frame_seed_mode": ["invalid"]},
            {"frame_source": []},
        ]:
            with self.subTest(patch=patch), self.assertRaises(ValueError):
                chained_img2img_gif.ChainedImageToImageGIF().generate(
                    **(inputs | patch)
                )

    def test_interrupted_sampling_leaves_no_gif(self) -> None:
        """Avoid publishing a partial GIF when a native frame sample fails."""
        with tempfile.TemporaryDirectory() as directory:
            folders: types.SimpleNamespace = types.SimpleNamespace(
                get_output_directory=lambda: directory
            )
            with mock.patch.dict(
                "sys.modules",
                {
                    "nodes": types.SimpleNamespace(
                        common_ksampler=mock.Mock(
                            side_effect=RuntimeError("Interrupted")
                        )
                    ),
                    "folder_paths": folders,
                },
            ):
                with self.assertRaisesRegex(RuntimeError, "Interrupted"):
                    chained_img2img_gif.ChainedImageToImageGIF().generate(
                        model=[object()],
                        positive=["first"],
                        negative=["shared"],
                        latent_image=[{"samples": torch.zeros(1, 4, 8, 8)}],
                        vae=[_Decoder()],
                        seed=[10],
                        steps=[2],
                        cfg=[3.0],
                        sampler_name=["euler"],
                        scheduler=["beta"],
                        denoise=[0.45],
                        fps=[6.0],
                        filename_prefix=["comic"],
                    )
            self.assertEqual(list(pathlib.Path(directory).iterdir()), [])


if __name__ == "__main__":
    unittest.main()
