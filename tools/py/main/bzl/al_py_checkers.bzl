load(":al_py_checker.bzl", "al_py_checker")

def al_py_checkers(
        name,
        srcs,
        isort_label = "//tools/isort",
        black_label = "//tools/black",
        mypy_label = "//tools/mypy",
        flake8_label = "//tools/flake8",
        pyproject_label = "//:pyproject"):
    """
    Generate -fix and -test targets for python checkers

    Args:
        srcs: list of source file labels
    """
    isort_args = [
        "--settings-path=$(rootpath {})".format(pyproject_label),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = "{}.isort".format(name),
        args_bin = isort_args,
        tool = isort_label,
        args_test = [isort_args[0], "--check-only"] + isort_args[1:],
        data = srcs + [pyproject_label],
    )

    black_args = [
        "--config=$(rootpath {})".format(pyproject_label),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = "black".format(name),
        args_bin = black_args,
        tool = black_label,
        args_test = [black_args[0], "--check"] + black_args[1:],
        data = srcs + [pyproject_label],
    )

    mypy_args = [
        "--config-file=$(rootpath {})".format(pyproject_label),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = "mypy".format(name),
        args_bin = mypy_args,
        tool = mypy_label,
        args_test = [mypy_args[0], "--check"] + mypy_args[1:],
        data = srcs + [pyproject_label],
        disable_fix = True,
    )

    al_py_checker(
        name = "flake8".format(name),
        disable_fix = True,
        tool = flake8_label,
        args_test = [
            "--toml-config=$(rootpath {})".format(pyproject_label),
        ] + ["$(rootpaths {})".format(src) for src in srcs],
        data = srcs + [pyproject_label],
    )
