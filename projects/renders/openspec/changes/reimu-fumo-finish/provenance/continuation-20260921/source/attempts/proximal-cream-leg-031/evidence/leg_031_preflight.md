# Leg 031 preflight

## Decision

Proceed with one bilateral **proximal cream-root emergence** unit at the
current foot placement. Preserve the retained 023/026 black rounded pods,
their transforms, their depth, the curled hem, and the current ground
contact. Reshape only the proximal part of each existing cream root so it
emerges upward, inward, and seatward as a soft lobe while its distal part
remains embedded in the black toe.

Do not recenter the feet in this unit. The current center separation is
`0.6267 Wh`, versus `0.563 Wh` measured in the canonical front and the
`0.560 +/- 0.040 Wh` landmark band. A reference-matching X translation was
tested and crossed the hem. One causal frontward seating adjustment made the
crossing worse. The spacing residual is preferable to changing the retained
hem or iterating an unbounded placement search during this unit.

## Evidence binding and method

- Candidate: `hand_029b_candidate.blend`
- Candidate SHA-256 before and after inspection:
  `9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d`
- Exact five-view receipt:
  `hand_029b_eevee_review/render_receipt.json`
- Canonical front: `references/canonical_front_25cm.png`
- Physical side: `references/physical_side.png`
- Normalization: evaluated `Hair028 crown and back hood` width,
  `Wh = 0.1174392551 m`; canonical front `Wh = 368 +/- 4 px`.

The asset faces `-Y`; front-camera rays travel `+Y`, the physical side is
compared to the `+X` side, and `+Z` is up. Bounds below are evaluated
world-space mesh bounds. Occlusion was checked by 72 by 72 first-hit ray
samples over each projected root bound. Contact was checked with evaluated
mesh BVHs: triangle-pair overlap plus a bidirectional vertex-to-surface
nearest sample. The side reference is oblique, so it controls construction
and a depth band rather than an exact overlay.

## Front and side proportions

| Measure | Candidate | Reference / contract | Finding |
| --- | ---: | ---: | --- |
| Left black toe width | `0.2467 Wh` | `0.209 Wh` image | About `+0.038 Wh` |
| Right black toe width | `0.2501 Wh` | `0.215 Wh` image | About `+0.035 Wh` |
| Black toe height | `0.2166 Wh` | `0.204-0.207 Wh` image | About `+0.01 Wh` |
| Toe-center separation | `0.6267 Wh` | `0.563 Wh` image | About `+0.064 Wh` |
| Black toe side depth | `0.3995-0.4073 Wh` | `0.32-0.42 Wh` contract | Pass |

The canonical black-toe masks measured 77 and 79 px wide, 75 and 76 px
high, with centers 207.2 px apart. Threshold placement and the soft image
edge contribute about `+/-0.01 Wh` uncertainty. The candidate pods are near
the width tolerance edge, but their rounded silhouette and physical-side
depth already belong to the retained representation and should not be
rebuilt.

In the canonical front, pale cloth is visible about `0.06-0.09 Wh` inward
past each black toe's inner edge and about `0.04-0.07 Wh` above its top.
Those are approximate occluded image landmarks, not new hard contract
limits. The physical side shows one broad proximal pale lobe continuing
behind and above the black toe, not a parallel-sided lower-leg tube.

## Existing objects and interfaces

The relevant objects are renderable, rigged meshes rather than hidden or
disabled placeholders:

- `Left black stuffed foot pod`: material `Feet black velour.002`, armature
  group `Leg_L`, geometry SHA-256
  `873d3a562a179b1aa86e869842d349421cafc0baf506d672306847aaef0fb2ae`.
- `Right black stuffed foot pod`: same material family, group `Leg_R`,
  geometry SHA-256
  `4ccea62076fbe96aa9425457acb28d4aff261ada29291c588bb564810f1a5693`.
