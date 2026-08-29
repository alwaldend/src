"""Pinned-Blender consumer for immutable live checkpoint snapshots.

The live Blender process only writes ``candidate.blend`` and a final
``snapshot_request.json`` into a newly reserved directory.  This driver claims
that directory with an exclusive ``RENDERING`` marker, reopens the candidate
with auto-execution disabled, and publishes ``READY`` only after every render
and provenance check succeeds.
"""

import argparse
import contextlib
import dataclasses
import datetime
import hashlib
import importlib
import json
import math
import os
import pathlib
import re
import sys
import tempfile
import time
import typing
import uuid

_SCRIPT_DIR = pathlib.Path(__file__).resolve().parent
_RUNFILES_ROOT = _SCRIPT_DIR.parents[3]
_TOOL_SOURCE_PATHS = {
    "checkpoint_render": pathlib.PurePosixPath(
        "projects/renders/cmd/fumo_review/checkpoint_render.py"
    ),
    "checkpoint_spec": pathlib.PurePosixPath(
        "projects/renders/cmd/fumo_review/checkpoint_spec.py"
    ),
    "png": pathlib.PurePosixPath("projects/renders/cmd/fumo_review/png.py"),
    "render_packet": pathlib.PurePosixPath(
        "projects/renders/cmd/fumo_review/render_packet.py"
    ),
    "render_spec": pathlib.PurePosixPath(
        "projects/renders/cmd/fumo_review/render_spec.py"
    ),
}
_PREIMPORT_TOOL_PATHS = {
    name: _RUNFILES_ROOT.joinpath(*relative.parts).resolve()
    for name, relative in _TOOL_SOURCE_PATHS.items()
}
_PREIMPORT_TOOL_HASHES = {
    name: hashlib.sha256(path.read_bytes()).hexdigest()
    for name, path in _PREIMPORT_TOOL_PATHS.items()
}
_LOADED_SOURCE_SHA256 = _PREIMPORT_TOOL_HASHES["checkpoint_render"]
if str(_RUNFILES_ROOT) not in sys.path:
    sys.path.insert(0, str(_RUNFILES_ROOT))

from projects.renders.cmd.fumo_review import (  # noqa: E402
    checkpoint_spec as checkpoint_spec_module,
)
from projects.renders.cmd.fumo_review import png as png_module  # noqa: E402
from projects.renders.cmd.fumo_review import (  # noqa: E402
    render_packet as render_packet_module,
)
from projects.renders.cmd.fumo_review import (  # noqa: E402
    render_spec as render_spec_module,
)
from projects.renders.cmd.fumo_review.checkpoint_spec import (  # noqa: E402
    CheckpointSpec,
    CheckpointSpecError,
    canonical_json,
    load_checkpoint_spec_with_hash,
    resolve_out_directory,
    resolve_workspace_file,
    sha256_file,
)
from projects.renders.cmd.fumo_review.png import (  # noqa: E402
    Image,
    PNGError,
    decode_png,
)
from projects.renders.cmd.fumo_review.render_packet import (  # noqa: E402
    RenderError as PacketRenderError,
)
from projects.renders.cmd.fumo_review.render_packet import (  # noqa: E402
    _canonicalize_png,
)

bpy: typing.Any = importlib.import_module("bpy")

_IMPORTED_TOOL_MODULES = {
    "checkpoint_spec": checkpoint_spec_module,
    "png": png_module,
    "render_packet": render_packet_module,
    "render_spec": render_spec_module,
}
for _tool_name, _tool_module in _IMPORTED_TOOL_MODULES.items():
    _module_file = pathlib.Path(
        typing.cast(str, _tool_module.__file__)
    ).resolve()
    if _module_file.suffix == ".pyc":
        _module_file = _module_file.with_suffix(".py")
    if _module_file != _PREIMPORT_TOOL_PATHS[_tool_name]:
        raise RuntimeError(
            f"imported {_tool_name} from unexpected source {_module_file}"
        )
for _tool_name, _tool_path in _PREIMPORT_TOOL_PATHS.items():
    _postimport_hash = hashlib.sha256(_tool_path.read_bytes()).hexdigest()
    if _postimport_hash != _PREIMPORT_TOOL_HASHES[_tool_name]:
        raise RuntimeError(f"{_tool_name} changed while renderer imported it")
_LOADED_TOOL_HASHES = dict(_PREIMPORT_TOOL_HASHES)

_REQUEST_NAME = "snapshot_request.json"
_CANDIDATE_NAME = "candidate.blend"
_LOCK_NAME = "RENDERING"
_READY_NAME = "READY"
_MAX_REQUEST_BYTES = 16 * 1024 * 1024
_MAX_DIMENSION = 640
_PREVIEW_SAMPLES = 32
_GEOMETRY_TYPES = {
    "CURVE",
    "CURVES",
    "FONT",
    "GREASEPENCIL",
    "MESH",
    "META",
    "POINTCLOUD",
    "SURFACE",
    "VOLUME",
}
_TIMELINE_KEYS = {
    "t0_source_packet_frozen_utc",
    "t1_candidate_saved_utc",
    "t2_snapshot_hashed_utc",
    "t3_first_pixels_written_utc",
    "t4_visual_verdict_complete_utc",
    "t5_goal_records_updated_utc",
}


class CheckpointRenderError(RuntimeError):
    """Raised when a snapshot cannot become a trustworthy render packet."""


@dataclasses.dataclass(frozen=True)
class _Artifact:
    path: str
    sha256: str


@dataclasses.dataclass(frozen=True)
class _Dependency:
    kind: str
    name: str
    path: str
    sha256: str


@dataclasses.dataclass(frozen=True)
class _Context:
    active_scene: str
    active_view_layer: str
    active_object: typing.Optional[str]
    selected_objects: typing.Tuple[str, ...]
    mode: str
    source_filepath: str


@dataclasses.dataclass(frozen=True)
class _Camera:
    name: str
    projection: str
    lens: float
    ortho_scale: float
    sensor_fit: str
    sensor_width: float
    sensor_height: float
    shift_x: float
    shift_y: float
    clip_start: float
    clip_end: float
    matrix_world: typing.Tuple[typing.Tuple[float, ...], ...]


@dataclasses.dataclass(frozen=True)
class _Object:
    name: str
    type: str
    data: typing.Optional[str]
    hide_render: bool


@dataclasses.dataclass(frozen=True)
class _ParentOpenAttestation:
    path: str
    sha256: str
    spec_sha256: str
    source_packet_sha256: str
    opened_utc: str
    parent_file_version: typing.Tuple[int, int, int]
    load_handlers_suspended: typing.Mapping[str, int]


@dataclasses.dataclass(frozen=True)
class _PinnedBlender:
    version: str
    toolchain: _Artifact


@dataclasses.dataclass(frozen=True)
class _Live:
    blender_version: str
    tool_sources: typing.Tuple[_Artifact, ...]
    context_before_snapshot: _Context
    context_after_snapshot: _Context
    dirty_before_snapshot: bool
    dirty_state_preserved: bool
    exact_parent_verified: bool
    save_handlers_suspended: typing.Mapping[str, int]
    parent_open_attestation: _ParentOpenAttestation
    pinned_blender: _PinnedBlender


@dataclasses.dataclass(frozen=True)
class _Scene:
    name: str
    boolean_requirements: typing.Mapping[str, bool]
    views: typing.Mapping[str, _Camera]
    view_layers: typing.Tuple[str, ...]
    objects: typing.Tuple[_Object, ...]


@dataclasses.dataclass(frozen=True)
class _Request:
    checkpoint: str
    spec: _Artifact
    parent: _Artifact
    source_packet: _Artifact
    candidate: _Artifact
    dependencies: typing.Tuple[_Dependency, ...]
    live: _Live
    scene: _Scene
    width: int
    height: int
    hypothesis: str
    vetoes: typing.Tuple[str, ...]
    timeline: typing.Mapping[str, typing.Optional[str]]
    timings_seconds: typing.Mapping[str, float]


@dataclasses.dataclass(frozen=True)
class _FrozenInput:
    role: str
    record: _Artifact
    path: pathlib.Path


@dataclasses.dataclass(frozen=True)
class _RendererSource:
    module: str
    path: pathlib.Path
    runtime_path: pathlib.Path
    sha256: str


@dataclasses.dataclass(frozen=True)
class _SilhouetteStats:
    dark_mask: bytearray
    dark_indices: typing.Tuple[int, ...]
    light: int
    border_light: int
    chromatic: int
    bounds: typing.Tuple[int, int, int, int]


@dataclasses.dataclass(frozen=True)
class _Run:
    workspace: pathlib.Path
    output_dir: pathlib.Path
    request_path: pathlib.Path
    request_sha256: str
    request: _Request
    spec: CheckpointSpec
    inputs: typing.Tuple[_FrozenInput, ...]
    tool_sources: typing.Tuple[_RendererSource, ...]


def _script_args() -> typing.List[str]:
    try:
        separator = sys.argv.index("--")
    except ValueError:
        return sys.argv[1:]
    return sys.argv[separator + 1 :]  # noqa: E203


def _resolve_workspace() -> pathlib.Path:
    configured = os.environ.get("BUILD_WORKSPACE_DIRECTORY")
    workspace = (
        pathlib.Path(configured).resolve()
        if configured
        else _RUNFILES_ROOT.resolve()
    )
    out_root = workspace / "out"
    if (
        not workspace.is_dir()
        or not out_root.is_dir()
        or out_root.is_symlink()
    ):
        raise CheckpointRenderError(
            "workspace must exist and contain repository-root out/: "
            f"{workspace}"
        )
    return workspace


def _resolve_request_path(
    workspace: pathlib.Path,
    option: pathlib.Path,
) -> pathlib.Path:
    unresolved = option if option.is_absolute() else workspace / option
    if unresolved.name != _REQUEST_NAME:
        raise CheckpointRenderError(
            f"--snapshot-request must name {_REQUEST_NAME!r}"
        )
    output_dir = resolve_out_directory(workspace, unresolved.parent)
    request_path = resolve_workspace_file(
        workspace,
        unresolved,
        "snapshot request",
    )
    out_root = (workspace / "out").resolve()
    _require_below(request_path, out_root, "snapshot request")
    if request_path.parent == out_root:
        raise CheckpointRenderError(
            "snapshot request output must be below repository-root out/"
        )
    if request_path.parent != output_dir:
        raise CheckpointRenderError(
            "snapshot request changed output-directory identity"
        )
    if not request_path.is_file():
        raise CheckpointRenderError(
            f"snapshot request is not a file: {request_path}"
        )
    return request_path


def _require_below(
    path: pathlib.Path,
    root: pathlib.Path,
    label: str,
) -> None:
    try:
        path.relative_to(root)
    except ValueError as error:
        raise CheckpointRenderError(
            f"{label} resolves outside {root}: {path}"
        ) from error


def _load_request(path: pathlib.Path) -> typing.Tuple[_Request, str]:
    try:
        payload = path.read_bytes()
    except OSError as error:
        raise CheckpointRenderError(
            f"cannot read snapshot request {path}: {error}"
        ) from error
    if len(payload) > _MAX_REQUEST_BYTES:
        raise CheckpointRenderError(
            f"snapshot request exceeds {_MAX_REQUEST_BYTES} bytes"
        )
    try:
        value = json.loads(
            payload.decode("utf-8"),
            object_pairs_hook=_unique_object,
            parse_constant=_reject_json_constant,
        )
    except (UnicodeError, json.JSONDecodeError) as error:
        raise CheckpointRenderError(
            f"invalid snapshot request {path}: {error}"
        ) from error
    request = _parse_request(value)
    return request, hashlib.sha256(payload).hexdigest()


def _unique_object(
    pairs: typing.List[typing.Tuple[str, object]],
) -> typing.Dict[str, object]:
    result: typing.Dict[str, object] = {}
    for key, value in pairs:
        if key in result:
            raise CheckpointRenderError(f"duplicate JSON key {key!r}")
        result[key] = value
    return result


