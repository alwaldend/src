load("@rules_pkg//pkg/private/tar:tar.bzl", "pkg_tar")
load("//bzl/vars:vars.bzl", "LABELS")
load(":al_run_tool.bzl", "al_run_tool")

def al_lua_library(
        name,
        srcs,
        check = [],
        stylua_config_label = LABELS.stylua_config,
        stylua_label = LABELS.stylua,
        pkg_tar_kwargs = {},
        visibility = ["//visibility:public"]):
    """
    Generate targets for a lua library

    Targets:
        ${name}: pkg_tar
        ${name}-stylua-fix: al_run_tool executable
        ${name}-stylua-test: al_run_tool test

    Args:
        name: library name
        srcs: library sources
        check: if set, only these files will be checked
        visibility: visibility
        pkg_tar_kwargs: kwargs for pkg_tar
    """
    stylua_args = [
        "--config-path=$(rootpath {})".format(stylua_config_label),
    ] + ["$(rootpaths {})".format(src) for src in (check or srcs)]
    data = (check or srcs) + [stylua_config_label]
    pkg_tar(
        name = name,
        srcs = srcs,
        visibility = visibility,
        **pkg_tar_kwargs
    )
    al_run_tool(
        name = "{}-stylua-fix".format(name),
        executable = True,
        tool = stylua_label,
        args = stylua_args,
        data = data,
    )
    al_run_tool(
        name = "{}-stylua-test".format(name),
        test = True,
        tool = stylua_label,
        args = [stylua_args[0], "--check"] + stylua_args[1:],
        size = "small",
        data = data,
    )
