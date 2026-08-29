"""Audit fixed-view Blender render packets and write review artifacts."""

import collections
import copy
import ctypes
import errno
import hashlib
import html
import json
import math
import os
import pathlib
import shutil
import stat
import tempfile
import typing
import urllib.parse

from projects.renders.cmd.fumo_review.png import (
    MAX_PNG_FILE_BYTES,
    Image,
    PNGError,
    decode_png_bytes,
    png_dimensions_bytes,
)


class ConfigError(ValueError):
    """Raised when the review configuration is invalid."""


Metric = typing.Union[int, float, None]

MAX_CONFIG_BYTES = 1024 * 1024
MAX_AUDITED_FILE_BYTES = 1024 * 1024 * 1024
MAX_VIEWS = 32
MAX_NAMED_FILES = 64
MAX_COMPONENTS = 256
MAX_TOTAL_IMAGE_PIXELS = 32 * 1024 * 1024
MAX_TOTAL_NAMED_FILE_BYTES = 1024 * 1024 * 1024
MAX_COMPONENT_PIXEL_EVALUATIONS = 256 * 1024 * 1024
_HASH_CHUNK_BYTES = 1024 * 1024
_RENAME_NOREPLACE = 1

SUPPORTED_METRICS = frozenset(
    {
        "pixel_count",
        "pixel_fraction",
        "bbox_x_min",
        "bbox_y_min",
        "bbox_x_max",
        "bbox_y_max",
        "bbox_width",
        "bbox_height",
        "connected_components",
        "max_horizontal_run",
        "max_horizontal_run_pixels",
    }
)


def _as_dict(value: object, label: str) -> typing.Dict[str, object]:
    if not isinstance(value, dict):
        raise ConfigError(f"{label} must be an object")
    result: typing.Dict[str, object] = {}
    for key, child in value.items():
        if not isinstance(key, str):
            raise ConfigError(f"{label} keys must be strings")
        result[key] = child
    return result


def _as_list(value: object, label: str) -> typing.List[object]:
    if not isinstance(value, list):
        raise ConfigError(f"{label} must be an array")
    return value


def _as_string(value: object, label: str) -> str:
    if not isinstance(value, str) or not value:
        raise ConfigError(f"{label} must be a non-empty string")
    return value


def _as_number(value: object, label: str) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise ConfigError(f"{label} must be a number")
    try:
        number = float(value)
    except OverflowError as error:
        raise ConfigError(f"{label} must be finite") from error
    if not math.isfinite(number):
        raise ConfigError(f"{label} must be finite")
    return number


def _string_list(value: object, label: str) -> typing.List[str]:
    result = []
    for index, child in enumerate(_as_list(value, label)):
        result.append(_as_string(child, f"{label}[{index}]"))
    return result


def _resolve_path(
    base: pathlib.Path, value: object, label: str
) -> pathlib.Path:
    path = pathlib.Path(_as_string(value, label))
    if not path.is_absolute():
        path = base / path
    return path.resolve()


def _read_limited_bytes(path: pathlib.Path, maximum: int, label: str) -> bytes:
    try:
        with path.open("rb") as source:
            metadata = os.fstat(source.fileno())
            if not stat.S_ISREG(metadata.st_mode):
                raise ConfigError(f"{label} is not a regular file: {path}")
            if metadata.st_size > maximum:
                raise ConfigError(
                    f"{label} exceeds {maximum} byte limit: {path}"
                )
            data = source.read(maximum + 1)
    except OSError as error:
        raise ConfigError(f"cannot read {label}: {error}") from error
    if len(data) > maximum:
        raise ConfigError(f"{label} exceeds {maximum} byte limit: {path}")
    return data


def _hash_file(
    path: pathlib.Path,
    maximum: int,
    label: str,
    expected_size: typing.Optional[int] = None,
) -> typing.Dict[str, object]:
    digest = hashlib.sha256()
    size = 0
    try:
        with path.open("rb") as source:
            metadata = os.fstat(source.fileno())
            if not stat.S_ISREG(metadata.st_mode):
                raise ConfigError(f"{label} is not a regular file: {path}")
            if expected_size is not None and metadata.st_size != expected_size:
                raise ConfigError(f"audited input changed: {label} at {path}")
            if metadata.st_size > maximum:
                raise ConfigError(
                    f"{label} exceeds {maximum} byte limit: {path}"
                )
            while True:
                block = source.read(_HASH_CHUNK_BYTES)
                if not block:
                    break
                size += len(block)
                if size > maximum:
                    raise ConfigError(
                        f"{label} exceeds {maximum} byte limit: {path}"
                    )
                digest.update(block)
            final_metadata = os.fstat(source.fileno())
    except OSError as error:
        raise ConfigError(f"cannot read {label}: {error}") from error
    if (
        metadata.st_dev != final_metadata.st_dev
        or metadata.st_ino != final_metadata.st_ino
        or metadata.st_size != final_metadata.st_size
        or metadata.st_mtime_ns != final_metadata.st_mtime_ns
        or metadata.st_ctime_ns != final_metadata.st_ctime_ns
        or size != final_metadata.st_size
    ):
        raise ConfigError(f"{label} changed while it was being hashed: {path}")
    return {
        "path": str(path),
        "sha256": digest.hexdigest(),
        "size": size,
        "limit": maximum,
    }


