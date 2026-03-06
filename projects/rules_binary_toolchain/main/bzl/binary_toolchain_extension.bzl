load(":binary_toolchain_root_repo.bzl", "binary_toolchain_root_repo")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchain:
            binary_toolchain_root_repo(
                name = tag.name,
                lock = tag.lock
            )
            root_module_direct_deps.append(tag.name)
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
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
