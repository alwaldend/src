from __future__ import annotations

import unittest

from render_inventory import render, source_sections, validate


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
    def test_render_replaces_only_host_sections(self) -> None:
        result = render(SOURCE, REPLACEMENTS)
        self.assertIn("master ansible_host=192.0.2.10", result)
        self.assertIn("compute ansible_host=192.0.2.11", result)
        self.assertNotIn("control01", result)
        self.assertIn("[baremetal:children]\ncontrol\nnetwork", result)

    def test_render_is_idempotent(self) -> None:
        first = render(SOURCE, REPLACEMENTS)
        second = render(first, REPLACEMENTS)
        self.assertEqual(first, second)

    def test_duplicate_source_section_is_rejected(self) -> None:
        with self.assertRaisesRegex(ValueError, "duplicate sections"):
            source_sections(SOURCE + "\n[control]\nother\n")

    def test_missing_replacement_is_rejected(self) -> None:
        replacements = dict(REPLACEMENTS)
        replacements.pop("storage")
        with self.assertRaisesRegex(ValueError, "missing replacement"):
            validate(replacements)

    def test_extra_replacement_is_rejected(self) -> None:
        replacements = dict(REPLACEMENTS)
        replacements["extra"] = ["host"]
        with self.assertRaisesRegex(ValueError, "unexpected replacement"):
            validate(replacements)

    def test_multiline_host_is_rejected(self) -> None:
        replacements = dict(REPLACEMENTS)
        replacements["control"] = ["master\nother"]
        with self.assertRaisesRegex(ValueError, "multiline"):
            validate(replacements)


if __name__ == "__main__":
    unittest.main()
