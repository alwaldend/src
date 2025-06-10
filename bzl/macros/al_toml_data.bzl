load("@rules_shell//shell:sh_test.bzl", "sh_test")
load("//bzl/vars:labels.bzl", "LABELS")
load(":al_run_tool.bzl", "al_run_tool")

def al_toml_data(name, srcs, visibility = ["//visibility:public"], tomlv_label = LABELS.tomlv, size = "small", **kwargs):
    """
    Create toml data targets

    Targets:
        ${name}: filegroup
        ${name}-test-load: al_run_tool test whether the file can be loaded

    Args:
        name: Target name
        srcs: Toml files
        visibility: visibility
        **kwargs: kwargs
    """
    native.filegroup(
        name = name,
        srcs = srcs,
        visibility = visibility,
    )
    al_run_tool(
        name = "{}-load-test".format(name),
        test = True,
        args = ["$(execpaths {})".format(src) for src in srcs],
        size = "small",
        data = srcs,
        visibility = visibility,
        tool = tomlv_label,
        **kwargs
    )
