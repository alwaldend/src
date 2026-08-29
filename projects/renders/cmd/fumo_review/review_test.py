"""Hermetic tests for the Fumo render-packet review tool."""

import hashlib
import json
import pathlib
import struct
import tempfile
import typing
import unittest
import zlib
from unittest import mock

from projects.renders.cmd.fumo_review import png as png_module
from projects.renders.cmd.fumo_review import review as review_module
from projects.renders.cmd.fumo_review.png import (
    MAX_PNG_DIMENSION,
    PNGError,
    decode_png,
    decode_png_bytes,
)
from projects.renders.cmd.fumo_review.review import (
    ConfigError,
    audit_config,
    write_outputs,
)


def _chunk(kind: bytes, payload: bytes) -> bytes:
    checksum = zlib.crc32(kind)
    checksum = zlib.crc32(payload, checksum) & 0xFFFFFFFF
    return (
        struct.pack(">I", len(payload))
        + kind
        + payload
        + struct.pack(">I", checksum)
    )


def _paeth(left: int, above: int, upper_left: int) -> int:
    estimate = left + above - upper_left
    candidates = (
        (abs(estimate - left), left),
        (abs(estimate - above), above),
        (abs(estimate - upper_left), upper_left),
    )
    return min(candidates, key=lambda item: item[0])[1]


def _filter_row(raw: bytes, previous: bytes, bpp: int, kind: int) -> bytes:
    encoded = bytearray()
    for index, value in enumerate(raw):
        left = raw[index - bpp] if index >= bpp else 0
        above = previous[index] if previous else 0
        upper_left = previous[index - bpp] if previous and index >= bpp else 0
        if kind == 0:
            predictor = 0
        elif kind == 1:
            predictor = left
        elif kind == 2:
            predictor = above
        elif kind == 3:
            predictor = (left + above) // 2
        else:
            predictor = _paeth(left, above, upper_left)
        encoded.append((value - predictor) & 0xFF)
    return bytes(encoded)


def _write_png(
    path: pathlib.Path,
    width: int,
    rows: list[bytes],
    color_type: int,
) -> None:
    channels = {0: 1, 2: 3, 4: 2, 6: 4}[color_type]
    encoded = bytearray()
    previous = b""
    for index, row in enumerate(rows):
        filter_type = index % 5
        encoded.append(filter_type)
        encoded.extend(_filter_row(row, previous, channels, filter_type))
        previous = row
    header = struct.pack(">IIBBBBB", width, len(rows), 8, color_type, 0, 0, 0)
    data = (
        b"\x89PNG\r\n\x1a\n"
        + _chunk(b"IHDR", header)
        + _chunk(b"IDAT", zlib.compress(bytes(encoded)))
        + _chunk(b"IEND", b"")
    )
    path.write_bytes(data)


def _write_bad_filter_png(path: pathlib.Path, width: int, height: int) -> None:
    header = struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0)
    encoded = b"".join(b"\x05" + bytes(width * 3) for _ in range(height))
    data = (
        b"\x89PNG\r\n\x1a\n"
        + _chunk(b"IHDR", header)
        + _chunk(b"IDAT", zlib.compress(encoded))
        + _chunk(b"IEND", b"")
    )
    path.write_bytes(data)


def _rgb_rows(width: int, height: int) -> list[bytes]:
    rows = []
    for y in range(height):
        row = bytearray()
        for x in range(width):
            row.extend(((x * 40 + y) % 256, (x + y * 30) % 256, 200 - y))
        rows.append(bytes(row))
    return rows


def _convert_rows(rows: list[bytes], color_type: int) -> list[bytes]:
    converted = []
    for row in rows:
        target = bytearray()
        for index in range(0, len(row), 3):
            pixel_end = index + 3
            red, green, blue = row[index:pixel_end]
            gray = red
            if color_type == 0:
                target.append(gray)
            elif color_type == 2:
                target.extend((red, green, blue))
            elif color_type == 4:
                target.extend((gray, 173))
            elif color_type == 6:
                target.extend((red, green, blue, 173))
        converted.append(bytes(target))
    return converted


