"""Prepare the site's plain Markdown for Hugo without changing source files."""

def _prepare_content_impl(ctx):
    output = ctx.actions.declare_file(ctx.label.name + ".tar")
    args = ctx.actions.args()
    args.add("--input", ctx.file.src)
    args.add("--output", output)
    ctx.actions.run(
        executable = ctx.executable._tool,
        arguments = [args],
        inputs = [ctx.file.src],
        outputs = [output],
        mnemonic = "PrepareSiteContent",
    )
    return [DefaultInfo(files = depset([output]))]

prepare_content = rule(
    implementation = _prepare_content_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".tar"], mandatory = True),
        "_tool": attr.label(
            default = "//projects/alwaldend.com/cmd/prepare_content",
            executable = True,
            cfg = "exec",
        ),
    },
)
