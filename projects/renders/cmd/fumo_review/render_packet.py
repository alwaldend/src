"""Blender-side driver for deterministic Fumo review packets."""

import argparse
import hashlib
import importlib
import json
import os
import pathlib
import struct
import sys
import tempfile
import typing

_SCRIPT_DIR = pathlib.Path(__file__).resolve().parent
_RUNFILES_ROOT = _SCRIPT_DIR.parents[3]
if str(_RUNFILES_ROOT) not in sys.path:
    sys.path.insert(0, str(_RUNFILES_ROOT))

from projects.renders.cmd.fumo_review.render_spec import (  # noqa: E402
    RGB,
    ComponentIDPass,
    RenderSpec,
    SpecError,
    inverse_srgb,
    load_spec,
    resolve_output_path,
)

bpy: typing.Any = importlib.import_module("bpy")


class RenderError(RuntimeError):
    """Raised when a validated spec cannot be rendered by Blender."""


def _script_args() -> typing.List[str]:
    try:
        separator = sys.argv.index("--")
    except ValueError:
        return []
    return sys.argv[separator + 1 :]  # noqa: E203


def _resolve_cli_path(path: pathlib.Path) -> pathlib.Path:
    if path.is_absolute():
        return path.resolve()
    workspace = os.environ.get("BUILD_WORKSPACE_DIRECTORY")
    base = pathlib.Path(workspace) if workspace else pathlib.Path.cwd()
    return (base / path).resolve()


