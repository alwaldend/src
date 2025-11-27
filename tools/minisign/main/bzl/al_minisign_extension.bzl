load("//tools/minisign/main/bzl:al_minisign_repo.bzl", "al_minisign_repo")

def _impl(ctx):
    direct_deps = []
    dev_deps = []

    for mod in ctx.modules:
        for tag in mod.tags.repo:
            al_minisign_repo(
                name = tag.name,
                version = tag.version,
                url = tag.url,
                integrity = tag.integrity,
                platform_to_compat = tag.platform_to_compat,
                platform_to_extension = tag.platform_to_extension,
            )
            direct_deps.append(tag.name)

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
            "url": attr.string(
                default = "https://github.com/jedisct1/minisign/releases/download/{version}/minisign-{version}-{platform}{extension}",
                doc = "Release url template",
            ),
            "integrity": attr.string_dict(
                doc = "Integrity, keys are files, values are integrities",
            ),
            "platform_to_extension": attr.string_dict(
                default = {
                    "linux": ".tar.gz",
                    "win64": ".zip",
                    "macos": ".zip",
                },
            ),
            "platform_to_compat": attr.string_list_dict(
                default = {
                    "linux": [
                        "@platforms//os:linux",
                        "@platforms//cpu:x86_64",
                    ],
                    "win64": [
                        "@platforms//os:windows",
                        "@platforms//cpu:x86_64",
                    ],
                    "macos": [
                        "@platforms//os:macos",
                        "@platforms//cpu:x86_64",
                    ],
                },
                doc = "Platorms, keys are minisign platforms, values are bazel platforms",
            ),
        }),
    },
)
