"""Documentation packaging helpers."""

load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")

def docs_filegroup(
        name,
        srcs = None,
        visibility = None,
        deps = [],
        prefix_root = "content/docs/",
        prefix = None):
    """Creates a documentation filegroup.

    Args:
        name: Target name.
        srcs: Documentation files. Defaults to all Markdown files in the
            package.
        visibility: Optional visibility for the aggregate target.
        deps: Documentation targets to include. Bare relative package names
            are normalized to their `docs` targets.
        prefix_root: Default archive root used when `prefix` is omitted.
        prefix: Optional explicit archive prefix.
    """
    package_name = native.package_name()
    package_prefix = package_name
    if package_prefix:
        package_prefix = "{}/".format(package_prefix)

    deps_normalized = []
    for dep in deps:
        is_relative = not dep.startswith("//") and not dep.startswith("@")
        if is_relative and ":" not in dep:
            dep = "//{}{}:docs".format(package_prefix, dep)
        deps_normalized.append(dep)

    if srcs == None:
        srcs = native.glob(["*.md"])
    if prefix == None:
        prefix = "{}{}".format(prefix_root, package_name)

    aggregate_kwargs = {
        "name": name,
        "srcs": [":{}.files".format(name)] + deps_normalized,
    }
    if visibility != None:
        aggregate_kwargs["visibility"] = visibility

    pkg_filegroup(**aggregate_kwargs)

    pkg_files(
        name = "{}.files".format(name),
        srcs = srcs,
        prefix = prefix,
    )
