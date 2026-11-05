# Reproduce028 upper-crown refinement

Observed candidate c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6.
Comparison baseline027c19ce0cb14c7d679750422702d6df97753480cf8a4db7cd73f1203b0f28bf7416.
Construction input02656efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882.

Restore the exact sources below to out/reimu_fumo_finish/desktop_astra in
the applicable linked workspace. Prior027 evidence head_027c_sources.md
provides the preceding recipe and actual-material baseline setup. Keep
both026 and027c immutable, including027c's render_receipt.json. Use the
repository-pinned Blender5.2.1 LTS build9e2066aef7ef through
//tools/blender:blender, background/factory-startup/disable-autoexec,
threads4/python-exit-code2, running head_028_writer.py then render_head_028.py.
Set TMPDIR and XDG_CACHE_HOME to existing task-specific tmp/cache directories.
Never overwrite an existing candidate, receipt, or review directory.

Regenerated Blender serialization may differ. Bind any new asset and its
receipts to fresh hashes and fresh visual/technical review rather than
claiming the recorded file hash or acceptance automatically. The formula
changes only upper-front depth; implementation metadata/object names identify
028. Evidence of actual preserved regions belongs to the independent audit,
not merely this source-level intention. Full goal remains unfinished.

## head_028_writer.py

SHA256 1bdff1a0a9fad55fe9b14d1c728e4862a0995612002a7f28a7a0dfe3f0fdac53

```python
"""Root sole writer; target list must be explicitly reviewed before running."""

import hashlib
import json
import math
import struct
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "hem_026_candidate.blend"
SOURCE_HASH = "56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882"
BASELINE = ROOT / "head_027c_candidate.blend"
BASELINE_HASH = "19ce0cb14c7d679750422702d6df97753480cf8a4db7cd73f1203b0f28bf7416"
OUTPUT = ROOT / "head_028_candidate.blend"
RECEIPT = ROOT / "head_028_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "Head_Gusseted_Cushion_020b",
    "Hair_Continuous_Traced_Fringe_020b",
    "A44 tiny neutral embroidered mouth dash",
    "A45 left flush composite eye applique",
    "A45 right flush composite eye applique",
    "A45 left drooped half-lid stitch",
    "A45 right drooped half-lid stitch",
    "A45 left fine upper expression stitch",
    "A45 right fine upper expression stitch",
    "A45 left tapered flexible cheek lock",
    "A45 right tapered flexible cheek lock",
    "Rear_Center_Cloth_021", "Rear_Left_Cloth_021", "Rear_Right_Cloth_021",
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
helper = ROOT / "head_028_draft.py"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
result = scope["build_head_028"]()
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
scene["candidate_status"] = "Head028 unreviewed construction study; no stage pass"
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
           "appearance_unchanged": True, "head_construction": result,
           "limitations": ["No visual, animation, whole-scene technical or final acceptance."]}
assert sha(SOURCE) == SOURCE_HASH
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({k: v for k, v in receipt.items() if k not in {"controls", "rig_pose", "appearance", "head_construction"}}))
```

## head_028_draft.py

SHA256 909a7fdb53cd1c75db5e73b9dcd509fc27e01f281c97fe3386fe161d80f799a3

