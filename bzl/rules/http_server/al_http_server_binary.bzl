load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")

def _impl(ctx):
    runfiles_symlinks = {}
    args = []
    python = ctx.toolchains["@rules_python//python:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))

    for src in ctx.attr.srcs:
        for files, _ in src[PackageFilegroupInfo].pkg_files:
            runfiles_symlinks.update(files.dest_src_map)

    runfiles = ctx.runfiles(
        transitive_files = depset(transitive = [python.py3_runtime.files]),
        symlinks = runfiles_symlinks,
    )

    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        find -maxdepth 3 -ls
        exec '{python3}' \
            -m http.server \
            {arguments} \
            "${{@}}"
    """.format(
        python3 = python.py3_runtime.interpreter.short_path,
        arguments = " ".join([
            shell.quote(ctx.expand_location(ctx.expand_make_variables("al_http_server_binary", arg, {})))
            for arg in args
        ]),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )

    default_info = DefaultInfo(
        runfiles = runfiles,
        executable = script,
    )
    return [default_info]

al_http_server_binary = rule(
    implementation = _impl,
    doc = "Run a http server",
    toolchains = ["@rules_python//python:toolchain_type"],
    executable = True,
    attrs = {
        "srcs": attr.label_list(
            providers = [PackageFilegroupInfo],
            doc = "Files to symlink",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Arguments",
        ),
    },
)
