def _impl(ctx):
    env = {
        "ORAS_BINARY": ctx.executable.oras.path,
    }
    runfiles = ctx.runfiles(files = ctx.files.oras)
    default_info = DefaultInfo(
        files = depset(ctx.files.oras),
        runfiles = runfiles,
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
            oras = ctx.executable.oras,
        ),
    ]

al_oras_toolchain = rule(
    doc = "oras toolchain",
    implementation = _impl,
    attrs = {
        "oras": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "oras binary",
        ),
    },
)
