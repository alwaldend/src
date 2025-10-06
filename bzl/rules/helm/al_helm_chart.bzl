load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load(":al_helm_chart_info.bzl", "AlHelmChartInfo")

def _impl(ctx):
    if not ctx.attr.source and not ctx.attr.package:
        fail("missing both `source` and `package`")

    files = []

    if ctx.attr.source:
        source = ctx.actions.declare_directory("{}.source".format(ctx.label.name))
        package = None
        files.append(source)
        deps = [dep[AlHelmChartInfo].package for dep in ctx.attr.deps]
        cmd = [
            "tar -xf '{}' -C '{}'".format(ctx.file.source.path, source.path),
            "mkdir '{}/charts'".format(source.path),
            "mv {} '{}/charts'".format(" ".join([shell.quote(dep.path) for dep in deps]), source.path),
        ]
        ctx.actions.run_shell(
            outputs = [source],
            inputs = [ctx.file.source] + deps,
            command = " && ".join(cmd),
        )
    else:
        source = None
        package = ctx.file.package
        files.append(package)

    default_info = DefaultInfo(
        files = depset(files, transitive = [dep.files for dep in ctx.attr.deps]),
    )
    al_helm_chart_info = AlHelmChartInfo(
        source = source,
        deps = depset(ctx.attr.deps),
        package = package,
    )
    return [default_info, al_helm_chart_info]

al_helm_chart = rule(
    implementation = _impl,
    doc = "Helm chart",
    provides = [AlHelmChartInfo],
    attrs = {
        "deps": attr.label_list(
            default = [],
            providers = [AlHelmChartInfo],
            doc = "Helm chart deps",
        ),
        "source": attr.label(
            allow_single_file = [".tar"],
            doc = "Helm chart source",
        ),
        "package": attr.label(
            allow_single_file = [".tgz"],
            doc = "Helm chart package",
        ),
    },
)
