# Reimu Fumo references and modeling contract

[Back to goal](README.md)

## Reference dossier

The current private attachment root is
`/var/home/simeonwarrenbot/.t3/userdata/attachments/`. Every source filename
begins with `136a19c9-b38c-4c43-959e-141b9e2224e3-`. The suffixes below plus
their hashes are the durable source locators for this workspace.

- **Canonical 25 cm reference PNG
  `82f07f2f-0207-4723-8a9b-d76b2e7e0d49.png`:** SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
  This is the primary authority for the exact Reimu Fumo variant, overall
  proportions, front silhouette, hair and bow design, expression, clothing,
  seated pose, and real-world scale. The complete plush is `0.25 m` tall.
- **Canonical slow-180 GIF:** user-supplied URL
  `https://c.tenor.com/wuTstMILarIAAAAC/tenor.gif`; verified working copy
  `out/reimu_fumo_reference_set_v2/canonical_turn_180.gif`; SHA-256
  `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`.
  It contains 30 coalesced `498 × 498` frames at `10 cs` per frame. Use all
  frames as the primary source for three-quarter, side, rear, depth order,
  seated volume, hair layering, bow drape, skirt pooling, and foot occlusion.
  The uploaded `da8f04ae...png`, SHA-256
  `13cb3e91c44b2e919c2fe5d1b3adf4daf02dccaefe8a9b257dc25be59528bd43`,
  is a flattened preview of this animation and is not the controlling motion
  source.

- **Clean front PNG `04ac3273-048f-4ee0-a605-12a7edd4c7bf.png`:**
  SHA-256
  `37813e03e04e4966f1dbe914e03a25a5f5ae561dcbf58b72677195c513ea48ca`.
  Primary front identity and neutral silhouette; isolated presentation reduces
  contact evidence.
- **Physical front PNG `c690c3e5-d072-4c0a-bda5-7d452a501519.png`:**
  SHA-256
  `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51`.
  Primary real-fabric, stuffing, front proportion, and seated-contact evidence.
- **Physical side PNG `5a8e0eba-24f7-4d02-97db-5608d16966f9.png`:**
  SHA-256
  `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`.
  Primary depth, rear hair, bow layering, skirt pooling, and leg/foot
  construction evidence; perspective is not orthographic.
- **Turn GIF `f8df027f-9508-4b59-8779-21310993ccfa.gif`:**
  SHA-256
  `b42368e921bd055d73fbbb7bf65c2509a9aaf190cab02f89824b92b4cb75ece4`.
  Use all four coalesced frames: `frame_000` three-quarter, `frame_001` front,
  `frame_002` side, and `frame_003` rear. Low resolution and a different
  physical pose limit exact measurements.
- **Sofa GIF `631a9796-a3a8-487c-baa6-89da2eb26598.gif`:**
  SHA-256
  `7c9173f91e6b6c801a1c77e50f9135e86fc89319f3c0262c10312320b1af8589`.
  Use `frame_000`, `frame_040`, and `frame_050` for seated silhouette;
  `frame_020`, `frame_030`, `frame_060`, and `frame_070` for fabric panels; and
  `frame_080` for the face. Motion blur, occlusion, pan/zoom, and sofa
  perspective limit silhouette measurements.

The new canonical 25 cm reference controls variant identity, overall front
proportions, and final scale. The verified 30-frame slow-180 GIF controls
three-quarter, side, rear, and depth construction. The older physical front and side
stills retain supporting authority for fabric construction and occluded depth;
the clean front retains supporting graphic landmarks. Older GIFs control only
otherwise missing side/rear relationships. There is still no high-resolution
true rear frame in the received files, so rear fidelity must be reviewed with
explicitly lower measurement confidence rather than claimed at front-view
precision.

### Multi-reference use rule

No candidate module may be accepted from one reference merely because that
view has the cleanest silhouette. Before modeling, assign every relevant
supplied still and GIF frame a documented role: controlling measurement,
silhouette, depth/contact, seam/overlap, fabric/pile, or qualitative
cross-check. Use all relevant views in those roles and state uncertainty or
contradictions instead of averaging them into a fictional shape.

