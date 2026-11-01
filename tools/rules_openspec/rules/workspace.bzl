"""Declare owner-local OpenSpec workspace sources as validated targets."""

def _openspec_workspace_impl(ctx):
    return [
        DefaultInfo(
            files = depset(ctx.files.srcs),
            default_runfiles = ctx.runfiles(files = ctx.files.srcs),
        ),
    ]

openspec_workspace = rule(
    implementation = _openspec_workspace_impl,
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = True,
            doc = "OpenSpec configuration, specifications, and change artifacts.",
        ),
    },
    doc = "Declare an owner-local OpenSpec workspace source tree.",
)
