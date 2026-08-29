# A77 next-module decision

## Verdict

**Revise.** Do not put a fabric shader pass on the rung-003 whole model yet.
Use A77 for one sparse, untextured, whole-context head/hair macro assembly.
Keep the shader design as a reusable non-critical-path coupon for after the
silhouette and construction gate passes.

## Why

The requested outcome is a reusable 25 cm Reimu Fumo that is recognizable in
front, side, rear, and three-quarter renders and already reads as constructed
plush before microdetail. The protected rung remains the immutable comparison
baseline.

The strongest case for a material pass is that the user repeatedly identifies
the current result as unlike real plush, and the physical references clearly
show distinct short-nap hair, quiet face fleece, woven clothing, and black
velour. A geometry-preserving material copy is cheap, reversible, and useful
later.

The decisive objection is visible in the current pixels before implementation
details are considered: the head/hair reads as one deep, uniformly rounded
helmet; the lower stack reads as an upright cone; and the panels do not form
the photographed layer order. The independent blind review scores macro form
`3/10`, construction `2/10`, and medium read `2/10`. Procedural fibers would
therefore make a fuzzy helmet and cone. This violates the large-to-small gate
and cannot satisfy the untextured plush-construction requirement.

## Alternatives

- Material pass now: low cost and reusable, but cannot change the failed
  silhouette or panel construction. Keep only as a later coupon.
- Another receiver loft, dense formula field, synthetic sculpt stroke, or
  flat visible-envelope stack: rejected because A68--A76 already falsified
  those representation families.
- Sparse direct BMesh/Edit-mode macro assembly: predictable broad controls,
  materially different from the failed dense/stroke/loft approaches, and can
  be rejected from three cheap fixed views before detail work.
- Lower-body module first: valid second priority, but it leaves the dominant
  identity-owning head/hair failure untouched.

## A77 boundary

Build from the exact rung-003 hash on an isolated copy. Replace only the
silhouette-owning head/hair base and the minimum face carrier/interface needed
to judge it. Keep bow, body, sleeves, skirt, feet, cameras, and lights as
whole-result witnesses. Use a low-density direct control cage or explicit
fabric panels; do not recreate a sphere, superellipsoid, loft, open card, or
dense scripted brush field. Render front, worse three-quarter, profile, and
rear before any materials or microdetail.

Stop after the first packet if it still reads as a helmet, egg, mattress,
face card, bald shell, or detached panel stack, or if macro silhouette is
below `6/10`. A passing first packet earns one bounded correction; a failure
changes representation again instead of accumulating polish.

## Evidence

- `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/five_view_gate/contact_sheet.png`
- `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`
- `projects/renders/blender/fumo/reimu_fumo/references/physical_front.png`
- `projects/renders/blender/fumo/reimu_fumo/references/physical_side.png`
- `out/reimu_fumo_attempt_077_next_module/blind_review.md`
- `out/reimu_fumo_attempt_077_next_module/process_audit.md`
- `out/reimu_fumo_attempt_076_live_sculpt/decision_review.md`
