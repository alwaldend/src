load("@bazel_skylib//lib:shell.bzl", "shell")
load(":al_helm_chart_info.bzl", "AlHelmChartInfo")

def _impl(ctx):
    helm = ctx.toolchains["//bzl/rules/helm:toolchain_type"]
    runfiles_files = [] + ctx.files.chart
    runfiles_symlinks = {}
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))

    cd = ctx.attr.cd
    if ctx.attr.chart:
        chart = ctx.attr.chart[AlHelmChartInfo]
        runfiles_symlinks.update(chart.files_info.dest_src_map)

    runfiles = ctx.runfiles(files = runfiles_files, symlinks = runfiles_symlinks)

    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        helm="${{PWD}}"/'{helm}'
        cd "{cd}"
        exec "${{helm}}" \
            {arguments} \
            "${{@}}"
    """.format(
        helm = helm.helm.short_path,
        cd = cd,
        arguments = " ".join([shell.quote(arg) for arg in args]),
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
    toolchains = ["//bzl/rules/helm:toolchain_type"],
    attrs = {
        "cd": attr.string(
            doc = "Cd to a directory before running bazel",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Helm arguments",
        ),
        "chart": attr.label(
            providers = [AlHelmChartInfo],
            doc = "Helm chart",
        ),
    },
)
