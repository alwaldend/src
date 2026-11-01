"""Site diagram outputs using the shared Alwaldend identity preset."""

load("//projects/alwaldend.com/assets:defs.bzl", "site_stylesheet")
load("//tools/mermaid:defs.bzl", "mermaid_svg", "mermaid_webp")

_SITE_THEME = "//data/alwaldend:mermaid_theme"

# Preserve existing configuration and Sass inputs for spatial consumers.
SITE_MERMAID_THEME = struct(
    config = "//data/alwaldend/assets/mermaid:config",
    palette = "//projects/alwaldend.com/assets:mermaid_palette",
    font = "//data/alwaldend:diagram_font",
)

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
        **kwargs
    )
