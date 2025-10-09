def _impl(ctx):
    args = ctx.actions.args()
    for arg in ctx.attr.arguments:
        args.add(ctx.expand_location(arg))
    ctx.actions.run(
        outputs = ctx.outputs.outs,
        inputs = ctx.files.srcs,
        executable = ctx.executable.git,
        arguments = [args],
        progress_message = "Running git action: %{label}",
    )
    return [
        DefaultInfo(
            files = depset(ctx.outputs.outs),
        ),
    ]

al_git_run_binary = rule(
    implementation = _impl,
    doc = "Run a git binary as a build action",
    attrs = {
        "arguments": attr.string_list(
            default = [],
            doc = "Arguments (Location is expanded)",
        ),
        "srcs": attr.label_list(
            default = [],
            doc = "Files to be made available",
        ),
        "outs": attr.output_list(
            mandatory = True,
            doc = "Outputs",
        ),
        "git": attr.label(
            mandatory = True,
            executable = True,
            cfg = "exec",
            doc = "Git binary",
        ),
    },
)
