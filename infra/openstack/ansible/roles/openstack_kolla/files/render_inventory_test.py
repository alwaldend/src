#!/usr/bin/env python3
"""Tests for the Kolla inventory renderer."""

from __future__ import annotations

import importlib.util
import unittest
from pathlib import Path

MODULE_PATH = Path(__file__).with_name("render_inventory.py")
SPEC = importlib.util.spec_from_file_location("render_inventory", MODULE_PATH)
assert SPEC is not None and SPEC.loader is not None
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)

render = MODULE.render

SOURCE = """# stock inventory
[control]
control01

[network]
network01

[compute]
compute01

[monitoring]
monitoring01

[storage]
storage01

[deployment]
localhost ansible_connection=local

[baremetal:children]
control
network
compute
storage
monitoring
"""

REPLACEMENTS = {
    "control": ["master ansible_host=192.0.2.10"],
    "network": ["master ansible_host=192.0.2.10"],
    "compute": ["compute ansible_host=192.0.2.11"],
    "monitoring": [],
    "storage": ["master ansible_host=192.0.2.10"],
    "deployment": ["localhost ansible_connection=local"],
}


class RenderInventoryTest(unittest.TestCase):
    def test_replaces_only_host_sections(self) -> None:
        rendered = render(SOURCE, REPLACEMENTS)
        self.assertIn("[baremetal:children]\ncontrol\nnetwork", rendered)
        self.assertNotIn("control01", rendered)
        self.assertNotIn("monitoring01", rendered)
        self.assertEqual(rendered.count("master ansible_host=192.0.2.10"), 3)

    def test_is_idempotent(self) -> None:
        rendered = render(SOURCE, REPLACEMENTS)
        self.assertEqual(render(rendered, REPLACEMENTS), rendered)

    def test_rejects_missing_source_section(self) -> None:
        source = SOURCE.replace("[storage]", "[storage-missing]")
        with self.assertRaisesRegex(ValueError, "source inventory is missing"):
            render(source, REPLACEMENTS)

    def test_rejects_duplicate_source_section(self) -> None:
        source = SOURCE + "\n[compute]\ncompute02\n"
        with self.assertRaisesRegex(ValueError, "duplicate sections"):
            render(source, REPLACEMENTS)

    def test_rejects_extra_replacement(self) -> None:
        replacements = dict(REPLACEMENTS, unexpected=[])
        with self.assertRaisesRegex(ValueError, "unexpected replacement"):
            render(SOURCE, replacements)

    def test_rejects_multiline_host_entry(self) -> None:
        replacements = dict(REPLACEMENTS)
        replacements["compute"] = ["compute\nsecond-line"]
        with self.assertRaisesRegex(ValueError, "multiline"):
            render(SOURCE, replacements)


if __name__ == "__main__":
    unittest.main()
