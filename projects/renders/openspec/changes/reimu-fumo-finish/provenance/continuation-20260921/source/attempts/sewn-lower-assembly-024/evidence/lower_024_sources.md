# Reproduction of rejected024

Reconstruct retained023 with the preceding evidence, then place these sources
beside it under ignored out/reimu_fumo_finish/desktop_astra. Execute the root
writer then renderer through pinned //tools/blender:blender in background,
factory-startup, disabled-autoexec, task-local TMPDIR/cache. Never overwrite
old candidates. Resulting Blender bytes may differ across serialization;
review any regenerated candidate against its actual hash, not this receipt.

## lower_024_writer.py

```python
"""Root-only immutable-input lower assembly writer."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "body_023_candidate.blend"
SOURCE_HASH = "61c5efe89e833f8c79b5327a439cb2e3688113a927c8dc094a98ed4441a718f1"
OUTPUT = ROOT / "lower_024_candidate.blend"
RECEIPT = ROOT / "lower_024_writer_receipt.json"
TARGETS = {
    "Skirt022 soft hem ruffle",
    "Left short hidden leg root", "Right short hidden leg root",
    "Left black stuffed foot pod", "Right black stuffed foot pod",
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def record(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    digest = hashlib.sha256()
    try:
        digest.update(struct.pack("<3I", len(mesh.vertices), len(mesh.edges), len(mesh.polygons)))
        for vertex in mesh.vertices:
            world = evaluated.matrix_world @ vertex.co
            assert all(math.isfinite(v) for v in world)
            digest.update(struct.pack("<3f", *world))
        for polygon in mesh.polygons:
            digest.update(struct.pack("<2I", len(polygon.vertices), polygon.material_index))
            for index in polygon.vertices:
                digest.update(struct.pack("<I", index))
        return {"geometry": digest.hexdigest(), "materials": [m.name if m else None for m in mesh.materials],
                "visibility": [obj.hide_render, obj.hide_viewport, obj.hide_get()],
                "parent": obj.parent.name if obj.parent else None,
                "modifiers": [(m.name, m.type) for m in obj.modifiers]}
    finally:
        evaluated.to_mesh_clear()


assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert sha(SOURCE) == SOURCE_HASH
assert not OUTPUT.exists() and not RECEIPT.exists()
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.frame_current == 1 and bpy.context.mode == "OBJECT"
assert all(scene.objects[name].visible_get() and not scene.objects[name].hide_render for name in TARGETS)
controls = {obj.name: record(obj) for obj in scene.objects
            if obj.type in {"MESH", "CURVE"} and obj.visible_get()
            and not obj.hide_render and obj.name not in TARGETS}
rig = scene.objects["ReimuFumoRig"]
rig_pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
helper = ROOT / "lower_024_draft.py"
namespace = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), namespace)
result = namespace["build_lower_024"]()
bpy.context.view_layer.update()
assert controls == {name: record(scene.objects[name]) for name in controls}
assert rig_pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert all(scene.objects[name].hide_render and scene.objects[name].hide_get() for name in TARGETS)
assert sha(SOURCE) == SOURCE_HASH
scene["candidate_status"] = "Lower024 intermediate construction study, not accepted"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
receipt = {
    "candidate": OUTPUT.name, "candidate_sha256": sha(OUTPUT),
    "source": SOURCE.name, "source_sha256": SOURCE_HASH,
    "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
    "writer_sha256": sha(Path(__file__)), "helper_sha256": sha(helper),
    "target_names": sorted(TARGETS), "control_count": len(controls),
    "controls": controls, "rig_pose": rig_pose,
    "controls_unchanged": True, "rig_pose_unchanged": True,
    "lower_assembly": result,
    "limitations": ["No visual, animated-rig, final-material, or whole-goal acceptance."]}
assert sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {"controls", "rig_pose", "lower_assembly"}}))

```

## lower_024_draft.py

