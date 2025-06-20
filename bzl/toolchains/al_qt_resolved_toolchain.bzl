def _impl(ctx):
    toolchain = ctx.toolchains["//bzl/toolchain-types:qt"]
    return [
        toolchain,
        platform_common.TemplateVariableInfo(toolchain.env),
    ]

al_qt_resolved_toolchain = rule(
    doc = "Resolved qt toolchain",
    implementation = _impl,
    toolchains = ["//bzl/toolchain-types:qt"],
)
