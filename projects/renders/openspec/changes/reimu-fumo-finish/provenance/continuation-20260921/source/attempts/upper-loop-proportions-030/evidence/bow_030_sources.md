# Bow030 reproducible source evidence

Candidate bow_030_candidate.blend SHA256
4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b.
Frozen input hand_029b_candidate.blend SHA256
9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d.
Run writer then renderer in pinned Blender5.2.1 LTS background mode with
factory startup, disabled autoexec, four threads and python-exit-code2.
Use fresh task-local outputs; do not overwrite prior candidate or receipt.
Regeneration may produce different Blender bytes and needs a fresh hash and
image review. Sources preserve materials, lighting, world and color settings.

## bow_030_draft.py

SHA256 cd2cdfb396304e573d937676f779688a172ba9fbfe8f23f479c5591085bd7556.

```python
"""Two-cage030 draft. Defines helpers only; no open/save or top-level mutation."""
import hashlib
import math
import struct
from pathlib import Path

import bpy
from mathutils import Vector

TARGETS = ('A154 Left loop macro cage', 'A155 Right loop macro cage')
AFFECTED_GEOMETRY = (
    'A42 Left constructed bow loop',
    'A42 Left narrow gathered loop ruffle',
    'A42 Left white zigzag applique',
    'A42 Left root fold 1', 'A42 Left root fold 2',
    'A42 Right constructed bow loop',
    'A42 Right narrow gathered loop ruffle',
    'A42 Right white zigzag applique',
    'A42 Right root fold 1', 'A42 Right root fold 2',
)
SOURCE_SHA256 = '9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d'


def _geometry(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        return [ev.matrix_world @ v.co for v in mesh.vertices]
    finally:
        ev.to_mesh_clear()


def _bounds(points):
    return [[min(p[k] for p in points), max(p[k] for p in points)]
            for k in range(3)]


def _cage_record(obj):
    return {'matrix_world': [list(row) for row in obj.matrix_world],
            'dimensions': [obj.data.points_u, obj.data.points_v, obj.data.points_w],
            'interpolation': [obj.data.interpolation_type_u,
                              obj.data.interpolation_type_v,
                              obj.data.interpolation_type_w],
            'points': [list(p.co_deform) for p in obj.data.points]}


def _digest(points):
    result = hashlib.sha256()
    for point in points:
        result.update(struct.pack('<3f', *point))
    return result.hexdigest()


def _linear(value, pairs, inverse=False):
    if inverse:
        pairs = [(b, a) for a, b in pairs]
    assert pairs[0][0] <= value <= pairs[-1][0], ('Outside global cage', value)
    for (a, b), (c, d) in zip(pairs, pairs[1:]):
        if a <= value <= c:
            return b + (d - b) * (value - a) / (c - a)
    raise AssertionError(value)


def _rows(obj):
    # U increases in local X; reorder every row from inner toward outer.
    side_left = obj.name == TARGETS[0]
    u_order = (3, 2, 1, 0) if side_left else (0, 1, 2, 3)
    return [[obj.matrix_world @ obj.data.points[w * 8 + u].co_deform
             for u in u_order] for w in range(4)]


def _outboard_weight(x, rows):
    inner = abs(rows[0][1].x)
    outer = abs(rows[0][3].x)
    return min(1., max(0., (abs(x) - inner) / (outer - inner)))


def _interpolated_row_z(x, row):
    return _linear(abs(x), [(abs(p.x), p.z) for p in row])


def build_bow_030():
    source = Path(bpy.data.filepath).resolve()
    assert source.name == 'hand_029b_candidate.blend'
    assert hashlib.sha256(source.read_bytes()).hexdigest() == SOURCE_SHA256
    assert bpy.context.scene.frame_current == 1 and bpy.context.mode == 'OBJECT'
    scene = bpy.context.scene
    assert all(name in scene.objects for name in TARGETS + AFFECTED_GEOMETRY)
    cages = [scene.objects[name] for name in TARGETS]
    for cage in cages:
        assert cage.type == 'LATTICE'
        assert (cage.data.points_u, cage.data.points_v, cage.data.points_w) == (4, 2, 4)
        assert all(getattr(cage.data, 'interpolation_type_' + axis) == 'KEY_LINEAR'
                   for axis in 'uvw')
        assert all(abs(cage.matrix_world[i][j]) < 1e-10
                   for i in range(3) for j in range(3) if i != j)
    assert {obj.name for obj in scene.objects
            if any(m.type == 'LATTICE' and m.object in cages for m in obj.modifiers)} == set(AFFECTED_GEOMETRY)
    global_cage = scene.objects['022 bow proportion cage']
    global_before = _cage_record(global_cage)
    assert (global_cage.data.points_u, global_cage.data.points_v) == (2, 2)
    assert all(getattr(global_cage.data, 'interpolation_type_' + axis) == 'KEY_LINEAR'
               for axis in 'uvw')
    mapping = []
    for w in range(global_cage.data.points_w):
        row = global_cage.data.points[w * 4:w * 4 + 4]
        before = [global_cage.matrix_world @ p.co for p in row]
        after = [global_cage.matrix_world @ p.co_deform for p in row]
        assert max(p.z for p in after) - min(p.z for p in after) < 1e-7
        assert all(abs(a.x-b.x) < 1e-7 and abs(a.y-b.y) < 1e-7
                   for a, b in zip(before, after))
        mapping.append((before[0].z, after[0].z))
    before_geometry = {name: _geometry(scene.objects[name]) for name in AFFECTED_GEOMETRY}
    controls = {obj.name: _digest(_geometry(obj)) for obj in scene.objects
                if obj.type in {'MESH', 'CURVE'} and obj.visible_get()
                and not obj.hide_render and obj.name not in AFFECTED_GEOMETRY}
    cage_before = {obj.name: _cage_record(obj) for obj in cages}
    rows = {obj.name: _rows(obj) for obj in cages}
    upper_names = [name for name in AFFECTED_GEOMETRY
                   if 'constructed bow loop' in name or 'loop ruffle' in name]
    all_upper = [p for name in upper_names for p in before_geometry[name]]
    upper_before = _bounds(all_upper)
    head_bounds = _bounds(_geometry(scene.objects['Hair028 crown and back hood']))
    wh = head_bounds[0][1] - head_bounds[0][0]
    assert .1173 < wh < .1176, ('Unexpected head datum', wh)
    width_target, height_target = 1.50 * wh, .61 * wh
    left_extreme = min(all_upper, key=lambda p: p.x)
    right_extreme = max(all_upper, key=lambda p: p.x)
    sensitivity_x = (_outboard_weight(left_extreme.x, rows[TARGETS[0]])
                     + _outboard_weight(right_extreme.x, rows[TARGETS[1]]))
    inward = (upper_before[0][1] - upper_before[0][0] - width_target) / sensitivity_x
    lowest = min(all_upper, key=lambda p: p.z)
    low_cage = TARGETS[0] if lowest.x < 0 else TARGETS[1]
    low_rows = rows[low_cage]
    z_before_global = _linear(lowest.z, mapping, inverse=True)
    row0_z = _interpolated_row_z(lowest.x, low_rows[0])
    row1_z = _interpolated_row_z(lowest.x, low_rows[1])
    t = (z_before_global - row0_z) / (row1_z - row0_z)
    assert 0. <= t <= 1., ('Lowest witness outside lower cage interval', t)
    sensitivity_z = _outboard_weight(lowest.x, low_rows) * (1. - .75 * t)
    # Lower two cage levels are below the global upper-Z compression pivot.
    assert abs((_linear(row1_z, mapping) - _linear(row0_z, mapping))
               / (row1_z - row0_z) - 1.) < 1e-4
    lowering = (height_target - (upper_before[2][1] - upper_before[2][0])) / sensitivity_z
    assert .012 < inward < .020 and .007 < lowering < .017
    highest_name, highest_index = max(
        ((name, i) for name in upper_names for i in range(len(before_geometry[name]))),
        key=lambda pair: before_geometry[pair[0]][pair[1]].z)
    changes = {}
    for cage in cages:
        inverse = cage.matrix_world.inverted()
        original = cage_before[cage.name]['points']
        root_columns = (2, 3) if cage.name == TARGETS[0] else (0, 1)
        moved = []
        for index, point in enumerate(cage.data.points):
            u, w = index % 4, index // 8
            if u in root_columns:
                continue
            p0 = Vector(original[index])
            world = cage.matrix_world @ p0
            weight = .5 if u in (1, 2) else 1.
            delta_x = inward * weight * (1 if cage.name == TARGETS[0] else -1)
            lower_weight = (1., .25, 0., 0.)[w]
            final_z = _linear(world.z, mapping) - lowering * weight * lower_weight
            moved_world = Vector((world.x + delta_x, world.y,
                                  _linear(final_z, mapping, inverse=True)))
            local = inverse @ moved_world
            # Exact existing local Y, including prior padding, is protected.
            point.co_deform.x = local.x
            point.co_deform.z = local.z if lower_weight else p0.z
            moved.append(index)
        cage.data.update_tag()
        changes[cage.name] = {'moved_control_indices': moved,
                              'fixed_root_columns_u': list(root_columns)}
    bpy.context.view_layer.update()
    after_geometry = {name: _geometry(scene.objects[name]) for name in AFFECTED_GEOMETRY}
    upper_after = _bounds([p for name in upper_names for p in after_geometry[name]])
    assert controls == {name: _digest(_geometry(scene.objects[name])) for name in controls}
    assert global_before == _cage_record(global_cage)
    by_object = {}
    for name in AFFECTED_GEOMETRY:
        a, b = before_geometry[name], after_geometry[name]
        assert len(a) == len(b)
        deltas = [pb-pa for pa, pb in zip(a, b)]
        roots = [i for i, p in enumerate(a) if abs(p.x) <= .00775]
        by_object[name] = {
            'bounds_before_xyz_m': _bounds(a), 'bounds_after_xyz_m': _bounds(b),
            'max_abs_delta_xyz_m': [max(abs(d[k]) for d in deltas) for k in range(3)],
            'root_vertex_count_within_knot_span': len(roots),
            'max_root_displacement_m': max((deltas[i].length for i in roots), default=0.),
        }
    for cage in cages:
        old = cage_before[cage.name]['points']
        new = _cage_record(cage)['points']
        assert all(a[1] == b[1] for a, b in zip(old, new)), 'Cage Y changed'
        for i in range(32):
            if i % 4 in changes[cage.name]['fixed_root_columns_u']:
                assert old[i] == new[i], 'Root column changed'
        changes[cage.name]['before'] = cage_before[cage.name]
        changes[cage.name]['after'] = _cage_record(cage)
    highest_delta = (after_geometry[highest_name][highest_index]
                     - before_geometry[highest_name][highest_index])
    assert abs(upper_after[2][1] - upper_before[2][1]) < .00025
    assert hashlib.sha256(source.read_bytes()).hexdigest() == SOURCE_SHA256
    return {
        'source_sha256': SOURCE_SHA256, 'modified_cages': list(TARGETS),
        'affected_geometry': list(AFFECTED_GEOMETRY), 'cage_changes': changes,
        'evaluated_head_width_m': wh, 'target_upper_span_m': width_target,
        'target_upper_height_m': height_target,
        'upper_bounds_before_xyz_m': upper_before, 'upper_bounds_after_xyz_m': upper_after,
        'upper_span_Wh': (upper_after[0][1] - upper_after[0][0]) / wh,
        'upper_height_Wh': (upper_after[2][1] - upper_after[2][0]) / wh,
        'outer_control_inward_m': inward, 'outer_lower_control_lowering_m': lowering,
        'prediction_method': 'Measured extrema, linear cage weights and inverse shared022 Z mapping; no optimizer',
        'highest_z_change_m': upper_after[2][1] - upper_before[2][1],
        'original_highest_vertex': {'object': highest_name, 'index': highest_index,
                                    'displacement_xyz_m': list(highest_delta)},
        'per_object': by_object, 'protected_control_count': len(controls),
        'protected_evaluated_controls_unchanged': True,
        'global_bow_cage_unchanged': True, 'cage_y_and_two_root_columns_exact': True,
        'limitations': ['Post-cage ruffle Solidify can change evaluated normal offsets',
                        'No material, animation, contact or visual acceptance',
                        'No blend save performed by helper'],
    }
```

