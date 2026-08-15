from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import launcher


class LauncherTest(unittest.TestCase):
    def setUp(self) -> None:
        self.config = {
            "openstack_release": "2026.1",
            "openstack_kolla_revision": "a" * 40,
            "openstack_requirements_revision": "b" * 40,
            "openstack_kolla_base_distro": "rocky",
            "openstack_kolla_container_engine": "docker",
            "openstack_kolla_python": "/usr/bin/python3",
            "openstack_kolla_internal_vip_address": "192.0.2.10",
            "openstack_kolla_keepalived_virtual_router_id": 53,
            "openstack_manage_storage": False,
            "openstack_storage_wipe": False,
            "openstack_fast_raid_device": "/dev/md/cinder-fast",
            "openstack_bulk_raid_device": "/dev/md/cinder-bulk",
            "openstack_fast_vg": "cinder-fast",
            "openstack_bulk_vg": "cinder-bulk",
            "openstack_glance_data_dir": "/var/lib/openstack/glance",
            "openstack_manage_compute_storage": False,
            "openstack_compute_storage_wipe": False,
            "openstack_nova_datadir": "/var/lib/nova",
            "openstack_compute_filesystem": "xfs",
            "openstack_compute_mount_options": "defaults,noatime",
            "openstack_compute_minimum_bytes": 1_800_000_000_000,
            "openstack_nova_reserved_host_disk_mb": 51_200,
            "openstack_supported_actions": list(launcher.SUPPORTED_ACTIONS),
        }
        self.inventory = {
            "all": {
                "children": {
                    "openstack": {
                        "children": {
                            "openstack_master": {
                                "hosts": {
                                    "master.example.com": {
                                        "ansible_user": "ansible",
                                        "ansible_become": True,
                                        "ansible_python_interpreter": (
                                            "/usr/bin/python3"
                                        ),
                                        "openstack_service_hostname": "master",
                                        "openstack_management_interface": "eno1",
                                        "openstack_external_interface": "eno1.20",
                                    }
                                }
                            },
                            "openstack_compute": {
                                "hosts": {
                                    "compute.example.com": {
                                        "ansible_host": "192.0.2.22",
                                        "ansible_user": "ansible",
                                        "ansible_become": True,
                                        "ansible_python_interpreter": (
                                            "/usr/bin/python3"
                                        ),
                                        "openstack_service_hostname": "compute",
                                        "openstack_management_interface": "eno1",
                                    }
                                }
                            },
                        }
                    }
                }
            }
        }

    def test_validate_config_accepts_pinned_inputs(self) -> None:
        launcher.validate_config(self.config)

    def test_validate_config_rejects_short_revision(self) -> None:
        self.config["openstack_kolla_revision"] = "short"
        with self.assertRaisesRegex(
            launcher.LauncherError,
            "full Git SHAs",
        ):
            launcher.validate_config(self.config)

    def test_source_urls_use_pinned_revisions(self) -> None:
        self.assertTrue(
            launcher.kolla_source(self.config).endswith("@" + "a" * 40)
        )
        self.assertIn(
            "b" * 40,
            launcher.kolla_constraints(self.config),
        )

    def test_parse_topology_builds_kolla_host_lines(self) -> None:
        topology = launcher.parse_topology(self.inventory)
        replacements = launcher.inventory_replacements(topology)
        master = replacements["control"][0]
        compute = replacements["compute"][0]
        self.assertIn("master.example.com", master)
        self.assertIn("neutron_external_interface=eno1.20", master)
        self.assertIn("ansible_host=192.0.2.22", compute)
        self.assertNotIn("neutron_external_interface", compute)
        self.assertEqual(replacements["storage"], [master])

    def test_parse_topology_rejects_duplicate_service_hostnames(self) -> None:
        compute = self.inventory["all"]["children"]["openstack"][
            "children"
        ]["openstack_compute"]["hosts"]["compute.example.com"]
        compute["openstack_service_hostname"] = "master"
        with self.assertRaisesRegex(
            launcher.LauncherError,
            "service hostnames must differ",
        ):
            launcher.parse_topology(self.inventory)

    def test_render_globals_uses_compute_datadir(self) -> None:
        result = launcher.render_globals(self.config, Path("/tmp/kolla"))
        self.assertIn('nova_instance_datadir_volume: "/var/lib/nova"', result)
        self.assertIn("keepalived_virtual_router_id: 53", result)
        self.assertIn('node_custom_config: "/tmp/kolla/config"', result)

    def test_write_private_text_replaces_content_and_mode(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "nested" / "value.json"
            launcher.write_private_text(path, json.dumps({"value": 1}))
            self.assertEqual(json.loads(path.read_text()), {"value": 1})
            self.assertEqual(path.stat().st_mode & 0o777, 0o600)
            self.assertEqual(path.parent.stat().st_mode & 0o777, 0o700)

    def test_prepare_workdir_uses_packaged_marker(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory) / "workspace"
            root.mkdir()
            (root / launcher.PREPARE_ROOT_MARKER).write_text("marker\n")
            with mock.patch.dict(
                "os.environ", {"RUNFILES_DIR": directory}, clear=False
            ):
                self.assertEqual(launcher.prepare_workdir(), root)

    def test_action_sequences_do_not_nest_ansible(self) -> None:
        self.assertEqual(
            launcher.KOLLA_ACTIONS["prechecks"],
            ("bootstrap-servers", "prechecks"),
        )
        self.assertNotIn("ansible-playbook", launcher.KOLLA_ACTIONS["deploy"])


if __name__ == "__main__":
    unittest.main()
