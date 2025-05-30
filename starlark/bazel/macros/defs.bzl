"""
Bazel macros
"""

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@bazel_skylib//rules:diff_test.bzl", "diff_test")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@rules_python//python:pip.bzl", "compile_pip_requirements")
load("@rules_python//python:py_binary.bzl", "py_binary")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")
load("//starlark/bazel/rules:defs.bzl", "al_genrule_executable", "al_genrule_regular", "al_genrule_test")

BLACK_SRC = "//python:black"
ISORT_SRC = "//python:isort"
MYPY_SRC = "//python:mypy"
PYPROJECT_SRC = "//:pyproject.toml"
RUN_ARGS_SRC = "//shell/scripts:run-args-lib"
GO_SRC = "@rules_go//go"
FLAKE8_SRC = "//python:flake8"
EDITORCONFIG_SRC = "//:.editorconfig"
STYLUA_SRC = "@com-alwaldend-src-cargo//:stylua__stylua"
SHFMT_SRC = "@cc_mvdan_sh_v3//cmd/shfmt:shfmt"
STYLUA_CONFIG_SRC = "//lua:stylua.toml"
SHELLCHECK_SRC = "@com-github-koalaman-shellcheck//:shellcheck"
INSTALL_FILE_SRC = "//python/install-file:lib"
VISIBILITY_PUBLIC = ["//visibility:public"]
REPLACE_SECTION_SRC = "//python/replace-section"

# def al_go_checkers(
#         name,
#         srcs = [],
#         go_src = GO_SRC,
#         run_args_src = RUN_ARGS_SRC):
#     """
#     Generate golang checkers for
#     """
#     kwargs = {"srcs": [run_args_src], "data": srcs + [go_src]}
#     native.sh_binary(
#         name = "{}-stylua-fix".format(name),
#         args = stylua_args,
#         **stylua_kwargs
#     )
#     native.sh_test(
#         name = "{}-stylua-test".format(name),
#         args = [stylua_args[0], "--check"] + stylua_args[1:],
#         size = "small",
#         **stylua_kwargs
#     )

def al_run_tool(**kwargs):
    """
    Run _al_run_tool with test = False

    Args:
        **kwargs: kwargs for `_al_run_tool`
    """
    kwargs["test"] = False
    _al_run_tool(**kwargs)

def al_run_tool_test(**kwargs):
    """
    Run _al_run_tool with test = True

    Args:
        **kwargs: kwargs for `_al_run_tool`
    """
    kwargs["test"] = True
    _al_run_tool(**kwargs)

def _al_run_tool(name = "", tool = "", test = False, run_args_src = RUN_ARGS_SRC, **kwargs):
    """
    Generate a target to run a tool

    Args:
        name: Target name (required)
        tool: Tool label to run (required)
        test: If true, generate native.sh_test, use native.sh_binary otherwise
        **kwargs: kwargs for rules
    """
    kwargs["name"] = name
    kwargs["data"] = kwargs.get("data", []) + [tool]
    kwargs["args"] = ["$(rootpath {})".format(tool)] + kwargs.get("args", [])
    kwargs["srcs"] = [run_args_src]

    if test:
        native.sh_test(**kwargs)
    else:
        native.sh_binary(**kwargs)

def al_genrule(test = False, executable = False, **kwargs):
    """
    Generate al_genrule target

    Args:
        test: If set, use al_genrule_test
        executable: if set, use al_genrule_executable
        **kwargs: kwargs for the rule
    """
    if test:
        al_genrule_test(**kwargs)
    elif executable:
        al_genrule_executable(**kwargs)
    else:
        al_genrule_regular(**kwargs)

def al_lua_library(
        name,
        srcs = [],
        check = [],
        stylua_config_src = STYLUA_CONFIG_SRC,
        stylua_src = STYLUA_SRC,
        pkg_tar_kwargs = {},
        visibility = VISIBILITY_PUBLIC):
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
        "--config-path=$(rootpath {})".format(stylua_config_src),
    ] + ["$(rootpaths {})".format(src) for src in (check or srcs)]
    data = (check or srcs) + [stylua_config_src]
    pkg_tar(
        name = name,
        srcs = srcs,
        visibility = visibility,
        **pkg_tar_kwargs
    )
    al_run_tool(
        name = "{}-stylua-fix".format(name),
        tool = stylua_src,
        args = stylua_args,
        data = data,
    )
    al_run_tool_test(
        name = "{}-stylua-test".format(name),
        tool = stylua_src,
        args = [stylua_args[0], "--check"] + stylua_args[1:],
        size = "small",
        data = data,
    )

