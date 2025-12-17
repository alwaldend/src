_BUILD_FILE = """\
load("@com_alwaldend_src//tools/git/main/bzl:al_git_toolchain.bzl", "al_git_toolchain")
load("@com_alwaldend_src//tools/git/main/bzl:al_git_binary.bzl", "al_git_binary")

toolchain(
    name = "git_toolchain",
    toolchain = ":git_toolchain_impl",
    toolchain_type = "@com_alwaldend_src//tools/git/main/bzl:toolchain_type",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "git_objects",
    srcs = glob([".git/objects/**"]),
)

filegroup(
    name = "git_refs",
    srcs = glob([".git/refs/**"]),
)

filegroup(
    name = "git_head",
    srcs = [".git/HEAD"],
)

filegroup(
    name = "git_invalidation",
    srcs = [":git_objects", ":git_refs", ":git_head"] + glob([".git/**"]),
    visibility = ["//visibility:public"],
)

al_git_toolchain(
    name = "git_toolchain_impl",
    git_path = "{git_path}",
    git_dir = "{git_dir}",
    git_root = "{git_root}",
    invalidation = [":git_invalidation"],
    visibility = ["//visibility:public"],
)

al_git_binary(
    name = "git",
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    git_path = ctx.which("git")
    git_dir = ctx.workspace_root.get_child(".git")
    ctx.symlink(git_dir, ".git")
    git_root = ctx.workspace_root
    ctx.file(
        "BUILD.bazel",
        _BUILD_FILE.format(
            git_path = git_path,
            git_dir = git_dir,
            git_root = git_root,
        ),
    )

al_git_repo = repository_rule(
    implementation = _impl,
    doc = "Git repo",
)
