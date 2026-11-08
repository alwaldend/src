"""Repository-specific AL authentication and injection target macros."""

load("@rules_al//:defs.bzl", "al_injector_call", "al_plugin")

def infra_xo_login(name, config_name = None, oidc_provider = "", labels = {}, **kwargs):
    """Declare Xen Orchestra login configuration."""
    if not oidc_provider:
        fail("infra_xo_login requires oidc_provider")
    al_plugin(
        name = name,
        config_name = config_name,
        labels = labels,
        executable = Label("//infra/xcp_ng/cmd/xo_login"),
        args = {
            "discovery_url": "https://vault.alwaldend.com:8200/v1/identity/oidc/provider/" + oidc_provider + "/.well-known/openid-configuration",
            "xoa_url": "https://xoa.xcp-ng.alwaldend.com",
        },
        **kwargs
    )

def infra_pve_provider_inputs(name, config_name = None, labels = {}, **kwargs):
    """Declare pve provider inputs configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_env = {"PM_API_URL": {"value": "https://pve.alwaldend.com:8006/api2/json"}},
        **kwargs
    )

def infra_ansible_keys(name, config_name = None, vault_ssh = {}, labels = {}, resource_name = "ansible_keys", **kwargs):
    """Declare ansible keys configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_vault_ssh = {resource_name: vault_ssh},
        res_env = {
            "ANSIBLE_PRIVATE_KEY_FILE": {"value": "{{ .Last.Data.private_key }}", "deps": [resource_name]},
            "SSH_AUTH_SOCK": {"value": "", "deps": [resource_name]},
        },
        **kwargs
    )

def infra_k3s_token(name, config_name = None, path = "", labels = {}, resource_name = "k3s_token", mount = "secrets", **kwargs):
    """Declare k3s token configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = {resource_name: {"mount": mount, "path": path}},
        res_env = {"K3S_TOKEN": {"value": "{{ .Last.Data.k3s_token }}", "deps": [resource_name]}},
        **kwargs
    )

def infra_kubernetes_login(name, config_name = None, oidc = {}, cluster_ca = "", labels = {}, resource_name = "kubernetes_login", **kwargs):
    """Declare kubernetes login configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_oidc = {resource_name: oidc},
        res_file = {
            resource_name + "_file": {
                "extra": {"cluster_ca": cluster_ca},
                "value": """apiVersion: v1
clusters:
  - cluster:
      certificate-authority-data: "{{ .Extra.cluster_ca }}"
      server: https://127.0.0.1:6443
    name: default
contexts:
  - context:
      cluster: default
      user: default
    name: default
current-context: default
kind: Config
users:
  - name: default
    user:
      token: "{{ .Last.Data.id_token }}"
""",
                "deps": [resource_name],
            },
        },
        res_env = {"KUBECONFIG": {"value": "{{ index .Last.Files 0 }}", "deps": [resource_name + "_file"]}},
        **kwargs
    )

def infra_yc_auth(name, config_name = None, path = "", labels = {}, resource_name = "yc_auth", filename = "service_account_key", **kwargs):
    """Declare yc auth configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = {resource_name: {"mount": "secrets", "path": path}},
        res_file = {filename: {"value": "{{ .Last.Data.service_account_key }}", "deps": [resource_name]}},
        res_env = {
            "YC_CLOUD_ID": {"value": "{{ .Last.Data.cloud_id }}", "deps": [resource_name]},
            "TF_VAR_cloud_id": {"value": "{{ .Last.Data.cloud_id }}", "deps": [resource_name]},
            "YC_FOLDER_ID": {"value": "{{ .Last.Data.folder_id }}", "deps": [resource_name]},
            "TF_VAR_folder_id": {"value": "{{ .Last.Data.folder_id }}", "deps": [resource_name]},
            "YC_SERVICE_ACCOUNT_KEY_FILE": {"value": "{{ index .Last.Files 0 }}", "deps": [filename]},
        },
        **kwargs
    )

def infra_yc_account(name, config_name = None, path = "", labels = {}, resource_name = "yc_account", **kwargs):
    """Declare yc account configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = {resource_name: {"mount": "secrets", "path": path}},
        res_env = {"TF_VAR_service_account_id": {"value": "{{ .Last.Data.service_account_id }}", "deps": [resource_name]}},
        **kwargs
    )

