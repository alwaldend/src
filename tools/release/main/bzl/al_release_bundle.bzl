"""Package declared payloads and versioning metadata for the release publisher."""

def _impl(ctx):
    bundle = ctx.actions.declare_directory(ctx.label.name)
    args = ctx.actions.args()
    args.add("generate")
    args.add("--project", ctx.attr.project)
    args.add("--version_file", ctx.info_file)
    args.add("--output_dir", bundle.path)
    args.add_all(ctx.files.srcs, before_each = "--add_file")
    ctx.actions.run(
        executable = ctx.executable.release_tool,
        arguments = [args],
        inputs = ctx.files.srcs + [ctx.info_file],
        outputs = [bundle],
        mnemonic = "ReleaseBundle",
        progress_message = "Packaging release bundle %{label}",
    )
    return [DefaultInfo(files = depset([bundle]))]

al_release_bundle = rule(
    implementation = _impl,
    doc = "Emit release.json and files/ using STABLE_VERSION from the versioning bootstrap.",
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = True,
            doc = "Payload files; basenames must be unique.",
        ),
        "project": attr.string(mandatory = True, doc = "Canonical project subdirectory."),
        "release_tool": attr.label(
            default = "//tools/release/main/go",
            executable = True,
            cfg = "exec",
            doc = "Manifest generator and bundle writer.",
        ),
    },
)
