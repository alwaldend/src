#!/usr/bin/env python3
"""Run host preparation and Kolla-Ansible without nested Ansible."""

from __future__ import annotations

import argparse
import json
import os
import shlex
import shutil
import stat
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Mapping, Sequence

from render_inventory import render as render_kolla_inventory

PREPARE_ROOT_MARKER = "openstack_prepare_root.marker"

SUPPORTED_ACTIONS = (
    "deploy",
    "post-deploy",
    "prechecks",
    "reconfigure",
)

KOLLA_ACTIONS: Mapping[str, tuple[str, ...]] = {
    "prechecks": ("bootstrap-servers", "prechecks"),
    "deploy": (
        "bootstrap-servers",
        "prechecks",
        "deploy",
        "post-deploy",
    ),
    "reconfigure": ("prechecks", "reconfigure", "post-deploy"),
    "post-deploy": ("post-deploy",),
}

ACTIONS_WITH_HOST_PREPARATION = frozenset(
    {"deploy", "prechecks", "reconfigure"}
)
ACTIONS_WITH_SERVICE_CHECKS = frozenset(
    {"deploy", "post-deploy", "reconfigure"}
)


class LauncherError(RuntimeError):
    """Report a launcher configuration or execution error."""


@dataclass(frozen=True)
class Host:
    """One host from the repository inventory."""

    inventory_name: str
    ansible_host: str
    ansible_user: str
    ansible_port: int
    python: str
    service_hostname: str
    management_interface: str
    external_interface: str | None = None


@dataclass(frozen=True)
class Topology:
    """The fixed master and compute topology."""

    master: Host
    compute: Host