def _reject_json_constant(value: str) -> typing.NoReturn:
    raise CheckpointRenderError(f"invalid JSON number {value}")


def _parse_request(value: object) -> _Request:
    root = _object(value, "snapshot request")
    _keys(
        root,
        required={
            "schema_version",
            "checkpoint",
            "spec",
            "parent",
            "source_packet",
            "candidate",
            "dependencies",
            "live",
            "scene",
            "dimensions",
            "hypothesis",
            "vetoes",
            "timeline",
            "timings_seconds",
        },
        label="snapshot request",
    )
    version = root["schema_version"]
    if isinstance(version, bool) or version != 1:
        raise CheckpointRenderError("schema_version must be the integer 1")
    dimensions = _dimensions(root["dimensions"])
    return _Request(
        checkpoint=_identifier(root["checkpoint"], "checkpoint"),
        spec=_artifact(root["spec"], "spec"),
        parent=_artifact(root["parent"], "parent"),
        source_packet=_artifact(root["source_packet"], "source_packet"),
        candidate=_artifact(root["candidate"], "candidate"),
        dependencies=_dependencies(root["dependencies"]),
        live=_live(root["live"]),
        scene=_scene(root["scene"]),
        width=dimensions[0],
        height=dimensions[1],
        hypothesis=_text(root["hypothesis"], "hypothesis"),
        vetoes=_vetoes(root["vetoes"]),
        timeline=_timeline(root["timeline"]),
        timings_seconds=_timings(root["timings_seconds"]),
    )


def _object(value: object, label: str) -> typing.Dict[str, object]:
    if not isinstance(value, dict):
        raise CheckpointRenderError(f"{label} must be an object")
    if not all(isinstance(key, str) for key in value):
        raise CheckpointRenderError(f"{label} keys must be strings")
    return typing.cast(typing.Dict[str, object], value)


def _keys(
    value: typing.Mapping[str, object],
    required: typing.Set[str],
    label: str,
) -> None:
    missing = required - set(value)
    unknown = set(value) - required
    if missing:
        raise CheckpointRenderError(
            f"{label} is missing keys {sorted(missing)}"
        )
    if unknown:
        raise CheckpointRenderError(
            f"{label} has unknown keys {sorted(unknown)}"
        )


def _text(value: object, label: str) -> str:
    if not isinstance(value, str) or value == "":
        raise CheckpointRenderError(f"{label} must be a non-empty string")
    if len(value) > 16384:
        raise CheckpointRenderError(f"{label} is unreasonably long")
    return value


def _optional_text(value: object, label: str) -> typing.Optional[str]:
    if value is None:
        return None
    return _text(value, label)


def _identifier(value: object, label: str) -> str:
    result = _text(value, label)
    if re.fullmatch(r"[A-Za-z0-9_]+", result) is None:
        raise CheckpointRenderError(
            f"{label} must contain only ASCII letters, digits, and underscores"
        )
    return result


def _repository_path(value: object, label: str) -> str:
    raw = _text(value, label)
    if "\\" in raw:
        raise CheckpointRenderError(f"{label} must use forward slashes")
    path = pathlib.PurePosixPath(raw)
    normalized = path.as_posix()
    if (
        path.is_absolute()
        or ".." in path.parts
        or normalized in {"", "."}
        or normalized != raw
    ):
        raise CheckpointRenderError(
            f"{label} must be a canonical repository-relative file path"
        )
    return normalized


def _sha256(value: object, label: str) -> str:
    result = _text(value, label)
    if re.fullmatch(r"[0-9a-f]{64}", result) is None:
        raise CheckpointRenderError(
            f"{label} must be a lowercase 64-character SHA-256 digest"
        )
    return result


def _artifact(value: object, label: str) -> _Artifact:
    raw = _object(value, label)
    _keys(raw, {"path", "sha256"}, label)
    return _Artifact(
        path=_repository_path(raw["path"], f"{label}.path"),
        sha256=_sha256(raw["sha256"], f"{label}.sha256"),
    )


def _dependencies(value: object) -> typing.Tuple[_Dependency, ...]:
    if not isinstance(value, list):
        raise CheckpointRenderError("dependencies must be an array")
    result: typing.List[_Dependency] = []
    for index, child in enumerate(value):
        label = f"dependencies[{index}]"
        raw = _object(child, label)
        _keys(raw, {"kind", "name", "path", "sha256"}, label)
        result.append(
            _Dependency(
                kind=_text(raw["kind"], f"{label}.kind"),
                name=_text(raw["name"], f"{label}.name"),
                path=_repository_path(raw["path"], f"{label}.path"),
                sha256=_sha256(raw["sha256"], f"{label}.sha256"),
            )
        )
    keys = [(item.kind, item.name, item.path) for item in result]
    if keys != sorted(keys):
        raise CheckpointRenderError(
            "dependencies must use canonical sort order"
        )
    if len(keys) != len(set(keys)):
        raise CheckpointRenderError(
            "dependencies must use canonical deduplication"
        )
    return tuple(result)


def _live(value: object) -> _Live:
    raw = _object(value, "live")
    _keys(
        raw,
        {
            "blender_version",
            "tool_sources",
            "context_before_snapshot",
            "context_after_snapshot",
            "dirty_before_snapshot",
            "dirty_state_preserved",
            "exact_parent_verified",
            "save_handlers_suspended",
            "parent_open_attestation",
            "pinned_blender",
        },
        "live",
    )
    tools = _tool_sources(raw["tool_sources"])
    handlers = _handler_counts(
        raw["save_handlers_suspended"],
        "live.save_handlers_suspended",
        {"save_pre", "save_post", "save_post_fail"},
    )
    return _Live(
        blender_version=_text(raw["blender_version"], "live.blender_version"),
        tool_sources=tools,
        context_before_snapshot=_context(
            raw["context_before_snapshot"],
            "live.context_before_snapshot",
        ),
        context_after_snapshot=_context(
            raw["context_after_snapshot"],
            "live.context_after_snapshot",
        ),
        dirty_before_snapshot=_boolean(
            raw["dirty_before_snapshot"],
            "live.dirty_before_snapshot",
        ),
        dirty_state_preserved=_true(
            raw["dirty_state_preserved"],
            "live.dirty_state_preserved",
        ),
        exact_parent_verified=_true(
            raw["exact_parent_verified"],
            "live.exact_parent_verified",
        ),
        save_handlers_suspended=handlers,
        parent_open_attestation=_parent_open_attestation(
            raw["parent_open_attestation"]
        ),
        pinned_blender=_pinned_blender(raw["pinned_blender"]),
    )


def _tool_sources(value: object) -> typing.Tuple[_Artifact, ...]:
    if not isinstance(value, list) or not value:
        raise CheckpointRenderError(
            "live.tool_sources must be a non-empty array"
        )
    if len(value) > 32:
        raise CheckpointRenderError(
            "live.tool_sources may contain at most 32 items"
        )
    result = tuple(
        _artifact(child, f"live.tool_sources[{index}]")
        for index, child in enumerate(value)
    )
    paths = [item.path for item in result]
    if len(paths) != len(set(paths)):
        raise CheckpointRenderError(
            "live.tool_sources contains duplicate paths"
        )
    return result


def _handler_counts(
    value: object,
    label: str,
    allowed: typing.Set[str],
) -> typing.Dict[str, int]:
    raw = _object(value, label)
    if set(raw) - allowed:
        raise CheckpointRenderError(f"{label} has unknown handler names")
    result: typing.Dict[str, int] = {}
    for name, count in raw.items():
        if isinstance(count, bool) or not isinstance(count, int) or count < 0:
            raise CheckpointRenderError(
                f"{label}.{name} must be a nonnegative integer"
            )
        result[name] = count
    return result


def _parent_open_attestation(value: object) -> _ParentOpenAttestation:
    label = "live.parent_open_attestation"
    raw = _object(value, label)
    _keys(
        raw,
        {
            "path",
            "sha256",
            "spec_sha256",
            "source_packet_sha256",
            "opened_utc",
            "parent_file_version",
            "load_handlers_suspended",
        },
        label,
    )
    opened_utc = _optional_timestamp(raw["opened_utc"], f"{label}.opened_utc")
    if opened_utc is None:
        raise CheckpointRenderError(f"{label}.opened_utc must be populated")
    return _ParentOpenAttestation(
        path=_repository_path(raw["path"], f"{label}.path"),
        sha256=_sha256(raw["sha256"], f"{label}.sha256"),
        spec_sha256=_sha256(raw["spec_sha256"], f"{label}.spec_sha256"),
        source_packet_sha256=_sha256(
            raw["source_packet_sha256"],
            f"{label}.source_packet_sha256",
        ),
        opened_utc=opened_utc,
        parent_file_version=_version_triplet(
            raw["parent_file_version"],
            f"{label}.parent_file_version",
        ),
        load_handlers_suspended=_handler_counts(
            raw["load_handlers_suspended"],
            f"{label}.load_handlers_suspended",
            {"load_pre", "load_post", "load_post_fail", "version_update"},
        ),
    )


def _pinned_blender(value: object) -> _PinnedBlender:
    label = "live.pinned_blender"
    raw = _object(value, label)
    _keys(raw, {"version", "toolchain"}, label)
    toolchain = _artifact(raw["toolchain"], f"{label}.toolchain")
    if toolchain.path != "tools/blender/binary_toolchain.json":
        raise CheckpointRenderError(
            f"{label}.toolchain must name tools/blender/binary_toolchain.json"
        )
    version = _text(raw["version"], f"{label}.version")
    _semantic_version(version, f"{label}.version")
    return _PinnedBlender(version=version, toolchain=toolchain)


def _version_triplet(
    value: object,
    label: str,
) -> typing.Tuple[int, int, int]:
    if not isinstance(value, list) or len(value) != 3:
        raise CheckpointRenderError(f"{label} must contain three integers")
    parts: typing.List[int] = []
    for index, part in enumerate(value):
        if isinstance(part, bool) or not isinstance(part, int) or part < 0:
            raise CheckpointRenderError(
                f"{label}[{index}] must be a nonnegative integer"
            )
        parts.append(part)
    return (parts[0], parts[1], parts[2])


def _semantic_version(
    value: str,
    label: str,
) -> typing.Tuple[int, int, int]:
    match = re.fullmatch(r"(\d+)\.(\d+)\.(\d+)", value)
    if match is None:
        raise CheckpointRenderError(
            f"{label} must be a three-component numeric version"
        )
    return (int(match.group(1)), int(match.group(2)), int(match.group(3)))


def _declared_live_version(value: str) -> typing.Tuple[int, int, int]:
    match = re.fullmatch(r"(\d+)\.(\d+)\.(\d+)(?:\s+.*)?", value)
    if match is None:
        raise CheckpointRenderError(
            "live.blender_version must begin with a numeric version triplet"
        )
    return (int(match.group(1)), int(match.group(2)), int(match.group(3)))


def _context(value: object, label: str) -> _Context:
    raw = _object(value, label)
    _keys(
        raw,
        {
            "active_scene",
            "active_view_layer",
            "active_object",
            "selected_objects",
            "mode",
            "source_filepath",
        },
        label,
    )
    selected = _text_array(
        raw["selected_objects"], f"{label}.selected_objects"
    )
    if tuple(sorted(selected)) != selected:
        raise CheckpointRenderError(f"{label}.selected_objects must be sorted")
    return _Context(
        active_scene=_text(raw["active_scene"], f"{label}.active_scene"),
        active_view_layer=_text(
            raw["active_view_layer"], f"{label}.active_view_layer"
        ),
        active_object=_optional_text(
            raw["active_object"], f"{label}.active_object"
        ),
        selected_objects=selected,
        mode=_text(raw["mode"], f"{label}.mode"),
        source_filepath=_text(
            raw["source_filepath"], f"{label}.source_filepath"
        ),
    )


