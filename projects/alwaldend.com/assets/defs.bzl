"""Compile site-owned Sass for public pages and offline application consumers."""

def _site_stylesheet_impl(ctx):
    arguments = ctx.actions.args()
    arguments.add_all(["--no-source-map", "--style=compressed"])
    inputs = [ctx.file.src] + ctx.files.deps
    for target, directory in ctx.attr.includes.items():
        root = "/".join([part for part in [target.label.workspace_root, target.label.package, directory] if part])
        arguments.add("--load-path=" + root)
        inputs.extend(target[DefaultInfo].files.to_list())
    arguments.add(ctx.file.src.path)
    arguments.add(ctx.outputs.out.path)
    ctx.actions.run(
        executable = ctx.executable._sass,
        arguments = [arguments],
        inputs = inputs,
        outputs = [ctx.outputs.out],
        tools = [ctx.attr._sass[DefaultInfo].files_to_run],
        mnemonic = "SiteStylesheet",
    )
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

site_stylesheet = rule(
    implementation = _site_stylesheet_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".scss"], mandatory = True),
        "deps": attr.label_list(allow_files = [".scss"]),
        "includes": attr.label_keyed_string_dict(allow_files = True),
        "out": attr.output(mandatory = True),
        "_sass": attr.label(
            default = "@com_github_sass_dart_sass//:sass_binary",
            executable = True,
            cfg = "exec",
        ),
    },
)
