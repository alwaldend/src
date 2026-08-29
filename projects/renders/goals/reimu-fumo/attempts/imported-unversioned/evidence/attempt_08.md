# Reimu Fumo attempt 8

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 8 — layered cushion and gathered continuous ribbon

**Candidates:** `out/reimu_fumo_attempt_008_head.blend`, SHA-256
`01ce5141f171550c74c18af801dde3f90e60389aacc6257d43bacf4c9fee4dc0`,
and `out/reimu_fumo_attempt_008_right_loop.blend`, SHA-256
`367e1c029dfcc440107edc5017bcb62d6cd4ae05cb9b6b912f9dc4535142858e`,
neutral construction stage, review packet
`attempt-008-layered-head-cloth-bow`.

**Failure targeted:** Attempt 7's exact outer contours were undermined by an
incorrect layer stack. A brown structural shell around a beige material region
created a helmet and inset mask. Separate visible and return bow panels plus
`0.295–0.305 Wh` fold bridges created rigid bars and a total loop depth of
about `0.38–0.39 Wh`, while the reference bow core is only
`0.03–0.05 Wh` thick.

**Hypothesis:** An exposed beige front/rear/gusset cushion covered by thin,
physically separate brown fleece panels will read as a sewn head with layered
hair. A single pre-folded U-shaped ribbon per loop, relaxed only at a gathered
root with Blender cloth sewing springs, will preserve the measured front
silhouette and remain plausibly thin in side view.

**Plan written before implementation:**

1. Start from factory-empty data and copy no rejected render mesh. Import the
   frozen Attempt 7 front and side outlines only as non-rendering control cages.
2. Build the structural head as one beige all-quad cushion with distinct front,
   rear, and perimeter-gusset regions. Keep it at
   `1.00 × 0.74 × 1.03 Wh`; limit central-face convexity to
   `0.015–0.025 Wh` and roll only the outer `0.08–0.12 Wh` into the cheeks.
3. Build crown/front, side-wrap, rear, fringe, cheek-lock, and nape hair as
   separate `0.015–0.025 Wh` felt layers whose outside follows the accepted
   cage. Keep the beige cushion beneath them and inset covered cushion regions
   by `0.010–0.015 Wh`. The hairline is a physical rounded cut edge, not a
   material boundary or raised curve.
4. Preserve the frozen hairline anchors but round its central point across
   `0.04–0.05 Wh`. Keep each cheek lock `0.16 ± 0.02 Wh` wide, with a blunt
   end and no part outside `x = ±0.50 Wh`. Add a separate rear panel about
   `0.99 × 1.12 Wh`, reaching `0.79 ± 0.04 Wh` total depth and ending in
   three to five shallow irregular nape points. Encode seams as
   `0.004–0.008 Wh` valleys across `0.012–0.020 Wh` compression bands; use no
   piping.
5. Prove one right bow loop before building the remaining bow. Use one
   continuous U-folded `35 × 11` all-quad surface: 16 stations from root to
   the visible outer edge, four stations around a `0.013 Wh` rolled fold, and
   15 return stations back to an open root. Keep visible/return separation at
   `0.026 Wh`, resample the measured front trace, and keep its outer fold broad.
6. Seat the loop behind the rear crown at its evaluated surface plus
   `0.008 Wh`. Add three isolated pinned anchors and loose sewing edges from
   both open loop-root boundaries to the anchors. Author three shallow root
   pleats below `0.009 Wh` amplitude before simulation; cloth may relax and
   gather them but may not generate the controlling outline.
7. Simulate frames 1–80 without pressure, with low gravity, collision against
   one smooth head proxy, self-collision, and `sewing_force_max` ramped from
   `0.5` at frame 1 to `3.0` at frame 15 and `8.0` from frame 30. Use quality
   12, mass `0.03`, tension/compression/shear stiffness `40/25/15`, bending
   stiffness `1.2`, and corresponding damping `7/7/7/2`. Extract the first of
   frames 70, 75, or 80 whose five-frame RMS movement is below `0.001 Wh`.
8. Reject the isolated loop at frame 1 if the measured front trace is wrong,
   the fold is pointed, side depth exceeds `0.05 Wh`, or the root is not seated.
   Reject at frame 30 if gathering produces spikes, changes front extrema by
   more than `0.03 Wh`, crosses the head, or collapses into a side-view bar.
9. Only after the right loop passes, build the left loop and both tails in the
   same continuous-ribbon language. Keep tails behind the head, with only
   `0.010–0.015 Wh` authored twist. Add a thin rectangular center wrap after
   simulation, with at least 60% hidden behind the crown; do not add a rounded
   knot cushion.
10. Render fixed front, side, rear, three-quarter, grazing-close-up, overlay,
    and black-silhouette evidence. Require head/hair width `1.00 ± 0.03 Wh`,
    cushion depth `0.74 Wh`, outer hair depth `0.79 ± 0.04 Wh`, bow span
    `1.64–1.74 Wh`, bow vertical extent `0.95–1.05 Wh`, and bow core depth
    `0.03–0.05 Wh`. Reject before body work unless every component score is at
    least 8/10 and no helmet, mask, piping, tusk, rabbit-ear, petal, fin, wing,
    pill, slab, floating-root, or tangent-root read remains.

