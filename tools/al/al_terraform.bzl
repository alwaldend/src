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
