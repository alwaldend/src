# A74 visible-envelope design review

## Verdict

**REVISE the proposed representation, then proceed.**  Do not build another
receiver-first loft, closed brown shell, or one-piece head envelope.  A72 and
A73 show that even a deterministic, manifold, pole-free volume with correct
gross dimensions still reads as a mattress from front/rear and an egg from
profile.  A46 and A34 further show that separately extruded cards or a ruled
front field merely exchange the mattress for a box, visor, and exposed closure
edges.

The smallest materially different A74 test is a **visible construction first**:

1. one thin, open crown/temple saddle whose front boundary is the observed
   asymmetric hairline;
2. two independent thin cheek-lock pockets;
3. one broad background rear/nape pocket with an unequal lobed free edge; and
4. one broad, off-centre foreground rear leaf crossing that background.

Those five brown objects are the module.  A shallow beige cushion may exist
only as a subordinate contact carrier and exposed face plane; it does not own
the brown silhouette and is not accepted independently.  The rear pockets must
overlap the open saddle beneath the bow so no global hair shell or hidden cape
is needed.

This is not a claim about the factory pattern.  It is the smallest reversible
digital construction that directly expresses the visible free edges,
T-junctions, depth reserve, and unequal rear lobes established by the supplied
references.

## Decision challenge

### Outcome and constraints

The outcome is a recognizable 25 cm canonical Reimu Fumo head/hair assembly,
in the protected whole-subject context, that already reads as constructed
plush in neutral clay.  The canonical front controls frontal identity and
scale; all 30 canonical turn frames control depth, transitions, and layer
order; the physical front/side and sofa/older GIFs control fabric behavior and
veto rigid construction.  No tracked Blender asset may be modified by this
design task.

### Strongest case for another loft or closed shell

A closed shared surface is easy to keep manifold, gives every attachment a
stable receiver, and can interpolate front, profile, and rear measurements in
one coordinate system.  It avoids floating cards and is convenient to rig.

### Decisive disconfirming evidence

- A73 used independent width/front/rear profiles, explicit section anchors,
  pole-free quad caps, a high rear maximum, and a lower-rear undercut.  It still
  scored `4/10`: rounded box/mattress in front/rear/top and egg in profile.
- A72 failed in the same visual family.  The repeated failure is therefore the
  volume representation, not missing exponent or topology polish.
- No supplied source exposes a naked receiver or a closed brown perimeter.
  The measured pixels belong to occluding hair layers.  Assigning them to a
  hidden shell repeats the session's unsupported inference.
- Physical side and the canonical rear arc show thin free hair edges and
  brown-on-brown T-junctions.  A smooth closed shell cannot make those without
  adding the independent panels that actually own the likeness.
- A46/A34 reject the opposite shortcut: independently extruded planar cards,
  artificial bridges, and ruled patches become a box/visor.  The A74 saddle
  and pockets must be genuinely curved 3D surfaces seated against one
  subordinate support, not planar cut-outs connected after the fact.

### Alternatives

| Representation | Reference fit | Main risk | Verdict |
| --- | --- | --- | --- |
| Retune A73 loft | Deterministic but hidden and weakly observed | Third mattress/egg cycle | Reject |
| One watertight brown helmet | Easy coverage and rigging | Smooth helmet, no free-edge construction | Reject |
| Independent flat/extruded cards | Easy tracing | Box, visor, floating roots, rigid felt | Reject |
| Full freehand sculpt | Can create softness | Weak reproducibility and prior ineffective stroke interface | Defer |
| Open saddle + seated fabric pockets | Directly owns observed boundaries and overlaps | Root coverage and curvature must be gated in both 3Q views | Proceed |

The conditions that would reverse this verdict are new reference evidence of a
visible closed shell seam, or a first-pixel failure showing the saddle cannot
cover the support without becoming a helmet.  In the latter case, stop at P0;
do not thicken or extend it into another shell.

## Controlling references and source-of-truth records

Do not select one favored image.  The builder should consume these existing
records rather than retype approximate landmarks:

