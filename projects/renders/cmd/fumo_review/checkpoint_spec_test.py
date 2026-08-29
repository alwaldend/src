"""Tests for Blender-free live-checkpoint capsule validation."""

import hashlib
import json
import math
import pathlib
import tempfile
import typing
import unittest

from projects.renders.cmd.fumo_review.checkpoint_spec import (
    CheckpointSpecError,
    canonical_json,
    load_checkpoint_spec,
    load_checkpoint_spec_with_hash,
    parse_checkpoint_spec,
    resolve_out_directory,
    resolve_workspace_artifact,
)


def _hash(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _valid_spec() -> typing.Dict[str, typing.Any]:
    return {
        "schema_version": 1,
        "checkpoint": "hair_stage3c_02",
        "scene": "Hair scratch",
        "resolution": {"width": 320, "height": 320},
        "parent": {"path": "out/parent.blend", "sha256": "a" * 64},
        "source_packet": {
            "path": "out/source/README.md",
            "sha256": "b" * 64,
        },
        "hypothesis": "The traced panel stays outside the support.",
        "views": {
            "rear": {"camera": "Review Rear"},
            "front": {"camera": "Review Front"},
        },
        "vetoes": ["exposed support", "symmetric fringe"],
        "scene_boolean_requirements": {
            "stage3d_clearance_gate_pass": True,
        },
    }


class CheckpointSpecTest(unittest.TestCase):
    """Exercise capsule, path, hash, and canonical-output validation."""

    def test_valid_spec_is_sorted_and_typed(self) -> None:
        spec = parse_checkpoint_spec(_valid_spec())
        self.assertEqual(spec.checkpoint, "hair_stage3c_02")
        self.assertEqual(
            (spec.resolution.width, spec.resolution.height), (320, 320)
        )
        self.assertEqual(list(spec.views), ["front", "rear"])
        self.assertEqual(spec.views["front"].camera, "Review Front")
        self.assertEqual(
            spec.vetoes,
            ("exposed support", "symmetric fringe"),
        )
        self.assertEqual(
            spec.scene_boolean_requirements,
            {"stage3d_clearance_gate_pass": True},
        )

    def test_fast_checkpoint_limits_resolution_and_views(self) -> None:
        value = _valid_spec()
        value["resolution"]["width"] = 641
        with self.assertRaisesRegex(CheckpointSpecError, "16 to 640"):
            parse_checkpoint_spec(value)
        value = _valid_spec()
        value["views"]["side"] = {"camera": "Review Side"}
        with self.assertRaisesRegex(CheckpointSpecError, "at most 2 views"):
            parse_checkpoint_spec(value)

    def test_source_timestamp_must_be_rfc3339(self) -> None:
        value = _valid_spec()
        value["t0_source_packet_frozen_utc"] = "2026-08-29 12:00:00"
        with self.assertRaisesRegex(CheckpointSpecError, "RFC 3339"):
            parse_checkpoint_spec(value)
        value["t0_source_packet_frozen_utc"] = "2026-08-29T12:00:00Z"
        self.assertEqual(
            parse_checkpoint_spec(value).t0_source_packet_frozen_utc,
            "2026-08-29T12:00:00Z",
        )

    def test_rejects_unknown_keys_and_unsafe_identifiers(self) -> None:
        value = _valid_spec()
        value["typo"] = True
        with self.assertRaisesRegex(CheckpointSpecError, "unknown keys"):
            parse_checkpoint_spec(value)
        del value["typo"]
        value["checkpoint"] = "hair-stage"
        with self.assertRaisesRegex(CheckpointSpecError, "ASCII letters"):
            parse_checkpoint_spec(value)
        value["checkpoint"] = "hair_stage "
        with self.assertRaisesRegex(CheckpointSpecError, "whitespace"):
            parse_checkpoint_spec(value)
        value["checkpoint"] = "hair_stage"
        value["hypothesis"] = "x" * 16385
        with self.assertRaisesRegex(CheckpointSpecError, "unreasonably long"):
            parse_checkpoint_spec(value)

    def test_rejects_duplicate_view_cameras(self) -> None:
        value = _valid_spec()
        value["views"]["rear"] = {"camera": "Review Front"}
        with self.assertRaisesRegex(CheckpointSpecError, "distinct cameras"):
            parse_checkpoint_spec(value)
        value = _valid_spec()
        value["views"]["front "] = value["views"].pop("front")
        with self.assertRaisesRegex(CheckpointSpecError, "whitespace"):
            parse_checkpoint_spec(value)

    def test_rejects_bad_artifacts_and_vetoes(self) -> None:
        value = _valid_spec()
        value["parent"]["path"] = "../parent.blend"
        with self.assertRaisesRegex(
            CheckpointSpecError, "inside the repository"
        ):
            parse_checkpoint_spec(value)
        value["parent"]["path"] = "out/parent.blend"
        value["parent"]["sha256"] = "ABC"
        with self.assertRaisesRegex(CheckpointSpecError, "SHA-256"):
            parse_checkpoint_spec(value)
        value["parent"]["sha256"] = "a" * 64
        value["vetoes"] = ["same", "same"]
        with self.assertRaisesRegex(CheckpointSpecError, "duplicates"):
            parse_checkpoint_spec(value)

    def test_rejects_invalid_scene_boolean_requirements(self) -> None:
        value = _valid_spec()
        value["scene_boolean_requirements"] = {"clearance-pass": True}
        with self.assertRaisesRegex(CheckpointSpecError, "ASCII letters"):
            parse_checkpoint_spec(value)
        value["scene_boolean_requirements"] = {"clearance_pass": 1}
        with self.assertRaisesRegex(CheckpointSpecError, "must be a boolean"):
            parse_checkpoint_spec(value)

    def test_load_rejects_duplicate_json_keys(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = pathlib.Path(directory, "spec.json")
            path.write_text(
                '{"schema_version": 1, "schema_version": 1}',
                encoding="utf-8",
            )
            with self.assertRaisesRegex(CheckpointSpecError, "duplicate JSON"):
                load_checkpoint_spec(path)

    def test_load_hashes_the_exact_parsed_bytes(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = pathlib.Path(directory, "spec.json")
            payload = (json.dumps(_valid_spec(), indent=2) + "\n").encode()
            path.write_bytes(payload)
            spec, digest = load_checkpoint_spec_with_hash(path)
            self.assertEqual(spec.checkpoint, "hair_stage3c_02")
            self.assertEqual(digest, _hash(payload))

    def test_resolve_artifact_checks_containment_and_hash(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            workspace = pathlib.Path(directory, "workspace")
            workspace.mkdir()
            data = b"candidate"
            artifact_path = workspace / "parent.blend"
            artifact_path.write_bytes(data)
            value = _valid_spec()
            value["parent"] = {
                "path": "parent.blend",
                "sha256": _hash(data),
            }
            spec = parse_checkpoint_spec(value)
            self.assertEqual(
                resolve_workspace_artifact(workspace, spec.parent),
                artifact_path,
            )
            value["parent"]["sha256"] = "0" * 64
            spec = parse_checkpoint_spec(value)
            with self.assertRaisesRegex(CheckpointSpecError, "hash mismatch"):
                resolve_workspace_artifact(workspace, spec.parent)

    def test_resolve_artifact_rejects_symlink_escape(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = pathlib.Path(directory)
            workspace = root / "workspace"
            outside = root / "outside"
            workspace.mkdir()
            outside.mkdir()
            target = outside / "source.md"
            target.write_bytes(b"source")
            (workspace / "linked").symlink_to(
                outside, target_is_directory=True
            )
            value = _valid_spec()
            value["source_packet"] = {
                "path": "linked/source.md",
                "sha256": _hash(b"source"),
            }
            spec = parse_checkpoint_spec(value)
            with self.assertRaisesRegex(CheckpointSpecError, "symbolic-link"):
                resolve_workspace_artifact(workspace, spec.source_packet)

    def test_resolve_artifact_rejects_internal_symlink(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            workspace = pathlib.Path(directory, "workspace")
            workspace.mkdir()
            target = workspace / "source.md"
            target.write_bytes(b"source")
            (workspace / "alias.md").symlink_to(target)
            value = _valid_spec()
            value["source_packet"] = {
                "path": "alias.md",
                "sha256": _hash(b"source"),
            }
            spec = parse_checkpoint_spec(value)
            with self.assertRaisesRegex(CheckpointSpecError, "symbolic-link"):
                resolve_workspace_artifact(workspace, spec.source_packet)

    def test_output_must_be_strictly_below_repository_out(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            workspace = pathlib.Path(directory, "workspace")
            (workspace / "out").mkdir(parents=True)
            expected = workspace / "out" / "attempt" / "checkpoint"
            self.assertEqual(
                resolve_out_directory(
                    workspace,
                    pathlib.Path("out/attempt/checkpoint"),
                ),
                expected,
            )
            with self.assertRaisesRegex(CheckpointSpecError, "must be below"):
                resolve_out_directory(workspace, pathlib.Path("out"))
            with self.assertRaisesRegex(CheckpointSpecError, "outside"):
                resolve_out_directory(workspace, pathlib.Path("elsewhere"))
            (workspace / "out" / "real").mkdir()
            (workspace / "out" / "alias").symlink_to(
                workspace / "out" / "real",
                target_is_directory=True,
            )
            with self.assertRaisesRegex(CheckpointSpecError, "symbolic-link"):
                resolve_out_directory(
                    workspace,
                    pathlib.Path("out/alias/checkpoint"),
                )

    def test_canonical_json_is_order_independent(self) -> None:
        first = canonical_json({"b": 2, "a": {"d": 4, "c": 3}})
        second = canonical_json({"a": {"c": 3, "d": 4}, "b": 2})
        self.assertEqual(first, second)
        self.assertEqual(json.loads(first), {"a": {"c": 3, "d": 4}, "b": 2})
        with self.assertRaises(ValueError):
            canonical_json({"bad": math.nan})


if __name__ == "__main__":
    unittest.main()
