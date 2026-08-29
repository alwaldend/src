"""Snapshot-only handoff from a live Blender scene to a pinned renderer.

The live process proves exact parent ancestry, freezes a candidate, and writes
``snapshot_request.json`` last. It deliberately does not render: rendering is
performed by the repository-pinned background Blender after reopening the
candidate with auto-execution disabled.
"""

import contextlib
import datetime
import hashlib
import importlib
import json
import math
import os
import pathlib
import tempfile
import time
import typing

from projects.renders.cmd.fumo_review import (
    checkpoint_spec as checkpoint_spec_module,
)
from projects.renders.cmd.fumo_review import render_spec as render_spec_module
from projects.renders.cmd.fumo_review.checkpoint_spec import (
    Artifact,
    CheckpointSpec,
    canonical_json,
    load_checkpoint_spec_with_hash,
    resolve_out_directory,
    resolve_workspace_artifact,
    resolve_workspace_file,
    sha256_file,
)

bpy: typing.Any = importlib.import_module("bpy")

_SCRIPT_DIR = pathlib.Path(__file__).resolve().parent
_DEFAULT_WORKSPACE = _SCRIPT_DIR.parents[3]
_LOADED_SOURCE_SHA256 = hashlib.sha256(
    pathlib.Path(__file__).read_bytes()
).hexdigest()
_SAVE_HANDLER_NAMES = ("save_pre", "save_post", "save_post_fail")
_LOAD_HANDLER_NAMES = (
    "load_pre",
    "version_update",
    "load_post",
    "load_post_fail",
)
_PARENT_ATTESTATION_KEY = "fumo_review_parent_open_attestation_v1"
_PINNED_BLENDER_CONFIG = pathlib.PurePosixPath(
    "tools/blender/binary_toolchain.json"
)
_MAX_REQUEST_BYTES = 16 * 1024 * 1024
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


def _invalidate_parent_attestation(_unused: typing.Any) -> None:
    bpy.app.driver_namespace.pop(_PARENT_ATTESTATION_KEY, None)


setattr(
    _invalidate_parent_attestation,
    "_fumo_review_parent_invalidator",
    True,
)
_persistent_invalidator = bpy.app.handlers.persistent(
    _invalidate_parent_attestation
)
if not any(
    getattr(handler, "_fumo_review_parent_invalidator", False)
    for handler in bpy.app.handlers.load_pre
):
    bpy.app.handlers.load_pre.append(_persistent_invalidator)


class LiveCheckpointError(RuntimeError):
    """Raised when a live snapshot cannot preserve its evidence contract."""