def _scene(value: object) -> _Scene:
    raw = _object(value, "scene")
    _keys(
        raw,
        {"name", "boolean_requirements", "views", "view_layers", "objects"},
        "scene",
    )
    views = _views(raw["views"])
    layers = _text_array(raw["view_layers"], "scene.view_layers")
    if not layers or tuple(sorted(layers)) != layers:
        raise CheckpointRenderError(
            "scene.view_layers must be a non-empty sorted array"
        )
    return _Scene(
        name=_text(raw["name"], "scene.name"),
        boolean_requirements=_boolean_map(
            raw["boolean_requirements"], "scene.boolean_requirements"
        ),
        views=views,
        view_layers=layers,
        objects=_objects(raw["objects"]),
    )


def _views(value: object) -> typing.Dict[str, _Camera]:
    raw = _object(value, "scene.views")
    if not 1 <= len(raw) <= 2:
        raise CheckpointRenderError("scene.views must define one or two views")
    result: typing.Dict[str, _Camera] = {}
    camera_names: typing.Set[str] = set()
    for raw_name in sorted(raw):
        name = _identifier(raw_name, "scene view name")
        camera = _camera(raw[raw_name], f"scene.views.{name}")
        if camera.name in camera_names:
            raise CheckpointRenderError(
                "scene views must use distinct cameras"
            )
        camera_names.add(camera.name)
        result[name] = camera
    return result


def _camera(value: object, label: str) -> _Camera:
    raw = _object(value, label)
    _keys(
        raw,
        {
            "name",
            "projection",
            "lens",
            "ortho_scale",
            "sensor_fit",
            "sensor_width",
            "sensor_height",
            "shift_x",
            "shift_y",
            "clip_start",
            "clip_end",
            "matrix_world",
        },
        label,
    )
    projection = _text(raw["projection"], f"{label}.projection")
    if projection not in {"PERSP", "ORTHO"}:
        raise CheckpointRenderError(
            f"{label}.projection must be PERSP or ORTHO"
        )
    sensor_fit = _text(raw["sensor_fit"], f"{label}.sensor_fit")
    if sensor_fit not in {"AUTO", "HORIZONTAL", "VERTICAL"}:
        raise CheckpointRenderError(f"{label}.sensor_fit is unsupported")
    clip_start = _positive_float(raw["clip_start"], f"{label}.clip_start")
    clip_end = _positive_float(raw["clip_end"], f"{label}.clip_end")
    if clip_end <= clip_start:
        raise CheckpointRenderError(
            f"{label}.clip_end must exceed {label}.clip_start"
        )
    return _Camera(
        name=_text(raw["name"], f"{label}.name"),
        projection=projection,
        lens=_positive_float(raw["lens"], f"{label}.lens"),
        ortho_scale=_positive_float(
            raw["ortho_scale"], f"{label}.ortho_scale"
        ),
        sensor_fit=sensor_fit,
        sensor_width=_positive_float(
            raw["sensor_width"], f"{label}.sensor_width"
        ),
        sensor_height=_positive_float(
            raw["sensor_height"], f"{label}.sensor_height"
        ),
        shift_x=_finite_float(raw["shift_x"], f"{label}.shift_x"),
        shift_y=_finite_float(raw["shift_y"], f"{label}.shift_y"),
        clip_start=clip_start,
        clip_end=clip_end,
        matrix_world=_matrix(raw["matrix_world"], f"{label}.matrix_world"),
    )


def _matrix(
    value: object,
    label: str,
) -> typing.Tuple[typing.Tuple[float, ...], ...]:
    if not isinstance(value, list) or len(value) != 4:
        raise CheckpointRenderError(f"{label} must be a 4 by 4 array")
    rows: typing.List[typing.Tuple[float, ...]] = []
    for row_index, row in enumerate(value):
        if not isinstance(row, list) or len(row) != 4:
            raise CheckpointRenderError(f"{label} must be a 4 by 4 array")
        rows.append(
            tuple(
                _finite_float(child, f"{label}[{row_index}][{column}]")
                for column, child in enumerate(row)
            )
        )
    return tuple(rows)


def _objects(value: object) -> typing.Tuple[_Object, ...]:
    if not isinstance(value, list) or not value:
        raise CheckpointRenderError("scene.objects must be a non-empty array")
    result: typing.List[_Object] = []
    for index, child in enumerate(value):
        label = f"scene.objects[{index}]"
        raw = _object(child, label)
        _keys(raw, {"name", "type", "data", "hide_render"}, label)
        result.append(
            _Object(
                name=_text(raw["name"], f"{label}.name"),
                type=_text(raw["type"], f"{label}.type"),
                data=_optional_text(raw["data"], f"{label}.data"),
                hide_render=_boolean(
                    raw["hide_render"], f"{label}.hide_render"
                ),
            )
        )
    names = [item.name for item in result]
    if names != sorted(names) or len(names) != len(set(names)):
        raise CheckpointRenderError(
            "scene.objects must be sorted by unique object name"
        )
    return tuple(result)


def _boolean_map(value: object, label: str) -> typing.Dict[str, bool]:
    raw = _object(value, label)
    if len(raw) > 64:
        raise CheckpointRenderError(f"{label} may contain at most 64 entries")
    result: typing.Dict[str, bool] = {}
    for raw_name in sorted(raw):
        name = _identifier(raw_name, f"{label} property name")
        result[name] = _boolean(raw[raw_name], f"{label}.{name}")
    return result


def _boolean(value: object, label: str) -> bool:
    if not isinstance(value, bool):
        raise CheckpointRenderError(f"{label} must be a boolean")
    return value


def _true(value: object, label: str) -> bool:
    result = _boolean(value, label)
    if not result:
        raise CheckpointRenderError(f"{label} must be true")
    return result


def _finite_float(value: object, label: str) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise CheckpointRenderError(f"{label} must be a number")
    result = float(value)
    if not math.isfinite(result):
        raise CheckpointRenderError(f"{label} must be finite")
    return result


def _positive_float(value: object, label: str) -> float:
    result = _finite_float(value, label)
    if result <= 0.0:
        raise CheckpointRenderError(f"{label} must be positive")
    return result


def _dimensions(value: object) -> typing.Tuple[int, int]:
    raw = _object(value, "dimensions")
    _keys(raw, {"width", "height"}, "dimensions")
    return (
        _dimension(raw["width"], "dimensions.width"),
        _dimension(raw["height"], "dimensions.height"),
    )


def _dimension(value: object, label: str) -> int:
    if isinstance(value, bool) or not isinstance(value, int):
        raise CheckpointRenderError(f"{label} must be an integer")
    if not 16 <= value <= _MAX_DIMENSION:
        raise CheckpointRenderError(
            f"{label} must be between 16 and {_MAX_DIMENSION}"
        )
    return value


def _vetoes(value: object) -> typing.Tuple[str, ...]:
    result = _text_array(value, "vetoes")
    if not result:
        raise CheckpointRenderError("vetoes must be non-empty")
    return result


def _text_array(value: object, label: str) -> typing.Tuple[str, ...]:
    if not isinstance(value, list):
        raise CheckpointRenderError(f"{label} must be an array")
    result = tuple(_text(child, f"{label} item") for child in value)
    if len(result) != len(set(result)):
        raise CheckpointRenderError(f"{label} must not contain duplicates")
    return result


def _timeline(
    value: object,
) -> typing.Dict[str, typing.Optional[str]]:
    raw = _object(value, "timeline")
    _keys(raw, _TIMELINE_KEYS, "timeline")
    result = {
        name: _optional_timestamp(raw[name], f"timeline.{name}")
        for name in sorted(raw)
    }
    for required in (
        "t1_candidate_saved_utc",
        "t2_snapshot_hashed_utc",
    ):
        if result[required] is None:
            raise CheckpointRenderError(
                f"timeline.{required} must be populated"
            )
    for pending in (
        "t3_first_pixels_written_utc",
        "t4_visual_verdict_complete_utc",
        "t5_goal_records_updated_utc",
    ):
        if result[pending] is not None:
            raise CheckpointRenderError(
                f"timeline.{pending} must still be null"
            )
    populated = [
        result[name]
        for name in (
            "t0_source_packet_frozen_utc",
            "t1_candidate_saved_utc",
            "t2_snapshot_hashed_utc",
        )
        if result[name] is not None
    ]
    if any(
        _timestamp_value(typing.cast(str, earlier))
        > _timestamp_value(typing.cast(str, later))
        for earlier, later in zip(populated, populated[1:])
    ):
        raise CheckpointRenderError(
            "timeline populated stages must be chronological"
        )
    return result


def _optional_timestamp(value: object, label: str) -> typing.Optional[str]:
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
        raise CheckpointRenderError(f"{label} must be an RFC 3339 timestamp")
    normalized = result[:-1] + "+00:00" if result.endswith("Z") else result
    try:
        parsed = datetime.datetime.fromisoformat(normalized)
    except ValueError as error:
        raise CheckpointRenderError(
            f"{label} must be an RFC 3339 timestamp"
        ) from error
    if parsed.tzinfo is None:
        raise CheckpointRenderError(f"{label} must include a UTC offset")
    return result


def _timestamp_value(value: str) -> datetime.datetime:
    normalized = value[:-1] + "+00:00" if value.endswith("Z") else value
    return datetime.datetime.fromisoformat(normalized)


def _timings(value: object) -> typing.Dict[str, float]:
    raw = _object(value, "timings_seconds")
    _keys(
        raw,
        {"candidate_save_operator", "candidate_hash_and_reverify"},
        "timings_seconds",
    )
    result: typing.Dict[str, float] = {}
    for raw_name in sorted(raw):
        name = _identifier(raw_name, "timing name")
        seconds = _finite_float(raw[raw_name], f"timings_seconds.{name}")
        if seconds < 0.0:
            raise CheckpointRenderError(
                f"timings_seconds.{name} must be nonnegative"
            )
        result[name] = seconds
    return result


def _validate_live_claims(
    workspace: pathlib.Path,
    request: _Request,
) -> None:
    before = request.live.context_before_snapshot
    after = request.live.context_after_snapshot
    if before != after:
        raise CheckpointRenderError(
            "live contexts before and after snapshot are not identical"
        )
    if before.active_scene != request.scene.name:
        raise CheckpointRenderError(
            "live active scene does not match the requested scene"
        )
    if before.active_view_layer not in request.scene.view_layers:
        raise CheckpointRenderError(
            "live active view layer is absent from the frozen scene"
        )
    if before.mode != "OBJECT":
        raise CheckpointRenderError(
            "live snapshot was not made in Object Mode"
        )
    object_names = {item.name for item in request.scene.objects}
    if (
        before.active_object is not None
        and before.active_object not in object_names
    ):
        raise CheckpointRenderError(
            "live active object is absent from the scene"
        )
    if not set(before.selected_objects) <= object_names:
        raise CheckpointRenderError(
            "live selected objects are absent from the frozen scene"
        )
    _validate_parent_attestation(workspace, request)
    render_version = _numeric_blender_version()
    if render_version != request.live.pinned_blender.version:
        raise CheckpointRenderError(
            "render Blender version does not match the declared pin: "
            f"expected {request.live.pinned_blender.version!r}, "
            f"got {render_version!r} ({bpy.app.version_string})"
        )
    live_version = _declared_live_version(request.live.blender_version)
    pinned_version = _semantic_version(
        request.live.pinned_blender.version,
        "live.pinned_blender.version",
    )
    if live_version > pinned_version:
        raise CheckpointRenderError(
            "live Blender version is newer than the declared renderer pin"
        )
    _validate_live_tool_sources(request)


