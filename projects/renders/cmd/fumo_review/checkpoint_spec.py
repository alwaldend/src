"""Pure validation for fast live-Blender checkpoint capsules."""

import dataclasses
import datetime
import hashlib
import json
import os
import pathlib
import re
import typing

from projects.renders.cmd.fumo_review.render_spec import Resolution

_LOADED_SOURCE_SHA256 = hashlib.sha256(
    pathlib.Path(__file__).read_bytes()
).hexdigest()


class CheckpointSpecError(ValueError):
    """Raised when a checkpoint capsule or path is invalid."""


@dataclasses.dataclass(frozen=True)
class Artifact:
    """One frozen repository-relative input and its expected hash."""

    path: str
    sha256: str


@dataclasses.dataclass(frozen=True)
class CheckpointView:
    """One controlling view rendered as beauty and silhouette."""

    name: str
    camera: str


@dataclasses.dataclass(frozen=True)
class CheckpointSpec:
    """Validated source-owned task capsule for one cheap visual gate."""

    checkpoint: str
    scene: str
    resolution: Resolution
    parent: Artifact
    source_packet: Artifact
    hypothesis: str
    views: typing.Mapping[str, CheckpointView]
    vetoes: typing.Tuple[str, ...]
    scene_boolean_requirements: typing.Mapping[str, bool]
    t0_source_packet_frozen_utc: typing.Optional[str]


def load_checkpoint_spec(path: pathlib.Path) -> CheckpointSpec:
    """Load and validate a schema-v1 checkpoint capsule."""
    return load_checkpoint_spec_with_hash(path)[0]


def load_checkpoint_spec_with_hash(
    path: pathlib.Path,
) -> typing.Tuple[CheckpointSpec, str]:
    """Load, hash, decode, and validate one immutable capsule byte string."""
    try:
        payload = path.read_bytes()
        value = json.loads(
            payload.decode("utf-8"),
            object_pairs_hook=_unique_object,
        )
    except (OSError, UnicodeError, json.JSONDecodeError) as error:
        raise CheckpointSpecError(
            f"cannot read spec {path}: {error}"
        ) from error
    return parse_checkpoint_spec(value), hashlib.sha256(payload).hexdigest()


def parse_checkpoint_spec(value: object) -> CheckpointSpec:
    """Validate a decoded schema-v1 checkpoint capsule."""
    root = _object(value, "spec")
    _keys(
        root,
        required={
            "schema_version",
            "checkpoint",
            "scene",
            "resolution",
            "parent",
            "source_packet",
            "hypothesis",
            "views",
            "vetoes",
        },
        optional={
            "scene_boolean_requirements",
            "t0_source_packet_frozen_utc",
        },
        label="spec",
    )
    version = root["schema_version"]
    if isinstance(version, bool) or version != 1:
        raise CheckpointSpecError("schema_version must be the integer 1")
    return CheckpointSpec(
        checkpoint=_identifier(root["checkpoint"], "checkpoint"),
        scene=_text(root["scene"], "scene"),
        resolution=_resolution(root["resolution"]),
        parent=_artifact(root["parent"], "parent"),
        source_packet=_artifact(root["source_packet"], "source_packet"),
        hypothesis=_text(root["hypothesis"], "hypothesis"),
        views=_views(root["views"]),
        vetoes=_vetoes(root["vetoes"]),
        scene_boolean_requirements=_boolean_requirements(
            root.get("scene_boolean_requirements", {})
        ),
        t0_source_packet_frozen_utc=_optional_rfc3339(
            root.get("t0_source_packet_frozen_utc"),
            "t0_source_packet_frozen_utc",
        ),
    )


def resolve_workspace_artifact(
    workspace: pathlib.Path,
    artifact: Artifact,
) -> pathlib.Path:
    """Resolve and hash-check an artifact contained by the workspace."""
    root = workspace.resolve()
    lexical = root.joinpath(*pathlib.PurePosixPath(artifact.path).parts)
    _reject_symlink_components(root, lexical, f"artifact {artifact.path!r}")
    candidate = lexical.resolve()
    _require_contained(candidate, root, f"artifact {artifact.path!r}")
    if not candidate.is_file():
        raise CheckpointSpecError(f"artifact does not exist: {artifact.path}")
    actual = sha256_file(candidate)
    if actual != artifact.sha256:
        raise CheckpointSpecError(
            f"artifact hash mismatch for {artifact.path}: "
            f"expected {artifact.sha256}, got {actual}"
        )
    return candidate


