"""Documentation packaging helpers."""

load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files", "strip_prefix")

def docs_filegroup(
        name,
        srcs = None,
        visibility = None,
        deps = [],
        prefix_root = "content/docs/",
        prefix = None,
        preserve_paths = False):
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
        preserve_paths: Keep each source's path relative to the package
            instead of flattening it to its basename. Enable this when the
            sources sit in subdirectories that hold identically named files,
            which would otherwise collide on one destination path. Sources
            outside this package cannot be preserved and keep their flattened
            name.
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

    files_kwargs = {
        "name": "{}.files".format(name),
        "srcs": srcs,
        "prefix": prefix,
    }
    if preserve_paths:
        files_kwargs["strip_prefix"] = strip_prefix.from_pkg()
    pkg_files(**files_kwargs)
