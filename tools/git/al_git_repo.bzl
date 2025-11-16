_BUILD_FILE = """\
load("@com_alwaldend_src//tools/git:al_git_toolchain.bzl", "al_git_toolchain")
load("@com_alwaldend_src//tools/git:al_git_binary.bzl", "al_git_binary")

toolchain(
    name = "git_toolchain",
    toolchain = ":git_toolchain_impl",
    toolchain_type = "@com_alwaldend_src//tools/git:toolchain_type",
    visibility = ["//visibility:public"],
)

al_git_toolchain(
    name = "git_toolchain_impl",
    git_path = "{git_path}",
    git_dir = "{git_dir}",
    git_root = "{git_root}",
    visibility = ["//visibility:public"],
)

al_git_binary(
    name = "git",
    visibility = ["//visibility:public"],
)

genrule(
    name = "current_rev",
    outs = ["current_rev.txt"],
    cmd = "$(execpath :git) rev-parse HEAD >$(@)",
    tags = ["no-cache"],
    visibility = ["//visibility:public"],
    tools = [":git"],
)
"""

def _impl(ctx):
    git_path = ctx.which("git")
    git_dir = ctx.workspace_root.get_child(".git")
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
