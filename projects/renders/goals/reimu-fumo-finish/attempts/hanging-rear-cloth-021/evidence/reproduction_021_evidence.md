# Reproduction sources for the retained intermediate rear revision

Candidate SHA `9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22`.
This is a retained intermediate module, not an approved final asset. Extract
these source blocks into ignored `out/reimu_fumo_finish/desktop_astra/`.
The head builder uses the tracked, hash-bound A202 donor; it writes the
experimental 020b input. In a newly verified writer, execute each builder with
`__file__` bound to its path. A regenerated input may have different Blender
metadata bytes; bind its new SHA as a new candidate rather than claiming exact
recovery or silently bypassing the source assertion. The rear script consumes
only the explicitly verified input. The renderer uses pinned Blender through
`bazel_agent bazel run //tools/blender:blender -- --background
--factory-startup --disable-autoexec --threads 4 --python <absolute-script>`.

Rendering intentionally retains material viewport colors. Earlier 019/020
clay scripts replaced them with Principled socket defaults despite linked
color inputs; that was unsuitable for color/face-boundary comparisons.
Matching 020b baseline views were rerendered with this same color handling.
These are simple diagnostic colors, not final fabric material acceptance.

## build_head_020b.py

```python
"""Author one reference-traced cushion and connected fringe in live Blender."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy
import numpy as np
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree

ROOT = Path(__file__).resolve().parent
WORKSPACE = ROOT.parents[2]
SOURCE = WORKSPACE / "projects/renders/assets/reimu_fumo/donors/a202/model.blend"
SOURCE_HASH = "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5"
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest() == SOURCE_HASH

# Append the donor scene into this disposable writer; never overwrite it.
old_scenes = list(bpy.data.scenes)
old_objects = set(bpy.data.objects)
holding_scene = bpy.data.scenes.new("Empty_Import_Boundary_020b")
bpy.context.window.scene = holding_scene
for old in old_scenes:
    bpy.data.scenes.remove(old)
for old in old_objects:
    bpy.data.objects.remove(old, do_unlink=True)
with bpy.data.libraries.load(str(SOURCE), link=False) as (available, loaded):
    loaded.scenes = ["Attempt41_Manual_Head_Maquette"]
scene = loaded.scenes[0]
bpy.context.window.scene = scene
bpy.data.scenes.remove(holding_scene)
scene.frame_set(1)
bpy.context.view_layer.update()
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(ROOT / "donor_appended_020b.blend"))

replaced = {
    "Head_Cushion_Manual_Target",
    "A44 continuous hair cap with smooth opening",
    "A44 left temple fringe panel",
    "A44 left temple transition panel",
    "A44 off-center main bang panel",
    "A44 right swept fringe panel",
    "A44 right temple transition panel",
    "Hair_Front_Smooth_Panel",
    "Subtle crown center seam",
    "A42 Main lock left seated seam",
    "A42 Main lock right seated seam",
}
fitted = {
    o.name for o in scene.objects
    if not o.hide_render and (
        o.name.startswith("A45 ") or "rear lock" in o.name
        or o.name == "A44 tiny neutral embroidered mouth dash"
    )
}
targets = replaced | fitted


def mesh_digest(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    digest = hashlib.sha256()
    for vertex in mesh.vertices:
        p = evaluated.matrix_world @ vertex.co
        digest.update(struct.pack("<3f", *p))
    evaluated.to_mesh_clear()
    return digest.hexdigest()


controls = {
    o.name: mesh_digest(o) for o in scene.objects
    if o.type == "MESH" and not o.hide_render and o.name not in targets
}
for name in replaced:
    obj = bpy.data.objects.get(name)
    if obj:
        obj.hide_render = True
        obj.hide_set(True)

collection = bpy.data.collections.new("Reference_Traced_Head_020b")
scene.collection.children.link(collection)
hair = bpy.data.materials["Hair brown clay"]
skin = bpy.data.materials["Face fabric clay"]
WH = 0.1165
RX = WH * 0.5
RZ = 0.0574
CZ = 0.1580
CROWN = CZ + RZ


def front_y(x, z):
    radius = min(1.0, math.hypot(x / RX, (z - CZ) / RZ))
    return -0.008 - 0.045 * max(0.0, 1.0 - radius**2)**0.40


def rear_y(x, z):
    radius = min(1.0, math.hypot(x / RX, (z - CZ) / RZ))
    return -0.008 + 0.040 * max(0.0, 1.0 - radius**2)**0.58


def make_mesh(name, vertices, faces, materials):
    mesh = bpy.data.meshes.new(name)
    mesh.from_pydata(vertices, [], faces)
    mesh.update()
    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    for material in materials:
        mesh.materials.append(material)
    for polygon in mesh.polygons:
        polygon.use_smooth = True
    return obj


# Shaped front and rear cloth panels meet at a soft perimeter.
# The long constant-depth annular wall is removed, not scaled down.
segments = 96
vertices = [(0, -0.053, CZ)]
faces = []
face_mats = []
rings = []
for radius in (0.12, 0.24, 0.36, 0.48, 0.60, 0.70, 0.79, 0.86, 0.92, 0.96, 0.985, 0.998, 1.0):
    ring = []
    for j in range(segments):
        angle = 2 * math.pi * j / segments
        x, z = RX * radius * math.cos(angle), CZ + RZ * radius * math.sin(angle)
        ring.append(len(vertices))
        vertices.append((x, front_y(x, z), z))
    rings.append(ring)
for j in range(segments):
    faces.append((0, rings[0][j], rings[0][(j + 1) % segments]))
    face_mats.append(1)
for k in range(1, len(rings)):
    for j in range(segments):
        nj = (j + 1) % segments
        faces.append((rings[k - 1][j], rings[k][j], rings[k][nj], rings[k - 1][nj]))
        face_mats.append(1)
last = rings[-1]
for radius in (0.998, 0.985, 0.96, 0.92, 0.86, 0.79, 0.70, 0.60, 0.48, 0.36, 0.24, 0.12):
    ring = []
    for j in range(segments):
        angle = 2 * math.pi * j / segments
        x = RX * radius * math.cos(angle)
        z = CZ + RZ * radius * math.sin(angle)
        ring.append(len(vertices))
        vertices.append((x, rear_y(x, z), z))
    for j in range(segments):
        nj = (j + 1) % segments
        faces.append((last[j], ring[j], ring[nj], last[nj]))
        face_mats.append(0)
    last = ring
back_center = len(vertices)
vertices.append((0, 0.032, CZ))
for j in range(segments):
    faces.append((last[j], back_center, last[(j + 1) % segments]))
    face_mats.append(0)
head = make_mesh("Head_Gusseted_Cushion_020b", vertices, faces, [hair, skin])
for polygon, material_index in zip(head.data.polygons, face_mats):
    polygon.material_index = material_index
subdivision = head.modifiers.new("Soft panel transitions", "SUBSURF")
subdivision.levels = subdivision.render_levels = 1

# Manually read normalized front silhouette witnesses from the controlling
# photo. Each pair is x relative to head center, and downward offset from crown.
contour = [
    (-0.500, 0.662), (-0.472, 0.610), (-0.449, 0.586),
    (-0.421, 0.650), (-0.335, 0.749), (-0.315, 0.762),
    (-0.309, 0.610), (-0.266, 0.480), (-0.199, 0.376),
    (-0.158, 0.469), (-0.106, 0.564), (-0.049, 0.620),
    (0.024, 0.659), (0.088, 0.677), (0.099, 0.580),
    (0.122, 0.467), (0.167, 0.379), (0.222, 0.445),
    (0.278, 0.571), (0.308, 0.699), (0.321, 0.772),
    (0.386, 0.702), (0.433, 0.626), (0.448, 0.599),
    (0.473, 0.640), (0.500, 0.662),
]


def interpolate_linear(x):
    for (x0, y0), (x1, y1) in zip(contour, contour[1:]):
        if x0 <= x <= x1:
            t = (x - x0) / (x1 - x0)
            return y0 + t * (y1 - y0)
    return contour[-1][1]


# Fit to the evaluated cushion and preserve the brown cloth underlayer.
for polygon in head.data.polygons:
    if polygon.material_index == 1:
        x, _, z = polygon.center
        lower = CROWN - WH * interpolate_linear(max(-0.5, min(0.5, x / WH)))
        if z > lower + 0.001:
            polygon.material_index = 0
bpy.context.view_layer.update()
head_bvh = BVHTree.FromObject(head, bpy.context.evaluated_depsgraph_get())


def seated_front(x, z):
    hit, normal, face, distance = head_bvh.ray_cast(Vector((x, -0.3, z)), Vector((0, 1, 0)))
    return hit.y if hit is not None else front_y(x, z)


def seated_rear(x, z):
    anchor_z = max(z, 0.126)
    hit, normal, face, distance = head_bvh.ray_cast(Vector((x, 0.3, anchor_z)), Vector((0, -1, 0)))
    if hit is None:
        hit, normal, face, distance = head_bvh.find_nearest(Vector((x, 0.04, anchor_z)))
    depth = hit.y + 0.001 + 0.003 * max(0, min(1, (0.126 - z) / 0.038))
    return hit, depth


xs = sorted(set([round(-0.5 + i / 128, 8) for i in range(129)] + [p[0] for p in contour]))
row_count = 18
vertices, faces = [], []
for x_normalized in xs:
    x = WH * x_normalized
    lower_z = CROWN - WH * interpolate_linear(x_normalized)
    upper_z = CZ + (CROWN - CZ) * math.sqrt(max(0, 1 - (2 * x_normalized) ** 2))
    for j in range(row_count + 1):
        t = j / row_count
        z = lower_z + t * (upper_z - lower_z)
        y = seated_front(x, z) - 0.00075 - 0.00020 * math.sin(math.pi * t)
        vertices.append((x, y, z))
for i in range(len(xs) - 1):
    for j in range(row_count):
        a = i * (row_count + 1) + j
        b = a + row_count + 1
        faces.append((a, b, b + 1, a + 1))
fringe = make_mesh("Hair_Continuous_Traced_Fringe_020b", vertices, faces, [hair])
subdivision = fringe.modifiers.new("Rounded traced contour", "SUBSURF")
subdivision.levels = subdivision.render_levels = 1
thickness = fringe.modifiers.new("Thin felt edge", "SOLIDIFY")
thickness.thickness = 0.0007
thickness.offset = -1



def bake_world_mesh(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = bpy.data.meshes.new_from_object(evaluated, depsgraph=bpy.context.evaluated_depsgraph_get())
    matrix = evaluated.matrix_world.copy()
    for v in mesh.vertices:
        v.co = matrix @ v.co
    obj.parent = None
    obj.modifiers.clear()
    obj.data = mesh
    obj.matrix_world = Matrix.Identity(4)
    return mesh


for name in sorted(fitted):
    obj = bpy.data.objects[name]
    mesh = bake_world_mesh(obj)
    # Retain depth around each donor's fitted midsurface. Reprojecting every
    # vertex to one surface would collapse the baked felt/embroidery thickness.
    points = np.array([tuple(v.co) for v in mesh.vertices])
    x, z = points[:, 0] / WH, (points[:, 2] - CZ) / WH
    basis = np.column_stack((np.ones(len(x)), x, z, x*x, x*z, z*z))
    fitted_depth = basis @ np.linalg.lstsq(basis, points[:, 1], rcond=None)[0]
    depth_residual = points[:, 1] - fitted_depth
    if "rear lock" in name:
        for v in mesh.vertices:
            v.co.x *= 0.83
            v.co.z = 0.090 + (v.co.z - 0.089) * 0.79
            seated, depth = seated_rear(v.co.x, v.co.z)
            if v.co.z > 0.142 and abs(seated.z - v.co.z) > 0.001:
                v.co.x, v.co.z = seated.x, seated.z
            v.co.y = depth
    elif "cheek lock" in name:
        for v in mesh.vertices:
            v.co.x *= 0.84
            v.co.z = 0.088 + (v.co.z - 0.0882) * 0.85
            t = max(0, min(1, (v.co.z - 0.122) / 0.020))
            v.co.y = (1 - t) * -0.036 + t * (seated_front(v.co.x, v.co.z) - 0.0012)
    else:
        for v in mesh.vertices:
            v.co.x *= 0.80
            v.co.z += -0.0007 if "mouth" in name else -0.0017
            v.co.y = seated_front(v.co.x, v.co.z) - 0.0007
    for v, residual in zip(mesh.vertices, depth_residual):
        v.co.y += float(residual) * 0.70
    mesh.update()

bpy.context.view_layer.update()
after_controls = {name: mesh_digest(bpy.data.objects[name]) for name in controls}
assert controls == after_controls, [name for name in controls if controls[name] != after_controls[name]]
scene["candidate_stage"] = "head module; unreviewed; rig not yet refitted"
scene["source_sha256"] = SOURCE_HASH
scene["authoring_method"] = "reference-traced connected control mesh"
bpy.ops.wm.save_as_mainfile(filepath=str(ROOT / "head_020b_candidate.blend"))
receipt = {
    "source_sha256": SOURCE_HASH,
    "candidate_sha256": hashlib.sha256((ROOT / "head_020b_candidate.blend").read_bytes()).hexdigest(),
    "source_preserved": hashlib.sha256(SOURCE.read_bytes()).hexdigest() == SOURCE_HASH,
    "new_objects": [head.name, fringe.name],
    "replaced_objects": sorted(replaced),
    "fitted_objects": sorted(fitted),
    "unchanged_control_meshes": len(controls),
    "controls": controls,
    "head_width_m": WH,
    "fringe_tip_normalized": [0.588, 0.677],
    "contour_witnesses": contour,
}
(ROOT / "head_020b_writer_receipt.json").write_text(json.dumps(receipt, indent=2))
result = {k: v for k, v in receipt.items() if k not in {"controls", "contour_witnesses"}}
```

