load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")
load(":al_git_library_info.bzl", "AlGitLibraryInfo")

def _impl(ctx):
    runfiles = ctx.runfiles(files = [ctx.file.git_archive])
    return [
        DefaultInfo(
            files = depset([ctx.file.git_archive]),
            runfiles = runfiles,
        ),
        AlGitLibraryInfo(
            git_archive = ctx.file.git_archive,
        ),
    ]

al_git_library = rule(
    implementation = _impl,
    doc = "Define git information",
    toolchains = ["//tools/git:toolchain_type"],
    provides = [AlGitLibraryInfo],
    attrs = {
        "git_archive": attr.label(
            mandatory = True,
            allow_single_file = [".tar"],
            doc = "Git directory",
        ),
    },
)
