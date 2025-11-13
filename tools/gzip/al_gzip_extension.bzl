load("//tools/gzip:al_gzip_repo.bzl", "al_gzip_repo")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for download in mod.tags.download:
            root_module_direct_deps.append(download.name)
            al_gzip_repo(
                name = download.name,
                url = download.url,
                integrity = download.integrity,
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_gzip_extension = module_extension(
    implementation = _impl,
    doc = "Extension for gzip repos",
    tag_classes = {
        "download": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "url": attr.string(
                mandatory = True,
                doc = "Url",
            ),
            "integrity": attr.string(
                doc = "Integrity",
            ),
        }),
    },
)