def open_checkpoint_parent(
    spec_path: typing.Union[str, pathlib.Path],
    workspace_root: typing.Optional[typing.Union[str, pathlib.Path]] = None,
) -> pathlib.Path:
    """Open and attest the capsule parent before any builder mutates it.

    Snapshot publication requires this controlled open. Suppressing Python load
    handlers prevents an add-on from changing the just-loaded parent before the
    attestation is issued. Any later normal file load invalidates it.
    """
    workspace = _resolve_workspace(workspace_root)
    spec_file = resolve_workspace_file(
        workspace,
        pathlib.Path(spec_path),
        "checkpoint spec",
    )
    spec, spec_hash = load_checkpoint_spec_with_hash(spec_file)
    parent_path = resolve_workspace_artifact(workspace, spec.parent)
    source_packet_path = resolve_workspace_artifact(
        workspace,
        spec.source_packet,
    )
    pinned_blender = _pinned_blender(workspace)
    pinned_record = typing.cast(
        typing.Dict[str, str], pinned_blender["toolchain"]
    )
    _invalidate_parent_attestation(None)
    handler_counts: typing.Dict[str, int] = {}
    with _suspended_handlers(_LOAD_HANDLER_NAMES, handler_counts):
        result = bpy.ops.wm.open_mainfile(
            filepath=str(parent_path),
            load_ui=False,
            use_scripts=False,
        )
    if "FINISHED" not in result:
        raise LiveCheckpointError(
            f"Blender did not open checkpoint parent: {parent_path}"
        )
    parent_file_version = tuple(int(value) for value in bpy.app.version_file)
    live_version = tuple(int(value) for value in bpy.app.version[:3])
    if parent_file_version[:2] > live_version[:2]:
        raise LiveCheckpointError(
            "checkpoint parent was authored by a newer Blender major/minor "
            f"than the live process: file={parent_file_version}, "
            f"live={live_version}"
        )
    if bpy.data.is_dirty:
        raise LiveCheckpointError(
            "controlled parent open produced a dirty Blender database"
        )
    _require_current_parent(parent_path, spec.parent.sha256)
    frozen = [
        _record(workspace, spec_file, spec_hash),
        _record(workspace, parent_path, spec.parent.sha256),
        _record(workspace, source_packet_path, spec.source_packet.sha256),
        pinned_record,
    ]
    _verify_records(workspace, frozen)
    attestation = {
        "path": spec.parent.path,
        "sha256": spec.parent.sha256,
        "spec_sha256": spec_hash,
        "source_packet_sha256": spec.source_packet.sha256,
        "parent_file_version": list(parent_file_version),
        "opened_utc": _utc_now(),
        "load_handlers_suspended": handler_counts,
    }
    bpy.app.driver_namespace[_PARENT_ATTESTATION_KEY] = attestation
    return parent_path


