"""Declare owner-local OpenSpec workspace sources as validated targets."""

def _openspec_workspace_impl(ctx):
    files = ctx.files.config + ctx.files.srcs
    return [
        DefaultInfo(
            files = depset(files),
            default_runfiles = ctx.runfiles(files = files),
        ),
    ]

openspec_workspace = rule(
    implementation = _openspec_workspace_impl,
    attrs = {
        "config": attr.label(
            allow_single_file = [".yaml"],
            doc = "The OpenSpec configuration file. Existing callers may include it in srcs.",
        ),
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = True,
            doc = "OpenSpec configuration, specifications, and change artifacts.",
        ),
    },
    doc = "Declare an owner-local OpenSpec workspace source tree.",
)
