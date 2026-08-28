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
    name = "git_invalidation",
    srcs = [".git"],
    visibility = ["//visibility:public"],
)

filegroup(
    name = "release_git_state",
    srcs = ["release_git_repository"],
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

def _execute(ctx, arguments, description):
    result = ctx.execute(arguments, quiet = True)
    if result.return_code:
        fail("{} failed:\n{}\n{}".format(description, result.stdout, result.stderr))

def _impl(ctx):
    git_path = ctx.which("git")
    if not git_path:
        fail("Could not find git")

    release_git_repository = ctx.path("release_git_repository")
    ctx.delete(release_git_repository)
    _execute(
        ctx,
        [git_path, "init", "--bare", "--quiet", release_git_repository],
        "Initializing release Git repository",
    )
    _execute(
        ctx,
        [
            git_path,
            "--git-dir={}".format(release_git_repository),
            "fetch",
            "--force",
            "--no-recurse-submodules",
            "--no-write-fetch-head",
            "--quiet",
            ctx.workspace_root,
            "+HEAD:refs/heads/release",
        ],
        "Fetching release Git state",
    )
    _execute(
        ctx,
        [
            git_path,
            "--git-dir={}".format(release_git_repository),
            "symbolic-ref",
            "HEAD",
            "refs/heads/release",
        ],
        "Setting release Git revision",
    )
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
    local = True,
)
