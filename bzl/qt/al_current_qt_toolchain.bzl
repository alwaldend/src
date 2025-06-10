load("@bazel_skylib//lib:paths.bzl", "paths")
load("//bzl/vars:toolchain_types.bzl", "TOOLCHAIN_TYPES")

def _current_qt_toolchain_impl(ctx):
    toolchain = ctx.toolchains[ctx.attr._toolchain]
    return [
        toolchain,
        platform_common.TemplateVariableInfo(toolchain.env),
    ]

al_current_qt_toolchain = rule(
    doc = "Get current selected qt toolchain",
    implementation = _current_qt_toolchain_impl,
    attrs = {
        "_toolchain": attr.string(default = str(Label(TOOLCHAIN_TYPES.qt))),
    },
    toolchains = [
        str(Label(TOOLCHAIN_TYPES.qt)),
    ],
)
