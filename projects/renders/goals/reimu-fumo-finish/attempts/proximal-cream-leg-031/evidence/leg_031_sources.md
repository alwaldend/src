# Leg031 rejected-candidate source evidence

Candidate ece85247dc07e9ac59388c20321b992c8638d4f3294ac4d4ef6436e975489b71
is rejected on pixels, despite geometric bounds and contact checks.
Source030b d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6
remains retained. Run writer then renderer in pinned5.2.1 LTS background
with factory startup, disabled autoexec, four threads, python-exit-code2.
Use fresh task-local output paths. Regenerated Blender bytes require fresh
hashes/review; do not replay this falsified strategy as new progress.

## leg_031_draft.py

SHA256 db707ea26be0b5018e0caaf181c4962dd3d9dcfb785351d409c19d2167f56921.

```python
"""Unexecuted Leg031 helper; root reviews before calling build_leg_031().

Only the two existing cream roots are replaced.  The evaluated source root
is a closed ellipsoid.  For each evaluated world point p, define

    q = clamp((p.y - ymin) / (ymax - ymin), 0, 1)
    u = clamp((q - 0.42) / (1 - 0.42), 0, 1)
    t = 6*u**5 - 15*u**4 + 10*u**3

The frontmost 42 percent of root depth has t=0 and is the exact protected
distal toe-intersection band.  The proximal warp is

    x' = cx + (x-cx) * (1 + 0.35*t) - side * 0.145*Wh*t
    y' = y + 0.055*Wh*t
    z' = cz + (z-cz) * (1 + 0.12*t) + 0.045*Wh*t

where side=-1 on the left and +1 on the right.  This is one monotone C2
field, not a ring extrusion: it retains the ellipsoid's changing section and
closed cap while widening and moving its proximal volume inward, up, and
seatward.  The uncertain physical-side exposure band is not forced.

This file defines the helper but does not call it, save a model, render, or
write a receipt.  A strict root/hem or root/skirt crossing is a falsifier;
technical clearance never constitutes visual acceptance.
"""

import bmesh
import hashlib
import json
import math
import struct
from pathlib import Path

import bpy
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree


ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "bow_030b_candidate.blend"
SOURCE_HASH = (
    "d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6"
)
TARGETS = (
    "Left short hidden leg root",
    "Right short hidden leg root",
)
CREATED_NAMES = (
    "Leg031 left proximal cream root",
    "Leg031 right proximal cream root",
)
PODS = (
    "Left black stuffed foot pod",
    "Right black stuffed foot pod",
)
HEM = "Hem026 curled cotton strip"
SKIRT = "Skirt022 joined gathered panels"
RIG = "ReimuFumoRig"
WH = 0.11743925511837006

PROTECTED_Q = 0.42
INWARD_SHIFT = 0.145 * WH
SEATWARD_SHIFT = 0.055 * WH
UPWARD_SHIFT = 0.045 * WH
LATERAL_WIDEN = 0.35
VERTICAL_WIDEN = 0.12
MAX_ROOT_HEM_CLEARANCE = 0.015 * WH
ROUNDTRIP_TOLERANCE = 5e-6
GEOMETRY_TYPES = {"MESH", "CURVE", "SURFACE", "FONT", "META"}


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _clamp(value, low=0.0, high=1.0):
    return max(low, min(high, value))


def _smootherstep(value):
    value = _clamp(value)
    return value**3 * (10.0 - 15.0 * value + 6.0 * value**2)


def _evaluated_geometry(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        points = [evaluated.matrix_world @ vertex.co
                  for vertex in mesh.vertices]
        faces = [tuple(face.vertices) for face in mesh.polygons]
        material_indices = [face.material_index for face in mesh.polygons]
        return points, faces, material_indices
    finally:
        evaluated.to_mesh_clear()


def _bounds(points):
    return [[min(point[axis] for point in points),
             max(point[axis] for point in points)]
            for axis in range(3)]


def _signed_volume(points, faces):
    volume = 0.0
    for face in faces:
        anchor = points[face[0]]
        for index in range(1, len(face) - 1):
            b = points[face[index]]
            c = points[face[index + 1]]
            volume += anchor.dot(b.cross(c)) / 6.0
    return volume


def _tree(obj):
    points, faces, _ = _evaluated_geometry(obj)
    return BVHTree.FromPolygons(points, faces, all_triangles=False)


def _overlap_count(a, b):
    return len(_tree(a).overlap(_tree(b)))


def _nearest_sample(a, b):
    """Bidirectional vertex-to-surface sample, not an exact mesh distance."""
    a_points, _, _ = _evaluated_geometry(a)
    b_points, _, _ = _evaluated_geometry(b)
    a_tree, b_tree = _tree(a), _tree(b)
    results = [a_tree.find_nearest(point) for point in b_points]
    results.extend(b_tree.find_nearest(point) for point in a_points)
    return min(result[3] for result in results if result is not None)


def _geometry_digest(obj):
    points, faces, material_indices = _evaluated_geometry(obj)
    digest = hashlib.sha256()
    digest.update(struct.pack("<2I", len(points), len(faces)))
    for point in points:
        assert all(math.isfinite(value) for value in point)
        digest.update(struct.pack("<3f", *point))
    for face, material_index in zip(faces, material_indices):
        digest.update(struct.pack("<2I", len(face), material_index))
        for index in face:
            digest.update(struct.pack("<I", index))
    return digest.hexdigest()


def _object_record(obj):
    record = {
        "type": obj.type,
        "parent": obj.parent.name if obj.parent else None,
        "matrix_world": [list(row) for row in obj.matrix_world],
        "visibility": [obj.hide_render, obj.hide_viewport, obj.hide_get()],
        "modifiers": [(modifier.name, modifier.type)
                      for modifier in obj.modifiers],
    }
    if obj.type in GEOMETRY_TYPES:
        record["evaluated_geometry_sha256"] = _geometry_digest(obj)
        record["materials"] = [material.name if material else None
                               for material in obj.data.materials]
    return record


def _pose_record(rig):
    return {bone.name: [list(row) for row in bone.matrix]
            for bone in rig.pose.bones}


def _scene_record(scene):
    view = scene.view_settings
    return {
        "scene": scene.name,
        "frame": scene.frame_current,
        "camera": scene.camera.name if scene.camera else None,
        "world": scene.world.name if scene.world else None,
        "view": [view.view_transform, view.look, view.exposure, view.gamma],
        "render_engine": scene.render.engine,
        "resolution": [scene.render.resolution_x,
                       scene.render.resolution_y,
                       scene.render.resolution_percentage],
    }


def _pose_inverse(rig, bone_name):
    bone = rig.pose.bones[bone_name]
    pose_delta = (
        rig.matrix_world
        @ bone.matrix
        @ bone.bone.matrix_local.inverted()
        @ rig.matrix_world.inverted()
    )
    return pose_delta.inverted()


def _warp(points, side):
    bounds = _bounds(points)
    ymin, ymax = bounds[1]
    depth = ymax - ymin
    assert depth > 1e-6
    cx = 0.5 * (bounds[0][0] + bounds[0][1])
    cz = 0.5 * (bounds[2][0] + bounds[2][1])
    output = []
    weights = []
    for point in points:
        q = _clamp((point.y - ymin) / depth)
        u = (q - PROTECTED_Q) / (1.0 - PROTECTED_Q)
        weight = _smootherstep(u)
        output.append(Vector((
            cx + (point.x - cx) * (1.0 + LATERAL_WIDEN * weight)
            - side * INWARD_SHIFT * weight,
            point.y + SEATWARD_SHIFT * weight,
            cz + (point.z - cz) * (1.0 + VERTICAL_WIDEN * weight)
            + UPWARD_SHIFT * weight,
        )))
        weights.append(weight)
    protected = [index for index, weight in enumerate(weights)
                 if weight == 0.0]
    assert protected
    protected_delta = max((output[index] - points[index]).length
                          for index in protected)
    assert protected_delta == 0.0
    assert all(math.isfinite(value) for point in output for value in point)
    return output, {
        "source_bounds_m": bounds,
        "source_depth_m": depth,
        "protected_q_max": PROTECTED_Q,
        "protected_vertex_count": len(protected),
        "protected_max_delta_m": protected_delta,
        "maximum_field_weight": max(weights),
    }


def _make_replacement(source, name, bone_name, desired_world, faces,
                      material_indices, rig, collection):
    materials = [slot.material for slot in source.material_slots]
    assert materials and all(material is not None for material in materials)
    assert all(index < len(materials) for index in material_indices)

    outward_faces = list(faces)
    desired_volume = _signed_volume(desired_world, outward_faces)
    winding_reversed = desired_volume < 0.0
    if winding_reversed:
        outward_faces = [tuple(reversed(face)) for face in outward_faces]

    inverse_pose = _pose_inverse(rig, bone_name)
    rest_points = [inverse_pose @ point for point in desired_world]
    mesh = bpy.data.meshes.new(name + " mesh")
    mesh.from_pydata(rest_points, [], outward_faces)
    for material in materials:
        mesh.materials.append(material)
    for polygon, material_index in zip(mesh.polygons, material_indices):
        polygon.material_index = material_index
        polygon.use_smooth = True
    mesh.update()

    check = bmesh.new()
    check.from_mesh(mesh)
    try:
        assert check.faces and all(edge.is_manifold for edge in check.edges)
        bmesh.ops.recalc_face_normals(check, faces=list(check.faces))
        check.to_mesh(mesh)
    finally:
        check.free()
    mesh.update()

    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    obj.parent = rig
    obj.matrix_world = Matrix.Identity(4)
    group = obj.vertex_groups.new(name=bone_name)
    group.add(list(range(len(mesh.vertices))), 1.0, "REPLACE")
    armature = obj.modifiers.new(
        "Inherited Leg031 bone attachment", "ARMATURE"
    )
    armature.object = rig
    obj["construction"] = "rounded proximal cream lobe from evaluated root"
    obj["status"] = "unreviewed Leg031 construction; no stage pass"
    bpy.context.view_layer.update()

    evaluated_points, evaluated_faces, _ = _evaluated_geometry(obj)
    assert len(evaluated_points) == len(desired_world)
    roundtrip_error = max((actual - expected).length
                          for actual, expected
                          in zip(evaluated_points, desired_world))
    assert roundtrip_error <= ROUNDTRIP_TOLERANCE, roundtrip_error
    volume = _signed_volume(evaluated_points, evaluated_faces)
    assert volume > 1e-9, volume
    assert [slot.material for slot in obj.material_slots] == materials
    return obj, {
        "source_signed_volume_m3": desired_volume,
        "replacement_signed_volume_m3": volume,
        "winding_reversed_by_guard": winding_reversed,
        "pose_roundtrip_max_error_m": roundtrip_error,
        "material_names": [material.name for material in materials],
        "material_ids": [material.as_pointer() for material in materials],
    }


def _reveal(root, pod, side):
    root_bounds = _bounds(_evaluated_geometry(root)[0])
    pod_bounds = _bounds(_evaluated_geometry(pod)[0])
    if side < 0:
        inward = root_bounds[0][1] - pod_bounds[0][1]
    else:
        inward = pod_bounds[0][0] - root_bounds[0][0]
    return {
        "root_bounds_m": root_bounds,
        "pod_bounds_m": pod_bounds,
        "inward_reveal_m": inward,
        "inward_reveal_Wh": inward / WH,
        "above_toe_m": root_bounds[2][1] - pod_bounds[2][1],
        "above_toe_Wh": (
            (root_bounds[2][1] - pod_bounds[2][1]) / WH
        ),
        "seatward_reveal_m": root_bounds[1][1] - pod_bounds[1][1],
        "seatward_reveal_Wh": (
            (root_bounds[1][1] - pod_bounds[1][1]) / WH
        ),
    }


def build_leg_031():
    """Build two in-memory replacements and return JSON-safe diagnostics."""
    assert bpy.app.version[:2] == (5, 2)
    assert Path(bpy.data.filepath).resolve() == SOURCE.resolve()
    assert _sha(SOURCE) == SOURCE_HASH
    assert bpy.context.scene.frame_current == 1
    assert bpy.context.mode == "OBJECT"
    assert not bpy.data.is_dirty, "Start from a clean frozen030b session"
    scene = bpy.context.scene
    required = TARGETS + CREATED_NAMES + PODS + (HEM, SKIRT, RIG)
    assert all(name in scene.objects for name in TARGETS + PODS +
               (HEM, SKIRT, RIG))
    assert all(name not in scene.objects for name in CREATED_NAMES)
    assert all(scene.objects[name].type == "MESH"
               for name in TARGETS + PODS + (HEM, SKIRT))
    assert all(not scene.objects[name].hide_render for name in TARGETS)
    assert all(not scene.objects[name].hide_get() for name in TARGETS)
    rig = scene.objects[RIG]
    assert rig.type == "ARMATURE"
    assert all(name in rig.pose.bones for name in ("Leg_L", "Leg_R"))

    expected_modifiers = ["ARMATURE", "LATTICE", "LATTICE"]
    expected_groups = ("Leg_L", "Leg_R")
    for source, bone_name in zip(
            (scene.objects[name] for name in TARGETS), expected_groups):
        assert [modifier.type for modifier in source.modifiers] == (
            expected_modifiers
        )
        assert [group.name for group in source.vertex_groups] == [bone_name]
        for vertex in source.data.vertices:
            assert len(vertex.groups) == 1
            assignment = vertex.groups[0]
            assert source.vertex_groups[assignment.group].name == bone_name
            assert abs(assignment.weight - 1.0) < 1e-6

    source_hashes = {name: _geometry_digest(scene.objects[name])
                     for name in TARGETS}
    control_names = [obj.name for obj in scene.objects
                     if obj.name not in TARGETS]
    controls = {name: _object_record(scene.objects[name])
                for name in control_names}
    pose = _pose_record(rig)
    scene_state = _scene_record(scene)
    material_ids = {material.name: material.as_pointer()
                    for material in bpy.data.materials}
    target_visibility = {
        name: (scene.objects[name].hide_render,
               scene.objects[name].hide_get())
        for name in TARGETS
    }
    collection = None
    created = []
    metrics = {}
    try:
        collection = bpy.data.collections.new("Leg031 proximal cream roots")
        scene.collection.children.link(collection)
        for side, target_name, created_name, pod_name, bone_name in zip(
                (-1, 1), TARGETS, CREATED_NAMES, PODS, expected_groups):
            source = scene.objects[target_name]
            source_points, faces, material_indices = _evaluated_geometry(
                source
            )
            desired, construction = _warp(source_points, side)
            replacement, replacement_metrics = _make_replacement(
                source, created_name, bone_name, desired, faces,
                material_indices, rig, collection,
            )
            created.append(replacement)
            baseline_root_pod = _overlap_count(source,
                                               scene.objects[pod_name])
            metrics[created_name] = {
                "side": "left" if side < 0 else "right",
                "source_name": target_name,
                "pod_name": pod_name,
                "construction": construction,
                "replacement": replacement_metrics,
                "source_root_pod_overlap_pairs": baseline_root_pod,
            }

        for name in TARGETS:
            scene.objects[name].hide_render = True
            scene.objects[name].hide_set(True)
        bpy.context.view_layer.update()

        failures = []
        for side, replacement, pod_name in zip((-1, 1), created, PODS):
            row = metrics[replacement.name]
            pod = scene.objects[pod_name]
            root_pod = _overlap_count(replacement, pod)
            root_hem = _overlap_count(replacement, scene.objects[HEM])
            root_skirt = _overlap_count(replacement, scene.objects[SKIRT])
            root_hem_clearance = _nearest_sample(replacement,
                                                 scene.objects[HEM])
            row["reveal"] = _reveal(replacement, pod, side)
            row["contact"] = {
                "root_pod_overlap_pairs": root_pod,
                "root_hem_overlap_pairs": root_hem,
                "root_skirt_overlap_pairs": root_skirt,
                "root_hem_nearest_sample_m": root_hem_clearance,
                "root_hem_nearest_sample_Wh": root_hem_clearance / WH,
            }
            if root_pod == 0:
                failures.append(replacement.name + ": detached from toe")
            if root_hem != 0:
                failures.append(replacement.name + ": crosses Hem026")
            if root_skirt != 0:
                failures.append(replacement.name + ": crosses Skirt022")
            if root_hem_clearance > MAX_ROOT_HEM_CLEARANCE:
                failures.append(replacement.name + ": floats below hem")

        assert controls == {name: _object_record(scene.objects[name])
                            for name in control_names}
        assert pose == _pose_record(rig)
        assert scene_state == _scene_record(scene)
        assert material_ids == {material.name: material.as_pointer()
                                for material in bpy.data.materials}
        assert source_hashes == {
            name: _geometry_digest(scene.objects[name]) for name in TARGETS
        }
        assert _sha(SOURCE) == SOURCE_HASH

        result = {
            "status": "Leg031 built in memory; no save or stage pass",
            "source": str(SOURCE),
            "source_sha256": SOURCE_HASH,
            "helper_sha256": _sha(Path(__file__).resolve()),
            "target_names": list(TARGETS),
            "created_names": [obj.name for obj in created],
            "hidden_names": list(TARGETS),
            "head_width_m": WH,
            "formula": {
                "protected_q": PROTECTED_Q,
                "inward_shift_Wh": INWARD_SHIFT / WH,
                "seatward_shift_Wh": SEATWARD_SHIFT / WH,
                "upward_shift_Wh": UPWARD_SHIFT / WH,
                "lateral_widen": LATERAL_WIDEN,
                "vertical_widen": VERTICAL_WIDEN,
                "field": (
                    "C2 smootherstep from protected depth to proximal cap"
                ),
            },
            "construction_metrics": metrics,
            "control_count": len(controls),
            "non_target_controls_unchanged": True,
            "rig_pose_unchanged": True,
            "scene_and_appearance_settings_unchanged": True,
            "existing_material_ids_reused": True,
            "source_roots_preserved_and_hidden": True,
            "spacing_residual_Wh": 0.6267,
            "limitations": [
                "No file was saved and no render was made.",
                "Zero BVH crossings do not establish a clean visible seam.",
                "Front, side, three-quarter, rear, and mirror review remains.",
                "Animation behavior beyond the frozen pose remains "
                "unverified.",
                "The deferred foot-spacing residual is unchanged.",
            ],
        }
        if failures:
            raise RuntimeError(json.dumps({
                "status": "Leg031 geometric falsifier",
                "failures": failures,
                "construction_metrics": metrics,
            }, sort_keys=True))
        json.dumps(result, sort_keys=True)
        return result
    except Exception:
        for name, (hide_render, hide_set) in target_visibility.items():
            scene.objects[name].hide_render = hide_render
            scene.objects[name].hide_set(hide_set)
        for obj in created:
            mesh = obj.data
            bpy.data.objects.remove(obj, do_unlink=True)
            if mesh.users == 0:
                bpy.data.meshes.remove(mesh)
        if collection is not None:
            bpy.data.collections.remove(collection)
        bpy.context.view_layer.update()
        raise
```

