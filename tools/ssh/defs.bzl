"""Run the host SSH client through AL's injection lifecycle."""

load("//projects/al/rules/al:al_binary.bzl", "al_binary_run")

def ssh_binary(name, args = [], **kwargs):
    """Create an AL command that resolves ssh from the host PATH.

    Args:
        name: Target name.
        args: Fixed SSH arguments; invocation arguments are appended.
        **kwargs: Additional al_binary_run attributes, including configs,
            run_args, and data for consumer-owned key injection.
    """
    al_binary_run(
        name = name,
        args = ["ssh"] + args,
        **kwargs
    )
