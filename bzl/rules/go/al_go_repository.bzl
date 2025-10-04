load("@gazelle//:deps.bzl", "go_repository")

def _impl(ctx):
    for mod in ctx.modules:
        for tag in mod.tags.go_repository:
            go_repository(
                name = tag.name,
                importpath = tag.importpath,
                sum = tag.sum,
                version = tag.version,
            )

al_go_repository = module_extension(
    implementation = _impl,
    doc = "Extension wrapper around go_repository (useless because you can just call use_repo_rule)",
    tag_classes = {
        "go_repository": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "importpath": attr.string(mandatory = True, doc = "importpath"),
            "sum": attr.string(doc = "checksum"),
            "version": attr.string(mandatory = True),
        }),
    },
)
