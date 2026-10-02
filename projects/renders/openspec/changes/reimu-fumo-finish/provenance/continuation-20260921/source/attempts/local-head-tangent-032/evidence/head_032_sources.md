# Head032 reproducible source evidence

Candidate head_032_candidate.blend SHA256
6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8.
Retained source030b d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6.
Run writer then renderer with pinned Blender5.2.1 LTS background, factory
startup, disabled autoexec, four threads and python-exit-code2. Fresh outputs
are required. Regenerated Blender bytes need new hashes and visual review.
The helper also verifies the frozen030b writer receipt and exact head inputs.
No new modifier, shader, shell or rig architecture is added.

## head_032_draft.py

SHA256 f234f504cd481057d5018fb8c7e1d6dae6d5fea060830a51e050d2a14b10d318.

```python
"""One local tangent-continuity draft. No open/save or top-level mutation."""
import bisect
import hashlib
import json
import math
import struct
from pathlib import Path

import bpy
from mathutils import Vector
from mathutils.bvhtree import BVHTree

TARGETS = ('Head028 sewn cushion', 'Hair028 crown and back hood',
           'Hair028 traced padded fringe')
BASELINE_SHA256 = 'd69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6'
RECEIPT_SHA256 = '3da6d2f75bfaaa432e306082a8587f4ca62c8d739e24ab6ceecd362937bd0028'
# Root may explicitly supply an exact later retained source path/hash, but
# only unchanged guarded head inputs permit such integration.
AUTHORIZED_ACTIVE_SOURCE = None
RX = .05816962569952011
RZ = .05732078477740287
CZ = .13590002432465553
RHO0 = .86


def _smooth(t):
    t = min(1., max(0., t))
    return t*t*(3.-2.*t)


def _mesh(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    try:
        return {'points': [ev.matrix_world @ v.co for v in mesh.vertices],
                'faces': [tuple(f.vertices) for f in mesh.polygons],
                'materials': [f.material_index for f in mesh.polygons]}
    finally:
        ev.to_mesh_clear()


def _record(obj):
    ev = obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh = ev.to_mesh()
    digest = hashlib.sha256()
    try:
        digest.update(struct.pack('<3I',len(mesh.vertices),len(mesh.edges),len(mesh.polygons)))
        for v in mesh.vertices:
            p = ev.matrix_world @ v.co
            digest.update(struct.pack('<3f',*p))
        for f in mesh.polygons:
            digest.update(struct.pack('<2I',len(f.vertices),f.material_index))
            for i in f.vertices:
                digest.update(struct.pack('<I',i))
        return {'geometry':digest.hexdigest(),
                'materials':[m.name if m else None for m in mesh.materials],
                'visibility':[obj.hide_render,obj.hide_viewport,obj.hide_get()],
                'parent':obj.parent.name if obj.parent else None,
                'modifiers':[[m.name,m.type] for m in obj.modifiers]}
    finally:
        ev.to_mesh_clear()


def _topology(obj):
    return {'faces':[(tuple(f.vertices),f.material_index) for f in obj.data.polygons],
            'weights':[[(g.group,g.weight) for g in v.groups] for v in obj.data.vertices],
            'groups':[g.name for g in obj.vertex_groups],
            'modifiers':[(m.name,m.type,getattr(getattr(m,'object',None),'name',None)) for m in obj.modifiers]}


def _head_map(obj):
    rig = bpy.context.scene.objects['ReimuFumoRig']
    bone = rig.pose.bones['Head']
    assert [m.type for m in obj.modifiers] in (['ARMATURE'],['SOLIDIFY','ARMATURE'])
    assert obj.modifiers[-1].object == rig
    group = obj.vertex_groups.get('Head')
    assert group is not None
    assert all(len(v.groups)==1 and v.groups[0].group==group.index and v.groups[0].weight==1.
               for v in obj.data.vertices)
    return (rig.matrix_world @ bone.matrix @ bone.bone.matrix_local.inverted()
            @ rig.matrix_world.inverted() @ obj.matrix_world)


def _rho_theta(p):
    x,z = p.x/RX,(p.z-CZ)/RZ
    return math.hypot(x,z),math.atan2(z,x)%(2*math.pi)


def _stats(values):
    values=sorted(values)
    return {'count':len(values),'minimum':values[0],'median':values[len(values)//2],
            'p95':values[int(.95*(len(values)-1))],'maximum':values[-1]}


def _normal(points,face):
    return (points[face[1]]-points[face[0]]).cross(points[face[2]]-points[face[0]]).normalized()


def _edges(faces):
    result={}
    for fi,f in enumerate(faces):
        for a,b in zip(f,f[1:]+f[:1]):
            result.setdefault(tuple(sorted((a,b))),[]).append(fi)
    return result


def _seam_jumps(data,edges):
    output={'above125':[],'above145':[],'fade125to145':[]}
    pts,faces=data['points'],data['faces']
    for edge,incident in edges.items():
        if len(incident)!=2 or min(incident)>=4896 or max(incident)<9792:
            continue
        z=sum(pts[i].z for i in edge)/2
        if z<=.125:
            continue
        normals=[_normal(pts,faces[i]) for i in incident]
        angle=math.degrees(math.acos(max(-1.,min(1.,normals[0].dot(normals[1])))))
        output['above125'].append(angle)
        output['above145' if z>.145 else 'fade125to145'].append(angle)
    return {key:_stats(values) for key,values in output.items()}


class _Field:
    def __init__(self,data,edges):
        self.data=data
        self.tree=BVHTree.FromPolygons(data['points'],data['faces'])
        self.front={i for f in data['faces'][:4896] for i in f}
        self.seam={i for e,inc in edges.items() if len(inc)==2 and min(inc)<4896
                   and max(inc)>=9792 for i in e}
        neighbors={i:set() for i in self.seam}
        for a,b in edges:
            if a in neighbors: neighbors[a].add(b)
            if b in neighbors: neighbors[b].add(a)
        rows=[]
        for i in self.seam:
            p=data['points'][i]
            rho,theta=_rho_theta(p)
            outer=[j for j in neighbors[i] if j not in self.front]
            if p.z>.12:
                assert len(outer)==1,('Ambiguous gusset tangent',i,outer)
                q=data['points'][outer[0]]
                qrho,unused=_rho_theta(q)
                assert qrho>rho
                slope=(q.y-p.y)/(qrho-rho)
            else:
                slope=None
            rows.append((theta,rho,p.y,slope,i))
        self.rows=sorted(rows)
        self.angles=[row[0] for row in self.rows]
        self.cache={}
        self.max_transfer_distance=0.

    def _profile(self,theta):
        key=round(theta,7)
        if key in self.cache: return self.cache[key]
        right=bisect.bisect_right(self.angles,theta)%len(self.rows)
        left=(right-1)%len(self.rows)
        a,b=self.rows[left],self.rows[right]
        ta,tb=a[0],b[0]
        if tb<=ta: tb+=2*math.pi
        query=theta if theta>=ta else theta+2*math.pi
        t=(query-ta)/(tb-ta)
        assert a[3] is not None and b[3] is not None
        rho1=a[1]*(1-t)+b[1]*t
        y1=a[2]*(1-t)+b[2]*t
        m1=a[3]*(1-t)+b[3]*t
        x,z=RX*RHO0*math.cos(theta),CZ+RZ*RHO0*math.sin(theta)
        hit=self.tree.ray_cast(Vector((x,-.2,z)),Vector((0,1,0)),.4)
        assert hit[0] is not None and hit[2]<4896,('Inner field receiver',theta,hit[2])
        normal=hit[1]
        y0=hit[0].y
        m0=-(normal.x*RX*math.cos(theta)+normal.z*RZ*math.sin(theta))/normal.y
        length=rho1-RHO0
        aa=2*y0-2*y1+length*(m0+m1)
        bb=-3*y0+3*y1-length*(2*m0+m1)
        cc=length*m0
        probes=[0.,1.]
        if abs(aa)>1e-12:
            stationary=-bb/(3*aa)
            if 0<stationary<1: probes.append(stationary)
        minimum=min((3*aa*u*u+2*bb*u+cc)/length for u in probes)
        assert y1>y0 and minimum>=-1e-6,('Non-monotone cubic; stop without reshaping',theta,y0,y1,m0,m1,minimum)
        profile={'theta':theta,'rho1':rho1,'y0':y0,'y1':y1,'m0':m0,'m1':m1,
                 'a':aa,'b':bb,'c':cc,'minimum_dY_drho':minimum}
        self.cache[key]=profile
        return profile

    def delta(self,p,receiver_y=None):
        if p.z<=.125: return 0.
        rho,theta=_rho_theta(p)
        if rho<=RHO0 or rho>=.98801: return 0.
        profile=self._profile(theta)
        if rho>=profile['rho1']-1e-7: return 0.
        u=(rho-RHO0)/(profile['rho1']-RHO0)
        target=((profile['a']*u+profile['b'])*u+profile['c'])*u+profile['y0']
        if receiver_y is None:
            hit=self.tree.ray_cast(Vector((p.x,-.2,p.z)),Vector((0,1,0)),.4)
            assert hit[0] is not None and hit[2]<4896,('Fringe receiver outside front',list(p),hit[2])
            receiver_y=hit[0].y
            self.max_transfer_distance=max(self.max_transfer_distance,abs(p.y-receiver_y))
        value=_smooth((p.z-.125)/.020)*(target-receiver_y)
        assert math.isfinite(value) and abs(value)<.008,('Nonlocal fairing displacement',value,list(p))
        return value


def build_head_032():
    root=Path(__file__).resolve().parent
    baseline=root/'bow_030b_candidate.blend'
    receipt_path=root/'bow_030b_writer_receipt.json'
    assert hashlib.sha256(baseline.read_bytes()).hexdigest()==BASELINE_SHA256
    assert hashlib.sha256(receipt_path.read_bytes()).hexdigest()==RECEIPT_SHA256
    reference=json.loads(receipt_path.read_text())
    active=Path(bpy.data.filepath).resolve()
    if AUTHORIZED_ACTIVE_SOURCE is None:
        assert active==baseline.resolve()
        active_hash=BASELINE_SHA256
    else:
        assert active==Path(AUTHORIZED_ACTIVE_SOURCE['path']).resolve()
        active_hash=AUTHORIZED_ACTIVE_SOURCE['sha256']
    assert hashlib.sha256(active.read_bytes()).hexdigest()==active_hash
    assert bpy.context.scene.frame_current==1 and bpy.context.mode=='OBJECT'
    scene=bpy.context.scene
    guards={name:record for name,record in reference['controls'].items()
            if name.startswith(('Head028 ','Hair028 '))}
    assert len(guards)==15 and all(name in guards for name in TARGETS)
    assert guards=={name:_record(scene.objects[name]) for name in guards},'Head inputs differ from frozen030b'
    controls={obj.name:_record(obj) for obj in scene.objects if obj.type in {'MESH','CURVE'}
              and obj.visible_get() and not obj.hide_render and obj.name not in TARGETS}
    targets=[scene.objects[name] for name in TARGETS]
    assert all(obj.type=='MESH' and obj.data.users==1 for obj in targets)
    before={obj.name:_mesh(obj) for obj in targets}
    topologies={obj.name:_topology(obj) for obj in targets}
    maps={obj.name:_head_map(obj) for obj in targets}
    base_before={obj.name:[v.co.copy() for v in obj.data.vertices] for obj in targets}
    rig=scene.objects['ReimuFumoRig']
    pose={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    core=before[TARGETS[0]]
    assert len(core['faces'])==11040
    assert len(targets[0].data.vertices)==len(core['points'])
    edges=_edges(core['faces'])
    field=_Field(core,edges)
    core_delta=[field.delta(p,p.y) if i in field.front and i not in field.seam else 0.
                for i,p in enumerate(core['points'])]
    # Recover the exact hood subset/vertex order from the retained027c/028
    # recipe, using the still-preserved donor only for its material-region tag.
    donor=scene.objects['Head_Gusseted_Cushion_020b']
    assert len(donor.data.polygons)==2496
    assert donor.modifiers[0].type=='SUBSURF' and donor.modifiers[0].levels==1
    # One Catmull-Clark level emits one child quad per source face corner,
    # with the source face's material tag. Hidden-object dependency-graph
    # evaluation is not needed; the exact hood topology check below verifies
    # this retained-source ordering before any target coordinate is edited.
    donor_materials=[face.material_index for face in donor.data.polygons for unused in face.vertices]
    assert len(donor_materials)==9792
    back=[i for i in range(4896,9792) if max(core['points'][j].z for j in core['faces'][i])>.089]
    selected=back+[i for i in range(4896) if donor_materials[i]==0]+list(range(9792,11040))
    selected_faces=[core['faces'][i] for i in selected]
    used=sorted({i for face in selected_faces for i in face})
    mapping={i:j for j,i in enumerate(used)}
    expected=[tuple(mapping[i] for i in face) for face in selected_faces]
    assert expected==[tuple(f.vertices) for f in targets[1].data.polygons],'Hood correspondence changed'
    assert len(used)==len(targets[1].data.vertices)
    hood_delta=[core_delta[i] for i in used]
    for j,i in enumerate(used):
        base_world=maps[TARGETS[1]]@targets[1].data.vertices[j].co
        assert (base_world-core['points'][i]).length<=.000551,'Hood baseline stand-off changed'
    fringe_delta=[field.delta(maps[TARGETS[2]]@v.co) for v in targets[2].data.vertices]
    deltas={TARGETS[0]:core_delta,TARGETS[1]:hood_delta,TARGETS[2]:fringe_delta}
    profile_before={str(round(p['theta'],7)):p for p in field.cache.values()}
    # All coefficients and monotonicity checks above precede the first edit.
    for obj in targets:
        inverse=maps[obj.name].to_3x3().inverted()
        for vertex,value in zip(obj.data.vertices,deltas[obj.name]):
            if value!=0.:
                vertex.co+=inverse@Vector((0.,value,0.))
        obj.data.update()
    bpy.context.view_layer.update()
    after={obj.name:_mesh(obj) for obj in targets}
    assert controls=={name:_record(scene.objects[name]) for name in controls},'Protected geometry changed'
    assert pose=={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    assert topologies=={obj.name:_topology(obj) for obj in targets},'Topology, weights or modifier stack changed'
    ca,cb=core['points'],after[TARGETS[0]]['points']
    assert all(a.x==b.x and a.z==b.z for a,b in zip(ca,cb)),'Core X/Z changed'
    zero_core={i for i,v in enumerate(core_delta) if v==0.}
    assert all(ca[i]==cb[i] for i in zero_core),'Protected zero-field core vertex moved'
    protected_faces=[i for i,f in enumerate(core['faces']) if i>=4896 or max(ca[j].z for j in f)<=.125]
    assert all(all(ca[j]==cb[j] for j in core['faces'][i]) for i in protected_faces)
    movement={}
    for obj in targets:
        old,new=before[obj.name]['points'],after[obj.name]['points']
        assert len(old)==len(new) and before[obj.name]['faces']==after[obj.name]['faces']
        delta=[b-a for a,b in zip(old,new)]
        low=[i for i,p in enumerate(old) if p.z<=.125]
        movement[obj.name]={
            'changed_base_vertices':sum(v!=0. for v in deltas[obj.name]),
            'base_depth_delta_range_m':[min(deltas[obj.name]),max(deltas[obj.name])],
            'evaluated_max_abs_delta_xyz_m':[max(abs(d[k]) for d in delta) for k in range(3)],
            'evaluated_below125_max_displacement_m':max((delta[i].length for i in low),default=0.),
            'base_zero_field_vertices_exact':all(v.co==base_before[obj.name][i] for i,v in enumerate(obj.data.vertices) if deltas[obj.name][i]==0.),
        }
    assert hashlib.sha256(active.read_bytes()).hexdigest()==active_hash
    return {'baseline_sha256':BASELINE_SHA256,'active_source':str(active),'active_source_sha256':active_hash,
            'head_input_guards':guards,'targets':list(TARGETS),'protected_control_count':len(controls),
            'protected_controls_unchanged':True,'rig_pose_weights_modifiers_topology_unchanged':True,
            'core_xz_and_zero_field_interfaces_exact':True,'protected_core_face_count':len(protected_faces),
            'core_seam_jump_before_deg':_seam_jumps(core,edges),
            'core_seam_jump_after_deg':_seam_jumps(after[TARGETS[0]],edges),
            'cubic_profile_count':len(profile_before),
            'minimum_cubic_dY_drho':min(p['minimum_dY_drho'] for p in profile_before.values()),
            'maximum_fringe_receiver_standoff_m':field.max_transfer_distance,
            'movement':movement,'profile_coefficients':profile_before,
            'construction':'Existing base vertices, worldY correction inverse-mapped through current Head pose; no new modifiers',
            'limitations':['Live Solidify normal offsets may move derived X/Z and points below125mm',
                           'Monotonic cubic witnesses are not a global collision or animation pass',
                           'No save and no visual acceptance']}
```

