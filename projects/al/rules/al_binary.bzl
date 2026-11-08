load("@bazel_skylib//rules:native_binary.bzl", "native_binary", "native_test")
load(":al_config.bzl", "al_config")

def _patch_kwargs(kwargs, name, testonly = False):
    """
    Patch kwargs for run binaries
    """
    data = [] + kwargs.pop("data", [])
    args_res = ["run"]
    configs = kwargs.pop("configs", [])
    if len(configs) > 1:
        combined = name + ".config"
        common = {key: kwargs[key] for key in ["tags", "target_compatible_with", "exec_compatible_with", "features"] if key in kwargs}
        al_config(
            name = combined,
            configs = configs,
            testonly = testonly or kwargs.get("testonly", False),
            visibility = ["//visibility:private"],
            **common
        )
        configs = [":" + combined]
    for config in configs:
        data.append(config)
        args_res.extend(("--config", "$(rootpath {})".format(config)))
    for label, value in kwargs.pop("plugin_labels", {}).items():
        if type(label) != "string" or type(value) != "string":
            fail("plugin_labels must map label names to string values")
        selector = label + "=" + value if value else label
        args_res.extend(("--plugin_label", selector))
    args_res.extend(kwargs.pop("run_args", []))
    args_res.append("--")
    args_res.extend(kwargs.pop("args", []))
    kwargs["args"] = args_res
    kwargs["data"] = data
    return kwargs

def al_binary_run(src = Label("//cmd/al:al"), **kwargs):
    """
    Create an AL command with configs and dictionary plugin_labels selectors.

    Empty selector values match the label's presence. run_args carries other
    AL command options; args carries the child executable and its arguments.
    """
    kwargs = _patch_kwargs(kwargs, kwargs["name"])
    native_binary(
        src = src,
        **kwargs
    )

def al_binary_run_test(name, **kwargs):
    """
    Create a native test for al_binary_run
    """
    kwargs = _patch_kwargs(kwargs, name, testonly = True)
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