def snapshot_checkpoint(
    spec_path: typing.Union[str, pathlib.Path],
    output_directory: typing.Union[str, pathlib.Path],
    workspace_root: typing.Optional[typing.Union[str, pathlib.Path]] = None,
) -> pathlib.Path:
    """Freeze one exact-parent candidate and return its render request path.

    A successful return means the output directory was atomically claimed,
    the live file matched the capsule's parent byte-for-byte, save handlers
    were suspended, state was preserved, and all frozen inputs were rehashed.
    It does not mean the candidate rendered or passed visual review.
    """
    workspace = _resolve_workspace(workspace_root)
    spec_file = resolve_workspace_file(
        workspace,
        pathlib.Path(spec_path),
        "checkpoint spec",
    )
    spec, spec_hash = load_checkpoint_spec_with_hash(spec_file)
    parent_path = resolve_workspace_artifact(workspace, spec.parent)
    source_packet_path = resolve_workspace_artifact(
        workspace,
        spec.source_packet,
    )
    parent_attestation = _require_exact_parent(
        parent_path,
        spec.parent.path,
        spec.parent.sha256,
        spec_hash,
        spec.source_packet.sha256,
    )
    pinned_blender = _pinned_blender(workspace)
    pinned_record = typing.cast(
        typing.Dict[str, str], pinned_blender["toolchain"]
    )
    scene = _validated_scene(spec)
    if bpy.context.mode != "OBJECT":
        raise LiveCheckpointError(
            "live snapshot requires Object Mode to preserve edit state"
        )

    dependencies = _external_dependencies(workspace)
    tools = _tool_sources(workspace)
    frozen_inputs = [
        _record(workspace, spec_file, spec_hash),
        _record(workspace, parent_path, spec.parent.sha256),
        _record(workspace, source_packet_path, spec.source_packet.sha256),
        *dependencies,
        *tools,
        pinned_record,
    ]
    _verify_records(workspace, frozen_inputs)

    output_dir = resolve_out_directory(
        workspace,
        pathlib.Path(output_directory),
    )
    candidate_path = output_dir / "candidate.blend"
    request_path = output_dir / "snapshot_request.json"

    context_before = _context_signature()
    scene_before = _scene_record(scene, spec)
    dirty_before = bool(bpy.data.is_dirty)
    handler_counts = _handler_count_snapshot(_SAVE_HANDLER_NAMES)
    preflight_utc = _utc_now()
    placeholder_candidate = _record(
        workspace,
        candidate_path,
        "0" * 64,
    )
    preflight_request = _snapshot_request(
        spec=spec,
        spec_record=frozen_inputs[0],
        parent_record=frozen_inputs[1],
        source_packet_record=frozen_inputs[2],
        candidate_record=placeholder_candidate,
        dependencies=dependencies,
        tools=tools,
        parent_attestation=parent_attestation,
        pinned_blender=pinned_blender,
        scene_record=scene_before,
        context_before=context_before,
        context_after=context_before,
        dirty_before=dirty_before,
        handler_counts=handler_counts,
        candidate_saved_utc=preflight_utc,
        snapshot_hashed_utc=preflight_utc,
        save_seconds=999999999.999999,
        hash_seconds=999999999.999999,
    )
    _request_payload(preflight_request)
    _reserve_output_directory(workspace, output_dir)
    save_started = time.perf_counter()
    with _suspended_handlers(_SAVE_HANDLER_NAMES, handler_counts):
        result = bpy.ops.wm.save_as_mainfile(
            filepath=str(candidate_path),
            copy=True,
        )
    save_seconds = time.perf_counter() - save_started
    if "FINISHED" not in result or not candidate_path.is_file():
        raise LiveCheckpointError(
            f"Blender did not save checkpoint copy: {candidate_path}"
        )
    _fsync_file(candidate_path)
    _fsync_directory(output_dir)
    candidate_saved_utc = _utc_now()

    context_after = _context_signature()
    scene_after = _scene_record(scene, spec)
    dirty_after = bool(bpy.data.is_dirty)
    if context_after != context_before:
        raise LiveCheckpointError(
            "saving the checkpoint copy changed the live Blender context"
        )
    if dirty_after != dirty_before:
        raise LiveCheckpointError(
            "saving the checkpoint copy changed Blender's dirty state"
        )
    if scene_after != scene_before:
        raise LiveCheckpointError(
            "saving the checkpoint copy changed scene or camera identity"
        )
    _require_exact_parent(
        parent_path,
        spec.parent.path,
        spec.parent.sha256,
        spec_hash,
        spec.source_packet.sha256,
    )

    hash_started = time.perf_counter()
    candidate_hash = sha256_file(candidate_path)
    candidate = _record(workspace, candidate_path, candidate_hash)
    _verify_records(workspace, [*frozen_inputs, candidate])
    hash_seconds = time.perf_counter() - hash_started
    snapshot_hashed_utc = _utc_now()
    request = _snapshot_request(
        spec=spec,
        spec_record=frozen_inputs[0],
        parent_record=frozen_inputs[1],
        source_packet_record=frozen_inputs[2],
        candidate_record=candidate,
        dependencies=dependencies,
        tools=tools,
        parent_attestation=parent_attestation,
        pinned_blender=pinned_blender,
        scene_record=scene_before,
        context_before=context_before,
        context_after=context_after,
        dirty_before=dirty_before,
        handler_counts=handler_counts,
        candidate_saved_utc=candidate_saved_utc,
        snapshot_hashed_utc=snapshot_hashed_utc,
        save_seconds=save_seconds,
        hash_seconds=hash_seconds,
    )
    _write_atomic(request_path, _request_payload(request))
    return request_path


def _resolve_workspace(
    workspace_root: typing.Optional[typing.Union[str, pathlib.Path]],
) -> pathlib.Path:
    if workspace_root is not None:
        result = pathlib.Path(workspace_root).resolve()
    elif os.environ.get("BUILD_WORKSPACE_DIRECTORY"):
        result = pathlib.Path(
            typing.cast(str, os.environ["BUILD_WORKSPACE_DIRECTORY"])
        ).resolve()
    else:
        result = _DEFAULT_WORKSPACE.resolve()
    if not result.is_dir() or not (result / "out").is_dir():
        raise LiveCheckpointError(
            f"workspace must exist and contain repository-root out/: {result}"
        )
    return result


