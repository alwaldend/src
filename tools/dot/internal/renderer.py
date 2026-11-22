"""Capture DOT from a command and render it using the shared site theme."""

import base64
import html
import json
import pathlib
import subprocess
import sys
import typing

_Theme = typing.TypedDict(
    "_Theme",
    {
        "colors": dict[str, str],
        "modes": dict[str, dict[str, str]],
        "typography": dict[str, dict[str, str]],
    },
)
"""Shared theme fields needed to render a diagram."""


def _graphviz_color(color: str) -> str:
    """Expand CSS shorthand into the six-digit form accepted by Graphviz."""
    if len(color) == 4:
        return "#" + "".join(digit * 2 for digit in color[1:])
    return color


def _render(
    dot: str, graph: str, theme_path: str, font_path: str, license_path: str
) -> str:
    """Lay out the DOT graph and embed its typography and dark-mode colors."""
    theme: _Theme = json.loads(
        pathlib.Path(theme_path).read_text(encoding="utf-8")
    )
    colors: dict[str, str] = theme["colors"]
    light: dict[str, str] = theme["modes"]["light"]
    dark: dict[str, str] = theme["modes"]["dark"]
    diagram: dict[str, str] = theme["typography"]["diagram"]
    canvas: str = _graphviz_color(colors[light["canvas"]])
    ink: str = _graphviz_color(colors[light["ink"]])
    family: str = diagram["family"]
    font_name: str = family.split(",", 1)[0]
    font_size: str = diagram["size"].removesuffix("px")
    font: str = base64.b64encode(pathlib.Path(font_path).read_bytes()).decode(
        "ascii"
    )
    license_text: str = html.escape(
        pathlib.Path(license_path).read_text(encoding="utf-8")
    )

    # Courier supplies built-in monospaced metrics for the embedded face.
    svg: str = subprocess.run(
        [
            dot,
            "-Tsvg",
            f"-Gbgcolor={canvas}",
            "-Gpad=0.2",
            f"-Ncolor={ink}",
            f"-Ecolor={ink}",
            f"-Nfontcolor={ink}",
            f"-Efontcolor={ink}",
            "-Nfontname=Courier",
            "-Efontname=Courier",
            f"-Nfontsize={font_size}",
            f"-Efontsize={font_size}",
            graph,
        ],
        check=True,
        stdout=subprocess.PIPE,
        encoding="utf-8",
    ).stdout
    style: str = f"""
<metadata>{license_text}</metadata>
<style>
@font-face {{ font-family: '{font_name}'; src: url(data:font/ttf;base64,{font}); }}
text {{ font-family: {family}; }}
@media (prefers-color-scheme: dark) {{
  [fill="{canvas}"] {{ fill: {colors[dark["canvas"]]}; }}
  [fill="{ink}"] {{ fill: {colors[dark["ink"]]}; }}
  [stroke="{ink}"] {{ stroke: {colors[dark["ink"]]}; }}
}}
</style>"""
    opening_tag_end: int = svg.index(">", svg.index("<svg ")) + 1
    return svg[:opening_tag_end] + style + svg[opening_tag_end:]


def _main(operation: str, output: str, *inputs: str) -> None:
    """Write the declared output of a DOT build action."""
    path: pathlib.Path = pathlib.Path(output)
    if operation == "export":
        with path.open("wb") as stream:
            subprocess.run(inputs, check=True, stdout=stream)
    elif operation == "svg":
        path.write_text(_render(*inputs), encoding="utf-8")
    else:
        raise ValueError(f"unknown DOT operation {operation!r}")


if __name__ == "__main__":
    _main(*sys.argv[1:])