```python
"""Unexecuted lower-assembly draft. Root calls build_lower_024() and saves.

The red skirt and all non-target geometry remain unchanged. The hem follows
its currently evaluated edge, including both retained proportion lattices.
"""

import bisect
import hashlib
import math
import struct
from collections import defaultdict
from pathlib import Path

import bmesh
import bpy
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "body_023_candidate.blend"
SOURCE_HASH = "61c5efe89e833f8c79b5327a439cb2e3688113a927c8dc094a98ed4441a718f1"
SKIRT = "Skirt022 joined gathered panels"
TARGETS = (
    "Skirt022 soft hem ruffle",
    "Left black stuffed foot pod", "Right black stuffed foot pod",
    "Left short hidden leg root", "Right short hidden leg root",
)
NEW_NAMES = (
    "Lower024 left stuffed leg and toe", "Lower024 right stuffed leg and toe",
    "Lower024 standing gathered hem",
)
COLLECTION_NAME = "Lower024 sewn seated assembly"
WH = 0.1165
TOE_DIAMETER = 0.0250
TOE_SEPARATION = 0.0652
POD_DEPTH = 0.0447
HEM_SPAN = 0.1241
TAU = 2 * math.pi


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _geometry(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        return ([ev.matrix_world @ v.co for v in mesh.vertices],
                [tuple(p.vertices) for p in mesh.polygons])
    finally:
        ev.to_mesh_clear()


def _bounds(points):
    return [[min(p[i] for p in points), max(p[i] for p in points)] for i in range(3)]


def _record(obj):
    points, faces = _geometry(obj)
    value = hashlib.sha256()
    for point in points:
        assert all(math.isfinite(v) for v in point)
        value.update(struct.pack("<3f", *point))
    for face in faces:
        value.update(struct.pack("<I", len(face)))
        for index in face:
            value.update(struct.pack("<I", index))
    return {"world_geometry_sha256": value.hexdigest(),
            "hide_render": obj.hide_render, "hide_set": obj.hide_get(),
            "modifiers": [(m.name, m.type) for m in obj.modifiers],
            "materials": [m.name if m else None for m in obj.data.materials],
            "parent": obj.parent.name if obj.parent else None}


def _tree(objects):
    points, faces = [], []
    for obj in objects:
        xyz, polygons = _geometry(obj)
        offset = len(points)
        points.extend(xyz)
        faces.extend(tuple(i + offset for i in face) for face in polygons)
    return BVHTree.FromPolygons(points, faces)


def _smooth(t):
    t = max(0.0, min(1.0, t))
    return t * t * (3 - 2 * t)


def _evaluated_hem(obj):
    """Read the final centerline on a temporary copy; never change obj."""
    copy = obj.copy()
    bpy.context.scene.collection.objects.link(copy)
    try:
        assert sum(m.type == "SOLIDIFY" for m in copy.modifiers) == 1
        for modifier in copy.modifiers:
            if modifier.type == "SOLIDIFY":
                modifier.show_viewport = modifier.show_render = False
        bpy.context.view_layer.update()
        points, faces = _geometry(copy)
        incidence = defaultdict(int)
        for face in faces:
            for a, b in zip(face, face[1:] + face[:1]):
                incidence[tuple(sorted((a, b)))] += 1
        adjacency = defaultdict(list)
        for (a, b), count in incidence.items():
            if count == 1:
                adjacency[a].append(b)
                adjacency[b].append(a)
        assert adjacency and all(len(v) == 2 for v in adjacency.values())
        unseen, loops = set(adjacency), []
        while unseen:
            start = min(unseen)
            indices, previous, current = [], None, start
            while current not in indices:
                indices.append(current)
                nxt = next(i for i in adjacency[current] if i != previous)
                previous, current = current, nxt
            assert current == start
            unseen.difference_update(indices)
            loops.append([points[i] for i in indices])
        assert len(loops) == 2, "Expected waist and hem boundaries"
        hem = min(loops, key=lambda row: sum(p.z for p in row) / len(row))
        center_y = (min(p.y for p in hem) + max(p.y for p in hem)) / 2
        hem.sort(key=lambda p: math.atan2(p.x, center_y - p.y) % TAU)
        assert len(hem) >= 160 and max(p.z for p in hem) < 0.040
        return hem
    finally:
        bpy.data.objects.remove(copy, do_unlink=True)
        bpy.context.view_layer.update()


def _pose_inverse(rig, bone_name):
    bone = rig.pose.bones[bone_name]
    return (rig.matrix_world @ bone.matrix @ bone.bone.matrix_local.inverted()
            @ rig.matrix_world.inverted()).inverted()


def _mesh(name, world, faces, indices, materials, bone_name, rig, collection):
    inverse = _pose_inverse(rig, bone_name)
    mesh = bpy.data.meshes.new(name)
    mesh.from_pydata([inverse @ p for p in world], [], faces)
    for material in materials:
        mesh.materials.append(material)
    mesh.update()
    for polygon, index in zip(mesh.polygons, indices):
        polygon.material_index = index
        polygon.use_smooth = True
    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    obj.parent = rig
    obj.matrix_world = Matrix.Identity(4)
    group = obj.vertex_groups.new(name=bone_name)
    group.add(list(range(len(world))), 1.0, "REPLACE")
    armature = obj.modifiers.new("Inherited seated bone attachment", "ARMATURE")
    armature.object = rig
    obj["status"] = "unreviewed construction candidate; motion unverified"
    return obj


def _leg(sign, name, ground, skirt_tree, cream, black, rig, collection):
    radius = TOE_DIAMETER / 2
    yaw = math.radians(8.0)
    back = Vector((-sign * math.sin(yaw), math.cos(yaw), 0))
    sideways = Vector((math.cos(yaw), sign * math.sin(yaw), 0))
    front_center = Vector((sign * TOE_SEPARATION / 2, -0.0815, ground + radius))
    # The almost planar black end panel shares its seam with the cream tube.
    # (depth from front center, radius fraction, material for arriving faces).
    rings = (
        (0.00005, 0.35, 1), (0.00020, 0.70, 1),
        (0.00070, 0.92, 1), (0.00140, 1.00, 1),
        (0.00260, 1.00, 1), (0.00450, 0.985, 1),
        (0.00520, 0.985, 0), (0.0080, 1.00, 0),
        (0.0150, 0.97, 0), (0.0240, 0.94, 0),
        (0.0330, 0.86, 0), (0.0395, 0.76, 0),
        (0.0430, 0.54, 0), (0.0444, 0.24, 0),
    )
    around = 96
    vertices, faces, materials = [front_center], [], []
    compression = []
    for depth, r, _ in rings:
        center = front_center + back * depth
        center.z += 0.0028 * _smooth((depth - 0.017) / (POD_DEPTH - 0.017))
        for j in range(around):
            angle = TAU * j / around
            point = center + sideways * (r * radius * math.cos(angle))
            point.z += r * radius * math.sin(angle)
            point.z = max(ground, point.z)
            original_z = point.z
            # Local stuffing compression under the unchanged red cloth.
            if point.z > center.z + 0.002 and depth >= 0.0052:
                hit, _, _, _ = skirt_tree.ray_cast(
                    Vector((point.x, point.y, ground - 0.002)),
                    Vector((0, 0, 1)), 0.10)
                if hit is not None:
                    point.z = min(point.z, hit.z - 0.0008)
                    assert point.z > center.z + 0.001, "Skirt crushes a leg crown"
            compression.append(original_z - point.z)
            vertices.append(point)
    for j in range(around):
        faces.append((0, 1 + j, 1 + (j + 1) % around))
        materials.append(1)
    for row in range(len(rings) - 1):
        for j in range(around):
            a = 1 + row * around + j
            b = 1 + row * around + (j + 1) % around
            faces.append((a, a + around, b + around, b))
            materials.append(rings[row + 1][2])
    back_center = front_center + back * POD_DEPTH
    back_center.z += 0.0028
    back_index = len(vertices)
    vertices.append(back_center)
    for j in range(around):
        a = 1 + (len(rings) - 1) * around + j
        b = 1 + (len(rings) - 1) * around + (j + 1) % around
        faces.append((a, back_index, b))
        materials.append(0)
    assert max(compression) <= 0.0040, max(compression)
    obj = _mesh(name, vertices, faces, materials, (cream, black),
                "Leg_L" if sign < 0 else "Leg_R", rig, collection)
    # Consistent winding of a closed stuffed surface, without altering shape.
    bm = bmesh.new()
    bm.from_mesh(obj.data)
    bmesh.ops.recalc_face_normals(bm, faces=list(bm.faces))
    assert all(edge.is_manifold for edge in bm.edges)
    bm.to_mesh(obj.data)
    bm.free()
    return obj, {"front_center_m": list(front_center), "toe_diameter_m": TOE_DIAMETER,
                 "axis_depth_m": POD_DEPTH, "yaw_degrees": sign * 8.0,
                 "maximum_crown_compression_m": max(compression),
                 "black_wrap_depth_m": 0.0045}


# Unequal seam lengths, with twelve gathers on approximately half the loop.
FOLD_WEIGHTS = (
    1.05, .84, 1.18, .94, .89, 1.12, .97, 1.19, .86, 1.03, .91, 1.07,
    .91, 1.15, .95, 1.04, .87, 1.16, 1.00, .88, 1.11, .93, 1.08, .96,
)
# Per fold: (fraction, outward opening, height relative to sewn edge).
FOLD_SHAPE = ((0, .44, .0010), (.18, .75, -.0020), (.44, 1.0, -.0062),
              (.66, .80, -.0030), (.83, .44, .0021), (1, .44, .0010))


def _ruffle(hem, leg_tree, ground, white, rig, collection):
    length = [0.0]
    for a, b in zip(hem, hem[1:] + hem[:1]):
        length.append(length[-1] + (b - a).length)
    total = length[-1]
    folds = [0.0]
    for weight in FOLD_WEIGHTS:
        folds.append(folds[-1] + weight / sum(FOLD_WEIGHTS))
    around, across = 768, 13
    seam, offsets, heights, fold_ids = [], [], [], []
    for j in range(around):
        fraction = j / around
        distance = fraction * total
        i = min(bisect.bisect_right(length, distance) - 1, len(hem) - 1)
        a, b = hem[i], hem[(i + 1) % len(hem)]
        s = (distance - length[i]) / (length[i + 1] - length[i])
        root = a.lerp(b, s)
        tangent = b - a
        outward = Vector((tangent.y, -tangent.x, 0)).normalized()
        k = min(bisect.bisect_right(folds, fraction) - 1, len(FOLD_WEIGHTS) - 1)
        phase = (fraction - folds[k]) / (folds[k + 1] - folds[k])
        left, right = next((a, b) for a, b in zip(FOLD_SHAPE, FOLD_SHAPE[1:])
                           if phase <= b[0])
        t = _smooth((phase - left[0]) / (right[0] - left[0]))
        opening = left[1] * (1 - t) + right[1] * t
        height = left[2] * (1 - t) + right[2] * t
        height *= (.88, 1.04, .95, 1.08, .91, 1.0)[k % 6]
        seam.append(root)
        offsets.append(Vector((outward.x * 0.0100 * opening,
                               outward.y * 0.0058 * opening, 0)))
        heights.append(height)
        fold_ids.append(k)
    # Solve the lateral spread once from the retained seam; never scale roots.
    low, high = 0.0, 3.0
    for _ in range(24):
        scale = (low + high) / 2
        xs = [p.x + scale * d.x for p, d in zip(seam, offsets)]
        if max(xs) - min(xs) < HEM_SPAN:
            low = scale
        else:
            high = scale
    lateral_scale = (low + high) / 2
    vertices, lift = [], []
    for row in range(across):
        u = row / (across - 1)
        for j in range(around):
            offset = offsets[j].copy()
            offset.x *= lateral_scale
            point = seam[j] + u * offset
            point.z += heights[j] * u**1.15 + 0.0008 * 4 * u * (1 - u)
            raw_z = point.z
            if row > 0:
                point.z = max(point.z, ground + 0.00035)
                hit, _, _, _ = leg_tree.ray_cast(
                    Vector((point.x, point.y, 0.08)), Vector((0, 0, -1)), 0.10)
                if hit is not None:
                    # Standing folds compress where their cloth rests on cream.
                    point.z = max(point.z, hit.z + 0.00065)
            lift.append(point.z - raw_z)
            vertices.append(point)
    faces = []
    for row in range(across - 1):
        for j in range(around):
            a = row * around + j
            b = row * around + (j + 1) % around
            faces.append((a, a + around, b + around, b))
    obj = _mesh(NEW_NAMES[2], vertices, faces, [0] * len(faces),
                (white,), "Body", rig, collection)
    # Dense smooth-shaded cloth retains its exact sewn edge; no subdivision
    # moves this root away from the evaluated red hem.
    solid = obj.modifiers.new("Folded cotton thickness", "SOLIDIFY")
    solid.thickness = 0.0005
    solid.offset = 0.0
    return obj, {
        "evaluated_red_seam_count": len(hem), "resampled_seam_count": around,
        "seam_perimeter_m": total, "gather_count": len(FOLD_WEIGHTS),
        "fold_weights": FOLD_WEIGHTS, "fold_cross_sections": FOLD_SHAPE,
        "lateral_spread_factor": lateral_scale, "midsurface_hem_span_m": HEM_SPAN,
        "maximum_support_lift_m": max(lift),
        "seam_witnesses_m": [list(p) for p in seam[::48]],
        "seam_method": "evaluated temporary midsurface copy; source unmodified",
    }


def build_lower_024():
    assert Path(bpy.data.filepath).resolve() == SOURCE.resolve()
    assert _sha(SOURCE) == SOURCE_HASH
    assert bpy.context.scene.frame_current == 1 and bpy.context.mode == "OBJECT"
    scene = bpy.context.scene
    assert all(name in scene.objects for name in TARGETS + (SKIRT, "Review floor"))
    assert all(not scene.objects[name].hide_render for name in TARGETS)
    assert all(name not in bpy.data.objects for name in NEW_NAMES)
    assert COLLECTION_NAME not in bpy.data.collections
    rig = scene.objects[SKIRT].parent
    assert rig and rig.type == "ARMATURE"
    assert all(name in rig.pose.bones for name in ("Body", "Leg_L", "Leg_R"))
    assert [g.name for g in scene.objects[TARGETS[1]].vertex_groups] == ["Leg_L"]
    assert [g.name for g in scene.objects[TARGETS[2]].vertex_groups] == ["Leg_R"]
    controls = {obj.name: _record(obj) for obj in scene.objects
                if obj.type in {"MESH", "CURVE"} and obj.name not in TARGETS}
    pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
    old_state = {name: (scene.objects[name].hide_render, scene.objects[name].hide_get())
                 for name in TARGETS}
    floor_top = _bounds(_geometry(scene.objects["Review floor"])[0])[2][1]
    ground = floor_top + 0.0001
    cream = scene.objects[TARGETS[3]].data.materials[0]
    black = scene.objects[TARGETS[1]].data.materials[0]
    white = scene.objects[TARGETS[0]].data.materials[0]
    hem = _evaluated_hem(scene.objects[SKIRT])
    red_tree = _tree([scene.objects[SKIRT]])
    collection = bpy.data.collections.new(COLLECTION_NAME)
    scene.collection.children.link(collection)
    created, foot_metrics = [], {}
    try:
        for sign, name in zip((-1, 1), NEW_NAMES[:2]):
            obj, metrics = _leg(sign, name, ground, red_tree, cream, black, rig, collection)
            created.append(obj)
            foot_metrics[name] = metrics
        bpy.context.view_layer.update()
        leg_tree = _tree(created)
        ruffle, hem_metrics = _ruffle(hem, leg_tree, ground, white, rig, collection)
        created.append(ruffle)
        bpy.context.view_layer.update()
        bounds = {obj.name: _bounds(_geometry(obj)[0]) for obj in created}
        for name in NEW_NAMES[:2]:
            assert abs(bounds[name][2][0] - ground) < 0.0002, bounds[name]
            assert 0.037 <= bounds[name][1][1] - bounds[name][1][0] <= 0.049
        hem_width = bounds[NEW_NAMES[2]][0][1] - bounds[NEW_NAMES[2]][0][0]
        assert abs(hem_width - HEM_SPAN) < 0.0010, hem_width
        assert bounds[NEW_NAMES[2]][2][0] > floor_top
        assert controls == {name: _record(scene.objects[name]) for name in controls}
        assert pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
        for name in TARGETS:
            scene.objects[name].hide_render = True
            scene.objects[name].hide_set(True)
        bpy.context.view_layer.update()
        assert controls == {name: _record(scene.objects[name]) for name in controls}
        assert _sha(SOURCE) == SOURCE_HASH
        return {
            "status": "unreviewed lower assembly; no save or acceptance",
            "source_sha256": SOURCE_HASH, "helper_sha256": _sha(Path(__file__)),
            "hidden_objects": list(TARGETS), "new_objects": list(NEW_NAMES),
            "controls": controls, "non_target_evaluated_geometry_preserved": True,
            "rig_pose_preserved": True, "head_width_m": WH,
            "floor_top_m": floor_top, "foot_ground_m": ground,
            "foot_center_separation_m": TOE_SEPARATION,
            "bounds_m": bounds, "foot_construction": foot_metrics,
            "hem_construction": hem_metrics,
            "materials_reused": [cream.name, black.name, white.name],
            "limitations": [
                "Clean reopen, intersection review, and final render remain required.",
                "Cap and tube share a seam; the leg rear root is concealed in the seat.",
                "Current pose is compensated; animated motion remains unverified.",
                "Hem span has cotton-thickness tolerance; no visual pass is claimed.",
            ],
        }
    except Exception:
        for name, state in old_state.items():
            scene.objects[name].hide_render = state[0]
            scene.objects[name].hide_set(state[1])
        for obj in list(collection.objects):
            mesh = obj.data
            bpy.data.objects.remove(obj, do_unlink=True)
            if mesh.users == 0:
                bpy.data.meshes.remove(mesh)
        bpy.data.collections.remove(collection)
        bpy.context.view_layer.update()
        raise

```