def _pinned_blender(workspace: pathlib.Path) -> typing.Dict[str, object]:
    config_path = resolve_workspace_file(
        workspace,
        pathlib.Path(_PINNED_BLENDER_CONFIG.as_posix()),
        "pinned Blender toolchain configuration",
    )
    try:
        payload = config_path.read_bytes()
        document = json.loads(payload.decode("utf-8"))
        toolchains = document["toolchains"]
        versions = {
            item["version"]
            for item in toolchains
            if item.get("name") == "blender"
        }
        archive_versions = {
            item["toolchain_version"]
            for item in document["archives"]
            if item.get("toolchain_name") == "blender"
        }
    except (KeyError, TypeError, UnicodeError, json.JSONDecodeError) as error:
        raise LiveCheckpointError(
            f"invalid pinned Blender configuration: {config_path}"
        ) from error
    if len(versions) != 1:
        raise LiveCheckpointError(
            "pinned Blender configuration must name exactly one version"
        )
    version = versions.pop()
    if archive_versions != {version}:
        raise LiveCheckpointError(
            "pinned Blender toolchain and archive versions disagree"
        )
    try:
        pinned_tuple = tuple(int(part) for part in version.split("."))
    except (AttributeError, ValueError) as error:
        raise LiveCheckpointError(
            f"invalid pinned Blender version: {version!r}"
        ) from error
    if len(pinned_tuple) != 3:
        raise LiveCheckpointError(
            f"pinned Blender version must have three components: {version!r}"
        )
    live_tuple = tuple(int(part) for part in bpy.app.version[:3])
    if live_tuple > pinned_tuple:
        raise LiveCheckpointError(
            f"live Blender {bpy.app.version_string} is newer than pinned "
            f"Blender {version}; the pinned renderer may not reopen its file"
        )
    return {
        "version": version,
        "toolchain": _record(
            workspace,
            config_path,
            hashlib.sha256(payload).hexdigest(),
        ),
    }


def _reserve_output_directory(
    workspace: pathlib.Path,
    output_dir: pathlib.Path,
) -> None:
    output_dir.parent.mkdir(parents=True, exist_ok=True)
    resolved = resolve_out_directory(workspace, output_dir)
    if resolved != output_dir:
        raise LiveCheckpointError(
            "checkpoint output changed identity while creating its parent"
        )
    try:
        output_dir.mkdir()
    except FileExistsError as error:
        raise LiveCheckpointError(
            "checkpoint output must be a new, nonexistent directory: "
            f"{output_dir}"
        ) from error
    _fsync_directory(output_dir.parent)


def _require_current_parent(
    parent_path: pathlib.Path,
    expected_hash: str,
) -> None:
    if not bpy.data.filepath:
        raise LiveCheckpointError(
            "live Blender file is unsaved; exact parent ancestry is unknown"
        )
    live_path = pathlib.Path(bpy.path.abspath(bpy.data.filepath)).resolve()
    if live_path != parent_path.resolve():
        raise LiveCheckpointError(
            "live Blender file does not match the frozen parent: "
            f"live={live_path}, parent={parent_path.resolve()}"
        )
    actual_hash = sha256_file(live_path)
    if actual_hash != expected_hash:
        raise LiveCheckpointError(
            "current parent file bytes do not match the capsule hash"
        )


def _require_exact_parent(
    parent_path: pathlib.Path,
    parent_record_path: str,
    expected_hash: str,
    spec_hash: str,
    source_packet_hash: str,
) -> typing.Dict[str, object]:
    _require_current_parent(parent_path, expected_hash)
    value = bpy.app.driver_namespace.get(_PARENT_ATTESTATION_KEY)
    if not isinstance(value, dict):
        raise LiveCheckpointError(
            "parent ancestry is unattested; call open_checkpoint_parent() "
            "before applying builder edits"
        )
    expected = {
        "path": parent_record_path,
        "sha256": expected_hash,
        "spec_sha256": spec_hash,
        "source_packet_sha256": source_packet_hash,
    }
    actual = {
        "path": value.get("path"),
        "sha256": value.get("sha256"),
        "spec_sha256": value.get("spec_sha256"),
        "source_packet_sha256": value.get("source_packet_sha256"),
    }
    if actual != expected:
        raise LiveCheckpointError(
            "controlled parent-open attestation does not match this capsule"
        )
    return dict(value)


