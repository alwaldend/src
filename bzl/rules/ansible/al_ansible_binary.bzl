load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    runfiles_files = []
    runfiles_symlinks = {}
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ansible = ctx.attr.ansible[ctx.attr.process_name]

    for i, (symlink, file) in enumerate(ctx.attr.symlink_files.items()):
        runfiles_symlinks[symlink] = file.files.to_list()[0]

    for i, (symlink, archive) in enumerate(ctx.attr.symlink_archives.items()):
        dir = ctx.actions.declare_directory("{}.symlink_archive_{}".format(ctx.label.name, i))
        runfiles_symlinks[symlink] = dir
        archive_file = archive.files.to_list()[0]
        ctx.actions.run_shell(
            outputs = [dir],
            inputs = [archive_file],
            command = "tar -xf '{}' -C '{}'".format(archive_file.path, dir.path),
        )

    runfiles = ctx.runfiles(
        files = runfiles_files,
        transitive_files = depset(transitive = [
            ansible[DefaultInfo].default_runfiles.files,
            ansible[DefaultInfo].data_runfiles.files,
        ]),
        symlinks = runfiles_symlinks,
    )

    # this merge doesn't work for some reason
    runfiles.merge_all([
        ansible[DefaultInfo].default_runfiles,
        ansible[DefaultInfo].data_runfiles,
    ])

    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        '{ansible}' \
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
            mandatory = True,
            doc = "Ansible arguments",
        ),
        "symlink_archives": attr.string_keyed_label_dict(
            default = {},
            allow_files = [".tar"],
            doc = "Archives to symlink",
        ),
        "symlink_files": attr.string_keyed_label_dict(
            default = {},
            allow_files = True,
            doc = "Files to symlink",
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
