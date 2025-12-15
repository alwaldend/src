def _impl(ctx):
    env = {
        "TRUFFLEHOG_BINARY": ctx.executable.trufflehog.path,
    }
    runfiles = ctx.runfiles(files = ctx.files.trufflehog)
    default_info = DefaultInfo(
        files = depset(ctx.files.trufflehog),
        runfiles = runfiles,
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
            trufflehog = ctx.executable.trufflehog,
        ),
    ]

al_trufflehog_toolchain = rule(
    doc = "Trufflehog toolchain",
    implementation = _impl,
    attrs = {
        "trufflehog": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "Trufflehog binary",
        ),
    },
)
