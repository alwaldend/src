# Sleeve construction diagnosis and bounded replacement proposal

Observed 2026-09-06, independently alongside the coordinator's hair work.
This is a design proposal, not a model, acceptance review, or authorization to
alter another subsystem. Only the named inspection script/JSON/log and this
report were written. One pinned Blender read-only process completed exit0;
no render or model was saved. Exact076 source remained
`c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.

## Finding from images

The sleeve mouth is short, nearly square, and uniformly rounded. It reads as
an elastomer tunnel, not a hanging sewn garment. This is apparent in076 side
and three-quarter renders before inspecting its implementation. The narrow
root at the torso and larger opening are broadly correct; simply making the
opening smaller would not correct the construction.

Canonical turn frames08/10/12 control aperture and layer order: the opening
is a tall irregular oval, with an upper region supported by the arm/hand and
a long unsupported lower region dropping toward the skirt. Its lower lip
is softer than its upper edge and the two sides do not make a rigid plane.
Physical_side controls the thin folded edge and shallow wrinkled cloth wall,
not exact aperture proportions: that physical variant has a shorter, more
rounded opening than the canonical turn. Physical_front and canonical_front
show broad front panels with comparatively little hand exposed.076 shows
noticeable projecting cream hand wedges in front.

Wh for numeric candidate normalization is the frozen076 head width
0.116195146m. The reference ranges below are manual image bands, NOT a
camera-calibrated acceptance overlay. Estimated turn image scale derives
from the head's near-front height; perspective/occlusion contributes roughly
0.05Wh uncertainty. Use the fixed side and3q pixels to disconfirm the design.

| Quantity | Reference evidence/target | Exact076 candidate |
| --- | --- | --- |
| Combined front sleeve span | Canonical landmark1.277Wh±0.05 |1.273Wh evaluated |
| Cuff aperture side height/width | Turn08–12 about1.5–1.8;±0.2 uncertain |0.915 at outer raw cuff |
| Cuff vertical extent | Turn approximately0.45–0.58Wh; oblique/uncertain |0.289Wh |
| Cuff depth extent | Turn approximately0.27–0.36Wh; uncertain |0.316Wh |
| Cloth folded-edge envelope | Physical side thin, variable; approximately0.006–0.015Wh | about0.022Wh nominal maximum:1.8mm turn+0.8mm shell |

Do not achieve the tall opening by blindly extending the whole sleeve down.
Front sleeve span already matches the canonical datum; a large underside
extension can collide with the skirt or cover the red hem. The likely useful
change is a more vertical/nonplanar aperture, less depth through its lower
half, and independently tensioned upper and slack lower panel boundaries.

## Exact existing interfaces

Inspection JSON SHA256:
`fba8ab08f9f6cde4f44342593b6395792ab3ce3a4160fdd6d32d6d4ea62864f8`.
Script SHA256:
`07fc2e96cb6c3d1cf512a508171689fa809de8a2d67b2b3321883f42f598c31b`.
Blender5.2.1 LTS/build9e2066aef7ef, background/factory/disable-autoexec.
Full evaluated bounds, transforms,16samples of six sleeve rings, materials,
modifiers, and intersection counts are in `sleeve_077_inspect.json`.

Per side, the relevant objects are:

- `Macro066 left/right supported bell sleeve`:3600 raw vertices,7200
  evaluated; single SOLIDIFY, thickness0.8mm, offset0. Evaluated positive
  volume,0nonmanifold edges,0self-overlap pairs. Ivory material
  `Fabric041 ivory cloth`.
- `Macro066 left/right sleeve red dashes`:342raw vertices,684evaluated;
 19 separate dash patches, each9×2 vertices, SOLIDIFY0.06mm. They must be
  remapped to the changed garment; do not retain world coordinates.
- `Macro066 left/right continuous stuffed arm`:1106vertices, no modifiers,
  `Fabric041 warm cream`. Center(±39.5,-6,56)mm; local scale(9,10,34.518)mm.
  Original067 start(±16,0,69)mm/end(±63,-12,43)mm describe its long axis.
- `Macro038 left/right collar`:protected. CombinedXYZ bounds approximately
  X±30.25mm,Y[-28.975,-14.783]mm,Z[67.428,81.188]mm.
- Torso exact name, from071 construction:
  `Macro038 leaning chest and seated hips`. The inspection's name filter
  omitted this object, so its current bounds/contact counts are UNOBSERVED;
  do not infer a zero core-intersection result.076 preserved the071 body.

Right sleeve evaluated bounds (left mirrorsX):
X[19.131,73.971],Y[-21.154,16.394],Z[27.721,76.930]mm.
Right sleeve original root ring0:
X[19.423,28.577],Y[-13,7],Z[63.439,76.561]mm.
Outer cuff ring40:
X[47.662,73.248],Y[-20.702,16.000],Z[28.148,61.734]mm.
Return ring44:
X[48.702,72.110],Y[-18.905,14.200],Z[29.754,60.322]mm.
Right arm evaluated bounds:
X[9.611,69.389],Y[-17.562,5.562],Z[37.471,74.529]mm.

Sleeve/arm triangle intersections are310left and296right. These raw pair
counts do not establish visible penetration severity, but they invalidate an
assumption that current sleeve/arm contacts are already clean. A replacement
must localize them and allow only intentional concealed shoulder seating,
not preserve them unquestioned.

## Cause and proposed replacement

067 builds every cross-section from the same radial superellipse
`sign(cosθ)|cosθ|^.66`, `sign(sinθ)|sinθ|^.76`. Both powers below1 flatten
the four quadrants toward a square. All cross-sections inherit one arm axis;
the lower cloth never becomes a separately suspended bag-like panel.068's
0.9mm180-degree normal turn makes a smooth thick perimeter all the way
around. Its technical intersection repair worked, but retained that visual
construction. Another amplitude/fold adjustment within these rings is not
the proposed reset.

Replace the sleeve with TWO longitudinal garment charts, front and back,
sharing an upper seam and an underarm seam. Use actual shoulder, arm envelope,
front silhouette, and cuff silhouette as independent boundary constraints.
This is still a hollow garment, but no radial axis determines its surface.

Root-callable construction sketch (intentionally not a complete builder):

1. For each sign, expose `panel(u,v,face)`,u∈[0,1] shoulder→cuff,
   v∈[0,1] upper-seam→underarm-seam,face∈{front,back}. Capture the old
   shoulder ring as the initial boundary; preserve it where it is concealed.
   Its two semicircles are the two chart root edges.
2. Author upper and lower seam curves independently. Keep the upper/front
   silhouette and combined sleeve span within the existing front bound;
   support the upper cloth over the actual arm instead of a radial radius.
   Give the lower seam a gravity-shaped sag, with its last third approaching
   the measured skirt clearance. The lower seam must NOT remain a scaled
   copy of the upper seam.
3. Author two cuff-edge curves inXYZ sharing upper/lower endpoints. Make
   their lower halves come closer inY than the top; put the greatestY
   opening around the supported upper/central region. Choose an irregular
   oval/dropped-teardrop side outline, not four superellipse quadrants.
   Make one cuff edge extend slightly farther inX and lower inZ than the
   other so the hem plane is flexible rather than one rigid oblique ring.
4. Coons-interpolate each front/back chart from its four boundaries:
   `S=(1-u)root(v)+u*cuff(v)+(1-v)upper(u)+v*lower(u)-bilinearCorners`.
   Add only broad panel bulge along ±Y, tapered to zero at all four edges.
   This has an interpretable cloth panel boundary and does not add radial
   lobes. Share seam vertices when joining charts; retain shoulder and cuff
   openings. Moderate48×40 grids are sufficient before local hem rows.
5. Seat only the supported shoulder/upper chart against evaluated arm/core
   surfaces with a small real fabric clearance. Do not ray-clamp every free
   vertex to the arm: that would erase the hanging lower bag. If preserved
   arm shape cannot fit cleanly, diagnose that specific collision before
   changing the hand; do not conceal it with a huge sleeve.
6. Give fabric roughly0.4–0.6mm genuine wall thickness (construction scale,
   not a material tweak). Make the cuff a narrow2–3mm turn-back following
   the chart's own outgoing tangent; use only the small bend needed to turn
   the cloth, then a mostly straight inner flap. That creates a folded hem,
   not a 360-degree rounded rubber rim. Check the evaluated doubled shell
   at the bend; normal-offset crossings remain possible at tight folds.
7. Reuse the two dash objects/materials. Place their strips along a constant
   garment-chart distance2–4mm inward from the cuff, offset by actual wall
   thickness plus a tiny applique clearance. The current19-patch layout can
   be preserved without inventing more surface detail.

## Decision gate and risk

The dominant risk is matching a tall side aperture while keeping canonical
front proportions and avoiding new lower-cuff/skirt intersections. For a
first bounded construction test, freeze sleeve root/arm/collar transforms,
front outerX extrema, and clear reference-derived aperture dimensions before
building. Inspect side and3q before proceeding to full views. A failure that
remains square/rigid means the boundary charts are still the wrong silhouette;
more fuzz or smaller shell thickness will not fix it.

Pre-save must require closed positive evaluated shell, zero nonadjacent
self-crossings, finite charts, no unsupported sleeve/hand penetration, and
localized intentional shoulder seating. Quantify cuff-vs-skirt contact from
the actual receiver. Post-save blind pixels must show a thin, asymmetrically
hanging cuff and a supported hand, not a hollow appliance. Nothing here
passes the current whole-asset construction gate.
