"""Search the packaged Danbooru vocabulary for prompt tags in ComfyUI."""

import collections.abc
import csv
import heapq
import pathlib
import typing


def _matches(
    source: typing.TextIO,
    query: str,
    include_deprecated: bool,
    sort_by: typing.Literal["relevance", "popularity"],
) -> collections.abc.Iterator[tuple[bool, bool, int, str]]:
    """Yield matching tag names with their sort priority."""
    rows: csv.DictReader = csv.DictReader(source)
    if not {"name", "post_count", "is_deprecated"}.issubset(
        rows.fieldnames or ()
    ):
        raise ValueError("Danbooru tags.csv has an unexpected header.")

    rank_by_relevance: bool = sort_by == "relevance"
    for row in rows:
        name: str = row["name"]
        folded: str = name.casefold()
        if query not in folded:
            continue
        if not include_deprecated and row["is_deprecated"] == "True":
            continue
        yield (
            rank_by_relevance and folded != query,
            rank_by_relevance and not folded.startswith(query),
            -int(row["post_count"]),
            name,
        )


def search_tags(
    vocabulary: pathlib.Path,
    query: str,
    limit: int,
    include_deprecated: bool = False,
    prompt_spaces: bool = False,
    sort_by: typing.Literal["relevance", "popularity"] = "relevance",
) -> list[str]:
    """Return ranked tags from a local vocabulary for the node and CLI."""
    if not 1 <= limit <= 100:
        raise ValueError("Limit must be between 1 and 100.")
    normalized: str = "_".join(query.casefold().split())
    if not normalized:
        return []
    with vocabulary.open(encoding="utf-8", newline="") as source:
        matches: list[tuple[bool, bool, int, str]] = heapq.nsmallest(
            limit,
            _matches(source, normalized, include_deprecated, sort_by),
        )
    return [
        name.replace("_", " ") if prompt_spaces else name
        for _, _, _, name in matches
    ]


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
            "optional": {
                "prompt_spaces": (
                    "BOOLEAN",
                    {
                        "default": False,
                        "tooltip": "Replace underscores with spaces for NoobAI prompts.",
                    },
                ),
                "sort_by": (
                    ["relevance", "popularity"],
                    {
                        "tooltip": "Popularity ranks all matching tags by post count."
                    },
                ),
            },
        }

    def search(
        self,
        query: str,
        limit: int,
        include_deprecated: bool,
        prompt_spaces: bool = False,
        sort_by: typing.Literal["relevance", "popularity"] = "relevance",
    ) -> tuple[str]:
        """Scan the pinned CSV with memory bounded by the result limit."""
        vocabulary: pathlib.Path = pathlib.Path(__file__).with_name("tags.csv")
        try:
            matches: list[str] = search_tags(
                vocabulary,
                query,
                limit,
                include_deprecated,
                prompt_spaces,
                sort_by,
            )
        except FileNotFoundError as error:
            raise RuntimeError(
                "Install the Danbooru tag search bundle so tags.csv is next "
                "to the search node."
            ) from error
        if not query.strip():
            return ("Enter a tag or part of a tag.",)
        return (", ".join(matches) or "No matching tags.",)


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
    "search_tags",
)
"""Expose the ComfyUI extension registration."""
