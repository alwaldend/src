load("@bazel_skylib//lib:paths.bzl", "paths")

def _drawio_toolchain_impl(ctx):
    env = {
        "DRAWIO": ctx.attr.compiler_path,
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(env = env),
    ]

al_drawio_toolchain = rule(
    doc = "Drawio toolchain",
    implementation = _drawio_toolchain_impl,
    attrs = {
        "compiler_path": attr.string(doc = "Compiler path", default = "~/.local/bin/drawio"),
    },
)
