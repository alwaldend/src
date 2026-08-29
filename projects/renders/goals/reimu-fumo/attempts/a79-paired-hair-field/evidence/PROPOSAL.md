# A79 compact rear-base paired-panel study

## Verdict

**NO-GO for the published exact net. Do not build it in Blender.**

The study establishes that a reference-authored shallow cut-pattern surface
can meet the rear silhouette, root-span, camera-density, two-axis-camber, and
non-concentricity requirements simultaneously. Its current square-to-outline
parameterization cannot meet A79's topology and paired-inner gates, however.
Continuing by densifying or smoothing this net would conceal rather than fix
that representation defect.

## What was frozen

- `rear_panel_control_net.json` contains literal `41 x 41` outer and inner
  control coordinates in millimetres.
- `rear_panel_literal_mesh.json` contains the exact `281 x 281` paired mesh
  coordinates, exact root-loop indices, provenance, and deterministic
  connectivity. There is no Blender modifier or receiver-derived visible
  outer surface.
- `prove_rear_panel.py` is a standard-library-only shape, projection,
  topology, thickness, root, and receiver-offset proof.
- `proof.json` is the complete result.
- `rear_panel_projection.png` (with a lossless PPM source) shows rear, side,
  and three-quarter wire
  projections. Brown is the reference outer net, red is the hidden root loop,
  and blue marks sampled outer/inner pairs.

Only inner vertices of the closed root loop register to the receiver proxy.
Every visible outer coordinate comes from the rear cut outline and independent
two-axis camber. The classifier excludes only a three-cell hidden root collar.

## Passing evidence

| Gate | Result |
| --- | ---: |
| rear width | `0.93999 Wh` |
| rear height | `1.14993 Wh` |
| maximum-width height | `0.6925` of rear height |
| broad unequal lobes | `5`; five unequal tip depths |
| transverse / longitudinal camber | `11.5 / 11.5 mm` |
| hidden root X / Z span | `108.69 / 105.99 mm` |
| hidden root median Y | `40.29 mm` |
| 64-sample proxy root coverage | `100%` in `+0.2..+0.8 mm`; zero negative |
| maximum normal boundary span at 512 px | `2.93 px` |
| maximum critical boundary span at 512 px | `1.64 px` |
| rear chord error at 512 px | `0.45 px` |
| best receiver-offset classifier fraction within `+/-1 mm` | `0.6497` (`<0.80`, so non-concentric) |
| closed paired pocket by edge counts | pass; Euler characteristic `2` |

These results answer the key representation question: reference silhouette
and receiver seating do **not** force a receiver-concentric visible field.

## Fatal failures

| Gate | Result | Why it is fatal |
| --- | ---: | --- |
| outer quad angles | `15.34..179.98 deg`; 2 foldovers | A single rectangular patch wrapped around a smooth closed outline becomes singular at its parameter corners. |
| outer adjacent-area ratio | `6.71` | Exceeds the selector's `4:1` ceiling. |
| inner quad angles | `0.06..179.89 deg`; 1,146 foldovers | Registering one internal square loop while retaining one rectangular correspondence field pinches the hidden skin. |
| inner aspect / adjacent-area ratio | `18.92 / 1179.37` | Not a repairable local tolerance miss. |
| paired thickness hard minimum | `1.42 mm`; normal ratio `0.794` | A small root-neighbourhood failure still violates the all-positive constructed-pocket contract. |
| inter-lobe relief | `9.32 mm` | The smoothed cut outline undershoots the frozen `13.2..18.48 mm` range. |

## Decision-review conclusion

The strongest case for this path is that its visible form is genuinely
reference-authored: it passes the receiver-offset veto by a wide margin while
matching the dominant rear measurements. The decisive rejection is
topological, not cosmetic. A smooth closed silhouette and an internal closed
registration loop should be represented by a multi-patch disk (or another
layout with a literal internal seam loop), not one square radial patch.

That replacement is larger than a bounded correction to this exact net. The
selector should prefer another candidate that already proves the same
non-concentric silhouette without these singularities. This study should be
retained as evidence for its cut outline, camber field, and root feasibility,
but its connectivity must not be merged into Variant A.

## Evidence hashes

- `prove_rear_panel.py`:
  `55b4c9f69192312c67682d7457f2a026b5f56b541863364f28057686b6afd1f7`
- `proof.json`:
  `717a749dfa78dcb1e7da495806a28b5e67ee33e6794b4d35c36fb100f0a7ea1b`
- `rear_panel_control_net.json`:
  `e29242e9dda614229a67758fc82cfca6fe2be3c8298928f425217d6cca1fb039`
- `rear_panel_literal_mesh.json`:
  `fb10f81a80c13d090bd7605f5bf045d7a90c89a2632a69e4295c5bb046e2db08`
- `rear_panel_projection.ppm`:
  `a4173d4742c7cb60866240f151d26128f4cdad37196c7d8863b319415464190b`
- `rear_panel_projection.png`:
  `9ea94aef28ddd5538fe2a3248ac057227a3e7fe3428fae5f7d4847f2043d4649`
