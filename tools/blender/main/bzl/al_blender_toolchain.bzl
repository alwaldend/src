def _impl(ctx):
    env = {
        "BLENDER_BINARY": ctx.executable.blender.path,
    }
    runfiles = ctx.runfiles(files = ctx.files.blender)
    default_info = DefaultInfo(
        files = depset(ctx.files.blender),
        runfiles = runfiles,
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
            blender = ctx.executable.blender,
        ),
    ]

al_blender_toolchain = rule(
    doc = "blender toolchain",
    implementation = _impl,
    attrs = {
        "blender": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "blender binary",
        ),
    },
)
