def _impl(ctx):
    out_name = ctx.attr.out
    if not out_name:
        out_name = "{}.out".format(ctx.label.name)
    out = ctx.actions.declare_directory(out_name)
    args = ctx.actions.args()
    args.add_all(["-xf", ctx.file.src.path, "-C", out.path])
    args.add_all(ctx.attr.arguments)
    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [out],
        arguments = [args],
        progress_message = "Extracting %{label}",
        command = 'tar "${@}"',
    )
    return [DefaultInfo(files = depset([out]))]

al_pkg_extract_dir = rule(
    implementation = _impl,
    doc = "Extract an archive into a TreeArtifact",
    attrs = {
        "src": attr.label(
            mandatory = True,
            allow_single_file = [".tar"],
            doc = "Archive to unpack",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Additional arguments",
        ),
        "out": attr.string(
            doc = "Output directory name",
        ),
    },
)