def _validate_parent_attestation(
    workspace: pathlib.Path,
    request: _Request,
) -> None:
    source_path = pathlib.Path(
        request.live.context_before_snapshot.source_filepath
    )
    if not source_path.is_absolute():
        raise CheckpointRenderError(
            "live source_filepath must be absolute for exact parent proof"
        )
    parent_path = _resolve_artifact_path(workspace, request.parent, "parent")
    if source_path.resolve() != parent_path:
        raise CheckpointRenderError(
            "live source_filepath does not name the frozen parent"
        )
    attestation = request.live.parent_open_attestation
    attested = (
        attestation.path,
        attestation.sha256,
        attestation.spec_sha256,
        attestation.source_packet_sha256,
    )
    expected_attestation = (
        request.parent.path,
        request.parent.sha256,
        request.spec.sha256,
        request.source_packet.sha256,
    )
    if attested != expected_attestation:
        raise CheckpointRenderError(
            "parent-open attestation does not match the frozen primary inputs"
        )
    live_version = _declared_live_version(request.live.blender_version)
    if attestation.parent_file_version[:2] > live_version[:2]:
        raise CheckpointRenderError(
            "parent file major/minor version is newer than live Blender"
        )
    candidate_saved_utc = typing.cast(
        str, request.timeline["t1_candidate_saved_utc"]
    )
    if _timestamp_value(attestation.opened_utc) > _timestamp_value(
        candidate_saved_utc
    ):
        raise CheckpointRenderError(
            "parent-open attestation occurs after the candidate save"
        )


def _validate_live_tool_sources(request: _Request) -> None:
    expected_tools = {
        "projects/renders/cmd/fumo_review/live_checkpoint.py",
        "projects/renders/cmd/fumo_review/checkpoint_spec.py",
        "projects/renders/cmd/fumo_review/render_spec.py",
    }
    actual_tools = {item.path for item in request.live.tool_sources}
    if actual_tools != expected_tools:
        raise CheckpointRenderError(
            "live.tool_sources must identify the snapshot writer and its "
            "repository imports"
        )


def _validate_request_directory(
    output_dir: pathlib.Path,
    request_path: pathlib.Path,
    candidate_path: pathlib.Path,
) -> None:
    expected = {_REQUEST_NAME, _CANDIDATE_NAME}
    actual = {child.name for child in output_dir.iterdir()}
    if actual != expected:
        raise CheckpointRenderError(
            "unclaimed checkpoint directory must contain exactly "
            f"{sorted(expected)}, got {sorted(actual)}"
        )
    candidate_entry = output_dir / _CANDIDATE_NAME
    if request_path.is_symlink() or candidate_entry.is_symlink():
        raise CheckpointRenderError(
            "snapshot request and candidate must not be symlinks"
        )
    if candidate_entry.resolve() != candidate_path:
        raise CheckpointRenderError(
            "candidate directory entry changed identity"
        )
    if not request_path.is_file() or not candidate_entry.is_file():
        raise CheckpointRenderError(
            "snapshot request and candidate must both be regular files"
        )


def _claim_output(
    output_dir: pathlib.Path,
    request_sha256: str,
) -> typing.Tuple[pathlib.Path, bytes]:
    lock_path = output_dir / _LOCK_NAME
    payload = canonical_json(
        {
            "schema_version": 1,
            "state": "incomplete_until_READY_exists",
            "request_sha256": request_sha256,
            "render_blender_version": bpy.app.version_string,
            "pid": os.getpid(),
            "started_utc": _utc_now(),
        }
    ).encode("utf-8")
    flags = os.O_WRONLY | os.O_CREAT | os.O_EXCL
    if hasattr(os, "O_NOFOLLOW"):
        flags |= os.O_NOFOLLOW
    try:
        descriptor = os.open(lock_path, flags, 0o644)
    except FileExistsError as error:
        raise CheckpointRenderError(
            f"checkpoint is already claimed or incomplete: {lock_path}"
        ) from error
    try:
        _write_descriptor(descriptor, payload)
        os.fsync(descriptor)
    finally:
        os.close(descriptor)
    _fsync_directory(output_dir)
    return lock_path, payload


def _validate_claimed_directory(output_dir: pathlib.Path) -> None:
    expected = {_REQUEST_NAME, _CANDIDATE_NAME, _LOCK_NAME}
    actual = {child.name for child in output_dir.iterdir()}
    if actual != expected:
        raise CheckpointRenderError(
            "claimed checkpoint directory changed before render: "
            f"expected {sorted(expected)}, got {sorted(actual)}"
        )


def _write_descriptor(descriptor: int, payload: bytes) -> None:
    offset = 0
    while offset < len(payload):
        written = os.write(descriptor, payload[offset:])
        if written <= 0:
            raise OSError("short write")
        offset += written


def _build_run(
    workspace: pathlib.Path,
    request_path: pathlib.Path,
    request_sha256: str,
    request: _Request,
) -> _Run:
    output_dir = request_path.parent
    candidate_path = _resolve_artifact_path(
        workspace, request.candidate, "candidate"
    )
    expected_candidate = (output_dir / _CANDIDATE_NAME).resolve()
    if candidate_path != expected_candidate or request.candidate.path != (
        expected_candidate.relative_to(workspace).as_posix()
    ):
        raise CheckpointRenderError(
            "candidate artifact must name the request's sibling candidate.blend"
        )
    inputs = _resolve_inputs(workspace, request)
    _validate_pinned_toolchain(inputs, request.live.pinned_blender.version)
    spec_path = _path_for_role(inputs, "spec")
    spec, spec_hash = load_checkpoint_spec_with_hash(spec_path)
    if spec_hash != request.spec.sha256:
        raise CheckpointRenderError(
            "checkpoint spec hash changed while loading"
        )
    _crosscheck_spec(request, spec)
    tools = _renderer_sources(workspace)
    return _Run(
        workspace=workspace,
        output_dir=output_dir,
        request_path=request_path,
        request_sha256=request_sha256,
        request=request,
        spec=spec,
        inputs=inputs,
        tool_sources=tools,
    )


def _resolve_inputs(
    workspace: pathlib.Path,
    request: _Request,
) -> typing.Tuple[_FrozenInput, ...]:
    records: typing.List[typing.Tuple[str, _Artifact]] = [
        ("spec", request.spec),
        ("parent", request.parent),
        ("source_packet", request.source_packet),
        ("candidate", request.candidate),
    ]
    records.extend(
        (
            f"dependency:{item.kind}:{item.name}",
            _Artifact(path=item.path, sha256=item.sha256),
        )
        for item in request.dependencies
    )
    records.extend(
        (f"live_tool_source:{index}", item)
        for index, item in enumerate(request.live.tool_sources)
    )
    records.append(
        (
            "pinned_blender_toolchain",
            request.live.pinned_blender.toolchain,
        )
    )
    result = []
    for role, record in records:
        path = _resolve_artifact_path(workspace, record, role)
        actual = sha256_file(path)
        if actual != record.sha256:
            raise CheckpointRenderError(
                f"frozen {role} hash mismatch for {record.path}: "
                f"expected {record.sha256}, got {actual}"
            )
        result.append(_FrozenInput(role=role, record=record, path=path))
    return tuple(result)


def _validate_pinned_toolchain(
    inputs: typing.Iterable[_FrozenInput],
    expected_version: str,
) -> None:
    path = _path_for_role(inputs, "pinned_blender_toolchain")
    try:
        value = json.loads(
            path.read_text(encoding="utf-8"),
            object_pairs_hook=_unique_object,
            parse_constant=_reject_json_constant,
        )
    except (OSError, UnicodeError, json.JSONDecodeError) as error:
        raise CheckpointRenderError(
            f"cannot read pinned Blender toolchain declaration: {error}"
        ) from error
    root = _object(value, "pinned Blender toolchain")
    toolchains = root.get("toolchains")
    archives = root.get("archives")
    if not isinstance(toolchains, list) or not isinstance(archives, list):
        raise CheckpointRenderError(
            "pinned Blender toolchain has no toolchains or archives array"
        )
    declared_versions: typing.Set[str] = set()
    for index, child in enumerate(toolchains):
        item = _object(child, f"pinned Blender toolchains[{index}]")
        if item.get("name") == "blender" and isinstance(
            item.get("version"), str
        ):
            declared_versions.add(typing.cast(str, item["version"]))
    archive_versions: typing.Set[str] = set()
    for index, child in enumerate(archives):
        item = _object(child, f"pinned Blender archives[{index}]")
        if item.get("toolchain_name") == "blender" and isinstance(
            item.get("toolchain_version"), str
        ):
            archive_versions.add(typing.cast(str, item["toolchain_version"]))
    if declared_versions != {expected_version} or archive_versions != {
        expected_version
    }:
        raise CheckpointRenderError(
            "pinned Blender version disagrees with its toolchain declaration: "
            f"request={expected_version!r}, "
            f"toolchains={sorted(declared_versions)!r}, "
            f"archives={sorted(archive_versions)!r}"
        )


def _resolve_artifact_path(
    workspace: pathlib.Path,
    artifact: _Artifact,
    label: str,
) -> pathlib.Path:
    unresolved = workspace.joinpath(
        *pathlib.PurePosixPath(artifact.path).parts
    )
    return resolve_workspace_file(workspace, unresolved, f"frozen {label}")


def _path_for_role(
    inputs: typing.Iterable[_FrozenInput],
    role: str,
) -> pathlib.Path:
    for item in inputs:
        if item.role == role:
            return item.path
    raise CheckpointRenderError(f"internal error: missing frozen input {role}")


def _crosscheck_spec(request: _Request, spec: CheckpointSpec) -> None:
    mismatches: typing.List[str] = []
    _compare(mismatches, "checkpoint", request.checkpoint, spec.checkpoint)
    _compare(mismatches, "scene", request.scene.name, spec.scene)
    _compare(mismatches, "width", request.width, spec.resolution.width)
    _compare(mismatches, "height", request.height, spec.resolution.height)
    _compare(mismatches, "parent.path", request.parent.path, spec.parent.path)
    _compare(
        mismatches, "parent.sha256", request.parent.sha256, spec.parent.sha256
    )
    _compare(
        mismatches,
        "source_packet.path",
        request.source_packet.path,
        spec.source_packet.path,
    )
    _compare(
        mismatches,
        "source_packet.sha256",
        request.source_packet.sha256,
        spec.source_packet.sha256,
    )
    _compare(mismatches, "hypothesis", request.hypothesis, spec.hypothesis)
    _compare(mismatches, "vetoes", request.vetoes, tuple(spec.vetoes))
    _compare(
        mismatches,
        "boolean requirements",
        dict(request.scene.boolean_requirements),
        dict(spec.scene_boolean_requirements),
    )
    request_views = {
        name: camera.name for name, camera in request.scene.views.items()
    }
    spec_views = {name: view.camera for name, view in spec.views.items()}
    _compare(mismatches, "views", request_views, spec_views)
    _compare(
        mismatches,
        "T0 source packet time",
        request.timeline["t0_source_packet_frozen_utc"],
        spec.t0_source_packet_frozen_utc,
    )
    if mismatches:
        details = "\n".join(f"  - {item}" for item in mismatches)
        raise CheckpointRenderError(
            f"snapshot request disagrees with checkpoint spec:\n{details}"
        )


def _compare(
    mismatches: typing.List[str],
    label: str,
    request_value: object,
    spec_value: object,
) -> None:
    if request_value != spec_value:
        mismatches.append(
            f"{label}: request={request_value!r}, spec={spec_value!r}"
        )


