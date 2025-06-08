load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load(":al_bzl_library.bzl", "al_bzl_library")
load(":al_md_data.bzl", "al_md_data")

def al_bzl_library_map(name, libs, common_deps = [], visibility = ["//visibility:public"]):
    """
    Create al_bzl_library targets from a map

    Targets:
        ${name}-stardoc: all stardoc markdown
        ${name}: all bzl_library targets

    Args:
        name: al_md_data name
        libs: map of al_bzl_library, keys are names, values are kwargs for
            al_bzl_library
        common_deps: list of common blz_library deps
        visibility: visibility
    """
    for key, config in libs.items():
        al_bzl_library(
            name = key,
            deps = common_deps + libs[key].get("deps", []),
            visibility = libs[key].get("visibility", visibility),
        )

    bzl_library(
        name = name,
        srcs = libs.keys(),
        visibility = visibility,
    )
    al_md_data(
        name = "{}-readme".format(name),
        srcs = [":README.md"],
        visibility = visibility,
    )
    al_md_data(
        name = "{}-stardoc".format(name),
        srcs = ["{}.md".format(key) for key in libs.keys()] + ["{}-readme".format(name)],
        visibility = visibility,
    )
