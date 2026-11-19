"""Site diagram outputs using the shared Alwaldend identity preset."""

load("//tools/mermaid:defs.bzl", "mermaid_svg", "mermaid_webp")

_SITE_THEME = "//internal/assets:mermaid_theme"

def mermaid_diagram(
        name,
        src,
        out,
        theme = _SITE_THEME,
        webp_aspect_ratio = [5, 2],
        webp_scale = 1,
        testonly = False,
        visibility = None):
    """Render light and dark SVG and WebP files from one Mermaid source.

    Args:
        name: Target that builds all four outputs. The _svg, _webp, and
            _dark_webp targets build the SVG pair or an individual WebP.
        src: Mermaid source label.
        out: Package-relative output stem, without a file extension.
        theme: Bundled native Mermaid theme.
        webp_aspect_ratio: WebP canvas ratio. Use [] for the natural canvas.
        webp_scale: WebP density multiplier.
        testonly: Whether only test targets may depend on the diagram.
        visibility: Visibility of the aggregate and individual targets.
    """
    mermaid_svg(
        name = name + "_svg",
        src = src,
        out = out + ".svg",
        dark_out = out + ".dark.svg",
        native = True,
        testonly = testonly,
        theme = theme,
        visibility = visibility,
    )

    for color_scheme, target_suffix, file_suffix in [
        ("light", "_webp", ".webp"),
        ("dark", "_dark_webp", ".dark.webp"),
    ]:
        mermaid_webp(
            name = name + target_suffix,
            src = src,
            out = out + file_suffix,
            aspect_ratio = webp_aspect_ratio,
            color_scheme = color_scheme,
            native = True,
            scale = webp_scale,
            testonly = testonly,
            theme = theme,
            visibility = visibility,
        )

    native.filegroup(
        name = name,
        srcs = [
            ":" + name + "_svg",
            ":" + name + "_webp",
            ":" + name + "_dark_webp",
        ],
        testonly = testonly,
        visibility = visibility,
    )
