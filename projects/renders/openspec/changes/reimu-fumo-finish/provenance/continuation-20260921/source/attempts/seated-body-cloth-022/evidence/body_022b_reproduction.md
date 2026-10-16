# Body022b reproducibility evidence

Audit identities: final technical JSON SHA-256
`37410489831ada543ad1cf569a67cd3237db8f7e5dd5a5dbc7c8fea873c2e7f9`;
first022 technical JSON SHA-256
`a68fa3ee0f0294c4cb9e1e420fcf49506882e8e9a52dfff7889b1830a890ed15`;
final render receipt SHA-256
`73153da72d119b156e1b79aa0a2154030b869e9bd9fe21f9c19f5ca2d5c5b18d`.

## Lattice spacing repair, before retry

The first disposable fixture failed before Reimu was opened. Intended
smooth displacement differed by 10.813 mm; minimum movement was only
3.873 mm. No model candidate was written.

A separate pinned Blender spacing diagnostic returned 256 points and
base-coordinate bounds X/Y [-0.5, 0.5], Z [-31.5, 31.5]. Requested W=65
is clamped to 64. The script had incorrectly treated all axes as unit
extent, stretching the height cage to 18.9 m rather than 0.30 m.

Use each observed base-coordinate span to normalize the cage transform.
Record actual lattice resolution. Predicted result: the disposable mesh
and curve follow the frozen height field within 0.1 mm, X/Y stay fixed,
and modifier removal restores exact original coordinates. This consumes
the single preauthorized setup repair; no visual parameter changes.

## Correct two implementation defects before another visual decision

Frozen first state 022 is rejected as a completed model. Independent
reference-first review identifies repeated horizontal skirt bands and
preexisting disconnected collar/tie, helmet hair, and rigid accessories.
Technical review identifies a mouth classification bug: the low mouth
fell below the Z threshold and received the body field, not head motion.

The skirt sampler uses a separate zero-slope smoothstep at every drape
row, creating repeated shelves despite monotone input heights. Replace
that sampler with monotone cubic Hermite interpolation carrying nonzero
slopes through those rows. Preserve every original drape control point,
waist and hem span, support clearance, and deformation parameter.
Assign the exact mouth name to the head field. No other role changes.

Save as a separate 022b candidate, rebuilt from protected 021. Keep 022
and its negative review intact. This is a source-bug correction within
the same body hypothesis, not a new design or a claimed stage pass.
Reuse the successful native fixture evidence; do not repeat the settled
capability test. Compare all five frozen views and repeat independent
review on corrected bytes. Defer collar/torso redesign to the next cycle.

## Body022 continuation review

Subject: body022b, SHA-256
`96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014`.
This is intermediate work, not a completed asset. Goal source RV71 before
the closing checkpoint. Root uses the ergonomics-review skill to preserve
these bounded lessons; no shared instructions or host settings are changed.

- `FUMO-LATTICE-EXTENT`: fixture-tested. Blender's 64-row native lattice
  has Z coordinates [-31.5,31.5], not a unit cube. One failed fixture plus
  one separate spacing diagnostic identified this. Normalize observed
  base-coordinate spans; corrected mesh/curve tests recover exactly.
- `FUMO-ROLE-LOW-MOUTH`: live saved evidence. Height-only role classification
  misplaced the low mouth by 5.5 mm relative to the head. Exact-name
  classification fixes it. Use explicit semantic sets for later module
  moves; a geometric threshold is insufficient for attachments.
- `FUMO-DRAPE-STATIONARY-ROWS`: live pixels. Zero derivatives at every
  piecewise-smooth drape row created horizontal terraces. Monotone cubic
  interpolation removed the bands while retaining the specified rows.
- `FUMO-REVIEW-COST`: two complete 5-view pinned packets took about 7 seconds
  each including save/reopen. Background scripting was adequate after the
  bounded live session expired. Re-proving live GUI state was unnecessary.
- `GOAL-EMPTY-EVIDENCE`: routed validation. An empty evidence directory
  disappears from Git/Bazel inputs. Add actual decision evidence when
  opening an attempt; do not rely on an empty directory for portability.
- `CORDIS-CATALOG-SIZE`: routed. Printing the full duplicated catalog caused
  a roughly 12k-token output truncation; subsequent calls projected only
  structured needed fields. Continue using bounded projections.

Independent pixel review, technical QA, and normalized measurements ran in
parallel on frozen files. CURRENT.md carries exact next work and latest
artifact identities. Do not revisit settled screenshots or start another
window investigation unless an actual required operation needs it.

## Body 022 and 022b technical audit

Read-only clean reopen with pinned Blender 5.2.1 LTS, build
`9e2066aef7ef`, at rest frame 1. No blend file was saved or rendered.

### Exact candidates

- 022: `0ff8cab54dfa5be530a12524f86e7238c13f90920ea3fd1864a27a6e6cdeef92`.
- 022b: `96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014`.
- Source 021: `9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22`.
- A202 donor: `a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5`.

Candidate, source, and donor hashes remained unchanged in both audits.

### 022 findings

All three terminal lattice roles are attached correctly to their declared
objects, with finite 2 x 2 x 64 cages. Their role counts are head 13, bow 15,
and body 54. Existing XY geometry is unchanged. Maximum field errors are
0.0000463 mm for the head, 0.388880 mm for the bow, and 0.047022 mm for the
body; bow/body fields are linearly sampled by the cages.

The declaration itself contains a regression: the mouth is classified as
body. It moves down 16.6 mm instead of the head's 22.1 mm, producing a
5.5000002 mm upward displacement relative to the face. All fourteen head and
face base meshes remain intrinsically unchanged; this does not prevent a
relative placement regression in their evaluated result.

Foot, floor, and inspected rig records exactly match source 021. The rig
record covers world transform, parent, visibility, modifiers, rest-bone
matrices/parents/lengths/deform flags, rest-frame pose matrices, and active
action identity. It is not a complete animated-rig audit.

Both new cloth bases have 320 expected boundary edges and Euler
characteristic 0. Their evaluated solidified results are closed manifold
annular cloth volumes, also Euler characteristic 0, with consistent winding,
finite coordinates, no zero-length edges, and no degenerate faces. All base
vertices have positive skin weights. This does not establish useful motion.

#### 022 placement measurements

| Measurement | Result |
| --- | ---: |
| Head minimum Z minus collar maximum Z | -5.204 mm |
| Minimum sampled left collar-to-head distance | 18.202 mm |
| Minimum sampled right collar-to-head distance | 18.160 mm |
| Minimum sampled skirt-to-left/right-foot distance | 1.577 / 1.577 mm |
| Minimum sampled hem-to-left/right-foot distance | 1.328 / 1.434 mm |
| Hem minimum Z above floor | 3.181 mm |
| Each foot minimum Z above floor | 0.100 mm |

Collar height overlap does not establish seated contact. Proximity values
sample every source vertex against target triangles; they do not prove an
exhaustive global minimum or collision-free geometry. No sampled skirt or
hem vertex lies more than 0.1 mm behind its nearest foot surface normal.

### 022b correction checks

