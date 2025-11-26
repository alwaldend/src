_BUILD = """\
load("@com_alwaldend_src//tools/helm/main/bzl:al_helm_chart.bzl", "al_helm_chart")
PACKAGES = {packages}

al_helm_chart(
    name = "charts",
    deps = PACKAGES,
    visibility = ["//visibility:public"],
)

[
    al_helm_chart(
        name = package,
        package = "{{}}.tgz".format(package),
        visibility = ["//visibility:public"],
    )
    for package in PACKAGES
]
"""

def _impl(ctx):
    downloads = []
    packages = {}
    for lock_file in ctx.attr.locks:
        lock = json.decode(ctx.read(lock_file))
        for dependency in lock.get("dependencies", []):
            name = "{}-{}".format(dependency["name"], dependency["version"])
            integrity = ctx.attr.integrity.get(name, "")
            download = ctx.download(
                url = "{}/{}.tgz".format(dependency["repository"], name),
                output = "{}.tgz".format(name),
                integrity = integrity,
                block = False,
            )
            packages[name] = None
            downloads.append((name, download))

    ctx.file(
        "BUILD.bazel",
        _BUILD.format(packages = packages.keys()),
    )

    integrity = {}
    for name, download in downloads:
        result = download.wait()
        integrity[name] = result.integrity

    if integrity == ctx.attr.integrity:
        return ctx.repo_metadata(reproducible = True)
    return ctx.repo_metadata(
        attrs_for_reproducibility = {
            "name": ctx.name,
            "locks": ctx.attr.locks,
            "integrity": integrity,
        },
    )

al_helm_deps_repo = repository_rule(
    implementation = _impl,
    doc = "Helm deps repo",
    attrs = {
        "locks": attr.label_list(
            doc = "Lock labels to parse",
        ),
        "integrity": attr.string_dict(
            doc = "Intergrity for locks, keys are packages, values are integrity",
        ),
    },
)
