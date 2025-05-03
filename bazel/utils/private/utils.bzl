load("@rules_python//python:py_binary.bzl", "py_binary")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

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
        outs = ["patched.tar"],
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

def pkg_tar_combined(name = None, tars = [], out = "out.tar", strip_components = 2, **genrule_kwargs):
    """
    combine several tars into one
    """
    cmd = "set -eux"
    cmd += "\n".join(["""
        mkdir -p '{dir}'
        tar -xf $(location {label}) --strip-components '{strip_components}' -C '{dir}'
    """.format(strip_components = strip_components, **tar) for tar in tars])
    cmd += """
        tar -cf $(location {output}) {dirs}
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
        outs = ["src.tar"],
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
            wheel_paths.append("$(location {})".format(wheel))
            srcs = srcs + [wheel]
            wheels_env_script = 'export PYTHONPATH="{}"\n'.format(":".join(wheel_paths))
        wheel_paths.append(["$${PYTHONPATH:-}"])
    kwargs["cmd"] = wheels_env_script + cmd
    kwargs["srcs"] = srcs

    native.genrule(**kwargs)
