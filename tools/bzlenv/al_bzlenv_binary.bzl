_SCRIPT = """\
#!/usr/bin/env sh
set -eu
export BZLENV_RUNFILES_DIR="$(dirname ${{PWD}})"
export BZLENV_BIN_DIR="${{PWD}}/bin"
ln -sf "${{PWD}}" "${{BUILD_WORKSPACE_DIRECTORY}}/.bzlenv"
envsubst '${{BZLENV_RUNFILES_DIR}} ${{BZLENV_BIN_DIR}} ${{BUILD_WORKSPACE_DIRECTORY}}' <'{activate}' >bin/activate
echo "${{PWD}}/bin/activate"
"""

def _impl(ctx):
    exec = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ctx.actions.write(
        output = exec,
        is_executable = True,
        content = _SCRIPT.format(activate = ctx.file.activate.short_path),
    )
    runfiles = ctx.runfiles(files = [ctx.file.activate])
    runfiles_symlinks = {}
    for tool_name, tool in ctx.attr.tools.items():
        tool_exec = tool[DefaultInfo].files_to_run.executable
        if not tool_exec:
            fail("Tool {} is missing an executable: {}".format(tool_name, tool))
        runfiles_symlinks["bin/{}".format(tool_name)] = tool_exec
        runfiles = runfiles.merge(tool[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge(ctx.runfiles(symlinks = runfiles_symlinks))

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
            default = "//tools/bzlenv:bzlenv_activate",
            allow_single_file = True,
            doc = "Activation script",
        ),
    },
)