def _validated_scene(spec: CheckpointSpec) -> typing.Any:
    scene = bpy.data.scenes.get(spec.scene)
    if scene is None:
        available = sorted(item.name for item in bpy.data.scenes)
        raise LiveCheckpointError(
            f"scene {spec.scene!r} does not exist; available: {available}"
        )
    errors: typing.List[str] = []
    if bpy.context.scene is not scene:
        errors.append(f"scene {scene.name!r} must be the active live scene")
    else:
        errors.extend(_active_view_layer_errors(scene, spec))
    errors.extend(_camera_errors(scene, spec))
    errors.extend(_boolean_requirement_errors(scene, spec))
    if errors:
        details = "\n".join(f"  - {error}" for error in errors)
        raise LiveCheckpointError(f"invalid live checkpoint scene:\n{details}")
    return scene


def _camera_errors(
    scene: typing.Any,
    spec: CheckpointSpec,
) -> typing.List[str]:
    errors: typing.List[str] = []
    for view in spec.views.values():
        camera = scene.objects.get(view.camera)
        if camera is None:
            errors.append(f"camera {view.camera!r} is missing")
        elif camera.type != "CAMERA":
            errors.append(
                f"camera {view.camera!r} names a {camera.type}, not a CAMERA"
            )
        elif camera.data.type not in {"PERSP", "ORTHO"}:
            errors.append(
                f"camera {view.camera!r} uses unsupported projection "
                f"{camera.data.type!r}"
            )
    return errors


def _boolean_requirement_errors(
    scene: typing.Any,
    spec: CheckpointSpec,
) -> typing.List[str]:
    errors: typing.List[str] = []
    for name, expected in spec.scene_boolean_requirements.items():
        if name not in scene:
            errors.append(
                f"required scene boolean property {name!r} is missing"
            )
            continue
        actual = scene[name]
        if not isinstance(actual, bool) or actual is not expected:
            errors.append(
                f"required scene boolean property {name!r} must be "
                f"{expected}, got {actual!r}"
            )
    return errors


def _active_view_layer_errors(
    scene: typing.Any,
    spec: CheckpointSpec,
) -> typing.List[str]:
    errors: typing.List[str] = []
    view_layer = bpy.context.view_layer
    scene_layer = scene.view_layers.get(view_layer.name)
    if (
        scene_layer is None
        or scene_layer.as_pointer() != view_layer.as_pointer()
    ):
        return [
            f"active view layer {view_layer.name!r} does not belong to "
            f"scene {scene.name!r}"
        ]
    if hasattr(view_layer, "use") and not view_layer.use:
        errors.append(f"active view layer {view_layer.name!r} is disabled")
    layer_objects = list(view_layer.objects)
    if not any(
        obj.type in _GEOMETRY_TYPES and not obj.hide_render
        for obj in layer_objects
    ):
        errors.append("active view layer has no renderable geometry")
    for view in spec.views.values():
        camera = scene.objects.get(view.camera)
        if camera is not None and camera not in layer_objects:
            errors.append(
                f"camera {view.camera!r} is excluded from active view layer"
            )
    return errors


@contextlib.contextmanager
def _suspended_handlers(
    names: typing.Iterable[str],
    counts: typing.MutableMapping[str, int],
) -> typing.Iterator[None]:
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


def _handler_count_snapshot(
    names: typing.Iterable[str],
) -> typing.Dict[str, int]:
    return {
        name: len(handlers)
        for name in names
        if (handlers := getattr(bpy.app.handlers, name, None)) is not None
    }


