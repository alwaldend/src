load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")

_EXECUTABLE = """\
#!/usr/bin/env sh
set -eu
cd "${{0}}.runfiles/{workspace}"
exec make {args} "${{@}}"
"""

def _impl(ctx):
    executable = ctx.actions.declare_file("{}.sh".format(ctx.label.name))
    runfiles_symlinks = {}
    for files, _ in ctx.attr.src[PackageFilegroupInfo].pkg_files:
        runfiles_symlinks.update(files.dest_src_map)
    ctx.actions.write(
        output = executable,
        is_executable = True,
        content = _EXECUTABLE.format(
            workspace = ctx.workspace_name,
            args = " ".join([shell.quote(arg) for arg in ctx.attr.args]),
        ),
    )
    return [
        DefaultInfo(
            executable = executable,
            runfiles = ctx.runfiles(symlinks = runfiles_symlinks),
        ),
    ]

al_make_install_binary = rule(
    doc = "Create a binary target for a make install",
    implementation = _impl,
    executable = True,
    attrs = {
        "arguments": attr.string_list(
            doc = "Arugments for the executable",
        ),
        "src": attr.label(
            mandatory = True,
            providers = [PackageFilegroupInfo],
            doc = "Make install target",
        ),
    },
)
