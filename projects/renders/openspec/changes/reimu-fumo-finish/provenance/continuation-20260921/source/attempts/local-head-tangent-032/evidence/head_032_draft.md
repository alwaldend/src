# Head032 draft handoff

`head_032_draft.py` exposes exact `TARGETS` and `build_head_032()`.
It defines functions/constants only at import, never opens or saves a model,
and edits only the three existing meshes' base vertex coordinates. Root owns
the next saved candidate and visual decision.

- Base: retained `bow_030b_candidate.blend`, SHA256
  `d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6`.
- Evaluated helper SHA256:
  `f234f504cd481057d5018fb8c7e1d6dae6d5fea060830a51e050d2a14b10d318`.
- Small result: `head_032_draft_metrics.json`.
- Complete unsaved evidence: `head_032_dryrun_fixed.json`.
- Runtime: pinned Blender 5.2.1 LTS, build `9e2066aef7ef`.

## Implemented scope

```text
Head028 sewn cushion
Hair028 crown and back hood
Hair028 traced padded fringe
```

The existing outer front band uses a cubic between rho=.86 and the actual
retained seam near rho=.988. Inner position/slope come from the baseline
front receiver triangle; the outer position remains the actual seam point,
with the slope derived from the first retained gusset band. The correction
fades from zero at Z125mm to full at Z145mm. All cubic coefficients are checked
for monotonic depth before the first model edit; no coefficient sweep or
shape-parameter retry occurred.

Depth/Y displacement is inverse-mapped through the current Head pose and
object transform into existing base vertices. No new modifier, Geometry
Nodes system, shell, subdivision, padding or armature was added. Existing
weights, Armature/Solidify order, materials and topology are retained.

The hood receives exactly its corresponding core-vertex delta. The retained
recipe's hood face/index mapping is validated against actual topology before
editing. The fringe receives the same core-receiver depth difference while
retaining its current stand-off/padding; no independent fringe tuning occurs.

## Measured result, not visual approval

| Witness | Before | After |
| --- | ---: | ---: |
| Full-strength seam above145mm, median normal jump | 61.900 degrees | 0.506 degrees |
| Same86 seam edges, maximum normal jump | 76.214 degrees | 0.992 degrees |
| All108 seam edges above125mm, median jump | 67.343 degrees | 0.577 degrees |
| Fade-band22 edges, median jump | 76.782 degrees | 14.941 degrees |

The fade region intentionally retains a large maximum jump near its zero
endpoint:76.234 degrees. This is not a claim that the entire head seam is
already tangent continuous. All3389 distinct evaluated cubic profiles have
positive depth derivative; the minimum is0.009555m per normalized radius.
That checks the proposed meridional cubic, not global collision freedom or
the appearance of the resulting curve.

Maximum core motion is4.663mm forward in Y. Core X/Z, every zero-field core
vertex, the actual seam vertices, all rear/gusset faces and the lower125mm
region are exactly unchanged. The protected-face assertion covers7806 core
faces. All83 visible non-target controls remain exact, including facial
features, locks, collar and body/bow. Rig pose, target weights, modifier
stacks, polygon topology and material indices are unchanged.

The live Solidify layers have qualified derived movement:

| Derived geometry | Maximum X movement | Maximum Z movement | Maximum movement below125mm |
| --- | ---: | ---: | ---: |
| Hood | 0.619mm | 0.414mm | 0.0211mm |
| Fringe | 0.583mm | 0.424mm | 0.00668mm |

These are normal-derived offset changes, not additional base X/Z edits.
All zero-field base vertices remain exact. Existing thickness and padding
parameters were not changed. The resulting cloth contact/overlap still needs
the coordinator's frozen audit and pixels.

## Guard and execution details

The helper verifies the baseline file hash and the frozen030b writer receipt
hash, then checks all15 visible head inputs against their exact evaluated
geometry/material/visibility/parent/modifier records. It snapshots all
non-target visible geometry and the rig again before editing. Default active
source must be030b. An explicit `AUTHORIZED_ACTIVE_SOURCE` path/hash override
exists for coordinator-only integration with a later retained file, but
cannot bypass the15 exact head-input guards; it is unnecessary now that030b
remains retained.

The first invocation stopped before any model edit because the hidden donor
was not fully evaluated by the active dependency graph. That failure remains
in `head_032_dryrun.json`. The causal correction derives immutable base-face
material tags through the known single subdivision level and verifies the
exact resulting hood topology, avoiding hidden evaluated geometry. Only then
did the single actual geometry evaluation run. No geometry parameter changed
between invocations and no blend was saved.

Next: root's guarded save and one fixed side/three-quarter comparison. The
test has removed the measured full-strength core kink; if the same ridge
remains visible, treat independent fringe-root overlap as a separate cause.
Reject a new bulge, shelf, pinch or contact loss even if these normal metrics
look good. No visual, rig-motion, physics or overall acceptance is claimed.
