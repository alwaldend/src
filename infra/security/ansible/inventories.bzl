"""Authoritative owner inventories consumed by the security workflow."""

INVENTORIES = {
    "download": "//infra/download/ansible:inventory",
    "flux": "//infra/flux/ansible:inventory",
    "forgejo": "//infra/forgejo/ansible:inventory",
    "dev": "//infra/dev/ansible:inventory",
    "harbor": "//infra/harbor/ansible:inventory",
    "ingress": "//infra/ingress/ansible:inventory",
    "openhands": "//infra/openhands/ansible:inventory",
    "pve": "//infra/pve/ansible:inventory",
    "simeonwarren": "//users/simeonwarren/ansible:inventory",
    "threexui": "//infra/threexui/ansible:inventory",
    "vault": "//infra/vault/ansible:inventory",
    "xcp_ng": "//infra/xcp_ng/ansible:inventory",
}
