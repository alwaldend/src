load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    out_dirs = []
    for out_dir in ctx.attr.out_dirs:
        out_dirs.append(ctx.actions.declare_directory(out_dir))
    args = ctx.actions.args()
    for arg in ctx.attr.arguments:
        args.add(ctx.expand_location(arg))
    ctx.actions.run(
        executable = ctx.executable.hugo,
        outputs = out_dirs + ctx.outputs.outs,
        arguments = [args],
        progress_message = "Running hugo action: %{label}",
    )
    return [
        DefaultInfo(
            files = depset(out_dirs + ctx.outputs.outs),
        ),
    ]

al_hugo_run_binary = rule(
    implementation = _impl,
    doc = "Run hugo binary as a build action",
    attrs = {
        "outs": attr.output_list(
            doc = "Output files",
        ),
        "out_dirs": attr.string_list(
            default = [],
            doc = "Output directories",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Hugo arguments",
        ),
        "hugo": attr.label(
            mandatory = True,
            executable = True,
            cfg = "exec",
            doc = "Hugo binary to use",
        ),
    },
)
