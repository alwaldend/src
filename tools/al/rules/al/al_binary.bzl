load("@rules_go//go:def.bzl", "go_binary")

def al_binary_run(
        name,
        args = [],
        run_args = [],
        configs = [],
        data = [],
        al = ["//tools/al/cmd/al:al_lib"],
        **kwargs):
    """
    Create a binary for `al run`
    """
    data = data + configs
    args_res = ["run"]
    for config in configs:
        args_res.extend(("--config", "$(rootpath {})".format(config)))
    args_res.extend(run_args)
    args_res.append("--")
    args_res.extend(args)
    go_binary(
        name = name,
        args = args_res,
        data = data,
        embed = al,
        **kwargs
    )

def al_binary_run_map(
        name,
        args = [],
        **kwargs):
    """
    Create several al_binary_run runners
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        al_binary_run(
            name = cur_name,
            args = args_value,
            **kwargs
        )
