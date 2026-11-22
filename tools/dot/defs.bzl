"""Build actions for exporting and rendering Graphviz DOT diagrams."""

def _dot_export_impl(ctx):
    output = ctx.outputs.out
    if not output.basename.endswith(".dot"):
        fail("out must name a .dot file")

    args = ctx.actions.args()
    args.add("export")
    args.add(output)
    args.add(ctx.executable.command)
    args.add_all(ctx.attr.arguments)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.attr._renderer[DefaultInfo].files_to_run,
        mnemonic = "DotExport",
        outputs = [output],
        progress_message = "Exporting DOT %{label}",
        tools = [ctx.attr.command[DefaultInfo].files_to_run],
    )
    return [DefaultInfo(files = depset([output]))]

def _dot_svg_impl(ctx):
    output = ctx.outputs.out
    if not output.basename.endswith(".svg"):
        fail("out must name an .svg file")

    inputs = [ctx.file.src, ctx.file._theme, ctx.file._font, ctx.file._font_license]
    args = ctx.actions.args()
    args.add("svg")
    args.add(output)
    args.add(ctx.executable._dot)
    args.add_all(inputs)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.attr._renderer[DefaultInfo].files_to_run,
        inputs = inputs,
        mnemonic = "DotSvg",
        outputs = [output],
        progress_message = "Rendering DOT as SVG %{label}",
        tools = [ctx.attr._dot[DefaultInfo].files_to_run],
    )
    return [DefaultInfo(files = depset([output]))]

_renderer = attr.label(
    cfg = "exec",
    default = Label("//tools/dot:renderer"),
    executable = True,
)

dot_export = rule(
    implementation = _dot_export_impl,
    doc = "Runs an executable that prints DOT and captures its standard output.",
    attrs = {
        "arguments": attr.string_list(doc = "Arguments to pass to the command."),
        "command": attr.label(
            cfg = "exec",
            doc = "Executable that prints DOT to standard output.",
            executable = True,
            mandatory = True,
        ),
        "out": attr.output(doc = "Generated .dot file.", mandatory = True),
        "_renderer": _renderer,
    },
)

dot_svg = rule(
    implementation = _dot_svg_impl,
    doc = "Renders a DOT file as a self-contained SVG using the common theme.",
    attrs = {
        "out": attr.output(doc = "Rendered .svg file.", mandatory = True),
        "src": attr.label(
            allow_single_file = [".dot"],
            doc = "DOT source file or target that generates it.",
            mandatory = True,
        ),
        "_dot": attr.label(
            cfg = "exec",
            default = Label("//third_party/org_graphviz_graphviz:dot"),
            executable = True,
        ),
        "_font": attr.label(
            allow_single_file = [".ttf"],
            default = Label("//internal/assets:diagram_font"),
        ),
        "_font_license": attr.label(
            allow_single_file = True,
            default = Label("//internal/assets:diagram_font_license"),
        ),
        "_renderer": _renderer,
        "_theme": attr.label(
            allow_single_file = [".json"],
            default = Label("//internal/assets:theme.json"),
        ),
    },
)
