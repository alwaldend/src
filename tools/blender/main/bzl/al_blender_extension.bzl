load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//tools/bzl/main/bzl:al_bzl_generate_repository.bzl", "al_bzl_generate_repository")
load(":al_blender_archives.bzl", "AL_BLENDER_ARCHIVES")

_BUILD_TOOLCHAIN = """
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

native_binary(
    name = "blender",
    src = glob(["*/blender"])[0],
    exec_compatible_with = {platforms},
    visibility = ["//visibility:public"],
)
"""

_BUILD_ROOT = """
load("@com_alwaldend_src//tools/blender/main/bzl:al_blender_toolchain.bzl", "al_blender_toolchain")
CHILDREN = {children}

[
    [
        toolchain(
            name = "toolchain_{{}}".format(child),
            toolchain = ":toolchain_{{}}_impl".format(child),
            toolchain_type = "@com_alwaldend_src//tools/blender/main/bzl:toolchain_type",
            exec_compatible_with = platforms,
            visibility = ["//visibility:public"],
        ),
        al_blender_toolchain(
            name = "toolchain_{{}}_impl".format(child),
            exec_compatible_with = platforms,
            blender = "@{{}}//:blender".format(child),
        )
    ]
    for child, platforms in CHILDREN
]
"""

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchains:
            children = []
            for archive in AL_BLENDER_ARCHIVES.get(tag.version, []):
                name = "_".join([tag.name] + archive["platforms"]).replace(":", "_")
                platforms = ["@platforms//{}".format(platform) for platform in archive["platforms"]]
                children.append((name, platforms))
                http_archive(
                    name = name,
                    urls = archive.get("urls"),
                    url = archive.get("url"),
                    strip_prefix = archive.get("strip_prefix"),
                    build_file_content = _BUILD_TOOLCHAIN.format(
                        name = name,
                        platforms = platforms,
                    ),
                    integrity = archive.get("integrity", ""),
                )
            root_module_direct_deps.append(tag.name)
            al_bzl_generate_repository(
                name = tag.name,
                files = {
                    "BUILD.bazel": _BUILD_ROOT.format(children = children),
                },
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_blender_extension = module_extension(
    implementation = _impl,
    doc = "blender extension",
    tag_classes = {
        "toolchains": tag_class(
            {
                "name": attr.string(
                    mandatory = True,
                    doc = "Name",
                ),
                "version": attr.string(
                    mandatory = True,
                    doc = "Release version",
                ),
            },
        ),
    },
)
