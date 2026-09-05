# Foot033 reproducible construction sources

Saved candidate SHA256 98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48.
Source is immutable Head032, SHA256 6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8.
Pinned Blender 5.2.1 LTS build 9e2066aef7ef. Root reviewed the helper,
then executed writer and renderer in one background invocation. No overwrites
or post-render saves. Reproduction requires the exact source and032 receipts
named in the helper beside these scripts; no network or GUI is needed.

The writer verifies helper SHA256 4b4ef4d8ef9a692382e9770f2503be7a2268b5b101bd964d77a8f5aa8b194a34.
The numerical certification fixes and bounded endpoint/cleanup repairs did
not alter the chosen affine material field. Final saved-file evidence, not
the older pre-repair dry-run hash, owns candidate verification.

## foot_033_draft.py

```python
"""Unexecuted Foot033 conformal two-material pod helper.

The helper partitions each retained evaluated foot-pod surface with one fixed
mirrored affine field.  It introduces only interpolated seam vertices on
donor triangle edges; it does not move the pod surface, foot placement, hem,
rig, or any appearance control.  ``build_foot_033()`` builds in memory and
returns JSON-safe construction and first-hit mask metrics.  It never saves or
renders.
"""

from collections import Counter, defaultdict, deque
import hashlib
import json
import math
import struct
from pathlib import Path

import bpy
from bpy_extras.object_utils import world_to_camera_view
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree


ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_032_candidate.blend"
SOURCE_HASH = (
    "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8"
)
WRITER_RECEIPT = ROOT / "head_032_writer_receipt.json"
WRITER_RECEIPT_HASH = (
    "9e98892445a0520e360276fd8811ea9415f992218a65fb7584506cb880fef4e4"
)
RENDER_RECEIPT = ROOT / "head_032_eevee_review" / "render_receipt.json"
RENDER_RECEIPT_HASH = (
    "1d5c9b6fc5d432376e8c78d06fb3b4a4f1eee81b05916346b5992a5c5d858013"
)

TARGETS = (
    "Left black stuffed foot pod",
    "Right black stuffed foot pod",
)
CREATED_NAMES = (
    "Foot033 left conformal two-material stuffed pod",
    "Foot033 right conformal two-material stuffed pod",
)
MESH_NAMES = tuple(name + " mesh" for name in CREATED_NAMES)
COLLECTION_NAME = "Foot033 conformal two-material pods"
BONES = ("Leg_L", "Leg_R")
SIDES = (-1, 1)
HEM = "Hem026 curled cotton strip"
RIG = "ReimuFumoRig"
BLACK_MATERIAL = "Feet black velour.002"
CREAM_MATERIAL = "Dress warm white cloth.002"
WH = 0.11743925511837006

EXPECTED_GEOMETRY = {
    TARGETS[0]: (
        "873d3a562a179b1aa86e869842d349421cafc0baf506d672306847aaef0fb2ae"
    ),
    TARGETS[1]: (
        "4ccea62076fbe96aa9425457acb28d4aff261ada29291c588bb564810f1a5693"
    ),
    HEM: "895787943ee693e8ab636890639e76a85d96318c8a4af3b64434eee96da6b020",
}

FIELD_U = 0.60
FIELD_P = 0.60
FIELD_OFFSET = 0.64
FIELD_EPSILON = 2e-9
MATERIAL_SIDE_TOLERANCE = 1e-6
ROUNDTRIP_TOLERANCE_M = 5e-6
BOUNDS_TOLERANCE_M = 2e-8
MIN_TWICE_TRIANGLE_AREA_M2 = 1e-14
CAMERA_DISTANCE_M = 5.0
GEOMETRY_TYPES = {"MESH", "CURVE", "SURFACE", "FONT", "META"}

ARMATURE_SCALAR_SETTINGS = (
    "show_expanded",
    "show_in_editmode",
    "show_on_cage",
    "show_render",
    "show_viewport",
    "use_apply_on_spline",
    "use_bone_envelopes",
    "use_deform_preserve_volume",
    "use_multi_modifier",
    "use_pin_to_last",
    "use_vertex_groups",
    "vertex_group",
    "invert_vertex_group",
)


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _evaluated_geometry(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        points = [evaluated.matrix_world @ vertex.co
                  for vertex in mesh.vertices]
        faces = [tuple(polygon.vertices) for polygon in mesh.polygons]
        material_indices = [polygon.material_index
                            for polygon in mesh.polygons]
        return points, faces, material_indices
    finally:
        evaluated.to_mesh_clear()


def _geometry_digest(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    digest = hashlib.sha256()
    try:
        # This deliberately matches head_032_writer.py's receipt record,
        # including the evaluated edge count in the digest header.
        digest.update(struct.pack(
            "<3I", len(mesh.vertices), len(mesh.edges), len(mesh.polygons)
        ))
        for vertex in mesh.vertices:
            point = evaluated.matrix_world @ vertex.co
            assert all(math.isfinite(value) for value in point)
            digest.update(struct.pack("<3f", *point))
        for face in mesh.polygons:
            digest.update(struct.pack(
                "<2I", len(face.vertices), face.material_index
            ))
            for index in face.vertices:
                digest.update(struct.pack("<I", index))
        return digest.hexdigest()
    finally:
        evaluated.to_mesh_clear()


def _receipt_control_record(obj):
    return {
        "geometry": _geometry_digest(obj),
        "materials": [slot.material.name if slot.material else None
                      for slot in obj.material_slots],
        "visibility": [obj.hide_render, obj.hide_viewport, obj.hide_get()],
        "parent": obj.parent.name if obj.parent else None,
        "modifiers": [[modifier.name, modifier.type]
                      for modifier in obj.modifiers],
    }


def _object_record(obj):
    record = {
        "type": obj.type,
        "parent": obj.parent.name if obj.parent else None,
        "matrix_world": [[float(value) for value in row]
                         for row in obj.matrix_world],
        "visibility": [obj.hide_render, obj.hide_viewport, obj.hide_get()],
        "modifiers": [(modifier.name, modifier.type)
                      for modifier in obj.modifiers],
    }
    if obj.type in GEOMETRY_TYPES:
        record["evaluated_geometry_sha256"] = _geometry_digest(obj)
        record["materials"] = [
            slot.material.name if slot.material else None
            for slot in obj.material_slots
        ]
    return record


def _pose_record(rig):
    return {
        bone.name: [[float(value) for value in row] for row in bone.matrix]
        for bone in rig.pose.bones
    }


def _scene_record(scene):
    view = scene.view_settings
    return {
        "frame": scene.frame_current,
        "camera": scene.camera.name if scene.camera else None,
        "world": scene.world.name if scene.world else None,
        "engine": scene.render.engine,
        "resolution": [scene.render.resolution_x,
                       scene.render.resolution_y,
                       scene.render.resolution_percentage],
        "view": [view.view_transform, view.look, view.exposure, view.gamma],
    }


def _armature_modifier_record(modifier):
    assert modifier.type == "ARMATURE"
    settings = {}
    for name in ARMATURE_SCALAR_SETTINGS:
        if hasattr(modifier, name):
            value = getattr(modifier, name)
            assert isinstance(value, (bool, int, float, str))
            settings[name] = value
    return {
        "name": modifier.name,
        "object": modifier.object.name if modifier.object else None,
        "settings": settings,
    }


def _copy_armature_modifier(source_modifier, target_obj):
    source_record = _armature_modifier_record(source_modifier)
    modifier = target_obj.modifiers.new(source_modifier.name, "ARMATURE")
    assert modifier.name == source_modifier.name
    modifier.object = source_modifier.object
    for name, value in source_record["settings"].items():
        setattr(modifier, name, value)
    assert _armature_modifier_record(modifier) == source_record
    return modifier, source_record


def _uv_state(mesh):
    layers = mesh.uv_layers
    return {
        "names": [layer.name for layer in layers],
        "active_index": layers.active_index,
        "active_name": layers.active.name if layers.active else None,
        "active_render": [layer.name for layer in layers
                          if getattr(layer, "active_render", False)],
    }


def _source_triangles(obj):
    """Return evaluated world triangles and per-loop UV values."""
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        mesh.calc_loop_triangles()
        points = [evaluated.matrix_world @ vertex.co
                  for vertex in mesh.vertices]
        uv_state = _uv_state(mesh)
        assert len(uv_state["names"]) == 1
        assert uv_state["active_index"] == 0
        assert uv_state["active_name"] == uv_state["names"][0]
        assert uv_state["active_render"] == uv_state["names"]
        uv_names = uv_state["names"]
        triangles = []
        for triangle in mesh.loop_triangles:
            corner_uv = {}
            for layer in mesh.uv_layers:
                corner_uv[layer.name] = [
                    tuple(float(value) for value in layer.data[loop].uv)
                    for loop in triangle.loops
                ]
            triangles.append({
                "vertices": tuple(triangle.vertices),
                "uv": corner_uv,
            })
        return {
            "points": points,
            "triangles": triangles,
            "uv_names": uv_names,
            "uv_state": uv_state,
            "source_polygon_count": len(mesh.polygons),
        }
    finally:
        evaluated.to_mesh_clear()


def _bounds(points):
    return [
        [min(point[axis] for point in points),
         max(point[axis] for point in points)]
        for axis in range(3)
    ]


def _signed_volume(points, faces):
    volume = 0.0
    for face in faces:
        anchor = points[face[0]]
        for index in range(1, len(face) - 1):
            b = points[face[index]]
            c = points[face[index + 1]]
            volume += anchor.dot(b.cross(c)) / 6.0
    return volume


def _field(point, bounds, side):
    center = [0.5 * (axis[0] + axis[1]) for axis in bounds]
    radius = [0.5 * (axis[1] - axis[0]) for axis in bounds]
    assert min(radius) > 1e-7
    inward = -side * (point.x - center[0]) / radius[0]
    proximal = (point.y - center[1]) / radius[1]
    zeta = (point.z - center[2]) / radius[2]
    value = zeta + FIELD_U * inward + FIELD_P * proximal - FIELD_OFFSET
    return value, (inward, proximal, zeta)


def _topology(points, faces):
    edge_faces = defaultdict(list)
    adjacency = defaultdict(set)
    used = set()
    for face_index, face in enumerate(faces):
        used.update(face)
        for a, b in zip(face, face[1:] + face[:1]):
            edge = tuple(sorted((a, b)))
            edge_faces[edge].append(face_index)
            adjacency[a].add(b)
            adjacency[b].add(a)
    components = 0
    unseen = set(used)
    while unseen:
        components += 1
        queue = [unseen.pop()]
        while queue:
            vertex = queue.pop()
            for neighbor in adjacency[vertex]:
                if neighbor in unseen:
                    unseen.remove(neighbor)
                    queue.append(neighbor)
    chi = len(used) - len(edge_faces) + len(faces)
    genus = (2 * components - chi) / 2.0
    return {
        "vertices": len(used),
        "edges": len(edge_faces),
        "faces": len(faces),
        "components": components,
        "euler": chi,
        "genus": genus,
        "non_two_incident_edges": sum(len(rows) != 2
                                      for rows in edge_faces.values()),
    }, edge_faces


def _clip_surface(snapshot, side):
    points = list(snapshot["points"])
    bounds = _bounds(points)
    raw_fields = [_field(point, bounds, side)[0] for point in points]
    fields = [0.0 if abs(value) <= FIELD_EPSILON else value
              for value in raw_fields]
    source_faces = [row["vertices"] for row in snapshot["triangles"]]
    source_volume = _signed_volume(points, source_faces)
    reversed_by_guard = source_volume < 0.0

    triangles = []
    for row in snapshot["triangles"]:
        if reversed_by_guard:
            triangles.append({
                "vertices": (row["vertices"][0], row["vertices"][2],
                             row["vertices"][1]),
                "uv": {
                    name: (values[0], values[2], values[1])
                    for name, values in row["uv"].items()
                },
            })
        else:
            triangles.append(row)
    source_faces = [row["vertices"] for row in triangles]
    source_volume = _signed_volume(points, source_faces)
    assert source_volume > 1e-12

    output_points = list(points)
    output_fields = list(fields)
    crossing_cache = {}
    output_faces = []
    output_materials = []
    output_uv = {name: [] for name in snapshot["uv_names"]}

    def intersection(a, b):
        if abs(a["field"]) <= FIELD_EPSILON:
            return a
        if abs(b["field"]) <= FIELD_EPSILON:
            return b
        key = tuple(sorted((a["source"], b["source"])))
        denominator = a["field"] - b["field"]
        assert abs(denominator) > FIELD_EPSILON
        t = a["field"] / denominator
        assert -FIELD_EPSILON <= t <= 1.0 + FIELD_EPSILON
        t = max(0.0, min(1.0, t))
        if key not in crossing_cache:
            point = points[a["source"]].lerp(points[b["source"]], t)
            crossing_cache[key] = len(output_points)
            output_points.append(point)
            output_fields.append(_field(point, bounds, side)[0])
        weights = tuple(
            (1.0 - t) * wa + t * wb
            for wa, wb in zip(a["weights"], b["weights"])
        )
        return {
            "index": crossing_cache[key],
            "source": key[0],
            "field": 0.0,
            "weights": weights,
        }

    def clipped(base, positive):
        polygon = []

        def append_unique(item):
            if not polygon or polygon[-1]["index"] != item["index"]:
                polygon.append(item)

        for a, b in zip(base, base[1:] + base[:1]):
            inside_a = (a["field"] >= -FIELD_EPSILON if positive
                        else a["field"] <= FIELD_EPSILON)
            inside_b = (b["field"] >= -FIELD_EPSILON if positive
                        else b["field"] <= FIELD_EPSILON)
            if inside_a:
                append_unique(a)
            if inside_a != inside_b:
                append_unique(intersection(a, b))
        if len(polygon) > 1 and (
                polygon[0]["index"] == polygon[-1]["index"]):
            polygon.pop()
        assert len({item["index"] for item in polygon}) == len(polygon)
        return polygon

    for triangle in triangles:
        vertex_ids = triangle["vertices"]
        base = []
        for local, vertex_id in enumerate(vertex_ids):
            weights = tuple(1.0 if index == local else 0.0
                            for index in range(3))
            base.append({
                "index": vertex_id,
                "source": vertex_id,
                "field": fields[vertex_id],
                "weights": weights,
            })
        for positive, material_index in ((False, 0), (True, 1)):
            polygon = clipped(base, positive)
            if len(polygon) < 3:
                continue
            for index in range(1, len(polygon) - 1):
                corners = (polygon[0], polygon[index], polygon[index + 1])
                face = tuple(corner["index"] for corner in corners)
                twice_area = (
                    (output_points[face[1]] - output_points[face[0]])
                    .cross(output_points[face[2]] - output_points[face[0]])
                    .length
                )
                assert twice_area > MIN_TWICE_TRIANGLE_AREA_M2, twice_area
                output_faces.append(face)
                output_materials.append(material_index)
                for name in snapshot["uv_names"]:
                    source_uv = triangle["uv"][name]
                    output_uv[name].append([
                        tuple(sum(weight * source_uv[corner_index][axis]
                                  for corner_index, weight
                                  in enumerate(corner["weights"]))
                              for axis in range(2))
                        for corner in corners
                    ])

    topology, edge_faces = _topology(output_points, output_faces)
    donor_topology, _ = _topology(points, source_faces)
    seam_edges = []
    seam_adjacency = defaultdict(set)
    for edge, adjacent in edge_faces.items():
        if len(adjacent) == 2 and (
                output_materials[adjacent[0]] !=
                output_materials[adjacent[1]]):
            seam_edges.append(edge)
            seam_adjacency[edge[0]].add(edge[1])
            seam_adjacency[edge[1]].add(edge[0])
    seam_components = 0
    unseen = set(seam_adjacency)
    while unseen:
        seam_components += 1
        queue = [unseen.pop()]
        while queue:
            vertex = queue.pop()
            for neighbor in seam_adjacency[vertex]:
                if neighbor in unseen:
                    unseen.remove(neighbor)
                    queue.append(neighbor)

    max_side_violation = 0.0
    for face, material_index in zip(output_faces, output_materials):
        values = [output_fields[index] for index in face]
        violation = (-min(values) if material_index == 1 else max(values))
        max_side_violation = max(max_side_violation, violation)
    output_bounds = _bounds(output_points)
    max_bounds_delta = max(
        abs(output_bounds[axis][bound] - bounds[axis][bound])
        for axis in range(3) for bound in range(2)
    )
    output_volume = _signed_volume(output_points, output_faces)

    assert topology["components"] == 1
    assert topology["non_two_incident_edges"] == 0
    assert topology["euler"] == donor_topology["euler"] == 2
    assert abs(topology["genus"]) <= 1e-9
    assert seam_components == 1
    assert seam_adjacency and all(len(neighbors) == 2
                                  for neighbors in seam_adjacency.values())
    # mathutils stores the interpolated seam coordinate as float32.  Preserve
    # and report its actual residual, with a tolerance far below one pixel or
    # any physical construction scale.
    assert max_side_violation <= MATERIAL_SIDE_TOLERANCE, max_side_violation
    assert max_bounds_delta <= BOUNDS_TOLERANCE_M, max_bounds_delta
    assert abs(output_volume - source_volume) <= 1e-10

    bottom_values = []
    for point, material_field in zip(output_points, output_fields):
        _, (_, _, zeta) = _field(point, bounds, side)
        if zeta <= -0.60 + FIELD_EPSILON:
            bottom_values.append(material_field)
    assert bottom_values and max(bottom_values) < 0.0

    center = [0.5 * (axis[0] + axis[1]) for axis in bounds]
    radius = [0.5 * (axis[1] - axis[0]) for axis in bounds]
    distal_top = Vector((center[0], center[1] - radius[1],
                         center[2] + radius[2]))
    distal_field = _field(distal_top, bounds, side)[0]
    assert distal_field < 0.0
    assert abs(distal_field + 0.24) <= MATERIAL_SIDE_TOLERANCE

    return {
        "points": output_points,
        "faces": output_faces,
        "materials": output_materials,
        "uv": output_uv,
        "uv_state": snapshot["uv_state"],
    }, {
        "side": "left" if side < 0 else "right",
        "source_bounds_m": bounds,
        "source_half_extents_m": radius,
        "source_polygon_count": snapshot["source_polygon_count"],
        "source_triangle_count": len(source_faces),
        "output_vertex_count": len(output_points),
        "output_triangle_count": len(output_faces),
        "crossing_vertex_count": len(crossing_cache),
        "source_signed_volume_m3": source_volume,
        "output_signed_volume_m3": output_volume,
        "winding_reversed_by_guard": reversed_by_guard,
        "donor_topology": donor_topology,
        "output_topology": topology,
        "seam_edge_count": len(seam_edges),
        "seam_vertex_count": len(seam_adjacency),
        "seam_components": seam_components,
        "non_degree_two_seam_vertices": sum(
            len(neighbors) != 2 for neighbors in seam_adjacency.values()
        ),
        "maximum_material_side_violation": max_side_violation,
        "maximum_bounds_delta_m": max_bounds_delta,
        "maximum_bottom_20_percent_field": max(bottom_values),
        "distal_center_top_field": distal_field,
        "uv_layers_preserved": list(snapshot["uv_names"]),
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


def _make_replacement(source, created_name, mesh_name, bone_name, partition,
                      rig, collection):
    black = source.material_slots[0].material
    cream = bpy.data.materials[CREAM_MATERIAL]
    assert black is bpy.data.materials[BLACK_MATERIAL]
    inverse_pose = _pose_inverse(rig, bone_name)
    rest_points = [inverse_pose @ point for point in partition["points"]]
    source_modifier = source.modifiers[0]
    mesh = None
    obj = None
    try:
        assert created_name not in bpy.data.objects
        assert mesh_name not in bpy.data.meshes
        mesh = bpy.data.meshes.new(mesh_name)
        assert mesh.name == mesh_name
        mesh.from_pydata(rest_points, [], partition["faces"])
        mesh.materials.append(black)
        mesh.materials.append(cream)
        for polygon, material_index in zip(mesh.polygons,
                                           partition["materials"]):
            polygon.material_index = material_index
            polygon.use_smooth = True
        for name, faces_uv in partition["uv"].items():
            layer = mesh.uv_layers.new(name=name)
            assert layer.name == name
            for polygon, corner_uv in zip(mesh.polygons, faces_uv):
                for loop_index, uv in zip(polygon.loop_indices, corner_uv):
                    layer.data[loop_index].uv = uv
        assert len(mesh.uv_layers) == 1
        mesh.uv_layers.active_index = partition["uv_state"]["active_index"]
        for layer in mesh.uv_layers:
            layer.active_render = (
                layer.name in partition["uv_state"]["active_render"]
            )
        assert _uv_state(mesh) == partition["uv_state"]
        mesh.update()

        obj = bpy.data.objects.new(created_name, mesh)
        assert obj.name == created_name
        collection.objects.link(obj)
        obj.parent = rig
        obj.matrix_world = Matrix.Identity(4)
        group = obj.vertex_groups.new(name=bone_name)
        assert group.name == bone_name
        group.add(list(range(len(mesh.vertices))), 1.0, "REPLACE")
        _, modifier_record = _copy_armature_modifier(source_modifier, obj)
        obj["construction"] = "retained pod surface split by Foot033 field"
        obj["status"] = "unreviewed Foot033 in-memory construction"
        bpy.context.view_layer.update()

        evaluated_points, evaluated_faces, evaluated_materials = (
            _evaluated_geometry(obj)
        )
        assert len(evaluated_points) == len(partition["points"])
        roundtrip = max(
            (actual - expected).length
            for actual, expected in zip(
                evaluated_points, partition["points"]
            )
        )
        assert roundtrip <= ROUNDTRIP_TOLERANCE_M, roundtrip
        assert evaluated_faces == partition["faces"]
        assert evaluated_materials == partition["materials"]
        assert _signed_volume(evaluated_points, evaluated_faces) > 1e-12
        assert obj.data.name == mesh_name
        assert _uv_state(obj.data) == partition["uv_state"]
        assert _armature_modifier_record(obj.modifiers[0]) == modifier_record
        return obj, {
            "bone": bone_name,
            "object_datablock_name": obj.name,
            "mesh_datablock_name": mesh.name,
            "material_names": [black.name, cream.name],
            "material_ids": [black.as_pointer(), cream.as_pointer()],
            "source_armature_modifier": modifier_record,
            "uv_state": _uv_state(mesh),
            "pose_roundtrip_max_error_m": roundtrip,
            "evaluated_bounds_m": _bounds(evaluated_points),
        }
    except Exception:
        if obj is not None and bpy.data.objects.get(obj.name) is obj:
            bpy.data.objects.remove(obj, do_unlink=True)
        if mesh is not None and bpy.data.meshes.get(mesh.name) is mesh:
            assert mesh.users == 0
            bpy.data.meshes.remove(mesh)
        bpy.context.view_layer.update()
        raise


def _pixel_components(pixels):
    unseen = set(pixels)
    sizes = []
    while unseen:
        start = unseen.pop()
        queue = deque([start])
        size = 0
        while queue:
            x, y = queue.popleft()
            size += 1
            for neighbor in ((x - 1, y), (x + 1, y),
                             (x, y - 1), (x, y + 1)):
                if neighbor in unseen:
                    unseen.remove(neighbor)
                    queue.append(neighbor)
        sizes.append(size)
    return sorted(sizes, reverse=True)


def _pixel_bbox(pixels):
    if not pixels:
        return None
    return [min(x for x, _ in pixels), min(y for _, y in pixels),
            max(x for x, _ in pixels), max(y for _, y in pixels)]


def _ray(camera, pixel_x, pixel_y, width, height, scene):
    aspect = (
        width * scene.render.pixel_aspect_x /
        (height * scene.render.pixel_aspect_y)
    )
    scale = camera.data.ortho_scale
    local_x = ((pixel_x + 0.5) / width - 0.5) * scale * aspect
    local_y = (0.5 - (pixel_y + 0.5) / height) * scale
    origin = camera.matrix_world @ Vector((local_x, local_y, 0.0))
    direction = (
        camera.matrix_world.to_quaternion() @ Vector((0.0, 0.0, -1.0))
    ).normalized()
    return origin, direction


def _projected_bbox(scene, camera, points, width, height):
    projected = [world_to_camera_view(scene, camera, point)
                 for point in points]
    xs = [value.x * width for value in projected]
    ys = [(1.0 - value.y) * height for value in projected]
    return [
        max(0, math.floor(min(xs)) - 2),
        max(0, math.floor(min(ys)) - 2),
        min(width - 1, math.ceil(max(xs)) + 2),
        min(height - 1, math.ceil(max(ys)) + 2),
    ]


def _first_hit_masks(scene, created, render_receipt):
    depsgraph = bpy.context.evaluated_depsgraph_get()
    width, height, percentage = render_receipt["settings"]["resolution"]
    assert percentage == 100 and (width, height) == (512, 512)
    output = {}
    created_set = {obj.name for obj in created}

    for view_name, view in render_receipt["renders"].items():
        camera_name = view["camera"]["camera"]
        camera = scene.objects[camera_name]
        expected_location = Vector(view["camera"]["location_m"])
        expected_rotation = Vector(view["camera"]["rotation_euler_rad"])
        assert (camera.location - expected_location).length <= 2e-6
        assert (Vector(camera.rotation_euler) - expected_rotation).length <= 2e-6
        assert camera.data.type == "ORTHO"
        assert abs(camera.data.ortho_scale - 0.292) <= 2e-6
        assert camera.parent is None
        assert (camera.scale - Vector((1.0, 1.0, 1.0))).length <= 1e-7
        assert abs(camera.data.shift_x) <= 1e-9
        assert abs(camera.data.shift_y) <= 1e-9
        assert abs(scene.render.pixel_aspect_x - 1.0) <= 1e-9
        assert abs(scene.render.pixel_aspect_y - 1.0) <= 1e-9
        view_rows = {}

        for obj in created:
            points, faces, _ = _evaluated_geometry(obj)
            tree = BVHTree.FromPolygons(points, faces, all_triangles=True)
            search = _projected_bbox(scene, camera, points, width, height)
            footprint = set()
            cream = set()
            black = set()
            labels = Counter()
            for pixel_y in range(search[1], search[3] + 1):
                for pixel_x in range(search[0], search[2] + 1):
                    origin, direction = _ray(
                        camera, pixel_x, pixel_y, width, height, scene
                    )
                    isolated = tree.ray_cast(
                        origin, direction, CAMERA_DISTANCE_M
                    )
                    if isolated[0] is None:
                        continue
                    pixel = (pixel_x, pixel_y)
                    footprint.add(pixel)
                    hit, _, _, polygon_index, hit_obj, _ = scene.ray_cast(
                        depsgraph, origin, direction,
                        distance=CAMERA_DISTANCE_M,
                    )
                    if not hit:
                        labels["miss"] += 1
                        continue
                    hit_name = hit_obj.name
                    if hit_name == obj.name:
                        assert 0 <= polygon_index < len(obj.data.polygons)
                        material_index = obj.data.polygons[
                            polygon_index
                        ].material_index
                        if material_index == 1:
                            labels["cream"] += 1
                            cream.add(pixel)
                        else:
                            assert material_index == 0
                            labels["black"] += 1
                            black.add(pixel)
                    elif hit_name in created_set:
                        labels["other_foot"] += 1
                    else:
                        labels["occluder:" + hit_name] += 1

            components = _pixel_components(cream)
            cream_bbox = _pixel_bbox(cream)
            black_bbox = _pixel_bbox(black)
            footprint_bbox = _pixel_bbox(footprint)
            pixels_per_wh = width * WH / camera.data.ortho_scale
            cream_width = (
                (cream_bbox[2] - cream_bbox[0] + 1) / pixels_per_wh
                if cream_bbox else 0.0
            )
            black_width = (
                black_bbox[2] - black_bbox[0] + 1
                if black_bbox else 0
            )
            footprint_width = (
                footprint_bbox[2] - footprint_bbox[0] + 1
                if footprint_bbox else 0
            )
            view_rows[obj.name] = {
                "projected_search_bbox_px": search,
                "isolated_footprint_pixels": len(footprint),
                "isolated_footprint_bbox_px": footprint_bbox,
                "first_hit_counts": dict(sorted(labels.items())),
                "cream_bbox_px": cream_bbox,
                "black_bbox_px": black_bbox,
                "cream_screen_horizontal_span_Wh": cream_width,
                "black_screen_horizontal_span_fraction_of_footprint": (
                    black_width / footprint_width if footprint_width else None
                ),
                "cream_connected_components": len(components),
                "cream_component_sizes_px": components,
                "cream_dominant_component_fraction": (
                    components[0] / len(cream) if cream else 0.0
                ),
                "cream_centroid_px": (
                    [sum(x for x, _ in cream) / len(cream),
                     sum(y for _, y in cream) / len(cream)]
                    if cream else None
                ),
            }
        output[view_name] = view_rows
    return output


def _paired_mask_summary(masks):
    def summary(view_name, object_name):
        row = masks[view_name][object_name]
        counts = row["first_hit_counts"]
        visible = counts.get("cream", 0) + counts.get("black", 0)
        return {
            "cream_first_hits": counts.get("cream", 0),
            "black_first_hits": counts.get("black", 0),
            "cream_fraction_of_visible_pod": (
                counts.get("cream", 0) / visible if visible else None
            ),
            "cream_components": row["cream_connected_components"],
            "dominant_component_fraction": (
                row["cream_dominant_component_fraction"]
            ),
            "cream_centroid_px": row["cream_centroid_px"],
            "cream_screen_horizontal_span_Wh": (
                row["cream_screen_horizontal_span_Wh"]
            ),
        }

    return {
        "front_left_right": {
            "left": summary("front", CREATED_NAMES[0]),
            "right": summary("front", CREATED_NAMES[1]),
        },
        "near_foot_three_quarter_pair": {
            "camera_plus_x_right_foot": summary(
                "three_quarter", CREATED_NAMES[1]
            ),
            "camera_minus_x_left_foot": summary(
                "three_quarter_mirror", CREATED_NAMES[0]
            ),
        },
        "interpretation": (
            "Raw paired evidence only; no automatic symmetry or visual pass."
        ),
    }


def build_foot_033():
    """Build the two in-memory partitions and return JSON diagnostics."""
    assert bpy.app.version[:2] == (5, 2)
    assert bpy.app.background
    assert Path(bpy.data.filepath).resolve() == SOURCE.resolve()
    assert _sha(SOURCE) == SOURCE_HASH
    assert _sha(WRITER_RECEIPT) == WRITER_RECEIPT_HASH
    assert _sha(RENDER_RECEIPT) == RENDER_RECEIPT_HASH
    assert bpy.context.scene.frame_current == 1
    assert bpy.context.mode == "OBJECT"
    assert not bpy.data.is_dirty, "Start from a clean frozen Head032 session"

    writer_receipt = json.loads(WRITER_RECEIPT.read_text())
    render_receipt = json.loads(RENDER_RECEIPT.read_text())
    assert writer_receipt["candidate_sha256"] == SOURCE_HASH
    assert render_receipt["candidate_sha256"] == SOURCE_HASH
    assert render_receipt["contract_sha256"] == (
        "4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4"
    )

    scene = bpy.context.scene
    assert all(name in scene.objects for name in TARGETS + (HEM, RIG))
    assert all(name not in bpy.data.objects for name in CREATED_NAMES)
    assert all(name not in bpy.data.meshes for name in MESH_NAMES)
    assert COLLECTION_NAME not in bpy.data.collections
    assert BLACK_MATERIAL in bpy.data.materials
    assert CREAM_MATERIAL in bpy.data.materials
    rig = scene.objects[RIG]
    assert rig.type == "ARMATURE"
    assert all(name in rig.pose.bones for name in BONES)

    for name in TARGETS + (HEM,):
        actual = _receipt_control_record(scene.objects[name])
        assert actual == writer_receipt["controls"][name], name
        assert actual["geometry"] == EXPECTED_GEOMETRY[name]
    for target, bone_name in zip(TARGETS, BONES):
        obj = scene.objects[target]
        assert obj.type == "MESH"
        assert [slot.material.name for slot in obj.material_slots] == [
            BLACK_MATERIAL
        ]
        assert obj.parent is rig
        assert [modifier.type for modifier in obj.modifiers] == ["ARMATURE"]
        assert obj.modifiers[0].object is rig
        assert [group.name for group in obj.vertex_groups] == [bone_name]
        assert all(
            len(vertex.groups) == 1
            and obj.vertex_groups[vertex.groups[0].group].name == bone_name
            and abs(vertex.groups[0].weight - 1.0) < 1e-6
            for vertex in obj.data.vertices
        )

    controls = {
        obj.name: _object_record(obj)
        for obj in scene.objects if obj.name not in TARGETS
    }
    target_records = {
        name: _receipt_control_record(scene.objects[name]) for name in TARGETS
    }
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
        collection = bpy.data.collections.new(COLLECTION_NAME)
        assert collection.name == COLLECTION_NAME
        scene.collection.children.link(collection)
        for target_name, created_name, mesh_name, bone_name, side in zip(
                TARGETS, CREATED_NAMES, MESH_NAMES, BONES, SIDES):
            source = scene.objects[target_name]
            snapshot = _source_triangles(source)
            partition, construction = _clip_surface(snapshot, side)
            replacement, replacement_metrics = _make_replacement(
                source, created_name, mesh_name, bone_name, partition, rig,
                collection,
            )
            created.append(replacement)
            metrics[created_name] = {
                "source_name": target_name,
                "construction": construction,
                "replacement": replacement_metrics,
            }

        assert tuple(obj.name for obj in created) == CREATED_NAMES
        assert tuple(obj.data.name for obj in created) == MESH_NAMES
        assert collection.name == COLLECTION_NAME

        for name in TARGETS:
            scene.objects[name].hide_render = True
            scene.objects[name].hide_set(True)
        bpy.context.view_layer.update()

        masks = _first_hit_masks(scene, created, render_receipt)
        paired_masks = _paired_mask_summary(masks)
        front = masks["front"]
        front_cream = {
            name: front[name]["first_hit_counts"].get("cream", 0)
            for name in CREATED_NAMES
        }
        visible = all(value > 0 for value in front_cream.values())
        decision = (
            "front_cream_nonzero_only_requires_root_visual_review"
            if visible else "reject_fixed_field_for_front_visibility"
        )

        assert controls == {
            name: _object_record(scene.objects[name]) for name in controls
        }
        assert pose == _pose_record(rig)
        assert scene_state == _scene_record(scene)
        assert material_ids == {
            material.name: material.as_pointer()
            for material in bpy.data.materials
        }
        after_targets = {
            name: _receipt_control_record(scene.objects[name])
            for name in TARGETS
        }
        for name in TARGETS:
            assert all(
                target_records[name][key] == after_targets[name][key]
                for key in target_records[name] if key != "visibility"
            )
            assert scene.objects[name].hide_render
            assert scene.objects[name].hide_get()
        assert _sha(SOURCE) == SOURCE_HASH

        result = {
            "status": (
                "Foot033 built and mask-tested in memory; no save, render, "
                "or stage pass"
            ),
            "decision": decision,
            "source": str(SOURCE),
            "source_sha256": SOURCE_HASH,
            "writer_receipt_sha256": WRITER_RECEIPT_HASH,
            "render_receipt_sha256": RENDER_RECEIPT_HASH,
            "helper_sha256": _sha(Path(__file__).resolve()),
            "target_names": list(TARGETS),
            "created_names": [obj.name for obj in created],
            "created_mesh_names": [obj.data.name for obj in created],
            "created_collection_name": collection.name,
            "hidden_names": list(TARGETS),
            "head_width_m": WH,
            "formula": {
                "coordinates": (
                    "u=-side*(x-cx)/rx; p=(y-cy)/ry; "
                    "zeta=(z-cz)/rz"
                ),
                "field": "F=zeta+0.60*u+0.60*p-0.64",
                "cream_rule": "F>0",
                "black_rule": "F<=0",
                "sweep_count": 0,
            },
            "construction_metrics": metrics,
            "first_hit_material_masks": masks,
            "paired_first_hit_summary": paired_masks,
            "front_cream_first_hits": front_cream,
            "recorded_non_target_object_controls_unchanged": True,
            "rig_pose_unchanged": True,
            "recorded_scene_settings_unchanged": True,
            "persistent_material_ids_reused": True,
            "source_pods_preserved_and_hidden_in_active_view_layer_and_render": True,
            "active_view_layer": bpy.context.view_layer.name,
            "tradeoff": (
                "Black cap width and .06-.09 Wh cream reveal are reported, "
                "not simultaneous hard gates inside the fixed envelope."
            ),
            "limitations": [
                "No file was saved and no image was rendered.",
                "Ray masks are pixel-center samples at the fixed 512 px views.",
                "Nonzero cream first hits do not constitute visual acceptance.",
                "The retained black ground band keeps the full pod-width mask.",
                "Animation beyond the frozen pose remains unverified.",
                "Material-node, light, and world appearance equality is not "
                "audited here; the root writer owns that guard.",
            ],
        }
        json.dumps(result, sort_keys=True)
        return result
    except Exception:
        for name, (hide_render, hide_get) in target_visibility.items():
            scene.objects[name].hide_render = hide_render
            scene.objects[name].hide_set(hide_get)
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

## foot_033_writer.py

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_032_candidate.blend"
SOURCE_HASH = "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8"
BASELINE = ROOT / "head_032_candidate.blend"
BASELINE_HASH = "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8"
OUTPUT = ROOT / "foot_033_candidate.blend"
RECEIPT = ROOT / "foot_033_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "Left black stuffed foot pod",
    "Right black stuffed foot pod",
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
helper = ROOT / "foot_033_draft.py"
assert sha(helper) == "4b4ef4d8ef9a692382e9770f2503be7a2268b5b101bd964d77a8f5aa8b194a34"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
result = scope["build_foot_033"]()
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
scene["candidate_status"] = "Foot033 unreviewed sewn panel study; no stage pass"
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
           "appearance_unchanged": True, "foot_construction": result,
           "limitations": ["No visual, animation, whole-scene technical or final acceptance."]}
assert sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {"controls", "rig_pose", "appearance", "foot_construction"}}))

```

## render_foot_033.py

```python
"""Frozen033 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "foot_033_candidate.blend"
OUT = ROOT / "foot_033_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "head_032_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "foot_033_writer_receipt.json").read_text())
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
    data = bpy.data.cameras.new("Frozen033_" + view)
    camera = bpy.data.objects.new("Frozen033_" + view, data)
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