def resolve_workspace_file(
    workspace: pathlib.Path,
    path: pathlib.Path,
    label: str,
) -> pathlib.Path:
    """Resolve a caller path to an existing file inside the workspace."""
    root = workspace.resolve()
    lexical = pathlib.Path(
        os.path.abspath(str(path if path.is_absolute() else root / path))
    )
    _reject_symlink_components(root, lexical, label)
    candidate = lexical.resolve()
    _require_contained(candidate, root, label)
    if not candidate.is_file():
        raise CheckpointSpecError(f"{label} is not a file: {candidate}")
    return candidate


def resolve_out_directory(
    workspace: pathlib.Path,
    path: pathlib.Path,
) -> pathlib.Path:
    """Resolve a checkpoint directory strictly below repository-root out/."""
    root = workspace.resolve()
    out_root = (root / "out").resolve()
    lexical = pathlib.Path(
        os.path.abspath(str(path if path.is_absolute() else root / path))
    )
    _reject_symlink_components(
        out_root, lexical, "checkpoint output directory"
    )
    candidate = lexical.resolve()
    _require_contained(candidate, out_root, "checkpoint output directory")
    if candidate == out_root:
        raise CheckpointSpecError(
            "checkpoint output directory must be below repository-root out/"
        )
    if candidate.exists() and not candidate.is_dir():
        raise CheckpointSpecError(
            f"checkpoint output path is not a directory: {candidate}"
        )
    return candidate


def sha256_file(path: pathlib.Path) -> str:
    """Return the SHA-256 digest of one file."""
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        while True:
            chunk = stream.read(1024 * 1024)
            if not chunk:
                break
            digest.update(chunk)
    return digest.hexdigest()


def canonical_json(value: object) -> str:
    """Encode stable, human-readable JSON for checkpoint evidence."""
    return (
        json.dumps(
            value,
            allow_nan=False,
            indent=2,
            sort_keys=True,
        )
        + "\n"
    )


def _unique_object(
    pairs: typing.List[typing.Tuple[str, object]],
) -> typing.Dict[str, object]:
    result: typing.Dict[str, object] = {}
    for key, value in pairs:
        if key in result:
            raise CheckpointSpecError(f"duplicate JSON key {key!r}")
        result[key] = value
    return result


def _object(value: object, label: str) -> typing.Dict[str, object]:
    if not isinstance(value, dict):
        raise CheckpointSpecError(f"{label} must be an object")
    result: typing.Dict[str, object] = {}
    for key, child in value.items():
        if not isinstance(key, str):
            raise CheckpointSpecError(f"{label} keys must be strings")
        result[key] = child
    return result


def _keys(
    value: typing.Mapping[str, object],
    required: typing.Set[str],
    optional: typing.Set[str],
    label: str,
) -> None:
    missing = required - set(value)
    unknown = set(value) - required - optional
    if missing:
        raise CheckpointSpecError(f"{label} is missing keys {sorted(missing)}")
    if unknown:
        raise CheckpointSpecError(
            f"{label} has unknown keys {sorted(unknown)}"
        )