def _external_dependencies(
    workspace: pathlib.Path,
) -> typing.List[typing.Dict[str, str]]:
    records: typing.List[typing.Dict[str, str]] = []
    path_map = bpy.data.file_path_map(include_libraries=True)
    for datablock, raw_paths in path_map.items():
        if _datablock_is_dirty(datablock):
            raise LiveCheckpointError(
                f"external datablock {datablock.name!r} has unsaved in-memory "
                "changes; save or repack it before checkpointing"
            )
        source = getattr(datablock, "source", "")
        if source in {
            "SEQUENCE",
            "TILED",
        } or _true_boolean_attribute(datablock, "is_sequence"):
            raise LiveCheckpointError(
                f"external datablock {datablock.name!r} uses an unbounded "
                "sequence or tiled source"
            )
        if not raw_paths:
            continue
        for raw in sorted(raw_paths):
            if not raw or raw == "<builtin>":
                continue
            records.append(
                _external_dependency_record(workspace, datablock, raw)
            )
    unique: typing.Dict[typing.Tuple[str, str, str], typing.Dict[str, str]]
    unique = {}
    for record in records:
        key = (record["kind"], record["name"], record["path"])
        previous = unique.get(key)
        if previous is not None and previous["sha256"] != record["sha256"]:
            raise LiveCheckpointError(
                f"external dependency changed while hashing: {record['path']}"
            )
        unique[key] = record
    return [unique[key] for key in sorted(unique)]


def _datablock_is_dirty(datablock: typing.Any) -> bool:
    return any(
        _true_boolean_attribute(datablock, name)
        for name in ("is_dirty", "is_modified")
    )


def _true_boolean_attribute(datablock: typing.Any, name: str) -> bool:
    value = getattr(datablock, name, False)
    return isinstance(value, bool) and value


def _external_dependency_record(
    workspace: pathlib.Path,
    datablock: typing.Any,
    raw: str,
) -> typing.Dict[str, str]:
    if any(token in raw for token in ("#", "<UDIM>", "<UVTILE>")):
        raise LiveCheckpointError(
            f"external dependency {datablock.name!r} uses an unbounded "
            f"sequence or tile path: {raw}"
        )
    kind = getattr(datablock, "id_type", datablock.bl_rna.identifier)
    try:
        absolute = pathlib.Path(
            bpy.path.abspath(
                raw,
                library=getattr(datablock, "library", None),
            )
        )
    except TypeError:
        absolute = pathlib.Path(bpy.path.abspath(raw))
    try:
        path = resolve_workspace_file(
            workspace,
            absolute,
            f"external {kind} dependency {datablock.name!r}",
        )
    except Exception as error:
        raise LiveCheckpointError(
            f"external {kind} dependency {datablock.name!r} must be one "
            f"existing normal file inside the workspace: {absolute}"
        ) from error
    record = _record(workspace, path, sha256_file(path))
    record.update({"kind": str(kind), "name": datablock.name})
    return record


def _tool_sources(
    workspace: pathlib.Path,
) -> typing.List[typing.Dict[str, str]]:
    loaded = (
        (pathlib.Path(__file__), _LOADED_SOURCE_SHA256),
        (
            pathlib.Path(typing.cast(str, checkpoint_spec_module.__file__)),
            checkpoint_spec_module._LOADED_SOURCE_SHA256,
        ),
        (
            pathlib.Path(typing.cast(str, render_spec_module.__file__)),
            render_spec_module._LOADED_SOURCE_SHA256,
        ),
    )
    result = []
    for path, loaded_digest in loaded:
        resolved = path.resolve()
        if resolved.parent != _SCRIPT_DIR:
            raise LiveCheckpointError(
                f"checkpoint module was loaded outside its package: {resolved}"
            )
        result.append(_record(workspace, resolved, loaded_digest))
    return result