def al_bzl_library(
        name,
        deps = [],
        srcs = ["defs.bzl"],
        out = "api.md",
        input = "defs.bzl",
        replace_section_src = REPLACE_SECTION_SRC,
        run_args_src = RUN_ARGS_SRC,
        visibility = VISIBILITY_PUBLIC):
    """
    Generate stardoc and bzl_library targets

    Args:
        name: library name
        deps: library dependencies
        srcs: library sources
        out: generated md file
        input: input file
    """
    bzl_library(
        name = name,
        srcs = srcs,
        deps = deps,
        visibility = visibility,
    )
    stardoc(
        name = "{}-stardoc".format(name),
        out = out,
        input = input,
        deps = [":{}".format(name)],
        visibility = visibility,
    )
    al_run_tool(
        name = "{}-stardoc-copy".format(name),
        tool = replace_section_src,
        args = [
            "-i",
            "-s",
            "STARDOC",
            "-f",
            "$(rootpath :{}-stardoc)".format(name),
            "$(rootpath :README.md)",
        ],
        data = [":{}-stardoc".format(name), ":README.md"],
    )
    diff_test(
        name = "{}-stardoc-test".format(name),
        file1 = "{}-stardoc".format(name),
        file2 = ":{}".format(out),
        size = "small",
    )

def al_py_checkers(
        srcs = [],
        black_src = BLACK_SRC,
        isort_src = ISORT_SRC,
        mypy_src = MYPY_SRC,
        flake8_src = FLAKE8_SRC,
        pyproject_src = PYPROJECT_SRC):
    """
    Generate -fix and -test targets for python checkers

    Args:
        src: source files label
    """
    kwargs = {
        "pyproject_src": pyproject_src,
        "srcs": srcs,
    }
    al_py_isort(name = "isort", isort_src = isort_src, **kwargs)
    al_py_black(name = "black", black_src = black_src, **kwargs)
    al_py_mypy(name = "mypy", mypy_src = mypy_src, **kwargs)
    al_py_flake8(name = "flake8", flake8_src = flake8_src, **kwargs)