## bow_030_writer.py

SHA256 0c1b75341579f49f0fea2674f05e595f9878cd8404132af15974ac16bbb7a566.

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "hand_029b_candidate.blend"
SOURCE_HASH = "9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d"
BASELINE = ROOT / "hand_029b_candidate.blend"
BASELINE_HASH = "9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d"
OUTPUT = ROOT / "bow_030_candidate.blend"
RECEIPT = ROOT / "bow_030_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "A154 Left loop macro cage",
    "A155 Right loop macro cage",
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


AFFECTED = frozenset(
    f"A42 {side} {part}"
    for side in ("Left", "Right")
    for part in ("constructed bow loop", "narrow gathered loop ruffle",
                 "white zigzag applique", "root fold 1", "root fold 2")
)


def cages():
    return {obj.name: {"matrix": [list(row) for row in obj.matrix_world],
                       "points": [list(p.co_deform) for p in obj.data.points]}
            for obj in bpy.context.scene.objects if obj.type == "LATTICE"}


assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert sha(SOURCE) == SOURCE_HASH and not OUTPUT.exists() and not RECEIPT.exists()
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.frame_current == 1 and bpy.context.mode == "OBJECT"
assert all(name in scene.objects for name in EXPECTED_TARGETS | AFFECTED)
controls = {obj.name: record(obj) for obj in scene.objects
            if obj.type in {"MESH", "CURVE"} and obj.visible_get()
            and not obj.hide_render and obj.name not in AFFECTED}
