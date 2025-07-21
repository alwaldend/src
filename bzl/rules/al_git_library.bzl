load("//bzl/providers:al_git_info.bzl", "AlGitInfo")

def _impl(ctx):
    dir = ctx.actions.declare_directory("{}-dir".format(ctx.label.name))
    files = [dir]
    git = ctx.toolchains["//bzl/toolchain-types:git"]
    runfiles = ctx.runfiles(files = files)

    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [dir],
        command = """
            tar -xf '{src}' -C '{output}' && \
                '{git}' -C '{output}' reset --hard
        """.format(src = ctx.file.src.path, output = dir.path, git = git.git_bin),
    )

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlGitInfo(
            git_dir = ctx.file.src,
        ),
    ]

al_git_library = rule(
    implementation = _impl,
    doc = "Define git information",
    toolchains = ["//bzl/toolchain-types:git"],
    provides = [AlGitInfo],
    attrs = {
        "src": attr.label(
            mandatory = True,
            allow_single_file = [".tar"],
            doc = "Git dir",
        ),
    },
)
