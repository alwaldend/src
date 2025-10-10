_BUILD_FILE = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files", "strip_prefix", "REMOVE_BASE_DIRECTORY")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@com_alwaldend_src//bzl/rules/git:al_git_library.bzl", "al_git_library")
load("@com_alwaldend_src//bzl/rules/git:al_git_binary.bzl", "al_git_binary")
load("@com_alwaldend_src//bzl/rules/git:al_git_run_binary.bzl", "al_git_run_binary")

al_git_binary(
    name = "git",
    src = ":git_lib",
    visibility = ["//visibility:public"],
)

al_git_run_binary(
    name = "git_source_archive",
    git = ":git",
    outs = ["git_source_archive.tar"],
    arguments = [
        "archive",
        "--output",
        "$(execpath :git_source_archive.tar)",
        "HEAD",
    ],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "git_archive",
    srcs = [":git_files"],
    visibility = ["//visibility:public"],
)

pkg_files(
    name = "git_files",
    srcs = glob([".git/**/*"]),
    visibility = ["//visibility:public"],
    strip_prefix = strip_prefix.from_pkg(),
)

al_git_library(
    name = "git_lib",
    git_archive = ":git_archive",
    visibility = ["//visibility:public"],
)
"""

def _impl(rctx):
    rctx.symlink(rctx.workspace_root.get_child(".git"), ".git")
    rctx.file("BUILD.bazel", _BUILD_FILE)

al_git_repo = repository_rule(
    implementation = _impl,
    doc = "Repository rule that exports .git directory",
)
