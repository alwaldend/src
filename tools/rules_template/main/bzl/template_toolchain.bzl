def _impl(ctx):
    templater = ctx.executable.templater
    template_variables = platform_common.TemplateVariableInfo(
        {
            "TEMPLATE_BIN": templater.path,
        },
    )
    default_info = DefaultInfo(
        files = depset([templater]),
        runfiles = ctx.attr.templater.default_runfiles.merge(ctx.runfiles(files = [templater])),
    )
    toolchain_info = platform_common.ToolchainInfo(
        default_info = default_info,
        templater = templater,
        template_variables = template_variables,
    )
    return [default_info, toolchain_info]

template_toolchain = rule(
    implementation = _impl,
    doc = "Template toolchain",
    attrs = {
        "templater": attr.label(
            doc = "Templater binary",
            executable = True,
            mandatory = True,
            cfg = "exec",
        ),
    },
)