def _renderer_sources(
    workspace: pathlib.Path,
) -> typing.Tuple[_RendererSource, ...]:
    module_paths = {
        "checkpoint_render": pathlib.Path(__file__).resolve(),
        "checkpoint_spec": _module_path(
            "projects.renders.cmd.fumo_review.checkpoint_spec"
        ),
        "png": _module_path("projects.renders.cmd.fumo_review.png"),
        "render_packet": _module_path(
            "projects.renders.cmd.fumo_review.render_packet"
        ),
        "render_spec": _module_path(
            "projects.renders.cmd.fumo_review.render_spec"
        ),
    }
    result: typing.List[_RendererSource] = []
    for name in sorted(module_paths):
        repository_path = workspace.joinpath(*_TOOL_SOURCE_PATHS[name].parts)
        if not repository_path.is_file():
            raise CheckpointRenderError(
                f"renderer source is missing from workspace: {repository_path}"
            )
        runtime_path = module_paths[name]
        loaded_hash = _LOADED_TOOL_HASHES[name]
        repository_hash = sha256_file(repository_path)
        runtime_hash = sha256_file(runtime_path)
        if runtime_hash != loaded_hash or repository_hash != loaded_hash:
            raise CheckpointRenderError(
                f"loaded {name} bytes differ from current repository source"
            )
        result.append(
            _RendererSource(
                module=name,
                path=repository_path.resolve(),
                runtime_path=runtime_path,
                sha256=loaded_hash,
            )
        )
    return tuple(result)


def _module_path(module_name: str) -> pathlib.Path:
    module = importlib.import_module(module_name)
    path = getattr(module, "__file__", None)
    if not isinstance(path, str):
        raise CheckpointRenderError(f"module {module_name} has no source path")
    result = pathlib.Path(path).resolve()
    if result.suffix == ".pyc" and result.with_suffix(".py").is_file():
        result = result.with_suffix(".py")
    if not result.is_file():
        raise CheckpointRenderError(
            f"module {module_name} source is not a file: {result}"
        )
    return result


def _verify_frozen(run: _Run, stage: str) -> None:
    actual_request = sha256_file(run.request_path)
    if actual_request != run.request_sha256:
        raise CheckpointRenderError(
            f"snapshot request hash changed at {stage}: "
            f"expected {run.request_sha256}, got {actual_request}"
        )
    for item in run.inputs:
        current_path = _resolve_artifact_path(
            run.workspace, item.record, item.role
        )
        if current_path != item.path:
            raise CheckpointRenderError(
                f"frozen {item.role} changed path identity at {stage}"
            )
        actual = sha256_file(current_path)
        if actual != item.record.sha256:
            raise CheckpointRenderError(
                f"frozen {item.role} hash changed at {stage}: "
                f"expected {item.record.sha256}, got {actual}"
            )
    _verify_renderer_sources(run, stage)


def _verify_renderer_sources(run: _Run, stage: str) -> None:
    for source in run.tool_sources:
        repository_hash = sha256_file(source.path)
        runtime_hash = sha256_file(source.runtime_path)
        if repository_hash != source.sha256 or runtime_hash != source.sha256:
            raise CheckpointRenderError(
                f"renderer source {source.module} changed at {stage}"
            )


def _open_candidate(path: pathlib.Path) -> typing.Dict[str, int]:
    handler_counts: typing.Dict[str, int] = {}
    with _suspended_load_handlers(handler_counts):
        result = bpy.ops.wm.open_mainfile(
            filepath=str(path),
            load_ui=False,
            use_scripts=False,
        )
    if "FINISHED" not in result:
        raise CheckpointRenderError(
            f"pinned Blender could not open candidate: {path}"
        )
    if not bpy.data.filepath:
        raise CheckpointRenderError("opened candidate has no Blender filepath")
    opened_path = pathlib.Path(bpy.path.abspath(bpy.data.filepath)).resolve()
    if opened_path != path:
        raise CheckpointRenderError(
            f"Blender opened an unexpected candidate: {opened_path}"
        )
    candidate_version = _blender_file_version()
    render_version = tuple(int(part) for part in bpy.app.version[:3])
    if candidate_version[:2] > render_version[:2]:
        raise CheckpointRenderError(
            "candidate file major/minor version is newer than render Blender"
        )
    return handler_counts


@contextlib.contextmanager
def _suspended_load_handlers(
    counts: typing.MutableMapping[str, int],
) -> typing.Iterator[None]:
    names = ("load_pre", "load_post", "load_post_fail", "version_update")
    captured: typing.List[
        typing.Tuple[typing.Any, typing.List[typing.Any]]
    ] = []
    try:
        for name in names:
            handlers = getattr(bpy.app.handlers, name, None)
            if handlers is None:
                continue
            original = list(handlers)
            counts[name] = len(original)
            captured.append((handlers, original))
            handlers.clear()
        yield
    finally:
        for handlers, original in captured:
            handlers.clear()
            handlers.extend(original)


def _validated_scene(run: _Run) -> typing.Tuple[typing.Any, typing.Any]:
    request = run.request
    scene = bpy.data.scenes.get(request.scene.name)
    if scene is None:
        available = sorted(item.name for item in bpy.data.scenes)
        raise CheckpointRenderError(
            f"scene {request.scene.name!r} is missing; available: {available}"
        )
    errors: typing.List[str] = []
    _validate_open_context(request, errors)
    _validate_boolean_requirements(scene, request, errors)
    _validate_scene_inventory(scene, request, errors)
    _validate_cameras(scene, request, errors)
    if errors:
        details = "\n".join(f"  - {error}" for error in errors)
        raise CheckpointRenderError(
            f"opened candidate failed snapshot preconditions:\n{details}"
        )
    view_layer = scene.view_layers.get(
        request.live.context_before_snapshot.active_view_layer
    )
    if view_layer is None:
        raise CheckpointRenderError("frozen active view layer is missing")
    if hasattr(view_layer, "use") and not view_layer.use:
        raise CheckpointRenderError("frozen active view layer is disabled")
    layer_objects = list(view_layer.objects)
    renderable = [
        obj
        for obj in layer_objects
        if obj.type in _GEOMETRY_TYPES and not obj.hide_render
    ]
    if not renderable:
        raise CheckpointRenderError(
            "active view layer has no renderable geometry"
        )
    for camera in request.scene.views.values():
        if scene.objects[camera.name] not in layer_objects:
            raise CheckpointRenderError(
                f"camera {camera.name!r} is excluded from active view layer"
            )
    _validate_open_dependencies(run)
    return scene, view_layer


def _validate_open_context(
    request: _Request,
    errors: typing.List[str],
) -> None:
    expected = request.live.context_before_snapshot
    actual = _current_context()
    comparisons = (
        ("active scene", actual.active_scene, expected.active_scene),
        (
            "active view layer",
            actual.active_view_layer,
            expected.active_view_layer,
        ),
        ("active object", actual.active_object, expected.active_object),
        (
            "selected objects",
            actual.selected_objects,
            expected.selected_objects,
        ),
        ("mode", actual.mode, expected.mode),
    )
    for label, actual_value, expected_value in comparisons:
        if actual_value != expected_value:
            errors.append(
                f"{label} changed across save/reopen: "
                f"expected {expected_value!r}, got {actual_value!r}"
            )


def _current_context() -> _Context:
    active = bpy.context.view_layer.objects.active
    return _Context(
        active_scene=bpy.context.scene.name,
        active_view_layer=bpy.context.view_layer.name,
        active_object=active.name if active else None,
        selected_objects=tuple(
            sorted(obj.name for obj in bpy.context.selected_objects)
        ),
        mode=bpy.context.mode,
        source_filepath=bpy.data.filepath,
    )


def _validate_boolean_requirements(
    scene: typing.Any,
    request: _Request,
    errors: typing.List[str],
) -> None:
    for name, expected in request.scene.boolean_requirements.items():
        if name not in scene:
            errors.append(f"required scene boolean {name!r} is missing")
            continue
        actual = scene[name]
        if not isinstance(actual, bool) or actual is not expected:
            errors.append(
                f"required scene boolean {name!r} must be {expected}, "
                f"got {actual!r}"
            )


def _validate_scene_inventory(
    scene: typing.Any,
    request: _Request,
    errors: typing.List[str],
) -> None:
    layers = tuple(sorted(layer.name for layer in scene.view_layers))
    if layers != request.scene.view_layers:
        errors.append(
            f"view layers changed: expected {request.scene.view_layers!r}, "
            f"got {layers!r}"
        )
    objects = tuple(
        _object_record(obj)
        for obj in sorted(scene.objects, key=lambda item: item.name)
    )
    if objects != request.scene.objects:
        errors.append("scene object inventory changed across save/reopen")


def _object_record(obj: typing.Any) -> _Object:
    data = getattr(obj, "data", None)
    return _Object(
        name=obj.name,
        type=obj.type,
        data=data.name if data else None,
        hide_render=bool(obj.hide_render),
    )


def _validate_cameras(
    scene: typing.Any,
    request: _Request,
    errors: typing.List[str],
) -> None:
    for view_name, expected in request.scene.views.items():
        camera = scene.objects.get(expected.name)
        if camera is None:
            errors.append(
                f"view {view_name!r} camera {expected.name!r} is missing"
            )
            continue
        if camera.type != "CAMERA":
            errors.append(
                f"view {view_name!r} names a {camera.type}, not a CAMERA"
            )
            continue
        actual = _camera_record(camera)
        difference = _camera_difference(expected, actual)
        if difference:
            errors.append(f"view {view_name!r} camera changed: {difference}")


def _camera_record(camera: typing.Any) -> _Camera:
    data = camera.data
    return _Camera(
        name=camera.name,
        projection=data.type,
        lens=float(data.lens),
        ortho_scale=float(data.ortho_scale),
        sensor_fit=data.sensor_fit,
        sensor_width=float(data.sensor_width),
        sensor_height=float(data.sensor_height),
        shift_x=float(data.shift_x),
        shift_y=float(data.shift_y),
        clip_start=float(data.clip_start),
        clip_end=float(data.clip_end),
        matrix_world=tuple(
            tuple(float(value) for value in row) for row in camera.matrix_world
        ),
    )


def _camera_difference(expected: _Camera, actual: _Camera) -> str:
    text_fields = ("name", "projection", "sensor_fit")
    for field in text_fields:
        if getattr(expected, field) != getattr(actual, field):
            return (
                f"{field} expected {getattr(expected, field)!r}, "
                f"got {getattr(actual, field)!r}"
            )
    float_fields = (
        "lens",
        "ortho_scale",
        "sensor_width",
        "sensor_height",
        "shift_x",
        "shift_y",
        "clip_start",
        "clip_end",
    )
    for field in float_fields:
        expected_value = getattr(expected, field)
        actual_value = getattr(actual, field)
        if not math.isclose(
            expected_value, actual_value, rel_tol=1e-7, abs_tol=1e-7
        ):
            return f"{field} expected {expected_value}, got {actual_value}"
    for row in range(4):
        for column in range(4):
            expected_value = expected.matrix_world[row][column]
            actual_value = actual.matrix_world[row][column]
            if not math.isclose(
                expected_value, actual_value, rel_tol=1e-7, abs_tol=1e-7
            ):
                return (
                    f"matrix_world[{row}][{column}] expected "
                    f"{expected_value}, got {actual_value}"
                )
    return ""


def _validate_open_dependencies(run: _Run) -> None:
    expected = tuple(run.request.dependencies)
    actual = _opened_dependencies(run.workspace)
    if actual != expected:
        raise CheckpointRenderError(
            "candidate external dependency inventory changed across save/reopen: "
            f"expected {_dependency_keys(expected)!r}, "
            f"got {_dependency_keys(actual)!r}"
        )


def _dependency_keys(
    dependencies: typing.Iterable[_Dependency],
) -> typing.List[typing.Tuple[str, str, str, str]]:
    return [
        (item.kind, item.name, item.path, item.sha256) for item in dependencies
    ]


