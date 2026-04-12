load(":binary_toolchain_lock.bzl", "binary_toolchain_lock_parse_archive")

_BUILD = """\
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

ARCHIVE = json.decode({archive})

[
    [
        native_binary(
            name = "{{}}_native_binary".format(binary["name"]),
            src = binary["path"],
            exec_compatible_with = ARCHIVE["toolchain"].get("exec_compatible_with"),
            visibility = ARCHIVE["toolchain"].get("visibility", ["//visibility:public"]),
        ),
        filegroup(
            name = "{{}}_filegroup".format(binary["name"]),
            srcs = [binary["path"]],
            visibility = ARCHIVE["toolchain"].get("visibility", ["//visibility:public"]),
        )
    ]
    for binary in ARCHIVE["binaries"]
]
"""

def _impl(ctx):
    archive = binary_toolchain_lock_parse_archive(json.decode(ctx.read(ctx.attr.archive_lock)))
    for action in archive.actions:
        if action.download_and_extract:
            download = ctx.download_and_extract(**action.download_and_extract)
            integrity = action.download_and_extract.get("integrity", "")
            if download.integrity != integrity:
                fail("invalid integrity: wanted '{}', got '{}'".format(integrity, download.integrity))
        elif action.download:
            download = ctx.download(**action.download)
            integrity = action.download.get("integrity", "")
            if download.integrity != integrity:
                fail("invalid integrity: wanted '{}', got '{}'".format(integrity, download.integrity))
        elif action.extract:
            ctx.extract(**action.extract)
        elif action.execute:
            execute = {} | action.execute
            arguments = execute.pop("arguments", [])
            ctx.execute(arguments, **execute)
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            archive = json.encode(json.encode_indent(archive, indent = "    ")),
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
