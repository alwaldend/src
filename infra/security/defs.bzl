"""infra security configuration targets."""

load("@rules_al//:defs.bzl", "al_injector_call")

def infra_security_maintenance_ssh(name, labels = {}, config_name = None, **kwargs):
    """Declare the maintenance ssh configuration target."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        substitutions = {"//data/ssh:ca_servers": "@AL_SSH_SERVER_CA@"},
        res_vault_ssh = {"maintenance_key": {"backend": "ssh/clients/sign/admins", "ttl": 7200}},
        res_file = {
            "maintenance_known_hosts": {
                "value": """@cert-authority *.alwaldend.com @AL_SSH_SERVER_CA@
""",
            },
            "maintenance_config": {
                "value": """Host *
    User ansible
    BatchMode yes
    ConnectTimeout 10
    ServerAliveInterval 15
    ServerAliveCountMax 2
    IdentityAgent none
    IdentitiesOnly yes
    IdentityFile "{{ .Res.maintenance_key.Data.private_key }}"
    CertificateFile "{{ .Res.maintenance_key.Data.public_key_signed }}"
    UserKnownHostsFile "{{ index .Res.maintenance_known_hosts.Files 0 }}"
    GlobalKnownHostsFile /dev/null
    StrictHostKeyChecking yes
    UpdateHostKeys no
    ControlMaster no
    ControlPath none
""",
                "deps": [
                    "maintenance_key",
                    "maintenance_known_hosts",
                ],
            },
        },
        res_env = {
            "HOST_MAINTENANCE_SSH_KNOWN_HOSTS": {"value": "{{ index .Last.Files 0 }}", "deps": ["maintenance_known_hosts"]},
            "HOST_MAINTENANCE_SSH_CONFIG": {"value": "{{ index .Last.Files 0 }}", "deps": ["maintenance_config"]},
            "ANSIBLE_SSH_ARGS": {"value": '-F "{{ index .Last.Files 0 }}"', "deps": ["maintenance_config"]},
            "ANSIBLE_PRIVATE_KEY_FILE": {"value": "{{ .Last.Data.private_key }}", "deps": ["maintenance_key"]},
        },
        **kwargs
    )
