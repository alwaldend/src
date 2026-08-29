# A79 paired-hair topology-density bound

## Verdict

Variant A's current `9 x 7`, `9 x 9`, and `7 x 7` grids are not eligible for
construction.  Their average projected edge spans are `21--51 px` in the
controlling directions at the frozen 512 px cameras.  Smooth shading cannot
repair those silhouettes.

Use a non-uniform, support-looped Catmull-Clark level-2 cage, apply/freeze the
result, and validate the evaluated paired pocket.  This is preferable to
directly authoring a 40,000--90,000-vertex paired tensor net.  Subdivision is
only the sampling mechanism: it does not authorize a receiver-derived shell,
uniform offset, live modifier in the saved candidate, or post-hoc smoothing
of a wrong form.

The hard evaluated sampling floor is:

| Part | Major evaluated spans | Minor evaluated spans |
| --- | ---: | ---: |
| front/crown | `136` around the crown/temples | `66` through profile depth |
| rear base | `94` over height | `76` over width |
| dominant leaf | `89` along the long diagonal | `43` over width |

Those values enforce at most `3 px` per ordinary evaluated span.  Every span
touching a crown extremum, temple turn, nape-lobe tip, leaf tip, or overlap
corner must instead be at most `2 px`.  Counts alone cannot prove that local
condition; the literal candidate must be projected through all frozen cameras
and measured after refinement.

## Evidence and scale

- All fixed review cameras are orthographic with a `292 mm` vertical scale on
  a `512 x 512` render.
- Therefore `1 px = 0.5703125 mm`, `3 px = 1.7109375 mm`, and
  `2 px = 1.140625 mm` in an unforeshortened image-plane direction.
- The audited project datum is `Wh = 132 mm`.
- The front/crown arc budget is the chord sum through the six canonical-front
  outer scanlines, `221.0125 mm`, multiplied by `1.05` so the sparse chord sum
  cannot claim the unobserved curve is shorter.  The density length is
  `232.0631 mm`.
- The front/crown profile direction uses the upper compact-field bound,
  `.85 Wh = 112.2 mm`.
- The rear-base directions use the upper frozen rear bounds:
  `1.21 Wh = 159.72 mm` high and `.98 Wh = 129.36 mm` wide.
- The dominant-leaf budget includes the reported uncertainty:
  `(1.05 + .07) Wh = 147.84 mm` root-to-tip and
  `(.50 + .05) Wh = 72.6 mm` wide.  Combining the upper height with the
  observed `.25 Wh` diagonal shift gives a `151.4783 mm` long-edge budget.

These are conservative sampling lengths, not hidden construction coordinates
and not proof of likeness.  The turntable measurements remain projected outer
bounds with their documented phase uncertainty.

The reproducible arithmetic is in `density_budget.py`; its exact output is in
`density_budget.json`.

## Evaluated segment floors

For a projected chain of length `L`, the straight-span count is:

```text
N = ceil(L / (pixel_limit * 292 / 512))
```

| Chain | Length (mm) | `N` at 3 px | `N` at 2 px everywhere |
| --- | ---: | ---: | ---: |
| front/crown outer arc | `232.063` | `136` | `204` |
| front/crown profile depth | `112.200` | `66` | `99` |
| rear-base height | `159.720` | `94` | `141` |
| rear-base width | `129.360` | `76` | `114` |
| leaf root-to-tip | `147.840` | `87` | `130` |
| leaf long diagonal | `151.478` | `89` | `133` |
| leaf width | `72.600` | `43` | `64` |

The 2 px column is a simple conservative fallback when the implementation
cannot allocate density non-uniformly.  It is not necessary over every smooth,
noncritical stretch.

## Grid and curve counts

The dimensions below are per skin.  The paired pocket has both an outer and
an independently authored inner skin plus its bridge.

### Dense explicit evaluated nets

An unrefined tensor grid needs at least:

| Part | Ordinary 3 px floor | Uniform 2 px fallback |
| --- | ---: | ---: |
| front/crown | `137 x 67` vertices | `205 x 100` vertices |
| rear base | `95 x 77` vertices | `142 x 115` vertices |
| dominant leaf | `90 x 44` vertices | `134 x 65` vertices |

The ordinary floors total `40,908` paired outer/inner vertices before any
extra roll loops; the uniform 2 px fallback totals `91,080`.  These counts are
computationally cheap, but directly authoring every coordinate is a poor form
control interface: it makes low-frequency camber, layer separation, and
outer/inner correspondence harder to reason about and makes accidental waves
easy to introduce.

