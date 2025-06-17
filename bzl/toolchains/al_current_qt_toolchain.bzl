load("@bazel_skylib//lib:paths.bzl", "paths")

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
        "_toolchain": attr.string(default = "//bzl/toolchain-types:al-qt"),
    },
    toolchains = [
        "//bzl/toolchain-types:al-qt",
    ],
)
