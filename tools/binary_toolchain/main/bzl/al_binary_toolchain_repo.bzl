_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

filegroup(
    name = "filegroup",
    srcs = ["{bin}"],
    visibility = ["//visibility:public"],
)

native_binary(
    name = "native_binary",
    src = ":filegroup",
    exec_compatible_with = {exec_compatible_with},
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    lock = json.decode(ctx.file(ctx.attr.lock))
    archive = lock.get("archives", {}).get(ctx.attr.version, {}).get(ctx.attr.archive_name, {})
    if not archive:
        fail("missing archive '{}' for version '{}'".format(ctx.attr.version, ctx.attr.archive_name))
    download = ctx.download_and_extract(
        url = archive["url"],
        output = "bin",
        integrity = archive["integrity"],
    )
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            bin = "bin/{}".format(archive["bin_path"]),
            exec_compatible_with = archive["exec_compatible_with"],
        ),
    )
    if download.integrity != archive["integrity"]:
        fail("invalid integrity: expected '{}', got '{}'".format(archive["integrity"], result.integrity))
    return ctx.repo_metadata(reproducible = True)

al_binary_toolchain_repo = repository_rule(
    implementation = _impl,
    doc = "Binary toolchain repo",
    attrs = {
        "lock": attr.label(
            mandatory = True,
            doc = "Lock file label",
        ),
        "version": attr.string(
            mandatory = True,
            doc = "Archive version",
        ),
        "archive_name": attr.string(
            mandatory = True,
            doc = "Archive name",
        ),
    },
)