The mouth is now explicitly one of fourteen head-role objects. Its rigid
translation differs from the requested -22.1 mm by at most
0.0000165 mm. The largest error across all fourteen head objects is
0.0000463 mm, with zero XY drift. All head base meshes and protected
foot/floor/inspected rig records still match source 021.

The corrected skirt and unchanged hem pass the same finite-coordinate,
consistent-winding, boundary-to-closed-topology, zero-edge, degenerate-face,
and positive-weight checks. The skirt evaluates to 62,080 vertices/faces;
the hem evaluates to 16,000 vertices/faces.

Contact sampling was not repeated for 022b. Neither audit establishes
reference alignment, visual acceptance, exhaustive collision-free contact,
or animated/full-rig acceptance.

Machine evidence: `body_022_technical_audit.json` and
`body_022b_technical_audit.json`. The 022b script reuses only explicitly
whitelisted helper function definitions from the 022 audit script; its
hash is recorded, and the earlier candidate inspection block is not rerun.

The following source and receipts bind the authored operation to frozen
candidate96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014.
The binary candidate and its021 input remain ignored local working assets;
this source record does not claim exact-byte regeneration in a fresh clone.

Run with repository-pinned Blender via bazel_agent, background, factory
startup, disable-autoexec, python-exit-code2, and task-owned temporary/cache
paths. Execute body_macro_022b.py then render_body_022b.py. Candidate paths
must be absent. The script requires immutable021 bytes and the prior fixture
receipt. Existing meshes/modifiers/rig are preserved; motion is unverified.

## body_skirt_022b.py

