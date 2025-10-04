load(":al_py_checker.bzl", "al_py_checker")

def al_py_checkers(
        srcs,
        isort_label = "//py:isort",
        black_label = "//py:black",
        mypy_label = "//py:mypy",
        flake8_label = "//py:flake8",
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
        name = "isort",
        args_bin = isort_args,
        tool = isort_label,
        args_test = [isort_args[0], "--check-only"] + isort_args[1:],
        data = srcs + [pyproject_label],
    )

    black_args = [
        "--config=$(rootpath {})".format(pyproject_label),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = "black",
        args_bin = black_args,
        tool = black_label,
        args_test = [black_args[0], "--check"] + black_args[1:],
        data = srcs + [pyproject_label],
    )

    mypy_args = [
        "--config-file=$(rootpath {})".format(pyproject_label),
    ] + ["$(rootpaths {})".format(src) for src in srcs]
    al_py_checker(
        name = "mypy",
        args_bin = mypy_args,
        tool = mypy_label,
        args_test = [mypy_args[0], "--check"] + mypy_args[1:],
        data = srcs + [pyproject_label],
        disable_fix = True,
    )

    al_py_checker(
        name = "flake8",
        disable_fix = True,
        tool = flake8_label,
        args_test = [
            "--toml-config=$(rootpath {})".format(pyproject_label),
        ] + ["$(rootpaths {})".format(src) for src in srcs],
        data = srcs + [pyproject_label],
    )
