load("@bazel_skylib//lib:paths.bzl", "paths")

def _current_qt_toolchain_impl(ctx):
    toolchain = ctx.toolchains[ctx.attr._toolchain]
    return [
        toolchain,
        platform_common.TemplateVariableInfo(toolchain.env),
    ]

def _preinstalled_qt_toolchain_impl(ctx):
    env = {
        "QMAKE": paths.join(ctx.attr.dir, ctx.attr.version, ctx.attr.platform, "bin", "qmake"),
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(env = env),
    ]

def register_qt_toolchains():
    """Register qt toolchains"""
    native.register_toolchains(
        "//bazel/qt:qt-toolchain",
    )

current_qt_toolchain = rule(
    doc = "Get current selected qt toolchain",
    implementation = _current_qt_toolchain_impl,
    attrs = {
        "_toolchain": attr.string(default = str(Label("//bazel/qt:qt-toolchain"))),
    },
    toolchains = [
        str(Label("//bazel/qt:qt-toolchain")),
    ],
)

preinstalled_qt_toolchain = rule(
    doc = "Preinstalled qt toolchain",
    implementation = _preinstalled_qt_toolchain_impl,
    attrs = {
        "dir": attr.string(doc = "Root qt directory", default = "/opt/qt"),
        "version": attr.string(doc = "Qt version", default = "6.9.0"),
        "platform": attr.string(doc = "Qt platform", default = "gcc_64"),
    },
)
