def _impl(ctx):
    git_info = ctx.actions.declare_file("{}.git_info.json".format(ctx.label.name))
    git_toolchain = ctx.toolchains["//tools/git/main/bzl:toolchain_type"]
    args = ctx.actions.args()
    args.add_all(["git_info"])
    args.add_all(["--git_root", git_toolchain.git_root])
    args.add_all(["--output", git_info])
    args.add_all(["--timeout", ctx.attr.timeout])
    ctx.actions.run(
        executable = ctx.executable.git_tool,
        arguments = [args],
        inputs = ctx.files.git_state,
        outputs = [git_info],
    )

    return [
        DefaultInfo(
            files = depset([git_info]),
        ),
    ]

al_git_info_file = rule(
    doc = "Generate git info file",
    toolchains = ["//tools/git/main/bzl:toolchain_type"],
    implementation = _impl,
    attrs = {
        "git_state": attr.label(
            default = "@com_alwaldend_src_tools_git//:git_state",
            doc = "Files that should invalidate the cache on new commit",
        ),
        "timeout": attr.int(
            default = 60,
            doc = "Timeout in seconds",
        ),
        "git_tool": attr.label(
            default = "//tools/git/main/go",
            cfg = "exec",
            executable = True,
            doc = "Git tool to use",
        ),
    },
)