## head_032_writer.py

SHA256 6daa08f765589a15e72778bce97cefcd4c64a3ada8978a9dd48cb0a5e086da01.

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
OUTPUT = ROOT / "head_032_candidate.blend"
RECEIPT = ROOT / "head_032_writer_receipt.json"
EXPECTED_TARGETS = frozenset({
    "Head028 sewn cushion",
    "Hair028 crown and back hood",
    "Hair028 traced padded fringe",
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
helper = ROOT / "head_032_draft.py"
scope = {"__file__": str(helper)}
exec(compile(helper.read_text(), str(helper), "exec"), scope)
assert frozenset(scope["TARGETS"]) == EXPECTED_TARGETS
result = scope["build_head_032"]()
bpy.context.view_layer.update()
assert controls == {name: record(scene.objects[name]) for name in controls}
assert pose == {bone.name: [list(row) for row in bone.matrix] for bone in rig.pose.bones}
assert look == appearance(scene), "Material nodes, lighting or color settings changed"
for name in EXPECTED_TARGETS:
    obj = scene.objects[name]
    assert obj.data.materials and all(
        m is not None and not m.is_evaluated
        and bpy.data.materials.get(m.name) == m
        for m in obj.data.materials), ("Invalid persistent material binding", name)
assert sha(SOURCE) == SOURCE_HASH
scene["candidate_status"] = "Head032 unreviewed tangent continuity study; no stage pass"
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

## render_head_032.py

SHA256 c87884bdb5c983ca5052b689c9d1dd5dad93e6f1854d2bd19dfe58dd628be1c6.

```python
"""Frozen032 existing-node views, settings bound to actual026 baseline."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "head_032_candidate.blend"
OUT = ROOT / "head_032_eevee_review"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
BASELINE = ROOT / "bow_030b_eevee_review/render_receipt.json"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


frozen = sha(SOURCE)
writer = json.loads((ROOT / "head_032_writer_receipt.json").read_text())
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
    data = bpy.data.cameras.new("Frozen032_" + view)
    camera = bpy.data.objects.new("Frozen032_" + view, data)
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


