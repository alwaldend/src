"""Exercise ordered native-sampler feedback and actual GIF output before GPU acceptance."""

import os
import pathlib
import shutil
import tempfile
import types
import unittest
from unittest import mock

import PIL.Image
import torch

from projects.comfyui_tools.nodes.chained_img2img_gif import (
    chained_img2img_gif,
)


class _Decoder:
    """Decode deterministic latent colors through the native VAE interface."""

    def decode(self, samples: torch.Tensor) -> torch.Tensor:
        """Return distinguishable RGB frames for an increasing latent value."""
        return samples[:, :3].movedim(1, -1) / 4


class _ChainedGIFTest(unittest.TestCase):
    """Check media properties, feedback, validation and interrupted writes."""

    def test_chained_feedback_and_gif(self) -> None:
        """Feed each latent forward and preserve prompt order in one real GIF."""
        with tempfile.TemporaryDirectory() as directory:
            calls: list[tuple[object, float, object, int]] = []

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
                    "nodes": types.SimpleNamespace(common_ksampler=sample),
                    "folder_paths": folders,
                },
            ):
                result: chained_img2img_gif._Result = (
                    chained_img2img_gif.ChainedImageToImageGIF().generate(
                        model=[object()],
                        positive=["first", "second", "third"],
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
                    ("second", 1.0, "shared", 11),
                    ("third", 2.0, "shared", 12),
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
