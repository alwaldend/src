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
        return [
            toolchain,
            platform_common.TemplateVariableInfo(toolchain.env),
        ]

    return rule(
        implementation = impl,
        toolchains = [toolchain_label],
        **kwargs
    )