## leg_031_writer.py

SHA256 c490169024929fd1f781062ed01d11a742960e244eb2bb94ee6b8d94fcd3645e.

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "bow_030b_candidate.blend"
SOURCE_HASH = "d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6"
BASELINE = ROOT / "bow_030b_candidate.blend"
BASELINE_HASH = "d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6"
OUTPUT = ROOT / "leg_031_candidate.blend"
RECEIPT = ROOT / "leg_031_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "Left short hidden leg root",
    "Right short hidden leg root",
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
helper = ROOT / "leg_031_draft.py"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
result = scope["build_leg_031"]()
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
scene["candidate_status"] = "Leg031 unreviewed proximal cream study; no stage pass"
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
           "appearance_unchanged": True, "leg_construction": result,
           "limitations": ["No visual, animation, whole-scene technical or final acceptance."]}
assert sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {"controls", "rig_pose", "appearance", "leg_construction"}}))

```

## render_leg_031.py

SHA256 47bcf1e66301c446760da024a1e5732b35e504a969653ced941c2ff438c6beca.

```python
"""Frozen031 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "leg_031_candidate.blend"
OUT = ROOT / "leg_031_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "bow_030b_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "leg_031_writer_receipt.json").read_text())
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
    data = bpy.data.cameras.new("Frozen031_" + view)
    camera = bpy.data.objects.new("Frozen031_" + view, data)
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


