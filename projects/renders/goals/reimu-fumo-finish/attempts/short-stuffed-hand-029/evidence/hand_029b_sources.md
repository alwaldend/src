# Hand 029b reproducible source evidence

Retained-candidate review is bound to hand_029b_candidate.blend SHA256
9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d.
This winding-corrected variant starts from frozen head_028_candidate.blend
c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6.
Run writer then renderer in pinned Blender 5.2.1 LTS background mode with
factory startup, disabled autoexec, four threads and python-exit-code 2.
Use fresh task-local output paths; scripts prohibit overwriting candidates.
Regenerated Blender bytes may differ: hash and review the new artifact.
Earlier hand_029_sources.md records the inherited-winding variant, not this fix.

## hand_029b_draft.py

SHA256 5c549b366396252945de9bd565dc9db669d880fb376515110d8889eadb08ab7a.

```python
"""Two-target distal emergence; import only defines helpers, root owns save."""
import hashlib
import math
import struct

import bpy
from mathutils import Vector

TARGETS = (
    'Sleeve44P L asymmetrically seated stuffed arm insert',
    'Sleeve44P R asymmetrically seated stuffed arm insert',
)
ANCHORS = {
    'L': ((-.060619596391916275, -.011930779553949833, .04682004079222679),
          (-.7327728867530823, -.1988484114408493, -.6507713198661804)),
    'R': ((.0605241023004055, -.011930594220757484, .04671471193432808),
          (.7304894924163818, -.19880202412605286, -.6533473134040833)),
}


def snapshot(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        return ([ev.matrix_world @ v.co for v in mesh.vertices],
                [tuple(p.vertices) for p in mesh.polygons],
                [p.material_index for p in mesh.polygons],
                [bpy.data.materials[m.name] for m in mesh.materials])
    finally:
        ev.to_mesh_clear()


def face_hash(points, faces):
    digest = hashlib.sha256()
    for face in faces:
        digest.update(struct.pack('<I', len(face)))
        for i in face:
            digest.update(struct.pack('<3f', *points[i]))
    return digest.hexdigest()


def bounds(points):
    return [[min(p[k] for p in points), max(p[k] for p in points)]
            for k in range(3)]


def build_hand_029():
    assert bpy.context.scene.frame_current == 1
    assert all(bpy.data.objects[n].visible_get() and not bpy.data.objects[n].hide_render
               for n in TARGETS)
    rig = bpy.data.objects['ReimuFumoRig']
    created, receipts = [], {}
    for side, name in zip(('L', 'R'), TARGETS):
        original = bpy.data.objects[name]
        points, faces, indices, materials = snapshot(original)
        assert len(points) == 338 and len(faces) == 336
        center, axis = [Vector(v) for v in ANCHORS[side]]
        axis.normalize()
        axial = [(p - center).dot(axis) for p in points]
        start, target = -.015, .005
        distal = max(axial)
        assert -.0065 < distal < -.0055
        delta = target - distal
        output = []
        for p, a in zip(points, axial):
            t = min(1., max(0., (a - start) / (distal - start)))
            weight = t*t*(3. - 2.*t)
            output.append(p.copy() if weight == 0. else p + axis*(delta*weight))
        retained = [face for face in faces if max(axial[i] for i in face) <= start]
        assert retained and face_hash(points, retained) == face_hash(output, retained)
        assert all(math.isfinite(v) for p in output for v in p)
        new_name = f'Hand029 {side} rounded stuffed insert'
        assert new_name not in bpy.data.objects
        mesh = bpy.data.meshes.new(new_name + ' mesh')
        volume = sum(output[f[0]].dot(output[f[j]].cross(output[f[j+1]])) / 6.
                     for f in faces for j in range(1, len(f)-1))
        assert abs(volume) > 1e-8
        oriented_faces = faces if volume > 0. else [tuple(reversed(f)) for f in faces]
        mesh.from_pydata(output, [], oriented_faces)
        for material in materials:
            assert material is not None and not material.is_evaluated
            mesh.materials.append(material)
        mesh.update()
        for polygon, material_index in zip(mesh.polygons, indices):
            polygon.use_smooth = True
            polygon.material_index = material_index
        obj = bpy.data.objects.new(new_name, mesh)
        bpy.context.scene.collection.objects.link(obj)
        bone = rig.pose.bones['Arm_' + side]
        deformation = (rig.matrix_world @ bone.matrix
                       @ bone.bone.matrix_local.inverted() @ rig.matrix_world.inverted())
        inverse = deformation.inverted()
        for vertex in mesh.vertices:
            vertex.co = inverse @ vertex.co
        group = obj.vertex_groups.new(name=bone.name)
        group.add(list(range(len(mesh.vertices))), 1., 'REPLACE')
        modifier = obj.modifiers.new('029 compensated Arm attachment', 'ARMATURE')
        modifier.object = rig
        original.hide_render = True
        original.hide_viewport = True
        original.hide_set(True)
        bpy.context.view_layer.update()
        actual = snapshot(obj)[0]
        assert face_hash(actual, retained) == face_hash(points, retained)
        actual_axial = [(p - center).dot(axis) for p in actual]
        assert abs(max(actual_axial) - target) < 1e-7
        created.append(new_name)
        receipts[side] = {
            'original': name, 'new': new_name, 'cuff_center_m': list(center),
            'axis': list(axis), 'proximal_cutoff_m': start,
            'source_axial_range_m': [min(axial), max(axial)],
            'new_axial_range_m': [min(actual_axial), max(actual_axial)],
            'source_bounds_m': bounds(points), 'new_bounds_m': bounds(actual),
            'preserved_face_count': len(retained),
            'preserved_source_indexed_face_coordinates_sha256': face_hash(points, retained),
            'winding_corrected': volume < 0.,
            'source_oriented_extended_volume_m3': volume,
            'written_outward_volume_m3': abs(volume),
            'distal_displacement_m': delta, 'arm_bone': bone.name,
            'method': 'Monotone smooth distal axial extension, unchanged transverse coordinates',
        }
    return {'created_names': created, 'hands': receipts,
            'scope': 'Only two inserts; all sleeve cloth/trim/shoulders preserved',
            'limitations': ['No contact or visual approval', 'Animation unverified']}
```

