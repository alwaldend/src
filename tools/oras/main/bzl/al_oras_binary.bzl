_SCRIPT = """\
#!/usr/bin/env sh
set -eu
for root in "" "${{RUNFILES_DIR:-}}/{workspace_name}/" "${{0}}.runfiles/{workspace_name}"; do
    if [ -x "${{root}}{bin}" ]; then
        exec "${{root}}{bin}" {arguments} "${{@}}"
    fi
done
echo "Could not find the binary"
exit 1
"""

def _impl(ctx):
    oras = ctx.toolchains["//tools/oras/main/bzl:toolchain_type"]
    script = ctx.actions.declare_file("{}/oras".format(ctx.label.name))
    runfiles = ctx.runfiles(files = [oras.oras])
    ctx.actions.write(
        output = script,
        content = _SCRIPT.format(
            workspace_name = ctx.workspace_name,
            bin = oras.oras.short_path,
            arguments = " ".join(
                [
                    '"{}"'.format(ctx.expand_make_variables(str(ctx.label), arg, {}))
                    for arg in ctx.attr.arguments
                ],
            ),
        ),
    )

    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
    ]

al_oras_binary = rule(
    implementation = _impl,
    doc = "Oras binary",
    executable = True,
    toolchains = ["//tools/oras/main/bzl:toolchain_type"],
    attrs = {
        "arguments": attr.string_list(
            doc = "Arguments",
        ),
    },
)
