load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(":al_trufflehog_archives.bzl", "AL_TRUFFLEHOG_ARCHIVES")

_BUILD_TOOLCHAIN = """
load("@com_alwaldend_src//tools/trufflehog/main/bzl:al_trufflehog_toolchain.bzl", "al_trufflehog_toolchain")
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

native_binary(
    name = "bin",
    src = "trufflehog",
    exec_compatible_with = {platforms},
    visibility = ["//visibility:public"],
)

alias(
    name = "{name}",
    actual = ":toolchain",
    visibility = ["//visibility:public"],
)

toolchain(
    name = "toolchain",
    toolchain = "toolchain_impl",
    toolchain_type = "@com_alwaldend_src//tools/trufflehog/main/bzl:toolchain_type",
    exec_compatible_with = {platforms},
    visibility = ["//visibility:public"],
)

al_hugo_toolchain(
    name = "toolchain_impl",
    exec_compatible_with = {platforms},
    hugo = ":bin",
)
"""

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchains:
            for archive in AL_TRUFFLEHOG_ARCHIVES.get(tag.version, []):
                name = "_".join([tag.name] + archive["platforms"]).replace(":", "_")
                root_module_direct_deps.append(name)
                http_archive(
                    name = name,
                    url = archive["url"],
                    strip_prefix = archive.get("strip_prefix"),
                    build_file_content = _BUILD_TOOLCHAIN.format(
                        name = name,
                        platforms = ["@platforms//{}".format(platform) for platform in archive["platforms"]],
                    ),
                    integrity = archive.get("integrity", ""),
                )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_trufflehog_extension = module_extension(
    implementation = _impl,
    doc = "Trufflehog extension",
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