```python
"""Draft: replace only the inspected skirt and hem in the sole live writer.

Not executed or visually validated. Root must review before live execution.
The torso, internal seat, sleeves, hands, feet, head, and bow stay unchanged.
"""

import hashlib
import json
import math
import os
import struct
from pathlib import Path

import bpy
from mathutils import Matrix, Vector
from mathutils.bvhtree import BVHTree
from mathutils.kdtree import KDTree

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_021_candidate.blend"
SOURCE_HASH = "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22"
CANDIDATE = ROOT / "body_022b_candidate.blend"
RECEIPT = ROOT / "body_022b_writer_receipt.json"
LOCK = ROOT / "body_022b_writer.lock"
COLLECTION = "Skirt_Cloth_022"
WH = 0.1165
RED_WIDTH = 0.1067
HEM_WIDTH = 0.1241
TAU = 2.0 * math.pi
PANELS = (
    "Garment42 front shallow lap panel",
    "Garment42 rear pooled dress panel",
    "Garment42 side gusset left",
    "Garment42 side gusset right",
)
RUFFLES = (
    "Garment42 front gathered ruffle",
    "Garment42 rear pooled ruffle",
    "Garment42 left side gathered ruffle",
    "Garment42 right side gathered ruffle",
)
# These exact 18 names were present in the pinned read-only inspection.
STITCHES = tuple(
    f"Garment42 hem stitch {side} {i:02d}"
    for side in ("front", "rear") for i in range(9)
)
TARGETS = PANELS + RUFFLES + STITCHES
SUPPORTS = (
    "Garment42 compact internal seat pad",
    "Left short hidden leg root",
    "Right short hidden leg root",
    "Left black stuffed foot pod",
    "Right black stuffed foot pod",
)
NEW_NAMES = ("Skirt022 joined gathered panels", "Skirt022 soft hem ruffle")
SCENE_NAME = "Attempt41_Manual_Head_Maquette"
GEOMETRY_TYPES = {"MESH", "CURVE", "SURFACE", "FONT", "META"}
# Root may set this in exec globals to integrate a separately authorized warp
# and make the combined candidate/receipt save itself.
DEFER_SAVE = bool(globals().get("BODY022_DEFER_SAVE", False))


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def object_record(obj):
    """Bind controls to evaluated world geometry, topology, and visibility."""
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    record = {
        "type": obj.type,
        "data": obj.data.name if obj.data else None,
        "parent": obj.parent.name if obj.parent else None,
        "matrix_world": [list(row) for row in evaluated.matrix_world],
        "hide_render": obj.hide_render,
        "hide_viewport": obj.hide_viewport,
        "hide_set": obj.hide_get(),
    }
    if obj.type in GEOMETRY_TYPES:
        mesh = evaluated.to_mesh()
        try:
            h = hashlib.sha256()
            h.update(struct.pack("<3I", len(mesh.vertices), len(mesh.edges),
                                 len(mesh.polygons)))
            for vertex in mesh.vertices:
                xyz = evaluated.matrix_world @ vertex.co
                assert all(math.isfinite(v) for v in xyz)
                h.update(struct.pack("<3f", *xyz))
            for edge in mesh.edges:
                h.update(struct.pack("<2I", *edge.vertices))
            for face in mesh.polygons:
                h.update(struct.pack("<2I", len(face.vertices),
                                     face.material_index))
                for vertex in face.vertices:
                    h.update(struct.pack("<I", vertex))
            record["world_evaluated_sha256"] = h.hexdigest()
            record["materials"] = [m.name if m else None
                                   for m in mesh.materials]
        finally:
            evaluated.to_mesh_clear()
    if obj.type == "ARMATURE":
        record["pose"] = {bone.name: [list(row) for row in bone.matrix]
                          for bone in obj.pose.bones}
    return record


def bounds(obj):
    evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = evaluated.to_mesh()
    try:
        points = [evaluated.matrix_world @ v.co for v in mesh.vertices]
        return [[min(p[i] for p in points), max(p[i] for p in points)]
                for i in range(3)]
    finally:
        evaluated.to_mesh_clear()


def union_bounds(objects):
    rows = [bounds(obj) for obj in objects]
    return [[min(row[i][0] for row in rows), max(row[i][1] for row in rows)]
            for i in range(3)]


def support_tree(objects):
    vertices, faces = [], []
    for obj in objects:
        evaluated = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
        mesh = evaluated.to_mesh()
        try:
            offset = len(vertices)
            vertices.extend(evaluated.matrix_world @ v.co for v in mesh.vertices)
            faces.extend(tuple(offset + i for i in face.vertices)
                         for face in mesh.polygons)
        finally:
            evaluated.to_mesh_clear()
    return BVHTree.FromPolygons(vertices, faces, all_triangles=False)


def support_height(tree, x, y):
    hit, _, _, _ = tree.ray_cast(Vector((x, y, 0.12)),
                                 Vector((0.0, 0.0, -1.0)), 0.13)
    return hit.z if hit is not None else None


def supported_z(tree, x, y, z):
    """Seat and leg pressure lifts cloth; the floor flattens low ruffles."""
    limit = 0.0010
    hit = support_height(tree, x, y)
    if hit is not None:
        limit = max(limit, hit + 0.0018)
    # Rounded maximum avoids a hard crease where free drape meets support.
    d = z - limit
    return 0.5 * (z + limit + math.sqrt(d * d + 0.0007**2))


def signed_power(value, exponent):
    return math.copysign(abs(value)**exponent, value)


def angular_gaussian(angle, center, width):
    distance = (angle - center + math.pi) % TAU - math.pi
    return math.exp(-0.5 * (distance / width)**2)


def smooth_rows(t, rows):
    """Monotone cubic height profile, without stationary terrace rows."""
    if t <= rows[0][0]:
        return rows[0][1:]
    if t >= rows[-1][0]:
        return rows[-1][1:]
    h = [b[0] - a[0] for a, b in zip(rows, rows[1:])]
    segment = next(i for i in range(len(h)) if t <= rows[i + 1][0])
    u = (t - rows[segment][0]) / h[segment]
    output = []
    for coordinate in range(1, len(rows[0])):
        values = [row[coordinate] for row in rows]
        slopes = [(b - a) / step for a, b, step
                  in zip(values, values[1:], h)]
        tangent = [slopes[0]]
        for i in range(1, len(rows) - 1):
            a, b = slopes[i - 1], slopes[i]
            if a * b <= 0.0:
                tangent.append(0.0)
            else:
                w1, w2 = 2 * h[i] + h[i - 1], h[i] + 2 * h[i - 1]
                tangent.append((w1 + w2) / (w1 / a + w2 / b))
        tangent.append(slopes[-1])
        i = segment
        output.append(
            (2 * u**3 - 3 * u**2 + 1) * values[i]
            + (u**3 - 2 * u**2 + u) * h[i] * tangent[i]
            + (-2 * u**3 + 3 * u**2) * values[i + 1]
            + (u**3 - u**2) * h[i] * tangent[i + 1])
    return tuple(output)


# (waist-to-hem fraction, front height, rear height).
# The tucked upper edge meets the unchanged bodice at z=0.0621 m.
DRAPE = (
    (0.00, 0.0621, 0.0621), (0.06, 0.0590, 0.0585),
    (0.18, 0.0495, 0.0460), (0.35, 0.0398, 0.0310),
    (0.53, 0.0334, 0.0195), (0.72, 0.0290, 0.0105),
    (0.88, 0.0274, 0.0069), (1.00, 0.0268, 0.0064),
)
# Five differently sized gathers follow seam tension. Their positions and
# amplitudes are deliberately independent, with no radial fluting oscillator.
GATHERS = (
    (0.42, 0.19, 0.0014), (1.09, 0.28, -0.0011),
    (2.48, 0.31, 0.0015), (3.80, 0.23, -0.0010),
    (5.57, 0.27, 0.0018),
)
# A few relaxed ruffle folds have unequal spans and heights.
RUFFLE_FOLDS = (
    (0.05, 0.13, 0.0030), (0.48, 0.15, 0.0038),
    (0.95, 0.12, 0.0024), (1.45, 0.18, 0.0032),
    (2.02, 0.15, 0.0025), (2.55, 0.19, 0.0030),
    (3.08, 0.14, 0.0020), (3.58, 0.17, 0.0028),
    (4.14, 0.14, 0.0024), (4.67, 0.20, 0.0035),
    (5.22, 0.14, 0.0026), (5.75, 0.16, 0.0036),
)


def skirt_point(t, angle, tree):
    s, c = math.sin(angle), math.cos(angle)
    # Rounded rectangular waist follows the inspected rectangular bodice.
    waist = Vector((0.0394 * signed_power(s, 0.36),
                    -0.0045 - 0.0300 * signed_power(c, 0.36), 0.0))
    # The seated outline is rounded and deep, with a flatter front lap edge.
    hem = Vector((0.5 * RED_WIDTH * signed_power(s, 0.76),
                  -0.0020 - 0.0590 * signed_power(c, 0.76), 0.0))
    spread = t**0.84
    point = waist.lerp(hem, spread)
    front_z, rear_z = smooth_rows(t, DRAPE)
    front_weight = ((1.0 + c) * 0.5)**1.12
    z = rear_z + (front_z - rear_z) * front_weight
    gather = sum(amplitude * angular_gaussian(angle, center, width)
                 for center, width, amplitude in GATHERS)
    envelope = 4.0 * t * (1.0 - t)
    point.y += gather * envelope * c
    z += gather * envelope
    # Keep the sewn waist in the bodice; support acts on the lower cloth.
    if t > 0.12:
        z = supported_z(tree, point.x, point.y, z)
    point.z = z
    return point


def ruffle_point(u, angle, tree):
    point = skirt_point(1.0, angle, tree)
    s, c = math.sin(angle), math.cos(angle)
    point.x += 0.5 * (HEM_WIDTH - RED_WIDTH) * u * signed_power(s, 0.76)
    point.y -= 0.0045 * u * signed_power(c, 0.76)
    fold = sum(height * angular_gaussian(angle, center, width)
               for center, width, height in RUFFLE_FOLDS)
    # The seam is shared exactly. Excess outer-edge length turns into folds.
    point.z += -0.0052 * u + fold * (u**1.25)
    if u > 0.0:
        point.z = supported_z(tree, point.x, point.y, point.z)
    return point


def skin_source(names):
    """Reuse the existing skirt's nearest rest-vertex bone weights."""
    objects = [bpy.data.objects[name] for name in names]
    rig = objects[0].parent
    assert rig and rig.type == "ARMATURE"
    assert all(obj.parent == rig for obj in objects)
    points, weights = [], []
    for obj in objects:
        for vertex in obj.data.vertices:
            points.append(obj.matrix_world @ vertex.co)
            weights.append({obj.vertex_groups[g.group].name: g.weight
                            for g in vertex.groups if g.weight > 0.0})
    assert all(weights), "An inspected source garment vertex has no skin weights"
    tree = KDTree(len(points))
    for index, point in enumerate(points):
        tree.insert(point, index)
    tree.balance()
    return rig, tree, weights


def make_cloth(name, sampler, rows, material, skin, collection):
    around = 160
    vertices = [sampler(i / (rows - 1), TAU * j / around)
                for i in range(rows) for j in range(around)]
    faces = []
    for i in range(rows - 1):
        for j in range(around):
            a = i * around + j
            b = i * around + (j + 1) % around
            faces.append((a, a + around, b + around, b))
    mesh = bpy.data.meshes.new(name)
    mesh.from_pydata(vertices, [], faces)
    mesh.materials.append(material)
    mesh.update()
    obj = bpy.data.objects.new(name, mesh)
    collection.objects.link(obj)
    rig, nearest, weights = skin
    obj.parent = rig
    obj.matrix_world = Matrix.Identity(4)
    groups = {key: obj.vertex_groups.new(name=key)
              for key in sorted({k for row in weights for k in row})}
    for index, point in enumerate(vertices):
        _, source_index, _ = nearest.find(point)
        for group, weight in weights[source_index].items():
            groups[group].add([index], weight, "REPLACE")
    for polygon in mesh.polygons:
        polygon.use_smooth = True
    sub = obj.modifiers.new("Relaxed sewn surface", "SUBSURF")
    sub.levels = sub.render_levels = 1
    solid = obj.modifiers.new("Thin cotton edge", "SOLIDIFY")
    solid.thickness = 0.00065 if name == NEW_NAMES[0] else 0.0005
    solid.offset = 0.0
    armature = obj.modifiers.new("Inherited skirt armature", "ARMATURE")
    armature.object = rig
    obj["construction"] = "shaped cloth around preserved seated support"
    obj["status"] = "unreviewed construction candidate; no stage pass"
    return obj


def run():
    assert bpy.app.background, "Root's approved writer is pinned background Blender"
    assert Path(bpy.data.filepath).resolve() == SOURCE.resolve()
    assert sha(SOURCE) == SOURCE_HASH
    assert bpy.context.scene.name == SCENE_NAME
    assert bpy.context.view_layer.name == "ViewLayer"
    assert bpy.context.scene.frame_current == 1
    assert bpy.context.mode == "OBJECT"
    assert not bpy.data.is_dirty, "Save or resolve pending live edits first"
    assert not CANDIDATE.exists(), CANDIDATE
    assert not RECEIPT.exists(), RECEIPT
    assert COLLECTION not in bpy.data.collections
    assert all(name not in bpy.data.objects for name in NEW_NAMES)
    scene = bpy.context.scene
    assert all(name in scene.objects for name in TARGETS + SUPPORTS)
    assert all(not scene.objects[name].hide_render for name in TARGETS)
    assert all(scene.objects[name].type == "MESH" for name in PANELS + RUFFLES)
    assert all(scene.objects[name].type == "CURVE" for name in STITCHES)
    material_red = bpy.data.materials["Dress red cloth.004"]
    material_white = bpy.data.materials["Dress warm white cloth.002"]
    descriptor = os.open(LOCK, os.O_CREAT | os.O_EXCL | os.O_WRONLY, 0o600)
    os.close(descriptor)
    collection = None
    saved = False
    before = {obj.name: object_record(obj) for obj in scene.objects
              if obj.name not in TARGETS}
    target_visibility = {name: (scene.objects[name].hide_render,
                                scene.objects[name].hide_get())
                         for name in TARGETS}
    before_bounds = {
        "skirt": union_bounds([scene.objects[name] for name in PANELS]),
        "hem": union_bounds([scene.objects[name] for name in RUFFLES]),
        "bodice": bounds(scene.objects["Garment42 compact bodice"]),
        "seat": bounds(scene.objects[SUPPORTS[0]]),
    }
    try:
        support = support_tree([scene.objects[name] for name in SUPPORTS])
        skins = skin_source(PANELS), skin_source(RUFFLES)
        collection = bpy.data.collections.new(COLLECTION)
        scene.collection.children.link(collection)
        main = make_cloth(NEW_NAMES[0],
                          lambda t, a: skirt_point(t, a, support),
                          49, material_red, skins[0], collection)
        hem = make_cloth(NEW_NAMES[1],
                         lambda u, a: ruffle_point(u, a, support),
                         13, material_white, skins[1], collection)
        bpy.context.view_layer.update()
        after_bounds = {"skirt": bounds(main), "hem": bounds(hem)}
        red_span = after_bounds["skirt"][0][1] - after_bounds["skirt"][0][0]
        hem_span = after_bounds["hem"][0][1] - after_bounds["hem"][0][0]
        depth = after_bounds["hem"][1][1] - after_bounds["hem"][1][0]
        assert abs(red_span - RED_WIDTH) < 0.0012, red_span
        assert abs(hem_span - HEM_WIDTH) < 0.0012, hem_span
        assert 0.1153 <= depth <= 0.1340, depth
        assert after_bounds["hem"][2][0] >= 0.0, after_bounds["hem"]
        assert after_bounds["skirt"][2][1] <= 0.0640
        assert before == {name: object_record(scene.objects[name])
                          for name in before}, "A non-target changed"
        for name in TARGETS:
            scene.objects[name].hide_render = True
            scene.objects[name].hide_set(True)
        bpy.context.view_layer.update()
        assert before == {name: object_record(scene.objects[name])
                          for name in before}, "Target hiding changed a control"
        assert sha(SOURCE) == SOURCE_HASH
        assert not CANDIDATE.exists() and not RECEIPT.exists()
        if DEFER_SAVE:
            return {
                "status": "skirt constructed in memory; save deferred to root",
                "source": str(SOURCE), "source_sha256": SOURCE_HASH,
                "candidate_written": False,
                "hidden_objects": list(TARGETS),
                "new_objects": list(NEW_NAMES),
                "controls_before_skirt": before,
                "world_evaluated_controls_preserved_by_skirt": True,
                "bounds_before_m": before_bounds,
                "bounds_after_skirt_m": after_bounds,
                "head_width_m": WH,
            }
        # Do not alter cameras, lighting, materials, frame, or scene settings.
        bpy.ops.wm.save_as_mainfile(filepath=str(CANDIDATE), check_existing=True)
        saved = True
        assert Path(bpy.data.filepath).resolve() == CANDIDATE.resolve()
        assert sha(SOURCE) == SOURCE_HASH
        receipt = {
            "status": "unreviewed skirt construction candidate; no stage pass",
            "source": str(SOURCE), "source_sha256": SOURCE_HASH,
            "candidate": str(CANDIDATE), "candidate_sha256": sha(CANDIDATE),
            "script_sha256": sha(Path(__file__).resolve()),
            "runtime": {"version": bpy.app.version_string,
                        "build_hash": bpy.app.build_hash.decode(),
                        "background": bpy.app.background},
            "hidden_objects": list(TARGETS), "new_objects": list(NEW_NAMES),
            "preserved_control_count": len(before), "controls": before,
            "world_evaluated_controls_preserved": True,
            "source_preserved": True,
            "bounds_before_m": before_bounds, "bounds_after_m": after_bounds,
            "head_width_m": WH,
            "measured_spans_m": {"red_width": red_span,
                                 "hem_width": hem_span, "hem_depth": depth},
            "construction": {"drape_rows": DRAPE, "gathers": GATHERS,
                             "ruffle_folds": RUFFLE_FOLDS,
                             "support_objects": list(SUPPORTS),
                             "support_clearance_m": 0.0018},
            "limitations": [
                "No render or visual acceptance is claimed by this script.",
                "Existing body height and foot proportions remain unchanged.",
                "Nearest source skirt weights are retained; motion needs review.",
                "Clean reopen, seam contact, and evaluated collision checks remain.",
            ],
        }
        with RECEIPT.open("x") as handle:
            handle.write(json.dumps(receipt, indent=2) + "\n")
        return {key: value for key, value in receipt.items()
                if key not in {"controls", "construction"}}
    except Exception:
        if not saved:
            for name, (hide_render, hide_set) in target_visibility.items():
                scene.objects[name].hide_render = hide_render
                scene.objects[name].hide_set(hide_set)
            if collection is not None:
                for obj in list(collection.objects):
                    mesh = obj.data
                    bpy.data.objects.remove(obj, do_unlink=True)
                    if mesh.users == 0:
                        bpy.data.meshes.remove(mesh)
                bpy.data.collections.remove(collection)
            bpy.context.view_layer.update()
        raise
    finally:
        LOCK.unlink()


result = run()
```