def _sha256(path: pathlib.Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        while True:
            chunk = stream.read(1024 * 1024)
            if not chunk:
                break
            digest.update(chunk)
    return digest.hexdigest()


def _canonicalize_png(path: pathlib.Path) -> None:
    """Remove Blender timing and host metadata while preserving image data."""
    signature = b"\x89PNG\r\n\x1a\n"
    data = path.read_bytes()
    if not data.startswith(signature):
        raise RenderError(f"Blender output is not a PNG: {path}")
    result = bytearray(signature)
    offset = len(signature)
    saw_end = False
    while offset < len(data):
        if offset + 12 > len(data):
            raise RenderError(f"truncated PNG chunk in {path}")
        length = struct.unpack_from(">I", data, offset)[0]
        end = offset + 12 + length
        if end > len(data):
            raise RenderError(f"truncated PNG chunk data in {path}")
        chunk_type = struct.unpack_from("4s", data, offset + 4)[0]
        if chunk_type not in {b"eXIf", b"iTXt", b"tEXt", b"tIME", b"zTXt"}:
            result.extend(data[offset:end])
        offset = end
        if chunk_type == b"IEND":
            saw_end = True
            break
    if not saw_end or offset != len(data):
        raise RenderError(f"invalid PNG chunk stream in {path}")

    temporary_name = ""
    try:
        with tempfile.NamedTemporaryFile(
            "wb",
            dir=path.parent,
            prefix=f".{path.name}.",
            suffix=".tmp",
            delete=False,
        ) as stream:
            temporary_name = stream.name
            stream.write(result)
        pathlib.Path(temporary_name).replace(path)
    finally:
        if temporary_name:
            pathlib.Path(temporary_name).unlink(missing_ok=True)


def _source_label(path: pathlib.Path) -> str:
    workspace = os.environ.get("BUILD_WORKSPACE_DIRECTORY")
    if workspace:
        try:
            return path.relative_to(
                pathlib.Path(workspace).resolve()
            ).as_posix()
        except ValueError:
            pass
    return path.as_posix()


def _open_source(path: pathlib.Path) -> None:
    if not path.is_file():
        raise RenderError(f"blend file does not exist: {path}")
    result = bpy.ops.wm.open_mainfile(
        filepath=str(path),
        load_ui=False,
        use_scripts=False,
    )
    if "FINISHED" not in result:
        raise RenderError(f"Blender could not open blend file: {path}")


def _validated_scene(spec: RenderSpec) -> typing.Any:
    scene = bpy.data.scenes.get(spec.scene)
    if scene is None:
        available = sorted(item.name for item in bpy.data.scenes)
        raise RenderError(
            f"scene {spec.scene!r} does not exist; available scenes: {available}"
        )

    errors: typing.List[str] = []
    camera_names = {view.camera for view in spec.beauty_views.values()} | {
        component_pass.camera
        for component_pass in spec.component_id_passes.values()
    }
    for camera_name in sorted(camera_names):
        camera = scene.objects.get(camera_name)
        if camera is None:
            errors.append(f"camera {camera_name!r} is missing from the scene")
        elif camera.type != "CAMERA":
            errors.append(
                f"camera {camera_name!r} names a {camera.type}, not a CAMERA"
            )

    object_names = {
        object_name
        for component_pass in spec.component_id_passes.values()
        for object_name in component_pass.objects
    }
    for object_name in sorted(object_names):
        obj = scene.objects.get(object_name)
        if obj is None:
            errors.append(f"component object {object_name!r} is missing")
        elif obj.type != "MESH":
            errors.append(
                f"component object {object_name!r} is a {obj.type}, not a MESH"
            )
        elif obj.hide_render:
            errors.append(
                f"component object {object_name!r} has hide_render enabled"
            )
    if errors:
        details = "\n".join(f"  - {error}" for error in errors)
        raise RenderError(f"invalid scene names:\n{details}")
    return scene


def _configure_png(scene: typing.Any, spec: RenderSpec, mode: str) -> None:
    render = scene.render
    render.resolution_x = spec.resolution.width
    render.resolution_y = spec.resolution.height
    render.resolution_percentage = 100
    render.use_file_extension = True
    render.use_overwrite = True
    render.use_placeholder = False
    if hasattr(render, "threads_mode"):
        render.threads_mode = "FIXED"
        render.threads = 1
    render.image_settings.file_format = "PNG"
    render.image_settings.color_mode = mode
    render.image_settings.color_depth = "8"
    render.image_settings.compression = 15
    render.use_stamp = False
    for prop in render.bl_rna.properties:
        if prop.identifier.startswith("use_stamp_"):
            setattr(render, prop.identifier, False)


def _render_still(
    scene: typing.Any,
    camera_name: str,
    path: pathlib.Path,
) -> None:
    camera = scene.objects[camera_name]
    scene.camera = camera
    path.parent.mkdir(parents=True, exist_ok=True)
    scene.render.filepath = str(path)
    result = bpy.ops.render.render(write_still=True, scene=scene.name)
    if "FINISHED" not in result or not path.is_file():
        raise RenderError(
            f"render from camera {camera_name!r} did not write {path}"
        )
    _canonicalize_png(path)


def _render_beauty(
    scene: typing.Any,
    spec: RenderSpec,
    output_dir: pathlib.Path,
) -> typing.List[typing.Dict[str, object]]:
    _configure_png(scene, spec, "RGBA")
    outputs: typing.List[typing.Dict[str, object]] = []
    for view in spec.beauty_views.values():
        path = resolve_output_path(output_dir, view.output)
        _render_still(scene, view.camera, path)
        outputs.append(
            {
                "kind": "beauty",
                "name": view.name,
                "camera": view.camera,
                "path": view.output,
                "sha256": _sha256(path),
            }
        )
    return outputs


def _emission_material(rgb: RGB) -> typing.Any:
    material = bpy.data.materials.new(
        name=f".fumo_review_id_{rgb[0]:02x}{rgb[1]:02x}{rgb[2]:02x}"
    )
    material.use_nodes = True
    nodes = material.node_tree.nodes
    nodes.clear()
    output = nodes.new("ShaderNodeOutputMaterial")
    emission = nodes.new("ShaderNodeEmission")
    linear = inverse_srgb(rgb)
    emission.inputs["Color"].default_value = (*linear, 1.0)
    emission.inputs["Strength"].default_value = 1.0
    material.node_tree.links.new(
        emission.outputs["Emission"], output.inputs["Surface"]
    )
    return material


def _mesh_states(
    scene: typing.Any,
) -> typing.List[typing.Tuple[typing.Any, bool, typing.List[typing.Any]]]:
    states = []
    for obj in scene.objects:
        slots = [(slot.link, slot.material) for slot in obj.material_slots]
        states.append((obj, bool(obj.hide_render), slots))
    return states


def _set_object_material(
    obj: typing.Any,
    material: typing.Any,
    appended_data: typing.Dict[int, typing.Any],
) -> None:
    if not obj.material_slots:
        pointer = obj.data.as_pointer()
        if pointer not in appended_data:
            obj.data.materials.append(material)
            appended_data[pointer] = obj.data
    for slot in obj.material_slots:
        slot.link = "OBJECT"
        slot.material = material


def _restore_mesh_states(
    states: typing.List[
        typing.Tuple[typing.Any, bool, typing.List[typing.Any]]
    ],
    appended_data: typing.Mapping[int, typing.Any],
) -> None:
    for obj, hide_render, slots in states:
        obj.hide_render = hide_render
        for current, saved in zip(obj.material_slots, slots):
            current.link = saved[0]
            current.material = saved[1]
    for mesh in appended_data.values():
        mesh.materials.pop(index=len(mesh.materials) - 1)


def _apply_component_mapping(
    scene: typing.Any,
    component_pass: ComponentIDPass,
    materials: typing.Mapping[RGB, typing.Any],
    appended_data: typing.Dict[int, typing.Any],
) -> None:
    geometry_types = {
        "CURVE",
        "CURVES",
        "FONT",
        "GREASEPENCIL",
        "META",
        "POINTCLOUD",
        "SURFACE",
        "VOLUME",
    }
    for obj in scene.objects:
        if obj.type != "MESH":
            if obj.type in geometry_types and not obj.hide_render:
                obj.hide_render = True
            continue
        if obj.hide_render:
            continue
        rgb = component_pass.objects.get(obj.name)
        if rgb is None and component_pass.fallback.mode == "hide":
            obj.hide_render = True
            continue
        if rgb is None:
            rgb = component_pass.fallback.rgb
        if rgb is None:
            raise RenderError(
                f"component-ID pass {component_pass.name!r} has no fallback"
            )
        _set_object_material(obj, materials[rgb], appended_data)


def _component_colors(spec: RenderSpec) -> typing.Set[RGB]:
    colors = {
        rgb
        for component_pass in spec.component_id_passes.values()
        for rgb in component_pass.objects.values()
    }
    colors.update(
        component_pass.fallback.rgb
        for component_pass in spec.component_id_passes.values()
        if component_pass.fallback.rgb is not None
    )
    return colors


def _configure_component_scene(scene: typing.Any, spec: RenderSpec) -> None:
    _configure_png(scene, spec, "RGB")
    scene.render.engine = "BLENDER_EEVEE"
    scene.render.film_transparent = False
    scene.render.use_compositing = False
    scene.render.use_sequencer = False
    if hasattr(scene.render, "use_freestyle"):
        scene.render.use_freestyle = False
    if hasattr(scene.render, "dither_intensity"):
        scene.render.dither_intensity = 0.0
    scene.display_settings.display_device = "sRGB"
    scene.view_settings.view_transform = "Standard"
    scene.view_settings.look = "None"
    scene.view_settings.exposure = 0.0
    scene.view_settings.gamma = 1.0
    scene.view_settings.use_curve_mapping = False


def _black_world() -> typing.Any:
    world = bpy.data.worlds.new(name=".fumo_review_black_world")
    world.use_nodes = True
    background = world.node_tree.nodes.get("Background")
    if background is None:
        raise RenderError("Blender did not create a Background world node")
    background.inputs["Color"].default_value = (0.0, 0.0, 0.0, 1.0)
    background.inputs["Strength"].default_value = 0.0
    return world


def _render_component_ids(
    scene: typing.Any,
    spec: RenderSpec,
    output_dir: pathlib.Path,
) -> typing.List[typing.Dict[str, object]]:
    if not spec.component_id_passes:
        return []
    _configure_component_scene(scene, spec)
    materials = {
        rgb: _emission_material(rgb) for rgb in sorted(_component_colors(spec))
    }
    original_world = scene.world
    world = _black_world()
    scene.world = world
    outputs: typing.List[typing.Dict[str, object]] = []
    try:
        for component_pass in spec.component_id_passes.values():
            states = _mesh_states(scene)
            appended_data: typing.Dict[int, typing.Any] = {}
            try:
                _apply_component_mapping(
                    scene, component_pass, materials, appended_data
                )
                path = resolve_output_path(output_dir, component_pass.output)
                _render_still(scene, component_pass.camera, path)
                outputs.append(
                    {
                        "kind": "component_id",
                        "name": component_pass.name,
                        "camera": component_pass.camera,
                        "path": component_pass.output,
                        "sha256": _sha256(path),
                    }
                )
            finally:
                _restore_mesh_states(states, appended_data)
    finally:
        scene.world = original_world
        bpy.data.worlds.remove(world)
        for material in materials.values():
            bpy.data.materials.remove(material)
    return outputs


def _component_manifest(spec: RenderSpec) -> typing.Dict[str, object]:
    result: typing.Dict[str, object] = {}
    for component_pass in spec.component_id_passes.values():
        fallback: typing.Dict[str, object] = {
            "mode": component_pass.fallback.mode
        }
        if component_pass.fallback.rgb is not None:
            fallback["rgb"] = list(component_pass.fallback.rgb)
        result[component_pass.name] = {
            "camera": component_pass.camera,
            "output": component_pass.output,
            "objects": {
                name: list(rgb) for name, rgb in component_pass.objects.items()
            },
            "fallback": fallback,
        }
    return result


def _write_manifest(
    output_dir: pathlib.Path,
    blend_file: pathlib.Path,
    blend_hash: str,
    spec_path: pathlib.Path,
    spec_hash: str,
    spec: RenderSpec,
    outputs: typing.List[typing.Dict[str, object]],
) -> pathlib.Path:
    camera_names = sorted(
        {view.camera for view in spec.beauty_views.values()}
        | {
            component_pass.camera
            for component_pass in spec.component_id_passes.values()
        }
    )
    manifest = {
        "schema_version": 1,
        "blender_version": bpy.app.version_string,
        "source": {
            "path": _source_label(blend_file),
            "sha256": blend_hash,
        },
        "spec": {
            "path": _source_label(spec_path),
            "sha256": spec_hash,
        },
        "scene": spec.scene,
        "dimensions": {
            "width": spec.resolution.width,
            "height": spec.resolution.height,
        },
        "camera_names": camera_names,
        "beauty_views": {
            view.name: {"camera": view.camera, "output": view.output}
            for view in spec.beauty_views.values()
        },
        "component_mapping": _component_manifest(spec),
        "outputs": outputs,
    }
    output_dir.mkdir(parents=True, exist_ok=True)
    manifest_path = output_dir / "manifest.json"
    payload = json.dumps(manifest, indent=2, sort_keys=True) + "\n"
    temporary_name = ""
    try:
        with tempfile.NamedTemporaryFile(
            "w",
            encoding="utf-8",
            dir=output_dir,
            prefix=".manifest.",
            suffix=".tmp",
            delete=False,
        ) as stream:
            temporary_name = stream.name
            stream.write(payload)
        pathlib.Path(temporary_name).replace(manifest_path)
    finally:
        if temporary_name:
            pathlib.Path(temporary_name).unlink(missing_ok=True)
    return manifest_path


def _validate_destinations(output_dir: pathlib.Path, spec: RenderSpec) -> None:
    if output_dir.exists() and not output_dir.is_dir():
        raise RenderError(f"output directory is not a directory: {output_dir}")
    output_dir.mkdir(parents=True, exist_ok=True)
    for output in spec.output_names():
        resolve_output_path(output_dir, output)


def main() -> None:
    """Parse Blender script arguments and render one complete packet."""
    parser = argparse.ArgumentParser(
        description="Render a deterministic review packet with pinned Blender."
    )
    parser.add_argument("--blend-file", required=True, type=pathlib.Path)
    parser.add_argument("--spec", required=True, type=pathlib.Path)
    parser.add_argument("--output-dir", required=True, type=pathlib.Path)
    options = parser.parse_args(_script_args())
    blend_file = _resolve_cli_path(options.blend_file)
    spec_path = _resolve_cli_path(options.spec)
    output_dir = _resolve_cli_path(options.output_dir)
    spec = load_spec(spec_path)
    _validate_destinations(output_dir, spec)
    blend_hash = _sha256(blend_file)
    spec_hash = _sha256(spec_path)
    _open_source(blend_file)
    scene = _validated_scene(spec)
    outputs = _render_beauty(scene, spec, output_dir)
    outputs.extend(_render_component_ids(scene, spec, output_dir))
    manifest = _write_manifest(
        output_dir,
        blend_file,
        blend_hash,
        spec_path,
        spec_hash,
        spec,
        outputs,
    )
    print(f"render packet: {manifest}")


if __name__ == "__main__":
    try:
        main()
    except (OSError, RenderError, SpecError) as error:
        raise RuntimeError(f"render_packet: {error}") from error
