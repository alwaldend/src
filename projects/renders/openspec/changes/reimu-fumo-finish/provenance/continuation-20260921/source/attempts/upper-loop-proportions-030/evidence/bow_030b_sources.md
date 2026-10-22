# Bow030b winding repair sources

Frozen result bow_030b_candidate.blend SHA256
 d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6.
Source bow_030_candidate.blend SHA256
4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b.
Only the inherited inward left loop is flipped using native Mesh.flip_normals.
Run writer then renderer with pinned5.2.1 LTS background/factory-startup,
disable-autoexec, four threads and python-exit-code2. Fresh output paths are
required. Reproduced Blender bytes require a fresh hash and image review.
The original030 proportion sources are in bow_030_sources.md.

## bow_030b_writer.py

SHA256 99896dcaeafcb5af7ea90b05034bfe26ed1627eca7e585e9c77d04456072efac.

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "bow_030_candidate.blend"
SOURCE_HASH = "4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b"
BASELINE = ROOT / "bow_030_candidate.blend"
BASELINE_HASH = "4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b"
OUTPUT = ROOT / "bow_030b_candidate.blend"
RECEIPT = ROOT / "bow_030b_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "A42 Left constructed bow loop",
})


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def value(data):
    if data is None or isinstance(data, (str, int, float, bool)):
        return data
    try:
        return list(data)
    except TypeError:
        return getattr(data, "name", str(type(data)))


def nodes(tree):
    if tree is None:
        return None
    output = []
    for node in tree.nodes:
        row = {"name": node.name, "type": node.bl_idname,
               "inputs": [(socket.name, value(socket.default_value)) for socket in node.inputs if hasattr(socket, "default_value")],
               "settings": {key: value(getattr(node, key)) for key in
                            ("operation", "blend_type", "noise_dimensions", "normalize") if hasattr(node, key)}}
        if hasattr(node, "color_ramp"):
            row["ramp"] = {"interpolation": node.color_ramp.interpolation,
                           "elements": [(item.position, list(item.color)) for item in node.color_ramp.elements]}
        output.append(row)
    return {"nodes": output, "links": [(link.from_node.name, link.from_socket.identifier,
                                         link.to_node.name, link.to_socket.identifier) for link in tree.links]}


def appearance(scene):
    return {"materials": {m.name: {"diffuse": list(m.diffuse_color), "nodes": nodes(m.node_tree)} for m in bpy.data.materials},
            "lights": {o.name: {"matrix": [list(row) for row in o.matrix_world],
                                  "energy": o.data.energy, "color": list(o.data.color),
                                  "type": o.data.type, "hide": o.hide_render}
                       for o in scene.objects if o.type == "LIGHT"},
            "world": {"name": scene.world.name, "color": list(scene.world.color),
                      "nodes": nodes(scene.world.node_tree)},
            "view": {key: getattr(scene.view_settings, key) for key in ("view_transform", "look", "exposure", "gamma")}}