affected_before = {name: record(scene.objects[name]) for name in AFFECTED}
cages_before = cages()
rig = scene.objects["ReimuFumoRig"]
pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
look = appearance(scene)
helper = ROOT / "bow_030_draft.py"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
assert frozenset(scope["AFFECTED_GEOMETRY"]) == AFFECTED
result = scope["build_bow_030"]()
bpy.context.view_layer.update()
assert controls == {name: record(scene.objects[name]) for name in controls}
assert pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert look == appearance(scene), "Appearance changed"
cages_after = cages()
assert cages_before.keys() == cages_after.keys()
for name in cages_before:
    if name not in EXPECTED_TARGETS:
        assert cages_before[name] == cages_after[name], name
    else:
        assert cages_before[name]["matrix"] == cages_after[name]["matrix"], name
        assert cages_before[name]["points"] != cages_after[name]["points"], name
        assert all(a[1] == b[1] for a, b in zip(
            cages_before[name]["points"], cages_after[name]["points"])), name
affected_after = {name: record(scene.objects[name]) for name in AFFECTED}
assert any(affected_before[n] != affected_after[n] for n in AFFECTED)
assert sha(SOURCE) == SOURCE_HASH
scene["candidate_status"] = "Bow030 unreviewed upper-loop proportion study; no stage pass"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
assert sha(SOURCE) == SOURCE_HASH
receipt = {"comparison_baseline": BASELINE.name, "comparison_baseline_sha256": BASELINE_HASH,
           "candidate": OUTPUT.name, "candidate_sha256": sha(OUTPUT),
           "source": SOURCE.name, "source_sha256": SOURCE_HASH,
           "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
           "writer_sha256": sha(Path(__file__)), "helper_sha256": sha(helper),
           "target_names": sorted(EXPECTED_TARGETS),
           "affected_geometry": sorted(AFFECTED),
           "control_count": len(controls), "controls": controls,
           "cages_before": cages_before, "cages_after": cages_after,
           "affected_before": affected_before, "affected_after": affected_after,
           "rig_pose": pose, "appearance": look,
           "controls_unchanged": True, "rig_pose_unchanged": True,
           "appearance_unchanged": True, "bow_construction": result,
           "limitations": ["No visual, animation, whole-scene technical or final acceptance."]}
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {
    "controls", "rig_pose", "appearance", "cages_before", "cages_after",
    "affected_before", "affected_after", "bow_construction"}}))

```

## render_bow_030.py

SHA256 ec01e4f8c3680eebbbc257293ab661b668941339a564b66373be0b8e0cb4885d.

```python
"""Frozen030 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "bow_030_candidate.blend"
OUT = ROOT / "bow_030_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "hand_029b_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "bow_030_writer_receipt.json").read_text())
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


