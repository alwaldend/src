load("//tools/git:al_git_repo.bzl", "al_git_repo")

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.local_git:
            root_module_direct_deps.append(tag.name)
            al_git_repo(
                name = tag.name,
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_git_extension = module_extension(
    implementation = _impl,
    doc = "Create git repos",
    tag_classes = {
        "local_git": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
        }),
    },
)
