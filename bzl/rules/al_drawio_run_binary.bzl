def _impl(ctx):
    """
    Reference: https://github.com/rlespinasse/docker-drawio-desktop-headless
    """
    drawio = ctx.toolchains["//bzl/toolchain-types:drawio"]
    script = ctx.actions.declare_file("{}-script.sh".format(ctx.label.name))
    script_content = """\
        #!/usr/bin/env sh
        set -eux
        Xvfb "${{DISPLAY}}" ${{XVFB_OPTIONS}} &
        sleep 0.5
        '{drawio}' --no-sandbox "${{@}}" || true
        if [ ! -f '{out}' ]; then
            echo "Output {out} was not built, check the log"
            exit 1
        fi
    """.format(
        drawio = drawio.drawio.path,
        out = ctx.outputs.out.path,
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
            "DISPLAY": ":{}".format(hash(str(ctx.label))).replace("-", ""),
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
