load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@rules_python//python:pip.bzl", "compile_pip_requirements")
load("@rules_python//python:py_binary.bzl", "py_binary")

BLACK_SRC = "//python:black"
ISORT_SRC = "//python:isort"
MYPY_SRC = "//python:mypy"
PYPROJECT_SRC = "//:pyproject"
RUN_ARGS_SRC = "//shell/scripts:run-args-lib"

def py_checkers(
        srcs = [],
        black_src = BLACK_SRC,
        isort_src = ISORT_SRC,
        mypy_src = MYPY_SRC,
        pyproject_src = PYPROJECT_SRC,
        run_args_src = RUN_ARGS_SRC):
    """
    Generate -fix and -test targets for isort

    Args:
        src: source files label
    """
    kwargs = {
        "pyproject_src": pyproject_src,
        "srcs": srcs,
        "run_args_src": run_args_src,
    }
    py_isort(name = "isort", isort_src = isort_src, **kwargs)
    py_black(name = "black", black_src = black_src, **kwargs)
    py_mypy(name = "mypy", mypy_src = mypy_src, **kwargs)

def py_isort(name = None, srcs = None, isort_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for isort
    """
    args = [
        "$(execpath {})".format(isort_src),
        "--settings-path=$(execpath {})".format(pyproject_src),
    ] + ["$(execpaths {})".format(src) for src in srcs]
    py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check-only"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [isort_src, run_args_src, pyproject_src],
    )

def py_black(name = None, srcs = None, black_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for black
    """
    args = [
        "$(execpath {})".format(black_src),
        "--config=$(execpath {})".format(pyproject_src),
    ] + ["$(execpaths {})".format(src) for src in srcs]
    py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [black_src, run_args_src, pyproject_src],
    )

def py_mypy(name = None, srcs = None, mypy_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for mypy
    """
    args = [
        "$(execpath {})".format(mypy_src),
        "--config-file=$(execpath {})".format(pyproject_src),
    ] + ["$(execpaths {})".format(src) for src in srcs]
    py_checker(
        name,
        args_bin = args,
        args_test = [args[0], "--check"] + args[1:],
        srcs = [run_args_src],
        data = srcs + [mypy_src, run_args_src, pyproject_src],
        disable_fix = True,
    )

def py_checker(name = None, args_bin = None, args_test = None, disable_fix = False, **kwargs):
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

def install_file(name = "", args = [], install_file_src = "//python/install-file:lib", visibility = ["//visibility:public"], **py_binary_kwargs):
    py_binary(
        name = name,
        srcs = [install_file_src],
        main = "install_file.py",
        args = args,
        visibility = visibility,
        **py_binary_kwargs
    )

def sh_script(name = "", visibility = ["//visibility:public"], **common_kwargs):
    """
    Generate sh_library and sh_binary targets

    Args:
        name: name
        **common_kwargs: kwargs for both targets
    """
    lib_name = "{}-lib".format(name)
    native.sh_library(name = lib_name, srcs = ["{}.sh".format(name)], visibility = visibility, **common_kwargs)
    native.sh_binary(name = name, srcs = [":{}".format(lib_name)], visibility = visibility, **common_kwargs)

def compile_pip_requirements_combined(name = "", srcs = [], **compile_kwargs):
    """
    Compile seveal requirement files

    Args:
        name: name
        srcs: requirement files to combine
        **compile_kwargs: kwargs for `compile_pip_requirements`
    """
    combine_files(
        name = "{}-src".format(name),
        srcs = srcs,
    )
    compile_pip_requirements(
        name = name,
        src = ":{}-src".format(name),
        **compile_kwargs
    )

def combine_files(name = "", srcs = [], **genrule_kwargs):
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

def apply_patches(name = "patched", src = "", patches = "", visibility = ["//visibility:public"]):
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
            cp $(execpaths {patches}) _patches/
            tar -xf "$(execpath {src})" -C _src --strip-components 2
            cd _src
            git apply -v ../_patches/*
            cd ../
            tar -cf "$(@)" -C "_src"
        """.format(src = src, patches = patches),
    )

def pkg_tar_combined(name = None, tars = [], strip_components = 2, **genrule_kwargs):
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

def genrule_src(name = "src", patterns = ["**"], visibility = ["//visibility:public"]):
    """
    Create "{name}-filegroup" filegroup and "{name}"
    genrule to create a tar archive out of it
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

def py_binary_shell(name = "", cmd = "", deps = [], srcs = [], shell_type = "python", **kwargs):
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

def genrule_with_wheels(wheels = None, srcs = [], cmd = [], **kwargs):
    """
    Regular genrule with wheels added to ${PYTHONPATH}
    """
    wheels_env_script = ""
    if wheels:
        wheel_paths = []
        for wheel in wheels:
            wheel_paths.append("$(execpath {})".format(wheel))
            srcs = srcs + [wheel]
            wheels_env_script = 'export PYTHONPATH="{}"\n'.format(":".join(wheel_paths))
        wheel_paths.append(["$${PYTHONPATH:-}"])
    kwargs["cmd"] = wheels_env_script + cmd
    kwargs["srcs"] = srcs

    native.genrule(**kwargs)
