def _impl(ctx):
    template_variables = platform_common.TemplateVariableInfo(
        {
            "DNSCONTROL_BIN": ctx.executable.dnscontrol.path,
        },
    )
    runfiles = ctx.runfiles().merge(ctx.attr.dnscontrol[DefaultInfo].default_runfiles)
    default_info = DefaultInfo(
        files = depset(ctx.files.dnscontrol),
        runfiles = runfiles,
    )
    toolchain_info = platform_common.ToolchainInfo(
        default_info = default_info,
        dnscontrol = ctx.executable.dnscontrol,
        template_variables = template_variables,
    )
    return [default_info, toolchain_info]

al_dnscontrol_toolchain = rule(
    implementation = _impl,
    doc = "Dnscontrol toolchain",
    attrs = {
        "dnscontrol": attr.label(
            doc = "Dnscontrol binary",
            allow_single_file = True,
            executable = True,
            mandatory = True,
            cfg = "exec",
        ),
    },
)
