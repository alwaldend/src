def _impl(ctx):
    env = {
        "MINISIGN_PATH": ctx.executable.minisign.path,
    }
    default_info = DefaultInfo(
        files = depset([ctx.executable.minisign]),
        runfiles = ctx.runfiles(files = [ctx.executable.minisign]),
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
            minisign = ctx.executable.minisign,
        ),
    ]

al_minisign_toolchain = rule(
    doc = "Minisign toolchain",
    implementation = _impl,
    attrs = {
        "minisign": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "Minisign binary",
        ),
    },
)
