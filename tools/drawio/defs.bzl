"""Sandboxed Drawio webapp extraction and browser rendering."""

def _webapp_impl(ctx):
    output = ctx.actions.declare_directory(ctx.label.name + ".webapp")
    args = ctx.actions.args()
    args.add("--drawio", ctx.executable._drawio.path)
    args.add("--out", output.path)
    ctx.actions.run(
        executable = ctx.executable._extract,
        arguments = [args],
        outputs = [output],
        tools = [ctx.attr._drawio[DefaultInfo].files_to_run],
        mnemonic = "DrawioWebapp",
    )
    return [DefaultInfo(files = depset([output]))]

_drawio_webapp = rule(
    implementation = _webapp_impl,
    attrs = {
        "_drawio": attr.label(default = "@com_drawio_desktop_bin//:drawio_binary", executable = True, cfg = "exec"),
        "_extract": attr.label(default = "//tools/drawio/cmd/extract", executable = True, cfg = "exec"),
    },
)

def _svg_impl(ctx):
    outputs = []
    manifest = {}
    for index, page in enumerate(ctx.attr.pages):
        basename = ctx.attr.pages[page]
        if "/" in basename or not basename.endswith(".svg"):
            fail("Each output must be a simple SVG basename")
        output = ctx.outputs.outs[index]
        outputs.append(output)
        manifest[page] = output.path
    manifest_file = ctx.actions.declare_file(ctx.label.name + ".json")
    ctx.actions.write(manifest_file, json.encode(manifest))
    ctx.actions.run(
        executable = ctx.executable._render,
        arguments = [ctx.file.webapp.path, ctx.executable._browser.path, ctx.file.src.path, manifest_file.path, ctx.file._font_anchor.path],
        inputs = [ctx.file.webapp, ctx.file.src, manifest_file] + ctx.files._fonts,
        outputs = outputs,
        tools = [ctx.attr._browser[DefaultInfo].files_to_run],
        env = {"BAZEL_BINDIR": ctx.bin_dir.path, "JS_BINARY__NO_CD_BINDIR": "1"},
        mnemonic = "DrawioSvg",
        progress_message = "Rendering Drawio pages %{label}",
    )
    return [DefaultInfo(files = depset(outputs))]

_drawio_svg = rule(
    implementation = _svg_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".drawio"], mandatory = True),
        "pages": attr.string_dict(mandatory = True),
        "outs": attr.output_list(mandatory = True),
        "webapp": attr.label(default = "//tools/drawio:webapp", allow_single_file = True),
        "_font_anchor": attr.label(default = "@com_alwaldend_src_tools_drawio_fonts//:LiberationSans-Regular.ttf", allow_single_file = True),
        "_fonts": attr.label(default = "@com_alwaldend_src_tools_drawio_fonts//:fonts"),
        "_render": attr.label(default = "//tools/drawio/cmd/render", executable = True, cfg = "exec"),
        "_browser": attr.label(default = "@com_alwaldend_src_tools_mermaid//:chrome_headless_shell_binary", executable = True, cfg = "exec"),
    },
)

def drawio_webapp(name, **kwargs):
    _drawio_webapp(name = name, **kwargs)

def drawio_svg(name, pages, **kwargs):
    _drawio_svg(name = name, pages = pages, outs = [name + "/" + basename for basename in pages.values()], **kwargs)
