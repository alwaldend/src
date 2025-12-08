load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    dnscontrol = ctx.toolchains["//tools/dnscontrol/main/bzl:toolchain_type"]
    runfiles_files = [ctx.file.creds]
    runfiles_symlinks = {
        "dnsconfig.js": ctx.file.config,
    } | {
        "zones/{}".format(file.basename): file
        for file in ctx.files.zones
    }
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(
        files = runfiles_files,
        symlinks = runfiles_symlinks,
    ).merge(dnscontrol.default_info.default_runfiles)

    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        set -eu
        creds="$(mktemp)"
        envsubst <'{creds}' >"${{creds}}"
        exec '{dnscontrol}' \
            {arguments} \
            --creds "${{creds}}" \
            "${{@}}"
    """.format(
        dnscontrol = dnscontrol.dnscontrol.short_path,
        workspace_name = ctx.workspace_name,
        creds = ctx.file.creds.short_path if ctx.file.creds else "",
        arguments = " ".join([shell.quote(arg) for arg in args]),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )

    default_info = DefaultInfo(
        runfiles = runfiles,
        executable = script,
    )
    return [default_info]

al_dnscontrol_binary = rule(
    implementation = _impl,
    executable = True,
    doc = "Dnscontrol binary",
    toolchains = ["//tools/dnscontrol/main/bzl:toolchain_type"],
    attrs = {
        "config": attr.label(
            mandatory = True,
            allow_single_file = ["js"],
            doc = "Dnsconfig.js",
        ),
        "zones": attr.label_list(
            allow_files = ["zone"],
            doc = "Zone files",
        ),
        "creds": attr.label(
            allow_single_file = ["json.tpl"],
            doc = "Creds.json template",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Dnscontrol arguments",
        ),
    },
)
