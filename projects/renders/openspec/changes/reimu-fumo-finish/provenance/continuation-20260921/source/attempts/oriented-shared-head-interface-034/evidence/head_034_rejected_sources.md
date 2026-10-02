# Rejected034 shared-interface construction sources

Both states were unsaved. Retained model remains Foot033 SHA256
98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48.
Diagnostics used exact unchanged Head032 for causal comparison.
These scripts require the hash-bound head_032_draft.py utility and input
files named in them. Cut JSON SHA256
b52af4da1cc8feedc9162b4e1c2c23b3d7912168e7a2b7a6c25ed42a074b2156,
guard JSON SHA256
eeb645a243c3f610a517ee16e3462027788395e32baf360fe29e88146a65a526.
Those raw task-local dependency files are not reproduced in this packet.
No claim that this bundle alone supplies every reproduction artifact.
The canonical032 source packet owns the utility source.

## head_034_shared_draft.py

```python
"""Captured two-skin shared interface. No open/save or import-time mutation."""
import bisect
import hashlib
import json
import math
from pathlib import Path

import bpy
from mathutils import Euler, Vector

TARGETS=('Hair028 traced padded fringe','Hair028 crown and back hood')
NEW_NAME='Hair034 shared crown fringe shell'
SOURCE_SHA='6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8'
CUT_SHA='b52af4da1cc8feedc9162b4e1c2c23b3d7912168e7a2b7a6c25ed42a074b2156'
GUARD_SHA='eeb645a243c3f610a517ee16e3462027788395e32baf360fe29e88146a65a526'
UTILITY_SHA='f234f504cd481057d5018fb8c7e1d6dae6d5fea060830a51e050d2a14b10d318'
AUTHORIZED_ACTIVE_SOURCE=None
REGIONS={0:'retained_fringe_outer',1:'shared_outer_bridge',2:'retained_hood_outer',
         3:'retained_front_cover',4:'inner_skin_or_allowance'}


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _capture(obj):
    ev=obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh=ev.to_mesh()
    try:
        mesh.calc_loop_triangles()
        normalmap=ev.matrix_world.to_3x3().inverted().transposed()
        return {'points':[ev.matrix_world@v.co for v in mesh.vertices],
                'normals':[(normalmap@v.normal).normalized() for v in mesh.vertices],
                'polygons':[tuple(p.vertices) for p in mesh.polygons],
                'triangles':[{'vertices':tuple(t.vertices),'polygon':t.polygon_index,
                              'material':mesh.polygons[t.polygon_index].material_index,
                              'smooth':mesh.polygons[t.polygon_index].use_smooth,
                              'uv':{layer.name:[tuple(layer.data[i].uv) for i in t.loops]
                                    for layer in mesh.uv_layers}}
                             for t in mesh.loop_triangles],
                'materials':[m.name if m else None for m in mesh.materials],
                'uv_layers':[l.name for l in mesh.uv_layers]}
    finally:
        ev.to_mesh_clear()


def _ordered(edges,points):
    graph={}
    for a,b in edges:
        graph.setdefault(a,[]).append(b)
        graph.setdefault(b,[]).append(a)
    ends=[i for i,adj in graph.items() if len(adj)==1]
    assert len(ends)==2 and all(len(adj)<=2 for adj in graph.values())
    path=[min(ends,key=lambda i:points[i].x)]
    previous=None
    while True:
        nxt=[i for i in graph[path[-1]] if i!=previous]
        if not nxt:
            break
        previous=path[-1]
        path.append(nxt[0])
    assert len(path)==len(graph)
    lengths=[0.]
    for a,b in zip(path,path[1:]):
        lengths.append(lengths[-1]+(points[b]-points[a]).length)
    return path,[d/lengths[-1] for d in lengths]


def _curve_sample(path,parameters,t,points):
    nearest=min(range(len(parameters)),key=lambda i:abs(parameters[i]-t))
    if abs(parameters[nearest]-t)<1e-10:
        i=path[nearest]
        return points[i].copy(),i,None
    j=min(len(path)-2,max(0,bisect.bisect_right(parameters,t)-1))
    u=(t-parameters[j])/(parameters[j+1]-parameters[j])
    a,b=path[j],path[j+1]
    return points[a].lerp(points[b],u),None,(a,b,u)


def _stats(values):
    values=sorted(values)
    return {'count':len(values),'min':values[0],'median':values[len(values)//2],'max':values[-1]} if values else None


class _Builder:
    def __init__(self,points):
        self.points=list(points)
        self.faces=[]
        self.info=[]
        self.edges={}

    def point(self,p):
        self.points.append(p.copy())
        return len(self.points)-1

    def add(self,face,info,anchor=False):
        face=tuple(face)
        assert len(set(face))==len(face)
        normal=(self.points[face[1]]-self.points[face[0]]).cross(self.points[face[2]]-self.points[face[0]])
        assert normal.length>1e-14,('Degenerate new face',info,face)
        if anchor:
            decisions=[]
            for a,b in zip(face,face[1:]+face[:1]):
                inc=self.edges.get(tuple(sorted((a,b))),[])
                if inc:
                    assert len(inc)==1,('Already closed seam edge',a,b,info)
                    decisions.append(inc[0][1]==(a,b))
            assert decisions and all(v==decisions[0] for v in decisions),('Inconsistent seam winding',info,decisions)
            if decisions[0]:
                face=tuple(reversed(face))
                info=dict(info)
                info['uv']={key:list(reversed(value)) for key,value in info.get('uv',{}).items()}
        fi=len(self.faces)
        self.faces.append(face)
        self.info.append(info)
        for a,b in zip(face,face[1:]+face[:1]):
            self.edges.setdefault(tuple(sorted((a,b))),[]).append((fi,(a,b)))
        return fi


def _ownership(scene,contract,mesh_regions):
    rows={}
    dg=bpy.context.evaluated_depsgraph_get()
    scale=contract['camera']['ortho_scale_m']
    for view,ys,xrange in [('side',(146,160,180,200,220),(180,235)),
                           ('three_quarter',(150,160,175,190,210),(225,307)),
                           ('three_quarter_mirror',(150,160,175,190,210),(100,280))]:
        spec=contract['fixed_views'][view]
        rotation=Euler(spec['rotation_euler_rad']).to_matrix()
        direction=rotation@Vector((0,0,-1))
        rows[view]=[]
        for y in ys:
            segments=[]
            for x in range(*xrange):
                origin=Vector(spec['location_m'])+rotation@Vector((((x+.5)/512-.5)*scale,
                                                                 (.5-(y+.5)/512)*scale,0))
                hit=scene.ray_cast(dg,origin,direction,distance=2.)
                name=hit[4].name if hit[0] else None
                region=mesh_regions[hit[3]] if name==NEW_NAME else None
                key=[name,REGIONS.get(region)]
                if not segments or segments[-1]['owner']!=key:
                    segments.append({'owner':key,'x':[x,x]})
                else:
                    segments[-1]['x'][1]=x
            rows[view].append({'y':y,'segments':segments})
    return rows


def build_head_034_shared():
    root=Path(__file__).resolve().parent
    baseline=root/'head_032_candidate.blend'
    cut_path=root/'head_034_shared_topology_probe.json'
    guard_path=root/'head_034_dryrun.json'
    utility=root/'head_032_draft.py'
    assert _sha(baseline)==SOURCE_SHA and _sha(cut_path)==CUT_SHA
    assert _sha(guard_path)==GUARD_SHA and _sha(utility)==UTILITY_SHA
    active=Path(bpy.data.filepath).resolve()
    active_hash=SOURCE_SHA
    if AUTHORIZED_ACTIVE_SOURCE is None:
        assert active==baseline.resolve()
    else:
        assert active==Path(AUTHORIZED_ACTIVE_SOURCE['path']).resolve()
        active_hash=AUTHORIZED_ACTIVE_SOURCE['sha256']
    assert _sha(active)==active_hash and bpy.context.scene.frame_current==1
    assert bpy.context.mode=='OBJECT' and NEW_NAME not in bpy.data.objects
    ns={'__file__':str(utility),'__name__':'shared_interface_utilities'}
    exec(compile(utility.read_text(),str(utility),'exec'),ns)
    scene=bpy.context.scene
    guards=json.loads(guard_path.read_text())['result']['head_input_guards']
    assert len(guards)==15 and guards=={name:ns['_record'](scene.objects[name]) for name in guards}
    controls={o.name:ns['_record'](o) for o in scene.objects if o.type in {'MESH','CURVE'}
              and o.visible_get() and not o.hide_render and o.name not in TARGETS}
    rig=scene.objects['ReimuFumoRig']
    pose={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    cut=json.loads(cut_path.read_text())
    snapshots={name:_capture(scene.objects[name]) for name in TARGETS}
    assert snapshots[TARGETS[0]]['materials']==snapshots[TARGETS[1]]['materials']
    mats=snapshots[TARGETS[0]]['materials']
    uv_names=sorted({key for s in snapshots.values() for key in s['uv_layers']})
    base_counts={name:len(scene.objects[name].data.vertices) for name in TARGETS}
    face_counts={name:len(scene.objects[name].data.polygons) for name in TARGETS}
    offsets={TARGETS[0]:0,TARGETS[1]:len(snapshots[TARGETS[0]]['points'])}
    all_points=snapshots[TARGETS[0]]['points']+snapshots[TARGETS[1]]['points']
    builder=_Builder(all_points)
    for name in TARGETS:
        obj=scene.objects[name]
        mapping=ns['_head_map'](obj)
        s=snapshots[name]
        assert len(s['points'])==2*base_counts[name]
        assert max((mapping@v.co-s['points'][i]).length for i,v in enumerate(obj.data.vertices))<2e-8
        for fi,p in enumerate(obj.data.polygons):
            assert set(v%base_counts[name] for v in s['polygons'][fi+face_counts[name]])==set(p.vertices)
    fp=snapshots[TARGETS[0]]['points']
    hp=snapshots[TARGETS[1]]['points']
    fpath,ft=_ordered(cut['fringe']['cut_edges'],fp)
    hpath,ht=_ordered(cut['hood']['gusset_receiver_cut_edges'],hp)
    assert len(fpath)==337 and len(hpath)==97
    knots=sorted({round(t,12) for t in ft+ht})
    splits={}
    boundary={}
    normals={}
    for name,path,ts in [(TARGETS[0],fpath,ft),(TARGETS[1],hpath,ht)]:
        snap=snapshots[name]
        for inner in (False,True):
            shift=base_counts[name] if inner else 0
            indices=[]
            nsamples=[]
            shifted=[i+shift for i in path]
            for t in knots:
                p,original,edge=_curve_sample(shifted,ts,t,snap['points'])
                if original is not None:
                    index=offsets[name]+original
                    normal=snap['normals'][original]
                else:
                    a,b,u=edge
                    index=builder.point(p)
                    ga,gb=offsets[name]+a,offsets[name]+b
                    fraction=u if ga<gb else 1.-u
                    splits.setdefault(tuple(sorted((ga,gb))),[]).append((fraction,index))
                    normal=snap['normals'][a].lerp(snap['normals'][b],u).normalized()
                indices.append(index)
                nsamples.append(normal.copy())
            boundary[(name,inner)]=indices
            normals[(name,inner)]=nsamples
    for edge,values in splits.items():
        values.sort()
    removed_by_name={TARGETS[0]:set(cut['fringe']['removed_faces']),TARGETS[1]:set(cut['hood']['removed_faces'])}
    root_edges={tuple(sorted(e)) for e in cut['fringe']['original_open_edges']}
    front_cover=set(cut['hood']['unchanged_front_cover_faces'])
    retained_map=[]
    removed_rims=0
    supported_triangles=0
    for name in TARGETS:
        snap=snapshots[name]
        nv,nf=base_counts[name],face_counts[name]
        removed=set(removed_by_name[name])|{i+nf for i in removed_by_name[name]}
        if name==TARGETS[0]:
            for fi,f in enumerate(snap['polygons']):
                if fi>=2*nf and tuple(sorted({i%nv for i in f})) in root_edges:
                    removed.add(fi)
                    removed_rims+=1
        for ti,tri in enumerate(snap['triangles']):
            if tri['polygon'] in removed:
                continue
            face=tuple(offsets[name]+i for i in tri['vertices'])
            region=(0 if name==TARGETS[0] else 3 if tri['polygon'] in front_cover else 2) if tri['polygon']<nf else 4
            info={'material':tri['material'],'smooth':tri['smooth'],'uv':tri['uv'],
                  'region':region,'source':[name,ti],'kind':'retained'}
            perimeter=[]
            uv={key:[] for key in uv_names}
            touched=False
            for j,(a,b) in enumerate(zip(face,face[1:]+face[:1])):
                perimeter.append(a)
                for key in uv_names:
                    uv[key].append(tri['uv'].get(key,[(0.,0.)]*3)[j])
                keyedge=tuple(sorted((a,b)))
                values=splits.get(keyedge,[])
                if a>b:
                    values=[(1.-u,i) for u,i in reversed(values)]
                for u,index in values:
                    touched=True
                    perimeter.append(index)
                    for key in uv_names:
                        source_uv=tri['uv'].get(key,[(0.,0.)]*3)
                        uv[key].append(tuple((1.-u)*source_uv[j][k]+u*source_uv[(j+1)%3][k] for k in range(2)))
            if not touched:
                fi=builder.add(face,info)
                retained_map.append((name,ti,fi))
            else:
                supported_triangles+=1
                center=builder.point(sum((builder.points[i] for i in face),Vector())/3.)
                center_uv={key:tuple(sum(v[k] for v in tri['uv'].get(key,[(0.,0.)]*3))/3. for k in range(2)) for key in uv_names}
                for j,a in enumerate(perimeter):
                    b=perimeter[(j+1)%len(perimeter)]
                    part=dict(info)
                    part['kind']='boundary_edge_subdivision'
                    part['uv']={key:[uv[key][j],uv[key][(j+1)%len(perimeter)],center_uv[key]] for key in uv_names}
                    builder.add((a,b,center),part)
    assert removed_rims==306
    # The two arcs have one shared longitudinal partition. Width rows are
    # an endpoint-tangent Hermite strip, not two coincident sheet envelopes.
    layers={}
    chord_lengths=[]
    for inner in (False,True):
        a=boundary[(TARGETS[0],inner)]
        b=boundary[(TARGETS[1],inner)]
        rows=[a]
        for level in range(1,8):
            t=level/8.
            row=[]
            for j,(ia,ib) in enumerate(zip(a,b)):
                pa,pb=builder.points[ia],builder.points[ib]
                chord=pb-pa
                if not inner and level==1:
                    chord_lengths.append(chord.length)
                na=normals[(TARGETS[0],inner)][j]
                nb=normals[(TARGETS[1],inner)][j]
                da=chord-na*chord.dot(na)
                db=chord-nb*chord.dot(nb)
                assert chord.length<.014 and min(da.length,db.length)>.00005,('Unusable endpoint tangent',j,chord.length)
                da=da.normalized()*chord.length
                db=db.normalized()*chord.length
                p=(pa*(2*t**3-3*t*t+1)+da*(t**3-2*t*t+t)
                   +pb*(-2*t**3+3*t*t)+db*(t**3-t*t))
                row.append(builder.point(p))
            rows.append(row)
        rows.append(b)
        layers[inner]=rows
        for level,(r0,r1) in enumerate(zip(rows,rows[1:])):
            for j in range(len(knots)-1):
                info={'material':0,'smooth':True,'uv':{key:[(knots[j],level/8.),(knots[j+1],level/8.),
                                                              (knots[j+1],(level+1)/8.),(knots[j],(level+1)/8.)] for key in uv_names},
                      'region':4 if inner else 1,'source':None,'kind':'inner_bridge' if inner else 'outer_bridge'}
                builder.add((r0[j],r0[j+1],r1[j+1],r1[j]),info,anchor=True)
    for a,b in cut['hood']['retained_front_cover_cut_edges']:
        outer_a,outer_b=offsets[TARGETS[1]]+a,offsets[TARGETS[1]]+b
        inner_a,inner_b=outer_a+base_counts[TARGETS[1]],outer_b+base_counts[TARGETS[1]]
        builder.add((outer_a,outer_b,inner_b,inner_a),{'material':0,'smooth':True,'uv':{},'region':4,
                    'source':None,'kind':'hidden_front_cover_allowance'},anchor=True)
    for j in (0,len(knots)-1):
        for level in range(8):
            face=(layers[False][level][j],layers[False][level+1][j],
                  layers[True][level+1][j],layers[True][level][j])
            builder.add(face,{'material':0,'smooth':True,'uv':{},'region':4,'source':None,'kind':'side_connector'},anchor=True)
    bad_edges={str(e):inc for e,inc in builder.edges.items() if len(inc)!=2 or inc[0][1]!=tuple(reversed(inc[1][1]))}
    assert not bad_edges,('Shell not closed consistently',len(bad_edges),list(bad_edges.items())[:8])
    thickness=[(builder.points[a]-builder.points[b]).length for ra,rb in zip(layers[False],layers[True]) for a,b in zip(ra,rb)]
    assert min(thickness)>.0002 and max(thickness)<.003,('Shell thickness range',min(thickness),max(thickness))
    # All array/topology checks precede creation or hiding of scene objects.
    used=sorted({i for f in builder.faces for i in f})
    remap={old:new for new,old in enumerate(used)}
    bone=rig.pose.bones['Head']
    deformation=rig.matrix_world@bone.matrix@bone.bone.matrix_local.inverted()@rig.matrix_world.inverted()
    local=[deformation.inverted()@builder.points[i] for i in used]
    mesh=bpy.data.meshes.new(NEW_NAME+' mesh')
    mesh.from_pydata(local,[],[tuple(remap[i] for i in f) for f in builder.faces])
    for name in mats:
        mesh.materials.append(bpy.data.materials[name])
    for p,info in zip(mesh.polygons,builder.info):
        p.material_index=info['material']
        p.use_smooth=info['smooth']
    for key in uv_names:
        layer=mesh.uv_layers.new(name=key)
        for p,info in zip(mesh.polygons,builder.info):
            values=info['uv'].get(key,[(0.,0.)]*len(p.vertices))
            for loop,value in zip(p.loop_indices,values):
                layer.data[loop].uv=value
    region_attr=mesh.attributes.new('head034_face_region','INT','FACE')
    for item,info in zip(region_attr.data,builder.info):
        item.value=info['region']
    obj=bpy.data.objects.new(NEW_NAME,mesh)
    scene.collection.objects.link(obj)
    group=obj.vertex_groups.new(name='Head')
    group.add(list(range(len(mesh.vertices))),1.,'REPLACE')
    modifier=obj.modifiers.new('034 live Head attachment','ARMATURE')
    modifier.object=rig
    obj['head034_face_regions']=json.dumps(REGIONS,sort_keys=True)
    for name in TARGETS:
        old=scene.objects[name]
        old.hide_render=True
        old.hide_viewport=True
        old.hide_set(True)
    bpy.context.view_layer.update()
    assert controls=={name:ns['_record'](scene.objects[name]) for name in controls}
    assert pose=={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    result=_capture(obj)
    # Existing off-support triangles keep their indexed coordinates and UV
    # corners; only boundary-adjacent triangles were split in their own plane.
    eval_mesh=ns['_mesh'](obj)
    errors=[]
    uv_exact=True
    for name,ti,fi in retained_map:
        tri=snapshots[name]['triangles'][ti]
        original=[snapshots[name]['points'][i] for i in tri['vertices']]
        face=eval_mesh['faces'][fi]
        errors.extend((eval_mesh['points'][j]-p).length for j,p in zip(face,original))
        uv_exact=uv_exact and builder.info[fi]['uv']==tri['uv']
    assert max(errors,default=0.)<2e-8 and uv_exact
    contract=json.loads((root.parents[2]/'projects/renders/assets/reimu_fumo/review_contract.json').read_text())
    ownership=_ownership(scene,contract,[info['region'] for info in builder.info])
    # Report source-region progression, not an object-name-only success.
    ownership_failures=[]
    for view,rows in ownership.items():
        for row in rows:
            seq=[s['owner'][1] for s in row['segments'] if s['owner'][0]==NEW_NAME]
            if any(v in {'inner_skin_or_allowance','retained_front_cover'} for v in seq):
                ownership_failures.append({'view':view,'y':row['y'],'reason':'Hidden region first-hit','sequence':seq})
            simple=[v for i,v in enumerate(seq) if i==0 or v!=seq[i-1]]
            if len(simple)!=len(set(simple)):
                ownership_failures.append({'view':view,'y':row['y'],'reason':'Repeated visible region','sequence':seq})
    assert _sha(active)==active_hash and _sha(baseline)==SOURCE_SHA
    return {'source_sha256':SOURCE_SHA,'active_source_sha256':active_hash,'head_input_guards':guards,
            'targets':list(TARGETS),'created_names':[NEW_NAME],'protected_control_count':len(controls),
            'protected_controls_exact':True,'rig_pose_exact':True,'live_binding':'One Armature, full unit Head weights',
            'thickness':'Captured existing two skins; no Solidify modifier or rig bake',
            'retained_off_support_triangles':len(retained_map),'off_support_coordinate_max_error_m':max(errors,default=0.),
            'retained_uv_corners_exact':uv_exact,'source_uv_layers':{name:s['uv_layers'] for name,s in snapshots.items()},
            'boundary_supported_source_triangles':supported_triangles,'common_arc_knots':len(knots),
            'removed_original_fringe_rims':removed_rims,'new_vertices':len(used),'new_faces':len(builder.faces),
            'nonmanifold_or_inconsistent_edges':0,'strip_chord_lengths_m':_stats(chord_lengths),
            'strip_paired_skin_separation_m':_stats(thickness),'face_region_counts':{REGIONS[k]:sum(i['region']==k for i in builder.info) for k in REGIONS},
            'ownership':ownership,'ownership_failures':ownership_failures,
            'suitability':'STOP before writer' if ownership_failures else 'Topology/ownership passed sampled checks; root visual and intersection review still required',
            'limitations':['No save or render, no visual acceptance.',
                           'Paired skin distances are not a complete thickness or self-intersection audit.',
                           'No animation acceptance; captured skins retain only the ordinary live Head attachment.',
                           'Smooth shading normal preservation at the new support ring requires visual review.']}
```

