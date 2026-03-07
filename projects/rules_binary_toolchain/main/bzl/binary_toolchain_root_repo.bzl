load(":binary_toolchain_lock.bzl", "binary_toolchain_parse_lock_attr")

_DEFS = """\
load("@rules_binary_toolchain//main/bzl:binary_toolchain_toolchain.bzl", "binary_toolchain_toolchain")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_binary.bzl", "binary_toolchain_binary")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_run_binary.bzl", "binary_toolchain_run_binary")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_test.bzl", "binary_toolchain_test")
load("@rules_binary_toolchain//main/bzl:binary_toolchain_resolved_toolchain.bzl", "binary_toolchain_resolved_toolchain")

{name}_toolchain_type = str(Label("//:toolchain_type"))
{name}_toolchain = binary_toolchain_toolchain
{name}_binary = binary_toolchain_binary({name}_toolchain_type)
{name}_run_binary = binary_toolchain_run_binary({name}_toolchain_type)
{name}_test = binary_toolchain_test({name}_toolchain_type)
{name}_resolved_toolchain = binary_toolchain_resolved_toolchain({name}_toolchain_type)
"""

_BUILD = """\
load(":defs.bzl", "{name}_binary", "{name}_resolved_toolchain", "{name}_toolchain")

ARCHIVES = {archives}

toolchain_type(
    name = "toolchain_type",
    visibility = ["//visibility:public"],
)

{name}_binary(
    name = "binary",
    visibility = ["//visibility:public"],
)

{name}_resolved_toolchain(
    name = "toolchain_resolved",
    visibility = ["//visibility:public"],
)

[
    [
        toolchain(
            name = "toolchain_{{}}".format(name),
            toolchain = ":toolchain_{{}}_impl".format(name),
            toolchain_type = ":toolchain_type",
            exec_compatible_with = archive["toolchain"].get("exec_compatible_with"),
            visibility = archive["toolchain"].get("visibility", ["//visibility:public"]),
        ),
        {name}_toolchain(
            name = "toolchain_{{}}_impl".format(name),
            binary = "@{{}}_{{}}//:binary".format("{original_name}", name),
            exec_compatible_with = archive["toolchain"].get("exec_compatible_with"),
            visibility = archive["toolchain"].get("visibility", ["//visibility:public"]),
        )
    ]
    for name, archive in ARCHIVES.items()
]
"""

def _impl(ctx):
    lock = binary_toolchain_parse_lock_attr(ctx)
    ctx.file(
        "defs.bzl",
        _DEFS.format(
            name = ctx.original_name,
        )
    )
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            name = ctx.original_name,
            archives = json.encode_indent(lock.version_archives, indent = "    "),
            original_name = ctx.original_name
        ),
    )
    return ctx.repo_metadata(
        reproducible = True
    )

binary_toolchain_root_repo = repository_rule(
    implementation = _impl,
    doc = "Root binary toolchain repository",
    attrs = {
        "lock": attr.label(
            mandatory = True,
            doc = "Lock file label",
        ),
    },
)