def _bytes_evidence(
    path: pathlib.Path,
    data: bytes,
    maximum: int,
) -> typing.Dict[str, object]:
    return {
        "path": str(path),
        "sha256": hashlib.sha256(data).hexdigest(),
        "size": len(data),
        "limit": maximum,
    }


def _evidence_sha256(
    evidence: typing.Iterable[typing.Mapping[str, object]],
) -> str:
    stable_records = [
        {
            "role": record["role"],
            "name": record["name"],
            "sha256": record["sha256"],
            "size": record["size"],
        }
        for record in evidence
    ]
    return hashlib.sha256(
        json.dumps(
            stable_records,
            sort_keys=True,
            separators=(",", ":"),
        ).encode("utf-8")
    ).hexdigest()


def _audit_result_sha256(result: typing.Mapping[str, object]) -> str:
    payload = {
        key: value
        for key, value in result.items()
        if key not in ("audit_result_sha256", "publication")
    }
    return hashlib.sha256(
        json.dumps(payload, sort_keys=True, separators=(",", ":")).encode(
            "utf-8"
        )
    ).hexdigest()


def _reject_json_constant(value: str) -> typing.NoReturn:
    raise ValueError(f"non-finite JSON number {value}")


def _reject_nonfinite_tree(value: object, label: str) -> None:
    pending = [(label, value)]
    while pending:
        child_label, child = pending.pop()
        if isinstance(child, float) and not math.isfinite(child):
            raise ConfigError(f"{child_label} contains a non-finite number")
        if isinstance(child, list):
            pending.extend(
                (f"{child_label}[{index}]", grandchild)
                for index, grandchild in enumerate(child)
            )
        elif isinstance(child, dict):
            pending.extend(
                (f"{child_label}.{key}", grandchild)
                for key, grandchild in child.items()
            )


def _component_mask(
    image: Image,
    rgb: typing.Tuple[int, int, int],
    tolerance: int,
) -> bytearray:
    mask = bytearray(image.width * image.height)
    for index in range(image.width * image.height):
        offset = index * 3
        if (
            abs(image.rgb[offset] - rgb[0]) <= tolerance
            and abs(image.rgb[offset + 1] - rgb[1]) <= tolerance
            and abs(image.rgb[offset + 2] - rgb[2]) <= tolerance
        ):
            mask[index] = 1
    return mask


def _connected_components(mask: bytearray, width: int, height: int) -> int:
    visited = bytearray(len(mask))
    component_count = 0
    queue: typing.Deque[int] = collections.deque()
    for start, selected in enumerate(mask):
        if not selected or visited[start]:
            continue
        component_count += 1
        visited[start] = 1
        queue.append(start)
        while queue:
            index = queue.popleft()
            for neighbor in _neighbors(index, width, height):
                if mask[neighbor] and not visited[neighbor]:
                    visited[neighbor] = 1
                    queue.append(neighbor)
    return component_count


def _neighbors(index: int, width: int, height: int) -> typing.Iterator[int]:
    x = index % width
    y = index // width
    if x > 0:
        yield index - 1
    if x + 1 < width:
        yield index + 1
    if y > 0:
        yield index - width
    if y + 1 < height:
        yield index + width


def component_metrics(
    image: Image, mask: bytearray
) -> typing.Dict[str, Metric]:
    """Calculate review metrics for a binary component mask."""
    pixel_count = 0
    max_run = 0
    x_min = image.width
    y_min = image.height
    x_max = 0
    y_max = 0
    for y in range(image.height):
        run = 0
        row_start = y * image.width
        for x in range(image.width):
            if mask[row_start + x]:
                pixel_count += 1
                run += 1
                max_run = max(max_run, run)
                x_min = min(x_min, x)
                y_min = min(y_min, y)
                x_max = max(x_max, x + 1)
                y_max = max(y_max, y + 1)
            else:
                run = 0

    result: typing.Dict[str, Metric] = {
        "pixel_count": pixel_count,
        "pixel_fraction": pixel_count / (image.width * image.height),
        "connected_components": _connected_components(
            mask, image.width, image.height
        ),
        "max_horizontal_run": max_run / image.width,
        "max_horizontal_run_pixels": max_run,
    }
    if not pixel_count:
        for name in (
            "bbox_x_min",
            "bbox_y_min",
            "bbox_x_max",
            "bbox_y_max",
            "bbox_width",
            "bbox_height",
        ):
            result[name] = None
        return result

    result.update(
        {
            "bbox_x_min": x_min / image.width,
            "bbox_y_min": y_min / image.height,
            "bbox_x_max": x_max / image.width,
            "bbox_y_max": y_max / image.height,
            "bbox_width": (x_max - x_min) / image.width,
            "bbox_height": (y_max - y_min) / image.height,
        }
    )
    return result


def _parse_rgb(value: object, label: str) -> typing.Tuple[int, int, int]:
    children = _as_list(value, label)
    if len(children) != 3:
        raise ConfigError(f"{label} must contain exactly three channels")
    channels = []
    for index, child in enumerate(children):
        number = _as_number(child, f"{label}[{index}]")
        if number != int(number) or not 0 <= number <= 255:
            raise ConfigError(
                f"{label}[{index}] must be an integer from 0 to 255"
            )
        channels.append(int(number))
    return channels[0], channels[1], channels[2]


