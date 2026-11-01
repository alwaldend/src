"""Package declared payloads and versioning metadata for the release publisher."""

load("@bazel_lib//lib:stamping.bzl", "STAMP_ATTRS", "maybe_stamp")
load(":release_files_info.bzl", "ReleaseFilesInfo")

def _impl(ctx):
    bundle = ctx.actions.declare_directory(ctx.label.name)
    args = ctx.actions.args()
    args.add("release")
    args.add("generate")
    args.add("--project", ctx.attr.project)
    status = ctx.file.status_file
    if not status:
        stamp = maybe_stamp(ctx)
        if stamp:
            status = stamp.stable_status_file
        else:
            status = ctx.actions.declare_file(ctx.label.name + ".development-status.txt")
            ctx.actions.write(status, "STABLE_VERSION 0.0.0-dev\nSTABLE_VERSION_CHANNEL development\n")
    args.add("--version_file", status)
    args.add("--output_dir", bundle.path)
    payloads = {}
    for file in ctx.files.srcs:
        if file.basename in payloads:
            fail("duplicate bundle payload: " + file.basename)
        payloads[file.basename] = file
    for dep in ctx.attr.deps:
        for name, file in dep[ReleaseFilesInfo].files.items():
            if name in payloads:
                fail("duplicate bundle payload: " + name)
            payloads[name] = file
    args.add_all(payloads.values(), before_each = "--add_file")
    ctx.actions.run(
        executable = ctx.executable.release_tool,
        arguments = [args],
        inputs = payloads.values() + [status],
        outputs = [bundle],
        mnemonic = "ReleaseBundle",
        progress_message = "Packaging release bundle %{label}",
    )
    return [DefaultInfo(files = depset([bundle]))]

release_bundle = rule(
    implementation = _impl,
    doc = "Emit release.json and files/ using STABLE_VERSION from the versioning bootstrap.",
    attrs = {
        "status_file": attr.label(allow_single_file = True, doc = "Explicit build metadata; defaults to Bazel stable workspace status."),
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Payload files; basenames must be unique.",
        ),
        "deps": attr.label_list(providers = [ReleaseFilesInfo], doc = "Existing project release payload sets."),
        "project": attr.string(mandatory = True, doc = "Canonical project subdirectory."),
        "release_tool": attr.label(
            default = Label("//cmd/al:al"),
            executable = True,
            cfg = "exec",
            doc = "Manifest generator and bundle writer.",
        ),
    } | STAMP_ATTRS,
)
