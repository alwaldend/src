"""Declare a package's staged Helm and OCI artifact mirror."""

load("//projects/al:defs.bzl", "al_binary_run")

_COMMAND = "//tools/artifact_mirror"
_HELM = "//tools/helm"
_ORAS = "//third_party/land_oras_oras:oras"

def third_party_artifact_mirror(inventory):
    """Declare pull and push targets for one upstream package.

    Args:
        inventory: Label of the package's chart and image pin file.
    """
    common_args = [
        "--inventory",
        "$(rootpath {})".format(inventory),
        "--oras",
        "$(rootpath {})".format(_ORAS),
    ]
    common_data = [inventory, _COMMAND, _ORAS]
    command = "$(rootpath {})".format(_COMMAND)

    al_binary_run(
        name = "pull",
        args = [command, "pull", "--helm", "$(rootpath {})".format(_HELM)] + common_args,
        data = common_data + [_HELM],
    )
    al_binary_run(
        name = "push",
        args = [command, "push", "--owner", "alwaldend"] + common_args,
        configs = ["//third_party:al"],
        data = common_data,
        plugin_labels = {"forgejo": ""},
    )
