# A96 first-pixel verdict

## Decision

`RESET`. The candidate fails the absolute category gate, so no blinded reviewer
and no correction are authorized.

The coordinator opened and inspected the candidate's front, exact-side, and
worst-three-quarter pixels before consulting structural measurements. All are
valid, non-black complete-subject renders. The failure is visible form, not
render validity.

## Bound candidate

- Blend SHA-256:
  `73099d33e19e7cc73be6b1184df3b3a41a437bcad02d481ae02bf8988f6699a7`
- Front SHA-256:
  `530623567a886536e4d5b30a673729713afd9182d0db4338aa42fc2efa593b7f`
- Side SHA-256:
  `a4e7dd72bfb389cfacdb8a0d457a2dce45ed056365d001e06d6560c282eb1441`
- Three-quarter SHA-256:
  `073cfc22f0061a3cdea3725ff58728d297101beb96c32a817c48d4e4ac93a5df`
- Baseline/candidate board SHA-256:
  `efe90195152041fc60560098cb6704f50ecf6528a30711047b7c6ff2fe5330d8`

The board is ordered by row as exact baseline on the left and candidate on the
right: front, exact side, then three-quarter.

## Pixel judgment

### Front

The two feet are taller and more visibly stuffed than the baseline, which is a
useful local signal. However, they remain exposed in front of an almost
horizontal hem and dark continuous under-skirt rail rather than being softly
occluded by pooled ruffles. The skirt remains a rigid trapezoid. This view does
not reach the required construction or contact category.

### Exact side

This is a decisive failure. The lower silhouette now extends farther rearward
as a thin red-and-white ramp/cape ending in a sharp vertical edge. The only
visible foot remains a forward ball, spatially separated from the trailing
cloth. The subject does not read as sitting on a compact stuffed body or pooled
garment.

### Three-quarter

The candidate exaggerates the baseline's triangular ramp and nearly horizontal
floor rail. The foot forms are larger, but their relationship to the hem is
less plausible: two foreground pods plus a long trailing sheet, not one seated
plush assembly.

## Scores for the edited lower assembly

| Category | Baseline | Candidate | Required |
|---|---:|---:|---:|
| Lower silhouette | 3.5/10 | 3.0/10 | at least 6/10 |
| Plush construction | 3.0/10 | 3.0/10 | at least 6/10 |
| Contact and occlusion | 3.0/10 | 2.5/10 | at least 6/10 |

The candidate is not clearly preferred to rung003 in every view. Its useful
foot-height improvement cannot offset the categorical side and three-quarter
regression.

## Gate consequence

The predeclared correction exception applies only after the shape category
passes and exactly one localized contact defect remains. Here the core shape
category fails in two controlling views, so refinement would violate the plan.
The attempt stops at the first pixels. The tracked canonical Blend and exact
rung003 parent remain untouched.

Because A96 demonstrated broad, deterministic coordinate control but still
failed to change the construction category, its stop condition requires an
explicit feasibility/strategy decision before another rebuild. It does not
authorize another lower-stack generator or another parameter variation.

