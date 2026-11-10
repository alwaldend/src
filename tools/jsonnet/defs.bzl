"""Generate declared JSON outputs with the pinned Jsonnet compiler."""

def _jsonnet_to_json_impl(ctx):
    if not ctx.outputs.outs:
        fail("At least one named JSON output is required")
    args = ctx.actions.args()
    args.add("--create-output-dirs")
    args.add("--multi", ctx.bin_dir.path + "/" + ctx.label.package)
    args.add(ctx.file.src.path)
    ctx.actions.run(
        executable = ctx.executable._compiler,
        arguments = [args],
        inputs = depset([ctx.file.src] + ctx.files.deps),
        outputs = ctx.outputs.outs,
        tools = [ctx.attr._compiler[DefaultInfo].files_to_run],
        mnemonic = "JsonnetToJson",
        progress_message = "Generating JSON from %{label}",
    )
    return [DefaultInfo(files = depset(ctx.outputs.outs))]

jsonnet_to_json = rule(
    implementation = _jsonnet_to_json_impl,
    doc = "Compile a Jsonnet object to explicit, package-relative JSON outputs.",
    attrs = {
        "src": attr.label(allow_single_file = [".jsonnet"], mandatory = True),
        "deps": attr.label_list(allow_files = True, doc = "All imported Jsonnet and data files."),
        "outs": attr.output_list(mandatory = True),
        "_compiler": attr.label(
            default = "//third_party/com_github_google_go_jsonnet:jsonnet",
            executable = True,
            cfg = "exec",
        ),
    },
)
