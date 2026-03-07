def _impl(ctx):
    toolchain = ctx.toolchains[ctx.attr.toolchain_type]
    outputs = [] + ctx.outputs.outs
    for out_dir in ctx.attr.out_dirs:
        dir = ctx.actions.declare_directory(out_dir)
        outputs.append(dir)
    ctx.actions.run(
        executable = toolchain.binary,
        inputs = ctx.files.srcs + ctx.files.data,
        outputs = outputs,
        arguments = [
            ctx.expand_make_variables(str(ctx.label), ctx.expand_location(arg), {})
            for arg in ctx.attr.arguments
        ],
    )
    runfiles = ctx.runfiles().merge_all(
        [
            ctx.runfiles(files = ctx.files.data)
        ] + [
            data[DefaultInfo].default_runfiles for data in ctx.attr.data
        ]
    )
    return [
        DefaultInfo(
            files = depset(ctx.outputs.outs),
            runfiles = runfiles,
        )
    ]

def binary_toolchain_run_binary(toolchain_type, kwargs = {}, attrs = {}):
    """
    Generate a run_binary rule for a binary_toolchain

    Args:
        toolchain_type: toolchain type label
        kwargs: override for rule kwargs
        attrs: override for attrs for the rule

    Returns:
        Rule
    """
    default_kwargs = {
        "implementation": _impl,
        "doc": "Run_binary rule for {}".format(toolchain_type),
        "toolchains": [toolchain_type],
    }
    default_attrs = {
        "srcs": attr.label_list(
            doc = "Inputs",
            allow_files = True,
        ),
        "data": attr.label_list(
            doc = "Inputs (also added to runfiles)",
            allow_files = True,
        ),
        "outs": attr.output_list(
            doc = "Outputs",
        ),
        "out_dirs": attr.string_list(
            doc = "Directory outputs",
        ),
        "toolchain_type": attr.string(
            doc = "Toolchain type",
            default = toolchain_type,
        ),
        "arguments": attr.string_list(
            doc = "Arguments (make variables and locations are expanded)",
        ),
    }
    kwargs = default_kwargs | kwargs
    kwargs["attrs"] = default_attrs | attrs
    res = rule(**kwargs)
    return res
