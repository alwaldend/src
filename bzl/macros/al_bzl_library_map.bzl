load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")

def al_bzl_library_map(name, visibility, libs = {}, deps = [], **kwargs):
    """
    Create bzl_library targets from a map

    Args:
        name (str): combined bzl_library target name
        libs (dict[str, list[str]]): bzl_library names
        deps (list[str]): other al_bzl_library_map targets
        visibiliy: visibiliy
        **kwargs: bzl_library kwargs
    """
    bzl_library(
        name = name,
        deps = libs.keys() + deps,
        visibility = visibility,
        **kwargs
    )
    pkg_tar(
        name = "{}-stardoc".format(name),
        srcs = ["{}-stardoc".format(lib) for lib in libs.keys()],
        out = "{}-stardoc.tar".format(name),
        package_dir = native.package_name().split("/")[-1],
        deps = ["{}-stardoc".format(dep) for dep in deps],
        visibility = visibility,
    )
    for lib_name, lib_deps in libs.items():
        bzl_library(
            name = lib_name,
            srcs = ["{}.bzl".format(lib_name)],
            deps = lib_deps,
            visibility = visibility,
            **kwargs
        )
        stardoc(
            name = "{}-stardoc".format(lib_name),
            out = "{}-stardoc.md".format(lib_name),
            input = "{}.bzl".format(lib_name),
            deps = [lib_name],
        )
