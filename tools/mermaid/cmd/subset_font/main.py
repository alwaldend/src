"""Subset a pinned font from rendered SVG input without rewriting SVG geometry."""

import html
import io
import sys

import fontTools.subset
import fontTools.ttLib


def main() -> None:
    """Write a deterministic TrueType subset to stdout for Mermaid's CSS input."""
    if len(sys.argv) != 2:
        raise ValueError("Expected one pinned font filename")
    text: str = html.unescape(sys.stdin.read())
    options: fontTools.subset.Options = fontTools.subset.Options()
    options.glyph_names = True
    font: fontTools.ttLib.TTFont = fontTools.ttLib.TTFont(
        sys.argv[1], recalcTimestamp=False
    )
    subsetter: fontTools.subset.Subsetter = fontTools.subset.Subsetter(
        options=options
    )
    subsetter.populate(text=text)
    subsetter.subset(font)
    output: io.BytesIO = io.BytesIO()
    font.save(output)
    font.close()
    sys.stdout.buffer.write(output.getvalue())


if __name__ == "__main__":
    main()
