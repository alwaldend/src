"""Search the packaged Danbooru vocabulary for prompt tags in ComfyUI."""

import collections.abc
import csv
import heapq
import pathlib
import typing


def _matches(
    source: typing.TextIO, query: str, include_deprecated: bool
) -> collections.abc.Iterator[tuple[bool, bool, int, str]]:
    """Yield matching tag names with their sort priority."""
    rows: csv.DictReader = csv.DictReader(source)
    if not {"name", "post_count", "is_deprecated"}.issubset(
        rows.fieldnames or ()
    ):
        raise ValueError("Danbooru tags.csv has an unexpected header.")

    for row in rows:
        name: str = row["name"]
        folded: str = name.casefold()
        if query not in folded:
            continue
        if not include_deprecated and row["is_deprecated"] == "True":
            continue
        yield (
            folded != query,
            not folded.startswith(query),
            -int(row["post_count"]),
            name,
        )


class _DanbooruTagSearch:
    """Return the highest-ranked matches as comma-separated prompt tags."""

    RETURN_TYPES: tuple[str] = ("STRING",)
    """Native string output suitable for Save Text or a prompt connection."""

    RETURN_NAMES: tuple[str] = ("tags",)
    """Label for the comma-separated tag list."""

    FUNCTION: str = "search"
    """ComfyUI execution entry point."""

    CATEGORY: str = "imagegen/workflow"
    """Menu location for the search utility."""

    @classmethod
    def INPUT_TYPES(cls) -> dict[str, object]:
        """Declare the query and result filters exposed by the App."""
        return {
            "required": {
                "query": ("STRING", {"default": "blue hair"}),
                "limit": ("INT", {"default": 20, "min": 1, "max": 100}),
                "include_deprecated": ("BOOLEAN", {"default": False}),
            },
        }

    def search(
        self, query: str, limit: int, include_deprecated: bool
    ) -> tuple[str]:
        """Scan the pinned CSV with memory bounded by the result limit."""
        if not 1 <= limit <= 100:
            raise ValueError("Limit must be between 1 and 100.")
        normalized: str = "_".join(query.casefold().split())
        if not normalized:
            return ("Enter a tag or part of a tag.",)

        vocabulary: pathlib.Path = pathlib.Path(__file__).with_name("tags.csv")
        try:
            with vocabulary.open(encoding="utf-8", newline="") as source:
                matches: list[tuple[bool, bool, int, str]] = heapq.nsmallest(
                    limit, _matches(source, normalized, include_deprecated)
                )
        except FileNotFoundError as error:
            raise RuntimeError(
                "Install the Danbooru tag search bundle so tags.csv is next "
                "to the search node."
            ) from error
        return (
            ", ".join(name for _, _, _, name in matches)
            or "No matching tags.",
        )


NODE_CLASS_MAPPINGS: dict[str, type[_DanbooruTagSearch]] = {
    "ImagegenDanbooruTagSearch": _DanbooruTagSearch,
}
"""Register the workflow's serialized node name with ComfyUI."""

NODE_DISPLAY_NAME_MAPPINGS: dict[str, str] = {
    "ImagegenDanbooruTagSearch": "Danbooru Tag Search",
}
"""Show a readable name for the search node in the editor."""

__all__: tuple[str, ...] = (
    "NODE_CLASS_MAPPINGS",
    "NODE_DISPLAY_NAME_MAPPINGS",
)
"""Expose the ComfyUI extension registration."""
