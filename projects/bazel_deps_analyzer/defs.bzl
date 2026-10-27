"""Convert Bazel target dependencies or captured JSON to Mermaid artifacts."""

def _bazel_deps_mermaid_from_json_impl(ctx):
    if ctx.attr.depth < -1:
        fail("depth must be -1 or nonnegative")
    output = ctx.actions.declare_file(ctx.attr.out or ctx.label.name + ".mmd")
    if not output.basename.endswith(".mmd"):
        fail("out must name a .mmd file")
    arguments = ctx.actions.args()
    arguments.add("--input", ctx.file.src.path)
    arguments.add("--output", output.path)
    arguments.add("--root", ctx.attr.root)
    arguments.add("--depth", ctx.attr.depth)
    arguments.add("--diagram-type", ctx.attr.diagram_type)
    arguments.add("--direction", ctx.attr.direction)
    if ctx.attr.include_files:
        arguments.add("--include-files")
    if ctx.attr.include_external:
        arguments.add("--include-external")
    if ctx.attr.external_leaves:
        arguments.add("--external-leaves")
    if ctx.attr.group_by_package:
        arguments.add("--group-by-package")
    ctx.actions.run(
        executable = ctx.executable._converter,
        arguments = [arguments],
        inputs = [ctx.file.src],
        outputs = [output],
        mnemonic = "BazelDepsMermaid",
        progress_message = "Generating dependency diagram %{label}",
    )
    return [DefaultInfo(files = depset([output]))]

bazel_deps_mermaid_from_json = rule(
    implementation = _bazel_deps_mermaid_from_json_impl,
    doc = "Converts captured Bazel streamed_jsonproto target records to Mermaid.",
    attrs = {
        "src": attr.label(allow_single_file = True, mandatory = True, doc = "Bazel streamed_jsonproto input."),
        "root": attr.string(mandatory = True, doc = "Complete root label as recorded in the JSON."),
        "out": attr.string(doc = "Output .mmd path; defaults to <name>.mmd."),
        "depth": attr.int(default = -1, doc = "Maximum dependency distance, or -1 for unlimited."),
        "diagram_type": attr.string(default = "flowchart", values = ["flowchart", "mindmap"], doc = "Dependency flowchart or ownership mindmap inventory."),
        "direction": attr.string(default = "LR", values = ["LR", "RL", "TB", "BT"], doc = "Mermaid flowchart direction; ignored for mindmap."),
        "include_files": attr.bool(default = False, doc = "Show source and generated file targets."),
        "include_external": attr.bool(default = False, doc = "Show external repository targets."),
        "external_leaves": attr.bool(default = False, doc = "Show first external targets as leaves without their dependencies."),
        "group_by_package": attr.bool(default = False, doc = "Group short target names by full package, nesting external packages by repository."),
        "_converter": attr.label(
            default = Label("//projects/bazel_deps_analyzer:bazel_deps_analyzer"),
            executable = True,
            cfg = "exec",
        ),
    },
)

def bazel_deps_mermaid(name, target, depth = -1, **kwargs):
    """Capture non-tool build inputs and convert them to Mermaid.

    Args:
        name: Name of the Mermaid conversion target.
        target: Root target label, resolved relative to the caller's package.
        depth: Maximum dependency distance, or -1 for the complete closure.
        **kwargs: Attributes forwarded to bazel_deps_mermaid_from_json.
    """
    if depth < -1:
        fail("depth must be -1 or nonnegative")
    root = native.package_relative_label(target)
    query_name = name + ".query.ndjson"
    label = str(root)
    if "'" in label or '"' in label or "\\" in label:
        fail("target label contains unsupported query quoting characters")
    native.genquery(
        name = query_name,
        expression = "deps('{}'{})".format(label, ", {}".format(depth) if depth >= 0 else ""),
        scope = [root],
        opts = [
            "--output=streamed_jsonproto",
            "--consistent_labels",
            "--implicit_deps",
            "--nonodep_deps",
            "--notool_deps",
            "--proto:output_rule_attrs=",
            "--proto:locations=false",
        ],
        visibility = ["//visibility:private"],
        tags = kwargs.get("tags", []),
        testonly = kwargs.get("testonly", False),
    )
    bazel_deps_mermaid_from_json(
        name = name,
        src = ":" + query_name,
        root = label,
        depth = depth,
        **kwargs
    )
