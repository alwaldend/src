load("@bazel_skylib//lib:paths.bzl", "paths")

def _impl(ctx):
    env = {
        "QT_QMAKE": paths.join(ctx.attr.dir, ctx.attr.version, ctx.attr.platform, "bin", "qmake"),
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            dir = ctx.attr.dir,
            version = ctx.attr.version,
            platform = ctx.attr.platform,
        ),
    ]

al_qt_local_toolchain = rule(
    doc = "Local qt toolchain",
    implementation = _impl,
    attrs = {
        "dir": attr.string(doc = "Root qt directory", mandatory = True),
        "version": attr.string(doc = "Qt version", mandatory = True),
        "platform": attr.string(doc = "Qt platform", mandatory = True),
    },
)
