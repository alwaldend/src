"""infra forgejo configuration targets."""

load("@rules_al//:defs.bzl", "al_inject")

def infra_forgejo_config(name, labels = {}, config_name = None, **kwargs):
    """Declare the config configuration target."""
    al_inject(
        name = name,
        config_name = config_name,
        labels = labels,
        res_op = {"forgejo_oidc_client": {"method": "read", "path": "identity/oidc/client/src_infra_dc1_forgejo1_provider"}},
        res_env = {
            "FORGEJO_OIDC_CLIENT_ID": {"value": "{{ .Last.Data.client_id }}", "deps": ["forgejo_oidc_client"]},
            "FORGEJO_OIDC_CLIENT_SECRET": {"value": "{{ .Last.Data.client_secret }}", "deps": ["forgejo_oidc_client"]},
            "LFS_JWT_SECRET": {"value": "{{ .Last.Data.lfs_jwt_secret }}", "deps": ["config"]},
            "SECURITY_INTERNAL_TOKEN": {"value": "{{ .Last.Data.security_internal_token }}", "deps": ["config"]},
            "OAUTH2_JWT_SECRET": {"value": "{{ .Last.Data.oauth2_jwt_secret }}", "deps": ["config"]},
        },
        res_kv = {"config": {"mount": "secrets", "path": "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/config"}},
        **kwargs
    )
