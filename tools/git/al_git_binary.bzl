load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    git_toolchain = ctx.toolchains["//tools/git:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    args = [git_toolchain.git_path, "--git-dir", git_toolchain.git_dir] + ctx.attr.arguments
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        exec {arguments} "${{@}}"
    """.format(
        args = args,
        arguments = " ".join([shell.quote(argument) for argument in args]),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    return [
        DefaultInfo(
            executable = script,
        ),
    ]

al_git_binary = rule(
    doc = "Run git binary",
    executable = True,
    toolchains = ["//tools/git:toolchain_type"],
    implementation = _impl,
    attrs = {
        "arguments": attr.string_list(
            doc = "Git arguments (templated)",
            default = [],
        ),
    },
)
