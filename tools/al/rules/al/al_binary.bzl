load("@bazel_skylib//rules:native_binary.bzl", "native_test")
load("@rules_go//go:def.bzl", "go_binary")

def _patch_kwargs(kwargs):
    """
    Patch kwargs for run binaries
    """
    data = [] + kwargs.pop("data", [])
    args_res = ["run"]
    for config in kwargs.pop("configs", []):
        data.append(config)
        args_res.extend(("--config", "$(rootpath {})".format(config)))
    args_res.extend(kwargs.pop("run_args", []))
    args_res.append("--")
    args_res.extend(kwargs.pop("args", []))
    kwargs["args"] = args_res
    kwargs["data"] = data
    return kwargs

def al_binary_run(
        name,
        al = ["//tools/al/cmd/al:al_lib"],
        **kwargs):
    """
    Create a binary for `al run`
    """
    kwargs = _patch_kwargs(kwargs)
    go_binary(
        name = name,
        embed = al,
        **kwargs
    )

def al_binary_run_test(name, **kwargs):
    """
    Create a native test for al_binary_run
    """
    kwargs = _patch_kwargs(kwargs)
    bin_kwargs = {} | kwargs
    bin_kwargs.pop("size", None)
    al_binary_run(
        name = "{}.bin".format(name),
        **bin_kwargs
    )
    native_test(
        name = name,
        src = "{}.bin".format(name),
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
