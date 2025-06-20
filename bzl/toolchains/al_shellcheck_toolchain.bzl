def _impl(ctx):
    env = {
        "SHELLCHECK": ctx.executable.shellcheck.path,
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            shellcheck_target = ctx.attr.shellcheck,
            shellcheck = ctx.executable.shellcheck,
        ),
    ]

al_shellcheck_toolchain = rule(
    doc = "Shellcheck toolchain",
    implementation = _impl,
    attrs = {
        "shellcheck": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "Shellcheck binary",
        ),
    },
)
