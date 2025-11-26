def _helm_toolchain_impl(ctx):
    helm = ctx.file.helm
    template_variables = platform_common.TemplateVariableInfo(
        {
            "HELM_BIN": helm.path,
        },
    )
    default_info = DefaultInfo(
        files = depset([helm]),
        runfiles = ctx.runfiles(files = [helm]),
    )
    toolchain_info = platform_common.ToolchainInfo(
        default_info = default_info,
        helm = helm,
        template_variables = template_variables,
    )
    return [default_info, toolchain_info]

al_helm_toolchain = rule(
    implementation = _helm_toolchain_impl,
    doc = "Helm toolchain",
    attrs = {
        "helm": attr.label(
            doc = "Helm binary",
            allow_single_file = True,
            mandatory = True,
            cfg = "exec",
        ),
    },
)
