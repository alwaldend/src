# Frozen head032 scoped technical audit

The three-target fairing preserves all 83 controls, recorded rig/appearance
and protected core interfaces. The measured front/gusset jump above 145 mm
falls below one degree; the 125–145 mm fade retains a 76.234-degree maximum.
Targeted exterior contact checks establish no new exposed crossing witness
in the tested directions. Inner-layer overlap increases and is qualified
below. No artistic, animation or whole-asset acceptance is inferred.

Candidate: `head_032_candidate.blend`, SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
Independent source: `bow_030b_candidate.blend`, SHA256
`d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6`.
Both clean-reopened in pinned Blender 5.2.1 LTS, build `9e2066aef7ef`,
background/autoexec disabled. Both hashes remain intact. No save or render.

## Preservation and integrity

All 83 evaluated control records match independently reopened 030b and the
receipt. Recorded rig/rest pose, appearance and visibility match. No new
visible objects or modifiers. All three targets preserve topology, material
regions, unit Head weights and modifier settings. Materials are persistent
database IDs. All evaluated meshes are closed, consistently wound, positive
volume and free of nonfinite vertices, edges below 1e-10 m or faces below
1e-14 m². Hood/fringe open bases remain expected inputs to Solidify.

Core X/Z is exact. Rear 4,896 faces, gusset 1,248, lower 125 mm: 3,708 and
underside 89 mm: 628 retain exact coordinates/winding hashes. All 157 front
seam vertices are stationary. Core depth motion peaks at 4.663127 mm.

All 639 collar/tie cloth points, nearest head points and unsigned gaps are
exact; gaps remain 0.297082–0.810112 mm. One left-collar sample picks adjacent
triangle 8015 instead of 8016 at an identical shared-boundary closest point.
Its signed-normal projection changes 0.000110 mm while its geometric gap
remains exactly 0.456368 mm. Therefore entire raw witness rows are not all
identical, but neck contact geometry is.

## Seam measurement and derived cloth motion

Normals use the same adjacent polygons and first-three-corner cross products
as the declared witness. All reported statistics reproduce the receipt.

| Region | Edges | Source median / maximum | Candidate median / maximum |
| --- | ---: | ---: | ---: |
| Above 145 mm | 86 | 61.900 / 76.214 degrees | 0.506 / 0.992 degrees |
| Fade 125–145 mm | 22 | 76.782 / 76.928 degrees | 14.941 / 76.234 degrees |

Base-world X/Z remains exact for all three targets. Evaluated hood maximum
X/Y/Z motion is 0.618845/4.727999/0.414044 mm; fringe is
0.582751/4.762724/0.424191 mm. Existing Solidify recomputes inner offsets.
Below 125 mm, evaluated maxima are 0.021120 mm hood and 0.006681 mm fringe.
Two hood base points themselves move 7.45 nm: their own Z is 124.942 mm,
but their mapped core receiver Z is 125.017777 mm. Thus a blanket assertion
that every hood-own-Z-below125 base point is bit-identical would be false.
Actual stored/world changed fringe base vertices number 2,833, versus 2,835
nonzero mathematical field values in the receipt; coordinate quantization
can erase a nonzero mathematical delta.

## Changed-band contact checks

Only moved target vertices above 10 nm plus their touching-triangle ring
were tested. Of 72 target/control pairs, 62 are separated by an AABB lower
bound exceeding1 mm. All 45 bow pairs are among these, with minimum bound
40.790748 mm. The four tested eye pairs have zero crossings in both
directions; other protected facial details are separated by the bound.

The six cheek/core/cover pairs already overlap in 030b. Counts change near
the retained cheek-root band; this is not a new disconnected contact pair.
Explicit outer-fringe/cheek tests remain zero in both directions.
Outer hood/cheek counts change, but all 29 added cheek-edge witnesses are
behind another scoped surface in front and both three-quarter ray directions:
at least 0.613356 mm among those tests. This is a bounded visibility witness,
not a statement about every camera or every possible overlap.

Full-shell hood-to-fringe crossings increase 171→954, reverse 855→3,126.
Their recorded layer pairs involve inner faces or rims. A separate explicit
outer-only test avoids first-hit masking and finds zero outer/outer crossings
in both directions, in both files. All 2,830 changed-fringe front rays remain
ahead of the hood: smallest projected front clearance 0.523515 mm source,
0.253435 mm candidate. These are Y-ray projections, not nearest distances;
some rays meet the rear hood where a front patch is absent. Increased
inner-shell overlap is real, not an assertion of globally disjoint cloth.

## Evidence and scope

Compact report in [machine evidence](head_032_machine_evidence.md), SHA256
`3215d3d4fcdfcca033eb0314bf32a5337a3970f9dc8ac5d2e4394eff81d878db`.
Separate neck/seam/cloth witnesses: `head_032_technical_witnesses.json`,
SHA256 `dac851184ec927bcc285b89ac836db1410cf029ef4fb9949211020152bf96e33`.
Outer-contact and hood-cutoff clarification:
`head_032_contact_clarification.json`, SHA256
`cfc10ed05967eb894ad9b3b5f8c79f6c8bb016d3151b78daaaddb3e514c765ea`.
The compact report and clarification are included; large separate raw
witness rows remain ignored and hash-bound, not claimed stored canonically.
No global contact, self-intersection, animation or acceptance scan was added.