def _text(value: object, label: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise CheckpointSpecError(f"{label} must be a non-empty string")
    if len(value) > 16384:
        raise CheckpointSpecError(f"{label} is unreasonably long")
    return value.strip()


def _optional_text(value: object, label: str) -> typing.Optional[str]:
    if value is None:
        return None
    return _text(value, label)


def _optional_rfc3339(value: object, label: str) -> typing.Optional[str]:
    if value is None:
        return None
    result = _text(value, label)
    if (
        re.fullmatch(
            r"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}"
            r"(?:\.\d{1,9})?(?:Z|[+-]\d{2}:\d{2})",
            result,
        )
        is None
    ):
        raise CheckpointSpecError(f"{label} must be an RFC 3339 timestamp")
    normalized = result[:-1] + "+00:00" if result.endswith("Z") else result
    try:
        parsed = datetime.datetime.fromisoformat(normalized)
    except ValueError as error:
        raise CheckpointSpecError(
            f"{label} must be an RFC 3339 timestamp"
        ) from error
    if parsed.tzinfo is None:
        raise CheckpointSpecError(f"{label} must include a UTC offset")
    return result


def _identifier(value: object, label: str) -> str:
    if not isinstance(value, str) or not value:
        raise CheckpointSpecError(f"{label} must be a non-empty string")
    if value != value.strip():
        raise CheckpointSpecError(
            f"{label} must not have leading or trailing whitespace"
        )
    result = value
    if re.fullmatch(r"[A-Za-z0-9_]+", result) is None:
        raise CheckpointSpecError(
            f"{label} must contain only ASCII letters, digits, and underscores"
        )
    return result


def _positive_dimension(value: object, label: str) -> int:
    if isinstance(value, bool) or not isinstance(value, int):
        raise CheckpointSpecError(f"{label} must be an integer")
    if not 16 <= value <= 640:
        raise CheckpointSpecError(f"{label} must be from 16 to 640")
    return value


def _resolution(value: object) -> Resolution:
    raw = _object(value, "resolution")
    _keys(
        raw,
        required={"width", "height"},
        optional=set(),
        label="resolution",
    )
    return Resolution(
        width=_positive_dimension(raw["width"], "resolution.width"),
        height=_positive_dimension(raw["height"], "resolution.height"),
    )


def _repository_path(value: object, label: str) -> str:
    raw = _text(value, label)
    if "\\" in raw:
        raise CheckpointSpecError(f"{label} must use forward slashes")
    path = pathlib.PurePosixPath(raw)
    if path.is_absolute() or ".." in path.parts:
        raise CheckpointSpecError(f"{label} must stay inside the repository")
    normalized = path.as_posix()
    if normalized in {"", "."}:
        raise CheckpointSpecError(f"{label} must name a file")
    return normalized


def _sha256(value: object, label: str) -> str:
    result = _text(value, label)
    if re.fullmatch(r"[0-9a-f]{64}", result) is None:
        raise CheckpointSpecError(
            f"{label} must be a lowercase 64-character SHA-256 digest"
        )
    return result


def _artifact(value: object, label: str) -> Artifact:
    raw = _object(value, label)
    _keys(raw, required={"path", "sha256"}, optional=set(), label=label)
    return Artifact(
        path=_repository_path(raw["path"], f"{label}.path"),
        sha256=_sha256(raw["sha256"], f"{label}.sha256"),
    )


def _views(value: object) -> typing.Dict[str, CheckpointView]:
    raw_views = _object(value, "views")
    if not raw_views:
        raise CheckpointSpecError("views must define at least one view")
    if len(raw_views) > 2:
        raise CheckpointSpecError(
            "a fast checkpoint may define at most 2 views"
        )
    result: typing.Dict[str, CheckpointView] = {}
    cameras: typing.Set[str] = set()
    for raw_name in sorted(raw_views):
        name = _identifier(raw_name, "view name")
        label = f"views.{name}"
        raw = _object(raw_views[raw_name], label)
        _keys(raw, required={"camera"}, optional=set(), label=label)
        camera = _text(raw["camera"], f"{label}.camera")
        if camera in cameras:
            raise CheckpointSpecError(
                f"views must use distinct cameras; duplicate {camera!r}"
            )
        cameras.add(camera)
        result[name] = CheckpointView(
            name=name,
            camera=camera,
        )
    return result


def _vetoes(value: object) -> typing.Tuple[str, ...]:
    if not isinstance(value, list) or not value:
        raise CheckpointSpecError("vetoes must be a non-empty array")
    result = tuple(_text(item, "veto") for item in value)
    if len(set(result)) != len(result):
        raise CheckpointSpecError("vetoes must not contain duplicates")
    return result


def _boolean_requirements(value: object) -> typing.Dict[str, bool]:
    raw = _object(value, "scene_boolean_requirements")
    if len(raw) > 32:
        raise CheckpointSpecError(
            "scene_boolean_requirements may define at most 32 properties"
        )
    result: typing.Dict[str, bool] = {}
    for raw_name in sorted(raw):
        name = _identifier(raw_name, "scene boolean requirement name")
        expected = raw[raw_name]
        if not isinstance(expected, bool):
            raise CheckpointSpecError(
                f"scene_boolean_requirements.{name} must be a boolean"
            )
        result[name] = expected
    return result


def _require_contained(
    candidate: pathlib.Path,
    root: pathlib.Path,
    label: str,
) -> None:
    try:
        candidate.relative_to(root)
    except ValueError as error:
        raise CheckpointSpecError(
            f"{label} resolves outside {root}"
        ) from error


def _reject_symlink_components(
    root: pathlib.Path,
    candidate: pathlib.Path,
    label: str,
) -> None:
    try:
        relative = candidate.relative_to(root)
    except ValueError as error:
        raise CheckpointSpecError(f"{label} is outside {root}") from error
    current = root
    for part in relative.parts:
        current /= part
        if current.is_symlink():
            raise CheckpointSpecError(
                f"{label} must not contain symbolic-link components: {current}"
            )
