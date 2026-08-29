# Reimu Fumo attempt 9

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 9 — structural head panels and art-directed folded ribbon

**Candidates:** `out/reimu_fumo_attempt_009_head.blend`, SHA-256
`9161f289277ec3e76fd7090c102f01206c7ae7b90b937e7338a4ee46c5550f8f`,
and `out/reimu_fumo_attempt_009_right_loop.blend`, SHA-256
`14d14faabd03739f0c7c0f415d96a6e27c44c2fb2474839ecdf9a73db39277fd`,
neutral construction stage, review packet
`attempt-009-structural-head-sculpted-bow`.

**Failure targeted:** Attempt 8 proved that separate visible layers and exact
dimensions do not prevent a helmet if one complete cushion is covered by one
complete hood. It also proved that cloth can satisfy stability and thickness
while preserving a petal silhouette. The representation must make the face a
true structural panel and must author the bow's fabric gesture directly.

**Hypothesis:** One manifold cushion with a full beige front panel, brown
perimeter gusset, and brown rear panel will explain its own soft depth without
a mask or shell. Thin hair pieces copied from that evaluated cushion will then
read as overlays. A low-resolution continuous U-fold whose broad fields,
radiating valleys, returned lip, gathered collar, and silhouette-changing
white trim are sculpted explicitly will read as ribbon rather than a card.

**Plan written before implementation:**

1. Start both isolated proofs from factory-empty Blender data and copy no
   rejected render mesh. Retain only the frozen front/side traces, cameras,
   and verified numeric constraints as non-rendering controls.
2. Build one manifold structural head with matching `11 × 11` front and rear
   all-quad panels, 40-vertex perimeters, and seven intermediate gusset rings.
   Assign beige to the complete front panel and brown to the complete gusset
   and rear. Use no hair hood, face insert, brown front rim, or added seam
   object.
3. Keep structural width `0.96 Wh`, height `1.03 Wh`, and depth `0.74 Wh`.
   Map the panels through the frozen front trace and every gusset row through
   the frozen side trace. Use paired seam-compression/recovery rings at depth
   interpolations `0.035/0.090` and `0.910/0.965`, with only
   `0.004–0.006 Wh` seam valleys. Turn the lower front seam beneath the chin
   so no brown horizontal chin rim is visible.
4. On one `StuffingNeutral` shape key, use broad controlled sculpt-like cage
   displacements only: flatten the central `0.60 Wh` face to
   `0.015–0.025 Wh` depth variation; add paired cheek fullness below
   `0.015 Wh`; compress crown, lower rear, and seam bands within their stated
   limits; and add at most `0.012 Wh` asymmetry. Preserve broad front and rear
   planes and reproject every silhouette extremum to the frozen control cage.
5. Run the structural-head-only pixel gate before hair. Require no straight
   silhouette run over `0.30 Wh`, a readable front-plane/seam/gusset/rear-plane
   transition in side view, exact depth, and no rounded-box, balloon, helmet,
   mask, or chin-rim read.
6. Only after that gate, duplicate evaluated cushion patches for a thin
   conforming fringe, cheek locks, and nape. Keep fringe projection at most
   `0.015 Wh`, round the central tip across `0.04–0.05 Wh`, keep lock widths
   `0.16 ± 0.02 Wh` with buried roots and blunt ends, and expose only the free
   lower nape points. Add no crown/rear hood or planar polygon card.
7. Build one right bow loop as a continuous `37 × 11` all-quad U-fold: 16
   visible stations, five stations around a `0.012–0.015 Wh` semicircular
   outer turn, and 16 return stations. Keep visible/return separation
   `0.024–0.030 Wh`, inset the return by only `0.008–0.012 Wh`, and use one
   subdivision level after shaping. Use no cloth, pressure, separate return,
   bridge, or inflated pocket.
8. Reject its unsmoothed silhouette before fold work unless it spans
   `x=0.05–0.70 Wh`, reaches `z=1.48–2.11 Wh`, has at least `0.28 Wh` of
   readable outer vertical edge, carries a `0.12–0.16 Wh` broad outer
   shoulder, and spreads laterally as a kite/trapezoid instead of a teardrop.
9. Shape `Bow_Default_Physical` with a `5 × 3 × 5` lattice and localized
   sculpt-like shape-key displacements. Create three slightly planar fields,
   two `0.015–0.022 Wh` radiating valleys that begin at the root and fade
   before the turn, one weaker upper compression fold, and an outer edge rolled
   rearward `0.008–0.012 Wh`. Do not globally smooth afterward.
10. Join the broad loop to a separate deterministic six-ring by nine-column
    gathered collar `0.10–0.13 Wh` long, narrowing from about `0.24 Wh` to
    `0.09–0.12 Wh` and carrying three `0.010–0.015 Wh` accordion offsets.
    Seat it inside a compressed rectangular wrap; add no rounded knot cushion.
11. Add simplified white turned trim to the exposed loop boundary in this
    macro gate. Give it `0.030–0.040 Wh` outward width, `0.007–0.010 Wh`
    thickness, and 12–16 broad gathers. Require the red core to remain
    `0.03–0.05 Wh` thick and the complete turned edge to reach
    `0.07–0.10 Wh` in side view.
