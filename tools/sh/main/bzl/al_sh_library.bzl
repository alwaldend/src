load("@rules_shell//shell:sh_binary.bzl", "sh_binary")
load("@rules_shell//shell:sh_library.bzl", "sh_library")
load("@rules_shell//shell:sh_test.bzl", "sh_test")

def al_sh_library(
        name,
        shfmt_src = "//tools/shfmt",
        editorconfig_src = "//:editorconfig",
        shellcheck_src = "//tools/shellcheck",
        run_args_src = "//tools/sh/main/sh:run_args_lib",
        visibility = ["//:__subpackages__"],
        test_data = [],
        **sh_kwargs):
    """
    Create targets for a shell library

    Targets:
    - ${name}.shfmt_fix: executable to run shfmt
    - ${name}.shfmt_test: test whether the script is formatted
    - ${name}.shellcheck_test: shellcheck test

    Args:
        name: target name
        **sh_kwargs: kwargs for sh targets
    """
    sh_library(
        name = name,
        visibility = visibility,
        **sh_kwargs
    )
    shfmt_args = [
        "$(rootpath {})".format(shfmt_src),
        "$(rootpath {})".format(name),
    ]
    shfmt_kwargs = {
        "srcs": [run_args_src],
        "data": [name, shfmt_src, editorconfig_src] + test_data,
    }
    sh_binary(
        name = "{}.shfmt_fix".format(name),
        args = [shfmt_args[0], "--write"] + shfmt_args[1:],
        **shfmt_kwargs
    )
    sh_test(
        name = "{}.shfmt_test".format(name),
        args = [shfmt_args[0], "--diff"] + shfmt_args[1:],
        size = "small",
        **shfmt_kwargs
    )
    sh_test(
        name = "{}.shellcheck_test".format(name),
        args = [
            "$(rootpath {})".format(shellcheck_src),
            "-x",
            "$(rootpath {})".format(name),
        ],
        size = "small",
        srcs = [run_args_src],
        data = [name, shellcheck_src] + test_data,
    )
