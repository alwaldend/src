# A93 authoring-interface decision review

## Verdict

Proceed with one disposable sparse Edit-mode proportional-deformation coupon.
Do not invoke `bpy.ops.sculpt.brush_stroke`, do not build a full character, and
do not resume analytic panel, loft, or primitive-stack generation.

This is the smallest test that can answer the current question: can the
available Blender path move broad, low-frequency regions in three independent
directions with predictable support and visible silhouette change? Passing
does not establish artistic quality; it authorizes exactly one low-resolution
manual-style whole-plush sculpt. Failure closes this modality.

## Actual outcome and constraints

The outcome remains a reusable, animatable, 25 cm Reimu Fumo matching all
canonical views. A91 and A92 both failed the controlling front/three-quarter
gate, so another full analytic character would violate the A90 recovery plan.
The tracked Blend and rung003 must remain byte-identical.

## Strongest reason to reject this plan

Programmatic proportional editing can become another analytic deformation
family: it may prove only that code can move vertices, not that the agent can
art-direct a plush. A broad numerical support mask can still make a bad shape.

The gate therefore binds before/after pixels, object-space displacements, and
an explicit human visual judgment of the silhouette. It also forbids using the
coupon as likeness evidence. If the later whole-plush sculpt still relies on
formula fields instead of sparse observed corrections, the coupon has not
authorized it.

## Disconfirming evidence against native sculpt strokes

- A68 used pinned Blender 5.2.1 in a foreground `VIEW_3D`; its useful support
  remained 35 vertices and larger strokes produced dimples.
- A71 repeated the path on uniform topology; ordinary Grab moved 11 of 3,863
  vertices and Elastic Grab moved only 40 vertices more than 0.10 mm.
- A76 established that the available programmable Blender surface exposes no
  genuine pointer drag distinct from `bpy.ops.sculpt.brush_stroke`.

Calling that operator again would vary transport, not the authoring method.

## Alternatives

| Route | Evidence value | Cost/risk | Decision |
| --- | --- | --- | --- |
| Scripted native sculpt stroke | Already falsified as broad control | Cheap but repetitive | Reject |
| Another complete analytic model | Direct pixels, but A91 just failed the representation | High wasted work | Reject |
| Sparse direct proportional edit coupon | Tests explicit broad support and non-collinear control | Low, reversible; artistic transfer remains uncertain | Proceed once |
| Human pointer sculpt | Best artistic loop in principle | No auditable pointer provider in the current interface | Reconsider only if capability changes |

## Reversible discriminator

On one connected soft low-resolution coupon, perform exactly three sequential
edits with intended dominant directions `+X`, `+Y`, and `+Z`. Each edit must:

- move at least 8% of vertices by more than 0.10 mm and at least 1% by more
  than 0.50 mm;
- form one dominant connected support region spanning at least 30% of the
  relevant coupon dimensions;
- change a meaningful silhouette region in a fixed controlling view;
- avoid a dimple, hard curvature ring, global affine motion, camera change,
  topology change, or zero-effect result; and
- leave the previous edits visible, proving independent sequential control.

Any miss stops the coupon without parameter tuning. A pass authorizes one
complete low-resolution sculpt and nothing later in the pipeline.
