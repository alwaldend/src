load("@rules_pkg//pkg:mappings.bzl", "pkg_files")
load("//tools/run_tool/main/bzl:al_run_tool.bzl", "al_run_tool")

def al_lua_library(
        name,
        srcs,
        check = [],
        stylua_config_label = "//tools/stylua:stylua_config",
        stylua_label = "//tools/stylua",
        pkg_kwargs = {},
        visibility = ["//:__subpackages__"]):
    """
    Generate targets for a lua library

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
    pkg_files(
        name = name,
        srcs = srcs,
        visibility = visibility,
        **pkg_kwargs
    )
    al_run_tool(
        name = "{}.stylua_fix".format(name),
        executable = True,
        tool = stylua_label,
        args = stylua_args,
        data = data,
    )
    al_run_tool(
        name = "{}.stylua_test".format(name),
        test = True,
        tool = stylua_label,
        args = [stylua_args[0], "--check"] + stylua_args[1:],
        size = "small",
        data = data,
    )