The canonical pair must be present in every attempt's durable reference packet.
Disposable ignored working blends may instead record the tracked packet path
and source hashes in their manifest so fixed-view copies remain byte-exact for
A/B review. A final promoted blend must include packed reference image
datablocks in a disabled `REFERENCE_ONLY` collection outside the reusable
asset collection. The canonical 25 cm front wins direct conflicts about
variant identity, proportions, and scale; the 180 GIF wins direct conflicts
about depth and layer order. Older references remain required supporting
evidence and uncertainty checks.

For the active front/top hair field, the clean and physical front stills
support identity and cross-check the asymmetric three-span fringe; the
physical front also supports pile, edge softness, and real stuffing. The
physical side supports crown/rear depth, bow occlusion, temple return, and
layered felt construction. All four older turn-GIF orientations cross-check
front-to-side-to-rear continuity and the rear lobe order, at lower dimensional
confidence. The selected sofa-GIF frames support seated felt behavior,
cheek-layer overlap, edge thickness, and soft contact. The canonical front and
30-frame turn remain the controlling pair. Rear views also act as a regression
that the accepted rear-panel V4 contour remains unchanged. A role marked
qualitative cannot override a higher-confidence measured view, but it can veto
an impossible transition or material construction.

The earlier source-trace dossier's raw O/F pixels remain valid for the
supporting physical-front variant, but its derived `Wh = 178 px` ratios do not.
Normalize that physical-front evidence at `Wh = 189 px`: its `90 px`
center-root span is `.476 Wh` and its roughly `52 px` center depth is
`.275 Wh`, not `.506 Wh` and `.292 Wh`. Those values support fabric and panel
construction only; they do not override or get averaged into the canonical
exact-variant registration below.

## Modeling strategy contract

### Current canonical authority after Attempt 36

The canonical exact-variant datum is `Wh = 368 ± 4 px`, center `x=485`, crown
top `y=231`, and central fringe tip `(.588,.677)` with `.101 Wh` viewer-right
offset. The `.603 × .603 Wh` measurement describes only the visible beige
pixel mask; it is not evidence for the complete or hidden face opening.
`1.098 Wh` is the lock-inclusive crown-to-lowest-lock composite regression.
Attempt 36 excluded locks and ended at `v=.990`, so it could not use
`1.098 Wh` as its own cage height. Turntable depth bands reserve
`.35 ± .05 Wh` for the independent rear leaf; the rejected Attempt 36
execution contract's `.25 Wh` is superseded.

The `Wh=189`, `1.03 Wh`, and tip `(.49,.57)` entries in the physical-front
landmark table remain useful supporting-variant evidence for fabric,
construction, and uncertainty checks. They do not control the canonical exact
variant and must not be averaged with its values.

Attempt 36's explicit `13 × 9` open cage is rejected at W0. No annulus,
separate beige insert, polar field, homothetic courses, analytic drum,
rectangular card, automatic rear disk, or replacement H-grid is authorized.
Attempts 45–46 subsequently proved that scalar-correct analytic shells and
weakly warped front/rear cards still render as helmets or open boxes. The live
strategy is the dependency-ordered working ladder in
[current_attempt.md](current_attempt.md), with the canonical turntable also
being tested as a local visual-hull guide rather than final geometry.

### Historical representation notes — subordinate to current authority

- **Authoritative front-hair interface reset:** V4 proves that a front-only
  offset crescent cannot own the source-visible crown; coupon V1 proves a
  narrow raised strip reads as a floating rail; and coupon V2 proves an
  arbitrary finite crown crop reads as a floating card whose artificial cut
  edges dominate. A new source-only audit finds no seam, root, overlap, depth
  step, or pile break enclosing the central bang. The smallest defensible
  module is therefore one continuous field spanning the visible crown,
  temple/root wrap, and complete asymmetric three-span fringe edge. Its unseen
  upper-rear closure is an inferred manufacturing interface, not a photographed
  seam. Cheek and nape locks remain later separate panels because their overlap
  is visible. `IFACE_Front_Hair_Field` remains registration only.
