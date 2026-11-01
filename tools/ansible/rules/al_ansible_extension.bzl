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

_BUILD_REPO = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup")

alias(
    name = "collections",
    actual = "{name}",
    visibility = ["//visibility:public"],
)

pkg_filegroup(
    name = "{name}",
    srcs = {srcs},
    prefix = "collections/ansible_collections",
    visibility = ["//visibility:public"],
)
"""

def _repo_impl(ctx):
    ctx.file(
        "BUILD.bazel",
        _BUILD_REPO.format(
            name = ctx.original_name,
            srcs = json.encode([str(val) for val in ctx.attr.collections]),
        ),
    )
    return ctx.repo_metadata(
        reproducible = True,
    )

_repo = repository_rule(
    implementation = _repo_impl,
    doc = "Ansible repository",
    attrs = {
        "collections": attr.label_list(
            default = [],
            doc = "List of collection labels",
        ),
    },
)

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.lock:
            lock = json.decode(ctx.read(tag.lock))
            names = []
            for data in lock.get("remote_collections", []):
                name = "{}_{}".format(tag.name, data["name"].replace("-", "_"))
                names.append("@{}".format(name))
                root_module_direct_deps.append(name)
                prefix = data["name"].replace("-", "/")
                http_archive(
                    name = name,
                    url = data["url"],
                    build_file_content = _BUILD.format(prefix = prefix, name = name),
                    integrity = data["integrity"],
                )
            _repo(
                name = tag.name,
                collections = names,
            )
            root_module_direct_deps.append(tag.name)

    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_ansible_extension = module_extension(
    implementation = _impl,
    doc = "Ansible extensions",
    tag_classes = {
        "lock": tag_class({
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
