load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@bazel_skylib//rules:diff_test.bzl", "diff_test")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@rules_python//python:pip.bzl", "compile_pip_requirements")
load("@rules_python//python:py_binary.bzl", "py_binary")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")

BLACK_SRC = "//python:black"
ISORT_SRC = "//python:isort"
MYPY_SRC = "//python:mypy"
PYPROJECT_SRC = "//:pyproject"
RUN_ARGS_SRC = "//shell/scripts:run-args-lib"
STYLUA_SRC = "@cargo//:stylua__stylua"
STYLUA_CONFIG_SRC = "//lua:stylua.toml"
INSTALL_FILE_SRC = "//python/install-file:lib"
VISIBILITY_PUBLIC = ["//visibility:public"]
REPLACE_SECTION_SRC = "//python/replace-section"

def al_lua_library(
        name,
        srcs = [],
        check = [],
        run_args_src = RUN_ARGS_SRC,
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
        "$(rootpath {})".format(stylua_src),
        "--config-path=$(rootpath {})".format(stylua_config_src),
    ] + ["$(rootpaths {})".format(src) for src in (check or srcs)]
    stylua_kwargs = {
        "srcs": [run_args_src],
        "data": (check or srcs) + [stylua_config_src, stylua_src],
    }
    pkg_tar(
        name = name,
        srcs = srcs,
        visibility = visibility,
        **pkg_tar_kwargs
    )
    native.sh_binary(
        name = "{}-stylua-fix".format(name),
        args = stylua_args,
        **stylua_kwargs
    )
    native.sh_test(
        name = "{}-stylua-test".format(name),
        args = [stylua_args[0], "--check"] + stylua_args[1:],
        size = "small",
        **stylua_kwargs
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
    native.sh_binary(
        name = "{}-stardoc-copy".format(name),
        srcs = [run_args_src],
        args = [
            "$(rootpath {})".format(replace_section_src),
            "-i",
            "-s",
            "STARDOC",
            "-f",
            "$(rootpath :{}-stardoc)".format(name),
            "$(rootpath :README.md)",
        ],
        data = [":{}-stardoc".format(name), replace_section_src, ":README.md"],
    )
    # diff_test(
    #     name = "{}-stardoc-test".format(name),
    #     file1 = "{}-stardoc".format(name),
    #     file2 = ":{}".format(out),
    # )

def al_py_checkers(
        srcs = [],
        black_src = BLACK_SRC,
        isort_src = ISORT_SRC,
        mypy_src = MYPY_SRC,
        pyproject_src = PYPROJECT_SRC,
        run_args_src = RUN_ARGS_SRC):
    """
    Generate -fix and -test targets for python checkers

    Args:
        src: source files label
    """
    kwargs = {
        "pyproject_src": pyproject_src,
        "srcs": srcs,
        "run_args_src": run_args_src,
    }
    al_py_isort(name = "isort", isort_src = isort_src, **kwargs)
    al_py_black(name = "black", black_src = black_src, **kwargs)
    al_py_mypy(name = "mypy", mypy_src = mypy_src, **kwargs)

def al_py_isort(name = None, srcs = None, isort_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for isort
    """
    args = [
        "$(rootpath {})".format(isort_src),
        "--settings-path=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check-only"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [isort_src, run_args_src, pyproject_src],
    )

def al_py_black(name = None, srcs = None, black_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for black
    """
    args = [
        "$(rootpath {})".format(black_src),
        "--config=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [black_src, run_args_src, pyproject_src],
    )

def al_py_mypy(name = None, srcs = None, mypy_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for mypy
    """
    args = [
        "$(rootpath {})".format(mypy_src),
        "--config-file=$(rootpath {})".format(pyproject_src),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [mypy_src, run_args_src, pyproject_src],
        disable_fix = True,
    )

def al_py_checker(
        name = None,
        args_bin = None,
        args_test = None,
        disable_fix = False,
        **kwargs):
    """
    Create -fix and -test targets for a python checker
    """
    if not disable_fix:
        native.sh_binary(
            name = "{}-fix".format(name),
            args = args_bin,
            **kwargs
        )
    native.sh_test(
        name = "{}-test".format(name),
        args = args_test,
        size = "small",
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

def al_sh_script(name = "", visibility = VISIBILITY_PUBLIC, **common_kwargs):
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
        tar -xf $(rootpath {label}) --strip-components '{strip_components}' -C '{dir}'
    """.format(strip_components = strip_components, **tar) for tar in tars])
    cmd += """
        tar -cf $(rootpath {output}) {dirs}
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
