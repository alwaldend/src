load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("//tools/bazel_contracts:al_bazel_contracts_repo.bzl", "al_bazel_contracts_repo")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.download:
            root_module_direct_deps.append(tag.name)
            al_bazel_contracts_repo(
                name = tag.name,
                prefix = tag.prefix,
                contracts = tag.contracts,
                base_url = tag.base_url,
                integrity = tag.integrity,
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_bazel_contracts_extension = module_extension(
    implementation = _impl,
    doc = "Create a repo of bazel contracts",
    tag_classes = {
        "download": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "integrity": attr.string_dict(
                mandatory = True,
                doc = "Integrity",
            ),
            "base_url": attr.string(
                mandatory = True,
                doc = "Base url",
            ),
            "prefix": attr.string(
                mandatory = True,
                doc = "Prefix for the contracts inside the repo",
            ),
            "contracts": attr.string_list(
                mandatory = True,
                doc = "List of contract paths",
            ),
        }),
    },
)
