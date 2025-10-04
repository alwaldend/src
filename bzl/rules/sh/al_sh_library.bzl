load("@rules_shell//shell:sh_binary.bzl", "sh_binary")
load("@rules_shell//shell:sh_library.bzl", "sh_library")
load("@rules_shell//shell:sh_test.bzl", "sh_test")

def al_sh_library(
        name,
        shfmt_src = "@cc_mvdan_sh_v3//cmd/shfmt:shfmt",
        editorconfig_src = "//:editorconfig",
        shellcheck_src = "@com-github-koalaman-shellcheck-linux-x86_64//:bin",
        run_args_src = "//sh/scripts:run-args.lib",
        visibility = ["//visibility:public"],
        **common_kwargs):
    """
    Create targets for a shell library

    Targets:
    - ${name}-shfmt-fix: executable to run shfmt
    - ${name}-shfmt-test: test whether the script is formatted
    - ${name}-shellcheck-test: shellcheck test

    Args:
        name: target name
        **common_kwargs: kwargs for both targets
    """
    lib_name = "{}.lib".format(name)
    sh_library(
        name = lib_name,
        srcs = ["{}.sh".format(name)],
        visibility = visibility,
        **common_kwargs
    )
    sh_binary(
        name = name,
        srcs = [":{}".format(lib_name)],
        visibility = visibility,
        **common_kwargs
    )
    shfmt_args = [
        "$(rootpath {})".format(shfmt_src),
        "$(rootpath {})".format(lib_name),
    ]
    shfmt_kwargs = {
        "srcs": [run_args_src],
        "data": [lib_name, shfmt_src, editorconfig_src],
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
        args = ["$(rootpath {})".format(shellcheck_src), "$(rootpath {})".format(lib_name)],
        size = "small",
        srcs = [run_args_src],
        data = [lib_name, shellcheck_src],
    )