- canonical exact front, SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`;
- canonical slow 180, all 30 frames, SHA-256
  `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`;
- clean front, physical front, physical side, older turn, and sofa sources
  listed with hashes in
  `out/reimu_fumo_attempt_073_profile_loft/reference_measurements/source_hashes.csv`;
- the complete reference packet:
  `out/reimu_fumo_attempt_073_profile_loft/reference_measurements/all_relevant_reference_packet.png`;
- observed source curves (not hidden interfaces):
  `out/reimu_fumo_attempt_053_head_seam_network_spec/head_hair_curves.json`;
- canonical scanline tables:
  `out/reimu_fumo_attempt_073_profile_loft/reference_measurements/`
  `front_outer_scanlines.csv`, `profile_outer_scanlines.csv`, and
  `rear_outer_scanlines.csv`;
- negative evidence: A34, A46, A69--A73.  None of their geometry is a parent.

The canonical front defines `Wh = 368 +/- 4 px = 132 mm`, horizontal center
`x=485 px`, and visible brown crown `y=231 px`.  Coordinates below use either
`u=(x-301)/368, v=(y-231)/368` from the canonical left hair bound and crown,
or signed `x=(pixel_x-485)/368` where explicitly stated.

## Source-derived controls

### A. Crown/temple saddle outer front boundary

Use the **independent left and right values** below.  They are copied from
`front_outer_scanlines.csv`; interpolation may be shape-preserving, but no
power, superellipse, mirrored average, or post-render rescale may replace
them.

| down `v/Wh` | left reach | right reach | full observed width | tolerance |
| ---: | ---: | ---: | ---: | ---: |
| .052 | .182 | .226 | .408 | .011 |
| .133 | .321 | .353 | .674 | .011 |
| .242 | .402 | .435 | .837 | .011 |
| .351 | .438 | .470 | .908 | .011 |
| .459 | .465 | .484 | .948 | .011 |
| .568 | .484 | .492 | .976 | .011 |

The raw 13-point left and right crown traces in
`head_hair_curves.json` remain the geometry source.  The six scanlines are
preflight and render gates, not a license to construct a six-ring loft.

### B. Free front hairline

These seven canonical pixels are the lower **free edge** of the same saddle,
not separate bang pillows and not a material boundary on a shell.

| point | source px `(x,y)` | normalized `(u,v)` | role |
| --- | --- | --- | --- |
| H0 | (365,438) | (.174,.563) | left temple return |
| H1 | (389,371) | (.239,.380) | left upper opening shoulder |
| H2 | (447,428) | (.397,.535) | broad left sweep |
| H3 | (517,480) | (.587,.677) | blunt off-centre main tip |
| H4 | (540,376) | (.649,.394) | central cleft/upper return |
| H5 | (604,435) | (.823,.554) | right sweep |
| H6 | (624,415) | (.878,.500) | right temple return |

The H3 lower third must remain at least `.12 Wh` broad.  No tangent at a tip
may form a narrow triangular spike.  Apparent front stand-off over beige is
`.02-.05 Wh`; unintended root separation is at most `.01 Wh`.

### C. Cheek locks

Use the raw, non-mirrored polylines
`canonical_front_left_lock_outer` and
`canonical_front_right_lock_outer` in `head_hair_curves.json`.  Their extrema
must preserve:

| quantity | target |
| --- | ---: |
| separate lock width | `.14-.18 Wh` each |
| crown to lowest lock | `1.098 +/- .03 Wh` |
| left lowest source point | `(356,635)` |
| right lowest source point | `(615,634)` |
| padded edge/core thickness | `.02-.04 Wh` |

The locks are two-skin pockets with a narrow edge strip.  They share root
positions with the saddle but keep independent tangents and free lower edges.
They may not carry the head's side depth.

### D. Profile ownership

Do not derive a hidden full profile.  Consume these observed curves directly:

- saddle/front field: `canonical_profile_front_visible`, frame 12, raw pixels
  `(207,72) (187,88) (175,118) (169,158) (170,198) (178,232)
  (192,260) (202,276)`;
- dominant rear leaf: `canonical_profile_rear_leaf_visible`, frame 12, raw
  pixels `(354,66) (377,75) (395,105) (408,145) (423,202) (442,263)
  (423,273) (402,255) (384,200) (369,130)`;
- opposite-side leaf cross-check: `canonical_profile_opposite_check`, frame
  26.  It validates unequal curvature and termination; it must not be averaged
  into the frame-12 leaf.

The compact saddle/close field occupies only the source-supported
`.77-.85 Wh` profile band.  Independent rear hair supplies `.36-.38 Wh` of
additional overhang.  Complete profile targets remain `1.14 +/- .05 Wh` on
side A and `1.19-.23 Wh` on side B with `.06 Wh` uncertainty.  The difference
is retained as visible asymmetry, not forced into one symmetric depth.

### E. Rear union and leaf courses

Frame 18 controls the outer rear union.  Use its independent scanline reaches,
not a centered oval:

| fraction down rear height | left reach | right reach | full width |
| ---: | ---: | ---: | ---: |
| .10 | .074 | .189 | .262 |
| .20 | .172 | .250 | .422 |
| .30 | .270 | .291 | .561 |
| .40 | .328 | .361 | .689 |
| .50 | .406 | .414 | .820 |
| .60 | .463 | .447 | .910 |
| .70 | .467 | .455 | .922 |
| .80 | .451 | .455 | .906 |
| .90 | .225 | .463 | .689 |

The union begins narrow below the bow, is widest around `.70 +/- .04` of its
height, and terminates in unequal lobes.  The dominant foreground leaf is
`.45-.50 Wh` wide, offset about `.10 Wh` screen-right in the fixed rear view,
and crosses the centerline.  Its long diagonal edge must move laterally
`.15-.25 Wh` from crown to lower tip.  The background pocket supplies the
other broad lobes; it is not a rectangular curtain.  Together they target
`.94 +/- .04 Wh` combined rear width, `1.16 +/- .05 Wh` height, and
`.10-.14 Wh` lower-edge relief.

Use the four frame-19 `canonical_rear_visible_panel_edges` polylines only as
overlap witnesses.  A74 has two rear pockets, so it need not reproduce four
literal seams.  It must reproduce the source-visible long diagonal overlap
and at least two brown-on-brown T-junction/free-tip witnesses across the
rear-three-quarter arc.  This is a deliberate smaller claim than inventing
five equal paddles.

## Geometry representation

### Visible module

| New object | Construction | Owns |
| --- | --- | --- |
| `A74_CrownTemple_Saddle` | One open low-density subdivision saddle, locally doubled only for thin fabric thickness | crown, upper front field, asymmetric hairline, upper temples, root cover |
| `A74_Left_Cheek_Lock` | independent lightly filled two-skin pocket | left lower front/side silhouette |
| `A74_Right_Cheek_Lock` | independent non-mirrored pocket | right lower front/side silhouette |
| `A74_Rear_Background_Pocket` | broad shallow pocket with 3 unequal lower lobes | continuous rear coverage and subordinate tips |
| `A74_Rear_Foreground_Leaf` | broad off-centre shallow pocket | diagonal overlap, main free overhang, deepest profile point |

The saddle is not a closed cap.  Its rear/root edge stays beneath the frozen
bow and the two rear pockets.  Its side returns and hairline are source-visible
free edges.  Each rear pocket shares only a short crown-root contact course;
the pockets keep independent tangents and free side/bottom boundaries.

Build the saddle from a sparse 3D patch cage whose boundary control vertices
are constrained in the relevant fixed camera projections.  Interior vertices
exist only to create broad low-frequency curvature and contact.  Do not sweep
one 2D contour into another, extrude camera-plane cards, bridge artificial
inner/outer skins across visible space, or close the whole brown envelope.
Before subdivision, no straight boundary segment may span more than `.12 Wh`;
after subdivision, do not add noise or high-frequency wrinkles.

### Subordinate support

`A74_Subordinate_Face_Support` is permitted only to carry the beige face and
provide contact.  It must remain behind all brown outer contours and may be
visible only through the canonical `0.603 +/- .03 Wh` wide and high face
exposure.  Its closure, rear material, and seam are unobserved and excluded
from acceptance.  A containment mask must report zero beige/support pixels
outside the face opening in front, both profiles, rear, and both 3Q views.

If keeping the support hidden requires widening, thickening, or extending the
saddle into a helmet, A74 is invalid and must stop at P0.  Modify the support,
not the source-controlled visible field.

## Explicit parent object boundary

Start from the protected rung-003 bytes and hide, by exact name, the same
audited 15-object head/hair boundary from A73:

`Head_Cushion_Manual_Target`,
`A44 continuous hair cap with smooth opening`,
`A44 left temple fringe panel`,
`A44 left temple transition panel`,
`A44 off-center main bang panel`,
`A44 right swept fringe panel`,
`A44 right temple transition panel`,
`A45 left tapered flexible cheek lock`,
`A45 right tapered flexible cheek lock`,
`A42 Left asymmetric rear lock`,
`A42 Off-center main rear lock`,
`A42 Short right rear lock`,
`A42 Main lock left seated seam`,
`A42 Main lock right seated seam`, and
`Subtle crown center seam`.

Do not delete or rename them.  Do not hide a collection or partial-name
match.  Keep the seven exact face witnesses and complete bow visible and byte/
transform unchanged.  The face objects are comparison witnesses, not approved
placement; moving them would confound the module verdict.  Preserve all 41
legacy exclusions already hidden in the parent.

This boundary is larger than the five-object A74 module because leaving any
legacy cap, fringe, temple turn, lock, or rear leaf visible would allow rejected
geometry to own the pixels and make the test uninterpretable.

## P0 gate: front saddle, locks, and support

P0 builds only the subordinate face support, saddle, and two cheek locks.  It
is a **representation veto**, not a complete head and not an approval artifact.
Keep the frozen bow and face witnesses for scale/context.  Render front, both
front three-quarters, and both profiles; also render rear as a support-leak
mask, not as a likeness score.

Pass P0 only when all are true:

1. Front outer scanlines and all seven hairline landmarks are within
   `.03 Wh`; H3 is `(u,v)=(.588,.677) +/- .03` and stays blunt/broad.
2. Visible beige exposure is `.603 +/- .03 Wh` in width and height, with zero
   support pixels outside the face opening and no parallel mask/plaque rim.
3. The saddle reads as a close, thin fabric field: contact/free-edge offset
   `.01-.03 Wh`, no root gap over `.01 Wh`, no visor over `.05 Wh`, and no
   second inflated crown silhouette.
4. Both 3Q views preserve continuous crown-to-temple seating and front layer
   order.  No card, hard crease band, artificial bridge, pinhole, triangular
   daylight, or mirror-symmetric bang cadence is visible.
5. Profile field depth remains inside `.77-.85 Wh` and the cheek locks remain
   thin independent panels; neither becomes a side wall.
6. Implementation-blind neutral-clay scores are at least `6/10` for front
   identity, constructed-plush read, and contact/occlusion, with no helmet,
   mask, box, slab, or egg veto.

Reset immediately if the saddle needs a closed rear wall to look coherent, if
the support is exposed, or if either 3Q view resembles A34's cap/cards.  Do not
add rear leaves to conceal a P0 failure.

## P1 gate: complete minimal visible envelope

P1 adds only the background rear pocket and foreground leaf.  Render fixed
front, both profiles, rear, both front 3Q, both rear 3Q, plus one uncropped
whole-subject presentation.  Every view must be side-by-side with its
controlling source and include a registered silhouette overlay/edge
difference.  The full 30-frame turn remains the transition gate; selected
frames are witnesses, not replacements.

Pass P1 only when all are true:

1. P0 front landmarks do not regress beyond `.03 Wh`.
2. Complete profile depth is `1.14 +/- .05 Wh` on side A and `1.19-.23 Wh`
   on side B within stated uncertainty, while the close saddle remains
   `.77-.85 Wh`; deepest reach belongs to the independent leaf.
3. Rear union is `.94 +/- .04 Wh` wide and `1.16 +/- .05 Wh` high; maximum
   width lies at `.70 +/- .04` of rear height; lower relief is `.10-.14 Wh`.
4. The broad foreground leaf is `.45-.50 Wh`, visibly off-centre, crosses the
   rear centerline, and creates the long diagonal overlap.  The two rear
   pockets create at least two readable T-junction/free-tip separations without
   becoming equal paddles, a symmetric split, or a full-width cape.
5. Brown coverage is continuous at crown, bow root, temples, and rear.  There
   is no beige/bald leak, exposed root tab, air channel, clipping, floating
   edge, or accidental tangency above `.01 Wh`.
6. Both profile and both 3Q arcs preserve local order: fringe/face, cheek lock,
   close saddle, rear pockets, and locally crossing bow tail.  The bow may not
   be forced into one global depth plane.
7. Whole-subject implementation-blind scores are at least `6/10` for macro
   silhouette, construction, identity, and contact.  The coupon must be
   preferred to protected rung 003 without any new major front-identity
   failure.

P1 is still an internal rung, not final sculpt approval.  The global 8/10
visual-quality gate remains unchanged for an approval candidate.

## Builder order and early returns

1. Copy the observed curves and hashes into an A74 machine-readable contract;
   static preflight verifies every value against the existing records.
2. Create the exact visibility boundary in a disposable copy and verify the
   protected source hash before and after.
3. Build only the subordinate support and P0 three brown objects.  Render
   512--640 px front plus the two regression-risk 3Q views first.
4. If those pixels fail, return early with contact sheet and mask evidence.
   Do not spend time on rear leaves, fibers, materials, seams, rigging, or bow
   reseating.
5. Only after P0 passes, add the two P1 rear pockets and run all fixed views
   plus the 30-frame transition review.
6. Never promote an isolated object.  Promotion requires the whole-subject P1
   packet because likeness depends on occlusion and layer order.

## Principal regression risks

- **Helmet by scope creep:** extending the saddle rearward to hide support
  converts it into the rejected shell.  Keep the saddle open and subordinate
  support contained.
- **Card relapse:** using source curves only in the front plane, then extruding
  them, repeats A34/A46.  Constrain the same sparse cage in front and both 3Q/
  profile cameras before adding thickness.
- **Rear cape/paddles:** two pockets must overlap asymmetrically and terminate
  independently; they may not share a straight bottom or equal widths.
- **Legacy-pixel contamination:** the whole audited 15-object boundary must be
  hidden by exact name even though A74 replaces it with only six new objects
  including the support.
- **Measurement theater:** front/rear union pixels do not reveal hidden seams.
  Keep every hidden root labeled as a reversible hypothesis and judge the
  visible result absolutely.

This plan is smaller than another full head rebuild but large enough to test
the actual visible construction.  Its first useful artifact should be a P0
whole-context render packet, not another isolated gray receiver.
