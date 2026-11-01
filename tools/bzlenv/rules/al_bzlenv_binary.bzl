load("@bazel_skylib//lib:paths.bzl", "paths")

_SCRIPT = """\
#!/usr/bin/env sh
set -eu
mkdir -p bin
export BZLENV_RUNFILES_DIR="$(dirname ${{PWD}})"
export BZLENV_PATH="{path}"
ln -sf "${{PWD}}" "${{BUILD_WORKSPACE_DIRECTORY}}/.bzlenv"
envsubst '${{BZLENV_RUNFILES_DIR}} ${{BZLENV_PATH}} ${{BUILD_WORKSPACE_DIRECTORY}}' <'{activate}' >bin/activate
echo "${{PWD}}/bin/activate"
"""

def _impl(ctx):
    exec = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = [ctx.file.activate])
    tool_directory = paths.join(ctx.label.package, "{}.bin".format(ctx.label.name))
    tool_links = {}
    for tool_name, tool in ctx.attr.tools.items():
        runfiles = runfiles.merge(tool[DefaultInfo].default_runfiles)
        executable = tool[DefaultInfo].files_to_run.executable
        if executable == None:
            fail("Tool {} ({}) must be executable".format(tool_name, tool.label))
        tool_links[paths.join(tool_directory, tool_name)] = executable
    runfiles = runfiles.merge(ctx.runfiles(symlinks = tool_links))
    ctx.actions.write(
        output = exec,
        is_executable = True,
        content = _SCRIPT.format(
            activate = ctx.file.activate.short_path,
            path = "${{0}}.runfiles/{}/{}".format(ctx.workspace_name, tool_directory),
        ),
    )

    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = exec,
        ),
    ]

al_bzlenv_binary = rule(
    doc = "Dev shell binary",
    implementation = _impl,
    executable = True,
    attrs = {
        "tools": attr.string_keyed_label_dict(
            cfg = "exec",
            doc = "Tools, keys are tool names, values are tool labels",
        ),
        "activate": attr.label(
            default = "//tools/bzlenv/internal/scripts:activate",
            allow_single_file = True,
            doc = "Activation script",
        ),
    },
)
