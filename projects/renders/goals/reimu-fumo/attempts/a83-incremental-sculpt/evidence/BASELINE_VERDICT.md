# C19 whole-model baseline viability verdict

## Verdict

**B — build a new coarse neutral sculpt. Exact C1b is categorically nonviable
as the direct-sculpt baseline for the whole-model macro pass.**

Preserve exact C1b, SHA-256
`d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`,
only as an immutable comparison and recovery artifact. Do not transfer its
macro silhouette into the next candidate, and do not continue correcting the
whole model by accumulating local shape keys or replacement coupons on C1b.

This does not say that every C1b component is unusable. The face graphics,
review cameras, lighting, and some accessory graphics remain useful witnesses.
It says that the visible form scaffold joining head, hair, bow, torso,
garment, sleeves, hem, and feet is too far from the reference to converge by a
sequence of protected local sculpts.

## Review scope and independence

The visual judgment used:

- the canonical front, all canonical turn directions, clean front, physical
  front and side, the supporting turn mosaic, and sofa samples;
- exact C1b fixed front and three-quarter renders from the latest C18 review
  packet; and
- no Blender scene, object names, topology, modifier, script, or generated
  candidate inspection.

The controlling authorities are
`projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`
for exact frontal identity and
`projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif`
for depth, rear silhouette, layer order, and seated volume. Physical sources
only veto implausible fabric construction.

Context limitation: the existing C18 textual review was visible before the
final whole-model audit, so this was not a perfectly context-naive review.
However, the verdict below was made from C1b pixels rather than C18 geometry
or implementation, and the ranked discrepancies are present in both frozen
C1b views independently of the C18 sleeve experiment.

## Absolute image review

- Same subject and exact variant without a label: **borderline recognition by
  color/graphics, no by physical-form likeness**.
- Intended medium: **no**. It reads as a smooth stylized figurine assembled
  from hard or inflated primitives, not as a 25 cm constructed plush.
- Overall reference likeness: **4/10**.
- Macro silhouette and proportions: **3/10**.
- Manufactured construction: **2/10**.
- Identity-defining forms: **4/10**.
- Contact and occlusion: **3/10**.
- Intended-medium read: **2/10**.
- Presentation readability: **8/10**.
- Major visible failure: **yes**.
- Absolute decision: **reject as a whole-model sculpt baseline**.

## Five largest whole-model discrepancies

### 1. Head and hair are one deep helmet instead of a shallow cushion with layered fabric hair

C1b's three-quarter silhouette is dominated by a nearly spherical brown mass:
the rear projects deeply, the highlight rolls continuously from crown to
nape, the maximum width remains high, and the lower edge resolves into a few
small spikes. The canonical turn instead shows a compact cushion, a close
crown field, and long overlapping rear/profile leaves whose free edges and
unequal tips own the lower and rear silhouette. The front also shows a broad,
smooth dome around the face rather than the reference's tapered, pile-softened
hair frame.

This is an identity-defining macro failure, not missing seams or surface
detail.

### 2. The bow is a rigid horizontal bar rather than a stuffed, crown-seated, drooping assembly

C1b's upper lobes form a thin, very wide, nearly straight symmetric bar with
hard end flares; the lower tails are similarly geometric and detached in
motion from the head. In the references the center is compressed into the
crown, each lobe has substantial padded height and asymmetric fold tension,
and the tails hang and twist around the hair volume. The side and rear turn
make clear that the bow participates in the head's depth and contact rather
than sitting as a flat logo behind it.

### 3. The lower body is a standing dress cone, not a compact seated plush

C1b's torso and dress form a tall, clean triangular shell that widens to a
straight floor-level hem. The references show a short stuffed core largely
hidden by a skirt that pools around the thighs, breaks asymmetrically at the
floor, and changes overlap as the plush rotates. The physical side and
canonical turn show rearward seated volume and broad fabric contact; C1b shows
neither convincing pelvis compression nor a supported sitting mass.

### 4. Both sleeves are arm-like tubes or rigid cones instead of broad collapsing cloth panels

The C1b sleeves descend as narrow geometric solids with circular openings and
uniform tension. Reference sleeves are large flattened bell panels: broad
front/back faces, a thin soft free edge, root folds, and asymmetric collapse
against the skirt and body. Their width is a major part of the frontal
silhouette. C1b therefore loses both the Fumo proportion and the manufactured
cloth logic.

### 5. Feet, hem, and ground contact have the wrong depth order

C1b places smooth dark balls in an almost straight frontal row beneath a
stiff hem; the three-quarter view reads multiple disconnected beads and a
skirt edge floating or bridging above them. The references show two short
stuffed foot pods attached to pale leg fabric, projected forward at different
depths, partly occluded by a pooled ruffled hem, with the body and garment
sharing a broad seated contact patch.

## Shared structural owners

The failures are not five independent polish tasks. They cluster under two
coupled macro owners:

