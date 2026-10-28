"""Analyze captured Bazel dependencies as target diagrams or package graphs."""

def _bazel_deps_mermaid_from_json_impl(ctx):
    if ctx.attr.depth < -1:
        fail("depth must be -1 or nonnegative")
    output = ctx.actions.declare_file(ctx.attr.out or ctx.label.name + ".mmd")
    if not output.basename.endswith(".mmd"):
        fail("out must name a .mmd file")
    arguments = ctx.actions.args()
    arguments.add_all(["analyze", "target"])
    arguments.add("--input", ctx.file.src.path)
    arguments.add("--output", output.path)
    arguments.add("--root", ctx.attr.root)
    arguments.add("--depth", ctx.attr.depth)
    arguments.add("--diagram-type", ctx.attr.diagram_type)
    arguments.add("--direction", ctx.attr.direction)
    if ctx.attr.build_deps_only:
        arguments.add("--build-deps-only")
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
        "build_deps_only": attr.bool(default = False, doc = "Keep explicit deps/srcs label values intersected with ruleInput; requires complete capture metadata."),
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

def _bazel_deps_packages_from_json_impl(ctx):
    if ctx.attr.package_depth < 1:
        fail("package_depth must be a positive integer")
    extension = ".json" if ctx.attr.format == "json" else ".mmd"
    output = ctx.actions.declare_file(ctx.attr.out or ctx.label.name + extension)
    if not output.basename.endswith(extension):
        fail("out must name a {} file for format {}".format(extension, ctx.attr.format))
    arguments = ctx.actions.args()
    arguments.add_all(["analyze", "package"])
    arguments.add("--input", ctx.file.src.path)
    arguments.add("--output", output.path)
    arguments.add("--package-depth", ctx.attr.package_depth)
    arguments.add("--direction", ctx.attr.direction)
    arguments.add("--format", ctx.attr.format)
    if ctx.attr.build_deps_only:
        arguments.add("--build-deps-only")
    if not ctx.attr.include_files:
        arguments.add("--include-files=false")
    if not ctx.attr.include_external:
        arguments.add("--include-external=false")
    if not ctx.attr.merge_reciprocal:
        arguments.add("--merge-reciprocal=false")
    ctx.actions.run(
        executable = ctx.executable._converter,
        arguments = [arguments],
        inputs = [ctx.file.src],
        outputs = [output],
        mnemonic = "BazelPackageGraph",
        progress_message = "Generating package graph %{label}",
    )
    return [DefaultInfo(files = depset([output]))]

bazel_deps_packages_from_json = rule(
    implementation = _bazel_deps_packages_from_json_impl,
    doc = "Aggregates every captured Bazel target into package nodes and directed relationships.",
    attrs = {
        "src": attr.label(allow_single_file = True, mandatory = True, doc = "One workspace's Bazel streamed_jsonproto capture."),
        "build_deps_only": attr.bool(default = False, doc = "Keep explicit deps/srcs label values intersected with ruleInput; requires complete capture metadata."),
        "out": attr.string(doc = "Output path; defaults to <name>.mmd or <name>.json according to format."),
        "package_depth": attr.int(default = 1, doc = "Positive number of local package path segments to retain."),
        "direction": attr.string(default = "LR", values = ["LR", "RL", "TB", "BT"], doc = "Mermaid flowchart direction."),
        "format": attr.string(default = "mermaid", values = ["mermaid", "json"], doc = "Mermaid presentation or shared JSON aggregate."),
        "include_files": attr.bool(default = True, doc = "Include source and generated files."),
        "include_external": attr.bool(default = True, doc = "Include external repositories in one aggregate."),
        "merge_reciprocal": attr.bool(default = True, doc = "Combine opposite Mermaid arrows; JSON always retains both directions."),
        "_converter": attr.label(
            default = Label("//projects/bazel_deps_analyzer:bazel_deps_analyzer"),
            executable = True,
            cfg = "exec",
        ),
    },
)

def bazel_deps_mermaid(name, target, depth = -1, build_deps_only = False, **kwargs):
    """Capture non-tool build inputs and convert them to Mermaid.

    Args:
        name: Name of the Mermaid conversion target.
        target: Root target label, resolved relative to the caller's package.
        depth: Maximum dependency distance, or -1 for the complete closure.
        build_deps_only: Capture complete attribute metadata and select explicit deps/srcs inputs.
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
            "--noimplicit_deps" if build_deps_only else "--implicit_deps",
            "--nonodep_deps",
            "--notool_deps",
            "--proto:output_rule_attrs=name,deps,srcs" if build_deps_only else "--proto:output_rule_attrs=",
            "--proto:locations=false",
        ] + ([
            "--proto:default_values=true",
            "--proto:flatten_selects=false",
            "--proto:rule_inputs_and_outputs=true",
        ] if build_deps_only else []),
        visibility = ["//visibility:private"],
        tags = kwargs.get("tags", []),
        testonly = kwargs.get("testonly", False),
    )
    bazel_deps_mermaid_from_json(
        name = name,
        src = ":" + query_name,
        root = label,
        depth = depth,
        build_deps_only = build_deps_only,
        **kwargs
    )
