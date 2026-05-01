load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")

def docs_filegroup(
        name,
        srcs,
        visibility,
        deps = [],
        prefix_root = "content/docs/"):
    """
    Create a documentation filegroup
    """
    package_name = native.package_name()
    if package_name:
        package_name = "{}/".format(package_name)
    deps_normalized = []
    for dep in deps:
        if not ":" in dep:
            dep = "//{}{}:docs".format(package_name, dep)
        deps_normalized.append(dep)
    pkg_filegroup(
        name = name,
        srcs = [":{}.files".format(name)] + deps_normalized,
        visibility = visibility,
    )

    pkg_files(
        name = "{}.files".format(name),
        srcs = srcs,
        prefix = "{}{}".format(prefix_root, native.package_name()),
    )
