"""Workspace environment adapter for executable formatters."""

def _tool_adapter_impl(ctx):
    tool = ctx.executable.tool
    launcher = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.symlink(output = launcher, target_file = ctx.file._script, is_executable = True)
    runfiles = ctx.runfiles(root_symlinks = {"repo_quality_tools/" + ctx.label.name: tool})
    return [DefaultInfo(
        executable = launcher,
        runfiles = runfiles.merge(ctx.attr.tool[DefaultInfo].default_runfiles),
    )]

tool_adapter = rule(
    implementation = _tool_adapter_impl,
    executable = True,
    attrs = {
        "tool": attr.label(cfg = "target", executable = True, mandatory = True),
        "_script": attr.label(allow_single_file = True, default = Label("//tools/repo_quality:tool_adapter.sh")),
    },
)
