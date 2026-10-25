"""Validate proxy configuration syntax without exposing enrollment secrets."""

import collections.abc
import pathlib

import ansible.errors
import tomllib
import yaml


def _validate_config(content: str, filename: str) -> str:
    """Return the original template only when it parses as a configuration map."""
    parsed: object
    try:
        if pathlib.PurePath(filename).suffix in (".yaml", ".yml"):
            parsed = yaml.safe_load(content)
        else:
            parsed = tomllib.loads(content)
    except (ValueError, yaml.YAMLError):
        raise ansible.errors.AnsibleFilterError(
            "Invalid Traefik configuration syntax; candidate was not installed."
        ) from None
    if not isinstance(parsed, dict):
        raise ansible.errors.AnsibleFilterError(
            "Traefik configuration must be a mapping."
        )
    return content


class FilterModule:
    """Expose the collection's controller-side configuration validator."""

    def filters(self) -> dict[str, collections.abc.Callable[[str, str], str]]:
        """Register the Ansible filter without importing remote host software."""
        return {"validate_traefik_config": _validate_config}