## head_034_shared_diagnostic.py

```python
"""Authorized unchanged reconstruction: new local hit/seam/intersection evidence."""
import hashlib
import json
import math
import traceback
from pathlib import Path

import bpy
from mathutils import Euler, Vector
from mathutils.bvhtree import BVHTree
from mathutils.geometry import intersect_ray_tri

ROOT=Path(__file__).resolve().parent
SOURCE=ROOT/'head_032_candidate.blend'
SOURCE_SHA='6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8'
HELPER=ROOT/'head_034_shared_draft.py'
HELPER_SHA='c9976b58dc310001b4478ddb399b6691bd801eb3a725b12358c6422000cb1c69'
OUT=ROOT/'head_034_shared_diagnostic.json'
assert not OUT.exists()
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest()==SOURCE_SHA
assert hashlib.sha256(HELPER.read_bytes()).hexdigest()==HELPER_SHA
bpy.ops.wm.open_mainfile(filepath=str(SOURCE),load_ui=False)
scope={'__file__':str(HELPER),'__name__':'unchanged_shared_reconstruction'}
exec(compile(HELPER.read_text(),str(HELPER),'exec'),scope)
construction=scope['build_head_034_shared']()
scene=bpy.context.scene
obj=scene.objects[scope['NEW_NAME']]
dg=bpy.context.evaluated_depsgraph_get()
ev=obj.evaluated_get(dg)
mesh=ev.to_mesh()
mesh.calc_loop_triangles()
points=[ev.matrix_world@v.co for v in mesh.vertices]
triangles=[tuple(t.vertices) for t in mesh.loop_triangles]
tri_polys=[t.polygon_index for t in mesh.loop_triangles]
regions=[mesh.attributes['head034_face_region'].data[p].value for p in tri_polys]
normalmap=ev.matrix_world.to_3x3().inverted().transposed()
normals=[(points[t[1]]-points[t[0]]).cross(points[t[2]]-points[t[0]]).normalized() for t in triangles]
corner_normals=[[(normalmap@mesh.corner_normals[j].vector).normalized() for j in t.loops]
                for t in mesh.loop_triangles]
tree=BVHTree.FromPolygons(points,triangles,all_triangles=True)
contract=json.loads((ROOT.parents[2]/'projects/renders/assets/reimu_fumo/review_contract.json').read_text())
spec=contract['fixed_views']['three_quarter_mirror']
rotation=Euler(spec['rotation_euler_rad']).to_matrix()
scale=contract['camera']['ortho_scale_m']
camera=Vector(spec['location_m'])
direction=rotation@Vector((0,0,-1))


def stats(values):
    values=sorted(values)
    return {'count':len(values),'min':values[0],'median':values[len(values)//2],
            'p95':values[int(.95*(len(values)-1))],'max':values[-1]} if values else None


def angle(a,b):
    return math.degrees(math.acos(max(-1.,min(1.,a.dot(b)))))


def project(p):
    q=rotation.transposed()@(p-camera)
    return [(q.x/scale+.5)*512-.5,(.5-q.y/scale)*512-.5]


def shading_normal(ti,p):
    a,b,c=[points[i] for i in triangles[ti]]
    v0,v1,v2=b-a,c-a,p-a
    d00,d01,d11=v0.dot(v0),v0.dot(v1),v1.dot(v1)
    d20,d21=v2.dot(v0),v2.dot(v1)
    denominator=d00*d11-d01*d01
    v=(d11*d20-d01*d21)/denominator
    w=(d00*d21-d01*d20)/denominator
    weights=(1.-v-w,v,w)
    return sum((n*t for n,t in zip(corner_normals[ti],weights)),Vector()).normalized()


pixel_hits=[]
for x in range(237,253):
    origin=camera+rotation@Vector((((x+.5)/512-.5)*scale,(.5-175.5/512)*scale,0))
    scene_hit=scene.ray_cast(dg,origin,direction,distance=2.)
    hits=[]
    start=origin.copy()
    for unused in range(12):
        p,n,ti,d=tree.ray_cast(start,direction,2.)
        if p is None:
            break
        hits.append({'triangle':ti,'polygon':tri_polys[ti],'vertices':list(triangles[ti]),
                     'region':scope['REGIONS'][regions[ti]],'point_m':list(p),
                     'camera_ray_depth_m':(p-origin).dot(direction),
                     'geometric_normal':list(n),'shading_normal':list(shading_normal(ti,p))})
        start=p+direction*.0000001
    pixel_hits.append({'pixel':[x,175], 'scene_first_object':scene_hit[4].name if scene_hit[0] else None,
                       'scene_first_polygon':scene_hit[3] if scene_hit[0] else None,'ordered_shell_hits':hits})

edges={}
for ti,t in enumerate(triangles):
    for a,b in zip(t,t[1:]+t[:1]):
        edges.setdefault(tuple(sorted((a,b))),[]).append(ti)
seams=[]
for edge,inc in edges.items():
    if len(inc)!=2 or {regions[i] for i in inc} not in ({0,1},{1,2}):
        continue
    pa,pb=[points[i] for i in edge]
    a,b=project(pa),project(pb)
    crossing=None
    if min(a[1],b[1])<=175<=max(a[1],b[1]) and abs(a[1]-b[1])>1e-10:
        t=(175-a[1])/(b[1]-a[1])
        crossing=a[0]+t*(b[0]-a[0])
    rows={'edge':list(edge),'triangles':inc,'polygons':[tri_polys[i] for i in inc],
          'regions':[scope['REGIONS'][regions[i]] for i in inc],
          'points_m':[list(pa),list(pb)],'projected_endpoints':[a,b],
          'normal_jump_deg':angle(normals[inc[0]],normals[inc[1]]),
          'normals':[list(normals[i]) for i in inc],'row175_crossing_x':crossing,
          'midpoint_z_m':(pa.z+pb.z)/2}
    seams.append(rows)
local_seams=[r for r in seams if max(p[0] for p in r['projected_endpoints'])>=235
             and min(p[0] for p in r['projected_endpoints'])<=254
             and max(p[1] for p in r['projected_endpoints'])>=171
             and min(p[1] for p in r['projected_endpoints'])<=179]
crossings=[r for r in seams if r['row175_crossing_x'] is not None and 235<=r['row175_crossing_x']<=254]


def segment_distance(p,a,b):
    v=b-a
    t=max(0.,min(1.,(p-a).dot(v)/v.length_squared)) if v.length_squared else 0.
    return (p-(a+v*t)).length


def coplanar_area(ta,tb,normal):
    axis=max(range(3),key=lambda k:abs(normal[k]))
    dims=[i for i in range(3) if i!=axis]
    a=[(points[i][dims[0]],points[i][dims[1]]) for i in ta]
    b=[(points[i][dims[0]],points[i][dims[1]]) for i in tb]
    def cross2(u,v):
        return u[0]*v[1]-u[1]*v[0]
    orientation=1. if cross2((b[1][0]-b[0][0],b[1][1]-b[0][1]),
                             (b[2][0]-b[0][0],b[2][1]-b[0][1]))>=0 else -1.
    output=a
    for q,r in zip(b,b[1:]+b[:1]):
        previous=output
        output=[]
        if not previous:
            break
        def side(p):
            return orientation*cross2((r[0]-q[0],r[1]-q[1]),(p[0]-q[0],p[1]-q[1]))
        for s,e in zip(previous,previous[1:]+previous[:1]):
            ds,de=side(s),side(e)
            if (ds>=0)!=(de>=0):
                t=ds/(ds-de)
                output.append((s[0]+t*(e[0]-s[0]),s[1]+t*(e[1]-s[1])))
            if de>=0:
                output.append(e)
    return abs(sum(p[0]*q[1]-p[1]*q[0] for p,q in zip(output,output[1:]+output[:1])))/2 if output else 0.


def contacts(ai,bi):
    ta,tb=triangles[ai],triangles[bi]
    shared=set(ta)&set(tb)
    na,nb=normals[ai],normals[bi]
    if abs(na.dot(nb))>.999999 and max(abs((points[i]-points[ta[0]]).dot(na)) for i in tb)<2e-8:
        area=coplanar_area(ta,tb,na)
        return {'kind':'coplanar_overlap','projected_area_m2':area} if area>1e-12 else None
    hits=[]
    for first,second in ((ta,tb),(tb,ta)):
        target=[points[i] for i in second]
        for ia,ib in zip(first,first[1:]+first[:1]):
            a,b=points[ia],points[ib]
            p=intersect_ray_tri(*target,b-a,a,True)
            if p is None:
                continue
            u=(p-a).dot(b-a)/(b-a).length_squared
            if u<-1e-7 or u>1+1e-7:
                continue
            if any((p-points[i]).length<2e-7 for i in shared):
                continue
            if len(shared)==2 and segment_distance(p,*[points[i] for i in shared])<2e-7:
                continue
            if all((p-q).length>2e-7 for q in hits):
                hits.append(p)
    return {'kind':'nonshared_segment_contact','points_m':[list(p) for p in hits]} if hits else None


outer=[i for i,r in enumerate(regions) if r==1]
outer_tree=BVHTree.FromPolygons(points,[triangles[i] for i in outer],all_triangles=True)
broad=outer_tree.overlap(tree)
pairs=set()
for local,other in broad:
    a=outer[local]
    if a!=other:
        pairs.add(tuple(sorted((a,other))))
intersection_evidence=[]
for a,b in sorted(pairs):
    contact=contacts(a,b)
    if contact:
        intersection_evidence.append({'triangles':[a,b],'polygons':[tri_polys[a],tri_polys[b]],
                                      'regions':[scope['REGIONS'][regions[a]],scope['REGIONS'][regions[b]]],
                                      'contact':contact})

report={'source_sha256':SOURCE_SHA,'unchanged_helper_sha256':HELPER_SHA,
        'purpose':'Distinguish projected welded boundary reversal from competing geometry; no shape changes.',
        'runtime':bpy.app.version_string,'build_hash':bpy.app.build_hash.decode(),
        'mirror_row175_hits':pixel_hits,'local_seam_edges':local_seams,'row175_seam_crossings':crossings,
        'upper_endpoint_seam_jump_statistics':{
            label:stats([r['normal_jump_deg'] for r in seams if r['midpoint_z_m']>.145
                         and region in r['regions']])
            for label,region in [('fringe_to_bridge','retained_fringe_outer'),('bridge_to_hood','retained_hood_outer')]},
        'all_seam_edges':seams,
        'outer_bridge_intersection_check':{'outer_bridge_triangles':len(outer),'broad_phase_unique_pairs':len(pairs),
                                          'reported_contact_pairs':len(intersection_evidence),'contacts':intersection_evidence,
                                          'shared_contact_tolerance_m':2e-7,'coplanar_projected_area_threshold_m2':1e-12},
        'limitations':['Intersection test is bounded to the actual outer bridge against this captured shell, not the whole scene.',
                       'Shared edge/vertex contacts within0.2 micrometer are allowed; coplanar area threshold1e-12m2.',
                       'No envelope-overshoot, whole-scene collision or visual audit.',
                       'Unchanged helper necessarily reconstructs its guards; their settled broad results are not re-audited here.']}
ev.to_mesh_clear()
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest()==SOURCE_SHA
assert hashlib.sha256(HELPER.read_bytes()).hexdigest()==HELPER_SHA
report['source_and_helper_preserved']=True
report['blend_saved']=False
report['rendered']=False
with OUT.open('x') as handle:
    handle.write(json.dumps(report,indent=2)+'\n')
print(json.dumps({'output':str(OUT),'crossing_count':len(crossings),
                  'seam_statistics':report['upper_endpoint_seam_jump_statistics'],
                  'intersection_pairs':len(intersection_evidence)}))
```

