"""Expand panel prompts into paired lists for native ComfyUI workflow mapping."""

import json
import re

WEB_DIRECTORY: str = "./web"
"""First-party frontend script directory served by ComfyUI."""


class PromptSequence:
    """Provide shared reference instructions and one filename per panel prompt."""

    RETURN_TYPES: tuple[str, str] = ("STRING", "STRING")
    """Native prompt and Save Image filename-prefix output types."""

    RETURN_NAMES: tuple[str, str] = ("prompt", "filename_prefix")
    """User-visible output names identifying the paired panel values."""

    OUTPUT_IS_LIST: tuple[bool, bool] = (True, True)
    """Ask ComfyUI to map downstream nodes over both ordered panel lists."""

    FUNCTION: str = "expand"
    """Framework entry point called to expand the widget input."""

    CATEGORY: str = "imagegen/workflow"
    """Menu location for this first-party workflow utility."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Declare multiline prompts, shared context and panel naming controls."""
        return {
            "required": {
                "prompts": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "A white robot rides a bicycle through a snowy park.\nA white robot stops beside a frozen pond.\nA white robot reads a book on a park bench.",
                    },
                ),
                "context": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "Use the character from image 1 and the illustration style, colors and linework from image 2. Create a new comic panel for this scene:",
                    },
                ),
                "separator": (["line", "blank line", "JSON"],),
                "filename_prefix": ("STRING", {"default": "comic"}),
            }
        }

    def expand(
        self,
        prompts: str,
        context: str,
        separator: str,
        filename_prefix: str,
    ) -> tuple[list[str], list[str]]:
        """Return ordered nonempty panels and safe paired output prefixes."""
        if not re.fullmatch(
            r"[A-Za-z0-9][A-Za-z0-9_-]{0,63}", filename_prefix
        ):
            raise ValueError(
                "Filename prefix must contain 1 to 64 letters, digits, underscores or hyphens."
            )
        segments: list[str]
        if separator == "JSON":
            try:
                decoded: object = json.loads(prompts)
            except json.JSONDecodeError:
                raise ValueError(
                    "Enter a JSON array of panel strings."
                ) from None
            if not isinstance(decoded, list) or not all(
                isinstance(panel, str) for panel in decoded
            ):
                raise ValueError("Enter a JSON array of panel strings.")
            segments = decoded
        elif separator in ("line", "blank line"):
            normalized: str = prompts.replace("\r\n", "\n").replace("\r", "\n")
            segments = (
                re.split(r"\n\s*\n", normalized)
                if separator == "blank line"
                else normalized.splitlines()
            )
        else:
            raise ValueError("Select line, blank line or JSON separation.")
        panels: list[str] = [text.strip() for text in segments if text.strip()]
        if not panels:
            raise ValueError("Enter at least one panel prompt.")
        shared: str = context.strip()
        texts: list[str] = [
            f"{shared}\n\n{panel}" if shared else panel for panel in panels
        ]
        prefixes: list[str] = [
            f"{filename_prefix}_{i:03d}" for i in range(1, len(panels) + 1)
        ]
        return texts, prefixes


NODE_CLASS_MAPPINGS: dict[str, type[PromptSequence]] = {
    "ImagegenPromptSequence": PromptSequence,
}
"""Framework registry for the sole enabled first-party node."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenPromptSequence": "Prompt Sequence",
}
"""Readable node title shown in the ComfyUI workflow editor."""
