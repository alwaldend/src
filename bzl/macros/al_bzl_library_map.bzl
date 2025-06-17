load("@bazel_skylib//:bzl_library.bzl", "bzl_library")

def al_bzl_library_map(name, libs, **kwargs):
    """
    Create bzl_library targets from a map

    Args:
        name (str): combined bzl_library target name
        libs (dict[str, list[str]]): bzl_library names
        **kwargs: bzl_library kwargs
    """
    bzl_library(
        name = name,
        deps = libs.keys(),
        **kwargs
    )
    for lib_name, lib_deps in libs.items():
        lib_srcs = ["{}.bzl".format(lib_name)]
        native.exports_files(lib_srcs)
        bzl_library(
            name = lib_name,
            srcs = lib_srcs,
            deps = lib_deps,
            **kwargs
        )
