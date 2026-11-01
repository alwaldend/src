def _impl(ctx):
    binary = ctx.executable.binary
    template_variables = platform_common.TemplateVariableInfo(
        {
            "BINARY_TOOLCHAIN_{}".format(binary.basename): binary.path,
        },
    )
    default_info = DefaultInfo(
        files = depset([binary]),
        runfiles = ctx.attr.binary.default_runfiles.merge(ctx.runfiles(files = [binary])),
    )
    toolchain_info = platform_common.ToolchainInfo(
        default_info = default_info,
        binary = binary,
        template_variables = template_variables,
    )
    return [default_info, toolchain_info, template_variables]

binary_toolchain_toolchain = rule(
    implementation = _impl,
    doc = "Binary toolchain",
    attrs = {
        "binary": attr.label(
            doc = "Binary",
            executable = True,
            mandatory = True,
            cfg = "exec",
        ),
    },
)