def _record(
    workspace: pathlib.Path,
    path: pathlib.Path,
    digest: str,
) -> typing.Dict[str, str]:
    try:
        relative = path.resolve().relative_to(workspace.resolve()).as_posix()
    except ValueError as error:
        raise LiveCheckpointError(
            f"checkpoint input resolves outside the workspace: {path}"
        ) from error
    return {"path": relative, "sha256": digest}


def _verify_records(
    workspace: pathlib.Path,
    records: typing.Iterable[typing.Mapping[str, str]],
) -> None:
    for record in records:
        artifact = Artifact(path=record["path"], sha256=record["sha256"])
        try:
            resolve_workspace_artifact(workspace, artifact)
        except Exception as error:
            raise LiveCheckpointError(
                f"frozen checkpoint input changed: {record['path']}"
            ) from error


def _context_signature() -> typing.Dict[str, object]:
    active = bpy.context.view_layer.objects.active
    return {
        "active_scene": bpy.context.scene.name,
        "active_view_layer": bpy.context.view_layer.name,
        "active_object": active.name if active else None,
        "selected_objects": sorted(
            obj.name for obj in bpy.context.selected_objects
        ),
        "mode": bpy.context.mode,
        "source_filepath": bpy.data.filepath,
    }


def _camera_record(camera: typing.Any) -> typing.Dict[str, object]:
    data = camera.data
    scalars = {
        "lens": float(data.lens),
        "ortho_scale": float(data.ortho_scale),
        "sensor_width": float(data.sensor_width),
        "sensor_height": float(data.sensor_height),
        "shift_x": float(data.shift_x),
        "shift_y": float(data.shift_y),
        "clip_start": float(data.clip_start),
        "clip_end": float(data.clip_end),
    }
    matrix = [[float(value) for value in row] for row in camera.matrix_world]
    if not all(math.isfinite(value) for value in scalars.values()) or not all(
        math.isfinite(value) for row in matrix for value in row
    ):
        raise LiveCheckpointError(
            f"camera {camera.name!r} contains a non-finite numeric value"
        )
    positive = (
        "lens",
        "ortho_scale",
        "sensor_width",
        "sensor_height",
        "clip_start",
        "clip_end",
    )
    if any(scalars[name] <= 0.0 for name in positive) or (
        scalars["clip_end"] <= scalars["clip_start"]
    ):
        raise LiveCheckpointError(
            f"camera {camera.name!r} has invalid positive or clip values"
        )
    return {
        "name": camera.name,
        "projection": data.type,
        "lens": scalars["lens"],
        "ortho_scale": scalars["ortho_scale"],
        "sensor_fit": data.sensor_fit,
        "sensor_width": scalars["sensor_width"],
        "sensor_height": scalars["sensor_height"],
        "shift_x": scalars["shift_x"],
        "shift_y": scalars["shift_y"],
        "clip_start": scalars["clip_start"],
        "clip_end": scalars["clip_end"],
        "matrix_world": matrix,
    }


def _scene_record(
    scene: typing.Any,
    spec: CheckpointSpec,
) -> typing.Dict[str, object]:
    return {
        "name": scene.name,
        "boolean_requirements": {
            name: scene[name] for name in spec.scene_boolean_requirements
        },
        "views": {
            view.name: _camera_record(scene.objects[view.camera])
            for view in spec.views.values()
        },
        "view_layers": sorted(layer.name for layer in scene.view_layers),
        "objects": [
            {
                "name": obj.name,
                "type": obj.type,
                "data": (
                    obj.data.name if getattr(obj, "data", None) else None
                ),
                "hide_render": bool(obj.hide_render),
            }
            for obj in sorted(scene.objects, key=lambda item: item.name)
        ],
    }