## body_macro_022b.py

```python
"""Pinned, non-destructive seated proportion study; sole writer is root."""

import hashlib
import json
import math
from pathlib import Path

import bpy
from mathutils import Vector

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_021_candidate.blend"
EXPECTED = "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22"
OUTPUT = ROOT / "body_022b_candidate.blend"
RECEIPT = ROOT / "body_022b_writer_receipt.json"
DONOR = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/donors/a202/model.blend"
DONOR_HASH = "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5"
WH = 0.1165
HEAD_SHIFT = 0.0221
BODY_SHIFT = 0.0166
FEET = {"Left black stuffed foot pod", "Right black stuffed foot pod"}
BOW_NAMES = {
    "A42 flattened gathered center tie",
    *(f"A42 {side} {part}" for side in ("Left", "Right") for part in (
        "constructed bow loop", "independent draped bow tail",
        "narrow gathered loop ruffle", "narrow gathered tail ruffle",
        "root fold 1", "root fold 2", "white zigzag applique")),
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def body_z(z):
    t = max(0.0, min(1.0, (z - 0.025) / 0.075))
    return z - BODY_SHIFT * t * t * (3.0 - 2.0 * t)


def lattice(name, function):
    data = bpy.data.lattices.new(name)
    data.points_u = data.points_v = 2
    data.points_w = 64
    data.interpolation_type_u = data.interpolation_type_v = "KEY_LINEAR"
    data.interpolation_type_w = "KEY_LINEAR"
    obj = bpy.data.objects.new(name, data)
    bpy.context.scene.collection.objects.link(obj)
    base = [[min(p.co_deform[i] for p in data.points),
             max(p.co_deform[i] for p in data.points)] for i in range(3)]
    extents = (0.6, 0.6, 0.3)
    obj.scale = tuple(extents[i] / (base[i][1] - base[i][0]) for i in range(3))
    obj.location = (0.0, 0.0, 0.15)
    bpy.context.view_layer.update()
    inverse = obj.matrix_world.inverted()
    for point in data.points:
        world = obj.matrix_world @ point.co_deform
        world.z = function(world.z)
        point.co_deform = inverse @ world
    obj.hide_render = True
    obj.hide_set(True)
    return obj


def attach(obj, cage):
    modifier = obj.modifiers.new("022 non-destructive rest proportion", "LATTICE")
    modifier.object = cage
    modifier.show_viewport = modifier.show_render = True
    return modifier


def points(obj):
    depsgraph = bpy.context.evaluated_depsgraph_get()
    evaluated = obj.evaluated_get(depsgraph)
    mesh = evaluated.to_mesh()
    try:
        return [evaluated.matrix_world @ vertex.co for vertex in mesh.vertices]
    finally:
        evaluated.to_mesh_clear()


def fixture_check():
    """Prove mesh/curve deformation and exact modifier-removal recovery."""
    bpy.ops.wm.read_factory_settings(use_empty=True)
    mesh = bpy.data.meshes.new("Disposable022Mesh")
    mesh.from_pydata([(-0.01, 0, 0.07), (0.01, 0, 0.07), (0, 0, 0.09)], [], [(0, 1, 2)])
    obj = bpy.data.objects.new("Disposable022Mesh", mesh)
    bpy.context.scene.collection.objects.link(obj)
    curve = bpy.data.curves.new("Disposable022Curve", "CURVE")
    curve.dimensions = "3D"
    curve.bevel_depth = 0.001
    spline = curve.splines.new("POLY")
    spline.points.add(1)
    spline.points[0].co = (-0.01, 0, 0.08, 1)
    spline.points[1].co = (0.01, 0, 0.08, 1)
    line = bpy.data.objects.new("Disposable022Curve", curve)
    bpy.context.scene.collection.objects.link(line)
    cage = lattice("Disposable022Lattice", body_z)
    checks = []
    for target in (obj, line):
        before = points(target)
        mod = attach(target, cage)
        bpy.context.view_layer.update()
        after = points(target)
        assert len(before) == len(after) and before
        error = max(abs(a.z - body_z(b.z)) for b, a in zip(before, after))
        xy = max(abs(a[i] - b[i]) for b, a in zip(before, after) for i in (0, 1))
        movement = min(b.z - a.z for b, a in zip(before, after))
        assert error < 0.0001 and xy < 1e-7 and movement > 0.008, (error, xy, movement)
        target.modifiers.remove(mod)
        bpy.context.view_layer.update()
        restored = points(target)
        assert all(tuple(a) == tuple(b) for a, b in zip(before, restored))
        checks.append({"type": target.type, "vertices": len(before),
                       "max_field_error_m": error, "max_xy_error_m": xy,
                       "minimum_displacement_m": movement,
                       "removal_recovers_exact_geometry": True})
    return checks


assert bpy.app.background and bpy.app.version[:2] == (5, 2)
assert not OUTPUT.exists() and not RECEIPT.exists()
assert sha(SOURCE) == EXPECTED and sha(DONOR) == DONOR_HASH
fixture = json.loads((ROOT / "body_022_writer_receipt.json").read_text())["fixture_checks"]
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
scene = bpy.context.scene
assert scene.name == "Attempt41_Manual_Head_Maquette"
assert scene.frame_current == 1 and bpy.context.mode == "OBJECT"
assert not bpy.data.is_dirty
skirt_path = ROOT / "body_skirt_022b.py"
namespace = {"__file__": str(skirt_path), "BODY022_DEFER_SAVE": True}
exec(compile(skirt_path.read_text(), str(skirt_path), "exec"), namespace)
skirt = namespace["result"]
record = namespace["object_record"]
protected = {name: record(scene.objects[name]) for name in FEET | {"Review floor", "ReimuFumoRig"}}
visible = [obj for obj in scene.objects
           if obj.type in {"MESH", "CURVE"} and obj.visible_get() and not obj.hide_render]
assert len(BOW_NAMES) == 15
assert all(scene.objects[name] in visible for name in BOW_NAMES)
roles = {"head": [], "bow": [], "body": []}
original = {}
for obj in visible:
    if obj.name in protected:
        continue
    xyz = points(obj)
    if not xyz:
        continue
    z_max = max(p.z for p in xyz)
    role = "body"
    if obj.name in BOW_NAMES:
        role = "bow"
    elif z_max > 0.125 or obj.name == "A44 tiny neutral embroidered mouth dash":
        assert min(p.z for p in xyz) > 0.07, obj.name
        role = "head"
    roles[role].append(obj.name)
    original[obj.name] = {"points": xyz, "modifiers": [(m.name, m.type) for m in obj.modifiers]}
assert len(roles["head"]) > 5 and len(roles["bow"]) > 4 and len(roles["body"]) > 10, roles
crown = max(p.z for name in roles["head"] for p in original[name]["points"])
bow_top = max(p.z for name in roles["bow"] for p in original[name]["points"])
pivot = 0.2205
factor = (crown + 0.220 * WH - pivot) / (bow_top - pivot)
assert 0.25 < factor < 0.9, (crown, bow_top, factor)


def bow_z(z):
    return min(z, pivot) + max(0.0, z - pivot) * factor - HEAD_SHIFT


functions = {"body": body_z, "head": lambda z: z - HEAD_SHIFT, "bow": bow_z}
cages = {role: lattice("022 " + role + " proportion cage", fn)
         for role, fn in functions.items()}
for role, names in roles.items():
    for name in names:
        attach(scene.objects[name], cages[role])
bpy.context.view_layer.update()
effects = {}
for role, names in roles.items():
    deviations, xy_errors = [], []
    for name in names:
        obj = scene.objects[name]
        assert [(m.name, m.type) for m in obj.modifiers][:-1] == original[name]["modifiers"]
        before, after = original[name]["points"], points(obj)
        assert len(before) == len(after)
        deviations.extend(abs(a.z - functions[role](b.z)) for b, a in zip(before, after))
        xy_errors.extend(abs(a[i] - b[i]) for b, a in zip(before, after) for i in (0, 1))
    # Linear lattice sampling introduces a small known interpolation error.
    assert max(deviations) < 0.0007 and max(xy_errors) < 1e-6, (role, max(deviations), max(xy_errors))
    effects[role] = {"count": len(names), "max_field_error_m": max(deviations), "max_xy_error_m": max(xy_errors)}
assert protected == {name: record(scene.objects[name]) for name in protected}
assert sha(SOURCE) == EXPECTED and sha(DONOR) == DONOR_HASH
scene["candidate_status"] = "Unreviewed body022 proportion study; no stage or rig acceptance"
bpy.context.preferences.filepaths.save_version = 0
bpy.ops.wm.save_as_mainfile(filepath=str(OUTPUT), check_existing=True)
assert sha(SOURCE) == EXPECTED and sha(DONOR) == DONOR_HASH
receipt = {
    "candidate": str(OUTPUT), "candidate_sha256": sha(OUTPUT),
    "source": str(SOURCE), "source_sha256": EXPECTED,
    "protected_donor_sha256": DONOR_HASH,
    "version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(),
    "background": bpy.app.background, "writer": "root-astra-pinned-native",
    "script_sha256": sha(Path(__file__)), "skirt_script_sha256": sha(skirt_path),
    "fixture_checks": fixture, "roles": roles, "effects": effects,
    "foot_floor_rig_controls_unchanged": True,
    "existing_geometry_and_modifiers_preserved": True,
    "skirt": {key: value for key, value in skirt.items() if key != "controls_before_skirt"},
    "head_shift_m": HEAD_SHIFT, "body_shift_m": BODY_SHIFT,
    "original_crown_m": crown, "bow_pivot_m": pivot,
    "original_bow_top_m": bow_top, "upper_bow_factor": factor,
    "limitations": ["Rest-frame correction; animated deformation is unverified.",
                    "No reference-fidelity or final-stage acceptance claimed."]}
with RECEIPT.open("x") as handle:
    handle.write(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({key: value for key, value in receipt.items() if key not in {"roles", "skirt"}}), flush=True)
```

