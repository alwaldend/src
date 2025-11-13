_BUILD = """
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

native_binary(
    name = "bin",
    src = "bin_src",
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    download = ctx.download(
        url = ctx.attr.url,
        output = "bin_src.gz",
        integrity = ctx.attr.integrity,
    )
    ctx.execute(
        ["gzip", "--decompress", "bin_src.gz"],
    )
    ctx.execute(
        ["chmod", "u+x", "bin_src"],
    )
    ctx.file("BUILD.bazel", _BUILD)

    return ctx.repo_metadata(
        reproducible = True,
    )

al_gzip_repo = repository_rule(
    doc = "Gzip repo",
    implementation = _impl,
    attrs = {
        "url": attr.string(
            mandatory = True,
            doc = "Url",
        ),
        "integrity": attr.string(
            doc = "Integrity",
        ),
    },
)
