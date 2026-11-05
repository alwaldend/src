# Frozen construction and render sources

These are the exact task-local scripts for rejected candidate
`b770864c3015561a7e3c780590ceb13f85307086d996bd942000bc9668af384a`.
They preserve the experiment, not an accepted production asset.

To reconstruct geometry, extract the Python blocks to the indicated names
under ignored `out/reimu_fumo_finish/desktop_astra/` in a feature worktree.
Use a newly verified local writer under PROCESS.md, execute `build_head.py`
with its `__file__` bound to that path, then run `render_head.py` using
`bazel_agent bazel run //tools/blender:blender -- --background
--factory-startup --disable-autoexec --threads 4 --python <absolute-script>`.
The source donor is tracked and its hash is asserted. New session metadata
can change .blend bytes; a regeneration is a new candidate, not exact-byte
recovery or inherited acceptance. Rendering requires an absent output folder.

## build_head.py

```python
"""Author one reference-traced cushion and connected fringe in live Blender."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy
import numpy as np
from mathutils import Matrix

ROOT = Path(__file__).resolve().parent
WORKSPACE = ROOT.parents[2]
SOURCE = WORKSPACE / "projects/renders/assets/reimu_fumo/donors/a202/model.blend"
SOURCE_HASH = "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5"
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest() == SOURCE_HASH

# Append the donor scene into this disposable writer; never overwrite it.
old_scenes = list(bpy.data.scenes)
old_objects = set(bpy.data.objects)
with bpy.data.libraries.load(str(SOURCE), link=False) as (available, loaded):
    loaded.scenes = ["Attempt41_Manual_Head_Maquette"]
scene = loaded.scenes[0]
bpy.context.window.scene = scene
for old in old_scenes:
    bpy.data.scenes.remove(old)
for old in old_objects:
    if old.name in bpy.data.objects and old.name not in scene.objects:
        bpy.data.objects.remove(old, do_unlink=True)
scene.frame_set(1)
bpy.context.view_layer.update()
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(ROOT / "donor_appended.blend"))

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

collection = bpy.data.collections.new("Reference_Traced_Head_019")
scene.collection.children.link(collection)
hair = bpy.data.materials["Hair brown clay"]
skin = bpy.data.materials["Face fabric clay"]
WH = 0.1165
RX = 0.0577
RZ = 0.05665
CZ = 0.1639
CROWN = 0.2211


def front_y(x, z):
    radius = min(1.0, math.hypot(x / RX, (z - CZ) / RZ))
    return -0.053 + 0.025 * radius**4


def rear_y(x, z):
    radius = min(1.0, math.hypot(x / (RX * 0.976), (z - CZ) / (RZ * 0.976)))
    return 0.034 - 0.008 * radius**4


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


# Two rounded fabric panels joined by a broad, softly bulged side gusset.
# Ring rows explicitly own seam depth and front/back flattening.
segments = 96
vertices = [(0, -0.053, CZ)]
faces = []
face_mats = []
rings = []
for radius in (0.16, 0.32, 0.48, 0.64, 0.78, 0.88, 0.95, 1.0):
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
for y, radius in ((-0.020, 1.010), (-0.007, 1.014), (0.007, 1.008), (0.019, 0.997), (0.026, 0.976)):
    ring = []
    for j in range(segments):
        angle = 2 * math.pi * j / segments
        ring.append(len(vertices))
        vertices.append((RX * radius * math.cos(angle), y, CZ + RZ * radius * math.sin(angle)))
    for j in range(segments):
        nj = (j + 1) % segments
        faces.append((last[j], ring[j], ring[nj], last[nj]))
        face_mats.append(0)
    last = ring
for radius in (0.93, 0.84, 0.70, 0.52, 0.34, 0.16):
    ring = []
    for j in range(segments):
        angle = 2 * math.pi * j / segments
        x = RX * 0.976 * radius * math.cos(angle)
        z = CZ + RZ * 0.976 * radius * math.sin(angle)
        ring.append(len(vertices))
        vertices.append((x, rear_y(x, z), z))
    for j in range(segments):
        nj = (j + 1) % segments
        faces.append((last[j], ring[j], ring[nj], last[nj]))
        face_mats.append(0)
    last = ring
back_center = len(vertices)
vertices.append((0, 0.034, CZ))
for j in range(segments):
    faces.append((last[j], back_center, last[(j + 1) % segments]))
    face_mats.append(0)
head = make_mesh("Head_Gusseted_Cushion_019", vertices, faces, [hair, skin])
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
        y = front_y(x, z) - 0.00075 - 0.00020 * math.sin(math.pi * t)
        vertices.append((x, y, z))
for i in range(len(xs) - 1):
    for j in range(row_count):
        a = i * (row_count + 1) + j
        b = a + row_count + 1
        faces.append((a, b, b + 1, a + 1))
fringe = make_mesh("Hair_Continuous_Traced_Fringe_019", vertices, faces, [hair])
subdivision = fringe.modifiers.new("Rounded traced contour", "SUBSURF")
subdivision.levels = subdivision.render_levels = 1
thickness = fringe.modifiers.new("Thin felt edge", "SOLIDIFY")
thickness.thickness = 0.0007
thickness.offset = -1
bevel = fringe.modifiers.new("Soft sewn edge", "BEVEL")
bevel.width = 0.00018
bevel.segments = 2


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
            v.co.z = 0.096 + (v.co.z - 0.089) * 0.79
            if v.co.z > 0.117:
                v.co.y = rear_y(v.co.x, v.co.z) + 0.0008
            else:
                v.co.y = 0.030 + 0.003 * (v.co.z - 0.096) / 0.021
    elif "cheek lock" in name:
        for v in mesh.vertices:
            v.co.x *= 0.84
            v.co.z = 0.094 + (v.co.z - 0.0882) * 0.85
            t = max(0, min(1, (v.co.z - 0.128) / 0.020))
            v.co.y = (1 - t) * -0.036 + t * (front_y(v.co.x, v.co.z) - 0.0012)
    else:
        for v in mesh.vertices:
            v.co.x *= 0.80
            v.co.z += 0.0053 if "mouth" in name else 0.0043
            v.co.y = front_y(v.co.x, v.co.z) - 0.0007
    for v, residual in zip(mesh.vertices, depth_residual):
        v.co.y += float(residual) * 0.70
    mesh.update()

bpy.context.view_layer.update()
after_controls = {name: mesh_digest(bpy.data.objects[name]) for name in controls}
assert controls == after_controls, [name for name in controls if controls[name] != after_controls[name]]
scene["candidate_stage"] = "head module; unreviewed; rig not yet refitted"
scene["source_sha256"] = SOURCE_HASH
scene["authoring_method"] = "reference-traced connected control mesh"
bpy.ops.wm.save_as_mainfile(filepath=str(ROOT / "head_019_candidate.blend"))
receipt = {
    "source_sha256": SOURCE_HASH,
    "candidate_sha256": hashlib.sha256((ROOT / "head_019_candidate.blend").read_bytes()).hexdigest(),
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
(ROOT / "head_019_writer_receipt.json").write_text(json.dumps(receipt, indent=2))
result = {k: v for k, v in receipt.items() if k not in {"controls", "contour_witnesses"}}
```

## render_head.py

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
output = ROOT / "head_019_review"
output.mkdir(exist_ok=False)
receipt = {
    "blender_version": bpy.app.version_string,
    "build_hash": bpy.app.build_hash.decode(),
    "contract_sha256": hashlib.sha256(CONTRACT.read_bytes()).hexdigest(),
    "engine": "BLENDER_WORKBENCH",
    "purpose": "construction review; material and rig acceptance not claimed",
    "inputs": {},
    "renders": {},
}
for label, file_name, views in (
    ("baseline", "donor_appended.blend", ["front", "three_quarter"]),
    ("candidate", "head_019_candidate.blend", list(contract["fixed_views"])),
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
    for material in bpy.data.materials:
        if material.use_nodes:
            for node in material.node_tree.nodes:
                if node.type == "BSDF_PRINCIPLED":
                    material.diffuse_color = node.inputs["Base Color"].default_value
                    break
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
