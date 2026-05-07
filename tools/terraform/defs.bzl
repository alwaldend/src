load("@rules_go//go:def.bzl", "go_binary")
load("//tools/al/rules/al:al_binary.bzl", "al_binary_run")

def terraform_runner(
        args = [],
        data = [],
        terraform = "//tools/terraform",
        terraform_runner = "//tools/terraform/runner",
        **kwargs):
    """
    Create a binary for terraform commands
    """
    al_binary_run(
        args = [
            "$(rootpath {})".format(terraform_runner),
            "--terraform",
            "$(rootpath {})".format(terraform),
            "--chdir",
            native.package_name() or ".",
        ] + args,
        data = data + [terraform, terraform_runner],
        **kwargs
    )

DEFAULT_TERRAFORM_RUNNERS = {
    "": [],
    "direct": ["--direct"],
    "init": ["--direct", "init"],
    "migrate": ["--direct", "init", "--migrate-state"],
    "fmt": ["--direct", "fmt", "--write", "--recursive"],
    "fmt_check": ["--direct", "fmt", "--check", "--recursive"],
    "plan": ["plan"],
    "output": ["output"],
    "import": ["import"],
    "destroy": ["destroy"],
    "apply": ["apply"],
    "show": ["show"],
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
