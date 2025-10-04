load("@bazel_skylib//lib:shell.bzl", "shell")
load(":al_git_info.bzl", "AlGitInfo")

def _impl(ctx):
    git_toolchain = ctx.toolchains["//bzl/rules/git:toolchain_type"]
    git_info = ctx.attr.src[AlGitInfo]
    script = ctx.actions.declare_file("{}-script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = ctx.files.src)
    runfiles.merge(ctx.attr.src[DefaultInfo].default_runfiles)

    # we cannot use a tree artifact here because bazel will create symlinks in
    # the sandbox
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        tar -xf "${{0}}.runfiles/{git_archive}"
        exec '{git}' --git-dir .git {arguments} "${{@}}"
    """.format(
        git = git_toolchain.git_bin,
        git_archive = "{}/{}".format(ctx.workspace_name, git_info.archive.short_path),
        arguments = " ".join([
            shell.quote(argument)
            for argument in ctx.attr.arguments
        ]),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = script,
        ),
    ]

al_git_binary = rule(
    doc = "Run git binary",
    executable = True,
    toolchains = ["//bzl/rules/git:toolchain_type"],
    implementation = _impl,
    attrs = {
        "src": attr.label(
            providers = [AlGitInfo],
            allow_single_file = True,
            mandatory = True,
            doc = "Git info",
        ),
        "arguments": attr.string_list(
            doc = "Git arguments (templated)",
            default = [],
        ),
    },
)
