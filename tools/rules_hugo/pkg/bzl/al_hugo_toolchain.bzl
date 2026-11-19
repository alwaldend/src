def _impl(ctx):
    env = {
        "HUGO": ctx.executable.hugo.path,
    }
    runfiles = ctx.runfiles(files = ctx.files.hugo)
    runfiles = runfiles.merge(ctx.attr.hugo[DefaultInfo].default_runfiles)
    hugo_files = [
        file
        for file in ctx.attr.hugo[DefaultInfo].files.to_list()
        if file.is_source
    ]
    runfiles = runfiles.merge(ctx.runfiles(files = hugo_files))
    default_info = DefaultInfo(
        files = depset(ctx.files.hugo + hugo_files),
        runfiles = runfiles,
    )
    env_file = ctx.actions.declare_file("{}.env.txt".format(ctx.label.name))
    ctx.actions.run_shell(
        outputs = [env_file],
        tools = [ctx.executable.hugo],
        command = "'{}' env >'{}'".format(ctx.executable.hugo.path, env_file.path),
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            default_info = default_info,
            env_file = env_file,
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