## head_034c_draft.py

```python
"""034c: only oriented co-normal endpoint derivatives change the shared strip."""
import bisect
import hashlib
import json
import math
from pathlib import Path

import bpy
from mathutils import Euler, Vector

TARGETS=('Hair028 traced padded fringe','Hair028 crown and back hood')
NEW_NAME='Hair034c shared crown fringe shell'
SOURCE_SHA='6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8'
CUT_SHA='b52af4da1cc8feedc9162b4e1c2c23b3d7912168e7a2b7a6c25ed42a074b2156'
GUARD_SHA='eeb645a243c3f610a517ee16e3462027788395e32baf360fe29e88146a65a526'
UTILITY_SHA='f234f504cd481057d5018fb8c7e1d6dae6d5fea060830a51e050d2a14b10d318'
AUTHORIZED_ACTIVE_SOURCE=None
PREFLIGHT={'phase':'not started','endpoint_directions':[],'derivative_vectors':[]}
REGIONS={0:'retained_fringe_outer',1:'shared_outer_bridge',2:'retained_hood_outer',
         3:'retained_front_cover',4:'inner_skin_or_allowance'}


def _sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _capture(obj):
    ev=obj.evaluated_get(bpy.context.evaluated_depsgraph_get())
    mesh=ev.to_mesh()
    try:
        mesh.calc_loop_triangles()
        normalmap=ev.matrix_world.to_3x3().inverted().transposed()
        return {'points':[ev.matrix_world@v.co for v in mesh.vertices],
                'normals':[(normalmap@v.normal).normalized() for v in mesh.vertices],
                'polygons':[tuple(p.vertices) for p in mesh.polygons],
                'triangles':[{'vertices':tuple(t.vertices),'polygon':t.polygon_index,
                              'material':mesh.polygons[t.polygon_index].material_index,
                              'smooth':mesh.polygons[t.polygon_index].use_smooth,
                              'uv':{layer.name:[tuple(layer.data[i].uv) for i in t.loops]
                                    for layer in mesh.uv_layers}}
                             for t in mesh.loop_triangles],
                'materials':[m.name if m else None for m in mesh.materials],
                'uv_layers':[l.name for l in mesh.uv_layers]}
    finally:
        ev.to_mesh_clear()


def _ordered(edges,points):
    graph={}
    for a,b in edges:
        graph.setdefault(a,[]).append(b)
        graph.setdefault(b,[]).append(a)
    ends=[i for i,adj in graph.items() if len(adj)==1]
    assert len(ends)==2 and all(len(adj)<=2 for adj in graph.values())
    path=[min(ends,key=lambda i:points[i].x)]
    previous=None
    while True:
        nxt=[i for i in graph[path[-1]] if i!=previous]
        if not nxt:
            break
        previous=path[-1]
        path.append(nxt[0])
    assert len(path)==len(graph)
    lengths=[0.]
    for a,b in zip(path,path[1:]):
        lengths.append(lengths[-1]+(points[b]-points[a]).length)
    return path,[d/lengths[-1] for d in lengths]


def _curve_sample(path,parameters,t,points):
    nearest=min(range(len(parameters)),key=lambda i:abs(parameters[i]-t))
    if abs(parameters[nearest]-t)<1e-10:
        i=path[nearest]
        return points[i].copy(),i,None
    j=min(len(path)-2,max(0,bisect.bisect_right(parameters,t)-1))
    u=(t-parameters[j])/(parameters[j+1]-parameters[j])
    a,b=path[j],path[j+1]
    return points[a].lerp(points[b],u),None,(a,b,u)


def _stats(values):
    values=sorted(values)
    return {'count':len(values),'min':values[0],'median':values[len(values)//2],'max':values[-1]} if values else None



def _co_normals(builder, path, side_sign, label):
    edge_data=[]
    for a,b in zip(path,path[1:]):
        incident=builder.edges[tuple(sorted((a,b)))]
        assert len(incident)==1, ('Not an exposed retained cut',label,a,b)
        fi,(start,end)=incident[0]
        face=builder.faces[fi]
        n=(builder.points[face[1]]-builder.points[face[0]]).cross(
            builder.points[face[2]]-builder.points[face[0]]).normalized()
        oriented=(builder.points[end]-builder.points[start]).normalized()
        desired=oriented.cross(n).normalized()*side_sign
        edge_data.append({'edge':[a,b],'retained_face':fi,'normal':n,'desired':desired})
    units=[]
    for j,i in enumerate(path):
        tangent=Vector()
        incident=[]
        if j:
            tangent+=(builder.points[i]-builder.points[path[j-1]]).normalized()
            incident.append(edge_data[j-1])
        if j+1<len(path):
            tangent+=(builder.points[path[j+1]]-builder.points[i]).normalized()
            incident.append(edge_data[j])
        tangent=tangent.normalized()
        raw=sum((e['desired'] for e in incident),Vector())
        raw-=tangent*raw.dot(tangent)
        unit=raw.normalized()
        signs=[unit.dot(e['desired']) for e in incident]
        entry={'boundary':label,'knot':j,'vertex':i,'point':list(builder.points[i]),
               'longitudinal_tangent':list(tangent),'unit_derivative':list(unit),
               'removed_longitudinal_component_dot':unit.dot(tangent),
               'incident_retained_edges':[{'edge':e['edge'],'face':e['retained_face'],
                                           'normal':list(e['normal']),
                                           'desired_conormal':list(e['desired'])} for e in incident],
               'desired_side_dot_products':signs,
               'valid':raw.length>1e-9 and min(signs)>1e-6
                       and abs(unit.dot(tangent))<1e-5}
        PREFLIGHT['endpoint_directions'].append(entry)
        units.append(unit)
    return units


def _strip_preflight(builder):
    jacobians=[]
    for fi,(face,info) in enumerate(zip(builder.faces,builder.info)):
        if info['kind'] not in {'outer_bridge','inner_bridge'}:
            continue
        p=[builder.points[i] for i in face]
        # A whole flipped cell must fail too, not only a bow-tie cell.
        reference=Vector(info['expected_normal'])
        row=[]
        for j in range(4):
            forward=p[(j+1)%4]-p[j]
            back=p[(j-1)%4]-p[j]
            cross=forward.cross(back)
            product=forward.length*back.length
            row.append({'corner':j,'area_twice_m2':cross.length,
                        'normalized_sine':cross.length/product if product else 0.,
                        'oriented_cosine':cross.normalized().dot(reference) if cross.length else 0.})
        jacobians.append({'face':fi,'kind':info['kind'],'vertices':list(face),
                          'points_m':[list(v) for v in p],'corners':row,
                          'valid':all(v['area_twice_m2']>1e-14 and v['normalized_sine']>1e-5
                                      and v['oriented_cosine']>1e-6 for v in row)})
    joins=[]
    for edge,inc in builder.edges.items():
        if len(inc)!=2:
            continue
        regions={builder.info[i]['region'] for i,unused in inc}
        if regions not in ({0,1},{1,2}):
            continue
        ns=[]
        for fi,unused in inc:
            f=builder.faces[fi]
            ns.append((builder.points[f[1]]-builder.points[f[0]]).cross(
                builder.points[f[2]]-builder.points[f[0]]).normalized())
        jump=math.degrees(math.acos(max(-1.,min(1.,ns[0].dot(ns[1])))))
        joins.append({'edge':list(edge),'faces':[i for i,unused in inc],
                      'regions':list(sorted(regions)),'first_triangle_normal_jump_deg':jump,
                      'normals':[list(n) for n in ns],
                      'points_m':[list(builder.points[i]) for i in edge]})
    PREFLIGHT['strip_jacobians']=jacobians
    PREFLIGHT['precreation_junction_normals']=joins
    PREFLIGHT['bad_jacobian_count']=sum(not r['valid'] for r in jacobians)
    PREFLIGHT['over90_junction_count']=sum(r['first_triangle_normal_jump_deg']>90 for r in joins)
    PREFLIGHT['junction_statistics']={
        label:_stats([r['first_triangle_normal_jump_deg'] for r in joins if r['regions']==regions])
        for label,regions in [('fringe_bridge',[0,1]),('bridge_hood',[1,2])]}
    assert PREFLIGHT['bad_jacobian_count']==0 and PREFLIGHT['over90_junction_count']==0, (
        'Fixed cut/correspondence still folds or collapses; no object created',
        PREFLIGHT['bad_jacobian_count'],PREFLIGHT['over90_junction_count'])


class _Builder:
    def __init__(self,points):
        self.points=list(points)
        self.faces=[]
        self.info=[]
        self.edges={}

    def point(self,p):
        self.points.append(p.copy())
        return len(self.points)-1

    def add(self,face,info,anchor=False):
        face=tuple(face)
        assert len(set(face))==len(face)
        normal=(self.points[face[1]]-self.points[face[0]]).cross(self.points[face[2]]-self.points[face[0]])
        assert normal.length>1e-14,('Degenerate new face',info,face)
        if anchor:
            decisions=[]
            for a,b in zip(face,face[1:]+face[:1]):
                inc=self.edges.get(tuple(sorted((a,b))),[])
                if inc:
                    assert len(inc)==1,('Already closed seam edge',a,b,info)
                    decisions.append(inc[0][1]==(a,b))
            assert decisions and all(v==decisions[0] for v in decisions),('Inconsistent seam winding',info,decisions)
            if decisions[0]:
                face=tuple(reversed(face))
                info=dict(info)
                info['uv']={key:list(reversed(value)) for key,value in info.get('uv',{}).items()}
        fi=len(self.faces)
        self.faces.append(face)
        self.info.append(info)
        for a,b in zip(face,face[1:]+face[:1]):
            self.edges.setdefault(tuple(sorted((a,b))),[]).append((fi,(a,b)))
        return fi


def _ownership(scene,contract,mesh_regions):
    rows={}
    dg=bpy.context.evaluated_depsgraph_get()
    scale=contract['camera']['ortho_scale_m']
    for view,ys,xrange in [('side',(146,160,180,200,220),(180,235)),
                           ('three_quarter',(150,160,175,190,210),(225,307)),
                           ('three_quarter_mirror',(150,160,175,190,210),(100,280))]:
        spec=contract['fixed_views'][view]
        rotation=Euler(spec['rotation_euler_rad']).to_matrix()
        direction=rotation@Vector((0,0,-1))
        rows[view]=[]
        for y in ys:
            segments=[]
            for x in range(*xrange):
                origin=Vector(spec['location_m'])+rotation@Vector((((x+.5)/512-.5)*scale,
                                                                 (.5-(y+.5)/512)*scale,0))
                hit=scene.ray_cast(dg,origin,direction,distance=2.)
                name=hit[4].name if hit[0] else None
                region=mesh_regions[hit[3]] if name==NEW_NAME else None
                key=[name,REGIONS.get(region)]
                if not segments or segments[-1]['owner']!=key:
                    segments.append({'owner':key,'x':[x,x]})
                else:
                    segments[-1]['x'][1]=x
            rows[view].append({'y':y,'segments':segments})
    return rows


def build_head_034c():
    PREFLIGHT['phase']='frozen source and retained geometry capture'
    root=Path(__file__).resolve().parent
    baseline=root/'head_032_candidate.blend'
    cut_path=root/'head_034_shared_topology_probe.json'
    guard_path=root/'head_034_dryrun.json'
    utility=root/'head_032_draft.py'
    assert _sha(baseline)==SOURCE_SHA and _sha(cut_path)==CUT_SHA
    assert _sha(guard_path)==GUARD_SHA and _sha(utility)==UTILITY_SHA
    active=Path(bpy.data.filepath).resolve()
    active_hash=SOURCE_SHA
    if AUTHORIZED_ACTIVE_SOURCE is None:
        assert active==baseline.resolve()
    else:
        assert active==Path(AUTHORIZED_ACTIVE_SOURCE['path']).resolve()
        active_hash=AUTHORIZED_ACTIVE_SOURCE['sha256']
    assert _sha(active)==active_hash and bpy.context.scene.frame_current==1
    assert bpy.context.mode=='OBJECT' and NEW_NAME not in bpy.data.objects
    ns={'__file__':str(utility),'__name__':'shared_interface_utilities'}
    exec(compile(utility.read_text(),str(utility),'exec'),ns)
    scene=bpy.context.scene
    guards=json.loads(guard_path.read_text())['result']['head_input_guards']
    assert len(guards)==15 and guards=={name:ns['_record'](scene.objects[name]) for name in guards}
    controls={o.name:ns['_record'](o) for o in scene.objects if o.type in {'MESH','CURVE'}
              and o.visible_get() and not o.hide_render and o.name not in TARGETS}
    rig=scene.objects['ReimuFumoRig']
    pose={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    cut=json.loads(cut_path.read_text())
    snapshots={name:_capture(scene.objects[name]) for name in TARGETS}
    assert snapshots[TARGETS[0]]['materials']==snapshots[TARGETS[1]]['materials']
    mats=snapshots[TARGETS[0]]['materials']
    uv_names=sorted({key for s in snapshots.values() for key in s['uv_layers']})
    base_counts={name:len(scene.objects[name].data.vertices) for name in TARGETS}
    face_counts={name:len(scene.objects[name].data.polygons) for name in TARGETS}
    offsets={TARGETS[0]:0,TARGETS[1]:len(snapshots[TARGETS[0]]['points'])}
    all_points=snapshots[TARGETS[0]]['points']+snapshots[TARGETS[1]]['points']
    builder=_Builder(all_points)
    for name in TARGETS:
        obj=scene.objects[name]
        mapping=ns['_head_map'](obj)
        s=snapshots[name]
        assert len(s['points'])==2*base_counts[name]
        assert max((mapping@v.co-s['points'][i]).length for i,v in enumerate(obj.data.vertices))<2e-8
        for fi,p in enumerate(obj.data.polygons):
            assert set(v%base_counts[name] for v in s['polygons'][fi+face_counts[name]])==set(p.vertices)
    fp=snapshots[TARGETS[0]]['points']
    hp=snapshots[TARGETS[1]]['points']
    fpath,ft=_ordered(cut['fringe']['cut_edges'],fp)
    hpath,ht=_ordered(cut['hood']['gusset_receiver_cut_edges'],hp)
    assert len(fpath)==337 and len(hpath)==97
    knots=sorted({round(t,12) for t in ft+ht})
    splits={}
    boundary={}
    normals={}
    for name,path,ts in [(TARGETS[0],fpath,ft),(TARGETS[1],hpath,ht)]:
        snap=snapshots[name]
        for inner in (False,True):
            shift=base_counts[name] if inner else 0
            indices=[]
            nsamples=[]
            shifted=[i+shift for i in path]
            for t in knots:
                p,original,edge=_curve_sample(shifted,ts,t,snap['points'])
                if original is not None:
                    index=offsets[name]+original
                    normal=snap['normals'][original]
                else:
                    a,b,u=edge
                    index=builder.point(p)
                    ga,gb=offsets[name]+a,offsets[name]+b
                    fraction=u if ga<gb else 1.-u
                    splits.setdefault(tuple(sorted((ga,gb))),[]).append((fraction,index))
                    normal=snap['normals'][a].lerp(snap['normals'][b],u).normalized()
                indices.append(index)
                nsamples.append(normal.copy())
            boundary[(name,inner)]=indices
            normals[(name,inner)]=nsamples
    for edge,values in splits.items():
        values.sort()
    removed_by_name={TARGETS[0]:set(cut['fringe']['removed_faces']),TARGETS[1]:set(cut['hood']['removed_faces'])}
    root_edges={tuple(sorted(e)) for e in cut['fringe']['original_open_edges']}
    front_cover=set(cut['hood']['unchanged_front_cover_faces'])
    retained_map=[]
    removed_rims=0
    supported_triangles=0
    for name in TARGETS:
        snap=snapshots[name]
        nv,nf=base_counts[name],face_counts[name]
        removed=set(removed_by_name[name])|{i+nf for i in removed_by_name[name]}
        if name==TARGETS[0]:
            for fi,f in enumerate(snap['polygons']):
                if fi>=2*nf and tuple(sorted({i%nv for i in f})) in root_edges:
                    removed.add(fi)
                    removed_rims+=1
        for ti,tri in enumerate(snap['triangles']):
            if tri['polygon'] in removed:
                continue
            face=tuple(offsets[name]+i for i in tri['vertices'])
            region=(0 if name==TARGETS[0] else 3 if tri['polygon'] in front_cover else 2) if tri['polygon']<nf else 4
            info={'material':tri['material'],'smooth':tri['smooth'],'uv':tri['uv'],
                  'region':region,'source':[name,ti],'kind':'retained'}
            perimeter=[]
            uv={key:[] for key in uv_names}
            touched=False
            for j,(a,b) in enumerate(zip(face,face[1:]+face[:1])):
                perimeter.append(a)
                for key in uv_names:
                    uv[key].append(tri['uv'].get(key,[(0.,0.)]*3)[j])
                keyedge=tuple(sorted((a,b)))
                values=splits.get(keyedge,[])
                if a>b:
                    values=[(1.-u,i) for u,i in reversed(values)]
                for u,index in values:
                    touched=True
                    perimeter.append(index)
                    for key in uv_names:
                        source_uv=tri['uv'].get(key,[(0.,0.)]*3)
                        uv[key].append(tuple((1.-u)*source_uv[j][k]+u*source_uv[(j+1)%3][k] for k in range(2)))
            if not touched:
                fi=builder.add(face,info)
                retained_map.append((name,ti,fi))
            else:
                supported_triangles+=1
                center=builder.point(sum((builder.points[i] for i in face),Vector())/3.)
                center_uv={key:tuple(sum(v[k] for v in tri['uv'].get(key,[(0.,0.)]*3))/3. for k in range(2)) for key in uv_names}
                for j,a in enumerate(perimeter):
                    b=perimeter[(j+1)%len(perimeter)]
                    part=dict(info)
                    part['kind']='boundary_edge_subdivision'
                    part['uv']={key:[uv[key][j],uv[key][(j+1)%len(perimeter)],center_uv[key]] for key in uv_names}
                    builder.add((a,b,center),part)
    assert removed_rims==306
    # The two arcs have one shared longitudinal partition. Width rows are
    # an endpoint-tangent Hermite strip, not two coincident sheet envelopes.
    PREFLIGHT['phase']='oriented boundary co-normal derivatives'
    endpoint_units={}
    for inner in (False,True):
        for name,side_sign in ((TARGETS[0],1.),(TARGETS[1],-1.)):
            endpoint_units[(name,inner)]=_co_normals(
                builder,boundary[(name,inner)],side_sign,name+(' inner' if inner else ' outer'))
    invalid=[r for r in PREFLIGHT['endpoint_directions'] if not r['valid']]
    PREFLIGHT['invalid_endpoint_count']=len(invalid)
    assert not invalid, ('No consistently outward transverse derivative on fixed cut',invalid[:2])
    layers={}
    chord_lengths=[]
    for inner in (False,True):
        a=boundary[(TARGETS[0],inner)]
        b=boundary[(TARGETS[1],inner)]
        rows=[a]
        for level in range(1,8):
            t=level/8.
            row=[]
            for j,(ia,ib) in enumerate(zip(a,b)):
                pa,pb=builder.points[ia],builder.points[ib]
                chord=pb-pa
                if not inner and level==1:
                    chord_lengths.append(chord.length)
                da=endpoint_units[(TARGETS[0],inner)][j]*chord.length
                db=endpoint_units[(TARGETS[1],inner)][j]*chord.length
                assert chord.length<.014 and min(da.length,db.length)>.00005,('Unusable endpoint tangent',j,chord.length)
                if level==1:
                    PREFLIGHT['derivative_vectors'].append({
                        'inner':inner,'knot':j,'fringe_point':list(pa),'hood_point':list(pb),
                        'chord':list(chord),'derivative_fringe':list(da),
                        'derivative_hood':list(db),'magnitude_m':chord.length})
                p=(pa*(2*t**3-3*t*t+1)+da*(t**3-2*t*t+t)
                   +pb*(-2*t**3+3*t*t)+db*(t**3-t*t))
                row.append(builder.point(p))
            rows.append(row)
        rows.append(b)
        layers[inner]=rows
        for level,(r0,r1) in enumerate(zip(rows,rows[1:])):
            for j in range(len(knots)-1):
                info={'material':0,'smooth':True,'uv':{key:[(knots[j],level/8.),(knots[j+1],level/8.),
                                                              (knots[j+1],(level+1)/8.),(knots[j],(level+1)/8.)] for key in uv_names},
                      'region':4 if inner else 1,'source':None,'kind':'inner_bridge' if inner else 'outer_bridge'}
                middle=(level+.5)/8.
                expected=((normals[(TARGETS[0],inner)][j]+normals[(TARGETS[0],inner)][j+1])*(1.-middle)
                          +(normals[(TARGETS[1],inner)][j]+normals[(TARGETS[1],inner)][j+1])*middle).normalized()
                info['expected_normal']=list(expected)
                builder.add((r0[j],r0[j+1],r1[j+1],r1[j]),info,anchor=True)
    for a,b in cut['hood']['retained_front_cover_cut_edges']:
        outer_a,outer_b=offsets[TARGETS[1]]+a,offsets[TARGETS[1]]+b
        inner_a,inner_b=outer_a+base_counts[TARGETS[1]],outer_b+base_counts[TARGETS[1]]
        builder.add((outer_a,outer_b,inner_b,inner_a),{'material':0,'smooth':True,'uv':{},'region':4,
                    'source':None,'kind':'hidden_front_cover_allowance'},anchor=True)
    for j in (0,len(knots)-1):
        for level in range(8):
            face=(layers[False][level][j],layers[False][level+1][j],
                  layers[True][level+1][j],layers[True][level][j])
            builder.add(face,{'material':0,'smooth':True,'uv':{},'region':4,'source':None,'kind':'side_connector'},anchor=True)
    bad_edges={str(e):inc for e,inc in builder.edges.items() if len(inc)!=2 or inc[0][1]!=tuple(reversed(inc[1][1]))}
    assert not bad_edges,('Shell not closed consistently',len(bad_edges),list(bad_edges.items())[:8])
    thickness=[(builder.points[a]-builder.points[b]).length for ra,rb in zip(layers[False],layers[True]) for a,b in zip(ra,rb)]
    assert min(thickness)>.0002 and max(thickness)<.003,('Shell thickness range',min(thickness),max(thickness))
    PREFLIGHT['phase']='local strip Jacobian and junction checks before object creation'
    _strip_preflight(builder)
    PREFLIGHT['phase']='precreation checks passed'
    # All array/topology checks precede creation or hiding of scene objects.
    used=sorted({i for f in builder.faces for i in f})
    remap={old:new for new,old in enumerate(used)}
    bone=rig.pose.bones['Head']
    deformation=rig.matrix_world@bone.matrix@bone.bone.matrix_local.inverted()@rig.matrix_world.inverted()
    local=[deformation.inverted()@builder.points[i] for i in used]
    mesh=bpy.data.meshes.new(NEW_NAME+' mesh')
    mesh.from_pydata(local,[],[tuple(remap[i] for i in f) for f in builder.faces])
    for name in mats:
        mesh.materials.append(bpy.data.materials[name])
    for p,info in zip(mesh.polygons,builder.info):
        p.material_index=info['material']
        p.use_smooth=info['smooth']
    for key in uv_names:
        layer=mesh.uv_layers.new(name=key)
        for p,info in zip(mesh.polygons,builder.info):
            values=info['uv'].get(key,[(0.,0.)]*len(p.vertices))
            for loop,value in zip(p.loop_indices,values):
                layer.data[loop].uv=value
    region_attr=mesh.attributes.new('head034_face_region','INT','FACE')
    for item,info in zip(region_attr.data,builder.info):
        item.value=info['region']
    obj=bpy.data.objects.new(NEW_NAME,mesh)
    scene.collection.objects.link(obj)
    group=obj.vertex_groups.new(name='Head')
    group.add(list(range(len(mesh.vertices))),1.,'REPLACE')
    modifier=obj.modifiers.new('034 live Head attachment','ARMATURE')
    modifier.object=rig
    obj['head034_face_regions']=json.dumps(REGIONS,sort_keys=True)
    for name in TARGETS:
        old=scene.objects[name]
        old.hide_render=True
        old.hide_viewport=True
        old.hide_set(True)
    bpy.context.view_layer.update()
    assert controls=={name:ns['_record'](scene.objects[name]) for name in controls}
    assert pose=={b.name:[list(row) for row in b.matrix] for b in rig.pose.bones}
    result=_capture(obj)
    # Existing off-support triangles keep their indexed coordinates and UV
    # corners; only boundary-adjacent triangles were split in their own plane.
    eval_mesh=ns['_mesh'](obj)
    errors=[]
    uv_exact=True
    for name,ti,fi in retained_map:
        tri=snapshots[name]['triangles'][ti]
        original=[snapshots[name]['points'][i] for i in tri['vertices']]
        face=eval_mesh['faces'][fi]
        errors.extend((eval_mesh['points'][j]-p).length for j,p in zip(face,original))
        uv_exact=uv_exact and builder.info[fi]['uv']==tri['uv']
    assert max(errors,default=0.)<2e-8 and uv_exact
    contract=json.loads((root.parents[2]/'projects/renders/assets/reimu_fumo/review_contract.json').read_text())
    ownership=_ownership(scene,contract,[info['region'] for info in builder.info])
    # Report source-region progression, not an object-name-only success.
    ownership_failures=[]
    for view,rows in ownership.items():
        for row in rows:
            seq=[s['owner'][1] for s in row['segments'] if s['owner'][0]==NEW_NAME]
            if any(v in {'inner_skin_or_allowance','retained_front_cover'} for v in seq):
                ownership_failures.append({'view':view,'y':row['y'],'reason':'Hidden region first-hit','sequence':seq})
            simple=[v for i,v in enumerate(seq) if i==0 or v!=seq[i-1]]
            if len(simple)!=len(set(simple)):
                ownership_failures.append({'view':view,'y':row['y'],'reason':'Repeated visible region','sequence':seq})
    assert _sha(active)==active_hash and _sha(baseline)==SOURCE_SHA
    return {'source_sha256':SOURCE_SHA,'active_source_sha256':active_hash,'head_input_guards':guards,
            'targets':list(TARGETS),'created_names':[NEW_NAME],'protected_control_count':len(controls),
            'protected_controls_exact':True,'rig_pose_exact':True,'live_binding':'One Armature, full unit Head weights',
            'thickness':'Captured existing two skins; no Solidify modifier or rig bake',
            'retained_off_support_triangles':len(retained_map),'off_support_coordinate_max_error_m':max(errors,default=0.),
            'retained_uv_corners_exact':uv_exact,'source_uv_layers':{name:s['uv_layers'] for name,s in snapshots.items()},
            'boundary_supported_source_triangles':supported_triangles,'common_arc_knots':len(knots),
            'removed_original_fringe_rims':removed_rims,'new_vertices':len(used),'new_faces':len(builder.faces),
            'nonmanifold_or_inconsistent_edges':0,'strip_chord_lengths_m':_stats(chord_lengths),
            'strip_paired_skin_separation_m':_stats(thickness),'face_region_counts':{REGIONS[k]:sum(i['region']==k for i in builder.info) for k in REGIONS},
            'ownership':ownership,'ownership_failures':ownership_failures,
            'suitability':'STOP before writer' if ownership_failures else 'Topology/ownership passed sampled checks; root visual and intersection review still required',
            'derivative_preflight':PREFLIGHT,
            'limitations':['No save or render, no visual acceptance.',
                           'Paired skin distances are not a complete thickness or self-intersection audit.',
                           'No animation acceptance; captured skins retain only the ordinary live Head attachment.',
                           'Smooth shading normal preservation at the new support ring requires visual review.']}
```

