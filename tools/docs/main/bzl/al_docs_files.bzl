load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files", "strip_prefix")

def al_docs_files(
        name,
        visibility,
        srcs = [],
        deps = [],
        prefix = None,
        root_prefix = "content/docs",
        renames = {"README.md": "_index.md"}):
    """
    Macro that creates docs files

    Args:
        name: name
        visibility: visibility
        srcs: docs sources
        deps: other pkg_files rules
        prefix: prefix for pkg_filegroup
        renames: renames for pkg_files
    """
    pkg_filegroup(
        name = name,
        srcs = [":{}.files".format(name)] + deps,
        visibility = visibility,
    )
    pkg_files(
        name = "{}.files".format(name),
        srcs = srcs,
        prefix = "{}/{}".format(root_prefix, prefix) if prefix else root_prefix,
        strip_prefix = strip_prefix.from_pkg(),
        renames = renames,
    )
