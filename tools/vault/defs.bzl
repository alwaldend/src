load("@rules_go//go:def.bzl", "go_binary")

def vault_runner(
        name,
        config,
        args = [],
        run_args = [],
        vault_env = "default:default",
        secrets = [],
        data = [],
        vault = "//tools/vault",
        al_lib = "//tools/al/cmd/al:al_lib"):
    """
    Create a vault runner
    """
    data = data + [config, vault]
    args = [
        "run",
        "--config",
        "$(rootpath {})".format(config),
        "--vault_env",
        vault_env,
    ] + run_args + [
        "--",
        "$(rootpath {})".format(vault),
    ] + args
    go_binary(
        name = name,
        args = args,
        data = data,
        embed = [al_lib],
    )

DEFAULT_VAULT_RUNNERS = {
    "": [],
    "kv_get": ["kv", "get", "-mount", "secrets"],
    "kv_put": ["kv", "put", "-mount", "secrets"],
    "kv_patch": ["kv", "patch", "-mount", "secrets"],
    "kv_delete": ["kv", "delete", "-mount", "secrets"],
    "login": ["login", "--method=userpass"],
    "token_lookup": ["token", "lookup"],
}

def vault_runners(
        name,
        args = DEFAULT_VAULT_RUNNERS,
        **kwargs):
    """
    Create several vault runners
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        vault_runner(
            name = cur_name,
            args = args_value,
            **kwargs
        )
