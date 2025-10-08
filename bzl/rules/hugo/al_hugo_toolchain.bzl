def _impl(ctx):
    env = {
        "HUGO": ctx.executable.hugo.path,
    }
    default_info = DefaultInfo(
        files = depset(ctx.files.hugo),
        runfiles = ctx.runfiles(files = ctx.files.hugo),
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
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
