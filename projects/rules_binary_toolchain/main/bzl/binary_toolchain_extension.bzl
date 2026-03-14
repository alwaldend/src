load(":binary_toolchain_lock.bzl", "binary_toolchain_lock_find_archives", "binary_toolchain_lock_parse")
load(":binary_toolchain_repo.bzl", "binary_toolchain_repo")
load(":binary_toolchain_root_repo.bzl", "binary_toolchain_root_repo")

def _impl(ctx):
    root_module_direct_deps = []
    root_module_direct_dev_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchains:
            lock = binary_toolchain_lock_parse(json.decode(ctx.read(tag.lock)))
            if ctx.is_dev_dependency(tag):
                dep_archive = root_module_direct_dev_deps
            else:
                dep_archive = root_module_direct_deps
            toolchains = {}
            for toolchain in tag.toolchains:
                toolchains[toolchain] = toolchain
            toolchains = toolchains | tag.toolchains_map
            for toolchain_name, bazel_name in toolchains.items():
                binary_toolchain_root_repo(
                    name = bazel_name,
                    toolchain_name = toolchain_name,
                    lock = tag.lock,
                )
                dep_archive.append(bazel_name)
                for archive in binary_toolchain_lock_find_archives(lock, toolchain_name):
                    name = "{}_{}".format(bazel_name, archive.archive_name)
                    if archive.archive_name in tag.archives.get(toolchain_name, []):
                        dep_archive.append(name)
                    binary_toolchain_repo(
                        name = name,
                        archive_lock = "@{}//:toolchain_{}_archive_lock.json".format(bazel_name, archive.archive_name),
                    )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = root_module_direct_dev_deps,
        reproducible = True,
    )

binary_toolchain_extension = module_extension(
    implementation = _impl,
    doc = "Binary toolchain extension",
    tag_classes = {
        "toolchains": tag_class(
            {
                "lock": attr.label(
                    mandatory = True,
                    doc = "Toolchain lock",
                ),
                "toolchains": attr.string_list(
                    doc = "List of toolchains names that will be exposed",
                ),
                "toolchains_map": attr.string_dict(
                    doc = "Map of toolchains that will be exposed with different names, keys are toolchain names, values are new names",
                ),
                "archives": attr.string_list_dict(
                    doc = "Map of archives that will be exposed, keys are toolchain names, values are list of archive names",
                ),
            },
        ),
    },
)
