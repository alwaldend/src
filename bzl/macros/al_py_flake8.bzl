load(":al_py_checker.bzl", "al_py_checker")

def al_py_flake8(name = None, srcs = None, flake8_src = None, pyproject_src = None, run_args_src = None):
    """
    Generate -fix and -test targets for isort
    """
