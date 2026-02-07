load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_helm_chart_info.bzl", "AlHelmChartInfo")

def _impl(ctx):
    helm = ctx.toolchains["//tools/helm/main/bzl:toolchain_type"]
    runfiles_files = []
    runfiles_symlinks = {}
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))

    for data in ctx.attr.data:
        if PackageFilegroupInfo in data:
            for files in data[PackageFilegroupInfo].pkg_files:
                runfiles_symlinks.update(files.dest_src_map)
        elif PackageFilesInfo in data:
            runfiles_symlinks.update(data[PackageFilesInfo].dest_src_map)

    cd = ctx.attr.cd

    runfiles = ctx.runfiles(files = runfiles_files, symlinks = runfiles_symlinks)
    runfiles = runfiles.merge(helm.default_info.default_runfiles)

    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        helm="${{PWD}}"/'{helm}'
        cd "{cd}"
        exec "${{helm}}" \
            {arguments} \
            "${{@}}"
    """.format(
        helm = helm.helm.short_path,
        cd = cd,
        arguments = " ".join(args),
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

al_helm_binary = rule(
    implementation = _impl,
    executable = True,
    doc = "Helm binary",
    toolchains = ["//tools/helm/main/bzl:toolchain_type"],
    attrs = {
        "cd": attr.string(
            doc = "Cd to a directory before running bazel",
            default = ".",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Helm arguments",
        ),
        "data": attr.label_list(
            providers = [[PackageFilegroupInfo], [PackageFilesInfo]],
            doc = "Helm chart",
        ),
    },
)
