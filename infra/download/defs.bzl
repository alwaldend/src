"""infra download configuration targets."""

load("@rules_al//:defs.bzl", "al_injector_call")

def infra_download_ssh(name, labels = {}, config_name = None, **kwargs):
    """Declare the ssh configuration target."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        substitutions = {"//data/ssh:ca_servers": "@AL_SSH_SERVER_CA@"},
        res_vault_ssh = {"ssh_key": {"backend": "ssh/clients/sign/admins", "ttl": 7200}},
        res_env = {
            "DOWNLOAD_SSH_PRIVATE_KEY": {"value": "{{ .Last.Data.private_key }}", "deps": ["ssh_key"]},
            "DOWNLOAD_SSH_CERTIFICATE": {"value": "{{ .Last.Data.public_key_signed }}", "deps": ["ssh_key"]},
            "DOWNLOAD_SSH_KNOWN_HOSTS": {"value": "{{ index .Last.Files 0 }}", "deps": ["ssh_known_hosts"]},
        },
        res_file = {
            "ssh_known_hosts": {
                "value": """@cert-authority *.download.alwaldend.com @AL_SSH_SERVER_CA@
""",
            },
        },
        **kwargs
    )
