load(":binary_toolchain_lock.bzl", "binary_toolchain_lock_find_archives", "binary_toolchain_lock_parse")

_DEFS_HEADER = """\
load("@rules_binary_toolchain//main/bzl:binary_toolchain_toolchain.bzl", "binary_toolchain_toolchain")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_binary.bzl", "binary_toolchain_binary")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_run_binary.bzl", "binary_toolchain_run_binary")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_test.bzl", "binary_toolchain_test")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_resolved_toolchain.bzl", "binary_toolchain_resolved_toolchain")

"""

_DEFS_BODY = """\
{name}_toolchain_type = str(Label("//:{name}_toolchain_type"))
{name}_toolchain = binary_toolchain_toolchain
{name}_binary = binary_toolchain_binary({name}_toolchain_type)
{name}_run_binary = binary_toolchain_run_binary({name}_toolchain_type)
{name}_test = binary_toolchain_test({name}_toolchain_type)
{name}_resolved_toolchain = binary_toolchain_resolved_toolchain({name}_toolchain_type)
"""

_BUILD_HEADER = """\
ARCHIVES = {archives}
"""

_BUILD_BIN = """\
load(":defs.bzl", "{name}_binary", "{name}_resolved_toolchain", "{name}_toolchain")

toolchain_type(
    name = "{name}_toolchain_type",
    visibility = ["//visibility:public"],
)

{name}_binary(
    name = "{name}_binary",
    visibility = ["//visibility:public"],
)

{name}_resolved_toolchain(
    name = "{name}_toolchain_resolved",
    visibility = ["//visibility:public"],
)

[
    [
        [
            toolchain(
                name = "{name}_toolchain_{{}}".format(archive["archive_name"]),
                toolchain = ":{name}_toolchain_{{}}_impl".format(archive["archive_name"]),
                toolchain_type = ":{name}_toolchain_type",
                exec_compatible_with = archive["toolchain"].get("exec_compatible_with"),
                visibility = archive["toolchain"].get("visibility", ["//visibility:public"]),
            ),
            {name}_toolchain(
                name = "{name}_toolchain_{{}}_impl".format(archive["archive_name"]),
                binary = "@{{}}_{{}}//:{name}_native_binary".format("{original_name}", archive["archive_name"]),
                exec_compatible_with = archive["toolchain"].get("exec_compatible_with"),
                visibility = archive["toolchain"].get("visibility", ["//visibility:public"]),
            )
        ]
        for binary in archive["binaries"]
        if binary["name"] == "{name}"
    ]
    for archive in ARCHIVES
]
"""

def _impl(ctx):
    lock = binary_toolchain_lock_parse(json.decode(ctx.read(ctx.attr.lock)))
    archives = binary_toolchain_lock_find_archives(lock, ctx.attr.toolchain_name)
    all_bins = {}
    for archive in archives:
        ctx.file(
            "toolchain_{}_archive_lock.json".format(archive.archive_name),
            json.encode_indent(archive, indent = "    "),
        )
        cur_bins = {}
        for bin in archive.binaries:
            if bin.name in cur_bins:
                fail("duplicate bin {}".format(bin.name))
            cur_bins[bin.name] = True
        all_bins.update(cur_bins)
    defs = [_DEFS_HEADER]
    for bin in all_bins.keys():
        defs.append(_DEFS_BODY.format(name = bin))
    ctx.file("defs.bzl", "\n".join(defs))
    build = [_BUILD_HEADER.format(archives = json.encode_indent(archives, indent = "    "))]
    for bin in all_bins.keys():
        build.append(
            _BUILD_BIN.format(
                name = bin,
                original_name = ctx.original_name,
            ),
        )
    ctx.file("BUILD.bazel", "\n".join(build))
    return ctx.repo_metadata(
        reproducible = True,
    )

binary_toolchain_root_repo = repository_rule(
    implementation = _impl,
    doc = "Root binary toolchain repository",
    attrs = {
        "toolchain_name": attr.string(
            doc = "Toolchain name this repo is created for",
        ),
        "lock": attr.label(
            mandatory = True,
            doc = "Lock file label",
        ),
    },
)
