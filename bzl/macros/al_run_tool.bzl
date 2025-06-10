load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("@rules_shell//shell:sh_binary.bzl", "sh_binary")
load("@rules_shell//shell:sh_test.bzl", "sh_test")
load("//bzl/vars:labels.bzl", "LABELS")

def al_run_tool(name, tool, executable = False, test = False, run_args_label = LABELS.run_args, **kwargs):
    """
    Generate either sh_test, sh_binary, or run_binary target

    Args:
        name: Target name (required)
        tool: Tool label to run (required)
        test: If True, generate sh_test
        executable: If True, generate sh_binary
        **kwargs: kwargs for rules
    """
    kwargs["name"] = name
    kwargs["data"] = kwargs.get("data", []) + [tool]
    kwargs["args"] = ["$(rootpath {})".format(tool)] + kwargs.get("args", [])
    kwargs["srcs"] = [run_args_label]

    if test:
        sh_test(**kwargs)
    elif executable:
        sh_binary(**kwargs)
    else:
        kwargs["srcs"] += kwargs.pop("data")
        kwargs["tool"] = tool
        kwargs["args"] = kwargs["args"][1:]
        run_binary(**kwargs)
