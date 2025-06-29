def _impl(ctx):
    drawio = ctx.toolchains["//bzl/toolchain-types:drawio"]
    ctx.run(
        executable = drawio.drawio,
        inputs = [ctx.file.src],
    )

al_drawio_diagrm = rule(
    implementation = _impl,
    doc = "Build a drawio diagram",
    toolchains = [
        "//bzl/toolchain-types:drawio",
    ],
    attrs = {
        "src": attr.label(allow_single_file = [".drawio"], doc = "Drawio diagram"),
        "out": attr.output(doc = "Drawio output"),
    },
)
