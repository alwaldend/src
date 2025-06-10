load("@rules_python//python:py_binary.bzl", "py_binary")
load("//bzl/vars:labels.bzl", "LABELS")

def al_py_binary_shell(
        name,
        deps = [],
        srcs = [],
        shell_type = "python",
        shell_label = LABELS.bazel_python_shell,
        **kwargs):
    """
    Create a py_binary target that allows you to run commands in proper python environment

    Args:
        name: target name
        deps: py_binary deps
        srcs: py_binary srcs
        shell_type: ${BAZEL_PYTHON_SHELL_TYPE}
        **kwargs: other py_binary kwargs
    """
    actual_kwargs = {"env": {}}
    actual_kwargs.update(kwargs)
    actual_kwargs["env"]["BAZEL_PYTHON_SHELL_TYPE"] = shell_type
    py_binary(
        name = name,
        srcs = [shell_label] + srcs,
        main = "bazel_python_shell.py",
        deps = deps,
        **actual_kwargs
    )
