load("@rules_go//go:def.bzl", "go_binary")

def terraform_runner(
        name,
        config,
        args = [],
        run_args = [],
        secrets = [],
        data = [],
        terraform = "//tools/terraform",
        terraform_runner = "//tools/terraform/runner",
        al_lib = "//tools/al/cmd/al:al_lib"):
    """
    Create binaries for terraform commands
    """
    data = data + [config, terraform, terraform_runner]
    args = [
        "run",
        "--config",
        "$(rootpath {})".format(config),
    ] + run_args + [
        "--",
        "$(rootpath {})".format(terraform_runner),
        "--terraform",
        "$(rootpath {})".format(terraform),
        "--chdir",
        native.package_name() or ".",
    ] + args
    go_binary(
        name = name,
        args = args,
        data = data,
        embed = [al_lib],
    )

DEFAULT_TERRAFORM_RUNNERS = {
    "": ["--direct"],
    "init": ["--direct", "init"],
    "migrate": ["--direct", "init", "--migrate-state"],
    "fmt": ["--direct", "fmt", "--write", "--recursive"],
    "fmt_check": ["--direct", "fmt", "--check", "--recursive"],
    "plan": ["plan"],
    "import": ["import"],
    "destroy": ["destroy"],
    "apply": ["apply"],
    "deploy": ["apply"],
    "deploy_y": ["apply", "-y"],
}

def terraform_runners(
        name,
        args = DEFAULT_TERRAFORM_RUNNERS,
        **kwargs):
    """
    Create several terraform runners
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        terraform_runner(
            name = cur_name,
            args = args_value,
            **kwargs
        )
