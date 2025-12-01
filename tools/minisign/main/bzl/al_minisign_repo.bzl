load("//tools/minisign/main/bzl:al_minisign_archives.bzl", "AL_MINISIGN_ARCHIVES")

_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")
load("@com_alwaldend_src//tools/minisign/main/bzl:al_minisign_toolchain.bzl", "al_minisign_toolchain")

_TOOLCHAINS = {toolchains}

[
    [
        toolchain(
            name = toolchain_config["name"],
            toolchain = "{{}}_impl".format(toolchain_config["name"]),
            toolchain_type = "@com_alwaldend_src//tools/minisign/main/bzl:toolchain_type",
            exec_compatible_with = toolchain_config["exec_compatible_with"],
            visibility = ["//visibility:public"],
        ),
        al_minisign_toolchain(
            name = "{{}}_impl".format(toolchain_config["name"]),
            minisign = ":{{}}_binary".format(toolchain_config["name"]),
            exec_compatible_with = toolchain_config["exec_compatible_with"],
            visibility = ["//visibility:public"],
        ),
        native_binary(
            name = "{{}}_binary".format(toolchain_config["name"]),
            src = toolchain_config["binary"],
            exec_compatible_with = toolchain_config["exec_compatible_with"],
            visibility = ["//visibility:public"],
        ),
    ]
    for toolchain_config in _TOOLCHAINS
]
"""

def _impl(ctx):
    downloads = []
    integrity = {}
    archive = AL_MINISIGN_ARCHIVES.get(ctx.attr.version, {}).get(ctx.attr.platform, {})
    if not archive:
        fail("missing platform '{}' for minisign version '{}'".format(ctx.attr.platform, ctx.attr.version))
    download = ctx.download_and_extract(
        url = archive["url"],
        integrity = archive["integrity"],
    )
    ctx.file("BUILD.bazel", _BUILD.format(toolchains = json.encode(archive["toolchains"])))
    return ctx.repo_metadata(reproducible = True)

al_minisign_repo = repository_rule(
    implementation = _impl,
    doc = "Minisign repo",
    attrs = {
        "version": attr.string(
            mandatory = True,
            doc = "Minisign version",
        ),
        "platform": attr.string(
            mandatory = True,
            doc = "Minisign platform",
        ),
    },
)
