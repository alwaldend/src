def _impl(ctx):
    env = {
        "DRAWIO": ctx.executable.drawio.path,
    }

    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            drawio = ctx.executable.drawio,
        ),
    ]

al_drawio_toolchain = rule(
    doc = "Drawio toolchain",
    implementation = _impl,
    attrs = {
        "drawio": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "Drawio binary",
        ),
    },
)
