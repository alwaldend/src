_BUILD_FILE = """\
load("{toolchain_bzl}", "git_toolchain")
load("{binary_bzl}", "git_binary")

toolchain(
    name = "git_toolchain",
    toolchain = ":git_toolchain_impl",
    toolchain_type = "{toolchain_type}",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "git_invalidation",
    srcs = ["release_git_state.bundle"],
    visibility = ["//visibility:public"],
)

filegroup(
    name = "release_git_state",
    srcs = ["release_git_state.bundle"],
    visibility = ["//visibility:public"],
)

git_toolchain(
    name = "git_toolchain_impl",
    git_path = "{git_path}",
    git_dir = "{git_dir}",
    git_root = "{git_root}",
    invalidation = [":git_invalidation"],
    visibility = ["//visibility:public"],
)

git_binary(
    name = "git",
    visibility = ["//visibility:public"],
)
"""

def _execute(ctx, arguments, description):
    result = ctx.execute(arguments, quiet = True)
    if result.return_code:
        fail("{} failed:\n{}\n{}".format(description, result.stdout, result.stderr))
    return result

def _resolve_git_path(ctx, git_path, name):
    result = _execute(
        ctx,
        [
            git_path,
            "-C",
            ctx.workspace_root,
            "rev-parse",
            "--path-format=absolute",
            "--git-path",
            name,
        ],
        "Resolving Git metadata path {}".format(name),
    )
    path = result.stdout.strip()
    if not path:
        fail("Git returned an empty metadata path for {}".format(name))
    return ctx.path(path)

def _watch_release_git_state(ctx, git_path):
    ctx.watch(ctx.workspace_root.get_child(".git"))
    ctx.watch(_resolve_git_path(ctx, git_path, "HEAD"))
    ctx.watch(_resolve_git_path(ctx, git_path, "packed-refs"))
    ctx.watch_tree(_resolve_git_path(ctx, git_path, "refs/tags"))

    symbolic_ref = ctx.execute(
        [git_path, "-C", ctx.workspace_root, "symbolic-ref", "-q", "HEAD"],
        quiet = True,
    )
    if symbolic_ref.return_code == 0:
        ref = symbolic_ref.stdout.strip()
        if not ref:
            fail("Git returned an empty symbolic HEAD")
        ctx.watch(_resolve_git_path(ctx, git_path, ref))
    elif symbolic_ref.return_code != 1:
        fail(
            "Resolving symbolic Git HEAD failed:\n{}\n{}".format(
                symbolic_ref.stdout,
                symbolic_ref.stderr,
            ),
        )

def _create_release_git_bundle(ctx, git_path):
    workspace_root = ctx.workspace_root
    head_result = _execute(
        ctx,
        [git_path, "-C", workspace_root, "rev-parse", "--verify", "HEAD^{commit}"],
        "Resolving the release Git revision",
    )
    head = head_result.stdout.strip()
    tags_result = _execute(
        ctx,
        [
            git_path,
            "-C",
            workspace_root,
            "tag",
            "--merged",
            head,
            "--sort=refname",
            "--format=%(objectname) %(refname)",
        ],
        "Resolving release Git tags",
    )

    bundle_refs = ["{} refs/heads/release".format(head)]
    revisions = [head]
    for line in tags_result.stdout.splitlines():
        fields = line.split(" ", 1)
        if len(fields) != 2 or not fields[1].startswith("refs/tags/"):
            fail("Unexpected release Git tag: {}".format(line))
        bundle_refs.append(line)
        revisions.append(fields[0])

    bundle_header = ctx.path("release_git_state.header")
    bundle_revisions = ctx.path("release_git_state.revisions")
    bundle = ctx.path("release_git_state.bundle")
    ctx.file(
        bundle_header,
        "# v2 git bundle\n{}\n\n".format("\n".join(bundle_refs)),
    )
    ctx.file(bundle_revisions, "{}\n".format("\n".join(revisions)))

    shell_path = ctx.which("sh")
    cat_path = ctx.which("cat")
    if not shell_path or not cat_path:
        fail("Could not find sh and cat")
    _execute(
        ctx,
        [
            shell_path,
            "-c",
            """\
set -eu
"$1" "$2" > "$4"
"$5" -C "$6" pack-objects \\
    --stdout \\
    --revs \\
    --no-use-bitmap-index \\
    --no-reuse-delta \\
    --no-reuse-object \\
    --threads=1 \\
    --window=10 \\
    --depth=50 \\
    --compression=9 \\
    --no-delta-islands \\
    --no-sparse \\
    --filter=blob:none \\
    < "$3" >> "$4"
""",
            "create_release_git_bundle",
            cat_path,
            bundle_header,
            bundle_revisions,
            bundle,
            git_path,
            workspace_root,
        ],
        "Creating deterministic release Git bundle",
    )
    ctx.delete(bundle_header)
    ctx.delete(bundle_revisions)

def _impl(ctx):
    git_path = ctx.which("git")
    if not git_path:
        fail("Could not find git")

    _watch_release_git_state(ctx, git_path)
    _create_release_git_bundle(ctx, git_path)
    git_dir = ctx.path(_execute(ctx, [git_path, "-C", ctx.workspace_root, "rev-parse", "--absolute-git-dir"], "Resolving Git directory").stdout.strip())
    git_root = ctx.path(_execute(ctx, [git_path, "-C", ctx.workspace_root, "rev-parse", "--show-toplevel"], "Resolving Git worktree").stdout.strip())
    ctx.file(
        "BUILD.bazel",
        _BUILD_FILE.format(
            toolchain_bzl = ctx.attr._toolchain_bzl,
            binary_bzl = ctx.attr._binary_bzl,
            toolchain_type = ctx.attr._toolchain_type,
            git_path = git_path,
            git_dir = git_dir,
            git_root = git_root,
        ),
    )

git_repo = repository_rule(
    implementation = _impl,
    doc = "Capture a local Git checkout and deterministic release history bundle.",
    attrs = {
        "_toolchain_bzl": attr.label(default = Label("//rules:git_toolchain.bzl")),
        "_binary_bzl": attr.label(default = Label("//rules:git_binary.bzl")),
        "_toolchain_type": attr.label(default = Label("//rules:git_toolchain_type")),
    },
    local = True,
)
