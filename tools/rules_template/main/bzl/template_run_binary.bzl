def _template_files_impl(ctx):
    toolchain = ctx.toolchains["//main/bzl:toolchain_type"]
    args = ctx.actions.args()
    for data in ctx.files.data:
        args.add("--data", data.path)
    for template in ctx.files.srcs:
        args.add("--template", template.path)
    for output in ctx.outputs.outs:
        args.add("--output", output.path)
    args.add_all(ctx.attr.args)
    ctx.actions.run(
        executable = toolchain.templater,
        arguments = ["template", args],
        outputs = ctx.outputs.outs,
        inputs = ctx.files.srcs + ctx.files.data,
        progress_message = "Templating %{label} to %{output}",
    )

template_run_binary = rule(
    implementation = _template_files_impl,
    doc = "Rule to template files",
    toolchains = ["//main/bzl:toolchain_type"],
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            mandatory = True,
            allow_empty = False,
            doc = "Template files",
        ),
        "outs": attr.output_list(
            mandatory = True,
            allow_empty = False,
            doc = "Output files",
        ),
        "data": attr.label_list(
            allow_files = True,
            mandatory = True,
            allow_empty = True,
            doc = "Data files",
        ),
        "args": attr.string_list(
            doc = "Extra arguments",
            default = [],
        ),
    },
)
