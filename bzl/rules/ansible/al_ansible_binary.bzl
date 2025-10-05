load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    runfiles_files = [ctx.file.config] + ctx.files.inventories
    if ctx.file.collections:
        collections = ctx.actions.declare_directory("{}.collections".format(ctx.label.name))
        runfiles_files.append(collections)
    else:
        collections = None
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    ansible = ctx.attr.ansible[ctx.attr.process_name]
    runfiles_symlinks = {}
    if ctx.file.config:
        runfiles_symlinks["ansible.cfg"] = ctx.file.config
    if collections:
        runfiles_symlinks["collections/ansible_collections"] = collections

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

    if collections:
        ctx.actions.run_shell(
            outputs = [collections],
            inputs = [ctx.file.collections],
            command = "tar -xf '{}' -C '{}'".format(ctx.file.collections.path, collections.path),
        )

    args = []
    for inventory in ctx.files.inventories:
        args.extend(["--inventory", inventory.short_path])
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
        "collections": attr.label(
            allow_single_file = [".tar"],
            doc = "Ansible collections",
        ),
        "inventories": attr.label_list(
            allow_files = True,
            default = [],
            doc = "Ansible inventories",
        ),
        "config": attr.label(
            allow_single_file = True,
            doc = "Ansible config",
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
