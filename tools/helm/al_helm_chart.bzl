load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_helm_chart_info.bzl", "AlHelmChartInfo")

def _impl(ctx):
    if not ctx.attr.source and not ctx.attr.package:
        fail("missing both `source` and `package`")

    files = []
    files_transitive = []
    files_info = PackageFilesInfo(
        attributes = {},
        dest_src_map = {},
    )

    if ctx.attr.source:
        source = ctx.attr.source[PackageFilegroupInfo]
        package = None
        for info, _ in ctx.attr.source[PackageFilegroupInfo].pkg_files:
            files.extend(info.dest_src_map.values())
            files_info.dest_src_map.update(info.dest_src_map)
            files_info.attributes.update(info.attributes)
        for dep in ctx.attr.deps:
            package = dep[AlHelmChartInfo].package
            files_transitive.append(package)
            files_info.dest_src_map["{}/{}".format("charts", package.basename)] = package
    else:
        source = None
        package = ctx.file.package
        files.append(package)

    default_info = DefaultInfo(
        files = depset(files, transitive = [depset(files_transitive)]),
    )
    al_helm_chart_info = AlHelmChartInfo(
        source = source,
        deps = depset(ctx.attr.deps),
        package = package,
        files_info = files_info,
    )
    return [default_info, al_helm_chart_info, files_info]

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
            providers = [PackageFilegroupInfo],
            doc = "Helm chart source",
        ),
        "package": attr.label(
            allow_single_file = [".tgz"],
            doc = "Helm chart package",
        ),
    },
)
