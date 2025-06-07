load("@rules_shell//shell:sh_binary.bzl", "sh_binary")
load("@rules_shell//shell:sh_library.bzl", "sh_library")
load("@rules_shell//shell:sh_test.bzl", "sh_test")
load("//bzl/vars:vars.bzl", "LABELS")

def al_sh_library(
        name,
        shfmt_src = LABELS.shfmt,
        editorconfig_src = LABELS.editorconfig,
        shellcheck_src = LABELS.shellcheck,
        run_args_src = LABELS.run_args,
        visibility = ["//visibility:public"],
        **common_kwargs):
    """
    Create targets for a shell library

    Targets:
        ${name}-shfmt-fix: executable to run shfmt
        ${name}-shfmt-test: test whether the script is formatted
        ${name}-shellcheck-test: shellcheck test

    Args:
        name: target name
        **common_kwargs: kwargs for both targets
    """
    lib_name = "{}-lib".format(name)
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
        name = "{}-shfmt-fix".format(name),
        args = [shfmt_args[0], "--write"] + shfmt_args[1:],
        **shfmt_kwargs
    )
    sh_test(
        name = "{}-shfmt-test".format(name),
        args = [shfmt_args[0], "--diff"] + shfmt_args[1:],
        size = "small",
        **shfmt_kwargs
    )
    sh_test(
        name = "{}-shellcheck-test".format(name),
        args = ["$(rootpath {})".format(shellcheck_src), "$(rootpath {})".format(lib_name)],
        size = "small",
        srcs = [run_args_src],
        data = [lib_name, shellcheck_src],
    )
