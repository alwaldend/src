load("@rules_go//go:def.bzl", "go_binary")
load("//tools/al/rules/al:al_binary.bzl", "al_binary_run")

def vault_runner(
        args = [],
        run_args = [],
        data = [],
        vault_env = "default:default",
        vault = "//tools/vault",
        **kwargs):
    """
    Create a vault runner
    """
    al_binary_run(
        run_args = ["--env_vault", vault_env] + run_args,
        args = ["$(rootpath {})".format(vault)] + args,
        data = data + [vault],
        **kwargs
    )

DEFAULT_VAULT_RUNNERS = {
    "": [],
    "kv_get": ["kv", "get", "-mount", "secrets"],
    "kv_put": ["kv", "put", "-mount", "secrets"],
    "kv_patch": ["kv", "patch", "-mount", "secrets"],
    "kv_delete": ["kv", "delete", "-mount", "secrets"],
    "login": ["login", "--method=userpass"],
    "token_lookup": ["token", "lookup"],
    "status": ["status"],
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