## build_rear_021.py

```python
"""Replace only the failed inherited rear locks with shaped cloth panels."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_020b_candidate.blend"
SOURCE_HASH = "ea9cf3ed9362b2ef3c90fc950151745984a1bb6d5606da430359dfd02b0ac54e"
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest() == SOURCE_HASH
assert bpy.data.filepath == str(SOURCE)
scene = bpy.context.scene
targets = {
    "A42 Left asymmetric rear lock",
    "A42 Off-center main rear lock",
    "A42 Short right rear lock",
}


def digest(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    h = hashlib.sha256()
    for vertex in mesh.vertices:
        h.update(struct.pack("<3f", *(evaluated.matrix_world @ vertex.co)))
    evaluated.to_mesh_clear()
    return h.hexdigest()


controls = {o.name: digest(o) for o in scene.objects
            if o.type == "MESH" and not o.hide_render and o.name not in targets}
for name in targets:
    bpy.data.objects[name].hide_render = True
    bpy.data.objects[name].hide_set(True)
collection = bpy.data.collections.new("Rear_Cloth_Panels_021")
scene.collection.children.link(collection)
hair = bpy.data.materials["Hair brown clay"]
profiles = {
    "Rear_Center_Cloth_021": [
        (0.000, 0.008, 0.207), (-0.001, 0.014, 0.202),
        (-0.003, 0.020, 0.190), (-0.004, 0.024, 0.173),
        (-0.003, 0.025, 0.154), (-0.002, 0.024, 0.135),
        (-0.004, 0.020, 0.117), (-0.004, 0.014, 0.105),
        (-0.002, 0.006, 0.096), (0.000, 0.001, 0.095),
    ],
    "Rear_Left_Cloth_021": [
        (-0.021, 0.005, 0.199), (-0.029, 0.009, 0.189),
        (-0.038, 0.011, 0.174), (-0.043, 0.012, 0.155),
        (-0.044, 0.011, 0.138), (-0.043, 0.009, 0.119),
        (-0.040, 0.006, 0.103), (-0.035, 0.001, 0.096),
    ],
    "Rear_Right_Cloth_021": [
        (0.022, 0.005, 0.198), (0.030, 0.009, 0.187),
        (0.039, 0.011, 0.173), (0.044, 0.011, 0.155),
        (0.044, 0.010, 0.137), (0.042, 0.008, 0.121),
        (0.038, 0.005, 0.109), (0.034, 0.001, 0.103),
    ],
}


def rear_depth(x, z):
    # Below the sewn attachment, continue the same depth without a jump.
    anchor_z = max(z, 0.127)
    radius = min(1.0, math.hypot(x / 0.05825, (anchor_z - 0.158) / 0.0574))
    surface = -0.008 + 0.040 * max(0.0, 1 - radius * radius)**0.58
    fall = max(0.0, min(1.0, (0.127 - z) / 0.035))
    return surface + 0.0015 + 0.004 * fall


created = []
for name, rows in profiles.items():
    vertices, faces = [], []
    columns = 9
    for center, width, height in rows:
        for j in range(columns):
            u = -1 + 2 * j / (columns - 1)
            x = center + width * u
            z = height - 0.001 * (1 - u * u)
            y = rear_depth(x, z) + 0.0006 * (1 - u * u)
            vertices.append((x, y, z))
    for i in range(len(rows) - 1):
        for j in range(columns - 1):
            a = i * columns + j
            faces.append((a, a + 1, a + columns + 1, a + columns))
    mesh = bpy.data.meshes.new(name)
    mesh.from_pydata(vertices, [], faces)
    mesh.materials.append(hair)
    mesh.update()
    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    for face in mesh.polygons:
        face.use_smooth = True
    sub = obj.modifiers.new("Soft cloth outline", "SUBSURF")
    sub.levels = sub.render_levels = 2
    solid = obj.modifiers.new("Felt thickness", "SOLIDIFY")
    solid.thickness = 0.0011
    solid.offset = 0.0
    created.append(obj.name)
bpy.context.view_layer.update()
assert controls == {name: digest(bpy.data.objects[name]) for name in controls}
scene["candidate_stage"] = "rear cloth correction; unreviewed; not an accepted asset"
bpy.context.preferences.filepaths.save_version = 0
path = ROOT / "head_021_candidate.blend"
bpy.ops.wm.save_as_mainfile(filepath=str(path))
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest() == SOURCE_HASH
receipt = {
    "source_sha256": SOURCE_HASH,
    "candidate_sha256": hashlib.sha256(path.read_bytes()).hexdigest(),
    "source_preserved": True,
    "new_objects": created,
    "hidden_objects": sorted(targets),
    "unchanged_control_meshes": len(controls),
    "controls": controls,
    "profiles": profiles,
}
(ROOT / "head_021_writer_receipt.json").write_text(json.dumps(receipt, indent=2))
result = {k: v for k, v in receipt.items() if k not in {"controls", "profiles"}}
```

## render_head_021.py

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
output = ROOT / "head_021_review"
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
    ("candidate", "head_021_candidate.blend", list(contract["fixed_views"])),
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