## head_034c_diagnostic.py

```python
"""One034c co-normal-only first state, rejecting causal precreation failures."""
import hashlib
import json
import math
import traceback
from pathlib import Path

import bpy
from mathutils import Euler, Vector
from mathutils.bvhtree import BVHTree
from mathutils.geometry import intersect_ray_tri

ROOT=Path(__file__).resolve().parent
SOURCE=ROOT/'head_032_candidate.blend'
SOURCE_SHA='6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8'
HELPER=ROOT/'head_034c_draft.py'
HELPER_SHA='68a035190a68d602813e93a501ad25208b68680483bdeb4fed02e3dff20c234d'
OUT=ROOT/'head_034c_diagnostic.json'
assert not OUT.exists()
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest()==SOURCE_SHA
assert hashlib.sha256(HELPER.read_bytes()).hexdigest()==HELPER_SHA
bpy.ops.wm.open_mainfile(filepath=str(SOURCE),load_ui=False)
scope={'__file__':str(HELPER),'__name__':'unchanged_shared_reconstruction'}
exec(compile(HELPER.read_text(),str(HELPER),'exec'),scope)
try:
    construction=scope['build_head_034c']()
except Exception:
    report={'source_sha256':SOURCE_SHA,'helper_sha256':HELPER_SHA,
            'execution_succeeded':False,'error':traceback.format_exc(),
            'preflight':scope.get('PREFLIGHT'),
            'runtime':bpy.app.version_string,'build_hash':bpy.app.build_hash.decode(),
            'created_object_present':scope['NEW_NAME'] in bpy.data.objects,
            'blend_saved':False,'rendered':False,
            'source_and_helper_preserved':hashlib.sha256(SOURCE.read_bytes()).hexdigest()==SOURCE_SHA
                and hashlib.sha256(HELPER.read_bytes()).hexdigest()==HELPER_SHA}
    with OUT.open('x') as handle:
        handle.write(json.dumps(report,indent=2)+'\\n')
    print(json.dumps({k:v for k,v in report.items() if k!='preflight'}))
    raise SystemExit(0)
scene=bpy.context.scene
obj=scene.objects[scope['NEW_NAME']]
dg=bpy.context.evaluated_depsgraph_get()
ev=obj.evaluated_get(dg)
mesh=ev.to_mesh()
mesh.calc_loop_triangles()
points=[ev.matrix_world@v.co for v in mesh.vertices]
triangles=[tuple(t.vertices) for t in mesh.loop_triangles]
tri_polys=[t.polygon_index for t in mesh.loop_triangles]
regions=[mesh.attributes['head034_face_region'].data[p].value for p in tri_polys]
normalmap=ev.matrix_world.to_3x3().inverted().transposed()
normals=[(points[t[1]]-points[t[0]]).cross(points[t[2]]-points[t[0]]).normalized() for t in triangles]
corner_normals=[[(normalmap@mesh.corner_normals[j].vector).normalized() for j in t.loops]
                for t in mesh.loop_triangles]
tree=BVHTree.FromPolygons(points,triangles,all_triangles=True)
contract=json.loads((ROOT.parents[2]/'projects/renders/assets/reimu_fumo/review_contract.json').read_text())
spec=contract['fixed_views']['three_quarter_mirror']
rotation=Euler(spec['rotation_euler_rad']).to_matrix()
scale=contract['camera']['ortho_scale_m']
camera=Vector(spec['location_m'])
direction=rotation@Vector((0,0,-1))


def stats(values):
    values=sorted(values)
    return {'count':len(values),'min':values[0],'median':values[len(values)//2],
            'p95':values[int(.95*(len(values)-1))],'max':values[-1]} if values else None


def angle(a,b):
    return math.degrees(math.acos(max(-1.,min(1.,a.dot(b)))))


def project(p):
    q=rotation.transposed()@(p-camera)
    return [(q.x/scale+.5)*512-.5,(.5-q.y/scale)*512-.5]


def shading_normal(ti,p):
    a,b,c=[points[i] for i in triangles[ti]]
    v0,v1,v2=b-a,c-a,p-a
    d00,d01,d11=v0.dot(v0),v0.dot(v1),v1.dot(v1)
    d20,d21=v2.dot(v0),v2.dot(v1)
    denominator=d00*d11-d01*d01
    v=(d11*d20-d01*d21)/denominator
    w=(d00*d21-d01*d20)/denominator
    weights=(1.-v-w,v,w)
    return sum((n*t for n,t in zip(corner_normals[ti],weights)),Vector()).normalized()


pixel_hits=[]
for x in range(237,253):
    origin=camera+rotation@Vector((((x+.5)/512-.5)*scale,(.5-175.5/512)*scale,0))
    scene_hit=scene.ray_cast(dg,origin,direction,distance=2.)
    hits=[]
    start=origin.copy()
    for unused in range(12):
        p,n,ti,d=tree.ray_cast(start,direction,2.)
        if p is None:
            break
        hits.append({'triangle':ti,'polygon':tri_polys[ti],'vertices':list(triangles[ti]),
                     'region':scope['REGIONS'][regions[ti]],'point_m':list(p),
                     'camera_ray_depth_m':(p-origin).dot(direction),
                     'geometric_normal':list(n),'shading_normal':list(shading_normal(ti,p))})
        start=p+direction*.0000001
    pixel_hits.append({'pixel':[x,175], 'scene_first_object':scene_hit[4].name if scene_hit[0] else None,
                       'scene_first_polygon':scene_hit[3] if scene_hit[0] else None,'ordered_shell_hits':hits})

edges={}
for ti,t in enumerate(triangles):
    for a,b in zip(t,t[1:]+t[:1]):
        edges.setdefault(tuple(sorted((a,b))),[]).append(ti)
seams=[]
for edge,inc in edges.items():
    if len(inc)!=2 or {regions[i] for i in inc} not in ({0,1},{1,2}):
        continue
    pa,pb=[points[i] for i in edge]
    a,b=project(pa),project(pb)
    crossing=None
    if min(a[1],b[1])<=175<=max(a[1],b[1]) and abs(a[1]-b[1])>1e-10:
        t=(175-a[1])/(b[1]-a[1])
        crossing=a[0]+t*(b[0]-a[0])
    rows={'edge':list(edge),'triangles':inc,'polygons':[tri_polys[i] for i in inc],
          'regions':[scope['REGIONS'][regions[i]] for i in inc],
          'points_m':[list(pa),list(pb)],'projected_endpoints':[a,b],
          'normal_jump_deg':angle(normals[inc[0]],normals[inc[1]]),
          'normals':[list(normals[i]) for i in inc],'row175_crossing_x':crossing,
          'midpoint_z_m':(pa.z+pb.z)/2}
    seams.append(rows)
local_seams=[r for r in seams if max(p[0] for p in r['projected_endpoints'])>=235
             and min(p[0] for p in r['projected_endpoints'])<=254
             and max(p[1] for p in r['projected_endpoints'])>=171
             and min(p[1] for p in r['projected_endpoints'])<=179]
crossings=[r for r in seams if r['row175_crossing_x'] is not None and 235<=r['row175_crossing_x']<=254]


def segment_distance(p,a,b):
    v=b-a
    t=max(0.,min(1.,(p-a).dot(v)/v.length_squared)) if v.length_squared else 0.
    return (p-(a+v*t)).length


def coplanar_area(ta,tb,normal):
    axis=max(range(3),key=lambda k:abs(normal[k]))
    dims=[i for i in range(3) if i!=axis]
    a=[(points[i][dims[0]],points[i][dims[1]]) for i in ta]
    b=[(points[i][dims[0]],points[i][dims[1]]) for i in tb]
    def cross2(u,v):
        return u[0]*v[1]-u[1]*v[0]
    orientation=1. if cross2((b[1][0]-b[0][0],b[1][1]-b[0][1]),
                             (b[2][0]-b[0][0],b[2][1]-b[0][1]))>=0 else -1.
    output=a
    for q,r in zip(b,b[1:]+b[:1]):
        previous=output
        output=[]
        if not previous:
            break
        def side(p):
            return orientation*cross2((r[0]-q[0],r[1]-q[1]),(p[0]-q[0],p[1]-q[1]))
        for s,e in zip(previous,previous[1:]+previous[:1]):
            ds,de=side(s),side(e)
            if (ds>=0)!=(de>=0):
                t=ds/(ds-de)
                output.append((s[0]+t*(e[0]-s[0]),s[1]+t*(e[1]-s[1])))
            if de>=0:
                output.append(e)
    return abs(sum(p[0]*q[1]-p[1]*q[0] for p,q in zip(output,output[1:]+output[:1])))/2 if output else 0.


def contacts(ai,bi):
    ta,tb=triangles[ai],triangles[bi]
    shared=set(ta)&set(tb)
    na,nb=normals[ai],normals[bi]
    if abs(na.dot(nb))>.999999 and max(abs((points[i]-points[ta[0]]).dot(na)) for i in tb)<2e-8:
        area=coplanar_area(ta,tb,na)
        return {'kind':'coplanar_overlap','projected_area_m2':area} if area>1e-12 else None
    hits=[]
    for first,second in ((ta,tb),(tb,ta)):
        target=[points[i] for i in second]
        for ia,ib in zip(first,first[1:]+first[:1]):
            a,b=points[ia],points[ib]
            p=intersect_ray_tri(*target,b-a,a,True)
            if p is None:
                continue
            u=(p-a).dot(b-a)/(b-a).length_squared
            if u<-1e-7 or u>1+1e-7:
                continue
            if any((p-points[i]).length<2e-7 for i in shared):
                continue
            if len(shared)==2 and segment_distance(p,*[points[i] for i in shared])<2e-7:
                continue
            if all((p-q).length>2e-7 for q in hits):
                hits.append(p)
    return {'kind':'nonshared_segment_contact','points_m':[list(p) for p in hits]} if hits else None


outer=[i for i,r in enumerate(regions) if r==1]
outer_tree=BVHTree.FromPolygons(points,[triangles[i] for i in outer],all_triangles=True)
broad=outer_tree.overlap(tree)
pairs=set()
for local,other in broad:
    a=outer[local]
    if a!=other:
        pairs.add(tuple(sorted((a,other))))
intersection_evidence=[]
for a,b in sorted(pairs):
    contact=contacts(a,b)
    if contact:
        intersection_evidence.append({'triangles':[a,b],'polygons':[tri_polys[a],tri_polys[b]],
                                      'regions':[scope['REGIONS'][regions[a]],scope['REGIONS'][regions[b]]],
                                      'contact':contact})

report={'source_sha256':SOURCE_SHA,'unchanged_helper_sha256':HELPER_SHA,
        'purpose':'One co-normal endpoint correction; unchanged cuts, correspondence, skins and off-support geometry.',
        'execution_succeeded':True,'derivative_preflight':scope['PREFLIGHT'],
        'runtime':bpy.app.version_string,'build_hash':bpy.app.build_hash.decode(),
        'mirror_row175_hits':pixel_hits,'local_seam_edges':local_seams,'row175_seam_crossings':crossings,
        'upper_endpoint_seam_jump_statistics':{
            label:stats([r['normal_jump_deg'] for r in seams if r['midpoint_z_m']>.145
                         and region in r['regions']])
            for label,region in [('fringe_to_bridge','retained_fringe_outer'),('bridge_to_hood','retained_hood_outer')]},
        'all_seam_edges':seams,
        'outer_bridge_intersection_check':{'outer_bridge_triangles':len(outer),'broad_phase_unique_pairs':len(pairs),
                                          'reported_contact_pairs':len(intersection_evidence),'contacts':intersection_evidence,
                                          'shared_contact_tolerance_m':2e-7,'coplanar_projected_area_threshold_m2':1e-12},
        'limitations':['Intersection test is bounded to the actual outer bridge against this captured shell, not the whole scene.',
                       'Shared edge/vertex contacts within0.2 micrometer are allowed; coplanar area threshold1e-12m2.',
                       'No envelope-overshoot, whole-scene collision or visual audit.',
                       'Frozen source guards are reconstruction prerequisites, not a new broad audit.']}
ev.to_mesh_clear()
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest()==SOURCE_SHA
assert hashlib.sha256(HELPER.read_bytes()).hexdigest()==HELPER_SHA
report['source_and_helper_preserved']=True
report['blend_saved']=False
report['rendered']=False
with OUT.open('x') as handle:
    handle.write(json.dumps(report,indent=2)+'\n')
print(json.dumps({'output':str(OUT),'crossing_count':len(crossings),
                  'seam_statistics':report['upper_endpoint_seam_jump_statistics'],
                  'intersection_pairs':len(intersection_evidence)}))
```