def _opened_dependencies(
    workspace: pathlib.Path,
) -> typing.Tuple[_Dependency, ...]:
    records: typing.List[_Dependency] = []
    path_map = bpy.data.file_path_map(include_libraries=True)
    for datablock, raw_paths in path_map.items():
        dirty = getattr(datablock, "is_dirty", False)
        modified = getattr(datablock, "is_modified", False)
        if (isinstance(dirty, bool) and dirty) or (
            isinstance(modified, bool) and modified
        ):
            raise CheckpointRenderError(
                f"opened external datablock {datablock.name!r} has "
                "unsaved in-memory changes"
            )
        source = getattr(datablock, "source", "")
        is_sequence = getattr(datablock, "is_sequence", False)
        if source in {"SEQUENCE", "TILED"} or (
            isinstance(is_sequence, bool) and is_sequence
        ):
            raise CheckpointRenderError(
                f"opened external datablock {datablock.name!r} uses an "
                "unbounded sequence or tiled source"
            )
        if not raw_paths:
            continue
        kind = str(getattr(datablock, "id_type", datablock.bl_rna.identifier))
        for raw in sorted(raw_paths):
            if not raw or raw == "<builtin>":
                continue
            if any(token in raw for token in ("#", "<UDIM>", "<UVTILE>")):
                raise CheckpointRenderError(
                    f"opened dependency {datablock.name!r} uses an "
                    f"unbounded sequence or tile path: {raw}"
                )
            path = _dependency_path(workspace, kind, datablock, raw)
            records.append(
                _Dependency(
                    kind=kind,
                    name=datablock.name,
                    path=path.relative_to(workspace).as_posix(),
                    sha256=sha256_file(path),
                )
            )
    unique: typing.Dict[typing.Tuple[str, str, str], _Dependency] = {}
    for record in records:
        key = (record.kind, record.name, record.path)
        previous = unique.get(key)
        if previous is not None and previous.sha256 != record.sha256:
            raise CheckpointRenderError(
                "opened external dependency changed while hashing: "
                f"{record.path}"
            )
        unique[key] = record
    return tuple(unique[key] for key in sorted(unique))


def _dependency_path(
    workspace: pathlib.Path,
    kind: str,
    datablock: typing.Any,
    raw: str,
) -> pathlib.Path:
    try:
        absolute = pathlib.Path(
            bpy.path.abspath(
                raw,
                library=getattr(datablock, "library", None),
            )
        )
    except TypeError:
        absolute = pathlib.Path(bpy.path.abspath(raw))
    return resolve_workspace_file(
        workspace,
        absolute,
        f"opened {kind} dependency {datablock.name!r}",
    )


def _configure_preview(
    scene: typing.Any,
    view_layer: typing.Any,
    request: _Request,
) -> typing.Dict[str, object]:
    engine = _select_eevee(scene)
    render = scene.render
    render.resolution_x = request.width
    render.resolution_y = request.height
    render.resolution_percentage = 100
    render.pixel_aspect_x = 1.0
    render.pixel_aspect_y = 1.0
    render.use_file_extension = True
    render.use_overwrite = True
    render.use_placeholder = False
    render.use_border = False
    render.use_crop_to_border = False
    render.film_transparent = False
    render.use_compositing = False
    render.use_sequencer = False
    render.use_stamp = False
    if hasattr(render, "use_multiview"):
        render.use_multiview = False
    if hasattr(render, "use_freestyle"):
        render.use_freestyle = False
    if hasattr(render, "dither_intensity"):
        render.dither_intensity = 0.0
    if hasattr(render, "use_motion_blur"):
        render.use_motion_blur = False
    if hasattr(render, "threads_mode"):
        render.threads_mode = "FIXED"
        render.threads = 1
    for prop in render.bl_rna.properties:
        if prop.identifier.startswith("use_stamp_"):
            setattr(render, prop.identifier, False)
    image_settings = render.image_settings
    image_settings.file_format = "PNG"
    image_settings.color_depth = "8"
    image_settings.compression = 15
    sample_settings = _bound_eevee_samples(scene, view_layer)
    return {
        "engine": engine,
        "maximum_dimension": _MAX_DIMENSION,
        "maximum_views": 2,
        "pixel_aspect": [1.0, 1.0],
        "compositing": False,
        "sequencer": False,
        "film_transparent": False,
        "dither_intensity": 0.0,
        "motion_blur": False,
        "multiview": False,
        "sample_settings": sample_settings,
        "beauty_color_management": {
            "display_device": scene.display_settings.display_device,
            "view_transform": scene.view_settings.view_transform,
            "look": scene.view_settings.look,
            "exposure": float(scene.view_settings.exposure),
            "gamma": float(scene.view_settings.gamma),
        },
        "silhouette_color_management": {
            "display_device": "sRGB",
            "view_transform": "Standard",
            "look": "None",
            "exposure": 0.0,
            "gamma": 1.0,
            "dither_intensity": 0.0,
        },
        "silhouette_world_override": {
            "forced": True,
            "use_sky": True,
            "color": [1.0, 1.0, 1.0, 1.0],
            "strength": 1.0,
        },
    }


def _select_eevee(scene: typing.Any) -> str:
    for engine in ("BLENDER_EEVEE_NEXT", "BLENDER_EEVEE"):
        try:
            scene.render.engine = engine
        except (TypeError, ValueError):
            continue
        if scene.render.engine == engine:
            return engine
    raise CheckpointRenderError("pinned Blender has no supported Eevee engine")


def _bound_eevee_samples(
    scene: typing.Any,
    view_layer: typing.Any,
) -> typing.Dict[str, object]:
    # Eevee uses the per-view-layer override when it exists. Force it instead
    # of trusting a candidate-provided value to keep the preview bounded.
    view_layer.samples = _PREVIEW_SAMPLES
    result: typing.Dict[str, object] = {
        "view_layer_samples": int(view_layer.samples)
    }
    owner = getattr(scene, "eevee", None)
    if owner is None:
        return result
    limits = {
        "taa_render_samples": _PREVIEW_SAMPLES,
        "taa_samples": _PREVIEW_SAMPLES,
        "volumetric_samples": _PREVIEW_SAMPLES,
        "volumetric_shadow_samples": 16,
        "shadow_step_count": 16,
    }
    for name, limit in limits.items():
        if not hasattr(owner, name):
            continue
        current = int(getattr(owner, name))
        bounded = max(1, min(current, limit))
        setattr(owner, name, bounded)
        result[name] = int(getattr(owner, name))
    for name in ("use_raytracing", "use_bokeh_jittered"):
        if hasattr(owner, name):
            setattr(owner, name, False)
            result[name] = bool(getattr(owner, name))
    return result


def _render_outputs(
    run: _Run,
    scene: typing.Any,
    view_layer: typing.Any,
    timeline: typing.MutableMapping[str, typing.Optional[str]],
) -> typing.Tuple[
    typing.List[typing.Dict[str, object]],
    typing.Dict[str, float],
    typing.Dict[str, object],
]:
    beauty_dir = run.output_dir / "beauty"
    silhouette_dir = run.output_dir / "silhouette"
    beauty_dir.mkdir()
    silhouette_dir.mkdir()
    preview = _configure_preview(scene, view_layer, run.request)
    outputs: typing.List[typing.Dict[str, object]] = []
    timings: typing.Dict[str, float] = {}

    beauty_started = time.perf_counter()
    scene.render.image_settings.color_mode = "RGBA"
    for name, camera in run.request.scene.views.items():
        output, written_utc = _render_view(
            scene,
            view_layer,
            camera.name,
            beauty_dir / f"{name}.png",
            "beauty",
            name,
            run.request.width,
            run.request.height,
        )
        outputs.append(output)
        if timeline["t3_first_pixels_written_utc"] is None:
            timeline["t3_first_pixels_written_utc"] = written_utc
    timings["beauty_renders"] = time.perf_counter() - beauty_started

    silhouette_started = time.perf_counter()
    with _silhouette_override(scene, view_layer):
        scene.render.image_settings.color_mode = "RGB"
        for name, camera in run.request.scene.views.items():
            output, _ = _render_view(
                scene,
                view_layer,
                camera.name,
                silhouette_dir / f"{name}.png",
                "silhouette",
                name,
                run.request.width,
                run.request.height,
            )
            outputs.append(output)
    timings["silhouette_renders"] = time.perf_counter() - silhouette_started
    return outputs, timings, preview


def _render_view(
    scene: typing.Any,
    view_layer: typing.Any,
    camera_name: str,
    path: pathlib.Path,
    kind: str,
    view_name: str,
    width: int,
    height: int,
) -> typing.Tuple[typing.Dict[str, object], str]:
    scene.camera = scene.objects[camera_name]
    scene.render.filepath = str(path)
    started = time.perf_counter()
    result = bpy.ops.render.render(
        write_still=True,
        scene=scene.name,
        layer=view_layer.name,
    )
    seconds = time.perf_counter() - started
    if "FINISHED" not in result or not path.is_file():
        raise CheckpointRenderError(
            f"{kind} render from {camera_name!r} did not write {path}"
        )
    written_utc = _utc_now()
    _canonicalize_png(path)
    _fsync_file(path)
    _fsync_directory(path.parent)
    image = decode_png(path)
    foreground = _validate_image(image, kind, path, width, height)
    return (
        {
            "kind": kind,
            "name": view_name,
            "camera": camera_name,
            "view_layer": view_layer.name,
            "path": path.relative_to(path.parent.parent).as_posix(),
            "sha256": sha256_file(path),
            "bytes": path.stat().st_size,
            "width": image.width,
            "height": image.height,
            "seconds": round(seconds, 6),
            "foreground_validation": foreground,
        },
        written_utc,
    )


