"""Bundled configuration, colors, and typography for Mermaid renderers."""

MermaidThemeInfo = provider(
    doc = "Concrete files defining a Mermaid theme and its licensed font.",
    fields = {
        "config": "Mermaid configuration JSON File.",
        "palette": "JSON File mapping light and dark to Mermaid theme-variable colors.",
        "font": "TrueType font File matching the configuration's font family.",
        "font_license": "License File for the bundled font.",
    },
)

def _mermaid_theme_impl(ctx):
    files = [ctx.file.config, ctx.file.palette, ctx.file.font, ctx.file.font_license]
    return [
        MermaidThemeInfo(
            config = ctx.file.config,
            palette = ctx.file.palette,
            font = ctx.file.font,
            font_license = ctx.file.font_license,
        ),
        DefaultInfo(files = depset(files)),
    ]

mermaid_theme = rule(
    implementation = _mermaid_theme_impl,
    doc = "Bundles a Mermaid configuration, light/dark palette, font, and license.",
    attrs = {
        "config": attr.label(
            allow_single_file = [".json"],
            mandatory = True,
            doc = "Mermaid configuration.",
        ),
        "palette": attr.label(
            allow_single_file = [".json"],
            mandatory = True,
            doc = "JSON object mapping light and dark to Mermaid theme-variable colors.",
        ),
        "font": attr.label(
            allow_single_file = [".ttf"],
            mandatory = True,
            doc = "Font matching the configured font family.",
        ),
        "font_license": attr.label(
            allow_single_file = True,
            mandatory = True,
            doc = "License belonging to font.",
        ),
    },
)
