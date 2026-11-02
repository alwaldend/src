def _impl(ctx):
    info = ctx.toolchains[Label("//rules:al_toolchain_type")]
    return [info.default_info, info, platform_common.TemplateVariableInfo(info.env)]

al_toolchain_resolved = rule(implementation = _impl, toolchains = [Label("//rules:al_toolchain_type")], doc = "Expose the AL execution toolchain.")
