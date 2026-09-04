# Reproducing the retained023 input for lower-assembly work

This evidence carries exact023 writer sources forward without modifying the
closed023 record. Reconstruct022b from its closed-attempt evidence first,
including body_022b_writer_receipt.json. Place these scripts beside022b in
ignored out/reimu_fumo_finish/desktop_astra. Use the pinned
//tools/blender:blender background target, factory startup, disabled autoexec,
and task-local TMPDIR/cache. Execute body_023_writer.py then render_body_023.py.
Do not overwrite existing candidates. Blender serialization may produce a
new file hash; bind any regeneration's review to its actual saved bytes.
The observed retained023 hash is
61c5efe89e833f8c79b5327a439cb2e3688113a927c8dc094a98ed4441a718f1.

## body_023_writer.py

```python
"""Sole pinned writer: narrower torso, then reviewed fitted cloth helper."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "body_022b_candidate.blend"
SOURCE_HASH = "96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014"
OUTPUT = ROOT / "body_023_candidate.blend"
RECEIPT = ROOT / "body_023_writer_receipt.json"
DONOR = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/donors/a202/model.blend"
DONOR_HASH = "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5"
COLLAR_TARGETS = {"Left soft collar", "Right soft collar", "Folded yellow cravat"}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def points(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        return [evaluated.matrix_world @ vertex.co for vertex in mesh.vertices]
    finally:
        evaluated.to_mesh_clear()


def bounds(obj):
    xyz = points(obj)
    return [[min(p[i] for p in xyz), max(p[i] for p in xyz)] for i in range(3)]


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


def xscale(z):
    t = min(1.0, max(0.0, (z - 0.028) / 0.023))
    return 1.0 - 0.31 * t * t * (3.0 - 2.0 * t)


def target_x(name, point):
    if name.startswith("Sleeve44P L "):
        return 0.90 * point.x + 0.0065
    if name.startswith("Sleeve44P R "):
        return 0.90 * point.x - 0.0065
    return point.x * xscale(point.z)


assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert sha(SOURCE) == SOURCE_HASH and sha(DONOR) == DONOR_HASH
assert not OUTPUT.exists() and not RECEIPT.exists()
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.name == "Attempt41_Manual_Head_Maquette" and scene.frame_current == 1
assert bpy.context.mode == "OBJECT"
prior = json.loads((ROOT / "body_022b_writer_receipt.json").read_text())
assert prior["candidate_sha256"] == SOURCE_HASH
names = set(prior["roles"]["body"]) - COLLAR_TARGETS
assert all(scene.objects[name].visible_get() and not scene.objects[name].hide_render for name in names)
controls = {obj.name: record(obj) for obj in scene.objects
            if obj.type in {"MESH", "CURVE"} and obj.visible_get()
            and not obj.hide_render and obj.name not in names | COLLAR_TARGETS}
rig = scene.objects["ReimuFumoRig"]
rig_pose = {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
before = {name: points(scene.objects[name]) for name in names}
bodice_before = bounds(scene.objects["Garment42 compact bodice"])

data = bpy.data.lattices.new("023 narrow waist field")
data.points_u = data.points_v = 2
data.points_w = 64
data.interpolation_type_u = data.interpolation_type_v = data.interpolation_type_w = "KEY_LINEAR"
cage = bpy.data.objects.new("023 narrow waist field", data)
scene.collection.objects.link(cage)
base = [[min(p.co_deform[i] for p in data.points), max(p.co_deform[i] for p in data.points)] for i in range(3)]
cage.scale = tuple(extent / (row[1] - row[0]) for extent, row in zip((0.4, 0.4, 0.12), base))
cage.location = (0.0, 0.0, 0.06)
bpy.context.view_layer.update()
inverse = cage.matrix_world.inverted()
for point in data.points:
    world = cage.matrix_world @ point.co_deform
    world.x *= xscale(world.z)
    point.co_deform = inverse @ world
cage.hide_render = True
cage.hide_set(True)
cages = {"body": cage}
for side, offset in (("L", 0.0065), ("R", -0.0065)):
    sleeve_data = bpy.data.lattices.new("023 sleeve root " + side)
    sleeve_data.points_u = sleeve_data.points_v = sleeve_data.points_w = 2
    sleeve_data.interpolation_type_u = sleeve_data.interpolation_type_v = sleeve_data.interpolation_type_w = "KEY_LINEAR"
    sleeve_cage = bpy.data.objects.new("023 sleeve root " + side, sleeve_data)
    scene.collection.objects.link(sleeve_cage)
    sleeve_cage.scale = (0.4, 0.4, 0.12)
    sleeve_cage.location = (0.0, 0.0, 0.06)
    bpy.context.view_layer.update()
    inverse = sleeve_cage.matrix_world.inverted()
    for point in sleeve_data.points:
        world = sleeve_cage.matrix_world @ point.co_deform
        world.x = 0.90 * world.x + offset
        point.co_deform = inverse @ world
    sleeve_cage.hide_render = True
    sleeve_cage.hide_set(True)
    cages[side] = sleeve_cage
for name in sorted(names):
    mod = scene.objects[name].modifiers.new("023 narrow waist, fixed height", "LATTICE")
    key = "L" if name.startswith("Sleeve44P L ") else "R" if name.startswith("Sleeve44P R ") else "body"
    mod.object = cages[key]
bpy.context.view_layer.update()
max_x_error, max_yz_error = 0.0, 0.0
for name in sorted(names):
    after = points(scene.objects[name])
    assert len(before[name]) == len(after)
    for a, b in zip(before[name], after):
        max_x_error = max(max_x_error, abs(b.x - target_x(name, a)))
        max_yz_error = max(max_yz_error, abs(b.y - a.y), abs(b.z - a.z))
assert max_x_error < 0.0002 and max_yz_error < 1e-7, (max_x_error, max_yz_error)
bodice_after = bounds(scene.objects["Garment42 compact bodice"])
assert abs(bodice_after[0][1] - bodice_after[0][0] - 0.0566) < 0.0005

# Root must inspect the helper before this process runs. It must not save.
helper = ROOT / "collar_023_draft.py"
namespace = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), namespace)
collar_receipt = namespace["build_collar_023"]()
bpy.context.view_layer.update()
assert controls == {name: record(scene.objects[name]) for name in controls}
assert rig_pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert sha(SOURCE) == SOURCE_HASH and sha(DONOR) == DONOR_HASH
scene["candidate_status"] = "Body023 intermediate construction candidate; unreviewed"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
receipt = {
    "candidate": OUTPUT.name, "candidate_sha256": sha(OUTPUT),
    "source": SOURCE.name, "source_sha256": SOURCE_HASH,
    "donor_sha256": sha(DONOR), "version": bpy.app.version_string,
    "build_hash": bpy.app.build_hash.decode(), "background": bpy.app.background,
    "writer_script_sha256": sha(Path(__file__)), "collar_script_sha256": sha(helper),
    "width_field": {"upper_scale": 0.69, "full_above_m": 0.051, "none_below_m": 0.028,
                    "sleeve_scale": 0.90, "sleeve_inward_shift_m": 0.0065,
                    "max_x_error_m": max_x_error, "max_yz_error_m": max_yz_error,
                    "objects": sorted(names)},
    "bodice_before_m": bodice_before, "bodice_after_m": bodice_after,
    "protected_control_count": len(controls), "controls": controls,
    "controls_unchanged": True, "rig_pose_unchanged": True,
    "collar": collar_receipt,
    "limitations": ["No visual, final material, or animation acceptance."]}
assert receipt["donor_sha256"] == DONOR_HASH and sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({key: value for key, value in receipt.items() if key != "controls"}), flush=True)

```

