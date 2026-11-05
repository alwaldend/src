# Frozen bow030 technical audit

The scoped proportion edit preserves all 76 non-target controls, the recorded
rig/rest pose and appearance, both inner root columns, cage Y, all other
lattices, and all 116 knot-root witnesses. No new crossing was found in the
tested lowered loop/ruffle patches. This is not a visual or whole-asset pass.
The left closed loop retains inherited inward winding; it is not described
as outward or repaired by this audit.

Candidate: `bow_030_candidate.blend`, SHA256
`4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b`.
Independent source: `hand_029b_candidate.blend`, SHA256
`9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d`.
Both clean-reopened with repository-pinned Blender 5.2.1 LTS,
build `9e2066aef7ef`, background/autoexec disabled. Audit runtime was 1.60 s;
both input hashes remained unchanged. No model save or render was performed.

## Exact preservation and dimensions

The 76 evaluated control records match both the independently reopened
source and writer receipt. No missing or unexpected visible geometry.
All ten dependencies preserve base coordinates, materials, weights where
applicable, modifier targets, evaluated polygon indices and resolved face
material names. All material slots are persistent database IDs.
Each cage changes exactly its 16 outboard points; both inner U columns and
every Y coordinate remain exact. All other lattice records, including
`022 bow proportion cage`, are identical.

The head-width datum is the evaluated hood width, 117.439255 mm.

| Measurement | Source029b | Candidate030 |
| --- | ---: | ---: |
| Upper loop plus ruffle width / head width | 1.764017 | 1.500822 |
| Upper loop plus ruffle height / head width | 0.535250 | 0.610381 |
| Complete bow width / head width | 2.006481 | 2.006481 |
| Bow top above hood crown / head width | 0.214212 | 0.214125 |

Upper bounds become X [-88.127568, 88.127866] mm and
Z [147.234872, 218.917549] mm. The previous highest ruffle vertex moves
13.629973 mm inward, +0.005387 mm in Y and -0.010133 mm in Z; its X is
intentionally not fixed. Maximum evaluated ruffle Y motion is 0.064254 mm
left and 0.061866 mm right. Non-ruffle dependency Y coordinates are exact.
Post-cage Solidify normal offsets explain this distinction from fixed cage Y.

All 116 source vertices within |X| <= 7.75 mm retain identical evaluated
coordinates. Their nearest knot witnesses and distances are unchanged;
maximum root displacement and knot-distance change are both exactly zero.
This preserves the existing connection, not a new claim of sewn topology.

## Mesh integrity and inherited qualifications

All ten evaluated meshes have zero nonfinite coordinates, zero edges shorter
than 1e-10 m, zero faces below 1e-14 m² and zero inconsistent shared-edge
winding. Counts and polygon indices remain unchanged.

Each loop is closed: 952 vertices, 1,900 edges, 950 faces, Euler 2.
Each ruffle retains its open 224-vertex base with 116 boundary edges;
Solidify produces a closed 448-vertex, 446-face evaluated shell. The four
root-fold curves and two zigzag curves retain their existing open ends;
they are not new closure requirements.

The left loop's evaluated signed volume is negative in source029b
(-3.4657073e-5 m³) and030 (-2.9432474e-5 m³). It is consistently wound but
inward, inherited from its unchanged base. The right remains positive
(source +3.3706871e-5, candidate +2.8877220e-5 m³). Root is considering a
separate, explicitly recorded winding-only correction;030 bytes are intact.

## Lowered-boundary contact regression check

Selected corresponding vertices lowered by more than 0.1 mm, plus every
triangle touching them: left loop 458 selected / 978 patch triangles;
right loop 478 / 1,020; each ruffle 312 / 640. Tested both triangulated-edge
directions against the unchanged head core, hood, all three rear hair cloths
and same-side bow tail and tail ruffle: 28 pairs in each candidate.

All tested source and candidate crossing counts are zero. Minimum sampled
candidate patch-vertex distance to the tested head/hair is 15.767617 mm.
The closest tested pair is the right loop to right tail ruffle:
1.464749 mm in source, 1.517663 mm in030. No new contact regression is
established by these tests. Nearest distances are vertex-to-triangle samples,
not exact mesh separation bounds; no coplanar, global collision, animation
or all-view visibility guarantee is claimed. Untested inherited contacts
remain outside this bounded proportion audit.

## Machine evidence

Audit JSON included in [machine evidence](bow_030_machine_evidence.md), SHA256
`7a0cbf817b55558b090b8b04e9dff34c0b3e4f4d565fe0f0177b1137d052b51b`,
contains the compact controls/cage/topology summary and all 28 pair results.
The separately bound task-local bow_030_contact_witnesses.json, SHA256
`5c2215d133f68a73cd3cf3fa577930ad540a77bbbc9f80ac611d9ee27b8162fe`,
contain all 116 root rows; there are no intersection-witness rows because
the scoped crossing tests returned zero. These are ignored task artifacts,
not claims that raw witnesses have been stored canonically.

Script: `audit_bow_030.py`, SHA256
`5bd5135597027b6d89c3764ec5522da8aa516dd879c509f666e963e66935c8b1`.