def _evaluate_thresholds(
    metrics: typing.Mapping[str, Metric],
    value: object,
    label: str,
) -> typing.List[typing.Dict[str, object]]:
    thresholds = _as_dict(value, label)
    gates: typing.List[typing.Dict[str, object]] = []
    for metric_name, raw_bounds in thresholds.items():
        if metric_name not in SUPPORTED_METRICS:
            raise ConfigError(f"{label} has unknown metric {metric_name!r}")
        bounds = _as_dict(raw_bounds, f"{label}.{metric_name}")
        unknown = set(bounds) - {"min", "max"}
        if unknown:
            raise ConfigError(
                f"{label}.{metric_name} has unknown bounds {sorted(unknown)}"
            )
        if not bounds:
            raise ConfigError(f"{label}.{metric_name} needs min or max")
        minimum = (
            _as_number(bounds["min"], f"{label}.{metric_name}.min")
            if "min" in bounds
            else None
        )
        maximum = (
            _as_number(bounds["max"], f"{label}.{metric_name}.max")
            if "max" in bounds
            else None
        )
        if minimum is not None and maximum is not None and minimum > maximum:
            raise ConfigError(f"{label}.{metric_name} min exceeds max")
        actual = metrics[metric_name]
        passed = actual is not None
        if passed and minimum is not None:
            passed = typing.cast(typing.Union[int, float], actual) >= minimum
        if passed and maximum is not None:
            passed = typing.cast(typing.Union[int, float], actual) <= maximum
        gates.append(
            {
                "metric": metric_name,
                "actual": actual,
                "min": minimum,
                "max": maximum,
                "passed": passed,
            }
        )
    return gates


def _decode(
    path: pathlib.Path,
    label: str,
    role: str,
    name: str,
    errors: typing.List[str],
    evidence: typing.List[typing.Dict[str, object]],
    remaining_pixels: int,
) -> typing.Tuple[typing.Optional[Image], int]:
    charged_pixels = 0
    try:
        data = _read_limited_bytes(path, MAX_PNG_FILE_BYTES, label)
        record = _bytes_evidence(path, data, MAX_PNG_FILE_BYTES)
        record.update({"role": role, "name": name})
        evidence.append(record)
        width, height = png_dimensions_bytes(data)
        charged_pixels = width * height
        if charged_pixels > remaining_pixels:
            raise PNGError(
                "PNG pixel count exceeds remaining packet limit: "
                f"{charged_pixels} > {remaining_pixels}"
            )
        return (
            decode_png_bytes(data, max_pixels=remaining_pixels),
            charged_pixels,
        )
    except (ConfigError, PNGError) as error:
        errors.append(f"{label}: {error}")
        return None, charged_pixels


def _audit_views(
    root: typing.Mapping[str, object],
    base: pathlib.Path,
    errors: typing.List[str],
    evidence: typing.List[typing.Dict[str, object]],
) -> typing.Tuple[
    typing.List[typing.Dict[str, object]],
    typing.Dict[str, Image],
    typing.Dict[str, Image],
    typing.Dict[str, typing.Dict[str, object]],
    typing.Optional[typing.Tuple[int, int]],
]:
    views_value = _as_dict(root.get("views", {}), "views")
    if len(views_value) > MAX_VIEWS:
        raise ConfigError(f"views may contain at most {MAX_VIEWS} entries")
    required_views = _string_list(
        root.get("required_views", []), "required_views"
    )
    if len(required_views) > MAX_VIEWS:
        raise ConfigError(
            f"required_views may contain at most {MAX_VIEWS} entries"
        )
    if not required_views:
        raise ConfigError("required_views must name at least one view")
    for name in required_views:
        if name not in views_value:
            errors.append(f"required view {name!r} is not configured")

    results: typing.List[typing.Dict[str, object]] = []
    loaded_views: typing.Dict[str, Image] = {}
    loaded_masks: typing.Dict[str, Image] = {}
    configs: typing.Dict[str, typing.Dict[str, object]] = {}
    common: typing.Optional[typing.Tuple[int, int]] = None
    total_pixels = 0
    for name, raw_view in views_value.items():
        view = _as_dict(raw_view, f"views.{name}")
        configs[name] = view
        image_path = _resolve_path(
            base, view.get("image"), f"views.{name}.image"
        )
        image, image_pixels = _decode(
            image_path,
            f"view {name!r}",
            "view_image",
            name,
            errors,
            evidence,
            MAX_TOTAL_IMAGE_PIXELS - total_pixels,
        )
        total_pixels += image_pixels
        if image:
            loaded_views[name] = image
            dimensions = (image.width, image.height)
            if common is None:
                common = dimensions
            elif dimensions != common:
                errors.append(
                    f"view {name!r} is {image.width}x{image.height}; "
                    f"expected {common[0]}x{common[1]}"
                )
        mask_path, mask_pixels = _load_view_mask(
            name,
            view,
            base,
            image,
            loaded_masks,
            errors,
            evidence,
            MAX_TOTAL_IMAGE_PIXELS - total_pixels,
        )
        total_pixels += mask_pixels
        image_record = next(
            (
                item
                for item in evidence
                if item["role"] == "view_image" and item["name"] == name
            ),
            None,
        )
        mask_record = next(
            (
                item
                for item in evidence
                if item["role"] == "view_mask" and item["name"] == name
            ),
            None,
        )
        results.append(
            {
                "name": name,
                "required": name in required_views,
                "image": str(image_path),
                "mask": str(mask_path) if mask_path else None,
                "width": image.width if image else None,
                "height": image.height if image else None,
                "decoded": image is not None,
                "image_sha256": (
                    image_record["sha256"] if image_record else None
                ),
                "image_size": image_record["size"] if image_record else None,
                "mask_sha256": (
                    mask_record["sha256"] if mask_record else None
                ),
                "mask_size": mask_record["size"] if mask_record else None,
            }
        )
    return results, loaded_views, loaded_masks, configs, common


