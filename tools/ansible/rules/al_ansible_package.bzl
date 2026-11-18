"""Repository Ansible packages with mapped runtime files and AL commands."""

load("@rules_al//:defs.bzl", "al_binary_run_map")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files", "strip_prefix")
load(":al_ansible_binary.bzl", "al_ansible_binary")

def al_ansible_package(
        name,
        srcs,
        configs = [],
        plugin_labels = {},
        args = None,
        packaged_srcs = None,
        renames = {},
        srcs_name = "srcs",
        files_name = "srcs_files",
        binary_visibility = None,
        **kwargs):
    """Package sources and expose one or more Ansible playbook commands.

    Args:
        name: Command prefix; the packaged executable is `<name>_bin`.
        srcs: Files to map relative to the caller's package.
        configs: AL configuration dependencies for command wrappers.
        plugin_labels: AL plugin selectors for command wrappers.
        args: Command suffix to Ansible argument list, without the executable.
            Defaults to the inventory.yaml/playbook_deploy.yaml invocation.
        packaged_srcs: Mapped source groups; defaults to the files_name target.
        renames: Destination overrides for srcs.
        srcs_name: Combined runtime source target name.
        files_name: Mapped source target name.
        binary_visibility: Optional visibility for the packaged executable.
        **kwargs: Common target attributes for the package and commands.
    """
    pkg_files(
        name = files_name,
        srcs = srcs,
        renames = renames,
        strip_prefix = strip_prefix.from_pkg(),
        **kwargs
    )
    pkg_filegroup(
        name = srcs_name,
        srcs = packaged_srcs if packaged_srcs != None else [":" + files_name],
        **kwargs
    )
    binary_name = name + "_bin"
    binary_kwargs = dict(kwargs)
    if binary_visibility != None:
        binary_kwargs["visibility"] = binary_visibility
    al_ansible_binary(
        name = binary_name,
        data = [
            ":" + srcs_name,
            Label("//infra:ansible_cfg"),
            Label("//projects/ansible_collection:collections"),
        ],
        process_name = "ansible_playbook",
        **binary_kwargs
    )
    if args == None:
        args = {"": ["--inventory", "inventory.yaml", "playbook_deploy.yaml"]}
    al_binary_run_map(
        name = name,
        args = {suffix: ["$(rootpath :{})".format(binary_name)] + arguments for suffix, arguments in args.items()},
        configs = configs,
        data = [":" + binary_name],
        plugin_labels = plugin_labels,
        **kwargs
    )
