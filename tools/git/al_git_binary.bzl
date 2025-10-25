load("@bazel_skylib//lib:shell.bzl", "shell")
load(":al_git_library_info.bzl", "AlGitLibraryInfo")

def _impl(ctx):
    git_toolchain = ctx.toolchains["//tools/git:toolchain_type"]
    git_info = ctx.attr.src[AlGitLibraryInfo]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(
        transitive_files = ctx.attr.src[DefaultInfo].default_runfiles.files,
    )

    script_content = """\
        #!/usr/bin/env sh
        set -eu
        if [ -f '{git_archive}' ]; then
            git_archive="${{PWD}}"/'{git_archive}'
        elif [ -f "${{0}}.runfiles/{workspace_name}/{git_archive}" ]; then
            git_archive="${{PWD}}/${{0}}.runfiles/{workspace_name}/{git_archive}"
        else
            echo "Missing runfiles"
            exit 1
        fi
        tar -xf "${{git_archive}}"
        exec '{git}' {arguments} "${{@}}"
    """.format(
        git = git_toolchain.git_bin,
        workspace_name = ctx.workspace_name,
        git_archive = git_info.git_archive.short_path,
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
    toolchains = ["//tools/git:toolchain_type"],
    implementation = _impl,
    attrs = {
        "src": attr.label(
            providers = [AlGitLibraryInfo],
            mandatory = True,
            doc = "Git info",
        ),
        "arguments": attr.string_list(
            doc = "Git arguments (templated)",
            default = [],
        ),
    },
)
