def al_resolved_toolchain(toolchain_label, **kwargs):
    """
    Create a resolved toolchain

    Args:
        toolchain: toolchain label
        **kwargs: rule kwargs

    Returns:
        Resolved toolchain rule
    """

    def impl(ctx):
        toolchain = ctx.toolchains[toolchain_label]
        if hasattr(toolchain, "default_info"):
            default_info = [toolchain.default_info]
        else:
            default_info = []
        return default_info + [
            toolchain,
            platform_common.TemplateVariableInfo(toolchain.env),
        ]

    return rule(
        implementation = impl,
        toolchains = [toolchain_label],
        **kwargs
    )
