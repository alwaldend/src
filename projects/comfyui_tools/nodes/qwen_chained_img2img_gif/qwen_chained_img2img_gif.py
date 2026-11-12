"""Recondition decoded animation frames with ComfyUI's native Qwen edit encoder.

The shared first-party GIF node owns sampling, decoding and media publication.
Resolve it through the native registry after custom-node loading has finished.
"""

import typing

import torch


class QwenChainedImageToImageGIF:
    """Edit each decoded previous frame with native Qwen image conditioning."""

    INPUT_IS_LIST: bool = True
    """Consume all frame instructions in one dependent queue job."""

    RETURN_TYPES: tuple[str] = ("IMAGE",)
    """Return the decoded frame batch using the shared GIF renderer."""

    FUNCTION: str = "generate"
    """Native ComfyUI execution entry point."""

    OUTPUT_NODE: bool = True
    """Publish the completed animation and its native preview."""

    CATEGORY: str = "imagegen/workflow"
    """Editor category shared with the first-party sequence nodes."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Accept raw frame instructions and the initial image for edit encoding."""
        import nodes

        schema: dict[str, object] = nodes.NODE_CLASS_MAPPINGS[
            "ImagegenChainedImageToImageGIF"
        ].INPUT_TYPES()
        required = typing.cast(dict[str, object], schema["required"])
        schema["required"] = {
            (
                "prompts"
                if name == "positive"
                else "negative_prompt"
                if name == "negative"
                else name
            ): (
                ("STRING", {"forceInput": True})
                if name in ("positive", "negative")
                else value
            )
            for name, value in required.items()
        } | {"clip": ("CLIP",), "reference_image": ("IMAGE",)}
        return schema

    def generate(
        self,
        model: list[object],
        prompts: list[str],
        negative_prompt: list[str],
        latent_image: list[dict[str, torch.Tensor]],
        vae: list[object],
        clip: list[object],
        reference_image: list[torch.Tensor],
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
    ) -> dict[str, object]:
        """Refresh image/latent conditioning per pose, sharing the GIF renderer."""
        import comfy_extras.nodes_qwen as nodes_qwen
        import node_helpers
        import nodes

        if len(clip) != 1 or len(reference_image) != 1:
            raise ValueError("Supply one Qwen text encoder and initial image.")
        if reference_image[0].ndim != 4 or reference_image[0].shape[0] != 1:
            raise ValueError("The initial reference must contain one image.")

        def frame_inputs(
            index: int, previous: torch.Tensor | None
        ) -> tuple[object, object, dict[str, torch.Tensor]]:
            """Use identical image pixels for edit references and latent encoding."""
            image: torch.Tensor = (
                reference_image[0]
                if previous is None or frame_source == ["original"]
                else previous
            )
            conditions: list[object] = []
            for text in (
                prompts[index],
                negative_prompt[0 if len(negative_prompt) == 1 else index],
            ):
                encoded = nodes_qwen.TextEncodeQwenImageEditPlus.execute(
                    clip=clip[0], prompt=text, vae=vae[0], image1=image
                ).result[0]
                conditions.append(
                    node_helpers.conditioning_set_values(
                        encoded,
                        {"reference_latents_method": "index_timestep_zero"},
                    )
                )
            latent: dict[str, torch.Tensor] = dict(latent_image[0])
            latent.update(nodes.VAEEncode().encode(vae[0], image)[0])
            return conditions[0], conditions[1], latent

        return nodes.NODE_CLASS_MAPPINGS[
            "ImagegenChainedImageToImageGIF"
        ]().generate(
            model=model,
            positive=list[object](prompts),
            negative=list[object](negative_prompt),
            latent_image=latent_image,
            vae=vae,
            seed=seed,
            steps=steps,
            cfg=cfg,
            sampler_name=sampler_name,
            scheduler=scheduler,
            denoise=denoise,
            fps=fps,
            filename_prefix=filename_prefix,
            unique_id=unique_id,
            frame_source=frame_source,
            frame_seed_mode=frame_seed_mode,
            _frame_inputs=frame_inputs,
        )


NODE_CLASS_MAPPINGS: dict[str, type[QwenChainedImageToImageGIF]] = {
    "ImagegenQwenChainedImageToImageGIF": QwenChainedImageToImageGIF,
}
"""Native registry for the model-specific decoded-image feedback node."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenQwenChainedImageToImageGIF": "Qwen Chained Image to Image GIF",
}
"""Editor title for the Qwen animation node."""
