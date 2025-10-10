def _impl(ctx):
    """
    Reference: https://github.com/rlespinasse/docker-drawio-desktop-headless
    """
    drawio = ctx.toolchains["//bzl/rules/drawio:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        {{
            set -x
            Xvfb "${{DISPLAY}}" ${{XVFB_OPTIONS}} &
            sleep 0.2
            timeout '{timeout}' '{drawio}' --no-sandbox "${{@}}" || true
            set +x
        }} >stderr.txt 2>&1
        if [ ! -f '{out}' ]; then
            cat stderr.txt
            echo "Output {out} was not built, check the log"
            exit 1
        fi
    """.format(
        drawio = drawio.drawio.path,
        out = ctx.outputs.out.path,
        timeout = ctx.attr.cmd_timeout,
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    display = ":{}".format(abs(hash(str(ctx.label))))
    ctx.actions.run(
        executable = script,
        inputs = ctx.files.srcs + [drawio.drawio],
        outputs = [ctx.outputs.out],
        env = {
            "DISPLAY": display,
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
    toolchains = ["//bzl/rules/drawio:toolchain_type"],
    attrs = {
        "arguments": attr.string_list(
            mandatory = True,
            doc = "Arguments, location statements are expanded",
        ),
        "cmd_timeout": attr.string(default = "1m", doc = "Drawio command timeout"),
        "srcs": attr.label_list(
            default = [],
            allow_files = True,
            doc = "Sources",
        ),
        "out": attr.output(mandatory = True, doc = "Output"),
    },
)
