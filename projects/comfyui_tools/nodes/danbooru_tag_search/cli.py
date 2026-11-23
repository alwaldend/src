"""Search the pinned Danbooru tag vocabulary without a ComfyUI server."""

import argparse
import pathlib
import typing

from python.runfiles import runfiles

from projects.comfyui_tools.nodes import danbooru_tag_search


def _main() -> None:
    """Read local tags and print a comma-separated list for prompt editing."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(
        description="Search the local Danbooru tag vocabulary."
    )
    parser.add_argument("vocabulary", help=argparse.SUPPRESS)
    parser.add_argument("query", help="Tag name or fragment")
    parser.add_argument(
        "--limit",
        type=int,
        default=20,
        help="Maximum number of results (1 to 100)",
    )
    parser.add_argument(
        "--raw-tags", action="store_true", help="Keep underscores"
    )
    parser.add_argument("--include-deprecated", action="store_true")
    parser.add_argument(
        "--sort", choices=("relevance", "popularity"), default="relevance"
    )
    args: argparse.Namespace = parser.parse_args()
    if not args.query.strip():
        parser.error("Enter a tag or part of a tag.")

    files: runfiles.Runfiles | None = runfiles.Create()
    vocabulary: str | None = (
        files.Rlocation(args.vocabulary) if files else None
    )
    if not vocabulary:
        parser.error("Could not locate the pinned tag vocabulary.")
    sort_by: typing.Literal["relevance", "popularity"] = args.sort
    try:
        tags: list[str] = danbooru_tag_search.search_tags(
            pathlib.Path(vocabulary),
            args.query,
            args.limit,
            include_deprecated=args.include_deprecated,
            prompt_spaces=not args.raw_tags,
            sort_by=sort_by,
        )
    except (OSError, ValueError) as error:
        parser.error(str(error))
    print(", ".join(tags))


if __name__ == "__main__":
    _main()
