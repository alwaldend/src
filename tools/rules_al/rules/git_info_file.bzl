load("//rules:git_current_state.bzl", "GIT_CURRENT_STATE")

def _impl(ctx):
    git_info = ctx.actions.declare_file("{}.git_info.json".format(ctx.label.name))
    git_toolchain = ctx.toolchains[Label("//rules:git_toolchain_type")]
    args = ctx.actions.args()
    args.add("release")
    args.add("git-info")
    args.add_all(["--git_root", git_toolchain.git_root])
    args.add_all(["--output", git_info])
    args.add_all(["--timeout", "{}s".format(ctx.attr.timeout)])
    ctx.actions.run(
        executable = ctx.executable.git_tool,
        arguments = [args],
        inputs = ctx.files.git_state + git_toolchain.git_invalidation,
        outputs = [git_info],
        execution_requirements = {"no-remote": "1"},
    )

    return [
        DefaultInfo(
            files = depset([git_info]),
        ),
    ]

git_info_file = rule(
    doc = "Generate git info file",
    toolchains = [Label("//rules:git_toolchain_type")],
    implementation = _impl,
    attrs = {
        "git_state": attr.label_list(
            default = GIT_CURRENT_STATE,
            doc = "Files that should invalidate the cache on new commit",
        ),
        "timeout": attr.int(
            default = 60,
            doc = "Timeout in seconds",
        ),
        "git_tool": attr.label(
            default = Label("//cmd/al:al"),
            cfg = "exec",
            executable = True,
            doc = "Git tool to use",
        ),
    },
)