## render_body_022b.py

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
output = ROOT / "body_022b_review"
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
    ("candidate", "body_022b_candidate.blend", list(contract["fixed_views"])),
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

## body_022b_writer_receipt.json

```json
{
  "candidate": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/body_022b_candidate.blend",
  "candidate_sha256": "96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014",
  "source": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_021_candidate.blend",
  "source_sha256": "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22",
  "protected_donor_sha256": "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5",
  "version": "5.2.1 LTS",
  "build_hash": "9e2066aef7ef",
  "background": true,
  "writer": "root-astra-pinned-native",
  "script_sha256": "d495d2a0836b2db015789be1d82272d6d08f06c13017d9a8217dc9e60a4ad861",
  "skirt_script_sha256": "38c4549c6629722d7b7e1202c7c10064792fc19467b29021e12cb2d4a0bda90f",
  "fixture_checks": [
    {
      "type": "MESH",
      "vertices": 3,
      "max_field_error_m": 1.2616574972562633e-05,
      "max_xy_error_m": 0.0,
      "minimum_displacement_m": 0.010749075561761856,
      "removal_recovers_exact_geometry": true
    },
    {
      "type": "CURVE",
      "vertices": 24,
      "max_field_error_m": 2.0995192973272125e-05,
      "max_xy_error_m": 0.0,
      "minimum_displacement_m": 0.01340349018573761,
      "removal_recovers_exact_geometry": true
    }
  ],
  "roles": {
    "head": [
      "A44 tiny neutral embroidered mouth dash",
      "A45 left flush composite eye applique",
      "A45 left drooped half-lid stitch",
      "A45 left fine upper expression stitch",
      "A45 right flush composite eye applique",
      "A45 right drooped half-lid stitch",
      "A45 right fine upper expression stitch",
      "A45 left tapered flexible cheek lock",
      "A45 right tapered flexible cheek lock",
      "Head_Gusseted_Cushion_020b",
      "Hair_Continuous_Traced_Fringe_020b",
      "Rear_Center_Cloth_021",
      "Rear_Left_Cloth_021",
      "Rear_Right_Cloth_021"
    ],
    "bow": [
      "A42 flattened gathered center tie",
      "A42 Left constructed bow loop",
      "A42 Left independent draped bow tail",
      "A42 Left narrow gathered loop ruffle",
      "A42 Left narrow gathered tail ruffle",
      "A42 Left root fold 1",
      "A42 Left root fold 2",
      "A42 Left white zigzag applique",
      "A42 Right constructed bow loop",
      "A42 Right independent draped bow tail",
      "A42 Right narrow gathered loop ruffle",
      "A42 Right narrow gathered tail ruffle",
      "A42 Right root fold 1",
      "A42 Right root fold 2",
      "A42 Right white zigzag applique"
    ],
    "body": [
      "Left soft collar",
      "Right soft collar",
      "Folded yellow cravat",
      "Garment42 compact bodice",
      "Garment42 compact internal seat pad",
      "Left short hidden leg root",
      "Right short hidden leg root",
      "Sleeve44P L asymmetrically seated stuffed arm insert",
      "Sleeve44P L front folded cuff edge",
      "Sleeve44P L front padded fabric panel",
      "Sleeve44P L front red running stitch 1",
      "Sleeve44P L front red running stitch 2",
      "Sleeve44P L front red running stitch 3",
      "Sleeve44P L front shoulder pleat 1",
      "Sleeve44P L front shoulder pleat 2",
      "Sleeve44P L front shoulder pleat 3",
      "Sleeve44P L gathered shoulder root seam",
      "Sleeve44P L lower pinched cuff fold",
      "Sleeve44P L lower sewn panel join",
      "Sleeve44P L rear folded cuff edge",
      "Sleeve44P L rear padded fabric panel",
      "Sleeve44P L rear red running stitch 1",
      "Sleeve44P L rear red running stitch 2",
      "Sleeve44P L rear red running stitch 3",
      "Sleeve44P L rear shoulder pleat 1",
      "Sleeve44P L rear shoulder pleat 2",
      "Sleeve44P L rear shoulder pleat 3",
      "Sleeve44P L upper pinched cuff fold",
      "Sleeve44P L upper sewn panel join",
      "Sleeve44P R asymmetrically seated stuffed arm insert",
      "Sleeve44P R front folded cuff edge",
      "Sleeve44P R front padded fabric panel",
      "Sleeve44P R front red running stitch 1",
      "Sleeve44P R front red running stitch 2",
      "Sleeve44P R front red running stitch 3",
      "Sleeve44P R front shoulder pleat 1",
      "Sleeve44P R front shoulder pleat 2",
      "Sleeve44P R front shoulder pleat 3",
      "Sleeve44P R gathered shoulder root seam",
      "Sleeve44P R lower pinched cuff fold",
      "Sleeve44P R lower sewn panel join",
      "Sleeve44P R rear folded cuff edge",
      "Sleeve44P R rear padded fabric panel",
      "Sleeve44P R rear red running stitch 1",
      "Sleeve44P R rear red running stitch 2",
      "Sleeve44P R rear red running stitch 3",
      "Sleeve44P R rear shoulder pleat 1",
      "Sleeve44P R rear shoulder pleat 2",
      "Sleeve44P R rear shoulder pleat 3",
      "Sleeve44P R upper pinched cuff fold",
      "Sleeve44P R upper sewn panel join",
      "Skirt022 joined gathered panels",
      "Skirt022 soft hem ruffle"
    ]
  },
  "effects": {
    "head": {
      "count": 14,
      "max_field_error_m": 4.627704619508677e-08,
      "max_xy_error_m": 0.0
    },
    "bow": {
      "count": 15,
      "max_field_error_m": 0.0003888799981287827,
      "max_xy_error_m": 0.0
    },
    "body": {
      "count": 53,
      "max_field_error_m": 4.7021535339863973e-05,
      "max_xy_error_m": 0.0
    }
  },
  "foot_floor_rig_controls_unchanged": true,
  "existing_geometry_and_modifiers_preserved": true,
  "skirt": {
    "status": "skirt constructed in memory; save deferred to root",
    "source": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_021_candidate.blend",
    "source_sha256": "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22",
    "candidate_written": false,
    "hidden_objects": [
      "Garment42 front shallow lap panel",
      "Garment42 rear pooled dress panel",
      "Garment42 side gusset left",
      "Garment42 side gusset right",
      "Garment42 front gathered ruffle",
      "Garment42 rear pooled ruffle",
      "Garment42 left side gathered ruffle",
      "Garment42 right side gathered ruffle",
      "Garment42 hem stitch front 00",
      "Garment42 hem stitch front 01",
      "Garment42 hem stitch front 02",
      "Garment42 hem stitch front 03",
      "Garment42 hem stitch front 04",
      "Garment42 hem stitch front 05",
      "Garment42 hem stitch front 06",
      "Garment42 hem stitch front 07",
      "Garment42 hem stitch front 08",
      "Garment42 hem stitch rear 00",
      "Garment42 hem stitch rear 01",
      "Garment42 hem stitch rear 02",
      "Garment42 hem stitch rear 03",
      "Garment42 hem stitch rear 04",
      "Garment42 hem stitch rear 05",
      "Garment42 hem stitch rear 06",
      "Garment42 hem stitch rear 07",
      "Garment42 hem stitch rear 08"
    ],
    "new_objects": [
      "Skirt022 joined gathered panels",
      "Skirt022 soft hem ruffle"
    ],
    "world_evaluated_controls_preserved_by_skirt": true,
    "bounds_before_m": {
      "skirt": [
        [
          -0.06414402276277542,
          0.06409301608800888
        ],
        [
          -0.05816975235939026,
          0.05012364313006401
        ],
        [
          0.007951358333230019,
          0.06764774024486542
        ]
      ],
      "hem": [
        [
          -0.06993462890386581,
          0.07034382224082947
        ],
        [
          -0.060614459216594696,
          0.05119286850094795
        ],
        [
          -0.0005807435372844338,
          0.02880212478339672
        ]
      ],
      "bodice": [
        [
          -0.04100000113248825,
          0.04100000113248825
        ],
        [
          -0.035999998450279236,
          0.027000000700354576
        ],
        [
          0.05999999865889549,
          0.10000000149011612
        ]
      ],
      "seat": [
        [
          -0.05250000208616257,
          0.052499983459711075
        ],
        [
          -0.04400000348687172,
          0.039999984204769135
        ],
        [
          0.010499999858438969,
          0.03550000116229057
        ]
      ]
    },
    "bounds_after_skirt_m": {
      "skirt": [
        [
          -0.05344659462571144,
          0.053423769772052765
        ],
        [
          -0.06106378883123398,
          0.057039447128772736
        ],
        [
          0.0061026280745863914,
          0.06233207881450653
        ]
      ],
      "hem": [
        [
          -0.06206637620925903,
          0.06209399923682213
        ],
        [
          -0.06557871401309967,
          0.061618294566869736
        ],
        [
          0.0024814445059746504,
          0.028632303699851036
        ]
      ]
    },
    "head_width_m": 0.1165
  },
  "head_shift_m": 0.0221,
  "body_shift_m": 0.0166,
  "original_crown_m": 0.2153976857662201,
  "bow_pivot_m": 0.2205,
  "original_bow_top_m": 0.25398585200309753,
  "upper_bow_factor": 0.6130256373444286,
  "limitations": [
    "Rest-frame correction; animated deformation is unverified.",
    "No reference-fidelity or final-stage acceptance claimed."
  ]
}
```

