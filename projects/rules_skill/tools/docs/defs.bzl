"""Documentation packaging helpers for the standalone module workspace."""

load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")

# Keep this fallback synchronized with //tools/docs:defs.bzl. The module is
# also a standalone workspace, where depending on the parent module would form
# a dependency cycle solely to build its documentation.
def docs_filegroup(
        name,
        srcs,
        visibility,
        deps = [],
        prefix_root = "content/docs/",
        prefix = None):
    """Creates a documentation filegroup."""
    package_name = native.package_name()
    if package_name:
        package_name = "{}/".format(package_name)
    deps_normalized = []
    for dep in deps:
        is_relative = not dep.startswith("//") and not dep.startswith("@")
        if is_relative and not ":" in dep:
            dep = "//{}{}:docs".format(package_name, dep)
        deps_normalized.append(dep)
    if prefix == None:
        prefix = "{}{}".format(prefix_root, native.package_name())
    pkg_filegroup(
        name = name,
        srcs = [":{}.files".format(name)] + deps_normalized,
        visibility = visibility,
    )

    pkg_files(
        name = "{}.files".format(name),
        srcs = srcs,
        prefix = prefix,
    )
