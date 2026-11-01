"""Executable targets that apply caller-declared user systemd scope properties."""

load("@aspect_bazel_lib//lib:copy_file.bzl", "COPY_FILE_TOOLCHAINS", "copy_file_action")

def _runfile_path(ctx, file):
    """Return the canonical runfiles key for a declared executable."""
    if file.short_path.startswith("../"):
        return file.short_path[3:]
    return ctx.workspace_name + "/" + file.short_path

def _systemd_run_impl(ctx):
    """Package an executable launcher, configuration, and application runfiles."""
    if not ctx.attr.properties:
        fail("properties must contain at least one systemd assignment")
    for name, value in ctx.attr.properties.items():
        if not name or "=" in name or not value:
            fail("properties require nonempty names and values")

    executable = ctx.actions.declare_file(ctx.label.name)
    configuration = ctx.actions.declare_file(ctx.label.name + ".json")
    copy_file_action(ctx, ctx.executable._runner, executable)
    targets = [ctx.attr.binary] + ctx.attr.data
    ctx.actions.write(configuration, json.encode({
        "schema": "systemd_run/v1",
        "binary": _runfile_path(ctx, ctx.executable.binary),
        "arguments": [ctx.expand_location(argument, targets = targets) for argument in ctx.attr.arguments],
        "properties": [name + "=" + ctx.attr.properties[name] for name in sorted(ctx.attr.properties)],
    }) + "\n")

    runfiles = ctx.runfiles(files = [executable, configuration] + ctx.files.data)
    runfiles = runfiles.merge_all([target[DefaultInfo].default_runfiles for target in targets])
    providers = [DefaultInfo(
        executable = executable,
        files = depset([executable, configuration]),
        runfiles = runfiles,
    )]
    if RunEnvironmentInfo in ctx.attr.binary:
        providers.append(ctx.attr.binary[RunEnvironmentInfo])
    return providers

systemd_run = rule(
    implementation = _systemd_run_impl,
    executable = True,
    doc = "Wrap a binary in a transient user systemd scope with BUILD-owned limits.",
    attrs = {
        "binary": attr.label(
            mandatory = True,
            executable = True,
            cfg = "target",
            doc = "Raw executable; configure its arguments on this wrapper.",
        ),
        "arguments": attr.string_list(
            doc = "Configured argument prefix, with location/runfiles expressions expanded from binary and data.",
        ),
        "data": attr.label_list(
            allow_files = True,
            doc = "Additional runfiles and labels available to argument expansion.",
        ),
        "properties": attr.string_dict(
            mandatory = True,
            doc = "Nonempty systemd NAME-to-VALUE assignments, such as MemoryMax = 36G.",
        ),
        "_runner": attr.label(
            default = Label("//cmd/systemd_run"),
            executable = True,
            cfg = "target",
        ),
    },
    toolchains = COPY_FILE_TOOLCHAINS,
)