def _load_view_mask(
    name: str,
    view: typing.Mapping[str, object],
    base: pathlib.Path,
    image: typing.Optional[Image],
    loaded_masks: typing.MutableMapping[str, Image],
    errors: typing.List[str],
    evidence: typing.List[typing.Dict[str, object]],
    remaining_pixels: int,
) -> typing.Tuple[typing.Optional[pathlib.Path], int]:
    if "mask" not in view:
        return None, 0
    path = _resolve_path(base, view["mask"], f"views.{name}.mask")
    mask_image, charged_pixels = _decode(
        path,
        f"mask for view {name!r}",
        "view_mask",
        name,
        errors,
        evidence,
        remaining_pixels,
    )
    if mask_image:
        loaded_masks[name] = mask_image
        if image and (
            mask_image.width != image.width
            or mask_image.height != image.height
        ):
            errors.append(
                f"mask for view {name!r} is "
                f"{mask_image.width}x{mask_image.height}; "
                f"view is {image.width}x{image.height}"
            )
    return path, charged_pixels


def _audit_files(
    root: typing.Mapping[str, object],
    base: pathlib.Path,
    errors: typing.List[str],
    evidence: typing.List[typing.Dict[str, object]],
) -> typing.List[typing.Dict[str, object]]:
    files_value = _as_dict(root.get("files", {}), "files")
    if len(files_value) > MAX_NAMED_FILES:
        raise ConfigError(
            f"files may contain at most {MAX_NAMED_FILES} entries"
        )
    required_files = _string_list(
        root.get("required_files", []), "required_files"
    )
    if len(required_files) > MAX_NAMED_FILES:
        raise ConfigError(
            f"required_files may contain at most {MAX_NAMED_FILES} entries"
        )
    for name in required_files:
        if name not in files_value:
            errors.append(f"required file {name!r} is not configured")
    results = []
    total_size = 0
    for name, raw_path in files_value.items():
        path = _resolve_path(base, raw_path, f"files.{name}")
        record: typing.Optional[typing.Dict[str, object]] = None
        try:
            remaining_bytes = MAX_TOTAL_NAMED_FILE_BYTES - total_size
            record = _hash_file(
                path,
                min(MAX_AUDITED_FILE_BYTES, remaining_bytes),
                f"file {name!r}",
            )
            record.update({"role": "named_file", "name": name})
            evidence.append(record)
            total_size += typing.cast(int, record["size"])
        except ConfigError as error:
            errors.append(str(error))
        results.append(
            {
                "name": name,
                "path": str(path),
                "required": name in required_files,
                "exists": record is not None,
                "sha256": record["sha256"] if record else None,
                "size": record["size"] if record else None,
            }
        )
    return results


def _component_result(
    component: typing.Mapping[str, object],
    label: str,
    view_configs: typing.Mapping[str, typing.Mapping[str, object]],
    loaded_views: typing.Mapping[str, Image],
    loaded_masks: typing.Mapping[str, Image],
    errors: typing.List[str],
) -> typing.Dict[str, object]:
    name = _as_string(component.get("name"), f"{label}.name")
    view_name = _as_string(component.get("view"), f"{label}.view")
    if view_name not in view_configs:
        raise ConfigError(f"{label}.view names unknown view {view_name!r}")
    rgb = _parse_rgb(component.get("rgb"), f"{label}.rgb")
    tolerance_number = _as_number(
        component.get("tolerance", 0), f"{label}.tolerance"
    )
    if (
        tolerance_number != int(tolerance_number)
        or not 0 <= tolerance_number <= 255
    ):
        raise ConfigError(
            f"{label}.tolerance must be an integer from 0 to 255"
        )
    tolerance = int(tolerance_number)
    use_mask = "mask" in view_configs[view_name]
    source = (
        loaded_masks.get(view_name)
        if use_mask
        else loaded_views.get(view_name)
    )
    if source is None:
        metrics: typing.Dict[str, Metric] = {
            metric: None for metric in SUPPORTED_METRICS
        }
        gates: typing.List[typing.Dict[str, object]] = []
        errors.append(f"component {name!r} has no decodable source image")
    else:
        selected = _component_mask(source, rgb, tolerance)
        metrics = component_metrics(source, selected)
        gates = _evaluate_thresholds(
            metrics,
            component.get("thresholds", {}),
            f"{label}.thresholds",
        )
    return {
        "name": name,
        "view": view_name,
        "rgb": list(rgb),
        "tolerance": tolerance,
        "source": "mask" if use_mask else "image",
        "metrics": metrics,
        "gates": gates,
        "passed": bool(gates)
        and all(typing.cast(bool, gate["passed"]) for gate in gates),
    }