```python
"""Unexecuted027 study. Import defines helpers only; coordinator owns save."""

import hashlib
import math
import struct
from pathlib import Path

import bpy
from mathutils import Vector
from mathutils.bvhtree import BVHTree

TARGETS = (
    'Head_Gusseted_Cushion_020b',
    'Hair_Continuous_Traced_Fringe_020b',
    'A44 tiny neutral embroidered mouth dash',
    'A45 left flush composite eye applique',
    'A45 right flush composite eye applique',
    'A45 left drooped half-lid stitch',
    'A45 right drooped half-lid stitch',
    'A45 left fine upper expression stitch',
    'A45 right fine upper expression stitch',
    'A45 left tapered flexible cheek lock',
    'A45 right tapered flexible cheek lock',
    'Rear_Center_Cloth_021', 'Rear_Left_Cloth_021', 'Rear_Right_Cloth_021',
)
SOURCE_HASH = '56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882'
WH = .1165
CUTOFF = .089
CZ = .13590002432465553
RX = .05816962569952011
RZ = .05732078477740287


def _smooth(value):
    value = min(1., max(0., value))
    return value * value * (3. - 2. * value)


def _blend(z):
    return _smooth((z - CUTOFF) / .018)


def _snapshot(obj, tag=False, omit_solidify=False):
    """Disposable evaluation copy: never change source data or modifiers."""
    copy = obj.copy()
    copy.data = obj.data.copy()
    bpy.context.scene.collection.objects.link(copy)
    data = copy.data
    try:
        if tag:
            assert len(data.polygons) == 2496
            for face in data.polygons:
                face.material_index = int(face.index < 1248)
        if omit_solidify:
            for modifier in copy.modifiers:
                if modifier.type == 'SOLIDIFY':
                    modifier.show_viewport = False
                    modifier.show_render = False
        bpy.context.view_layer.update()
        ev = copy.evaluated_get(bpy.context.evaluated_depsgraph_get())
        mesh = ev.to_mesh()
        try:
            return {
                'points': [ev.matrix_world @ v.co for v in mesh.vertices],
                'faces': [tuple(p.vertices) for p in mesh.polygons],
                'indices': [p.material_index for p in mesh.polygons],
                'materials': [bpy.data.materials[m.name] if m else None
                              for m in mesh.materials],
            }
        finally:
            ev.to_mesh_clear()
    finally:
        bpy.data.objects.remove(copy, do_unlink=True)
        if data.users == 0:
            bpy.data.meshes.remove(data)


def _tree(data):
    return BVHTree.FromPolygons(data['points'], data['faces'])


def _hit(tree, x, z, rear=False, allowance=.0025):
    direction = Vector((0, -1 if rear else 1, 0))
    hit = tree.ray_cast(Vector((x, .2 if rear else -.2, z)), direction, .4)
    if hit[0] is not None:
        return hit[0], hit[2], 0.
    # Only silhouette edge queries may lack an exact front projection.
    # A nearest 3D receiver is bounded and reported, not called contact.
    nearest = tree.find_nearest(Vector((x, -.023 if not rear else .005, z)))
    assert nearest[0] is not None, ('No receiver', x, z, rear)
    lateral = math.hypot(nearest[0].x - x, nearest[0].z - z)
    assert lateral <= allowance, ('Outside receiver domain', x, z, lateral)
    return nearest[0], nearest[2], lateral


def _bounds(points):
    return [[min(p[k] for p in points), max(p[k] for p in points)]
            for k in range(3)]


def _hash_faces(points, faces):
    digest = hashlib.sha256()
    for face in faces:
        digest.update(struct.pack('<I', len(face)))
        for i in face:
            digest.update(struct.pack('<3f', *points[i]))
    return digest.hexdigest()


def _make(name, points, faces, materials, indices=None):
    assert name not in bpy.data.objects, ('No overwrite', name)
    mesh = bpy.data.meshes.new(name + ' mesh')
    mesh.from_pydata(points, [], faces)
    mesh.materials.clear()
    for material in materials:
        mesh.materials.append(material)
    mesh.update()
    for p in mesh.polygons:
        p.use_smooth = True
        if indices is not None:
            p.material_index = indices[p.index]
    obj = bpy.data.objects.new(name, mesh)
    bpy.context.scene.collection.objects.link(obj)
    obj['construction_attempt'] = '028; unreviewed'
    return obj


def _bind_head(obj):
    """Preserve rest-frame world coordinates, with one Head bone influence."""
    rig = bpy.data.objects['ReimuFumoRig']
    pose = rig.pose.bones['Head']
    deformation = (rig.matrix_world @ pose.matrix
                   @ pose.bone.matrix_local.inverted()
                   @ rig.matrix_world.inverted())
    inverse = deformation.inverted()
    for vertex in obj.data.vertices:
        vertex.co = inverse @ vertex.co
    group = obj.vertex_groups.new(name='Head')
    group.add(list(range(len(obj.data.vertices))), 1., 'REPLACE')
    modifier = obj.modifiers.new('028 compensated Head attachment', 'ARMATURE')
    modifier.object = rig
    modifier.use_vertex_groups = True
    obj['attachment_limitation'] = 'Rest frame checked; animation unverified'


def _shell(obj, thickness):
    modifier = obj.modifiers.new('028 sewn cloth thickness', 'SOLIDIFY')
    modifier.thickness = thickness
    modifier.offset = -1.
    modifier.use_even_offset = True


def _seam(p):
    vertical = min(1., abs((p.z - CZ) / RZ))
    width = .004 + .007 * (1. - vertical * vertical)
    middle = -.020 - .004 * (1. - vertical * vertical)
    return middle - width / 2., middle + width / 2., width


def _upper_profile(z):
    z0, z1 = .125, CZ + RZ
    y0, y1 = -.03915 - .36 * (z0 - .09168), -.022
    length = z1 - z0
    t = (z - z0) / length
    return ((2*t**3 - 3*t**2 + 1)*y0
            + (t**3 - 2*t**2 + t)*length*(-.36)
            + (-2*t**3 + 3*t**2)*y1
            + (t**3 - t**2)*length*.85)


def _panel_point(p, front):
    if p.z <= CUTOFF:
        return p.copy()
    blend = _blend(p.z)
    r = min(1., math.hypot(p.x / RX, (p.z - CZ) / RZ))
    front_y, back_y, unused = _seam(p)
    radial = 1. - .012 * _smooth((r - .92) / .08) * blend
    x = p.x * radial
    z = CZ + (p.z - CZ) * radial
    if front:
        # Softly sloped mouth-to-eye region, with mild transverse stuffing.
        broad = (_upper_profile(p.z) if p.z > .125
                 else max(-.0525, -.03915 - .36 * (p.z - .09168)))
        broad += .0018 * min(1., (p.x / .04) ** 2)
        edge = _smooth((r - .72) / .28)
        desired = broad * (1. - edge) + front_y * edge
    else:
        desired = back_y + (.031984440982341766 - back_y) * (
            max(0., 1. - r * r) ** .62)
    return Vector((x, p.y + blend * (desired - p.y),
                   p.z + blend * (z - p.z)))


def _core(old, tagged):
    points, faces = old['points'], old['faces']
    assert faces == tagged['faces'] and len(points) == 9794
    front_faces = [i for i, value in enumerate(tagged['indices']) if value == 1]
    front_vertices = {j for i in front_faces for j in faces[i]}
    back_faces = [i for i, value in enumerate(tagged['indices']) if value == 0]
    back_vertices = {j for i in back_faces for j in faces[i]}
    seam = front_vertices & back_vertices
    assert len(seam) == 192
    seam = sorted(seam, key=lambda i: math.atan2(
        (points[i].z - CZ) / RZ, points[i].x / RX))
    output = [_panel_point(p, i in front_vertices) for i, p in enumerate(points)]
    rear_map = {}
    for i in seam:
        if points[i].z <= CUTOFF:
            rear_map[i] = i
        else:
            rear_map[i] = len(output)
            output.append(_panel_point(points[i], False))
    new_faces = [tuple(rear_map.get(i, i) if tagged['indices'][fi] == 0 else i
                       for i in face) for fi, face in enumerate(faces)]
    # Below the retained patch keep source skin/hair assignment, too.
    indices = [1 - old['indices'][i] if max(points[j].z for j in face) <= CUTOFF
               else 0 for i, face in enumerate(faces)]
    rings = [list(seam)]
    for level in range(1, 8):
        t = level / 8.
        ring = []
        for i in seam:
            p = points[i]
            if p.z <= CUTOFF:
                ring.append(i)
                continue
            blend = _blend(p.z)
            a, b, unused = _seam(p)
            radial = 1. - .012 * (1. - 4. * t * (1. - t)) * blend
            q = Vector((p.x * radial,
                        p.y + blend * ((1. - t) * a + t * b - p.y),
                        p.z + blend * ((CZ + (p.z - CZ) * radial) - p.z)))
            ring.append(len(output))
            output.append(q)
        rings.append(ring)
    rings.append([rear_map[i] for i in seam])
    gusset_start = len(new_faces)
    for a, b in zip(rings, rings[1:]):
        for j in range(len(seam)):
            k = (j + 1) % len(seam)
            face = tuple(dict.fromkeys((a[j], a[k], b[k], b[j])))
            if len(face) < 3:
                continue
            center = sum((output[i] for i in face), Vector()) / len(face)
            normal = (output[face[1]] - output[face[0]]).cross(
                output[face[2]] - output[face[0]])
            if normal.dot(Vector((center.x, 0., center.z - CZ))) < 0:
                face = tuple(reversed(face))
            new_faces.append(face)
            indices.append(0)
    preserved = [i for i, face in enumerate(faces)
                 if max(points[j].z for j in face) <= CUTOFF]
    before = _hash_faces(points, [faces[i] for i in preserved])
    after = _hash_faces(output, [new_faces[i] for i in preserved])
    assert before == after, 'Exact retained underside faces changed'
    data = {'points': output, 'faces': new_faces}
    obj = _make('Head028 sewn cushion', output, new_faces,
                [bpy.data.materials['Face fabric clay'],
                 bpy.data.materials['Hair brown clay']], indices)
    obj['preserved_world_z_max_m'] = CUTOFF
    obj['front_seam_vertices'] = len(seam)
    obj['gusset_band_count'] = 8
    hood_faces = [new_faces[i] for i in back_faces
                  if max(output[j].z for j in new_faces[i]) > CUTOFF]
    hood_faces += [new_faces[i] for i in front_faces if old['indices'][i] == 0]
    hood_faces += new_faces[gusset_start:]
    used = sorted({i for face in hood_faces for i in face})
    mapping = {i: j for j, i in enumerate(used)}
    normals = {i: Vector() for i in used}
    for face in hood_faces:
        normal = (output[face[1]] - output[face[0]]).cross(
            output[face[2]] - output[face[0]])
        for i in face:
            normals[i] += normal
    hood_points = []
    for i in used:
        normal = normals[i].normalized()
        padding = .00055 * _smooth((output[i].z - CUTOFF) / .008)
        hood_points.append(output[i] + normal * padding)
    hood = _make('Hair028 crown and back hood', hood_points,
                 [tuple(mapping[i] for i in face) for face in hood_faces],
                 [bpy.data.materials['Hair brown clay']])
    _shell(hood, .0007)
    return obj, hood, data, {
        'preserved_face_count': len(preserved),
        'preserved_face_coordinate_sha256': before,
        'preservation_cutoff_world_z_m': CUTOFF,
        'seam_width_m': [.004, .011],
        'seam_projection_inset_fraction': .012,
        'gusset_bands': 8, 'seam_samples': len(seam),
        'source_bounds_xyz_m': _bounds(points),
        'new_core_bounds_xyz_m': _bounds(output),
    }


def _fit_details(old, old_tree, new_tree, name, rear=False):
    output = []
    max_lateral = 0.
    overhang_count = 0
    shifts = []
    for p in old['points']:
        factor = _smooth((p.z - .112) / .044) if rear else 1.
        if factor == 0.:
            output.append(p.copy())
            shifts.append(0.)
            continue
        old_hit = old_tree.ray_cast(Vector((p.x, .2, p.z)), Vector((0, -1, 0)), .4) if rear else None
        new_hit = new_tree.ray_cast(Vector((p.x, .2, p.z)), Vector((0, -1, 0)), .4) if rear else None
        if rear and (old_hit[0] is None or new_hit[0] is None):
            a = old_tree.find_nearest(p)[0]
            assert a is not None
            b = _panel_point(a, False)
            lateral_a = lateral_b = math.hypot(a.x - p.x, a.z - p.z)
            overhang_count += 1
        else:
            a, unused, lateral_a = _hit(old_tree, p.x, p.z, rear)
            b, unused, lateral_b = _hit(new_tree, p.x, p.z, rear)
        delta = (b.y - a.y + (.00055 if rear else 0.)) * factor
        output.append(Vector((p.x, p.y + delta, p.z)))
        shifts.append(delta)
        max_lateral = max(max_lateral, lateral_a, lateral_b)
    obj = _make(name, output, old['faces'], old['materials'], old['indices'])
    return obj, {'max_xz_shift_m': 0., 'depth_shift_m': [min(shifts), max(shifts)],
                 'overhang_displacement_transfer_count': overhang_count,
                 'max_surface_query_lateral_fallback_m': max_lateral,
                 'method': 'Per-XZ surface depth delta retains depth residual and topology'}


def _fringe(old, old_tree, new_tree):
    output = []
    max_lateral = 0.
    overhang_count = 0
    for p in old['points']:
        # Reconstruct one continuous depth field across the traced outline.
        # The inherited receiver-hit/overhang boundary is not a cloth seam.
        r = min(1., math.hypot(p.x / RX, (p.z - CZ) / RZ))
        old_depth = -.008 - .045 * max(0., 1. - r * r) ** .40
        receiver = _panel_point(Vector((p.x, old_depth, p.z)), True)
        # Unequal broad raised lobes, tapered into the actual crown seam.
        crown = CZ + RZ * math.sqrt(max(0., 1. - (p.x / RX) ** 2))
        seam_distance = max(0., crown - p.z)
        root_taper = _smooth(seam_distance / .008)
        tip_taper = _smooth(max(0., p.z - .1035) / .007)
        lobe = .0004 + .00125 * (.6 + .4 * math.cos(p.x / WH * 5. * math.pi) ** 2)
        pad = lobe * root_taper * tip_taper
        output.append(Vector((p.x, receiver.y - .00085 - pad, p.z)))
    obj = _make('Hair028 traced padded fringe', output, old['faces'],
                old['materials'], old['indices'])
    _shell(obj, .0011)
    return obj, {'max_xz_shift_m': 0., 'padding_m': [.0004, .00165],
                 'root_taper_m': .008, 'cloth_thickness_m': .0011,
                 'overhang_displacement_transfer_count': overhang_count,
                 'depth_method': 'Continuous analytic panel field plus .85mm cloth stand-off; contact unverified',
                 'max_surface_query_lateral_fallback_m': max_lateral}


def _cheek(old, old_tree, new_tree, side):
    points, faces = old['points'], old['faces']
    tree = _tree(old)
    edges = set()
    for face in faces:
        for a, b in zip(face, face[1:] + face[:1]):
            edges.add(tuple(sorted((a, b))))
    zmin, zmax = min(p.z for p in points), max(p.z for p in points)
    rows = {zmin + (zmax - zmin) * i / 40. for i in range(41)}
    rows.update([.10156, .10249, .11098, .11197])
    rows.update([min(points, key=lambda p: p.x).z,
                 max(points, key=lambda p: p.x).z])
    rows = sorted(z for z in rows if zmin <= z <= zmax)
    columns = 17
    front, back, material_indices = [], [], []
    witnesses = []
    for z in rows:
        crossings = [p.x for p in points if abs(p.z - z) < 1e-7]
        for a, b in edges:
            pa, pb = points[a], points[b]
            if min(pa.z, pb.z) <= z <= max(pa.z, pb.z) and abs(pa.z - pb.z) > 1e-9:
                crossings.append(pa.x + (pb.x - pa.x) * (z - pa.z) / (pb.z - pa.z))
        assert crossings, ('No cheek outline slice', side, z)
        lo, hi = min(crossings), max(crossings)
        if hi - lo < .00005:
            lo -= .000025
            hi += .000025
        for j in range(columns):
            u = -1. + 2. * j / (columns - 1)
            x = lo + (hi - lo) * j / (columns - 1)
            a = tree.ray_cast(Vector((x, -.2, z)), Vector((0, 1, 0)), .4)
            b = tree.ray_cast(Vector((x, .2, z)), Vector((0, -1, 0)), .4)
            if a[0] is not None and b[0] is not None:
                middle = (a[0].y + b[0].y) / 2.
            else:
                near = min(points, key=lambda p: (p.x - x) ** 2 + (p.z - z) ** 2)
                middle = near.y
            root = _smooth((z - .112) / .012)
            if root > 0:
                receiver, unused, lateral = _hit(new_tree, x, z)
                middle = middle * (1. - root) + (receiver.y - .0007) * root
                if z >= .120 and j == columns // 2:
                    witnesses.append({'x': x, 'z': z, 'receiver_y': receiver.y,
                                      'cloth_mid_y': middle, 'lateral_fallback': lateral})
            t = (z - zmin) / (zmax - zmin)
            fill = .0022 * max(0., 1. - u * u) ** .7 * math.sin(math.pi * t) ** .5
            fill *= 1. - .7 * _smooth((z - .117) / .009)
            half = .0004 + fill
            front.append(Vector((x, middle - half, z)))
            back.append(Vector((x, middle + half, z)))
    output = front + back
    count = len(front)
    new_faces = []
    for i in range(len(rows) - 1):
        for j in range(columns - 1):
            a = i * columns + j
            quad = (a, a + 1, a + columns + 1, a + columns)
            center = sum((front[k] for k in quad), Vector()) / 4.
            hit = tree.ray_cast(Vector((center.x, -.2, center.z)), Vector((0, 1, 0)), .4)
            mi = old['indices'][hit[2]] if hit[0] is not None else 0
            new_faces.append(quad)
            material_indices.append(mi)
            new_faces.append(tuple(k + count for k in reversed(quad)))
            material_indices.append(mi)
    boundary = list(range(columns))
    boundary += [i * columns + columns - 1 for i in range(1, len(rows))]
    boundary += [(len(rows) - 1) * columns + j for j in range(columns - 2, -1, -1)]
    boundary += [i * columns for i in range(len(rows) - 2, 0, -1)]
    for a, b in zip(boundary, boundary[1:] + boundary[:1]):
        new_faces.append((a, a + count, b + count, b))
        material_indices.append(0)
    obj = _make(f'Hair028 {side} padded cheek lock', output, new_faces,
                old['materials'], material_indices)
    return obj, {'source_bounds_xyz_m': _bounds(points),
                 'new_bounds_xyz_m': _bounds(output), 'rows': len(rows),
                 'columns': columns, 'edge_thickness_m': .0008,
                 'maximum_center_thickness_m': .0052,
                 'root_witnesses': witnesses,
                 'construction': 'Two stuffed cloth panels sharing a sewn rolled edge'}


def build_head_028():
    source = Path(bpy.data.filepath).resolve()
    assert source.name == 'hem_026_candidate.blend'
    assert hashlib.sha256(source.read_bytes()).hexdigest() == SOURCE_HASH
    assert bpy.context.scene.frame_current == 1
    assert all(name in bpy.context.scene.objects for name in TARGETS)
    assert all(bpy.data.objects[name].visible_get() and not bpy.data.objects[name].hide_render
               for name in TARGETS)
    assert not any(obj.name.startswith(('Head028 ', 'Hair028 ')) for obj in bpy.data.objects)
    snapshots = {name: _snapshot(bpy.data.objects[name]) for name in TARGETS}
    print('028 phase: source snapshots complete', flush=True)
    tagged = _snapshot(bpy.data.objects[TARGETS[0]], tag=True)
    old = snapshots[TARGETS[0]]
    old_tree = _tree(old)
    core, hood, core_data, core_receipt = _core(old, tagged)
    new_tree = _tree(core_data)
    created = [core, hood]
    receipt = {'core': core_receipt, 'details': {}, 'cheek_locks': {}}
    fringe_mid = _snapshot(bpy.data.objects[TARGETS[1]], omit_solidify=True)
    obj, info = _fringe(fringe_mid, old_tree, new_tree)
    created.append(obj)
    receipt['fringe'] = info
    print('028 phase: core, hood and fringe complete', flush=True)
    detail_names = {
        'A44 tiny neutral embroidered mouth dash': 'Head028 mouth embroidery',
    }
    for side in ('left', 'right'):
        detail_names.update({
            f'A45 {side} flush composite eye applique': f'Head028 {side} eye applique',
            f'A45 {side} drooped half-lid stitch': f'Head028 {side} half lid embroidery',
            f'A45 {side} fine upper expression stitch': f'Head028 {side} upper expression embroidery',
        })
    for old_name, new_name in detail_names.items():
        obj, info = _fit_details(snapshots[old_name], old_tree, new_tree, new_name)
        created.append(obj)
        receipt['details'][new_name] = info
    for side in ('left', 'right'):
        obj, info = _cheek(snapshots[f'A45 {side} tapered flexible cheek lock'],
                           old_tree, new_tree, side)
        created.append(obj)
        receipt['cheek_locks'][side] = info
    print('028 phase: face details and cheek locks complete', flush=True)
    for key in ('Center', 'Left', 'Right'):
        name = f'Hair028 rear {key.lower()} cloth'
        obj, info = _fit_details(snapshots[f'Rear_{key}_Cloth_021'],
                                 old_tree, new_tree, name, rear=True)
        created.append(obj)
        receipt['details'][name] = info
    print('028 phase: rear cloth complete', flush=True)
    for obj in created:
        print('028 evaluating attachment: ' + obj.name, flush=True)
        _bind_head(obj)
        bpy.context.view_layer.update()
    print('028 phase: Head attachments complete', flush=True)
    # Originals are kept recoverable; the sole intended scene mutation scope.
    for name in TARGETS:
        obj = bpy.data.objects[name]
        print('028 hiding original: ' + name, flush=True)
        obj.hide_render = True
        obj.hide_viewport = True
        obj.hide_set(True)
        bpy.context.view_layer.update()
    bpy.context.view_layer.update()
    print('028 phase: visibility update complete', flush=True)
    receipt.update({
        'source_sha256': SOURCE_HASH, 'hidden_exact_targets': list(TARGETS),
        'created_names': [obj.name for obj in created],
        'created_bounds_xyz_m': {obj.name: _bounds(_snapshot(obj)['points']) for obj in created},
        'Wh_m': WH, 'material_nodes_changed': False,
        'contact_statement': 'Original underside faces wholly below89mm retained exactly; no new contact acceptance inferred',
        'attachment': 'Full Head weights with inverse current-pose compensation',
        'limitations': ['Unreviewed geometry study', 'No animation acceptance',
                        'Root and continuity witnesses require frozen candidate audit'],
    })
    assert hashlib.sha256(source.read_bytes()).hexdigest() == SOURCE_HASH
    return receipt
```

## render_head_028.py

SHA256 aedc310550683794a05ad5d137bb37c0bb400fb528c4d5e83414e02e24a44f25

```python
"""Frozen028 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_028_candidate.blend"
OUT = ROOT / "head_028_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "head_027c_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "head_028_writer_receipt.json").read_text())
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
    data = bpy.data.cameras.new("Frozen028_" + view)
    camera = bpy.data.objects.new("Frozen028_" + view, data)
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