## body_022_writer_receipt.json

```json
{
  "candidate": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/body_022_candidate.blend",
  "candidate_sha256": "0ff8cab54dfa5be530a12524f86e7238c13f90920ea3fd1864a27a6e6cdeef92",
  "source": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_021_candidate.blend",
  "source_sha256": "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22",
  "protected_donor_sha256": "a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5",
  "version": "5.2.1 LTS",
  "build_hash": "9e2066aef7ef",
  "background": true,
  "writer": "root-astra-pinned-native",
  "script_sha256": "556d5978badbefae15eb0be82c2999fbe25e7d2e78f9fe1416efd496b886c242",
  "skirt_script_sha256": "4ea9500da90737748362b0bade8113ce63bcfa5d23545852dd273f64558aac69",
  "fixture_checks": [
    {
      "type": "MESH",
      "vertices": 3,
      "max_field_error_m": 1.2616574972562633e-05,
      "max_xy_error_m": 0.0,
      "minimum_displacement_m": 0.010749075561761856,
      "removal_recovers_exact_geometry": true
    },
    {
      "type": "CURVE",
      "vertices": 24,
      "max_field_error_m": 2.0995192973272125e-05,
      "max_xy_error_m": 0.0,
      "minimum_displacement_m": 0.01340349018573761,
      "removal_recovers_exact_geometry": true
    }
  ],
  "roles": {
    "head": [
      "A45 left flush composite eye applique",
      "A45 left drooped half-lid stitch",
      "A45 left fine upper expression stitch",
      "A45 right flush composite eye applique",
      "A45 right drooped half-lid stitch",
      "A45 right fine upper expression stitch",
      "A45 left tapered flexible cheek lock",
      "A45 right tapered flexible cheek lock",
      "Head_Gusseted_Cushion_020b",
      "Hair_Continuous_Traced_Fringe_020b",
      "Rear_Center_Cloth_021",
      "Rear_Left_Cloth_021",
      "Rear_Right_Cloth_021"
    ],
    "bow": [
      "A42 flattened gathered center tie",
      "A42 Left constructed bow loop",
      "A42 Left independent draped bow tail",
      "A42 Left narrow gathered loop ruffle",
      "A42 Left narrow gathered tail ruffle",
      "A42 Left root fold 1",
      "A42 Left root fold 2",
      "A42 Left white zigzag applique",
      "A42 Right constructed bow loop",
      "A42 Right independent draped bow tail",
      "A42 Right narrow gathered loop ruffle",
      "A42 Right narrow gathered tail ruffle",
      "A42 Right root fold 1",
      "A42 Right root fold 2",
      "A42 Right white zigzag applique"
    ],
    "body": [
      "Left soft collar",
      "Right soft collar",
      "Folded yellow cravat",
      "Garment42 compact bodice",
      "Garment42 compact internal seat pad",
      "Left short hidden leg root",
      "Right short hidden leg root",
      "A44 tiny neutral embroidered mouth dash",
      "Sleeve44P L asymmetrically seated stuffed arm insert",
      "Sleeve44P L front folded cuff edge",
      "Sleeve44P L front padded fabric panel",
      "Sleeve44P L front red running stitch 1",
      "Sleeve44P L front red running stitch 2",
      "Sleeve44P L front red running stitch 3",
      "Sleeve44P L front shoulder pleat 1",
      "Sleeve44P L front shoulder pleat 2",
      "Sleeve44P L front shoulder pleat 3",
      "Sleeve44P L gathered shoulder root seam",
      "Sleeve44P L lower pinched cuff fold",
      "Sleeve44P L lower sewn panel join",
      "Sleeve44P L rear folded cuff edge",
      "Sleeve44P L rear padded fabric panel",
      "Sleeve44P L rear red running stitch 1",
      "Sleeve44P L rear red running stitch 2",
      "Sleeve44P L rear red running stitch 3",
      "Sleeve44P L rear shoulder pleat 1",
      "Sleeve44P L rear shoulder pleat 2",
      "Sleeve44P L rear shoulder pleat 3",
      "Sleeve44P L upper pinched cuff fold",
      "Sleeve44P L upper sewn panel join",
      "Sleeve44P R asymmetrically seated stuffed arm insert",
      "Sleeve44P R front folded cuff edge",
      "Sleeve44P R front padded fabric panel",
      "Sleeve44P R front red running stitch 1",
      "Sleeve44P R front red running stitch 2",
      "Sleeve44P R front red running stitch 3",
      "Sleeve44P R front shoulder pleat 1",
      "Sleeve44P R front shoulder pleat 2",
      "Sleeve44P R front shoulder pleat 3",
      "Sleeve44P R gathered shoulder root seam",
      "Sleeve44P R lower pinched cuff fold",
      "Sleeve44P R lower sewn panel join",
      "Sleeve44P R rear folded cuff edge",
      "Sleeve44P R rear padded fabric panel",
      "Sleeve44P R rear red running stitch 1",
      "Sleeve44P R rear red running stitch 2",
      "Sleeve44P R rear red running stitch 3",
      "Sleeve44P R rear shoulder pleat 1",
      "Sleeve44P R rear shoulder pleat 2",
      "Sleeve44P R rear shoulder pleat 3",
      "Sleeve44P R upper pinched cuff fold",
      "Sleeve44P R upper sewn panel join",
      "Skirt022 joined gathered panels",
      "Skirt022 soft hem ruffle"
    ]
  },
  "effects": {
    "head": {
      "count": 13,
      "max_field_error_m": 4.627704619508677e-08,
      "max_xy_error_m": 0.0
    },
    "bow": {
      "count": 15,
      "max_field_error_m": 0.0003888799981287827,
      "max_xy_error_m": 0.0
    },
    "body": {
      "count": 54,
      "max_field_error_m": 4.7021535339863973e-05,
      "max_xy_error_m": 0.0
    }
  },
  "foot_floor_rig_controls_unchanged": true,
  "existing_geometry_and_modifiers_preserved": true,
  "skirt": {
    "status": "skirt constructed in memory; save deferred to root",
    "source": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_021_candidate.blend",
    "source_sha256": "9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22",
    "candidate_written": false,
    "hidden_objects": [
      "Garment42 front shallow lap panel",
      "Garment42 rear pooled dress panel",
      "Garment42 side gusset left",
      "Garment42 side gusset right",
      "Garment42 front gathered ruffle",
      "Garment42 rear pooled ruffle",
      "Garment42 left side gathered ruffle",
      "Garment42 right side gathered ruffle",
      "Garment42 hem stitch front 00",
      "Garment42 hem stitch front 01",
      "Garment42 hem stitch front 02",
      "Garment42 hem stitch front 03",
      "Garment42 hem stitch front 04",
      "Garment42 hem stitch front 05",
      "Garment42 hem stitch front 06",
      "Garment42 hem stitch front 07",
      "Garment42 hem stitch front 08",
      "Garment42 hem stitch rear 00",
      "Garment42 hem stitch rear 01",
      "Garment42 hem stitch rear 02",
      "Garment42 hem stitch rear 03",
      "Garment42 hem stitch rear 04",
      "Garment42 hem stitch rear 05",
      "Garment42 hem stitch rear 06",
      "Garment42 hem stitch rear 07",
      "Garment42 hem stitch rear 08"
    ],
    "new_objects": [
      "Skirt022 joined gathered panels",
      "Skirt022 soft hem ruffle"
    ],
    "world_evaluated_controls_preserved_by_skirt": true,
    "bounds_before_m": {
      "skirt": [
        [
          -0.06414402276277542,
          0.06409301608800888
        ],
        [
          -0.05816975235939026,
          0.05012364313006401
        ],
        [
          0.007951358333230019,
          0.06764774024486542
        ]
      ],
      "hem": [
        [
          -0.06993462890386581,
          0.07034382224082947
        ],
        [
          -0.060614459216594696,
          0.05119286850094795
        ],
        [
          -0.0005807435372844338,
          0.02880212478339672
        ]
      ],
      "bodice": [
        [
          -0.04100000113248825,
          0.04100000113248825
        ],
        [
          -0.035999998450279236,
          0.027000000700354576
        ],
        [
          0.05999999865889549,
          0.10000000149011612
        ]
      ],
      "seat": [
        [
          -0.05250000208616257,
          0.052499983459711075
        ],
        [
          -0.04400000348687172,
          0.039999984204769135
        ],
        [
          0.010499999858438969,
          0.03550000116229057
        ]
      ]
    },
    "bounds_after_skirt_m": {
      "skirt": [
        [
          -0.053396739065647125,
          0.053371306508779526
        ],
        [
          -0.06103269383311272,
          0.057019203901290894
        ],
        [
          0.006100062746554613,
          0.06234709545969963
        ]
      ],
      "hem": [
        [
          -0.06206637620925903,
          0.06209399923682213
        ],
        [
          -0.06557871401309967,
          0.061618294566869736
        ],
        [
          0.0024814445059746504,
          0.028632303699851036
        ]
      ]
    },
    "head_width_m": 0.1165
  },
  "head_shift_m": 0.0221,
  "body_shift_m": 0.0166,
  "original_crown_m": 0.2153976857662201,
  "bow_pivot_m": 0.2205,
  "original_bow_top_m": 0.25398585200309753,
  "upper_bow_factor": 0.6130256373444286,
  "limitations": [
    "Rest-frame correction; animated deformation is unverified.",
    "No reference-fidelity or final-stage acceptance claimed."
  ]
}
```