def infra_yc_bucket_auth(name, config_name = None, path = "", labels = {}, resource_name = "yc_bucket_auth", **kwargs):
    """Declare yc bucket auth configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = {resource_name: {"mount": "secrets", "path": path}},
        res_env = {
            "AWS_ACCESS_KEY_ID": {"value": "{{ .Last.Data.access_key }}", "deps": [resource_name]},
            "AWS_SECRET_ACCESS_KEY": {"value": "{{ .Last.Data.secret_key }}", "deps": [resource_name]},
        },
        **kwargs
    )

def infra_rclone_config(name, config_name = None, path = "", labels = {}, resource_name = "rclone_config", **kwargs):
    """Declare rclone config configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_file = {
            resource_name: {
                "value": """                    [remote]
                    type = s3
                    provider = AWS
                    env_auth = true
                    region = ru-central1
                    endpoint = storage.yandexcloud.net
                """,
            },
        },
        res_kv = {"rclone_bucket": {"mount": "secrets", "path": path}},
        res_env = {
            "RCLONE_CONFIG": {"value": "{{ index .Last.Files 0 }}", "deps": [resource_name]},
            "RCLONE_S3_BUCKET": {"value": "{{ .Last.Data.bucket }}", "deps": ["rclone_bucket"]},
        },
        **kwargs
    )

def infra_mikrotik(name, config_name = None, path = "", host = "", labels = {}, resource_name = "mikrotik", mount = "secrets", **kwargs):
    """Declare mikrotik configuration."""
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = {resource_name: {"mount": mount, "path": path}},
        res_env = {
            "MIKROTIK_HOST": {"value": host, "deps": [resource_name]},
            "MIKROTIK_USER": {"value": "{{ .Last.Data.mikrotik_username }}", "deps": [resource_name]},
            "MIKROTIK_PASSWORD": {"value": "{{ .Last.Data.mikrotik_password }}", "deps": [resource_name]},
        },
        **kwargs
    )

def infra_dns(name, config_name = None, labels = {}, global_records = True, dc1 = False, vault_auth = None, **kwargs):
    """Construct credentials for the selected DNS provider views.

    Args:
        labels: Plugin-call selection labels.
        global_records: Include the public DNS provider.
        dc1: Include the local DNS provider.
        vault_auth: Optional named Vault authentication for secret reads.
        **kwargs: Additional plugin-call attributes.

    Returns:
        Declares a plugin-call target with named resources.
    """
    kv = {}
    env = {}
    if global_records:
        kv["dns_cloudflare"] = {"path": "cloudflare.com/dns_token", "mount": "secrets"}
        env.update({
            "CLOUDFLARE_API_TOKEN": {"deps": ["dns_cloudflare"], "value": "{{ .Last.Data.cloudflare_api_token }}"},
            "TF_VAR_dns_cloudflare_zone_id": {"deps": ["dns_cloudflare"], "value": '{{ with index .Last.Data "cloudflare_zone_id" }}{{ . }}{{ end }}'},
        })
    if dc1:
        kv["dns_mikrotik"] = {"path": "alwaldend.com/vault1/approles/src_infra_dns/mikrotik", "mount": "secrets"}
        env.update({
            "TF_VAR_dns_routeros_hosturl": {"deps": ["dns_mikrotik"], "value": "https://router1.dc1.alwaldend.com"},
            "TF_VAR_dns_routeros_username": {"deps": ["dns_mikrotik"], "value": "{{ .Last.Data.username }}"},
            "TF_VAR_dns_routeros_password": {"deps": ["dns_mikrotik"], "value": "{{ .Last.Data.password }}"},
        })
    if vault_auth != None:
        for settings in kv.values():
            settings["vault_auth"] = vault_auth
    al_injector_call(
        name = name,
        config_name = config_name,
        labels = labels,
        res_kv = kv,
        res_env = env,
        **kwargs
    )
