"""Preserve Ansible's CLI identity for its native SSH password helper."""

import runpy
import sys


def _main() -> None:
    """Set the expected name before Ansible initializes its askpass dispatcher."""
    sys.argv[0] = "ansible"
    runpy.run_module("ansible.cli.adhoc", run_name="__main__")


if __name__ == "__main__":
    _main()