class PNGTest(unittest.TestCase):
    def test_decodes_supported_color_types_and_all_filters(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            rgb_rows = _rgb_rows(3, 5)
            for color_type in (0, 2, 4, 6):
                path = directory / f"type_{color_type}.png"
                rows = _convert_rows(rgb_rows, color_type)
                _write_png(path, 3, rows, color_type)
                image = decode_png(path)
                self.assertEqual((image.width, image.height), (3, 5))
                if color_type in (0, 4):
                    self.assertEqual(image.pixel(1, 2), (42, 42, 42))
                else:
                    self.assertEqual(image.pixel(1, 2), (42, 61, 198))

    def test_rejects_a_bad_crc(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            path = pathlib.Path(temporary) / "broken.png"
            _write_png(path, 1, [b"\x01\x02\x03"], 2)
            data = bytearray(path.read_bytes())
            data[-1] ^= 0xFF
            path.write_bytes(data)
            with self.assertRaisesRegex(PNGError, "CRC mismatch"):
                decode_png(path)

    def test_rejects_huge_dimensions_before_scanline_allocation(self) -> None:
        header = struct.pack(
            ">IIBBBBB",
            MAX_PNG_DIMENSION + 1,
            1,
            8,
            2,
            0,
            0,
            0,
        )
        data = (
            b"\x89PNG\r\n\x1a\n"
            + _chunk(b"IHDR", header)
            + _chunk(b"IEND", b"")
        )
        with self.assertRaisesRegex(PNGError, "dimensions exceed limit"):
            decode_png_bytes(data)

    def test_rejects_compressed_bomb_with_capped_decompression(self) -> None:
        header = struct.pack(">IIBBBBB", 1, 1, 8, 2, 0, 0, 0)
        data = (
            b"\x89PNG\r\n\x1a\n"
            + _chunk(b"IHDR", header)
            + _chunk(b"IDAT", zlib.compress(b"\x00" * (1024 * 1024)))
            + _chunk(b"IEND", b"")
        )
        with self.assertRaisesRegex(PNGError, "exceeds expected"):
            decode_png_bytes(data)


class AuditTest(unittest.TestCase):
    def _packet(
        self, directory: pathlib.Path, rear_width: int = 4
    ) -> pathlib.Path:
        render_rows = [bytes([40, 50, 60] * 4) for _ in range(3)]
        rear_rows = [bytes([40, 50, 60] * rear_width) for _ in range(3)]
        mask_rows = [
            bytes([0, 0, 0] * 4),
            bytes([255, 0, 0, 255, 0, 0, 0, 0, 0, 255, 0, 0]),
            bytes([0, 0, 0] * 4),
        ]
        _write_png(directory / "front.png", 4, render_rows, 2)
        _write_png(directory / "rear.png", rear_width, rear_rows, 2)
        _write_png(directory / "front_ids.png", 4, mask_rows, 2)
        (directory / "model.blend").write_bytes(b"test placeholder")
        config = {
            "schema_version": 1,
            "title": "Fixture review",
            "required_views": ["front", "rear"],
            "views": {
                "front": {"image": "front.png", "mask": "front_ids.png"},
                "rear": {"image": "rear.png"},
            },
            "required_files": ["blend"],
            "files": {"blend": "model.blend"},
            "components": [
                {
                    "name": "red panels",
                    "view": "front",
                    "rgb": [255, 0, 0],
                    "thresholds": {
                        "pixel_count": {"min": 3, "max": 3},
                        "pixel_fraction": {"min": 0.25, "max": 0.25},
                        "bbox_width": {"min": 1, "max": 1},
                        "connected_components": {"min": 2, "max": 2},
                        "max_horizontal_run": {"min": 0.5, "max": 0.5},
                    },
                }
            ],
        }
        config_path = directory / "review.json"
        config_path.write_text(json.dumps(config), encoding="utf-8")
        return config_path

    def test_audits_packet_and_writes_review_artifacts(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory))
            self.assertTrue(result["passed"])
            evidence = typing.cast(
                typing.List[typing.Dict[str, object]], result["evidence"]
            )
            for record in evidence:
                data = pathlib.Path(
                    typing.cast(str, record["path"])
                ).read_bytes()
                self.assertEqual(record["size"], len(data))
                self.assertEqual(
                    record["sha256"], hashlib.sha256(data).hexdigest()
                )
            stable_evidence = [
                {
                    key: record[key]
                    for key in ("role", "name", "sha256", "size")
                }
                for record in evidence
            ]
            expected_evidence_hash = hashlib.sha256(
                json.dumps(
                    stable_evidence,
                    sort_keys=True,
                    separators=(",", ":"),
                ).encode("utf-8")
            ).hexdigest()
            self.assertEqual(result["evidence_sha256"], expected_evidence_hash)
            components = typing.cast(
                typing.List[typing.Dict[str, object]], result["components"]
            )
            metrics = typing.cast(
                typing.Dict[str, object], components[0]["metrics"]
            )
            self.assertEqual(metrics["connected_components"], 2)
            output = directory / "report"
            results_path, html_path = write_outputs(result, output)
            self.assertTrue(results_path.is_file())
            published = json.loads(results_path.read_text(encoding="utf-8"))
            self.assertEqual(len(published["config_sha256"]), 64)
            self.assertEqual(len(published["evidence_sha256"]), 64)
            snapshot = (
                output / published["publication"]["view_images"]["front"]
            )
            self.assertTrue(snapshot.is_file())
            snapshot_bytes = snapshot.read_bytes()
            replacement_rows = [bytes([70, 80, 90] * 4) for _ in range(3)]
            _write_png(directory / "front.png", 4, replacement_rows, 2)
            self.assertEqual(snapshot.read_bytes(), snapshot_bytes)
            self.assertNotEqual(
                (directory / "front.png").read_bytes(), snapshot_bytes
            )
            report = html_path.read_text(encoding="utf-8")
            self.assertIn("Fixture review", report)
            self.assertIn("Byte evidence", report)
            self.assertIn("sha256:", report)
            self.assertIn("Visual review remains mandatory", report)
            self.assertIn(
                "This review report was generated by an LLM.", report
            )

    def test_publication_reads_each_evidence_file_once(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory))
            with (
                mock.patch.object(
                    review_module,
                    "_hash_file",
                    wraps=review_module._hash_file,
                ) as hash_file,
                mock.patch.object(
                    review_module,
                    "_copy_verified",
                    wraps=review_module._copy_verified,
                ) as copy_verified,
            ):
                write_outputs(result, directory / "report")
            hashed_paths = [call.args[0] for call in hash_file.call_args_list]
            self.assertEqual(len(hashed_paths), len(set(hashed_paths)))
            self.assertEqual(len(hashed_paths), 3)
            self.assertEqual(copy_verified.call_count, 2)

    def test_component_work_limit_preflights_before_second_mask(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            config = json.loads(config_path.read_text(encoding="utf-8"))
            second = dict(config["components"][0])
            second["name"] = "second red panel gate"
            config["components"].append(second)
            config_path.write_text(json.dumps(config), encoding="utf-8")
            with (
                mock.patch.object(
                    review_module,
                    "MAX_COMPONENT_PIXEL_EVALUATIONS",
                    12,
                ),
                mock.patch.object(
                    review_module,
                    "_component_mask",
                    wraps=review_module._component_mask,
                ) as component_mask,
            ):
                with self.assertRaisesRegex(
                    ConfigError, "component pixel evaluations exceed limit"
                ):
                    audit_config(config_path)
            self.assertEqual(component_mask.call_count, 1)

    def test_packet_pixel_limit_preflights_before_decompression(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            with (
                mock.patch.object(
                    review_module,
                    "MAX_TOTAL_IMAGE_PIXELS",
                    12,
                ),
                mock.patch.object(
                    png_module,
                    "_decode_scanlines",
                    wraps=png_module._decode_scanlines,
                ) as decode_scanlines,
            ):
                result = audit_config(self._packet(directory))
            self.assertFalse(result["passed"])
            self.assertEqual(decode_scanlines.call_count, 1)
            errors = typing.cast(typing.List[str], result["errors"])
            self.assertTrue(
                any("remaining packet limit" in error for error in errors)
            )

    def test_failed_decode_still_consumes_packet_pixel_budget(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            config = json.loads(config_path.read_text(encoding="utf-8"))
            config["views"]["front"].pop("mask")
            config["components"] = []
            config_path.write_text(json.dumps(config), encoding="utf-8")
            _write_bad_filter_png(directory / "front.png", 4, 3)
            _write_bad_filter_png(directory / "rear.png", 4, 3)
            with (
                mock.patch.object(
                    review_module,
                    "MAX_TOTAL_IMAGE_PIXELS",
                    12,
                ),
                mock.patch.object(
                    png_module,
                    "_decode_scanlines",
                    wraps=png_module._decode_scanlines,
                ) as decode_scanlines,
            ):
                result = audit_config(config_path)
            self.assertFalse(result["passed"])
            self.assertEqual(decode_scanlines.call_count, 1)
            errors = typing.cast(typing.List[str], result["errors"])
            self.assertTrue(
                any("unsupported PNG filter type" in error for error in errors)
            )
            self.assertTrue(
                any("remaining packet limit" in error for error in errors)
            )

    def test_snapshot_names_do_not_collide_after_sanitizing(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            config = json.loads(config_path.read_text(encoding="utf-8"))
            config["required_views"] = ["a-b", "a_b"]
            config["views"] = {
                "a-b": {"image": "front.png"},
                "a_b": {"image": "front.png"},
            }
            config["components"] = []
            config_path.write_text(json.dumps(config), encoding="utf-8")
            result = audit_config(config_path)
            self.assertFalse(result["passed"])
            self.assertEqual(
                result["errors"], ["no component gates are configured"]
            )
            results_path, _ = write_outputs(result, directory / "report")
            published = json.loads(results_path.read_text(encoding="utf-8"))
            snapshots = published["publication"]["view_images"]
            self.assertNotEqual(snapshots["a-b"], snapshots["a_b"])
            self.assertTrue(
                (directory / "report" / snapshots["a-b"]).is_file()
            )
            self.assertTrue(
                (directory / "report" / snapshots["a_b"]).is_file()
            )

    def test_dimension_mismatch_fails_packet(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory, rear_width=3))
            self.assertFalse(result["passed"])
            errors = typing.cast(typing.List[str], result["errors"])
            self.assertTrue(any("expected 4x3" in error for error in errors))

    def test_replaced_model_cannot_reuse_passing_image_audit(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory))
            self.assertTrue(result["passed"])
            (directory / "model.blend").write_bytes(b"model B mismatch")
            output = directory / "report"
            with self.assertRaisesRegex(ConfigError, "audited input changed"):
                write_outputs(result, output)
            self.assertFalse(output.exists())

    def test_replaced_image_cannot_preserve_pass_report(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory))
            replacement_rows = [bytes([70, 80, 90] * 4) for _ in range(3)]
            _write_png(directory / "front.png", 4, replacement_rows, 2)
            output = directory / "report"
            with self.assertRaisesRegex(ConfigError, "audited input changed"):
                write_outputs(result, output)
            self.assertFalse(output.exists())

    def test_modified_semantic_verdict_cannot_be_published(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory, rear_width=3))
            self.assertFalse(result["passed"])
            result["passed"] = True
            with self.assertRaisesRegex(ConfigError, "semantic audit result"):
                write_outputs(result, directory / "report")

    def test_missing_declared_input_cannot_collide_with_output(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            config = json.loads(config_path.read_text(encoding="utf-8"))
            config["files"]["future"] = "report/results.json"
            config_path.write_text(json.dumps(config), encoding="utf-8")
            result = audit_config(config_path)
            self.assertFalse(result["passed"])
            with self.assertRaisesRegex(ConfigError, "collides"):
                write_outputs(result, directory / "report")

    def test_rejects_non_finite_threshold_number(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            config = json.loads(config_path.read_text(encoding="utf-8"))
            config["components"][0]["thresholds"]["pixel_count"]["max"] = (
                float("nan")
            )
            config_path.write_text(json.dumps(config), encoding="utf-8")
            with self.assertRaisesRegex(ConfigError, "non-finite JSON number"):
                audit_config(config_path)

    def test_rejects_config_output_collision(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            config_path = self._packet(directory)
            result = audit_config(config_path)
            with self.assertRaisesRegex(ConfigError, "collides"):
                write_outputs(result, config_path)

    def test_rejects_dangling_output_symlink(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            result = audit_config(self._packet(directory))
            output = directory / "report"
            output.symlink_to(directory / "missing-target")
            with self.assertRaisesRegex(ConfigError, "symlink"):
                write_outputs(result, output)


if __name__ == "__main__":
    unittest.main()