- **Attempt 27 hair-base clarification:** implementation review proves that a
  front-edge-to-upper-rear rectangle still requires two artificial crown-side
  cuts. Source review instead authorizes one continuous _visible hair-base
  assembly_ owning the reliable front edge through physical point `(302,181)`,
  every visible crown/rear brown sample, and the irregular rear silhouette.
  This is a module boundary, not evidence for one literal factory pattern.
  Hidden seams, exact side correspondence, attachment depth, and rear-lower
  manufacture remain inferred. Cheek/nape locks remain separate because their
  overlap is visible.
- **Rear-camera registration clarification:** the fixed rear camera reverses
  world X on screen. Rear-frame viewer-left therefore maps to positive world X
  and viewer-right to negative world X. The audited deepest viewer-left lobe
  at rear pixel `(84,224)` is stored at global world-X column `55`, not `29`;
  this is a camera-registration correction rather than a mirrored design
  change.
- Normalize measurements to outer head/hair width `Wh = 1`. The physical-front
  `LANDMARKS.md` table remains a supporting-variant uncertainty, exclusion,
  and construction record wherever it uses `Wh=189`, `1.03 Wh`, or tip
  `(.49,.57)`; the current canonical authority above controls exact-variant
  identity. Do not average conflicting variants or alter a frozen candidate
  gate after seeing its result. Do not apply the physical-front `1.03 Wh`
  envelope to the mostly occluded beige support cushion; its provisional
  construction band is `0.84–0.90 Wh`, and likeness must be judged on the
  observable coupled form.
- The schema-2 source packet proves visible projection and layer order, not
  hidden manufacture: brown owns every visible crown/rear sample; beige is
  confined to the face opening; the front boundary is one continuous
  asymmetric edge with three unequal semantic spans; and cheek/nape layers
  overlap at hidden roots. Do not claim that the pixels prove a literal
  one-piece brown cushion, beige applique, diagnostic-loop seam, or occluded
  root closure. A shallow unseen beige support remains permissible only if it
  never appears outside the face opening.
- For the next reversible macro trial, build the coupled head volume from one
  connected brown front annulus, a nearly flush visible beige face insert
  sharing its inner loop, and one brown all-quad multi-patch rear disk sharing
  the outer loop. This is a conditionally selected modeling hypothesis, not a
  source fact. It may survive only if it closes one shallow sewn volume with no
  pole/fan, hard equator, cap/curtain split, vertical side wall, raised mask,
  exposed beige crown/rear, or constant-depth extrusion.
- Preserve the physical-front broad center low band as part of the single
  connected asymmetric fringe edge, not a narrow point or three inflated
  wedges. Gate the bare coupled volume in front and rear before adding separate
  softly padded unequal cheek and nape locks; then gate their hidden overlap
  and side/rear depth order before any material or detail work.
- Build each bow loop as a folded closed ribbon pocket with a gathered root,
  return leaf, broad outer fold, low-frequency creases, and asymmetric drape.
  Tails are twisted double-sided strips; trim follows the real panel boundary.
- Build the torso as a compact seated bean cushion without human anatomy;
  sleeves as front/back bell envelopes with a narrow gusset; the skirt as broad
  gathered front/back panels with a deeper pooled rear; and each black foot as
  a flattened sewn pod partly hidden by the ruffle.
- Keep a physically plausible seated Fumo as the canonical rest presentation.
  Use one real armature with named control, mechanism, and deform bones; deform
  high-resolution detail from the low-resolution panel cages.
- Never use an egg or sphere as the final head, a separate full hair shell, a
  lathed cone skirt, a flared tube sleeve, a flat card bow, or a final-form
  disconnected primitive.