def record(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    digest = hashlib.sha256()
    try:
        digest.update(struct.pack("<3I", len(mesh.vertices), len(mesh.edges), len(mesh.polygons)))
        for vertex in mesh.vertices:
            world = ev.matrix_world @ vertex.co
            assert all(math.isfinite(v) for v in world)
            digest.update(struct.pack("<3f", *world))
        for face in mesh.polygons:
            digest.update(struct.pack("<2I", len(face.vertices), face.material_index))
            for index in face.vertices:
                digest.update(struct.pack("<I", index))
        return {"geometry": digest.hexdigest(), "materials": [m.name if m else None for m in mesh.materials],
                "visibility": [obj.hide_render, obj.hide_viewport, obj.hide_get()],
                "parent": obj.parent.name if obj.parent else None,
                "modifiers": [(m.name, m.type) for m in obj.modifiers]}
    finally:
        ev.to_mesh_clear()


def point_record(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        return [list(ev.matrix_world @ v.co) for v in mesh.vertices]
    finally:
        ev.to_mesh_clear()


def volume(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        mesh.calc_loop_triangles()
        points = [ev.matrix_world @ v.co for v in mesh.vertices]
        return sum(points[t.vertices[0]].dot(
            points[t.vertices[1]].cross(points[t.vertices[2]])) / 6.
            for t in mesh.loop_triangles)
    finally:
        ev.to_mesh_clear()


assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert sha(SOURCE) == SOURCE_HASH and not OUTPUT.exists() and not RECEIPT.exists()
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.frame_current == 1 and bpy.context.mode == "OBJECT"
target = scene.objects["A42 Left constructed bow loop"]
assert target.type == "MESH" and target.data.users == 1
assert [m.type for m in target.modifiers] == ["LATTICE", "ARMATURE", "LATTICE"]
controls = {obj.name: record(obj) for obj in scene.objects
            if obj.type in {"MESH", "CURVE"} and obj.visible_get()
            and not obj.hide_render and obj.name not in EXPECTED_TARGETS}
look = appearance(scene)
rig = scene.objects["ReimuFumoRig"]
pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
points_before = point_record(target)
faces_before = [(list(f.vertices), f.material_index) for f in target.data.polygons]
weights_before = [[(g.group, g.weight) for g in v.groups] for v in target.data.vertices]
volume_before = volume(target)
assert volume_before < -1e-9
target.data.flip_normals()
target.data.update()
bpy.context.view_layer.update()
assert points_before == point_record(target), "Winding repair moved evaluated points"
assert all(set(indices) == set(face.vertices) and material == face.material_index
           for (indices, material), face in zip(faces_before, target.data.polygons))
assert weights_before == [[(g.group, g.weight) for g in v.groups] for v in target.data.vertices]
assert controls == {name: record(scene.objects[name]) for name in controls}
assert pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert look == appearance(scene)
volume_after = volume(target)
assert volume_after > 1e-9
scene["candidate_status"] = "Bow030b left loop winding correction; no stage pass"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
assert sha(SOURCE) == SOURCE_HASH
receipt = {"candidate": OUTPUT.name, "candidate_sha256": sha(OUTPUT),
           "source": SOURCE.name, "source_sha256": SOURCE_HASH,
           "writer_sha256": sha(Path(__file__)),
           "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
           "target_names": sorted(EXPECTED_TARGETS), "control_count": len(controls),
           "controls": controls, "rig_pose": pose, "appearance": look,
           "controls_unchanged": True, "rig_pose_unchanged": True,
           "appearance_unchanged": True, "evaluated_points_unchanged": True,
           "material_regions_unchanged": True, "weights_unchanged": True,
           "volume_before_m3": volume_before, "volume_after_m3": volume_after,
           "flipped_face_count": len(faces_before),
           "limitations": ["Evaluated triangulation/contact carry requires independent audit.",
                          "No visual, animation or full-stage acceptance."]}
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {
    "controls", "rig_pose", "appearance"}}))

```

## render_bow_030b.py

SHA256 99f614de764d8085fc2b95b101a972be19142dae09beff931af5018108181691.

```python
"""Frozen030 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "bow_030b_candidate.blend"
OUT = ROOT / "bow_030b_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "bow_030_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "bow_030b_writer_receipt.json").read_text())
assert writer["candidate_sha256"] == frozen
baseline = json.loads(BASELINE.read_text())
contract = json.loads(CONTRACT.read_text())
assert sha(CONTRACT) == baseline["contract_sha256"]
OUT.mkdir(exist_ok=False)
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
scene.frame_set(1)
assert {key: getattr(scene.view_settings, key) for key in baseline["settings"]["view_settings"]} == baseline["settings"]["view_settings"]
scene.render.engine = "BLENDER_EEVEE"
scene.eevee.taa_render_samples = baseline["settings"]["taa_render_samples"]
scene.render.resolution_x, scene.render.resolution_y, scene.render.resolution_percentage = contract["camera"]["resolution"]
scene.render.image_settings.file_format = "PNG"
assert {key: getattr(scene.render.image_settings, key) for key in baseline["settings"]["image_format"]} == baseline["settings"]["image_format"]
assert scene.render.film_transparent == baseline["settings"]["film_transparent"]
receipt = {"candidate": SOURCE.name, "candidate_sha256": frozen,
           "contract_sha256": sha(CONTRACT), "baseline_receipt_sha256": sha(BASELINE),
           "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
           "settings": baseline["settings"], "renders": {},
           "purpose": "Same-settings actual-node construction comparison, no final acceptance"}
started = time.monotonic()
for view in ("front", "side", "three_quarter", "rear", "three_quarter_mirror"):
    spec = contract["fixed_views"][view]
    data = bpy.data.cameras.new("Frozen030_" + view)
    camera = bpy.data.objects.new("Frozen030_" + view, data)
    scene.collection.objects.link(camera)
    camera.location = spec["location_m"]
    camera.rotation_euler = spec["rotation_euler_rad"]
    data.type = contract["camera"]["projection"]
    data.ortho_scale = contract["camera"]["ortho_scale_m"]
    scene.camera = camera
    path = OUT / ("candidate_" + view + ".png")
    scene.render.filepath = str(path)
    bpy.ops.render.render(write_still=True)
    receipt["renders"][view] = {"path": path.name, "sha256": sha(path), "camera": spec}
    (OUT / "render_receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
assert sha(SOURCE) == frozen
receipt["candidate_preserved"] = True
receipt["elapsed_seconds"] = time.monotonic() - started
(OUT / "render_receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps(receipt))


```


