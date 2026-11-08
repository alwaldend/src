"""Vault Terraform backend configuration targets."""

load("@rules_al//:defs.bzl", "al_plugin")

def vault_tf_backend(name, path, labels = {}, mount = "secrets", vault_auth = "", vault_conn = "", config_name = None, **kwargs):
    """Construct a backend call for the shared tf_backend plugin registration.

    Args:
        name: Bazel target name.
        path: Vault KV path owning this backend's state and lock records.
        labels: Call selection labels identifying the Terraform stage.
        mount: Vault KV mount containing the backend records.
        vault_auth: Optional named Vault authentication method.
        vault_conn: Optional named Vault connection.
        config_name: Optional runtime plugin-call name.
        **kwargs: Additional plugin-call target attributes.
    """
    args = {"vault_secret": path, "vault_secret_mount": mount}
    if vault_auth:
        args["vault_auth"] = vault_auth
    if vault_conn:
        args["vault_conn"] = vault_conn
    al_plugin(name = name, config_name = config_name, plugin = "tf_backend", labels = labels, args = args, **kwargs)
