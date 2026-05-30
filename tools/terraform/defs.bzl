load("@bazel_skylib//rules:native_binary.bzl", "native_test")
load("@rules_go//go:def.bzl", "go_binary", "go_test")
load("//projects/al/rules/al:al_binary.bzl", "al_binary_run", "al_binary_run_test")

def terraform_binary(
        args = [],
        data = [],
        terraform = "//tools/terraform",
        terraform_runner = "//tools/terraform/runner",
        is_test = False,
        **kwargs):
    """
    Create a binary for terraform commands
    """
    rule = al_binary_run
    if is_test:
        rule = al_binary_run_test
    rule(
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

DEFAULT_TERRAFORM_BINARIES = {
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

def terraform_binary_map(
        name,
        args = DEFAULT_TERRAFORM_BINARIES,
        **kwargs):
    """
    Create several terraform binaries
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        terraform_binary(
            name = cur_name,
            args = args_value,
            **kwargs
        )

DEFAULT_TERRAFORM_TESTS = {
    "fmt": ["--direct", "fmt", "--check", "--recursive"],
}

def terraform_test_map(
        name,
        args = DEFAULT_TERRAFORM_TESTS,
        size = "small",
        **kwargs):
    """
    Create several terraform tests
    """
    for args_name, args_value in args.items():
        if not args_name:
            fail("empty test name")
        terraform_binary(
            name = "{}.{}_test".format(name, args_name),
            size = size,
            args = args_value,
            is_test = True,
            **kwargs
        )