def _validate_image(
    image: Image,
    kind: str,
    path: pathlib.Path,
    width: int,
    height: int,
) -> typing.Dict[str, object]:
    if image.width != width or image.height != height:
        raise CheckpointRenderError(
            f"wrong PNG dimensions for {path}: "
            f"expected {width}x{height}, got {image.width}x{image.height}"
        )
    if len(image.rgb) != width * height * 3:
        raise CheckpointRenderError(f"invalid RGB payload length for {path}")
    minimum, maximum = _channel_extrema(image.rgb)
    span = max(maximum[index] - minimum[index] for index in range(3))
    threshold = max(16, width * height // 10000)
    first = image.rgb[0:3]
    changed = sum(
        max(abs(pixel[index] - first[index]) for index in range(3)) >= 8
        for pixel in _pixels(image.rgb)
    )
    if span < 8 or changed < threshold:
        raise CheckpointRenderError(
            f"{kind} output has no informative foreground: {path}"
        )
    result: typing.Dict[str, object] = {
        "channel_min": list(minimum),
        "channel_max": list(maximum),
        "pixels_different_from_first": changed,
        "minimum_informative_pixels": threshold,
    }
    if kind == "silhouette":
        result.update(_validate_silhouette(image, path, threshold))
    return result


def _channel_extrema(
    rgb: bytes,
) -> typing.Tuple[typing.Tuple[int, int, int], typing.Tuple[int, int, int]]:
    minimum = [255, 255, 255]
    maximum = [0, 0, 0]
    for pixel in _pixels(rgb):
        for channel in range(3):
            minimum[channel] = min(minimum[channel], pixel[channel])
            maximum[channel] = max(maximum[channel], pixel[channel])
    return (
        (minimum[0], minimum[1], minimum[2]),
        (maximum[0], maximum[1], maximum[2]),
    )


def _pixels(rgb: bytes) -> typing.Iterator[typing.Tuple[int, int, int]]:
    for offset in range(0, len(rgb), 3):
        yield (rgb[offset], rgb[offset + 1], rgb[offset + 2])


def _validate_silhouette(
    image: Image,
    path: pathlib.Path,
    threshold: int,
) -> typing.Dict[str, object]:
    stats = _silhouette_stats(image)
    pixel_count = image.width * image.height
    minimum_dark = max(threshold, 32, pixel_count // 1000)
    minimum_light = max(threshold, pixel_count // 20)
    if len(stats.dark_indices) < minimum_dark or stats.light < minimum_light:
        raise CheckpointRenderError(
            f"silhouette lacks black foreground or white background: {path}"
        )
    border_pixels = 2 * image.width + 2 * image.height - 4
    if stats.border_light * 5 < border_pixels * 4:
        raise CheckpointRenderError(
            f"silhouette does not have a predominantly white border: {path}"
        )
    corner_light = sum(
        min(image.pixel(x, y)) >= 223
        for x, y in (
            (0, 0),
            (image.width - 1, 0),
            (0, image.height - 1),
            (image.width - 1, image.height - 1),
        )
    )
    if corner_light < 3:
        raise CheckpointRenderError(
            f"silhouette background polarity is not white: {path}"
        )
    if stats.bounds == (0, 0, image.width, image.height):
        raise CheckpointRenderError(
            f"silhouette dark foreground spans the entire frame: {path}"
        )
    if (
        stats.bounds[2] - stats.bounds[0] < 2
        or stats.bounds[3] - stats.bounds[1] < 2
    ):
        raise CheckpointRenderError(
            f"silhouette foreground is not spatially informative: {path}"
        )
    largest_component = _largest_dark_component(
        stats.dark_mask,
        stats.dark_indices,
        image.width,
        image.height,
    )
    if largest_component < minimum_dark:
        raise CheckpointRenderError(
            f"silhouette has no materially connected foreground: {path}"
        )
    if stats.chromatic > max(4, pixel_count // 1000):
        raise CheckpointRenderError(
            f"silhouette is not neutral black-on-white: {path}"
        )
    return {
        "black_pixels": len(stats.dark_indices),
        "white_pixels": stats.light,
        "black_fraction": round(len(stats.dark_indices) / pixel_count, 8),
        "white_fraction": round(stats.light / pixel_count, 8),
        "white_border_pixels": stats.border_light,
        "border_pixels": border_pixels,
        "white_corners": corner_light,
        "black_bbox": list(stats.bounds),
        "largest_black_component": largest_component,
        "minimum_black_component": minimum_dark,
        "chromatic_pixels": stats.chromatic,
    }


def _silhouette_stats(image: Image) -> _SilhouetteStats:
    dark_mask = bytearray(image.width * image.height)
    dark_indices: typing.List[int] = []
    light = 0
    border_light = 0
    chromatic = 0
    bounds = [image.width, image.height, -1, -1]
    for y in range(image.height):
        for x in range(image.width):
            pixel = image.pixel(x, y)
            if max(pixel) <= 32:
                index = y * image.width + x
                dark_mask[index] = 1
                dark_indices.append(index)
                bounds[0] = min(bounds[0], x)
                bounds[1] = min(bounds[1], y)
                bounds[2] = max(bounds[2], x + 1)
                bounds[3] = max(bounds[3], y + 1)
            if min(pixel) >= 223:
                light += 1
                if x in {0, image.width - 1} or y in {
                    0,
                    image.height - 1,
                }:
                    border_light += 1
            if max(pixel) - min(pixel) > 2:
                chromatic += 1
    return _SilhouetteStats(
        dark_mask=dark_mask,
        dark_indices=tuple(dark_indices),
        light=light,
        border_light=border_light,
        chromatic=chromatic,
        bounds=(bounds[0], bounds[1], bounds[2], bounds[3]),
    )


def _largest_dark_component(
    dark_mask: bytearray,
    dark_indices: typing.Iterable[int],
    width: int,
    height: int,
) -> int:
    visited = bytearray(len(dark_mask))
    largest = 0
    for start in dark_indices:
        if visited[start]:
            continue
        visited[start] = 1
        stack = [start]
        size = 0
        while stack:
            index = stack.pop()
            size += 1
            x = index % width
            neighbors = []
            if x > 0:
                neighbors.append(index - 1)
            if x + 1 < width:
                neighbors.append(index + 1)
            if index >= width:
                neighbors.append(index - width)
            if index + width < width * height:
                neighbors.append(index + width)
            for neighbor in neighbors:
                if dark_mask[neighbor] and not visited[neighbor]:
                    visited[neighbor] = 1
                    stack.append(neighbor)
        largest = max(largest, size)
    return largest


@contextlib.contextmanager
def _silhouette_override(
    scene: typing.Any,
    view_layer: typing.Any,
) -> typing.Iterator[None]:
    material = _black_material()
    world = _white_world()
    original_world = scene.world
    original_override = view_layer.material_override
    original_world_override = view_layer.world_override
    original_use_sky = view_layer.use_sky
    original_display_device = scene.display_settings.display_device
    original_view_transform = scene.view_settings.view_transform
    original_look = scene.view_settings.look
    original_exposure = scene.view_settings.exposure
    original_gamma = scene.view_settings.gamma
    original_curve_mapping = scene.view_settings.use_curve_mapping
    try:
        scene.world = world
        view_layer.world_override = world
        view_layer.use_sky = True
        scene.render.film_transparent = False
        scene.render.use_compositing = False
        scene.render.use_sequencer = False
        if hasattr(scene.render, "dither_intensity"):
            scene.render.dither_intensity = 0.0
        scene.display_settings.display_device = "sRGB"
        scene.view_settings.view_transform = "Standard"
        scene.view_settings.look = "None"
        scene.view_settings.exposure = 0.0
        scene.view_settings.gamma = 1.0
        scene.view_settings.use_curve_mapping = False
        view_layer.material_override = material
        yield
    finally:
        view_layer.material_override = original_override
        view_layer.world_override = original_world_override
        view_layer.use_sky = original_use_sky
        scene.world = original_world
        scene.display_settings.display_device = original_display_device
        scene.view_settings.view_transform = original_view_transform
        scene.view_settings.look = original_look
        scene.view_settings.exposure = original_exposure
        scene.view_settings.gamma = original_gamma
        scene.view_settings.use_curve_mapping = original_curve_mapping
        bpy.data.worlds.remove(world)
        bpy.data.materials.remove(material)


def _black_material() -> typing.Any:
    material = bpy.data.materials.new(name=".fumo_checkpoint_black")
    material.use_nodes = True
    nodes = material.node_tree.nodes
    nodes.clear()
    output = nodes.new("ShaderNodeOutputMaterial")
    emission = nodes.new("ShaderNodeEmission")
    emission.inputs["Color"].default_value = (0.0, 0.0, 0.0, 1.0)
    emission.inputs["Strength"].default_value = 1.0
    material.node_tree.links.new(
        emission.outputs["Emission"], output.inputs["Surface"]
    )
    return material


def _white_world() -> typing.Any:
    world = bpy.data.worlds.new(name=".fumo_checkpoint_white")
    world.use_nodes = True
    background = world.node_tree.nodes.get("Background")
    if background is None:
        raise CheckpointRenderError("Blender did not create a Background node")
    background.inputs["Color"].default_value = (1.0, 1.0, 1.0, 1.0)
    background.inputs["Strength"].default_value = 1.0
    return world


def _camera_manifest(
    scene: typing.Any,
    view_layer: typing.Any,
    request: _Request,
) -> typing.Dict[str, object]:
    with bpy.context.temp_override(scene=scene, view_layer=view_layer):
        depsgraph = bpy.context.evaluated_depsgraph_get()
    result: typing.Dict[str, object] = {}
    for view_name, expected in request.scene.views.items():
        camera = scene.objects[expected.name]
        snapshot = _camera_record(camera)
        evaluated_camera = camera.evaluated_get(depsgraph)
        actual = _camera_record(evaluated_camera)
        projection = evaluated_camera.calc_matrix_camera(
            depsgraph,
            x=request.width,
            y=request.height,
            scale_x=1.0,
            scale_y=1.0,
        )
        result[view_name] = {
            "name": actual.name,
            "projection": actual.projection,
            "lens": actual.lens,
            "ortho_scale": actual.ortho_scale,
            "sensor_fit": actual.sensor_fit,
            "sensor_width": actual.sensor_width,
            "sensor_height": actual.sensor_height,
            "shift_x": actual.shift_x,
            "shift_y": actual.shift_y,
            "clip_start": actual.clip_start,
            "clip_end": actual.clip_end,
            "matrix_world": [list(row) for row in actual.matrix_world],
            "projection_matrix": [
                [float(value) for value in row] for row in projection
            ],
            "snapshot_record": {
                "name": snapshot.name,
                "projection": snapshot.projection,
                "lens": snapshot.lens,
                "ortho_scale": snapshot.ortho_scale,
                "sensor_fit": snapshot.sensor_fit,
                "sensor_width": snapshot.sensor_width,
                "sensor_height": snapshot.sensor_height,
                "shift_x": snapshot.shift_x,
                "shift_y": snapshot.shift_y,
                "clip_start": snapshot.clip_start,
                "clip_end": snapshot.clip_end,
                "matrix_world": [list(row) for row in snapshot.matrix_world],
            },
            "evaluated_for_render": True,
        }
    return result


def _manifest(
    run: _Run,
    scene: typing.Any,
    view_layer: typing.Any,
    outputs: typing.List[typing.Dict[str, object]],
    preview: typing.Mapping[str, object],
    timeline: typing.Mapping[str, typing.Optional[str]],
    timings: typing.Mapping[str, float],
    renderer_load_handlers: typing.Mapping[str, int],
    verdict_hash: str,
    completed_utc: str,
) -> typing.Dict[str, object]:
    return {
        "schema_version": 1,
        "checkpoint": run.request.checkpoint,
        "completion_marker": _READY_NAME,
        "request": {
            "path": _relative(run.workspace, run.request_path),
            "sha256": run.request_sha256,
            "hash_verified_before_render": True,
            "hash_verified_after_render": True,
        },
        "inputs": {
            "spec": _input_manifest(run, "spec"),
            "parent": _input_manifest(run, "parent"),
            "source_packet": _input_manifest(run, "source_packet"),
            "candidate": _input_manifest(run, "candidate"),
            "dependencies": [
                {
                    "kind": item.kind,
                    "name": item.name,
                    "path": item.path,
                    "sha256": item.sha256,
                }
                for item in run.request.dependencies
            ],
            "live_tool_sources": [
                {"path": item.path, "sha256": item.sha256}
                for item in run.request.live.tool_sources
            ],
            "verification_stages": [
                "before_open",
                "after_open_before_render",
                "after_render",
                "before_ready",
            ],
        },
        "blender_versions": {
            "live": run.request.live.blender_version,
            "render": _numeric_blender_version(),
            "render_version_string": bpy.app.version_string,
            "render_build_hash": _blender_build_hash(),
            "candidate_file": list(_blender_file_version()),
            "pinned": run.request.live.pinned_blender.version,
            "toolchain": _input_manifest(run, "pinned_blender_toolchain"),
        },
        "renderer": {
            "candidate_opened_with_scripts": False,
            "load_handlers_suspended": dict(renderer_load_handlers),
            "tool_sources": [
                {
                    "module": item.module,
                    "path": _relative(run.workspace, item.path),
                    "sha256": item.sha256,
                }
                for item in run.tool_sources
            ],
            "preview": dict(preview),
        },
        "scene": {
            "name": scene.name,
            "boolean_requirements": dict(
                run.request.scene.boolean_requirements
            ),
            "view_layer": {
                "name": view_layer.name,
                "material_override": (
                    view_layer.material_override.name
                    if view_layer.material_override
                    else None
                ),
                "world_override": (
                    view_layer.world_override.name
                    if view_layer.world_override
                    else None
                ),
                "use_sky": bool(view_layer.use_sky),
            },
            "cameras": _camera_manifest(scene, view_layer, run.request),
            "object_count": len(run.request.scene.objects),
            "external_dependencies_match_request": True,
            "parent_open_attestation": {
                "path": run.request.live.parent_open_attestation.path,
                "sha256": run.request.live.parent_open_attestation.sha256,
                "spec_sha256": (
                    run.request.live.parent_open_attestation.spec_sha256
                ),
                "source_packet_sha256": (
                    run.request.live.parent_open_attestation.source_packet_sha256
                ),
                "opened_utc": (
                    run.request.live.parent_open_attestation.opened_utc
                ),
                "parent_file_version": list(
                    run.request.live.parent_open_attestation.parent_file_version
                ),
                "load_handlers_suspended": dict(
                    run.request.live.parent_open_attestation.load_handlers_suspended
                ),
            },
        },
        "dimensions": {
            "width": run.request.width,
            "height": run.request.height,
        },
        "hypothesis": run.request.hypothesis,
        "vetoes": list(run.request.vetoes),
        "outputs": outputs,
        "verdict_template": {
            "path": "VERDICT.md",
            "sha256": verdict_hash,
        },
        "timeline": dict(timeline),
        "render_completed_utc": completed_utc,
        "timings_seconds": {
            "live": {
                name: round(value, 6)
                for name, value in run.request.timings_seconds.items()
            },
            "render": {
                name: round(value, 6)
                for name, value in sorted(timings.items())
            },
        },
        "human_verdict": None,
    }


def _input_manifest(run: _Run, role: str) -> typing.Dict[str, object]:
    for item in run.inputs:
        if item.role == role:
            return {
                "path": item.record.path,
                "sha256": item.record.sha256,
                "hash_verified_before_render": True,
                "hash_verified_after_render": True,
            }
    raise CheckpointRenderError(f"internal error: no manifest input {role}")


def _relative(workspace: pathlib.Path, path: pathlib.Path) -> str:
    return path.resolve().relative_to(workspace.resolve()).as_posix()


def _verdict_template(request: _Request) -> str:
    lines = [
        f"# {request.checkpoint} visual verdict",
        "",
        "Status: `pending implementation-blind human review`",
        "",
        f"Candidate SHA-256: `{request.candidate.sha256}`",
        "",
        "## Frozen hypothesis",
        "",
        request.hypothesis,
        "",
        "## Hard vetoes",
        "",
    ]
    lines.extend(f"- {veto}" for veto in request.vetoes)
    lines.extend(["", "## Controlling views", ""])
    for name in request.scene.views:
        lines.append(
            f"- {name}: [beauty](beauty/{name}.png), "
            f"[silhouette](silhouette/{name}.png)"
        )
    lines.extend(
        [
            "",
            "## Absolute verdict",
            "",
            "- Result: `PASS` or `REJECT`",
            "- Major visible failure: `yes` or `no`",
            "- Evidence-backed observations:",
            "- Strategy decision:",
            "- T4 visual verdict complete (UTC):",
            "- T5 goal records updated (UTC):",
            "",
            "A quick checkpoint may reject but cannot accept the final model.",
            "",
            "This verdict template was generated by an LLM.",
            "",
        ]
    )
    return "\n".join(lines)


def _write_atomic_exclusive(path: pathlib.Path, payload: bytes) -> None:
    temporary = _write_complete_temporary(path, payload)
    try:
        os.link(temporary, path, follow_symlinks=False)
        _fsync_directory(path.parent)
    except FileExistsError as error:
        raise CheckpointRenderError(
            f"refusing to overwrite checkpoint output: {path}"
        ) from error
    finally:
        temporary.unlink(missing_ok=True)


def _verify_packet_before_ready(
    run: _Run,
    outputs: typing.Iterable[typing.Mapping[str, object]],
    verdict_hash: str,
    manifest_hash: str,
) -> None:
    expected_views = set(run.request.scene.views)
    root_names = {child.name for child in run.output_dir.iterdir()}
    expected_root = {
        _CANDIDATE_NAME,
        _REQUEST_NAME,
        _LOCK_NAME,
        "beauty",
        "silhouette",
        "VERDICT.md",
        "manifest.json",
    }
    if root_names != expected_root:
        raise CheckpointRenderError(
            "checkpoint directory changed before READY publication: "
            f"expected {sorted(expected_root)}, got {sorted(root_names)}"
        )
    for directory_name in ("beauty", "silhouette"):
        directory = run.output_dir / directory_name
        names = {path.name for path in directory.iterdir()}
        expected_names = {f"{name}.png" for name in expected_views}
        if names != expected_names:
            raise CheckpointRenderError(
                f"{directory_name} outputs changed before READY publication"
            )
    for output in outputs:
        relative = output.get("path")
        digest = output.get("sha256")
        if not isinstance(relative, str) or not isinstance(digest, str):
            raise CheckpointRenderError("invalid in-memory output record")
        path = (run.output_dir / relative).resolve()
        _require_below(path, run.output_dir, "render output")
        if not path.is_file() or sha256_file(path) != digest:
            raise CheckpointRenderError(
                f"render output changed before READY publication: {relative}"
            )
        image = decode_png(path)
        if (
            image.width != run.request.width
            or image.height != run.request.height
        ):
            raise CheckpointRenderError(
                f"render output dimensions changed before READY: {relative}"
            )
    if sha256_file(run.output_dir / "VERDICT.md") != verdict_hash:
        raise CheckpointRenderError(
            "VERDICT.md changed before READY publication"
        )
    if sha256_file(run.output_dir / "manifest.json") != manifest_hash:
        raise CheckpointRenderError(
            "manifest.json changed before READY publication"
        )


def _write_complete_temporary(
    path: pathlib.Path,
    payload: bytes,
) -> pathlib.Path:
    descriptor, name = tempfile.mkstemp(
        dir=path.parent,
        prefix=f".{path.name}.",
        suffix=f".{uuid.uuid4().hex}.tmp",
    )
    temporary = pathlib.Path(name)
    try:
        _write_descriptor(descriptor, payload)
        os.fsync(descriptor)
    except Exception:
        temporary.unlink(missing_ok=True)
        raise
    finally:
        os.close(descriptor)
    return temporary


def _publish_ready(
    run: _Run,
    lock_path: pathlib.Path,
    lock_payload: bytes,
    manifest_hash: str,
    verdict_hash: str,
) -> pathlib.Path:
    ready_path = run.output_dir / _READY_NAME
    ready_payload = canonical_json(
        {
            "schema_version": 1,
            "checkpoint": run.request.checkpoint,
            "request_sha256": run.request_sha256,
            "manifest": {"path": "manifest.json", "sha256": manifest_hash},
            "verdict_template": {
                "path": "VERDICT.md",
                "sha256": verdict_hash,
            },
            "completed_utc": _utc_now(),
        }
    ).encode("utf-8")
    temporary = _write_complete_temporary(ready_path, ready_payload)
    if lock_path.read_bytes() != lock_payload:
        raise CheckpointRenderError("RENDERING ownership marker changed")
    lock_path.unlink()
    _fsync_directory(run.output_dir)
    published = False
    try:
        os.link(temporary, ready_path, follow_symlinks=False)
        published = True
        temporary.unlink()
        _fsync_directory(run.output_dir)
    except Exception as error:
        # READY was absent immediately before the owned RENDERING marker was
        # removed. Any READY now present is either our failed publication or a
        # conflicting writer; remove it so a failed invocation never leaves a
        # packet that looks complete. Keep the hidden complete temporary for
        # diagnosis.
        if published or ready_path.exists() or ready_path.is_symlink():
            try:
                ready_path.unlink()
                _fsync_directory(run.output_dir)
            except OSError as cleanup_error:
                raise CheckpointRenderError(
                    "READY publication failed and the conflicting marker "
                    "could not be removed"
                ) from cleanup_error
        if isinstance(error, FileExistsError):
            raise CheckpointRenderError(
                "unexpected READY appeared during atomic publication"
            ) from error
        raise
    return ready_path


def _fsync_directory(path: pathlib.Path) -> None:
    flags = os.O_RDONLY
    if hasattr(os, "O_DIRECTORY"):
        flags |= os.O_DIRECTORY
    descriptor = os.open(path, flags)
    try:
        os.fsync(descriptor)
    finally:
        os.close(descriptor)


def _fsync_file(path: pathlib.Path) -> None:
    flags = os.O_RDONLY
    if hasattr(os, "O_NOFOLLOW"):
        flags |= os.O_NOFOLLOW
    descriptor = os.open(path, flags)
    try:
        os.fsync(descriptor)
    finally:
        os.close(descriptor)


def _utc_now() -> str:
    return (
        datetime.datetime.now(datetime.timezone.utc)
        .isoformat(timespec="microseconds")
        .replace("+00:00", "Z")
    )


def _numeric_blender_version() -> str:
    return ".".join(str(int(part)) for part in bpy.app.version[:3])


def _blender_file_version() -> typing.Tuple[int, int, int]:
    value = tuple(int(part) for part in bpy.app.version_file[:3])
    if len(value) != 3:
        raise CheckpointRenderError(
            "Blender did not report a three-component file version"
        )
    return (value[0], value[1], value[2])


def _blender_build_hash() -> str:
    value = getattr(bpy.app, "build_hash", "")
    if isinstance(value, bytes):
        return value.decode("ascii", errors="replace")
    return str(value)


def _execute(request_option: pathlib.Path) -> pathlib.Path:
    total_started = time.perf_counter()
    workspace = _resolve_workspace()
    request_path = _resolve_request_path(workspace, request_option)
    request, request_hash = _load_request(request_path)
    output_dir = request_path.parent
    candidate_path = _resolve_artifact_path(
        workspace, request.candidate, "candidate"
    )
    _validate_live_claims(workspace, request)
    _validate_request_directory(output_dir, request_path, candidate_path)
    lock_path, lock_payload = _claim_output(output_dir, request_hash)
    _validate_claimed_directory(output_dir)

    preflight_started = time.perf_counter()
    run = _build_run(
        workspace,
        request_path,
        request_hash,
        request,
    )
    _verify_frozen(run, "before_open")
    preflight_seconds = time.perf_counter() - preflight_started

    open_started = time.perf_counter()
    renderer_load_handlers = _open_candidate(
        _path_for_role(run.inputs, "candidate")
    )
    open_seconds = time.perf_counter() - open_started
    scene, view_layer = _validated_scene(run)
    _verify_frozen(run, "after_open_before_render")

    timeline = dict(run.request.timeline)
    outputs, render_timings, preview = _render_outputs(
        run, scene, view_layer, timeline
    )
    postflight_started = time.perf_counter()
    _verify_frozen(run, "after_render")
    _validate_open_dependencies(run)
    postflight_seconds = time.perf_counter() - postflight_started

    verdict = _verdict_template(run.request).encode("utf-8")
    verdict_path = run.output_dir / "VERDICT.md"
    _write_atomic_exclusive(verdict_path, verdict)
    verdict_hash = sha256_file(verdict_path)
    completed_utc = _utc_now()
    timings = {
        "preflight": preflight_seconds,
        "open_candidate": open_seconds,
        **render_timings,
        "postflight_hashes": postflight_seconds,
        "through_manifest_inputs": time.perf_counter() - total_started,
    }
    manifest = _manifest(
        run,
        scene,
        view_layer,
        outputs,
        preview,
        timeline,
        timings,
        renderer_load_handlers,
        verdict_hash,
        completed_utc,
    )
    manifest_path = run.output_dir / "manifest.json"
    _write_atomic_exclusive(
        manifest_path,
        canonical_json(manifest).encode("utf-8"),
    )
    manifest_hash = sha256_file(manifest_path)
    _verify_frozen(run, "before_ready")
    _validate_open_dependencies(run)
    _verify_packet_before_ready(
        run,
        outputs,
        verdict_hash,
        manifest_hash,
    )
    _publish_ready(
        run,
        lock_path,
        lock_payload,
        manifest_hash,
        verdict_hash,
    )
    return manifest_path


def main() -> None:
    """Render one exact snapshot request under repository-root out/."""
    parser = argparse.ArgumentParser(
        description="Render a split live checkpoint with pinned Blender."
    )
    parser.add_argument(
        "--snapshot-request",
        required=True,
        type=pathlib.Path,
    )
    options = parser.parse_args(_script_args())
    _execute(options.snapshot_request)


if __name__ == "__main__":
    try:
        main()
    except (
        CheckpointRenderError,
        CheckpointSpecError,
        OSError,
        PacketRenderError,
        PNGError,
    ) as error:
        raise RuntimeError(f"checkpoint_render: {error}") from error