## render_lower_024.py

```python
"""Frozen-input, fixed-camera construction review; never resave the blend."""

import hashlib
import json
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
WORKSPACE = ROOT.parents[2]
CONTRACT = WORKSPACE / "projects/renders/assets/reimu_fumo/review_contract.json"
contract = json.loads(CONTRACT.read_text())
output = ROOT / "lower_024_review"
output.mkdir(exist_ok=False)
receipt = {
    "blender_version": bpy.app.version_string,
    "build_hash": bpy.app.build_hash.decode(),
    "contract_sha256": hashlib.sha256(CONTRACT.read_bytes()).hexdigest(),
    "engine": "BLENDER_WORKBENCH",
    "purpose": "construction review using original material viewport colors; no final material or rig acceptance",
    "inputs": {},
    "renders": {},
}
for label, file_name, views in (
    ("candidate", "lower_024_candidate.blend", list(contract["fixed_views"])),
):
    path = ROOT / file_name
    frozen_hash = hashlib.sha256(path.read_bytes()).hexdigest()
    receipt["inputs"][label] = {"path": str(path), "sha256": frozen_hash}
    bpy.ops.wm.open_mainfile(filepath=str(path))
    scene = bpy.context.scene
    scene.render.engine = "BLENDER_WORKBENCH"
    scene.display.shading.light = "STUDIO"
    scene.display.shading.studio_light = "paint.sl"
    scene.display.shading.color_type = "MATERIAL"
    scene.display.shading.show_shadows = True
    scene.display.shading.show_cavity = True
    scene.display.shading.cavity_type = "WORLD"
    scene.display.shading.curvature_ridge_factor = 1.0
    scene.display.shading.curvature_valley_factor = 1.0
    scene.display.shading.background_type = "WORLD"
    if scene.world is None:
        scene.world = bpy.data.worlds.new("ReviewNeutralWorld")
    scene.world.color = (0.16, 0.16, 0.16)
    scene.view_settings.view_transform = "Standard"
    scene.render.resolution_x, scene.render.resolution_y, scene.render.resolution_percentage = contract["camera"]["resolution"]
    scene.render.image_settings.file_format = "PNG"
    for view in views:
        spec = contract["fixed_views"][view]
        data = bpy.data.cameras.new("FrozenReview_" + view)
        camera = bpy.data.objects.new("FrozenReview_" + view, data)
        scene.collection.objects.link(camera)
        camera.location = spec["location_m"]
        camera.rotation_euler = spec["rotation_euler_rad"]
        data.type = contract["camera"]["projection"]
        data.ortho_scale = contract["camera"]["ortho_scale_m"]
        scene.camera = camera
        scene.render.filepath = str(output / (label + "_" + view + ".png"))
        bpy.ops.render.render(write_still=True)
        receipt["renders"][label + "_" + view] = {
            "sha256": hashlib.sha256(Path(scene.render.filepath).read_bytes()).hexdigest(),
            "camera": spec,
        }
    assert hashlib.sha256(path.read_bytes()).hexdigest() == frozen_hash
(output / "render_receipt.json").write_text(json.dumps(receipt, indent=2))
print(json.dumps(receipt), flush=True)


```


