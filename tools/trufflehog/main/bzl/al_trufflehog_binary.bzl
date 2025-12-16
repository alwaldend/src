_SCRIPT = """\
#!/usr/bin/env sh
set -eu
for root in "" "${{RUNFILES_DIR:-}}/" "${{0}}.runfiles/{workspace_name}"; do
    if [ -x "${{root}}{bin}" ]; then
        exec "${{root}}{bin}" {arguments} "${{@}}"
    fi
done
echo "Could not find the binary"
exit 1
"""

def _impl(ctx):
    trufflehog = ctx.toolchains["//tools/trufflehog/main/bzl:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = [trufflehog.trufflehog])
    ctx.actions.write(
        output = script,
        content = _SCRIPT.format(
            workspace_name = ctx.workspace_name,
            bin = trufflehog.trufflehog.short_path,
            arguments = " ".join(['"{}"'.format(arg) for arg in ctx.attr.arguments]),
        ),
    )

    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
    ]

al_trufflehog_binary = rule(
    implementation = _impl,
    doc = "Trufflehog binary",
    executable = True,
    toolchains = ["//tools/trufflehog/main/bzl:toolchain_type"],
    attrs = {
        "arguments": attr.string_list(
            doc = "Arguments",
        ),
    },
)
