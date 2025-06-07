load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")
load("//bzl/rules:al_write_script.bzl", "al_write_script")
load("//bzl/vars:vars.bzl", "LABELS")

def al_bzl_library(
        name,
        src = None,
        stardoc_out = None,
        deps = [],
        visibility = ["//visibility:public"],
        replace_section_label = LABELS.replace_section):
    """
    Generate targets for a bzl library

    Targets:
        ${name}: bzl_library
        ${name}-stardoc: stardoc target

    Args:
        name: library name
        src: source file (default: `${name}.bzl`)
        stadoc_out: output stardoc (default: `${name}.md`)
        deps: bzl_library deps
    """
    if not src:
        src = "{}.bzl".format(name)
    if not stardoc_out:
        stardoc_out = "{}.md".format(name)
    bzl_library(
        name = name,
        srcs = [src],
        deps = deps,
        visibility = visibility,
    )
    stardoc(
        name = "{}-stardoc".format(name),
        out = stardoc_out,
        input = src,
        deps = [name],
        visibility = visibility,
    )
