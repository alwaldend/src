def _impl(ctx):
    env = {
        "HUGO": ctx.executable.hugo.path,
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            hugo_target = ctx.attr.hugo,
            hugo = ctx.executable.hugo,
        ),
    ]

al_hugo_toolchain = rule(
    doc = "Hugo toolchain",
    implementation = _impl,
    attrs = {
        "hugo": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "Hugo binary",
        ),
    },
)
