def _impl(ctx):
    exec = ctx.actions.declare_file("{}.shell.sh".format(ctx.label.name))
    ctx.actions.symlink(target_file = ctx.executable.dev_shell, output = exec)
    runfiles = ctx.runfiles()
    for tool_name, tool in ctx.attr.tools.items():
        tool_exec = tool[DefaultInfo].files_to_run.executable
        if not tool_exec:
            fail("Tool {} is missing an executable: {}".format(tool_name, tool))
        runfiles = runfiles.merge_all(
            [
                ctx.runfiles(symlinks = {"bin/{}".format(tool_name): tool_exec}),
                tool[DefaultInfo].default_runfiles,
            ],
        )

    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = exec,
        ),
    ]

al_dev_shell_binary = rule(
    doc = "Dev shell binary",
    implementation = _impl,
    executable = True,
    attrs = {
        "tools": attr.string_keyed_label_dict(
            cfg = "exec",
            doc = "Tools, keys are tool names, values are tool labels",
        ),
        "dev_shell": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/dev_shell:dev_shell_script",
            doc = "Dev shell binary",
        ),
    },
)