**Work performed:** Built the head and right-loop proofs in separate
factory-empty Blender files so failure in one could not contaminate the other.
The first hair version used independent crown, side, and rear panels and was
immediately rejected as rigid armor. The second copied the evaluated beige
cushion into one thin open conformal hood, with separate conforming fringe,
cheek locks, and nape panels. It corrected material assignment and produced
fixed front, side, rear, three-quarter, and grazing-close-up renders.

The first continuous U-fold loop was allowed to simulate freely and lost about
43% of its X span by frame 30 and 67% by frame 80. The second constrained the
broad lobe but remained unstable and reached `0.139 Wh` total root depth. The
third kept the broad lobe deterministic and confined cloth motion to its root
zone. It passed the numeric stability, silhouette-extrema, and folded-core
depth gates. Only that third loop entered the final Attempt 8 evidence. No
candidate entered the tracked asset.

**Evidence:** The final head/hair width is `0.9979 Wh`; cushion depth is
`0.7400 Wh`; outer hair depth is `0.7715 Wh`; cheek-lock widths are
`0.1680/0.1641 Wh`; and rear hair is `0.9835 × 1.1226 Wh`. Every frozen local
numeric gate passes. Read-only Blender MCP inspection found one closed
522-vertex head cage with no boundary, non-manifold, or degenerate geometry.
The hair surfaces are deliberately open before their Solidify modifiers and
have no degenerate faces.

The final loop stabilizes at frame 75 with five-frame RMS movements
`0.000233 Wh` and `0.000221 Wh`. Its X extent is `0.052–0.737 Wh`, Z extent is
`1.487–2.103 Wh`, and median local folded-core depth is `0.032 Wh`. Read-only
Blender MCP inspection confirms 340 non-degenerate base faces, a separate
closed collision proxy, no pressure, and Solidify/Subdivision only after the
extracted surface.

The implementation-blind pixels-only reviewer did not recognize either proof
as the reference subject. Head scores were likeness 2/10, hair likeness 2/10,
silhouette 2/10, sewn construction 1/10, intended-medium read 1/10, contact
and occlusion 3/10, and presentation 6/10. Loop scores were likeness 2/10,
silhouette 3/10, folded/sewn construction 1/10, intended-medium read 1/10,
attachment 1/10, and presentation 5/10. The reviewer reported helmet, mask,
shell, flat-card, fin, wing, petal, slab, floating/tangent-root, and molded-toy
reads.

**Criterion results:** All targeted scalar bounds and the loop's final motion
stability pass, proving that contour, depth, and simulation can be controlled.
Reference likeness, complete measured silhouette, plush construction, and
presentation quality fail the absolute image gate by large margins. Reusable
structure, animation readiness, and final technical integrity remain
unverified. Repository delivery still applies only to the rejected migrated
baseline.

**Decision:** Reject both isolated proofs. Do not reuse the hood, fringe,
locks, nape, cloth simulation, or final loop surface. Retain only the frozen
traces, head control envelope, and evidence that a deterministic continuous
U-fold can satisfy the thickness and stability constraints.

### Progress and approach audit after attempt 8

- **Improved mechanically:** every targeted head/hair scalar bound passed.
  The third loop changed an unstable collapsing cloth sheet into a stable
  `0.032 Wh` folded core with correct extrema. Background Blender MCP execution
  also reduced the component feedback loop to minutes rather than requiring a
  complete character rebuild.
- **Regressed visually:** the independent hair stack became armor; the
  replacement conformal hood became a monolithic helmet and mask. The final
  loop was thinner and more stable but remained an upright teardrop petal with
  a blade root. Absolute image scores fell to 1–3/10.
- **Absolute result:** neither component is a viable modeling base. Numeric
  compliance again failed to predict likeness or constructed-plush read.
- **Head representation failure:** putting a complete beige cushion under a
  physically separate brown shell still creates the same shell-and-opening
  image as a material inset. The next head must expose the beige front as the
  complete structural front panel and use brown only for the rear/perimeter
  gusset plus actual overlaid fringe, locks, and nape.
- **Bow representation failure:** cloth simulation is unsuitable for the broad
  identity silhouette. Even localized simulation adds cost without producing
  the reference folds. Preserve continuous U-fold topology and core depth, but
  sculpt the broad lobe and its returned fold directly. Use a separate short
  gathered root collar whose folds are explicitly authored and seated.
- **Highest-leverage problem:** prove a non-box structural head in which seam
  flow and panel colors explain the depth, and prove one laterally spread bow
  lobe with an unmistakable returned fold and broad gathered root.
- **Next approach:** build both macro forms from low-resolution sculptable
  cages, use measured outlines as constraints rather than rendered shells,
  and put the reference-critical white turned trim into the bow macro gate so
  the physical folded edge can be judged rather than inferred.
