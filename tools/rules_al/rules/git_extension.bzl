load("//rules:git_repo.bzl", "git_repo")

def _impl(ctx):
    for mod in ctx.modules:
        for tag in mod.tags.local_git:
            git_repo(
                name = tag.name,
            )
    return ctx.extension_metadata(reproducible = True)

git_extension = module_extension(
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
