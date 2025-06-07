load("@rules_python//python:py_binary.bzl", "py_binary")
load("//bzl/vars:vars.bzl", "LABELS")

def al_install_file(
        name,
        args = [],
        install_file_label = LABELS.install_file,
        visibility = ["//visibility:public"],
        **py_binary_kwargs):
    """
    Create py_binary target to install file

    Args:
        name: target name
        args: install-file args
    """
    py_binary(
        name = name,
        srcs = [install_file_label],
        main = "install_file.py",
        args = args,
        visibility = visibility,
        **py_binary_kwargs
    )
