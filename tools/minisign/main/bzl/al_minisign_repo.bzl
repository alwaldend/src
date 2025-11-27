_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")
load("@com_alwaldend_src//tools/minisign/main/bzl:al_minisign_toolchain.bzl", "al_minisign_toolchain")

toolchain(
    name = "toolchain",
    toolchain = "toolchain_impl",
    toolchain_type = "@com_alwaldend_src//tools/minisign/main/bzl:toolchain_type",
    visibility = ["//visibility:public"],
)

al_minisign_toolchain(
    name = "toolchain_impl",
    minisign = ":minisign",
    exec_compatible_with = {exec_compatible_with},
    visibility = ["//visibility:public"],
)

native_binary(
    name = "minisign",
    src = glob(["**/x86_64/**/minisign", "**/x86_64/**/minisign.exe", "minisign"], allow_empty = True)[0],
    exec_compatible_with = {exec_compatible_with},
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    downloads = []
    integrity = {}
    for platform, platform_compat in ctx.attr.platform_to_compat.items():
        download = ctx.download_and_extract(
            url = ctx.attr.url.format(
                version = ctx.attr.version,
                platform = platform,
                extension = ctx.attr.platform_to_extension[platform],
            ),
            output = platform,
            integrity = ctx.attr.integrity.get(platform, ""),
        )
        integrity[platform] = download.integrity
        ctx.file(
            "{}/BUILD.bazel".format(platform),
            _BUILD.format(exec_compatible_with = json.encode(platform_compat)),
        )
    ctx.file("BUILD.bazel", "")

    if integrity == ctx.attr.integrity:
        return ctx.repo_metadata(reproducible = True)
    return ctx.repo_metadata(
        attrs_for_reproducibility = {
            "name": ctx.name,
            "version": ctx.attr.version,
            "url": ctx.attr.url,
            "integrity": integrity,
            "platform_to_extension": ctx.attr.platform_to_extension,
            "platform_to_compat": ctx.attr.platform_to_compat,
        },
    )

al_minisign_repo = repository_rule(
    implementation = _impl,
    doc = "Minisign repo",
    attrs = {
        "version": attr.string(
            mandatory = True,
            doc = "Minisign version",
        ),
        "url": attr.string(
            mandatory = True,
            doc = "Release url template",
        ),
        "integrity": attr.string_dict(
            doc = "Integrity, keys are files, values are integrities",
        ),
        "platform_to_extension": attr.string_dict(
            mandatory = True,
            doc = "Platform extensions, keys are platforms, values are extensions",
        ),
        "platform_to_compat": attr.string_list_dict(
            mandatory = True,
            doc = "Platorms, keys are minisign platforms, values are bazel platforms",
        ),
    },
)
