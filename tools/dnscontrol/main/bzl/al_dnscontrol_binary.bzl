load("@bazel_skylib//lib:shell.bzl", "shell")
load("//tools/secrets/main/bzl:al_secrets_binary_info.bzl", "AlSecretsBinaryInfo")

_SCRIPT = """\
#!/usr/bin/env sh
set -eu
echo "Running dnscontrol script {dnscontrol}" >&2
export RUNFILES_DIR="${{RUNFILES_DIR:-$(dirname "${{PWD}}")}}"
for root in "" "${{0}}.runfiles/{workspace_name}/" "${{RUNFILES_DIR:-}}/{workspace_name}/"; do
    if [ -x "${{root}}{dnscontrol}" ]; then
        for source_script in {source_scripts}; do
            . "${{root}}${{source_script}}"
        done
        creds="$(mktemp)"
        envsubst <'{creds}' >"${{creds}}"
        if "${{root}}{dnscontrol}" \
            {arguments} \
            --creds "${{creds}}" \
            "${{@}}"; then
            rm "${{creds}}"
            exit
        else
            rm "${{creds}}"
            exit 1
        fi
    fi
done
echo "Could not find the binary"
exit 1
"""

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
    source_scripts = []
    runfiles = ctx.runfiles(
        files = runfiles_files,
        symlinks = runfiles_symlinks,
    )
    runfiles = runfiles.merge(dnscontrol.default_info.default_runfiles)
    for attr in ctx.attr.secrets:
        info = attr[AlSecretsBinaryInfo]
        runfiles = runfiles.merge(info.default_info.default_runfiles)
        source_scripts.append(info.source_script.short_path)

    args = []
    args.extend(ctx.attr.arguments)
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = _SCRIPT.format(
            dnscontrol = dnscontrol.dnscontrol.short_path,
            workspace_name = ctx.workspace_name,
            creds = ctx.file.creds.short_path if ctx.file.creds else "",
            arguments = " ".join([shell.quote(arg) for arg in args]),
            source_scripts = " ".join([shell.quote(source_script) for source_script in source_scripts]),
        ),
    )

    default_info = DefaultInfo(
        runfiles = runfiles,
        executable = script,
    )
    return [
        default_info,
        RunEnvironmentInfo(
            inherited_environment = ctx.attr.env_inherit,
        ),
    ]

_COMMON = {
    "implementation": _impl,
    "doc": "Dnscontrol binary",
    "toolchains": ["//tools/dnscontrol/main/bzl:toolchain_type"],
    "attrs": {
        "config": attr.label(
            mandatory = True,
            allow_single_file = ["js"],
            doc = "Dnsconfig.js",
        ),
        "secrets": attr.label_list(
            cfg = "exec",
            providers = [AlSecretsBinaryInfo],
            doc = "Secrets",
        ),
        "env_inherit": attr.string_list(
            doc = "Inherit env",
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
}

al_dnscontrol_binary = rule(
    executable = True,
    **_COMMON
)

al_dnscontrol_binary_test = rule(
    test = True,
    **_COMMON
)
