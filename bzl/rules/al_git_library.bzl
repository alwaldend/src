load("//bzl/providers:al_git_info.bzl", "AlGitInfo")

def _impl(ctx):
    git = ctx.toolchains["//bzl/toolchain-types:git"]
    files = [ctx.file.src]
    runfiles = ctx.runfiles(files = files)

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlGitInfo(
            archive = ctx.file.src,
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