12. Render each earliest gate plus final fixed front, side, rear,
    three-quarter, grazing-close-up, overlay, and silhouettes. Reject before
    the other loop, tails, or body unless every applicable blind score reaches
    8/10 and no box, balloon, helmet, mask, chin rim, shell, card, petal, fin,
    wing, slit, bar, slab, piping, floating root, or tangent root remains.

**Work performed:** Built both proofs from factory-empty data. The head used
one shared-vertex manifold with a complete beige front panel, brown perimeter
gusset, brown rear, compressed shared seam rings, and one low-frequency
`StuffingNeutral` shape key. Gate A rendered that structural cushion before
hair. Gate B then built a front-surface-derived fringe. Five root experiments
were pixel-checked; radial interpolation caused crossings, a crown-row bridge
caused a flap, and a normal offset caused a breakthrough. The final fixed-X
front parameterization removed those defects. Gate C added curved/tapered
cheek locks and a nape whose root rows were hidden behind the rear panel.

The bow used a continuous `37 × 11` visible/turn/return mesh, one
`Bow_Default_Physical` shape key, a `5 × 3 × 5` lattice, a closed deterministic
gather collar, compressed wrap, and simplified 14-gather trim. Gate A rendered
the untrimmed silhouette, Gate B enabled folds and collar, and Gate C enabled
the turned trim. Every stage was retained as temporary evidence only; no
candidate entered the tracked asset.

**Evidence:** The structural head evaluates to exactly
`0.960 × 0.740 × 1.030 Wh`. The final fringe projection is `0.013499 Wh`; lock
widths are about `0.149 Wh`, their inner edges are at approximately
`x = ±0.335 Wh`, and their outer edges remain inside `x = ±0.484 Wh`. The
nape exposes only its lower points. Read-only Blender MCP inspection found all
five evaluated meshes closed, with zero boundary edges, non-manifold edges,
or degenerate faces.

The two independent pixels-only head reviews both rejected the candidate.
The stricter absolute review scored head likeness 4/10, hair likeness 3/10,
silhouette 4/10, sewn construction 2/10, intended-medium read 2/10, contact
and occlusion 3/10, and presentation 7/10. A construction-focused review
scored silhouette 5/10, sewn construction 3/10, intended-medium read 3/10,
hair framing 4/10, and contact 2/10. Both reported a dominant rounded rear
block, beige crown-to-cheek crevice, attached fringe plate, exposed lock-root
ledges, slab-like locks, shallow beige side presence, and floating nape tips.

The bow's red core is `0.0315 Wh` thick and its trim brings the turned edge to
`0.073 Wh`; all bounds, topology, and trim-frequency gates pass. Read-only
Blender MCP inspection found every evaluated bow mesh closed with zero
boundary, non-manifold, or degenerate geometry. The internal pixels-only review
still scored the loop 3/10 and rejected at Gate A. The loop remained a smooth
upright petal/card; configured valleys were not legible; the collar became a
rectangular tab; and trim became a front-view zigzag and side-view slit.

**Criterion results:** The structural color/layer change raised isolated head
likeness from 2/10 to 4/10 and removed the explicit chin rim. Exact bounds,
manifold integrity, and several local attachment measurements pass. Reference
likeness, complete silhouette, plush construction, and presentation quality
still fail the absolute image gate. Reusable structure, animation readiness,
and final technical integrity remain unverified. Repository delivery still
applies only to the rejected migrated baseline.

**Decision:** Reject both rendered candidates. Retain the evidence that a full
beige front panel is better than a mask or hood and that the nape root can be
hidden, but copy no Attempt 9 render mesh. The head and bow both require new
occupied-volume scaffolds and new panel flow.

### Progress and approach audit after attempt 9

- **Improved:** head likeness doubled from 2/10 to 4/10. The front view first
  read as a structural two-material cushion rather than an inset face, the chin
  rim disappeared, the final fringe had no breakthrough, and the nape band was
  reduced to lower points. All evaluated meshes were technically closed and
  clean.
- **Regressed or unchanged:** no visual score reached 8/10. The head's side and
  three-quarter views still looked like a deep brown rounded block carrying a
  thin face plate. Locks remained slabs, and the bow repeated the petal/card
  read despite a different fold and trim implementation.
- **Absolute result:** the head representation is directionally better but not
  a viable approval checkpoint. The bow representation is exhausted after two
  consecutive continuous-U attempts.
- **Head representation failure:** assigning the complete front surface beige
  did not give that front panel visible side ownership. The brown gusset still
  occupied nearly all depth. Separate conformal hair surfaces remained too
  thin to establish padded felt pieces and could not conceal their hard roots.
- **Bow representation failure:** starting from any parametric sheet or U-grid
  makes the broad lobe visually average into a smooth sheet. Adding folds,
  collars, or trim afterward cannot create the missing occupied fold volume.
- **Highest-leverage problem:** sculpt the three-dimensional space occupied by
  the head's face-side-rear stuffing and bow's broad folds before deriving any
  visible fabric surface. Panel seams must follow that accepted volume rather
  than define it indirectly.
- **Next approach:** rebuild the head so the beige face rolls around
  `0.12–0.18 Wh` into each side before a narrower gusset and shallow rear
  plane. Build hair as lightly padded front/back panels with shared or deeply
  overlapped roots. For the bow, sculpt a hidden fold-volume wedge, then
  retopologize and weld explicit front, turned-edge, and back panels around it
  with radial pleats inside the joined root.