1. **Cranial attachment and visible-envelope scaffold** — discrepancies 1
   and 2 share the receiver depth, crown plane, bow-root seat, and ownership
   of the side/rear silhouette. Local hair edits cannot compact the receiver
   and reseat the bow without changing the same attachment field. Existing
   C11--C16 evidence further shows that the continuous cap loses either
   coverage or layered construction when pushed far enough to change the
   macro pixels.
2. **Seated torso/pelvis and garment-contact scaffold** — discrepancies 3, 4,
   and 5 share shoulder placement, torso height, hip/contact volume, skirt
   support, sleeve roots, foot anchors, and occlusion order. A local sleeve or
   foot replacement cannot make the model sit while the dress remains a tall
   cone; conversely, changing the skirt envelope necessarily changes sleeve
   contact and foot exposure. C17 and C18 already show that isolated sleeve
   representation swaps remain wrong in the unchanged C1b body context.

These two owners meet at the collar/shoulder interface. The required changes
would move most visible macro extrema and contacts, so calling them a series
of local C1b sculpts would preserve the name of the baseline while replacing
its actual structure.

## Why option A is rejected

A direct local-sculpt sequence on C1b is appropriate only when the frozen
baseline already has a viable macro scaffold and isolated owners can reach
the failed pixels without violating protected views. C1b fails that test:

- at least four identity-defining silhouettes are materially wrong together;
- the largest failures share coupled attachment/contact owners rather than
  separable surface patches;
- the head/hair branch has already exhausted local, whole-field, and layered
  boundary edits without simultaneously preserving coverage and achieving a
  fabric-panel read; and
- isolated sleeve experiments changed the local object but could not create
  the required seated-body and garment relationship.

Further C1b sculpt cycles would optimize surfaces around the wrong depth,
attachment, and contact scaffold. That is churn, not a lower-risk path.

## Recommended option B: new coarse neutral sculpt

Start from an empty, task-owned candidate file while keeping C1b linked only
as a hidden comparison source. Build a deliberately coarse, animatable owner
scaffold in two modules, each in full-plush context:

1. **Head module:** shallow gusseted face/head cushion, compact crown field,
   one broad rear/profile hair envelope plus one independent foreground leaf,
   and padded bow proxies seated into the crown.
2. **Seated-body module:** short stuffed torso/pelvis mass, pooled skirt
   envelope, broad flattened sleeve-panel proxies, and two foot pods with
   explicit hem occlusion and ground contact.

Keep each owner separate and transformable. Use low-resolution subdivision
surfaces or coarse sculpt meshes with enough topology for broad planar
flattening and contact compression; do not merge them into one static
sculpt. Defer eyes, embroidery, ruffles, seams, fibers, fine locks, detailed
folds, materials, and production rig transfer.

This is a new coarse sculpt, not a one-shot finished model. Work one owner at
a time, but judge every owner in the complete front/side/rear/three-quarter
silhouette so a local improvement cannot hide a global regression.

## First decisive artifact

Produce one immutable **C19-0 neutral macro blockout packet** from a single
candidate state:

- 512 px fixed front, both front three-quarters, both profiles, and rear;
- a neutral clay pass with no identity graphics except optional flat eye
  landmarks for camera alignment;
- an identical owner-ID pass distinguishing head cushion, crown hair, rear
  hair, foreground leaf, each bow role, torso/pelvis, skirt, both sleeves, and
  both feet;
- reference/candidate silhouette overlays aligned by head width, crown, eye
  line, and ground plane; and
- an uncropped normal-scale contact sheet showing the entire plush in every
  view.

The packet is decisive because it tests the two new structural owners before
detail investment. It must be made from one unchanged geometry state; do not
mix best views from different variants.

## Stop gate

Reject C19-0 immediately and do not add detail if any of these is true at
normal scale:

- a critical canonical-view landmark differs by more than `3%` of head width,
  or a major silhouette extremum/gap differs by more than `5%` after camera
  alignment;
- the head still reads as a sphere, egg, helmet, box, or uninterrupted
  crown-to-nape dome in either profile or three-quarter;
- the bow remains a flat horizontal bar or lacks visible crown seating and
  lobe/tail depth order;
- the torso/skirt remains a standing cone or lacks a broad seated contact and
  pooled hem in profile/rear;
- a sleeve reads as a tube, limb, paddle, or inflated mitten rather than a
  broad panel with a readable free edge;
- feet read as disconnected balls, more than two foot masses, or lack the
  reference hem/leg/foot occlusion order; or
- owner-ID reveals that an underlying support object, rather than the intended
  constructed panel, owns any failed outer silhouette.

If exactly one owner fails while the other module passes, freeze the passing
module and rebuild only the failed owner. If either whole module still fails
after two bounded coarse-shape cycles, stop and change that module's scaffold
or attachment interface before any additional sculpting. Only after the full
C19-0 packet passes the macro gate may the model proceed to manufactured
panel construction, then stuffing/folds, then materials and microdetail.

## Final recommendation

Use **B**. Freeze C1b as the unaccepted comparator, not the clay to keep
sculpting. The simplest path that still reaches the references is a new,
coarse, modular neutral blockout whose two structural owners are proven in a
single six-direction packet before any detailed asset work resumes.
