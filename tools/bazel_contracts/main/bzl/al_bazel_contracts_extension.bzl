load(":al_bazel_contracts_list.bzl", "AL_BAZEL_CONTRACTS_LIST")
load(":al_bazel_contracts_repo.bzl", "al_bazel_contracts_repo")
load(":al_bazel_contracts_versions.bzl", "AL_BAZEL_CONTRACTS_VERSIONS")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.download:
            root_module_direct_deps.append(tag.name)
            version = AL_BAZEL_CONTRACTS_VERSIONS.get(tag.version, {})
            al_bazel_contracts_repo(
                name = tag.name,
                contracts = tag.contracts or version["contracts"],
                base_url = tag.base_url or version["base_url"],
                integrity = tag.integrity or version.get("integrity", {}),
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
            "version": attr.string(
                doc = "Version",
            ),
            "integrity": attr.string_dict(
                doc = "Integrity",
            ),
            "base_url": attr.string(
                doc = "Base url",
            ),
            "prefix": attr.string(
                doc = "Prefix for the contracts inside the repo",
            ),
            "contracts": attr.string_list(
                default = AL_BAZEL_CONTRACTS_LIST,
                doc = "List of contract paths",
            ),
        }),
    },
)
