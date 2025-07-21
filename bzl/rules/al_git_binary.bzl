load("@bazel_skylib//lib:shell.bzl", "shell")
load("//bzl/providers:al_git_info.bzl", "AlGitInfo")

def _impl(ctx):
    git = ctx.toolchains["//bzl/toolchain-types:git"]
    script = ctx.actions.declare_file("{}-script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = [ctx.file.src])

    script_content = """\
        #!/usr/bin/env sh
        set -eu
        '{git}' -C '{git_dir}' {arguments} "${{@}}"
    """.format(
        git = git.git_bin,
        git_dir = ctx.file.src.short_path,
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
    toolchains = ["//bzl/toolchain-types:git"],
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
