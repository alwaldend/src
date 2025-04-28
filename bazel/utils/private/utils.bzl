def py_binary_shell(name = "", cmd = "", deps = [], **kwargs):
    native.py_binary(
        name = name,
        srcs = ["//python/bazel-shell:shell"],
	main = "main.py",
	deps = deps,
    )

def custom_genrule(wheels = None, srcs = [], cmd = [],  **kwargs):
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
