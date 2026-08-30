"""Build actions for rendering Mermaid source diagrams."""

def _mermaid_svg_impl(ctx):
    source = ctx.file.src
    output = ctx.outputs.out

    if not output.basename.endswith(".svg"):
        fail("out must name an .svg file, got: {}".format(output.basename))

    args = ctx.actions.args()
    args.add("-i")
    args.add(source.path)
    args.add("-o")
    args.add(output.path)
    args.add("-b")
    args.add("white")

    ctx.actions.run(
        arguments = [args],
        env = {
            "BAZEL_BINDIR": ctx.bin_dir.path,
            # Declared File.path values are relative to the action execroot.
            "JS_BINARY__NO_CD_BINDIR": "1",
            "PUPPETEER_EXECUTABLE_PATH": ctx.executable._browser.path,
        },
        executable = ctx.executable._mmdc,
        inputs = [source],
        mnemonic = "MermaidSvg",
        outputs = [output],
        progress_message = "Rendering Mermaid SVG %{label}",
        tools = [
            ctx.attr._browser[DefaultInfo].files_to_run,
            ctx.attr._mmdc[DefaultInfo].files_to_run,
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
        "_browser": attr.label(
            cfg = "exec",
            default = "@com_alwaldend_src_tools_mermaid//:chrome_headless_shell_binary",
            executable = True,
        ),
        "_mmdc": attr.label(
            cfg = "exec",
            default = "//tools/mermaid/cmd/mmdc:mmdc_raw",
            executable = True,
        ),
    },
)
