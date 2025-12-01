load("//tools/minisign/main/bzl:al_minisign_archives.bzl", "AL_MINISIGN_ARCHIVES")
load("//tools/minisign/main/bzl:al_minisign_repo.bzl", "al_minisign_repo")

def _impl(ctx):
    direct_deps = []
    dev_deps = []

    for mod in ctx.modules:
        for tag in mod.tags.repo:
            archive = AL_MINISIGN_ARCHIVES.get(tag.version)
            if not archive:
                fail(
                    "missing minisign version '{}', available: {}".format(
                        tag.version,
                        AL_MINISIGN_ARCHIVES.keys(),
                    ),
                )
            for platform, _ in archive.items():
                name = "{}_{}".format(tag.name, platform)
                al_minisign_repo(
                    name = name,
                    version = tag.version,
                    platform = platform,
                )
                direct_deps.append(name)

    return ctx.extension_metadata(
        root_module_direct_deps = direct_deps,
        root_module_direct_dev_deps = dev_deps,
        reproducible = True,
    )

al_minisign_extension = module_extension(
    implementation = _impl,
    doc = "Minisign extension",
    tag_classes = {
        "repo": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Repo name",
            ),
            "version": attr.string(
                mandatory = True,
                doc = "Minisign version",
            ),
        }),
    },
)
