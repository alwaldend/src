def _template_files_impl(ctx):
    args = ctx.actions.args()
    for data in ctx.files.data:
        args.add("--data", data.path)
    for template in ctx.files.srcs:
        args.add("--template", template.path)
    for output in ctx.outputs.outs:
        args.add("--output", output.path)
    ctx.actions.run(
        executable = ctx.executable.templater,
        arguments = ["template", args],
        outputs = ctx.outputs.outs,
        inputs = ctx.files.srcs + ctx.files.data,
        progress_message = "Templating %{label} to %{output}",
    )

al_template_files = rule(
    implementation = _template_files_impl,
    doc = "Load data files, then template the template and write the output",
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
        "templater": attr.label(
            executable = True,
            default = "//tools/template_files",
            doc = "Templater to use",
            cfg = "exec",
        ),
    },
)
