"""Repository Terraform operation maps with AL configuration composition."""

load("@rules_al//:defs.bzl", "al_binary_run", "al_binary_run_test")
load("@rules_terraform//terraform:defs.bzl", "terraform_binary_map", "terraform_target_binary_map", "terraform_test_map")

def al_terraform_binary_map(configs, plugin_labels = {}, run_args = [], **kwargs):
    """Create Terraform commands using the supplied AL configuration targets.

    Args:
        configs: AL configuration dependencies, including plugin runfiles.
        plugin_labels: Label/value selectors; empty values match presence.
        run_args: Other AL command options.
        **kwargs: Arguments to terraform_binary_map.
    """
    terraform_binary_map(
        wrapper = al_binary_run,
        wrapper_kwargs = {"configs": configs, "plugin_labels": plugin_labels, "run_args": run_args},
        **kwargs
    )

def al_terraform_target_binary_map(configs, plugin_labels = {}, run_args = [], **kwargs):
    """Create scoped Terraform commands retaining their saved-plan apply guard.

    Args:
        configs: AL configuration dependencies, including plugin runfiles.
        plugin_labels: Label/value selectors; empty values match presence.
        run_args: Other AL command options.
        **kwargs: Arguments to terraform_target_binary_map.
    """
    terraform_target_binary_map(
        wrapper = al_binary_run,
        wrapper_kwargs = {"configs": configs, "plugin_labels": plugin_labels, "run_args": run_args},
        **kwargs
    )

def al_terraform_test_map(configs, plugin_labels = {}, run_args = [], **kwargs):
    """Create Terraform checks using the supplied AL configuration targets.

    Args:
        configs: AL configuration dependencies, including plugin runfiles.
        plugin_labels: Label/value selectors; empty values match presence.
        run_args: Other AL command options.
        **kwargs: Arguments to terraform_test_map.
    """
    terraform_test_map(
        wrapper = al_binary_run_test,
        wrapper_kwargs = {"configs": configs, "plugin_labels": plugin_labels, "run_args": run_args},
        **kwargs
    )

def al_terraform_package(
        name,
        configs,
        data = [],
        providers = [],
        plugin_labels = {},
        args = None,
        tests = True,
        tests_name = None,
        tests_use_al = False,
        test_kwargs = {},
        command_kwargs = {},
        **kwargs):
    """Declare Terraform operations and checks with shared package inputs.

    Args:
        name: Operation prefix, such as tf or tf_setup.
        configs: AL configuration dependencies for operations.
        data: Shared runtime configuration and module inputs.
        providers: Shared pinned provider targets.
        plugin_labels: AL selectors for operations; not inherited by checks.
        args: Optional operation map; None retains the standard operations.
        tests: Whether to declare package checks.
        tests_name: Test prefix; defaults to `<name>_tests`.
        tests_use_al: Wrap checks with AL using configs, retaining legacy checks.
        test_kwargs: Check-specific overrides, including args, data or providers.
        command_kwargs: Operation-specific overrides, such as visibility.
        **kwargs: Shared Terraform map and common target attributes.
    """
    shared = dict(kwargs, data = data, providers = providers)
    commands = dict(shared, configs = configs, plugin_labels = plugin_labels)
    if args != None:
        commands["args"] = args
    al_terraform_binary_map(name = name, **(commands | command_kwargs))

    if not tests:
        return
    checks = dict(shared)
    if tests_use_al:
        checks["configs"] = configs
    checks.update(test_kwargs)
    test_map = al_terraform_test_map if tests_use_al else terraform_test_map
    test_map(name = tests_name if tests_name != None else name + "_tests", **checks)
