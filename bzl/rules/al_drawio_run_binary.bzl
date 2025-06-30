def _impl(ctx):
    drawio = ctx.toolchains["//bzl/toolchain-types:drawio"]
    script = ctx.actions.declare_file("{}-script.sh".format(ctx.label.name))
    script_content = """\
        #!/usr/bin/env sh
        set -eux
        Xvfb "${{DISPLAY}}" ${{XVFB_OPTIONS}} &
        '{drawio}' "${{@}}" || true
    """.format(
        drawio = drawio.drawio.path,
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    ctx.actions.run(
        executable = script,
        inputs = ctx.files.srcs + [drawio.drawio],
        outputs = [ctx.outputs.out],
        env = {
            "DISPLAY": ":42",
            "ELECTRON_DISABLE_SECURITY_WARNINGS": "true",
            "XVFB_OPTIONS": "-nolisten unix",
            "ELECTRON_ENABLE_LOGGING": "false",
            "APPIMAGE_EXTRACT_AND_RUN": "1",  # https://github.com/AppImage/AppImageKit/wiki/FUSE#docker
        },
        arguments = [ctx.expand_location(arg) for arg in ctx.attr.arguments],
    )

al_drawio_run_binary = rule(
    implementation = _impl,
    doc = "Run drawio a a build action",
    toolchains = ["//bzl/toolchain-types:drawio"],
    attrs = {
        "arguments": attr.string_list(
            mandatory = True,
            doc = "Arguments, location statements are expanded",
        ),
        "srcs": attr.label_list(
            default = [],
            allow_files = True,
            doc = "Sources",
        ),
        "out": attr.output(mandatory = True, doc = "Output"),
    },
)
