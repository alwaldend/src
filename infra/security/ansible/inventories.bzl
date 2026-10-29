"""Authoritative owner inventories consumed by the security workflow."""

INVENTORIES = {
    "download": "//infra/download/ansible:inventory",
    "flux": "//infra/flux/ansible:inventory",
    "forgejo": "//infra/forgejo/ansible:inventory",
    "forgejo_runner": "//infra/forgejo_runner/ansible:inventory",
    "harbor": "//infra/harbor/ansible:inventory",
    "host_bot": "//users/simeonwarren/host_bot/ansible:inventory",
    "ingress": "//infra/ingress/ansible:inventory",
    "openhands": "//infra/openhands/ansible:inventory",
    "pve": "//infra/pve/ansible:inventory",
    "simeonwarren": "//users/simeonwarren/ansible:inventory",
    "threexui": "//infra/threexui/ansible:inventory",
    "vault": "//infra/vault/ansible:inventory",
    "xcp_ng": "//infra/xcp_ng/ansible:inventory",
}
