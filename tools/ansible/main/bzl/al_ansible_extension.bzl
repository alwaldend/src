load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

_BUILD = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_files", "strip_prefix")

alias(
    name = "pkg_files",
    actual = "{name}",
    visibility = ["//visibility:public"],
)

pkg_files(
    name = "{name}",
    srcs = glob(["**"]),
    prefix = "{prefix}",
    strip_prefix = strip_prefix.from_pkg(),
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.remote_collections:
            lock = json.decode(ctx.read(tag.lock))
            for data in lock.get("remote_collections", []):
                name = "{}_{}".format(tag.name, data["name"].replace("-", "_"))
                root_module_direct_deps.append(name)
                prefix = data["name"].replace("-", "/")
                http_archive(
                    name = name,
                    url = data["url"],
                    build_file_content = _BUILD.format(prefix = prefix, name = name),
                    integrity = data["integrity"],
                )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_ansible_extension = module_extension(
    implementation = _impl,
    doc = "Ansible extensions",
    tag_classes = {
        "remote_collections": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name prefix",
            ),
            "lock": attr.label(
                doc = "Ansible lock path",
                mandatory = True,
            ),
        }),
    },
)
