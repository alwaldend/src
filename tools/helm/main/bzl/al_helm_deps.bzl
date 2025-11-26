load("//tools/helm/main/bzl:al_helm_deps_repo.bzl", "al_helm_deps_repo")

def _impl(ctx):
    direct_deps = []
    dev_deps = []

    for mod in ctx.modules:
        for tag in mod.tags.from_locks:
            al_helm_deps_repo(
                name = tag.name,
                locks = tag.locks,
                integrity = tag.integrity,
            )
            direct_deps.append(tag.name)

    return ctx.extension_metadata(
        root_module_direct_deps = direct_deps,
        root_module_direct_dev_deps = dev_deps,
        reproducible = True,
    )

al_helm_deps = module_extension(
    implementation = _impl,
    doc = "Extension to download helm dependencies",
    tag_classes = {
        "from_locks": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Repo name",
            ),
            "locks": attr.label_list(
                doc = "Helm lock labels",
            ),
            "integrity": attr.string_dict(
                doc = "Intergrity for locks, keys are packages, values are integrity",
            ),
        }),
    },
)
