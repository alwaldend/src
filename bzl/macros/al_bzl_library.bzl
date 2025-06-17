load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")
load("//bzl/rules:al_write_script.bzl", "al_write_script")

def al_bzl_library(
        name,
        visibility,
        src = None,
        deps = []):
    """
    Generate targets for a bzl library

    Targets:
    - ${name}: bzl_library
    - ${name}-stardoc: stardoc target (if src is passed)

    Args:
        name (str): library name
        src: (Optinal[str]): bzl source file
        visibility (list[str]): visibility
        deps (list[str]): bzl_library deps
    """
    stardoc_out = "{}.md".format(name)
    bzl_library(
        name = name,
        srcs = [src] if src else [],
        deps = deps,
        visibility = visibility,
    )
    if src:
        stardoc(
            name = "{}-stardoc".format(name),
            out = stardoc_out,
            input = src,
            deps = [name],
            visibility = visibility,
        )