## hand_029b_writer.py

SHA256 d51504cc93fb426e9b7e7715e15d54bbff6d7b7cfb8f030e4b113663ceaa46a3.

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_028_candidate.blend"
SOURCE_HASH = "c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6"
BASELINE = ROOT / "head_028_candidate.blend"
BASELINE_HASH = "c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6"
OUTPUT = ROOT / "hand_029b_candidate.blend"
RECEIPT = ROOT / "hand_029b_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "Sleeve44P L asymmetrically seated stuffed arm insert",
    "Sleeve44P R asymmetrically seated stuffed arm insert",
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


assert sha(BASELINE) == BASELINE_HASH
assert EXPECTED_TARGETS, "Root has not frozen the helper target set"
assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert sha(SOURCE) == SOURCE_HASH and not OUTPUT.exists() and not RECEIPT.exists()
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.frame_current == 1 and bpy.context.mode == "OBJECT"
assert all(name in scene.objects for name in EXPECTED_TARGETS)
controls = {obj.name: record(obj) for obj in scene.objects
            if obj.type in {"MESH", "CURVE"} and obj.visible_get()
            and not obj.hide_render and obj.name not in EXPECTED_TARGETS}
rig = scene.objects["ReimuFumoRig"]
pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
look = appearance(scene)
helper = ROOT / "hand_029b_draft.py"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
result = scope["build_hand_029"]()
bpy.context.view_layer.update()
assert controls == {name: record(scene.objects[name]) for name in controls}
assert pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert look == appearance(scene), "Material nodes, lighting or color settings changed"
for name in result["created_names"]:
    obj = scene.objects[name]
    assert obj.data.materials and all(
        m is not None and not m.is_evaluated
        and bpy.data.materials.get(m.name) == m
        for m in obj.data.materials), ("Invalid persistent material binding", name)
assert sha(SOURCE) == SOURCE_HASH
scene["candidate_status"] = "Hand029 unreviewed construction study; no stage pass"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
assert sha(BASELINE) == BASELINE_HASH
receipt = {"comparison_baseline": BASELINE.name, "comparison_baseline_sha256": BASELINE_HASH,
           "candidate": OUTPUT.name, "candidate_sha256": sha(OUTPUT),
           "source": SOURCE.name, "source_sha256": SOURCE_HASH,
           "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
           "writer_sha256": sha(Path(__file__)), "helper_sha256": sha(helper),
           "target_names": sorted(EXPECTED_TARGETS), "control_count": len(controls),
           "controls": controls, "rig_pose": pose, "appearance": look,
           "controls_unchanged": True, "rig_pose_unchanged": True,
           "appearance_unchanged": True, "hand_construction": result,
           "limitations": ["No visual, animation, whole-scene technical or final acceptance."]}
assert sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {"controls", "rig_pose", "appearance", "hand_construction"}}))
```

## render_hand_029b.py

SHA256 a4f846d367077cd73031e003b4a7e0afa8871e56da9bfd0ec1ca25ee7ddc64fa.

```python
"""Frozen029 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "hand_029b_candidate.blend"
OUT = ROOT / "hand_029b_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "head_028_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "hand_029b_writer_receipt.json").read_text())
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
    data = bpy.data.cameras.new("Frozen029_" + view)
    camera = bpy.data.objects.new("Frozen029_" + view, data)
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


