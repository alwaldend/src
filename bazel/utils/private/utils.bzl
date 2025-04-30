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

def py_binary_shell(name = "", cmd = "", deps = [], **kwargs):
    """
    Create a shell with linked python deps
    """
    native.py_binary(
        name = name,
        srcs = ["//python/bazel-shell:shell"],
	main = "main.py",
	deps = deps,
    )

def genrule_with_wheels(wheels = None, srcs = [], cmd = [],  **kwargs):
    """
    Regular genrule with wheels added to ${PYTHONPATH}
    """
    wheels_env_script = ""
    if wheels:
        wheel_paths = []
        for wheel in wheels:
            wheel_paths.append("$(location {})".format(wheel))
            srcs  = srcs + [wheel]
            wheels_env_script = 'export PYTHONPATH="{}"\n'.format(":".join(wheel_paths))
        wheel_paths.append(["$${PYTHONPATH:-}"])
    kwargs["cmd"] = wheels_env_script + cmd
    kwargs["srcs"] = srcs

    native.genrule(**kwargs)
