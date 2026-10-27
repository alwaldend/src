"""Check embedded fonts in actual Bazel-rendered native SVG variants."""

import base64
import io
import json
import os
import pathlib
import re
import sys
from xml.etree import ElementTree

import fontTools.ttLib


def main() -> None:
    """Compare rendered glyph coverage and geometry with the original font."""
    original_path: pathlib.Path = pathlib.Path(sys.argv[1])
    original: fontTools.ttLib.TTFont = fontTools.ttLib.TTFont(original_path)
    original_map: dict[int, str] = original.getBestCmap()
    report: list[dict[str, str | int]] = []
    for filename in sys.argv[2:]:
        svg: str = pathlib.Path(filename).read_text()
        match: re.Match[str] | None = re.search(
            r"data:font/ttf;base64,([A-Za-z0-9+/=]+)", svg
        )
        assert match is not None, "SVG must carry its own font"
        data: bytes = base64.b64decode(match.group(1), validate=True)
        assert len(data) < original_path.stat().st_size // 4, (
            "Native output still embeds substantially more than its used glyphs"
        )
        embedded: fontTools.ttLib.TTFont = fontTools.ttLib.TTFont(
            io.BytesIO(data)
        )
        assert embedded["name"].getDebugName(1) == original[
            "name"
        ].getDebugName(1)
        document: ElementTree.Element = ElementTree.fromstring(svg)
        text: str = "".join(
            "".join(node.itertext())
            for node in document.iter()
            if node.tag.rsplit("}", 1)[-1] in ("text", "foreignObject")
        )
        assert text.strip(), "The rendered fixture must contain labels"
        codepoints: set[int] = set(map(ord, text)) & original_map.keys()
        embedded_map: dict[int, str] = embedded.getBestCmap()
        assert codepoints <= embedded_map.keys(), "Rendered glyphs are missing"
        for codepoint in codepoints:
            assert original["glyf"][original_map[codepoint]].getCoordinates(
                original["glyf"]
            ) == embedded["glyf"][embedded_map[codepoint]].getCoordinates(
                embedded["glyf"]
            )
            assert (
                original["hmtx"][original_map[codepoint]]
                == embedded["hmtx"][embedded_map[codepoint]]
            )
        report.append(
            {
                "svg": filename,
                "font_bytes": len(data),
                "glyphs": len(codepoints),
            }
        )
        embedded.close()
    original.close()
    output: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    (output / "native-fonts.json").write_text(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
