load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")
load("//tools/secrets/main/bzl:al_secrets_binary_info.bzl", "AlSecretsBinaryInfo")
load(":al_ansible_scripts.bzl", "AL_ANSIBLE_SCRIPTS")

_SCRIPT = """\
#!/usr/bin/env sh
set -eu
export RUNFILES_DIR="${{RUNFILES_DIR:-$(dirname "${{PWD}}")}}"
for root in "" "${{0}}.runfiles/{workspace_name}/" "${{RUNFILES_DIR:-}}/{workspace_name}/"; do
    if [ -x "${{root}}{bin}" ]; then
        for source_script in {source_scripts}; do
            . "${{root}}${{source_script}}"
            rm "${{root}}${{source_script}}"
        done
        exec "${{root}}/{bin}" {arguments} "${{@}}"
    fi
done
echo "Could not find the binary"
exit 1
"""

def _impl(ctx):
    runfiles_files = []
    runfiles_symlinks = {}
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ansible = ctx.attr.ansible[ctx.attr.process_name]
    source_scripts = []

    for symlink in ctx.attr.data:
        for info, _ in symlink[PackageFilegroupInfo].pkg_files:
            runfiles_symlinks.update(info.dest_src_map)

    runfiles = ctx.runfiles(files = runfiles_files, symlinks = runfiles_symlinks)
    runfiles = runfiles.merge_all([ansible[DefaultInfo].default_runfiles])

    for attr in ctx.attr.secrets:
        info = attr[AlSecretsBinaryInfo]
        runfiles = runfiles.merge(info.default_info.default_runfiles)
        source_scripts.append(info.source_script.short_path)

    args = []
    args.extend(ctx.attr.arguments)
    script_content = _SCRIPT.format(
        bin = ansible[DefaultInfo].files_to_run.executable.short_path,
        arguments = " ".join([shell.quote(arg) for arg in args]),
        source_scripts = " ".join([shell.quote(source_script) for source_script in source_scripts]),
        workspace_name = ctx.workspace_name,
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )

    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = script,
        ),
    ]

al_ansible_binary = rule(
    implementation = _impl,
    executable = True,
    doc = "Run an ansible command",
    attrs = {
        "process_name": attr.string(
            mandatory = True,
            doc = "Process name",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Ansible arguments",
        ),
        "data": attr.label_list(
            doc = "Data",
            default = [],
            providers = [PackageFilegroupInfo],
        ),
        "secrets": attr.label_list(
            cfg = "exec",
            providers = [AlSecretsBinaryInfo],
            doc = "Secrets",
        ),
        "ansible": attr.string_keyed_label_dict(
            default = {
                cli: "//tools/ansible:{}".format(cli)
                for cli in AL_ANSIBLE_SCRIPTS
            },
            doc = "Ansible executable",
            cfg = "exec",
        ),
    },
)
