"""Native Mermaid diagrams using the site's shared color palette."""

load("//projects/alwaldend.com/assets:defs.bzl", "site_stylesheet")
load("//tools/mermaid:defs.bzl", "mermaid_svg", "mermaid_webp")

_SITE_THEME = "//projects/alwaldend.com/pkg/mermaid:theme.json"
_SITE_PALETTE = "//projects/alwaldend.com/assets:mermaid_palette"
_SITE_FONT = "@com_alwaldend_src_tools_drawio_fonts//:LiberationMono-Regular.ttf"

# Keep the existing diagram API while sharing the Sass compiler action.
mermaid_palette = site_stylesheet

def mermaid_site_svg(name, src, out, **kwargs):
    """Render a native SVG and its .dark.svg sibling for site image embeds."""
    mermaid_svg(
        name = name,
        src = src,
        out = out,
        dark_out = out.removesuffix(".svg") + ".dark.svg",
        native = True,
        theme = _SITE_THEME,
        palette = _SITE_PALETTE,
        label_font = _SITE_FONT,
        **kwargs
    )

def mermaid_site_webp(name, src, out, **kwargs):
    """Render a site diagram in the selected color scheme (light by default)."""
    mermaid_webp(
        name = name,
        src = src,
        out = out,
        native = True,
        theme = _SITE_THEME,
        palette = _SITE_PALETTE,
        label_font = _SITE_FONT,
        **kwargs
    )