## collar_023_draft.py

```python
"""Unexecuted collar/tie helper, called after root's body023 X narrowing.

Defines build_collar_023(); no top-level modeling, save, or receipt write.
Only the three named old garments change visibility. Root owns execution.
"""

import hashlib
import math
import struct
from pathlib import Path

import bpy
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "body_022b_candidate.blend"
SOURCE_HASH = "96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014"
TARGETS = ("Left soft collar", "Right soft collar", "Folded yellow cravat")
NEW_NAMES = ("Collar023 left seated cloth flap",
             "Collar023 right seated cloth flap", "Cravat023 soft gathered tie")
COLLECTION_NAME = "Collar023 surface fitted cloth"
BODICE = "Garment42 compact bodice"
HEAD = "Head_Gusseted_Cushion_020b"
TIE_WIDTH = 0.0200
TIE_HEIGHT = 0.0240
_FRONT_FALLBACKS = []


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _surface(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        points = [evaluated.matrix_world @ vertex.co for vertex in mesh.vertices]
        faces = [tuple(face.vertices) for face in mesh.polygons]
        bounds = [[min(p[i] for p in points), max(p[i] for p in points)]
                  for i in range(3)]
        return BVHTree.FromPolygons(points, faces), bounds
    finally:
        evaluated.to_mesh_clear()


def _digest(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    value = hashlib.sha256()
    try:
        for vertex in mesh.vertices:
            value.update(struct.pack("<3f", *(evaluated.matrix_world @ vertex.co)))
        for polygon in mesh.polygons:
            value.update(struct.pack("<I", len(polygon.vertices)))
            for index in polygon.vertices:
                value.update(struct.pack("<I", index))
        return value.hexdigest()
    finally:
        evaluated.to_mesh_clear()


def _smooth(value):
    value = max(0.0, min(1.0, value))
    return value * value * (3.0 - 2.0 * value)


def _front(tree, x, z, free_edge=False):
    for step in range(9 if free_edge else 1):
        inset = 0.00025 * step
        support_x = math.copysign(max(0.0, abs(x) - inset), x)
        point, normal, _, _ = tree.ray_cast(Vector((support_x, -0.20, z)),
                                            Vector((0, 1, 0)), 0.40)
        if point is not None and normal.y < -0.15:
            if step:
                _FRONT_FALLBACKS.append({"x_m": x, "requested_z_m": z,
                                         "support_x_m": point.x,
                                         "lateral_overhang_m": inset})
            return point
    raise AssertionError(("Missing front support within bounded free edge", x, z, free_edge))


def _neck(tree, x, y):
    point, normal, _, _ = tree.ray_cast(Vector((x, y, 0.04)),
                                        Vector((0, 0, 1)), 0.10)
    assert point is not None and normal.z < -0.20, ("Missing underside", x, y)
    # A thin sewn root lies just outside the evaluated head underside.
    return point + normal * 0.00055, point


def _cloth(name, grid, material, collection, rig):
    rows, columns = len(grid), len(grid[0])
    world = [point for row in grid for point in row]
    # Each inspected old garment has the single Body weight group. Invert
    # only that affine pose when authoring; retain the Body armature modifier.
    bone = rig.pose.bones["Body"]
    deform = (rig.matrix_world @ bone.matrix @ bone.bone.matrix_local.inverted()
              @ rig.matrix_world.inverted())
    inverse = deform.inverted()
    vertices = [inverse @ point for point in world]
    faces = []
    for row in range(rows - 1):
        for column in range(columns - 1):
            a = row * columns + column
            faces.append((a, a + columns, a + columns + 1, a + 1))
    mesh = bpy.data.meshes.new(name)
    mesh.from_pydata(vertices, [], faces)
    mesh.materials.append(material)
    mesh.update()
    # All visible fabric fronts must face -Y, including the mirrored flap.
    if sum(face.normal.y * face.area for face in mesh.polygons) > 0:
        mesh.flip_normals()
    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    obj.parent = rig
    obj.matrix_world = Matrix.Identity(4)
    group = obj.vertex_groups.new(name="Body")
    group.add(list(range(len(vertices))), 1.0, "REPLACE")
    for polygon in mesh.polygons:
        polygon.use_smooth = True
    soften = obj.modifiers.new("Soft sewn panel", "SUBSURF")
    soften.levels = soften.render_levels = 1
    thickness = obj.modifiers.new("Thin filled cloth", "SOLIDIFY")
    thickness.thickness = 0.0006
    thickness.offset = 0.0
    armature = obj.modifiers.new("Inherited Body attachment", "ARMATURE")
    armature.object = rig
    obj["stage"] = "unreviewed fitted garment; no rig or fidelity acceptance"
    return obj


def build_collar_023():
    assert Path(bpy.data.filepath).resolve() == SOURCE.resolve()
    assert _sha(SOURCE) == SOURCE_HASH
    assert bpy.context.scene.frame_current == 1 and bpy.context.mode == "OBJECT"
    scene = bpy.context.scene
    assert all(name in scene.objects for name in TARGETS + (BODICE, HEAD))
    assert all(not scene.objects[name].hide_render for name in TARGETS)
    assert all(name not in bpy.data.objects for name in NEW_NAMES)
    assert COLLECTION_NAME not in bpy.data.collections
    body_tree, body_bounds = _surface(scene.objects[BODICE])
    head_tree, head_bounds = _surface(scene.objects[HEAD])
    width = body_bounds[0][1] - body_bounds[0][0]
    assert 0.054 <= width <= 0.060, ("Call after bodice narrowing", width)
    assert abs(sum(body_bounds[0])) < 0.001
    half = width * 0.5
    top = body_bounds[2][1]
    bottom = body_bounds[2][0]
    neck_y = _front(body_tree, 0.0, top - 0.003).y + 0.0125
    rig = scene.objects[BODICE].parent
    assert rig and rig.type == "ARMATURE" and "Body" in rig.pose.bones
    controls = {obj.name: _digest(obj) for obj in scene.objects
                if obj.type in {"MESH", "CURVE"} and obj.name not in TARGETS}
    old_state = {name: (scene.objects[name].hide_render, scene.objects[name].hide_get())
                 for name in TARGETS}
    white = scene.objects[TARGETS[0]].data.materials[0]
    yellow = scene.objects[TARGETS[2]].data.materials[0]
    collection = bpy.data.collections.new(COLLECTION_NAME)
    scene.collection.children.link(collection)
    roots, chest_samples, created = [], [], []
    try:
        for sign, name in zip((-1, 1), NEW_NAMES[:2]):
            grid = []
            for row in range(21):
                v = row / 20.0
                strip = []
                for column in range(17):
                    u = column / 16.0
                    root_x = sign * half * (0.52 + 0.26 * u)
                    root, hit = _neck(head_tree, root_x, neck_y)
                    lower_x = sign * half * (0.10 + 0.875 * u)
                    lower_z = top - 0.0114 + 0.0035 * u
                    x = root.x * (1.0 - v) + lower_x * v
                    z = root.z * (1.0 - v) + lower_z * v
                    support_v = max(0.5, v)
                    support_x = root.x * (1 - support_v) + lower_x * support_v
                    support_z = root.z * (1 - support_v) + lower_z * support_v
                    front = _front(body_tree, support_x, support_z,
                                   free_edge=u > 0.70) if row else root
                    settle = _smooth(v / 0.5)
                    y = root.y * (1.0 - settle) + (front.y - 0.0008) * settle
                    # Broad soft filling, fading to the sewn perimeter.
                    y -= 0.00045 * 4 * u * (1 - u) * 4 * v * (1 - v)
                    point = Vector((x, y, z))
                    strip.append(point)
                    if row == 0:
                        roots.append({"object": name, "surface": list(hit),
                                      "root": list(point), "gap_m": (point - hit).length})
                    if row in (8, 16) and column in (2, 8, 14):
                        chest_samples.append({"object": name, "cloth": list(point),
                                              "chest": list(front),
                                              "support_station_v": support_v,
                                              "forward_clearance_m": front.y - point.y})
                grid.append(strip)
            created.append(_cloth(name, grid, white, collection, rig))

        tie_root, tie_hit = _neck(head_tree, 0.0, neck_y)
        tie_top = tie_root.z
        assert tie_top - TIE_HEIGHT >= bottom + 0.001
        # (drop fraction, half width). Rounded broad blade, gathered neck.
        profile = ((0, 0.0006), (0.07, 0.0017), (0.18, 0.0028),
                   (0.31, 0.0033), (0.52, 0.0067), (0.78, 0.0097),
                   (0.91, TIE_WIDTH / 2), (1, 0.0083))
        grid = []
        for row in range(41):
            v = row / 40.0
            a, b = next((a, b) for a, b in zip(profile, profile[1:]) if v <= b[0])
            t = _smooth((v - a[0]) / (b[0] - a[0]))
            radius = a[1] * (1 - t) + b[1] * t
            strip = []
            for column in range(17):
                u = -1 + 2 * column / 16.0
                x = radius * u
                z = tie_top - TIE_HEIGHT * v + 0.0010 * v**5 * abs(u)**2
                front = _front(body_tree, x, min(z, top - 0.001))
                settle = _smooth(v / 0.10)
                y = tie_root.y * (1 - settle) + (front.y - 0.00085) * settle
                # One unequal hanging fold gives cloth body without a thick prism.
                fold = math.exp(-0.5 * ((u + 0.25) / 0.35)**2)
                y -= (0.0006 * (1 - u * u) + 0.0007 * fold) * _smooth(v / 0.18)
                point = Vector((x, y, z))
                strip.append(point)
                if row in (12, 28, 36) and column == 8:
                    chest_samples.append({"object": NEW_NAMES[2], "cloth": list(point),
                                          "chest": list(front),
                                          "forward_clearance_m": front.y - point.y})
            grid.append(strip)
        roots.append({"object": NEW_NAMES[2], "surface": list(tie_hit),
                      "root": list(tie_root), "gap_m": (tie_root - tie_hit).length})
        created.append(_cloth(NEW_NAMES[2], grid, yellow, collection, rig))
        bpy.context.view_layer.update()
        assert controls == {name: _digest(scene.objects[name]) for name in controls}
        for name in TARGETS:
            scene.objects[name].hide_render = True
            scene.objects[name].hide_set(True)
        bpy.context.view_layer.update()
        assert controls == {name: _digest(scene.objects[name]) for name in controls}
        assert _sha(SOURCE) == SOURCE_HASH
        return {
            "status": "unreviewed collar/tie construction; no save or stage pass",
            "source_sha256": SOURCE_HASH, "helper_sha256": _sha(Path(__file__)),
            "hidden_objects": list(TARGETS), "new_objects": list(NEW_NAMES),
            "non_target_world_geometry_preserved": True,
            "control_count": len(controls), "controls": controls,
            "narrowed_bodice_bounds_m": body_bounds, "head_bounds_m": head_bounds,
            "new_bounds_m": {obj.name: _surface(obj)[1] for obj in created},
            "outline": {"collar_root_halfwidth_fraction": [0.52, 0.78],
                        "collar_lower_halfwidth_fraction": [0.10, 0.975],
                        "collar_lower_z_relative_bodice_top_m": [-0.0114, -0.0079],
                        "neck_y_m": neck_y, "tie_width_m": TIE_WIDTH,
                        "tie_height_m": TIE_HEIGHT, "tie_top_m": tie_top,
                        "cloth_thickness_m": 0.0006},
            "base_root_witnesses": roots, "base_chest_witnesses": chest_samples,
            "free_edge_front_supports": list(_FRONT_FALLBACKS),
            "rig": "single Body weights; current affine pose compensated",
            "limitations": [
                "Witnesses describe the base surface; verify evaluated final gaps.",
                "Short hidden roots tuck into the overlapping head/bodice neckline.",
                "Animated motion and reference-fidelity acceptance remain unverified.",
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

## render_body_023.py

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
output = ROOT / "body_023_review"
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
    ("candidate", "body_023_candidate.blend", list(contract["fixed_views"])),
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


