load("@bazel_skylib//lib:shell.bzl", "shell")

_SCRIPT = """\
#!/usr/bin/env sh
set -eu
exec {arguments} "${{@}}"
"""

def _impl(ctx):
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    env = {}
    if len(ctx.files.srcs) == 0:
        fail("missing iso files")
    elif len(ctx.files.srcs) == 1:
        env["ISO_PATH"] = "${{0}}/{}".format(ctx.files.srcs[0].short_path)
    else:
        env["ISO_PATHS"] = " ".join(ctx.files.srcs)
    script_content = _SCRIPT.format(
        arguments = " ".join(
            [
                shell.quote(
                    ctx.expand_make_variables(
                        str(ctx.label),
                        ctx.expand_location(argument),
                        env,
                    ),
                )
                for argument in ctx.attr.arguments
            ],
        ),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    return [
        DefaultInfo(
            executable = script,
            files = depset(ctx.files.srcs),
        ),
    ]

al_iso_binary = rule(
    doc = "Run a script using ISO",
    executable = True,
    implementation = _impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "Iso files",
            allow_files = ["iso"],
            default = [],
        ),
        "arguments": attr.string_list(
            doc = "Arguments (make variables and locations are expanded)",
            default = [],
        ),
    },
)
