load("@bazel_skylib//lib:paths.bzl", "paths")

def _preinstalled_qt_toolchain_impl(ctx):
    env = {
        "QMAKE": paths.join(ctx.attr.dir, ctx.attr.version, ctx.attr.platform, "bin", "qmake"),
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(env = env),
    ]

al_preinstalled_qt_toolchain = rule(
    doc = "Preinstalled qt toolchain",
    implementation = _preinstalled_qt_toolchain_impl,
    attrs = {
        "dir": attr.string(doc = "Root qt directory", default = "/opt/qt"),
        "version": attr.string(doc = "Qt version", default = "6.9.0"),
        "platform": attr.string(doc = "Qt platform", default = "gcc_64"),
    },
)