If separate boundary curves drive a more efficient non-tensor surface, their
evaluated curve segment floors are the chain values in the preceding table.
Interior topology may coarsen only away from the silhouette, root, overlap,
and bridge.  Any 2:1 transition or extraordinary vertex stays off those
identity-critical courses and must pass the selector's anti-pinching gate.

### Applied Catmull-Clark

Two binary refinement levels produce nominally four evaluated spans per cage
span.  The mathematical level-2 cage floors are:

| Part | Minimum level-2 cage | Resulting evaluated grid |
| --- | ---: | ---: |
| front/crown | `35 x 18` vertices | `137 x 69` vertices |
| rear base | `25 x 20` vertices | `97 x 77` vertices |
| dominant leaf | `24 x 12` vertices | `93 x 45` vertices |

Those are strict planning floors, not recommended authoring counts.  They
leave no budget for Catmull-Clark arc-length change or for the 2 px critical
subchains.  Use at least these reserved, non-uniform cages:

| Part | Recommended level-2 cage | Resulting evaluated grid |
| --- | ---: | ---: |
| front/crown | `40 x 20` vertices | `157 x 77` vertices |
| rear base | `28 x 23` vertices | `109 x 89` vertices |
| dominant leaf | `27 x 14` vertices | `105 x 53` vertices |

That approximately 15% major-axis reserve can allocate 2 px spacing over up
to about 30% of a chain while keeping all remaining spans at or below 3 px.
The allocation must be based on projected arc length, not uniform parameter
steps.  If a critical portion is longer, add controls until the actual
projection passes.  If local allocation is unavailable, use the uniform 2 px
level-2 cages: `52 x 26`, `37 x 30`, and `35 x 17`, respectively.

One refinement level requires much denser authoring cages
(`69 x 34`, `48 x 39`, `46 x 23` at the ordinary floor).  Three levels can
reduce them to `18 x 10`, `13 x 11`, and `13 x 7`, but that makes the visible
form depend too heavily on subdivision shrinkage and support-loop placement.
It reopens the canopy/helmet failure route.  Level 2 is the defensible balance.

## Required support and bake protocol

1. Author the visible outer cage from the reference silhouette, camber, and
   layer order.  Author the inner cage separately; do not generate either as a
   receiver offset.
2. Give each skin an inboard perimeter support course and split the bridge so
   the edge roll is controlled on both sides.  The support courses are form
   controls, not substitutes for the projected boundary counts above.
3. Put exact anchors at the crown extrema, both temple turns, every rear lobe
   tip/notch, the leaf tip, and overlap corners.  Adjacent evaluated spans are
   at most `2 px`; ordinary spans are at most `3 px` in every frozen view.
4. Precompensate and reproject the cage so two refinement levels move no
   silhouette sample by more than `1 px` and create no new extremum.
5. Apply/freeze refinement into a fresh mesh.  The saved candidate has no live
   Subdivision modifier.
6. Validate the exact applied mesh: closed two-manifold pocket, positive
   volume, no self-intersection, no correspondence crossing, area-weighted
   oriented thickness `2--5 mm` over at least 95%, no connected out-of-band
   patch over one pixel, and the full signed-root and retained-layer gates.
7. Reopen the frozen file and repeat projection, topology, thickness, and root
   tests before rendering.

## Chord-error interpretation

For a circular arc of radius `R` and chord length `c`, the interpolation
sagitta is:

```text
s = R - sqrt(R^2 - (c / 2)^2)
```

At the 3 px span limit (`c = 1.71094 mm`), a radius of at least
`0.92676 mm` already implies `s <= 1 px`.  Padded-felt silhouette rolls should
normally exceed that radius, so the straight-span gate is more restrictive
than geometric chord sagitta.  Deliberate sharp tips and corners are anchored
vertices rather than fitted circular arcs and use the 2 px adjacent-span gate.

This does **not** mean the reference chord-error gate passes automatically.
Density bounds only interpolation error.  After camera alignment, calculate
maximum and P95 normal distance from the evaluated candidate boundary to the
controlling reference boundary; require `<=1 px` at identity-critical extrema
and `<=2 px` elsewhere.  A dense curve in the wrong place still fails.

## Existing-grid diagnosis

Using the same conservative direction lengths, the current proposal averages:

| Existing grid | Major average span | Minor average span |
| --- | ---: | ---: |
| front/crown `9 x 7` | `50.9 px` | `32.8 px` |
| rear base `9 x 9` | `35.0 px` | `28.4 px` |
| dominant leaf `7 x 7` | `44.3 px` | `21.2 px` |

That is a representation defect, not a polish issue.  Do not build or render
the coarse version.
