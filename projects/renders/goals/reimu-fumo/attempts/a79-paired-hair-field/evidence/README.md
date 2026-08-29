# A79 render harness

This directory contains the frozen, no-candidate render harness for the A79
paired hair field. The renderer opens a disposable candidate read-only in the
repository-pinned Blender 5.2.1, renders beauty and material-semantic views in
one process, rehashes every protected input, and never calls a Blender save
operator.

## Fixed view and reference map

| View | Camera | Canonical evidence |
| --- | --- | --- |
| Front | Review_front_Camera | exact 25 cm front; turn frames 02--04 with frame 03 at front phase |
| First 3Q | Review_three_quarter_Camera | frames 06--09 |
| Positive profile | Review_side_Camera | frames 10--12; frames 10/11 nearest +90 degrees |
| Rear | Review_rear_Camera | frames 17--20; frame 18 nearest rear |
| Negative profile | Review_side_mirror_Camera | frames 25--26, only when the camera exists |
| Opposite 3Q | Review_three_quarter_mirror_Camera | frames 27--29 |

The turn advances about 12 degrees per frame with a one-frame phase
uncertainty. The rung003 source has no Review_side_mirror_Camera, so the dry
run records the negative profile as an explicit optional omission. The
composer still shows the directly observed frames 25--26 and never mirrors
the other profile.

## Entry points

- batch_render_spec.json pins the source hashes, Blender version, resolution,
  cameras, frame neighborhoods, semantic materials, and one-pixel thresholds.
- run_batch_render.sh is the executable Bazel/pinned-Blender wrapper.
- render_batch.py reuses projects/renders/cmd/fumo_review/render_packet.py for
  deterministic PNG settings, camera rendering, ID emissions, and PNG
  canonicalization.
- evaluate_semantic_masks.py extracts brown and beige masks and compares a
  candidate to rung003 for beige leakage and a background seam wider than one
  pixel.
- compose_comparison.py creates a reference/rung003/candidate board. A missing
  candidate becomes a clearly labeled baseline-only row.
- validate_spec.py validates hashes, camera availability, all controlling
  frames, and the effective repository render spec without starting Blender.
- capture_manifest.py hashes the harness, generated evidence, and protected
  external inputs.

## Dry verification

The checked dry run lives under verification/. It is intentionally based on
the read-only rung003 camera inventory and the already-published rung003 beauty
packet. It does not open Blender and does not render a candidate.

Run the same no-Blender checks from the workspace root:

    PYTHONDONTWRITEBYTECODE=1 python3 \
      out/reimu_fumo_attempt_079_paired_hair_field/render_harness/unit_checks.py

    PYTHONDONTWRITEBYTECODE=1 python3 \
      out/reimu_fumo_attempt_079_paired_hair_field/render_harness/validate_spec.py \
      --workspace "$PWD" \
      --spec out/reimu_fumo_attempt_079_paired_hair_field/render_harness/batch_render_spec.json \
      --camera-inventory out/reimu_fumo_attempt_078_head_rest/interface_inventory/interface_inventory.json \
      --output out/reimu_fumo_attempt_079_paired_hair_field/render_harness/verification/spec_validation.json

## Future candidate packet

Only after a candidate is frozen under the A79 output tree, run:

    out/reimu_fumo_attempt_079_paired_hair_field/render_harness/run_batch_render.sh \
      out/reimu_fumo_attempt_079_paired_hair_field/CANDIDATE.blend \
      out/reimu_fumo_attempt_079_paired_hair_field/render_harness/runs/CANDIDATE_PACKET

Use a newly named output directory for every packet. Then run the semantic
evaluator with the separately rendered rung003 packet as --baseline-packet,
and pass both packets to compose_comparison.py. Do not use a candidate-only
profile when rung003 lacks the same fixed camera; that column remains
non-comparable until a common pinned camera exists in both sources.

## Timing estimate

The prior rung003 five-beauty-view packet measured 73.53 seconds at 320 px.
The pinned 512 px packet adds five unlit semantic renders, so budget about
three to five minutes for the five available views. If a common optional
negative-profile camera exists, add roughly 20--45 seconds. Composer,
diagnostics, and manifest capture should remain below five seconds combined.