- `Left short hidden leg root`: material `Dress warm white cloth.002`, group
  `Leg_L`, geometry SHA-256
  `31c898e3edc83b33ac1c7637157a8912f148c5178690adeb332868596fb48be5`.
- `Right short hidden leg root`: same material family, group `Leg_R`,
  geometry SHA-256
  `acbbf3b2f10d8fa154384dfb187c7812fda14d3e1630164cc47dcbca541eea1e`.

The roots are closed, rounded 1,962-vertex / 2,016-face volumes. Each root
is `0.1533 Wh` wide, `0.2044 Wh` deep, and `0.1185 Wh` high. They are not
cylinders, but their placement hides them:

- each root's X bound sits fully inside its toe by `0.045-0.050 Wh` on both
  sides;
- each root rises only `0.0219 Wh` above the toe crown;
- each root extends only `0.0347 Wh` seatward beyond the toe, or about 9%
  of the pod depth;
- front first-hit root coverage is only 7.74% left and 2.72% right;
- side first-hit coverage is 26% for the frontmost right root; the left root
  is hidden behind the right assembly in profile.

The roots are therefore present but geometrically occluded. At the current
placement, root-to-pod BVH overlap counts are 148 left and 150 right. Neither
root crosses `Hem026 curled cotton strip` or `Skirt022 joined gathered
panels`. Root-to-hem sampled clearances are `0.00422 Wh` left and
`0.00298 Wh` right.

The retained pods also have zero triangle crossings with the hem. Their
nearest sampled hem clearances are `0.00317 Wh` left and `0.00178 Wh` right,
at these world-space witness pairs:

- left pod `(-0.0343007, -0.0595349, 0.0246358)` to hem
  `(-0.0342388, -0.0595405, 0.0250033)`;
- right pod `(0.0269969, -0.0654789, 0.0209035)` to hem
  `(0.0268549, -0.0655105, 0.0210535)`.

These near-contact witnesses, not merely the unchanged floor Z, define the
026 hem support relationship.

## Why 024/025 differed

`leg_future_hypothesis.md` records that the rejected construction narrowed
toward a concealed rear while raising its center only 2.8 mm. It therefore
read as an exposed side cylinder or toe-ended cone. The present source roots
already have a closed, rounded topology. The needed change is to reveal and
shape their **proximal** volume while retaining the distal black pod, not to
replace either assembly with another tube family.

## Recenter falsifier and protected gate

Moving the complete left and right assemblies inward by `0.003917009 m`
(`0.03335 Wh`) each would produce the reference target separation of
`0.560 Wh`. With all intrinsic meshes and the hem fixed, that pure X test
created 22 left and 102 right pod-to-hem triangle-pair overlaps.

One causal seating alternative was then tested: the same X correction plus
`-0.0015 m` (`-0.0128 Wh`) in Y, frontward. It produced 43 left and 108 right
pod-to-hem overlaps. It did not cheaply preserve the retained contact, so no
placement sweep is justified for 031.

Acceptance gates for the cream-root unit are:

1. Preserve both black-pod geometry hashes, transforms, side-depth band, and
   rounded toe silhouettes.
2. Preserve zero pod-to-hem triangle crossings and sampled pod-to-hem
   clearance no greater than `0.01 Wh`; keep the witness on the matching pod
   crown rather than merely preserving ground Z.
3. Preserve a positive root-to-pod intersection at each distal seam and no
   visible gap greater than `0.02 Wh`.
4. Keep root-to-hem triangle crossings at zero and sampled clearance no
   greater than `0.015 Wh` so the pale lobe remains seated under the skirt.
5. In the frozen front and side reviews, expose a bilateral pale proximal
   lobe consistent with the measured front landmarks and the physical-side
   stuffed-lobe construction. Reject any parallel-sided tube or pointed
   toe-ended cone reading.
6. Record the unchanged `0.6267 Wh` center-separation residual explicitly;
   do not claim the spacing contract is repaired in 031.
