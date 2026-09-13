"""Rules for producing malformed skill_library analysis fixtures."""

def _fixture_files_impl(ctx):
    for output in ctx.outputs.outs:
        ctx.actions.write(output = output, content = "fixture\n")
    return [DefaultInfo(files = depset(ctx.outputs.outs))]

fixture_files = rule(
    implementation = _fixture_files_impl,
    attrs = {
        "outs": attr.output_list(mandatory = True),
    },
)
