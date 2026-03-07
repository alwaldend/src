load(":binary_toolchain_lock.bzl", "binary_toolchain_parse_lock_attr")

_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

TOOLCHAIN = {toolchain}

native_binary(
    name = "binary",
    src = "{bin_path}",
    exec_compatible_with = TOOLCHAIN.get("exec_compatible_with"),
    visibility = TOOLCHAIN.get("visibility", ["//visibility:public"]),
)
"""

def _impl(ctx):
    lock = binary_toolchain_parse_lock_attr(ctx)
    archive = lock.version_archives.get(ctx.attr.archive_name)
    if not archive:
        fail("missing archive {} for version {}".format(ctx.attr.archive_name, lock.version))
    if archive.download_and_extract:
        download = ctx.download_and_extract(**archive.download_and_extract)
        integrity = archive.download_and_extract.get("integrity", "")
        if download.integrity != integrity:
            fail("invalid integrity: wanted '{}', got '{}'".format(download.integrity, integrity))
    elif archive.download:
        download = ctx.download(**archive.download)
        integrity = archive.download.get("integrity", "")
        if download.integrity != integrity:
            fail("invalid integrity: wanted '{}', got '{}'".format(download.integrity, integrity))
    else:
        fail("missing both download and download_and_extract")
    if archive.extract:
        ctx.extract(**archive.extract)
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            bin_path = archive.bin_path,
            toolchain = json.encode_indent(archive.toolchain),
        ),
    )
    return ctx.repo_metadata(reproducible = True)

binary_toolchain_repo = repository_rule(
    implementation = _impl,
    doc = "Binary toolchain repo",
    attrs = {
        "lock": attr.label(
            mandatory = True,
            doc = "Lock file label",
        ),
        "archive_name": attr.string(
            mandatory = True,
            doc = "Archive name",
        ),
    },
)
