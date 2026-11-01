load(":al_rfc_repo.bzl", "al_rfc_repo")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.download:
            root_module_direct_deps.append(tag.name)
            al_rfc_repo(
                name = tag.name,
                rfcs = tag.rfcs,
                integrity = tag.integrity,
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_rfc_extension = module_extension(
    implementation = _impl,
    doc = "Rfc extension",
    tag_classes = {
        "download": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "rfcs": attr.string_list(
                doc = "Rfcs",
            ),
            "integrity": attr.string_dict(
                doc = "Rfc integrity",
            ),
        }),
    },
)