def al_py_flake8(name = None, srcs = None, flake8_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for isort
    """
    al_py_checker(
        name,
        disable_fix = True,
        tool = flake8_src,
        args_test = [
            "--toml-config=$(rootpath {})".format(pyproject_src),
        ] + ["$(rootpaths {})".format(src) for src in srcs],
        data = srcs + [pyproject_src],
    )

def al_py_isort(name = None, srcs = None, isort_src = None, pyproject_src = None):
    """
    Generate -fix and -test targets for isort
    """
    args = [
        "--settings-path=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = name,
        args_bin = args,
        tool = isort_src,
        args_test = [args[0], "--check-only"] + args[1:],
        data = srcs + [pyproject_src],
    )

def al_py_black(name = None, srcs = None, black_src = None, pyproject_src = None):
    """
    Generate -fix and -test targets for black
    """
    args = [
        "--config=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name,
        args_bin = args,
        tool = black_src,
        args_test = [args[0], "--check"] + args[1:],
        data = srcs + [pyproject_src],
    )

def al_py_mypy(name = None, srcs = None, mypy_src = None, pyproject_src = None):
    """
    Generate -fix and -test targets for mypy
    """
    args = [
        "--config-file=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name,
        args_bin = args,
        tool = mypy_src,
        args_test = [args[0], "--check"] + args[1:],
        data = srcs + [pyproject_src],
        disable_fix = True,
    )

def al_py_checker(
        name = None,
        tool = None,
        args_bin = None,
        args_test = None,
        test_size = "small",
        disable_fix = False,
        **kwargs):
    """
    Create -fix and -test targets for a python checker

    Args:
        name: Name prefix
        tool: Tool label
        args_bin: Args for the binary target
        args_test: Args for the test
        disable_fix: If set, do not create fix target
        **kwargs: Kwargs for both targets
    """
    if not disable_fix:
        al_run_tool(
            name = "{}-fix".format(name),
            args = args_bin,
            tool = tool,
            **kwargs
        )
    al_run_tool_test(
        name = "{}-test".format(name),
        args = args_test,
        tool = tool,
        size = test_size,
        **kwargs
    )

def al_install_file(
        name = None,
        args = [],
        install_file_src = INSTALL_FILE_SRC,
        visibility = VISIBILITY_PUBLIC,
        **py_binary_kwargs):
    py_binary(
        name = name,
        srcs = [install_file_src],
        main = "install_file.py",
        args = args,
        visibility = visibility,
        **py_binary_kwargs
    )

def al_sh_script(
        name = "",
        shfmt_src = SHFMT_SRC,
        editorconfig_src = EDITORCONFIG_SRC,
        shellcheck_src = SHELLCHECK_SRC,
        run_args_src = RUN_ARGS_SRC,
        visibility = VISIBILITY_PUBLIC,
        **common_kwargs):
    """
    Generate sh_library and sh_binary targets

    Args:
        name: name
        **common_kwargs: kwargs for both targets
    """
    lib_name = "{}-lib".format(name)
    native.sh_library(
        name = lib_name,
        srcs = ["{}.sh".format(name)],
        visibility = visibility,
        **common_kwargs
    )
    native.sh_binary(
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
    native.sh_binary(
        name = "{}-shfmt-fix".format(name),
        args = [shfmt_args[0], "--write"] + shfmt_args[1:],
        **shfmt_kwargs
    )
    native.sh_test(
        name = "{}-shfmt-test".format(name),
        args = [shfmt_args[0], "--diff"] + shfmt_args[1:],
        size = "small",
        **shfmt_kwargs
    )
    native.sh_test(
        name = "{}-shellcheck-test".format(name),
        args = ["$(rootpath {})".format(shellcheck_src), "$(rootpath {})".format(lib_name)],
        size = "small",
        srcs = [run_args_src],
        data = [lib_name, shellcheck_src],
    )

def al_compile_pip_requirements_combined(name = "", srcs = [], **compile_kwargs):
    """
    Compile seveal requirement files

    Args:
        name: name
        srcs: requirement files to combine
        **compile_kwargs: kwargs for `compile_pip_requirements`
    """
    al_combine_files(
        name = "{}-src".format(name),
        srcs = srcs,
    )
    compile_pip_requirements(
        name = name,
        src = ":{}-src".format(name),
        **compile_kwargs
    )

def al_combine_files(name = "", srcs = [], **genrule_kwargs):
    """
    Combine several files into one

    Args:
        name: target name
        srcs: sources to combine
    """
    native.genrule(
        name = name,
        srcs = srcs,
        outs = ["{}-output".format(name)],
        cmd = "cat $(SRCS) >$(@)",
        **genrule_kwargs
    )

def al_apply_patches(
        name = "patched",
        src = "",
        patches = "",
        visibility = VISIBILITY_PUBLIC):
    """
    Apply patches

    Args:
        name: target name
        src: source archive label
        patches: patches label
        visibility: visibility
    """
    native.genrule(
        name = name,
        srcs = [src, patches],
        outs = ["{}.tar".format(name)],
        visibility = visibility,
        cmd = """
            set -eux
            mkdir _src _patches
            cp $(rootpaths {patches}) _patches/
            tar -xf "$(rootpath {src})" -C _src --strip-components 2
            cd _src
            git apply -v ../_patches/*
            cd ../
            tar -cf "$(@)" -C "_src"
        """.format(src = src, patches = patches),
    )

def al_pkg_tar_combined(name = None, tars = [], strip_components = 2, **genrule_kwargs):
    """
    Combine several tars into one
    """
    cmd = "set -eux"
    out = "{}.tar".format(name)
    cmd += "\n".join(["""
        mkdir -p '{dir}'
        tar -xf $(execpath {label}) --strip-components '{strip_components}' -C '{dir}'
    """.format(strip_components = strip_components, **tar) for tar in tars])
    cmd += """
        tar -cf $(execpath {output}) {dirs}
    """.format(output = out, dirs = " ".join(["'{}'".format(tar["dir"]) for tar in tars]))
    native.genrule(
        name = name,
        srcs = [tar["label"] for tar in tars],
        outs = [out],
        cmd = cmd,
        **genrule_kwargs
    )

def al_genrule_src(name = "src", patterns = ["**"], visibility = VISIBILITY_PUBLIC):
    """
    Create a filegroup and a genrule generating a tar archive
    """
    filegroup = name + "-filegroup"
    native.filegroup(
        name = filegroup,
        srcs = native.glob(patterns),
        visibility = visibility,
    )
    native.genrule(
        name = name,
        srcs = native.glob(patterns),
        outs = ["{}.tar".format(name)],
        cmd = "tar -cf ${@} ${<}",
        visibility = visibility,
    )

def al_py_binary_shell(name = "", cmd = "", deps = [], srcs = [], shell_type = "python", **kwargs):
    """
    Create a shell with linked python deps
    """
    actual_kwargs = {"env": {}}
    actual_kwargs.update(kwargs)
    actual_kwargs["env"]["BAZEL_PYTHON_SHELL_TYPE"] = shell_type
    py_binary(
        name = name,
        srcs = ["//python/bazel-python-shell:library"] + srcs,
        main = "bazel_python_shell.py",
        deps = deps,
        **actual_kwargs
    )

def al_genrule_with_wheels(wheels = None, srcs = [], cmd = [], **kwargs):
    """
    Regular genrule with wheels added to ${PYTHONPATH}
    """
    wheels_env_script = ""
    if wheels:
        wheel_paths = []
        for wheel in wheels:
            wheel_paths.append("$(rootpath {})".format(wheel))
            srcs = srcs + [wheel]
            wheels_env_script = 'export PYTHONPATH="{}"\n'.format(":".join(wheel_paths))
        wheel_paths.append(["$${PYTHONPATH:-}"])
    kwargs["cmd"] = wheels_env_script + cmd
    kwargs["srcs"] = srcs

    native.genrule(**kwargs)
