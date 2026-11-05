# Chin 071: image-only reference diagnosis

Observed 2026-09-05 19:01:53 UTC. Candidate: fabric_070. Reviewer inspected no model code or geometry. This is a bounded lower-face diagnosis, not a complete asset approval review; camera parameters and hidden chin-center geometry are unavailable. Parent owns acceptance. Branch/worktree identity verified before writing this ignored report.

Evidence: `fabric_070_fast_review/{front,side,three_quarter}.png`, `fabric_070_close_compare/pile_on.png`; reference `physical_front.png`, `physical_side.png`, `canonical_front_25cm.png`; turntable `side_062_frame_{06,08,10,12,26,28}.png`. Paths are under this report's directory except the three named still references, which are under `projects/renders/assets/reimu_fumo/references/`.

## Finding

The lower face looks tucked back, but a measured excessive chin-center setback is **not established**. Candidate profile rolls backward below the eyes. The physical plush does too, particularly frame 10. Hair hides the actual chin-center/underside transition in both side controls; the visible cheek edge must not be reported as a measured chin.

The clearer mismatch is front seating: candidate has a short mouth-to-underside span and a conspicuous red crescent between the face and collar. Its central underside appears nearly level or slightly raised relative to the adjoining cheeks. Both frontal controls have a gently convex cushion underside seated directly onto the white collar. The close three-quarter candidate makes the exposed red region especially conspicuous, amplifying the recessed impression.

## Controls and rough landmarks

Use `physical_front.png` as primary for cushion underside/collar seating, with canonical front corroborating graphic spacing. Use frame 10 as primary for the visible lower-cheek profile; physical side and frame 12 chiefly establish occlusion uncertainty. These are visual pixel estimates, not aligned-overlay acceptance measurements. `Wh` is frontal hair/head width excluding bow and hanging side locks.

| Landmark | Candidate | Reference | Limitation |
| --- | --- | --- | --- |
| Mouth to central face underside | ~0.07–0.09 Wh | ~0.10–0.12 Wh | ~0.01–0.02 Wh edge/crop uncertainty |
| Eye-height face edge to last visible lower-cheek edge, profile setback | ~0.16–0.20 Wh | frame 10 ~0.14–0.21 Wh | Different yaw/pitch; tuft occlusion; **not chin-center depth** |
| Red gap above white collar, central front | ~0.02–0.04 Wh | No comparable broad red crescent visible | Cloth overlap rather than bare neck; depth unmeasured |

Head pitch cannot be isolated from camera elevation in these images. The eye/mouth arrangement supplies no compelling evidence for rotating the whole head. Cheek rollback, chin underside height, and collar exposure are separate observations.

## One bounded correction

Fill/extend the central lower-face underside downward by approximately 0.025–0.035 Wh, smoothly fading toward the lateral cheeks, preserving eye/mouth positions and the upper face. Expected effect: a gently convex stuffed-cushion bottom and reduced red exposure in front/three-quarter. Preserve the rounded side transition; do not add a projecting anatomical chin or globally shift/pitch the head.

Risk: collar intersection or an excessively long/square lower face. Judge at unchanged cameras; reject if the underside bulges ahead of the mouth plane, clips the collar, or changes cheek width materially. Profile forward extrusion is currently less supported than restoring underside seating.

Evidence SHA-256 prefixes: candidate front `30d4aee2b36b5bce`; candidate side `3deeaebbcde2088a`; frame 10 `37c1e2866fdbe97c`; physical front `f8c7d0f9911dbff1`.