def _snapshot_request(
    *,
    spec: CheckpointSpec,
    spec_record: typing.Mapping[str, str],
    parent_record: typing.Mapping[str, str],
    source_packet_record: typing.Mapping[str, str],
    candidate_record: typing.Mapping[str, str],
    dependencies: typing.List[typing.Dict[str, str]],
    tools: typing.List[typing.Dict[str, str]],
    parent_attestation: typing.Mapping[str, object],
    pinned_blender: typing.Mapping[str, object],
    scene_record: typing.Mapping[str, object],
    context_before: typing.Mapping[str, object],
    context_after: typing.Mapping[str, object],
    dirty_before: bool,
    handler_counts: typing.Mapping[str, int],
    candidate_saved_utc: str,
    snapshot_hashed_utc: str,
    save_seconds: float,
    hash_seconds: float,
) -> typing.Dict[str, object]:
    return {
        "schema_version": 1,
        "checkpoint": spec.checkpoint,
        "spec": dict(spec_record),
        "parent": dict(parent_record),
        "source_packet": dict(source_packet_record),
        "candidate": dict(candidate_record),
        "dependencies": dependencies,
        "live": {
            "blender_version": bpy.app.version_string,
            "tool_sources": tools,
            "context_before_snapshot": dict(context_before),
            "context_after_snapshot": dict(context_after),
            "dirty_before_snapshot": dirty_before,
            "dirty_state_preserved": True,
            "exact_parent_verified": True,
            "save_handlers_suspended": dict(handler_counts),
            "parent_open_attestation": dict(parent_attestation),
            "pinned_blender": dict(pinned_blender),
        },
        "scene": dict(scene_record),
        "dimensions": {
            "width": spec.resolution.width,
            "height": spec.resolution.height,
        },
        "hypothesis": spec.hypothesis,
        "vetoes": list(spec.vetoes),
        "timeline": {
            "t0_source_packet_frozen_utc": spec.t0_source_packet_frozen_utc,
            "t1_candidate_saved_utc": candidate_saved_utc,
            "t2_snapshot_hashed_utc": snapshot_hashed_utc,
            "t3_first_pixels_written_utc": None,
            "t4_visual_verdict_complete_utc": None,
            "t5_goal_records_updated_utc": None,
        },
        "timings_seconds": {
            "candidate_save_operator": round(save_seconds, 6),
            "candidate_hash_and_reverify": round(hash_seconds, 6),
        },
    }


def _write_atomic(path: pathlib.Path, payload: str) -> None:
    temporary_name = ""
    try:
        with tempfile.NamedTemporaryFile(
            "w",
            encoding="utf-8",
            dir=path.parent,
            prefix=f".{path.name}.",
            suffix=".tmp",
            delete=False,
        ) as stream:
            temporary_name = stream.name
            stream.write(payload)
            stream.flush()
            os.fsync(stream.fileno())
        pathlib.Path(temporary_name).replace(path)
        _fsync_directory(path.parent)
    finally:
        if temporary_name:
            pathlib.Path(temporary_name).unlink(missing_ok=True)


def _fsync_file(path: pathlib.Path) -> None:
    with path.open("rb") as stream:
        os.fsync(stream.fileno())


def _fsync_directory(path: pathlib.Path) -> None:
    descriptor = os.open(path, os.O_RDONLY)
    try:
        os.fsync(descriptor)
    finally:
        os.close(descriptor)


def _request_payload(request: typing.Mapping[str, object]) -> str:
    try:
        payload = canonical_json(request)
    except (TypeError, ValueError) as error:
        raise LiveCheckpointError(
            f"checkpoint request is not canonical JSON: {error}"
        ) from error
    size = len(payload.encode("utf-8"))
    if size > _MAX_REQUEST_BYTES:
        raise LiveCheckpointError(
            f"checkpoint request is {size} bytes; maximum is "
            f"{_MAX_REQUEST_BYTES}"
        )
    return payload


def _utc_now() -> str:
    return (
        datetime.datetime.now(datetime.timezone.utc)
        .isoformat(timespec="microseconds")
        .replace("+00:00", "Z")
    )


__all__ = [
    "LiveCheckpointError",
    "open_checkpoint_parent",
    "snapshot_checkpoint",
]
