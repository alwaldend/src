def _impl(ctx):
    env = {}
    out = ctx.actions.declare_file("{}.bin/al".format(ctx.label.name))
    ctx.actions.symlink(
        target_file = ctx.executable.al,
        output = out,
        is_executable = True,
    )
    runfiles = ctx.runfiles(files = ctx.files.al)
    runfiles = runfiles.merge(ctx.attr.al[DefaultInfo].default_runfiles)
    default_info = DefaultInfo(
        files = depset([out]),
        runfiles = runfiles,
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            default_info = default_info,
            executable = out,
            env = env,
        ),
    ]

al_toolchain = rule(
    doc = "Al toolchain",
    implementation = _impl,
    attrs = {
        "al": attr.label(
            mandatory = True,
            executable = True,
            cfg = "exec",
            doc = "Al binary",
        ),
    },
)
