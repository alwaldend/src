"""Tests for Blender-free render-packet specification validation."""

import json
import pathlib
import tempfile
import typing
import unittest

from projects.renders.cmd.fumo_review.render_spec import (
    SpecError,
    inverse_srgb,
    load_spec,
    parse_spec,
    resolve_output_path,
)


def _valid_spec() -> typing.Dict[str, typing.Any]:
    return {
        "schema_version": 1,
        "scene": "Review Scene",
        "resolution": {"width": 320, "height": 240},
        "beauty_views": {
            "front": {"camera": "Review Front", "output": "front.png"},
            "side": {"camera": "Review Side", "output": "side.png"},
        },
        "component_id_passes": {
            "front_ids": {
                "camera": "Review Front",
                "output": "masks/front_ids.png",
                "objects": {
                    "Body": [255, 0, 0],
                    "Bow": [0, 255, 0],
                },
                "fallback": {"mode": "color", "rgb": [0, 0, 0]},
            },
            "side_ids": {
                "camera": "Review Side",
                "output": "masks/side_ids.png",
                "objects": {"Body": [255, 0, 0]},
                "fallback": {"mode": "hide"},
            },
        },
    }


class RenderSpecTest(unittest.TestCase):
    """Exercise schema, color, and output-path validation."""

    def test_valid_spec_is_sorted_and_typed(self) -> None:
        spec = parse_spec(_valid_spec())
        self.assertEqual(spec.scene, "Review Scene")
        self.assertEqual(
            (spec.resolution.width, spec.resolution.height), (320, 240)
        )
        self.assertEqual(list(spec.beauty_views), ["front", "side"])
        self.assertEqual(
            spec.component_id_passes["front_ids"].objects["Bow"], (0, 255, 0)
        )
        self.assertEqual(
            spec.output_names(),
            [
                "front.png",
                "side.png",
                "masks/front_ids.png",
                "masks/side_ids.png",
            ],
        )

    def test_rejects_unknown_schema_and_keys(self) -> None:
        value = _valid_spec()
        self.assertIsInstance(value, dict)
        value["schema_version"] = 2
        with self.assertRaisesRegex(SpecError, "integer 1"):
            parse_spec(value)
        value["schema_version"] = 1
        value["typo"] = True
        with self.assertRaisesRegex(SpecError, "unknown keys"):
            parse_spec(value)

    def test_rejects_invalid_colors(self) -> None:
        value = _valid_spec()
        self.assertIsInstance(value, dict)
        passes = value["component_id_passes"]
        self.assertIsInstance(passes, dict)
        front = passes["front_ids"]
        self.assertIsInstance(front, dict)
        objects = front["objects"]
        self.assertIsInstance(objects, dict)
        objects["Body"] = [256, 0, 0]
        with self.assertRaisesRegex(SpecError, "0 to 255"):
            parse_spec(value)

    def test_rejects_implicit_or_malformed_fallback(self) -> None:
        value = _valid_spec()
        self.assertIsInstance(value, dict)
        passes = value["component_id_passes"]
        self.assertIsInstance(passes, dict)
        side = passes["side_ids"]
        self.assertIsInstance(side, dict)
        side["fallback"] = {"mode": "color"}
        with self.assertRaisesRegex(SpecError, "missing keys"):
            parse_spec(value)

    def test_rejects_unsafe_and_duplicate_outputs(self) -> None:
        value = _valid_spec()
        self.assertIsInstance(value, dict)
        views = value["beauty_views"]
        self.assertIsInstance(views, dict)
        front = views["front"]
        self.assertIsInstance(front, dict)
        front["output"] = "../front.png"
        with self.assertRaisesRegex(SpecError, "beneath output"):
            parse_spec(value)
        front["output"] = "masks/../front.png"
        with self.assertRaisesRegex(SpecError, "beneath output"):
            parse_spec(value)
        front["output"] = "masks//front_ids.png"
        with self.assertRaisesRegex(SpecError, "duplicate output"):
            parse_spec(value)

    def test_load_rejects_duplicate_json_keys(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = pathlib.Path(directory, "spec.json")
            path.write_text(
                '{"schema_version": 1, "schema_version": 1}',
                encoding="utf-8",
            )
            with self.assertRaisesRegex(SpecError, "duplicate JSON key"):
                load_spec(path)

    def test_inverse_srgb_known_values(self) -> None:
        converted = inverse_srgb((0, 128, 255))
        self.assertEqual(converted[0], 0.0)
        self.assertAlmostEqual(converted[1], 0.21586050011389926)
        self.assertEqual(converted[2], 1.0)

    def test_resolve_output_rejects_symlink_escape(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = pathlib.Path(directory, "packet")
            outside = pathlib.Path(directory, "outside")
            root.mkdir()
            outside.mkdir()
            pathlib.Path(root, "linked").symlink_to(
                outside, target_is_directory=True
            )
            with self.assertRaisesRegex(SpecError, "outside output"):
                resolve_output_path(root, "linked/view.png")

    def test_json_round_trip_fixture(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = pathlib.Path(directory, "spec.json")
            path.write_text(json.dumps(_valid_spec()), encoding="utf-8")
            self.assertEqual(load_spec(path).scene, "Review Scene")


if __name__ == "__main__":
    unittest.main()