def _component_source(
    component: typing.Mapping[str, object],
    label: str,
    view_configs: typing.Mapping[str, typing.Mapping[str, object]],
    loaded_views: typing.Mapping[str, Image],
    loaded_masks: typing.Mapping[str, Image],
) -> typing.Optional[Image]:
    view_name = _as_string(component.get("view"), f"{label}.view")
    view = view_configs.get(view_name)
    if view is None:
        return None
    return (
        loaded_masks.get(view_name)
        if "mask" in view
        else loaded_views.get(view_name)
    )


def _audit_components(
    root: typing.Mapping[str, object],
    view_configs: typing.Mapping[str, typing.Mapping[str, object]],
    loaded_views: typing.Mapping[str, Image],
    loaded_masks: typing.Mapping[str, Image],
    errors: typing.List[str],
) -> typing.List[typing.Dict[str, object]]:
    values = _as_list(root.get("components", []), "components")
    if len(values) > MAX_COMPONENTS:
        raise ConfigError(
            f"components may contain at most {MAX_COMPONENTS} entries"
        )
    results = []
    names: typing.Set[str] = set()
    pixel_evaluations = 0
    for index, raw_component in enumerate(values):
        label = f"components[{index}]"
        component = _as_dict(raw_component, label)
        source = _component_source(
            component,
            label,
            view_configs,
            loaded_views,
            loaded_masks,
        )
        if source is not None:
            next_evaluations = pixel_evaluations + source.width * source.height
            if next_evaluations > MAX_COMPONENT_PIXEL_EVALUATIONS:
                raise ConfigError(
                    "component pixel evaluations exceed limit: "
                    f"{next_evaluations} > "
                    f"{MAX_COMPONENT_PIXEL_EVALUATIONS}"
                )
            pixel_evaluations = next_evaluations
        result = _component_result(
            component,
            label,
            view_configs,
            loaded_views,
            loaded_masks,
            errors,
        )
        name = typing.cast(str, result["name"])
        if name in names:
            raise ConfigError(f"duplicate component name {name!r}")
        names.add(name)
        results.append(result)
    if not results:
        errors.append("no component gates are configured")
    return results


def audit_config(config_path: pathlib.Path) -> typing.Dict[str, object]:
    """Audit a packet described by ``config_path`` and return its result."""
    resolved_config = config_path.resolve()
    try:
        config_bytes = _read_limited_bytes(
            resolved_config, MAX_CONFIG_BYTES, "config"
        )
        config_text = config_bytes.decode("utf-8")
        root_value: object = json.loads(
            config_text,
            parse_constant=_reject_json_constant,
        )
    except (
        ConfigError,
        UnicodeDecodeError,
        json.JSONDecodeError,
        RecursionError,
        ValueError,
    ) as error:
        raise ConfigError(f"cannot read config: {error}") from error
    root = _as_dict(root_value, "config")
    _reject_nonfinite_tree(root, "config")
    schema_version = root.get("schema_version", 1)
    if schema_version != 1:
        raise ConfigError("schema_version must be 1")
    title_value = root.get("title", "Fumo render review")
    title = _as_string(title_value, "title")
    base = resolved_config.parent
    errors: typing.List[str] = []
    config_evidence = _bytes_evidence(
        resolved_config, config_bytes, MAX_CONFIG_BYTES
    )
    config_evidence.update({"role": "config", "name": "config"})
    evidence = [config_evidence]

    (
        view_results,
        loaded_views,
        loaded_masks,
        view_configs,
        common_dimensions,
    ) = _audit_views(root, base, errors, evidence)
    file_results = _audit_files(root, base, errors, evidence)
    component_results = _audit_components(
        root,
        view_configs,
        loaded_views,
        loaded_masks,
        errors,
    )
    passed = not errors and all(
        typing.cast(bool, component["passed"])
        for component in component_results
    )
    evidence_sha256 = _evidence_sha256(evidence)
    result: typing.Dict[str, object] = {
        "schema_version": 1,
        "title": title,
        "config": str(resolved_config),
        "config_sha256": config_evidence["sha256"],
        "config_size": config_evidence["size"],
        "evidence_sha256": evidence_sha256,
        "evidence": evidence,
        "passed": passed,
        "dimensions": (
            {"width": common_dimensions[0], "height": common_dimensions[1]}
            if common_dimensions
            else None
        ),
        "errors": errors,
        "views": view_results,
        "files": file_results,
        "components": component_results,
    }
    result["audit_result_sha256"] = _audit_result_sha256(result)
    return result


def _evidence_limit(role: str) -> int:
    if role == "config":
        return MAX_CONFIG_BYTES
    if role in ("view_image", "view_mask"):
        return MAX_PNG_FILE_BYTES
    if role == "named_file":
        return MAX_AUDITED_FILE_BYTES
    raise ConfigError(f"unknown evidence role {role!r}")


