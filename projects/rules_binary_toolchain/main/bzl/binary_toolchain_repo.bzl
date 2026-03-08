load(":binary_toolchain_lock.bzl", "binary_toolchain_lock_parse_archive")

_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

ARCHIVE = {archive}

[
    native_binary(
        name = "{{}}_binary".format(binary["name"]),
        src = binary["path"],
        exec_compatible_with = ARCHIVE["toolchain"].get("exec_compatible_with"),
        visibility = ARCHIVE["toolchain"].get("visibility", ["//visibility:public"]),
    )
    for binary in ARCHIVE["binaries"]
]
"""

def _impl(ctx):
    archive = binary_toolchain_lock_parse_archive(json.decode(ctx.read(ctx.attr.archive_lock)))
    if archive.download_and_extract:
        download = ctx.download_and_extract(**archive.download_and_extract)
        integrity = archive.download_and_extract.get("integrity", "")
        if download.integrity != integrity:
            fail("invalid integrity: wanted '{}', got '{}'".format(integrity, download.integrity))
    elif archive.download:
        download = ctx.download(**archive.download)
        integrity = archive.download.get("integrity", "")
        if download.integrity != integrity:
            fail("invalid integrity: wanted '{}', got '{}'".format(integrity, download.integrity))
    else:
        fail("missing both download and download_and_extract")
    if archive.extract:
        ctx.extract(**archive.extract)
    if archive.execute:
        for execute in archive.execute:
            execute = {} | execute
            arguments = execute.pop("arguments", [])
            ctx.execute(arguments, **execute)
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            archive = json.encode_indent(archive, indent = "    "),
        ),
    )
    return ctx.repo_metadata(reproducible = True)

binary_toolchain_repo = repository_rule(
    implementation = _impl,
    doc = "Binary toolchain repo",
    attrs = {
        "archive_lock": attr.label(
            mandatory = True,
            doc = "Archive lock",
        ),
    },
)
