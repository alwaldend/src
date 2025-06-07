def al_genrule_with_wheels(name, wheels, srcs = [], cmd = [], **kwargs):
    """
    Regular genrule with wheels added to ${PYTHONPATH}

    Args:
        name: genrule name
        wheels: list of wheel labels
        srcs: srcs for the genrule
        cmd: genrule cmd
        **kwargs: other genrule kwargs
    """
    wheels_env_script = ""
    if wheels:
        wheel_paths = []
        for wheel in wheels:
            wheel_paths.append("$(rootpath {})".format(wheel))
            srcs = srcs + [wheel]
            wheels_env_script = 'export PYTHONPATH="{}"\n'.format(":".join(wheel_paths))
        wheel_paths.append(["$${PYTHONPATH:-}"])
    kwargs["name"] = name
    kwargs["cmd"] = wheels_env_script + cmd
    kwargs["srcs"] = srcs

    native.genrule(**kwargs)