def _validate_audit_result(result: typing.Mapping[str, object]) -> None:
    expected_result_hash = _as_string(
        result.get("audit_result_sha256"), "audit_result_sha256"
    )
    if _audit_result_sha256(result) != expected_result_hash:
        raise ConfigError("semantic audit result was modified")


def _snapshot_relative_path(
    name: str, digest: str, evidence_index: int
) -> pathlib.Path:
    safe_name = "".join(
        character if character.isalnum() else "_" for character in name
    )[:64]
    if not safe_name:
        safe_name = "view"
    name_digest = hashlib.sha256(name.encode("utf-8")).hexdigest()[:12]
    return pathlib.Path(
        "evidence",
        "views",
        f"{evidence_index:03d}_{safe_name}_{name_digest}_{digest[:16]}.png",
    )


def _validate_evidence(
    result: typing.Mapping[str, object],
    snapshot_dir: typing.Optional[pathlib.Path] = None,
) -> typing.Tuple[
    typing.List[typing.Dict[str, object]], typing.Dict[str, str]
]:
    _validate_audit_result(result)
    raw_evidence = result.get("evidence")
    if not isinstance(raw_evidence, list) or not raw_evidence:
        raise ConfigError("audit result has no byte evidence")
    evidence: typing.List[typing.Dict[str, object]] = []
    view_snapshots: typing.Dict[str, str] = {}
    for index, raw_record in enumerate(raw_evidence):
        record = _as_dict(raw_record, f"evidence[{index}]")
        role = _as_string(record.get("role"), f"evidence[{index}].role")
        name = _as_string(record.get("name"), f"evidence[{index}].name")
        path = pathlib.Path(
            _as_string(record.get("path"), f"evidence[{index}].path")
        )
        expected_hash = _as_string(
            record.get("sha256"), f"evidence[{index}].sha256"
        )
        expected_size = record.get("size")
        if (
            isinstance(expected_size, bool)
            or not isinstance(expected_size, int)
            or expected_size < 0
        ):
            raise ConfigError(f"evidence[{index}].size must be non-negative")
        if expected_size > _evidence_limit(role):
            raise ConfigError(f"evidence[{index}].size exceeds its role limit")
        if snapshot_dir is not None and role == "view_image":
            relative = _snapshot_relative_path(name, expected_hash, index)
            _copy_verified(record, snapshot_dir / relative)
            view_snapshots[name] = relative.as_posix()
        else:
            current = _hash_file(
                path,
                min(_evidence_limit(role), expected_size),
                f"evidence {role} {name!r}",
                expected_size=expected_size,
            )
            if current["sha256"] != expected_hash:
                raise ConfigError(
                    f"audited input changed: {role} {name!r} at {path}"
                )
        evidence.append(record)

    expected_aggregate = _as_string(
        result.get("evidence_sha256"), "evidence_sha256"
    )
    actual_aggregate = _evidence_sha256(evidence)
    if actual_aggregate != expected_aggregate:
        raise ConfigError("audit evidence record set was modified")
    return evidence, view_snapshots


def _lexical_absolute(path: pathlib.Path) -> pathlib.Path:
    return pathlib.Path(os.path.abspath(os.fspath(path)))


def _is_within(path: pathlib.Path, directory: pathlib.Path) -> bool:
    try:
        path.relative_to(directory)
        return True
    except ValueError:
        return False


def _reject_symlink_components(path: pathlib.Path) -> None:
    for candidate in (path, *path.parents):
        if os.path.lexists(candidate) and candidate.is_symlink():
            raise ConfigError(
                f"output path may not traverse a symlink: {candidate}"
            )


def _validate_output_target(
    result: typing.Mapping[str, object], output_dir: pathlib.Path
) -> pathlib.Path:
    output = _lexical_absolute(output_dir)
    _reject_symlink_components(output)
    input_paths = [typing.cast(str, result["config"])]
    for view in typing.cast(
        typing.List[typing.Dict[str, object]], result.get("views", [])
    ):
        input_paths.append(typing.cast(str, view["image"]))
        mask = view.get("mask")
        if mask is not None:
            input_paths.append(typing.cast(str, mask))
    input_paths.extend(
        typing.cast(str, item["path"])
        for item in typing.cast(
            typing.List[typing.Dict[str, object]], result.get("files", [])
        )
    )
    for raw_input_path in input_paths:
        input_path = _lexical_absolute(pathlib.Path(raw_input_path))
        if input_path == output or _is_within(input_path, output):
            raise ConfigError(
                f"output directory collides with audited input: {input_path}"
            )
    if os.path.lexists(output):
        if output.is_symlink():
            raise ConfigError(
                f"output directory may not be a symlink: {output}"
            )
        raise ConfigError(f"output directory already exists: {output}")
    output.parent.mkdir(parents=True, exist_ok=True)
    _reject_symlink_components(output)
    if not output.parent.is_dir():
        raise ConfigError(f"output parent is not a directory: {output.parent}")
    return output


