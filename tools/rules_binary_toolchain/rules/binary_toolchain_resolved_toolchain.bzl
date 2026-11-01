def binary_toolchain_resolved_toolchain(toolchain_type, kwargs = {}):
    """
    Create a resolved binary_toolchain that can be used in genrules

    Args:
        toolchain_type: toolchain type label

    Returns:
        Toolchain rule
    """

    def _impl(ctx):
        res = []
        toolchain = ctx.toolchains[toolchain_type]
        res.append(toolchain)
        if hasattr(toolchain, "default_info"):
            res.append(toolchain.default_info)
        if hasattr(toolchain, "template_variables"):
            res.append(toolchain.template_variables)
        return res

    default_kwargs = {
        "implementation": _impl,
        "toolchains": [toolchain_type],
    }
    kwargs = default_kwargs | kwargs
    return rule(**kwargs)
