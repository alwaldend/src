load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")

def _impl(ctx):
    runfiles_files = []
    runfiles_symlinks = {}
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ansible = ctx.attr.ansible[ctx.attr.process_name]

    for symlink in ctx.attr.symlinks:
        for info, _ in symlink[PackageFilegroupInfo].pkg_files:
            runfiles_symlinks.update(info.dest_src_map)

    runfiles = ctx.runfiles(files = runfiles_files, symlinks = runfiles_symlinks)
    runfiles = runfiles.merge_all([ansible[DefaultInfo].default_runfiles])

    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        exec '{ansible}' \
            {arguments} \
            "${{@}}"
    """.format(
        ansible = ansible[DefaultInfo].files_to_run.executable.short_path,
        arguments = " ".join([shell.quote(arg) for arg in args]),
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
        "symlinks": attr.label_list(
            doc = "Symlinks",
            default = [],
            providers = [PackageFilegroupInfo],
        ),
        "ansible": attr.string_keyed_label_dict(
            default = {
                cli: "//bzl/rules/ansible:{}".format(cli)
                for cli in [
                    "ansible",
                    "ansible_config",
                    "ansible_console",
                    "ansible_doc",
                    "ansible_galaxy",
                    "ansible_inventory",
                    "ansible_playbook",
                    "ansible_pull",
                    "ansible_test",
                    "ansible_vault",
                ]
            },
            doc = "Ansible executable",
            cfg = "exec",
        ),
    },
)