def _rename_noreplace(source: pathlib.Path, destination: pathlib.Path) -> None:
    if source.parent != destination.parent:
        raise ConfigError("atomic publication requires one parent directory")
    libc = ctypes.CDLL(None, use_errno=True)
    try:
        renameat2 = libc.renameat2
    except AttributeError as error:
        raise ConfigError(
            "atomic no-replace publication is unavailable on this host"
        ) from error
    renameat2.argtypes = [
        ctypes.c_int,
        ctypes.c_char_p,
        ctypes.c_int,
        ctypes.c_char_p,
        ctypes.c_uint,
    ]
    renameat2.restype = ctypes.c_int
    parent_fd = os.open(
        source.parent,
        os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW,
    )
    try:
        status = renameat2(
            parent_fd,
            os.fsencode(source.name),
            parent_fd,
            os.fsencode(destination.name),
            _RENAME_NOREPLACE,
        )
    finally:
        os.close(parent_fd)
    if status == 0:
        return
    error_number = ctypes.get_errno()
    if error_number == errno.EEXIST:
        raise ConfigError(
            f"output directory appeared during publish: {destination}"
        )
    raise OSError(error_number, os.strerror(error_number), destination)


def _copy_verified(
    record: typing.Mapping[str, object], destination: pathlib.Path
) -> None:
    source_path = pathlib.Path(typing.cast(str, record["path"]))
    expected_hash = typing.cast(str, record["sha256"])
    expected_size = typing.cast(int, record["size"])
    role = typing.cast(str, record["role"])
    digest = hashlib.sha256()
    size = 0
    destination.parent.mkdir(parents=True, exist_ok=True)
    try:
        with (
            source_path.open("rb") as source,
            destination.open("xb") as target,
        ):
            metadata = os.fstat(source.fileno())
            if not stat.S_ISREG(metadata.st_mode):
                raise ConfigError(f"invalid evidence source: {source_path}")
            if metadata.st_size != expected_size:
                raise ConfigError(
                    f"audited input changed: {role} at {source_path}"
                )
            while True:
                block = source.read(_HASH_CHUNK_BYTES)
                if not block:
                    break
                size += len(block)
                if size > expected_size:
                    raise ConfigError(
                        f"evidence source is too large: {source_path}"
                    )
                digest.update(block)
                target.write(block)
            final_metadata = os.fstat(source.fileno())
            target.flush()
            os.fsync(target.fileno())
    except OSError as error:
        raise ConfigError(
            f"cannot snapshot evidence {source_path}: {error}"
        ) from error
    if (
        metadata.st_dev != final_metadata.st_dev
        or metadata.st_ino != final_metadata.st_ino
        or metadata.st_size != final_metadata.st_size
        or metadata.st_mtime_ns != final_metadata.st_mtime_ns
        or metadata.st_ctime_ns != final_metadata.st_ctime_ns
        or size != expected_size
        or digest.hexdigest() != expected_hash
    ):
        raise ConfigError(f"audited input changed: {role} at {source_path}")


def _format_metric(value: object) -> str:
    if value is None:
        return "n/a"
    if isinstance(value, float):
        return f"{value:.5f}"
    return str(value)


def _publication_view_snapshots(
    result: typing.Mapping[str, object],
) -> typing.Mapping[str, object]:
    publication = result.get("publication")
    if not isinstance(publication, dict):
        return {}
    snapshots = publication.get("view_images")
    return snapshots if isinstance(snapshots, dict) else {}


