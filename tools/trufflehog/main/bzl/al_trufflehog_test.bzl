_SCRIPT = """\
#!/usr/bin/env sh
set -eu
exec "{trufflehog}" {arguments} "${{@}}"
"""

def _impl(ctx):
    trufflehog = ctx.toolchains["//tools/trufflehog/main/bzl:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = [trufflehog.trufflehog])
    ctx.actions.write(
        output = script,
        content = _SCRIPT.format(
            trufflehog = trufflehog.trufflehog.short_path,
            arguments = " ".join([
                '"{}"'.format(ctx.expand_make_variables("trufflehog_test", arg, {}))
                for arg in ctx.attr.arguments
            ]),
        ),
    )

    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
    ]

al_trufflehog_test = rule(
    implementation = _impl,
    doc = "Trufflehog binary",
    test = True,
    toolchains = ["//tools/trufflehog/main/bzl:toolchain_type"],
    attrs = {
        "arguments": attr.string_list(
            doc = "Arguments",
        ),
    },
)
