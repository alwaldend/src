load("@bazel_skylib//rules:native_binary.bzl", "native_binary", "native_test")
load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("@rules_shell//shell:sh_binary.bzl", "sh_binary")
load("@rules_shell//shell:sh_test.bzl", "sh_test")
load("//bzl/vars:labels.bzl", "LABELS")

def al_run_tool(name, tool, executable = False, test = False, **kwargs):
    """
    Generate either native_test, native_binary, or run_binary target

    Args:
        name: Target name (required)
        tool: Tool label to run (required)
        test: If True, generate native_test
        executable: If True, generate native_binary
        **kwargs: kwargs for rules
    """
    if test:
        native_test(name = name, src = tool, **kwargs)
    elif executable:
        native_binary(name = name, src = tool, **kwargs)
    else:
        kwargs["srcs"] += kwargs.pop("data")
        run_binary(name = name, tool = tool, **kwargs)