def _render_html_validated(
    result: typing.Mapping[str, object],
    evidence: typing.Sequence[typing.Mapping[str, object]],
) -> str:
    """Render HTML from evidence already validated by the caller."""
    passed = typing.cast(bool, result["passed"])
    verdict = "PASS" if passed else "FAIL"
    title = html.escape(typing.cast(str, result["title"]))
    parts = [
        "<!doctype html>",
        '<html lang="en"><head><meta charset="utf-8">',
        f"<title>{title}: {verdict}</title>",
        "<style>",
        "body{font:15px system-ui,sans-serif;margin:2rem;background:#17181b;",
        "color:#eee}a{color:#9ecbff}.pass{color:#80e6a3}.fail{color:#ff8f8f}",
        ".views{display:grid;grid-template-columns:",
        "repeat(auto-fit,minmax(260px,1fr));",
        "gap:1rem}.view{background:#24262b;padding:1rem;border-radius:.5rem}",
        ".view img{width:100%;height:auto;background:#111}",
        "table{border-collapse:collapse;",
        "width:100%;margin:1rem 0}th,td{border:1px solid #555;padding:.4rem;",
        "text-align:left}code{color:#f3d38a}</style></head><body>",
        f'<h1>{title} <strong class="{verdict.lower()}">'
        f"{verdict}</strong></h1>",
    ]
    errors = typing.cast(typing.List[str], result["errors"])
    if errors:
        parts.append('<section><h2>Packet errors</h2><ul class="fail">')
        parts.extend(f"<li>{html.escape(error)}</li>" for error in errors)
        parts.append("</ul></section>")

    files = typing.cast(typing.List[typing.Dict[str, object]], result["files"])
    if files:
        parts.append("<section><h2>Named files</h2><ul>")
        for item in files:
            path = typing.cast(str, item["path"])
            status = "PASS" if item["exists"] else "FAIL"
            name = html.escape(typing.cast(str, item["name"]))
            escaped_path = html.escape(path)
            parts.append(
                f'<li class="{status.lower()}">{status} '
                f"{name} <code>{escaped_path}</code></li>"
            )
        parts.append("</ul></section>")

    parts.append("<section><h2>Byte evidence</h2><ul>")
    for record in evidence:
        role = html.escape(typing.cast(str, record["role"]))
        name = html.escape(typing.cast(str, record["name"]))
        digest = html.escape(typing.cast(str, record["sha256"]))
        size = typing.cast(int, record["size"])
        parts.append(
            f"<li><code>{role}:{name}</code> {size} bytes "
            f"<code>sha256:{digest}</code></li>"
        )
    aggregate = html.escape(typing.cast(str, result["evidence_sha256"]))
    parts.append(
        f"</ul><p>Evidence set: <code>sha256:{aggregate}</code></p></section>"
    )

    views = typing.cast(typing.List[typing.Dict[str, object]], result["views"])
    view_snapshots = _publication_view_snapshots(result)
    parts.append('<section><h2>Views</h2><div class="views">')
    for view in views:
        name = html.escape(typing.cast(str, view["name"]))
        report_image = view_snapshots.get(typing.cast(str, view["name"]))
        if report_image is not None:
            url = urllib.parse.quote(
                typing.cast(str, report_image), safe="/._-"
            )
            image_markup = (
                f'<a href="{url}"><img src="{url}" alt="{name}"></a>'
            )
        else:
            image_markup = "<p>No immutable image snapshot is available.</p>"
        dimensions = f"{view['width']}x{view['height']}"
        parts.append(
            f'<article class="view"><h3>{name}</h3>'
            f"{image_markup}"
            f"<p>{html.escape(dimensions)}</p></article>"
        )
    parts.append("</div></section>")

    components = typing.cast(
        typing.List[typing.Dict[str, object]], result["components"]
    )
    parts.append("<section><h2>Component gates</h2>")
    for component in components:
        component_passed = typing.cast(bool, component["passed"])
        component_verdict = "PASS" if component_passed else "FAIL"
        name = html.escape(typing.cast(str, component["name"]))
        parts.append(
            f'<h3>{name} <span class="{component_verdict.lower()}">'
            f"{component_verdict}</span></h3>"
        )
        parts.append(
            "<table><thead><tr><th>Metric</th><th>Actual</th>"
            "<th>Minimum</th><th>Maximum</th><th>Gate</th></tr></thead><tbody>"
        )
        gates = typing.cast(
            typing.List[typing.Dict[str, object]], component["gates"]
        )
        for gate in gates:
            gate_passed = typing.cast(bool, gate["passed"])
            gate_verdict = "PASS" if gate_passed else "FAIL"
            metric = html.escape(typing.cast(str, gate["metric"]))
            parts.append(
                "<tr>"
                f"<td><code>{metric}</code></td>"
                f"<td>{_format_metric(gate['actual'])}</td>"
                f"<td>{_format_metric(gate['min'])}</td>"
                f"<td>{_format_metric(gate['max'])}</td>"
                f'<td class="{gate_verdict.lower()}">{gate_verdict}</td>'
                "</tr>"
            )
        parts.append("</tbody></table>")
    parts.append(
        "</section><p><strong>Visual review remains mandatory.</strong> "
        "These gates reject measurable regressions; they do not judge "
        "likeness, "
        "fabric construction, or appeal.</p>"
        "<p>This review report was generated by an LLM.</p></body></html>"
    )
    return "\n".join(parts)


def render_html(
    result: typing.Mapping[str, object], _output_dir: pathlib.Path
) -> str:
    """Validate evidence and render an inspectable HTML report."""
    evidence, _ = _validate_evidence(result)
    return _render_html_validated(result, evidence)


def write_outputs(
    result: typing.Mapping[str, object], output_dir: pathlib.Path
) -> typing.Tuple[pathlib.Path, pathlib.Path]:
    """Atomically publish a new, evidence-bound audit report directory."""
    _validate_audit_result(result)
    output = _validate_output_target(result, output_dir)
    stage = pathlib.Path(
        tempfile.mkdtemp(prefix=f".{output.name}.stage-", dir=output.parent)
    )
    os.chmod(stage, 0o700)
    published = copy.deepcopy(dict(result))
    try:
        evidence, view_snapshots = _validate_evidence(result, stage)

        published["publication"] = {
            "schema_version": 1,
            "view_images": view_snapshots,
        }

        results_stage = stage / "results.json"
        html_stage = stage / "review.html"
        results_stage.write_text(
            json.dumps(published, indent=2, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        html_stage.write_text(
            _render_html_validated(published, evidence),
            encoding="utf-8",
        )
        for path in (results_stage, html_stage):
            with path.open("rb") as source:
                os.fsync(source.fileno())
        directory_fd = os.open(stage, os.O_RDONLY | os.O_DIRECTORY)
        try:
            os.fsync(directory_fd)
        finally:
            os.close(directory_fd)
        _rename_noreplace(stage, output)
    except Exception:
        if stage.exists():
            shutil.rmtree(stage)
        raise
    return output / "results.json", output / "review.html"
