load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    collections = ctx.actions.declare_directory("{}.collections".format(ctx.label.name))
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ansible = ctx.attr.ansible[ctx.attr.process_name]
    runfiles = ctx.runfiles(
        files = [collections] + ctx.files.ansible,
        symlinks = {"collections": collections},
    )
    runfiles.merge(ansible[DefaultInfo].default_runfiles)
    runfiles.merge(ansible[DefaultInfo].data_runfiles)

    ctx.actions.run_shell(
        outputs = [collections],
        inputs = [ctx.file.collections],
        command = "tar -xf '{}' -C '{}'".format(ctx.file.collections.path, collections.path),
    )

    script_content = """\
        #!/usr/bin/env sh
        set -eux
        echo "${{PWD}}"
        find ../ -name "*python*"
        '{ansible}' \
            {arguments} \
            "${{@}}"
    """.format(
        collections = collections.short_path,
        ansible = ansible[DefaultInfo].files_to_run.executable.short_path,
        arguments = " ".join([
            shell.quote(argument)
            for argument in ctx.attr.arguments
        ]),
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
        "collections": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Ansible collections",
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
