load("@rules_al//:defs.bzl", "al_binary_run")

def vault_binary(
        args = [],
        run_args = [],
        plugin_labels = {},
        data = [],
        vault_env = "default_default",
        vault = "//tools/vault",
        **kwargs):
    """
    Create a vault runner
    """
    al_binary_run(
        plugin_labels = dict({"vault_env": vault_env}, **plugin_labels),
        run_args = run_args,
        args = ["$(rootpath {})".format(vault)] + args,
        data = data + [vault],
        **kwargs
    )

DEFAULT_VAULT_BINARY_MAP_ARGS = {
    "": [],
    "kv_get": ["kv", "get", "-mount", "secrets"],
    "kv_put": ["kv", "put", "-mount", "secrets"],
    "kv_patch": ["kv", "patch", "-mount", "secrets"],
    "kv_delete": ["kv", "delete", "-mount", "secrets"],
    "login": ["login", "--method=userpass"],
    "token_lookup": ["token", "lookup"],
    "status": ["status"],
}

def vault_binary_map(
        name,
        args = DEFAULT_VAULT_BINARY_MAP_ARGS,
        **kwargs):
    """
    Create several vault runners
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        vault_binary(
            name = cur_name,
            args = args_value,
            **kwargs
        )
