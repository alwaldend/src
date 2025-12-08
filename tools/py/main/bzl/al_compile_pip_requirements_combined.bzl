load("@rules_python//python:pip.bzl", "compile_pip_requirements")
load("//tools/txt/main/bzl:al_combine_files.bzl", "al_combine_files")

def al_compile_pip_requirements_combined(name, srcs, **kwargs):
    """
    Create compile_pip_requirements target for several requirement files

    Args:
        name: compile_pip_requirements name
        srcs: list of labels of requirement files to combine
        **kwargs: kwargs for `compile_pip_requirements`
    """
    al_combine_files(
        name = "{}.src".format(name),
        srcs = srcs,
    )
    compile_pip_requirements(
        name = name,
        src = ":{}.src".format(name),
        **kwargs
    )
