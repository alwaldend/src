load(":binary_toolchain_root_repo.bzl", "binary_toolchain_root_repo")
load(":binary_toolchain_lock.bzl", "binary_toolchain_parse_lock")
load(":binary_toolchain_repo.bzl", "binary_toolchain_repo")

def _impl(ctx):
    root_module_direct_deps = []
    root_module_direct_dev_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchain:
            lock = binary_toolchain_parse_lock(json.decode(ctx.read(tag.lock)))
            for archive_name, _ in lock.version_archives.items():
                name = "{}_{}".format(tag.name, archive_name)
                binary_toolchain_repo(
                    name = name,
                    lock = tag.lock,
                    archive_name = archive_name,
                )
            binary_toolchain_root_repo(
                name = tag.name,
                lock = tag.lock
            )
            if ctx.is_dev_dependency(tag):
                root_module_direct_dev_deps.append(tag.name)
            else:
                root_module_direct_deps.append(tag.name)
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = root_module_direct_dev_deps,
        reproducible = True,
    )

binary_toolchain_extension = module_extension(
    implementation = _impl,
    doc = "Binary toolchain extension",
    tag_classes = {
        "toolchain": tag_class(
            {
                "name": attr.string(
                    mandatory = True,
                    doc = "Repo name",
                ),
                "lock": attr.label(
                    mandatory = True,
                    doc = "Toolchain lock",
                ),
            },
        ),
    },
)
