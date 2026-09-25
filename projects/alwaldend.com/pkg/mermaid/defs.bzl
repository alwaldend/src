"""Native Mermaid diagrams using the site's shared color palette."""

load("//tools/mermaid:defs.bzl", "mermaid_svg")

def _mermaid_palette_impl(ctx):
    ctx.actions.run(
        executable = ctx.executable._sass,
        arguments = ["--no-source-map", "--style=compressed", ctx.file.src.path, ctx.outputs.out.path],
        inputs = [ctx.file.src] + ctx.files.deps,
        outputs = [ctx.outputs.out],
        tools = [ctx.attr._sass[DefaultInfo].files_to_run],
        mnemonic = "MermaidPalette",
    )
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

mermaid_palette = rule(
    implementation = _mermaid_palette_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".scss"], mandatory = True),
        "deps": attr.label_list(allow_files = [".scss"]),
        "out": attr.output(mandatory = True),
        "_sass": attr.label(
            default = "@com_github_sass_dart_sass//:sass_binary",
            executable = True,
            cfg = "exec",
        ),
    },
)

def mermaid_site_svg(name, src, out, **kwargs):
    """Render a native SVG and its .dark.svg sibling for site image embeds."""
    mermaid_svg(
        name = name,
        src = src,
        out = out,
        dark_out = out.removesuffix(".svg") + ".dark.svg",
        native = True,
        theme = "//projects/alwaldend.com/pkg/mermaid:theme.json",
        palette = "//projects/alwaldend.com/assets:mermaid_palette",
        label_font = "@com_alwaldend_src_tools_drawio_fonts//:LiberationSans-Regular.ttf",
        **kwargs
    )
