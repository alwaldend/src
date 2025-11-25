load("//tools/ansible/main/bzl:al_ansible_binary.bzl", "al_ansible_binary")
load("//tools/ansible/main/bzl:al_ansible_scripts.bzl", "AL_ANSIBLE_SCRIPTS")

def al_ansible_binaries(name, args = {}, visibility = None, **kwargs):
    """
    Generate a binary for earch ansible command

    Args:
        name: name prefix
        visibility: visibility
        args: specific kwargs
        **kwargs: common arguments
    """

    for script in AL_ANSIBLE_SCRIPTS:
        al_ansible_binary(
            name = "{}.{}".format(name, script),
            process_name = script,
            visibility = visibility,
            **(kwargs | args.get(script, {}))
        )
