"""Expand panel prompts into paired lists for native ComfyUI workflow mapping."""

import collections.abc
import json
import re
import typing
import uuid


class _Model(typing.Protocol):
    """Native model patcher interface used to attach local sampling progress."""

    def clone(self) -> "_Model":
        """Return an independent patcher sharing the existing model weights."""
        ...

    def add_wrapper_with_key(
        self,
        kind: str,
        key: str,
        wrapper: collections.abc.Callable[..., object],
    ) -> None:
        """Add one wrapper through ComfyUI's supported model patcher API."""
        ...


class PromptSequence:
    """Provide ordered panel prompts, filenames and one shared negative prompt."""

    RETURN_TYPES: tuple[str, str, str, str] = (
        "STRING",
        "STRING",
        "STRING",
        "MODEL",
    )
    """Native prompt/prefix/negative strings and the optional tracked model."""

    RETURN_NAMES: tuple[str, str, str, str] = (
        "prompt",
        "filename_prefix",
        "negative_prompt",
        "model",
    )
    """User-visible output names identifying the paired panel values."""

    OUTPUT_IS_LIST: tuple[bool, bool, bool, bool] = (True, True, False, False)
    """Map panel lists while sharing the scalar negative and tracked model."""

    FUNCTION: str = "execute"
    """Framework entry point called to expand the widget input."""

    CATEGORY: str = "imagegen/workflow"
    """Menu location for this first-party workflow utility."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Append optional negative text without shifting saved widget values."""
        return {
            "required": {
                "prompts": (
                    "STRING",
                    {
                        "multiline": True,
                        "tooltip": "In line mode, enter one prompt per line. Blank lines are ignored and can be used for visual spacing.",
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
            },
            "optional": {
                "negative_prompt": (
                    "STRING",
                    {
                        "multiline": True,
                        "default": "",
                        "tooltip": "Shared negative prompt for every panel. Leave empty to use an empty negative prompt.",
                    },
                ),
                "model": ("MODEL",),
            },
        }

    def expand(
        self,
        prompts: str,
        context: str,
        separator: str,
        filename_prefix: str,
        negative_prompt: str = "",
    ) -> tuple[list[str], list[str], str]:
        """Return panel lists and unchanged shared negative text, without encoding."""
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
        return texts, prefixes, negative_prompt

    @classmethod
    def IS_CHANGED(cls, **inputs: object) -> float:
        """Create a new run identity even when a queued graph's text is unchanged."""
        return float("nan")

    def execute(
        self,
        prompts: str,
        context: str,
        separator: str,
        filename_prefix: str,
        negative_prompt: str = "",
        model: _Model | None = None,
    ) -> tuple[list[str], list[str], str, _Model | None]:
        """Keep panels in one job, add a run prefix and report native sampler progress."""
        texts: list[str]
        prefixes: list[str]
        negative: str
        texts, prefixes, negative = self.expand(
            prompts, context, separator, filename_prefix, negative_prompt
        )
        run: str = uuid.uuid4().hex
        prefixes = [
            f"{filename_prefix}_{run}_{index:03d}"
            for index in range(1, len(texts) + 1)
        ]
        if model is not None:
            import comfy.patcher_extension
            import server

            model = model.clone()
            panel: int = 0

            def sample(
                executor: collections.abc.Callable[..., object],
                *args: object,
                **kwargs: object,
            ) -> object:
                """Label each complete native sampling invocation without changing it."""
                nonlocal panel
                panel += 1
                server.PromptServer.instance.send_progress_text(
                    f"Panel {panel} of {len(texts)}",
                    server.PromptServer.instance.last_node_id,
                    server.PromptServer.instance.client_id,
                )
                return executor(*args, **kwargs)

            model.add_wrapper_with_key(
                comfy.patcher_extension.WrappersMP.OUTER_SAMPLE,
                "imagegen_sequence_progress",
                sample,
            )
        return texts, prefixes, negative, model


NODE_CLASS_MAPPINGS: dict[str, type[PromptSequence]] = {
    "ImagegenPromptSequence": PromptSequence,
}
"""Framework registry for the sole enabled first-party node."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenPromptSequence": "Prompt Sequence",
}
"""Readable node title shown in the ComfyUI workflow editor."""
