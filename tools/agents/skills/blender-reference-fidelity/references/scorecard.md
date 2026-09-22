# Reference-fidelity scorecard

Keep one scorecard per asset in temporary output. Fill it with measurements,
not adjectives.

## Reference map

| View          | Source/frame | Controls                                | Perspective uncertainty |
| ------------- | ------------ | --------------------------------------- | ----------------------- |
| Front         |              | Widths, vertical landmarks, symmetry    |                         |
| Side          |              | Depth, face plane, attachment seating   |                         |
| Rear          |              | Back silhouette, seams, rear attachment |                         |
| Three-quarter |              | Overlap, transitions, volume continuity |                         |

When references conflict, name one primary source for each row. Record other
sources as variant or uncertainty evidence; do not average them into an
unbuildable hybrid.

## Camera calibration

| View          | Projection/lens | Object rotation | Crop/scale | Alignment error |
| ------------- | --------------- | --------------- | ---------- | --------------: |
| Front         |                 |                 |            |                 |
| Side          |                 |                 |            |                 |
| Rear          |                 |                 |            |                 |
| Three-quarter |                 |                 |            |                 |

An overlay cannot pass a gate while its camera-alignment error exceeds the
landmark tolerance being judged.

## Normalized targets

Use head width (`Wh`) or another declared invariant as `1.000`.

| Landmark                             | Target | Tol. | Candidate | Delta | Pass? |
| ------------------------------------ | -----: | ---: | --------: | ----: | ----- |
| Overall height                       |        |      |           |       |       |
| Head height                          |        |      |           |       |       |
| Head depth                           |        |      |           |       |       |
| Face width and height                |        |      |           |       |       |
| Eye centers and projection           |        |      |           |       |       |
| Body width and height                |        |      |           |       |       |
| Major accessory span and angle       |        |      |           |       |       |
| Lowest and widest garment points     |        |      |           |       |       |
| Contact points and visible gaps      |        |      |           |       |       |
| Head front/rear planes and gusset    |        |      |           |       |       |
| Face opening and applique projection |        |      |           |       |       |
| Accessory root seating and thickness |        |      |           |       |       |
| Sleeve panel profile                 |        |      |           |       |       |
| Hem gathering and panel transition   |        |      |           |       |       |
| Foot occlusion by garment            |        |      |           |       |       |

Add subject-specific identity landmarks rather than forcing every model into
these example rows.

## Cycle record

- Accepted baseline:
- Dominant measured failure:
- Current representation:
- Why that representation can or cannot match the reference:
- One-cycle hypothesis:
- Components changed:
- Controlling view:
- Regression-risk view:
- Fast comparison paths:
- Four-view comparison paths:
- Landmark rows improved:
- Landmark rows regressed:
- Decision: accept / reject / rebuild subsystem
- Evidence for decision:

## Absolute image review

Complete this section before consulting implementation diagnostics or calling
the candidate ready for approval.

- Reviewer and context isolation:
- Unlabeled same-subject recognition: yes / no
- Intended-medium read:
- Five largest visible discrepancies:
- Overall reference likeness score (0--10):
- Macro silhouette and proportions score (0--10):
- Manufactured or anatomical construction score (0--10):
- Identity-defining features score (0--10):
- Contact, attachment, and occlusion score (0--10):
- Intended-medium score (0--10):
- Presentation readability score (0--10):
- Major visible failure present: yes / no
- Absolute decision: pass / reject

## Gate decision

Pass only when every critical row is within its declared tolerance, all
identity-defining silhouettes agree in every applicable view, the absolute
image review meets every threshold in `visual-quality-gate.md`, and no earlier
detail tier remains unresolved. Record uncertainty rather than silently
relaxing a target.