def parse_args(argv: Sequence[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--action", choices=SUPPORTED_ACTIONS, required=True)
    parser.add_argument("--config", type=Path, required=True)
    parser.add_argument("--inventory", type=Path, required=True)
    parser.add_argument("--prepare-bin", type=Path)
    return parser.parse_args(argv)


def read_json_object(path: Path) -> dict[str, Any]:
    try:
        document = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        raise LauncherError(f"Cannot read JSON object from {path}: {error}") from error
    if not isinstance(document, dict):
        raise LauncherError(f"{path} must contain a JSON object")
    return document


def required_string(
    document: Mapping[str, Any],
    key: str,
    *,
    allow_change_me: bool = False,
) -> str:
    value = document.get(key)
    if not isinstance(value, str) or not value.strip():
        raise LauncherError(f"{key} must be a non-empty string")
    if not allow_change_me and "CHANGE_ME" in value:
        raise LauncherError(f"Replace the placeholder in {key}")
    return value


def required_int(document: Mapping[str, Any], key: str) -> int:
    value = document.get(key)
    if isinstance(value, bool) or not isinstance(value, int):
        raise LauncherError(f"{key} must be an integer")
    return value


def required_bool(document: Mapping[str, Any], key: str) -> bool:
    value = document.get(key)
    if not isinstance(value, bool):
        raise LauncherError(f"{key} must be a Boolean")
    return value


def validate_config(config: Mapping[str, Any]) -> None:
    required_string(config, "openstack_release")
    kolla_revision = required_string(config, "openstack_kolla_revision")
    requirements_revision = required_string(
        config, "openstack_requirements_revision"
    )
    required_string(config, "openstack_kolla_base_distro")
    required_string(config, "openstack_kolla_container_engine")
    required_string(config, "openstack_kolla_python")
    required_string(config, "openstack_kolla_internal_vip_address")
    required_string(config, "openstack_fast_raid_device")
    required_string(config, "openstack_bulk_raid_device")
    required_string(config, "openstack_fast_vg")
    required_string(config, "openstack_bulk_vg")
    required_string(config, "openstack_glance_data_dir")
    required_string(config, "openstack_nova_datadir")
    required_string(config, "openstack_compute_filesystem")
    required_string(config, "openstack_compute_mount_options")
    required_bool(config, "openstack_manage_storage")
    required_bool(config, "openstack_storage_wipe")
    required_bool(config, "openstack_manage_compute_storage")
    required_bool(config, "openstack_compute_storage_wipe")

    router_id = required_int(
        config, "openstack_kolla_keepalived_virtual_router_id"
    )
    if not 1 <= router_id <= 255:
        raise LauncherError(
            "openstack_kolla_keepalived_virtual_router_id must be 1 through 255"
        )

    minimum_bytes = required_int(
        config, "openstack_compute_minimum_bytes"
    )
    if minimum_bytes < 1_800_000_000_000:
        raise LauncherError(
            "openstack_compute_minimum_bytes must be at least 1.8 TB"
        )

    reserved_mb = required_int(
        config, "openstack_nova_reserved_host_disk_mb"
    )
    if reserved_mb < 1024:
        raise LauncherError(
            "openstack_nova_reserved_host_disk_mb must reserve at least 1 GB"
        )

    if len(kolla_revision) != 40 or len(requirements_revision) != 40:
        raise LauncherError("OpenStack source revisions must be full Git SHAs")

    actions = config.get("openstack_supported_actions")
    if actions != list(SUPPORTED_ACTIONS):
        raise LauncherError(
            "openstack_supported_actions must match the launcher actions"
        )


def kolla_source(config: Mapping[str, Any]) -> str:
    revision = required_string(config, "openstack_kolla_revision")
    return (
        "git+https://opendev.org/openstack/kolla-ansible.git@"
        f"{revision}"
    )


def kolla_constraints(config: Mapping[str, Any]) -> str:
    revision = required_string(config, "openstack_requirements_revision")
    return (
        "https://opendev.org/openstack/requirements/raw/commit/"
        f"{revision}/upper-constraints.txt"
    )


def inventory_hosts(
    inventory: Mapping[str, Any],
    group: str,
) -> Mapping[str, Any]:
    try:
        hosts = inventory["all"]["children"]["openstack"]["children"][
            group
        ]["hosts"]
    except (KeyError, TypeError) as error:
        raise LauncherError(f"Inventory group {group} is missing") from error
    if not isinstance(hosts, dict) or len(hosts) != 1:
        raise LauncherError(f"Inventory group {group} must contain one host")
    return hosts


def parse_host(
    inventory_name: str,
    variables: Mapping[str, Any],
    *,
    master: bool,
) -> Host:
    if not isinstance(variables, dict):
        raise LauncherError(f"Host {inventory_name} must contain variables")

    ansible_host = variables.get("ansible_host", inventory_name)
    if not isinstance(ansible_host, str) or not ansible_host:
        raise LauncherError(f"Host {inventory_name} has an invalid ansible_host")

    ansible_port = variables.get("ansible_port", 22)
    if isinstance(ansible_port, bool) or not isinstance(ansible_port, int):
        raise LauncherError(f"Host {inventory_name} has an invalid ansible_port")
    if not 1 <= ansible_port <= 65535:
        raise LauncherError(f"Host {inventory_name} has an invalid ansible_port")

    host = Host(
        inventory_name=inventory_name,
        ansible_host=ansible_host,
        ansible_user=required_string(variables, "ansible_user"),
        ansible_port=ansible_port,
        python=required_string(
            variables, "ansible_python_interpreter"
        ),
        service_hostname=required_string(
            variables, "openstack_service_hostname"
        ),
        management_interface=required_string(
            variables, "openstack_management_interface"
        ),
        external_interface=(
            required_string(variables, "openstack_external_interface")
            if master
            else None
        ),
    )
    if host.ansible_user != "ansible":
        raise LauncherError(
            f"Host {inventory_name} must use the Vault-signed ansible user"
        )
    return host


def parse_topology(inventory: Mapping[str, Any]) -> Topology:
    master_hosts = inventory_hosts(inventory, "openstack_master")
    compute_hosts = inventory_hosts(inventory, "openstack_compute")
    master_name, master_variables = next(iter(master_hosts.items()))
    compute_name, compute_variables = next(iter(compute_hosts.items()))
    if master_name == compute_name:
        raise LauncherError("The master and compute inventory hosts must differ")

    topology = Topology(
        master=parse_host(master_name, master_variables, master=True),
        compute=parse_host(compute_name, compute_variables, master=False),
    )
    if topology.master.service_hostname == topology.compute.service_hostname:
        raise LauncherError("The master and compute service hostnames must differ")
    return topology


def quote_inventory_value(value: str | int) -> str:
    text = str(value)
    if not text or any(character.isspace() for character in text):
        raise LauncherError(f"Invalid whitespace in inventory value: {text!r}")
    return text


def kolla_host_line(host: Host, *, master: bool) -> str:
    values: list[tuple[str, str | int]] = [
        ("ansible_host", host.ansible_host),
        ("ansible_user", host.ansible_user),
        ("ansible_port", host.ansible_port),
        ("ansible_become", "true"),
        ("ansible_python_interpreter", host.python),
        ("network_interface", host.management_interface),
        ("api_interface", host.management_interface),
        ("tunnel_interface", host.management_interface),
        ("migration_interface", host.management_interface),
    ]
    if master:
        if host.external_interface is None:
            raise LauncherError("The master external interface is missing")
        values.append(
            ("neutron_external_interface", host.external_interface)
        )
    fields = [quote_inventory_value(host.inventory_name)]
    fields.extend(
        f"{key}={quote_inventory_value(value)}" for key, value in values
    )
    return " ".join(fields)


def inventory_replacements(topology: Topology) -> dict[str, list[str]]:
    master_line = kolla_host_line(topology.master, master=True)
    compute_line = kolla_host_line(topology.compute, master=False)
    return {
        "control": [master_line],
        "network": [master_line],
        "compute": [compute_line],
        "monitoring": [],
        "storage": [master_line],
        "deployment": ["localhost ansible_connection=local"],
    }


def print_command(command: Sequence[str]) -> None:
    print(f"+ {shlex.join(command)}", flush=True)


def run(
    command: Sequence[str | Path],
    *,
    env: Mapping[str, str] | None = None,
    capture: bool = False,
    cwd: Path | None = None,
) -> subprocess.CompletedProcess[str]:
    command_text = [str(value) for value in command]
    print_command(command_text)
    try:
        return subprocess.run(
            command_text,
            check=True,
            env=dict(env) if env is not None else None,
            text=True,
            stdout=subprocess.PIPE if capture else None,
            stderr=None,
            cwd=cwd,
        )
    except (OSError, subprocess.CalledProcessError) as error:
        raise LauncherError(
            f"Command failed: {shlex.join(command_text)}"
        ) from error


def ensure_private_directory(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)
    path.chmod(0o700)


def write_private_text(path: Path, content: str) -> None:
    ensure_private_directory(path.parent)
    descriptor, temporary_name = tempfile.mkstemp(
        prefix=f".{path.name}.",
        dir=path.parent,
        text=True,
    )
    temporary = Path(temporary_name)
    try:
        with os.fdopen(descriptor, "w", encoding="utf-8") as stream:
            stream.write(content)
        temporary.chmod(0o600)
        os.replace(temporary, path)
        path.chmod(0o600)
    finally:
        temporary.unlink(missing_ok=True)


def validate_private_file(path: Path, label: str) -> None:
    try:
        details = path.stat()
    except OSError as error:
        raise LauncherError(f"Cannot read {label}: {path}") from error
    if not stat.S_ISREG(details.st_mode):
        raise LauncherError(f"{label} is not a regular file: {path}")
    mode = stat.S_IMODE(details.st_mode)
    if mode not in {0o400, 0o600}:
        raise LauncherError(
            f"{label} must use mode 0400 or 0600, not {mode:04o}"
        )
    if not os.access(path, os.R_OK):
        raise LauncherError(f"{label} is not readable: {path}")


def credential_path(environment_name: str, label: str) -> Path:
    value = os.environ.get(environment_name, "")
    if not value:
        raise LauncherError(
            f"{environment_name} was not injected by the Bazel wrapper"
        )
    path = Path(value)
    validate_private_file(path, label)
    return path


def state_paths(config: Mapping[str, Any]) -> dict[str, Path]:
    home = Path.home()
    revision = required_string(config, "openstack_kolla_revision")
    release = required_string(config, "openstack_release")
    state = (
        home
        / ".local"
        / "state"
        / "alwaldend"
        / "openstack"
        / f"{release}-{revision[:12]}"
    )
    return {
        "state": state,
        "config": state / "etc-kolla",
        "collections": state / "collections",
        "inventory": state / "multinode",
        "sections": state / "inventory-sections.json",
        "venv": (
            home
            / ".cache"
            / "alwaldend"
            / "openstack"
            / f"kolla-ansible-{revision[:12]}"
        ),
    }


def installed_kolla_revision(venv: Path) -> str:
    python = venv / "bin" / "python"
    if not python.is_file():
        return ""
    script = (
        "import importlib.metadata, json; "
        "dist = importlib.metadata.distribution('kolla-ansible'); "
        "data = json.loads(dist.read_text('direct_url.json') or '{}'); "
        "print(data.get('vcs_info', {}).get('commit_id', ''))"
    )
    try:
        result = run([python, "-c", script], capture=True)
    except LauncherError:
        return ""
    return result.stdout.strip()


def ensure_toolchain(
    config: Mapping[str, Any],
    paths: Mapping[str, Path],
) -> None:
    python = Path(required_string(config, "openstack_kolla_python"))
    if not python.is_file() or not os.access(python, os.X_OK):
        raise LauncherError(f"Local deployment Python is not executable: {python}")
    if shutil.which("git") is None:
        raise LauncherError("Git is required on the Bazel execution host")

    venv = paths["venv"]
    marker = venv / ".openstack-toolchain.json"
    expected = {
        "kolla_revision": config["openstack_kolla_revision"],
        "requirements_revision": config["openstack_requirements_revision"],
    }
    marker_value: Any = None
    if marker.is_file():
        try:
            marker_value = json.loads(marker.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            marker_value = None

    ready = (
        marker_value == expected
        and installed_kolla_revision(venv)
        == config["openstack_kolla_revision"]
        and (venv / "bin" / "kolla-ansible").is_file()
        and (venv / "bin" / "openstack").is_file()
    )
    if ready:
        return

    if venv.exists():
        shutil.rmtree(venv)
    ensure_private_directory(venv.parent)
    run([python, "-m", "venv", venv])
    pip = venv / "bin" / "pip"
    run(
        [
            pip,
            "install",
            "--disable-pip-version-check",
            kolla_source(config),
        ]
    )
    if installed_kolla_revision(venv) != config["openstack_kolla_revision"]:
        raise LauncherError("The installed Kolla revision does not match the pin")
    run(
        [
            pip,
            "install",
            "--disable-pip-version-check",
            "--constraint",
            kolla_constraints(config),
            "python-cinderclient",
            "python-openstackclient",
        ]
    )
    write_private_text(marker, json.dumps(expected, sort_keys=True) + "\n")


def validate_password_document(password_file: Path, venv: Path) -> None:
    script = """
import pathlib
import sys
import yaml

path = pathlib.Path(sys.argv[1])
document = yaml.safe_load(path.read_text(encoding="utf-8"))
assert isinstance(document, dict)
for key in ("database_password", "keystone_admin_password", "rabbitmq_password"):
    assert isinstance(document.get(key), str) and document[key]
key = document.get("kolla_ssh_key")
assert isinstance(key, dict)
assert isinstance(key.get("private_key"), str) and key["private_key"]
assert isinstance(key.get("public_key"), str) and key["public_key"]
""".strip()
    run([venv / "bin" / "python", "-c", script, password_file])


def kolla_environment(
    private_key: Path,
    paths: Mapping[str, Path],
) -> dict[str, str]:
    environment = os.environ.copy()
    environment.update(
        {
            "ANSIBLE_COLLECTIONS_PATH": str(paths["collections"]),
            "ANSIBLE_CONFIG": str(paths["state"] / "ansible.cfg"),
            "ANSIBLE_PRIVATE_KEY_FILE": str(private_key),
            "KOLLA_CONFIG_PATH": str(paths["config"]),
            "PATH": (
                f"{paths['venv'] / 'bin'}{os.pathsep}"
                f"{environment.get('PATH', '')}"
            ),
        }
    )
    return environment


def prepare_state_directories(paths: Mapping[str, Path]) -> None:
    for directory in (
        paths["state"],
        paths["config"],
        paths["collections"],
    ):
        ensure_private_directory(directory)


def write_ansible_config(paths: Mapping[str, Path]) -> None:
    content = f"""[defaults]
collections_path = {paths['collections']}
host_key_checking = True
interpreter_python = auto_silent
retry_files_enabled = False

[ssh_connection]
pipelining = True
"""
    write_private_text(paths["state"] / "ansible.cfg", content)


def ensure_kolla_dependencies(
    config: Mapping[str, Any],
    paths: Mapping[str, Path],
    environment: Mapping[str, str],
) -> None:
    marker = paths["collections"] / (
        f".kolla-deps-{config['openstack_kolla_revision'][:12]}"
    )
    if marker.is_file():
        return
    run(
        [paths["venv"] / "bin" / "kolla-ansible", "install-deps"],
        env=environment,
    )
    write_private_text(marker, "installed\n")


def locate_stock_inventory(venv: Path) -> Path:
    script = (
        "from kolla_ansible import utils; "
        "print(utils.get_data_files_path('ansible', 'inventory', 'multinode'))"
    )
    result = run([venv / "bin" / "python", "-c", script], capture=True)
    path = Path(result.stdout.strip())
    if not path.is_file():
        raise LauncherError(f"Kolla stock inventory is missing: {path}")
    return path


def yaml_string(value: Any) -> str:
    return json.dumps(str(value))


def render_globals(
    config: Mapping[str, Any],
    config_dir: Path,
) -> str:
    internal_vip = yaml_string(
        config["openstack_kolla_internal_vip_address"]
    )
    return f"""---
workaround_ansible_issue_8743: true

config_strategy: COPY_ALWAYS
kolla_base_distro: {yaml_string(config['openstack_kolla_base_distro'])}
openstack_release: {yaml_string(config['openstack_release'])}
kolla_container_engine: {yaml_string(config['openstack_kolla_container_engine'])}
node_custom_config: {yaml_string(config_dir / 'config')}

# Interface variables are supplied per host in the rendered inventory.
kolla_internal_vip_address: {internal_vip}
keepalived_virtual_router_id: {config['openstack_kolla_keepalived_virtual_router_id']}
disable_firewall: "yes"

neutron_plugin_agent: ovn
neutron_ovn_distributed_fip: "no"
enable_neutron_provider_networks: "no"

enable_openstack_core: "yes"
enable_heat: "no"
enable_horizon: "yes"
enable_cinder: "yes"
enable_cinder_backup: "no"
enable_cinder_backend_lvm: "yes"
enable_cinder_backend_iscsi: "yes"
cinder_backend_lvm_name: fast
cinder_volume_group: {yaml_string(config['openstack_fast_vg'])}
cinder_target_helper: lioadm

nova_compute_virt_type: kvm
nova_instance_datadir_volume: {yaml_string(config['openstack_nova_datadir'])}
glance_file_datadir_volume: {yaml_string(config['openstack_glance_data_dir'])}

openstack_logging_debug: "False"
enable_central_logging: "no"
enable_fluentd: "no"
enable_grafana: "no"
enable_prometheus: "no"
"""


def render_cinder_volume(config: Mapping[str, Any]) -> str:
    return f"""[DEFAULT]
enabled_backends = fast,bulk

[bulk]
volume_group = {config['openstack_bulk_vg']}
volume_driver = cinder.volume.drivers.lvm.LVMVolumeDriver
volume_backend_name = bulk
target_helper = lioadm
target_protocol = iscsi
"""


def render_nova_compute(config: Mapping[str, Any]) -> str:
    return f"""[DEFAULT]
instances_path = /var/lib/nova/instances
reserved_host_disk_mb = {config['openstack_nova_reserved_host_disk_mb']}
"""


def render_inputs(
    config: Mapping[str, Any],
    topology: Topology,
    paths: Mapping[str, Path],
) -> None:
    for directory in (
        paths["state"],
        paths["config"],
        paths["config"] / "config",
        paths["config"] / "config" / "cinder",
        paths["config"] / "config" / "nova",
        paths["collections"],
    ):
        ensure_private_directory(directory)

    replacements = inventory_replacements(topology)
    write_private_text(
        paths["sections"],
        json.dumps(replacements, indent=2, sort_keys=True) + "\n",
    )
    stock_inventory = locate_stock_inventory(paths["venv"])
    rendered_inventory = render_kolla_inventory(
        stock_inventory.read_text(encoding="utf-8"),
        replacements,
    )
    write_private_text(paths["inventory"], rendered_inventory)
    write_private_text(
        paths["config"] / "globals.yml",
        render_globals(config, paths["config"]),
    )
    write_private_text(
        paths["config"] / "config" / "cinder.conf",
        "[DEFAULT]\ndefault_volume_type = fast\n",
    )
    write_private_text(
        paths["config"] / "config" / "cinder" / "cinder-volume.conf",
        render_cinder_volume(config),
    )
    write_private_text(
        paths["config"] / "config" / "nova" / "nova-compute.conf",
        render_nova_compute(config),
    )


def prepare_workdir() -> Path:
    runfiles_dir = os.environ.get("RUNFILES_DIR", "")
    if not runfiles_dir:
        raise LauncherError("RUNFILES_DIR is required for host preparation")
    candidates = sorted(
        Path(runfiles_dir).glob(f"*/{PREPARE_ROOT_MARKER}")
    )
    if len(candidates) != 1:
        raise LauncherError(
            "Cannot locate the packaged Ansible working directory"
        )
    return candidates[0].parent


def run_host_preparation(action: str, prepare_bin: Path | None) -> None:
    if action not in ACTIONS_WITH_HOST_PREPARATION:
        return
    if prepare_bin is None:
        raise LauncherError(f"Action {action} requires --prepare-bin")
    if not prepare_bin.is_file() or not os.access(prepare_bin, os.X_OK):
        raise LauncherError(f"Host preparation binary is not executable: {prepare_bin}")
    run([prepare_bin], cwd=prepare_workdir())


def stage_password_file(source: Path, destination: Path) -> None:
    ensure_private_directory(destination.parent)
    destination.unlink(missing_ok=True)
    shutil.copyfile(source, destination)
    destination.chmod(0o600)


def run_kolla_actions(
    action: str,
    password_file: Path,
    paths: Mapping[str, Path],
    environment: Mapping[str, str],
) -> None:
    staged = paths["config"] / "passwords.yml"
    stage_password_file(password_file, staged)
    try:
        for command in KOLLA_ACTIONS[action]:
            run(
                [
                    paths["venv"] / "bin" / "kolla-ansible",
                    command,
                    "--configdir",
                    paths["config"],
                    "--inventory",
                    paths["inventory"],
                    "--become",
                ],
                env=environment,
            )
    finally:
        staged.unlink(missing_ok=True)


def private_clouds_file(paths: Mapping[str, Path]) -> Path:
    clouds = paths["config"] / "clouds.yaml"
    validate_private_file(clouds, "Kolla clouds.yaml")
    return clouds


def openstack_json(
    arguments: Sequence[str],
    *,
    paths: Mapping[str, Path],
    environment: Mapping[str, str],
) -> Any:
    result = run(
        [
            paths["venv"] / "bin" / "openstack",
            "--os-cloud",
            "admin",
            *arguments,
            "--format",
            "json",
        ],
        env=environment,
        capture=True,
    )
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError as error:
        raise LauncherError("OpenStack CLI returned invalid JSON") from error


def mapping_value(document: Mapping[str, Any], key: str) -> Any:
    for candidate, value in document.items():
        if candidate.casefold() == key.casefold():
            return value
    return None


def reconcile_volume_types(
    paths: Mapping[str, Path],
    environment: Mapping[str, str],
) -> None:
    rows = openstack_json(
        ["volume", "type", "list"],
        paths=paths,
        environment=environment,
    )
    if not isinstance(rows, list):
        raise LauncherError("Volume type list must be a JSON list")
    existing = {
        mapping_value(row, "Name")
        for row in rows
        if isinstance(row, dict)
    }
    openstack = paths["venv"] / "bin" / "openstack"
    for name in ("fast", "bulk"):
        if name not in existing:
            run(
                [
                    openstack,
                    "--os-cloud",
                    "admin",
                    "volume",
                    "type",
                    "create",
                    "--property",
                    f"volume_backend_name={name}",
                    name,
                ],
                env=environment,
            )
        run(
            [
                openstack,
                "--os-cloud",
                "admin",
                "volume",
                "type",
                "set",
                "--property",
                f"volume_backend_name={name}",
                name,
            ],
            env=environment,
        )
        document = openstack_json(
            ["volume", "type", "show", name],
            paths=paths,
            environment=environment,
        )
        if not isinstance(document, dict):
            raise LauncherError(f"Volume type {name} must be a JSON object")
        properties = mapping_value(document, "properties")
        if isinstance(properties, str):
            valid = f"volume_backend_name='{name}'" in properties
        elif isinstance(properties, dict):
            valid = properties.get("volume_backend_name") == name
        else:
            valid = False
        if not valid:
            raise LauncherError(f"Volume type {name} uses the wrong backend")


def validate_services(
    topology: Topology,
    paths: Mapping[str, Path],
    environment: Mapping[str, str],
) -> None:
    cinder = openstack_json(
        [
            "volume",
            "service",
            "list",
            "--service",
            "cinder-volume",
        ],
        paths=paths,
        environment=environment,
    )
    if not isinstance(cinder, list):
        raise LauncherError("Cinder service list must be a JSON list")
    active_cinder = {
        mapping_value(row, "Host")
        for row in cinder
        if isinstance(row, dict)
        and mapping_value(row, "Binary") == "cinder-volume"
        and mapping_value(row, "Status") == "enabled"
        and mapping_value(row, "State") == "up"
    }
    expected_cinder = {
        f"{topology.master.service_hostname}@fast",
        f"{topology.master.service_hostname}@bulk",
    }
    if active_cinder != expected_cinder:
        raise LauncherError(
            f"Unexpected Cinder services: {sorted(active_cinder)}"
        )

    nova = openstack_json(
        ["compute", "service", "list", "--service", "nova-compute"],
        paths=paths,
        environment=environment,
    )
    if not isinstance(nova, list):
        raise LauncherError("Nova service list must be a JSON list")
    active_nova = [
        row
        for row in nova
        if isinstance(row, dict)
        and mapping_value(row, "Binary") == "nova-compute"
        and mapping_value(row, "Status") == "enabled"
        and mapping_value(row, "State") == "up"
    ]
    if len(active_nova) != 1:
        raise LauncherError(f"Unexpected Nova services: {active_nova}")
    if (
        mapping_value(active_nova[0], "Host")
        != topology.compute.service_hostname
    ):
        raise LauncherError("Nova compute runs on the wrong host")


def main(argv: Sequence[str] | None = None) -> int:
    args = parse_args(argv)
    config = read_json_object(args.config)
    inventory = read_json_object(args.inventory)
    validate_config(config)
    topology = parse_topology(inventory)

    password_file = credential_path(
        "KOLLA_PASSWORDS_FILE", "Kolla password document"
    )
    private_key = credential_path(
        "ANSIBLE_PRIVATE_KEY_FILE", "Ansible private key"
    )
    paths = state_paths(config)
    ensure_toolchain(config, paths)
    validate_password_document(password_file, paths["venv"])
    prepare_state_directories(paths)
    write_ansible_config(paths)
    environment = kolla_environment(private_key, paths)
    ensure_kolla_dependencies(config, paths, environment)

    run_host_preparation(args.action, args.prepare_bin)
    render_inputs(config, topology, paths)
    run_kolla_actions(args.action, password_file, paths, environment)

    if args.action in ACTIONS_WITH_SERVICE_CHECKS:
        clouds = private_clouds_file(paths)
        client_environment = dict(environment)
        client_environment["OS_CLIENT_CONFIG_FILE"] = str(clouds)
        reconcile_volume_types(paths, client_environment)
        validate_services(topology, paths, client_environment)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except LauncherError as error:
        print(f"error: {error}", file=sys.stderr)
        raise SystemExit(2) from error
