load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _shellcheck_impl(ctx):
    key = "{}:{}".format(ctx.os.name, ctx.os.arch)
    for mod in ctx.modules:
        for install in mod.tags.install:
            url = install.urls[key]
            integrity = install.integrity[key]
            http_archive(
                name = install.name,
                url = url,
                integrity = integrity,
                build_file_content = """
filegroup(
    name = "shellcheck",
    srcs = glob(["*/shellcheck"]),
    visibility = ["//visibility:public"],
)
                """,
            )

shellcheck = module_extension(
    implementation = _shellcheck_impl,
    doc = "Extension to download shellcheck",
    tag_classes = {
        "install": tag_class({
            "name": attr.string(mandatory = True, doc = "name"),
            "urls": attr.string_dict(
                mandatory = True,
                doc = "urls",
            ),
            "integrity": attr.string_dict(
                mandatory = True,
                doc = "integrity for the urls",
            ),
        }),
    },
)
