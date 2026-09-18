"""Build actions for rendering Mermaid source diagrams."""

def _mermaid_svg_impl(ctx):
    source = ctx.file.src
    theme = ctx.file.theme
    output = ctx.outputs.out

    if not output.basename.endswith(".svg"):
        fail("out must name an .svg file, got: {}".format(output.basename))

    args = ctx.actions.args()
    args.add("--input", source.path)
    args.add("--output", output.path)
    args.add("--config", theme.path)
    args.add("--browser", ctx.executable._browser.path)
    args.add("--font-anchor", ctx.file._font_anchor.path)
    args.add("--background", "white")

    font_dependencies = ctx.files._fonts

    ctx.actions.run(
        arguments = [args],
        env = {
            "BAZEL_BINDIR": ctx.bin_dir.path,
            # Declared File.path values are relative to the action execroot.
            "JS_BINARY__NO_CD_BINDIR": "1",
        },
        executable = ctx.executable._render,
        inputs = [source, theme, ctx.file._font_anchor] + font_dependencies,
        mnemonic = "MermaidSvg",
        outputs = [output],
        progress_message = "Rendering Mermaid SVG %{label}",
        tools = [
            ctx.attr._browser[DefaultInfo].files_to_run,
            ctx.attr._render[DefaultInfo].files_to_run,
        ],
    )

    return [DefaultInfo(files = depset([output]))]

mermaid_svg = rule(
    implementation = _mermaid_svg_impl,
    doc = "Renders one Mermaid source file to a declared SVG output.",
    attrs = {
        "src": attr.label(
            allow_single_file = [".mmd"],
            mandatory = True,
            doc = "Mermaid source file.",
        ),
        "out": attr.output(
            mandatory = True,
            doc = "Rendered .svg output.",
        ),
        "theme": attr.label(
            allow_single_file = [".json"],
            default = "//tools/mermaid:theme.json",
            doc = "Mermaid configuration applied by default; a diagram's own " +
                  "Mermaid directives still override it.",
        ),
        "_font_anchor": attr.label(
            allow_single_file = [".ttf"],
            default = "@com_alwaldend_src_tools_drawio_fonts//:LiberationSans-Regular.ttf",
        ),
        "_fonts": attr.label(
            default = "@com_alwaldend_src_tools_drawio_fonts//:fonts",
        ),
        "_browser": attr.label(
            cfg = "exec",
            default = "@com_alwaldend_src_tools_mermaid//:chrome_headless_shell_binary",
            executable = True,
        ),
        "_render": attr.label(
            cfg = "exec",
            default = "//tools/mermaid/cmd/render",
            executable = True,
        ),
    },
)
