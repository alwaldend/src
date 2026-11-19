# Head attempt 019: reset volume, preserve contour evidence

The first complete saved state materially removes the donor's raised angular
forehead plates and wide lateral bulges. It does not pass the frozen module
gate: the side volume still reads as a deep cylinder/helmet, the front-to-side
transition is too abrupt, and the new lower head exposes a neck band above
the unchanged collar. No module, stage, or final asset is accepted.

Candidate: `out/reimu_fumo_finish/desktop_astra/head_019_candidate.blend`,
SHA-256 `b770864c3015561a7e3c780590ceb13f85307086d996bd942000bc9668af384a`.
Latest fixed review render:
`out/reimu_fumo_finish/desktop_astra/head_019_review/candidate_three_quarter.png`,
SHA-256 `b5e938585165bffcdfd9ec1cb0751f28a34a691e8800cabc73e30832a60b1a1b`.
These are rejected refinement evidence, not the reusable deliverable.

The tracked donor SHA remains
`a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5`.
The appended pre-edit snapshot SHA is
`dd1994a210cb65e4cdeaa4d42b0227b0a002524f9baef6c8da030fe9bee24fe5`.
Source bytes and all 33 evaluated visible non-target control meshes remain
unchanged, independently checked after pinned reopen.

## Technical result

Reviewer `head_019_technical`, pinned Blender 5.2.1 LTS build 9e2066aef7ef:
both new meshes have finite coordinates, no degenerate faces, and no
zero-length edges. Cushion and evaluated solidified fringe are closed
manifolds with Euler characteristic 2. The fringe's unsolidified base has
340 expected boundary edges. Normalized cushion dimensions are width
1.00206 Wh, depth 0.74677 Wh, height 0.98382 Wh. Passing depth bounds did not
prevent the visible cylinder read: the distribution of depth is defective.
An inward-wound left rear lock is inherited from the donor, not newly caused.
This module audit does not certify whole-asset dependencies or rigging.

Audit JSON SHA:
`1c977daba567d28b5ce81ad2f814c7145714adb904c871467f77e71e53c46aac`.
Render receipt SHA:
`18974749d083ca36041e9dd65b47b8d18625282e9a4876c0ed152c9445f3472f`.

## Independent visual review

Reviewer `head_019_blind_01`, independent image-only reviewer, inspected
candidate front, side, rear, both three-quarter views; donor front and
three-quarter; canonical_front_25cm, physical_front, and physical_side.
The candidate review preceded baseline comparison. No scripts, writer
receipts, or other judgments were supplied. No controlling rear/turntable
frame or calibrated overlay was reviewed; landmark tolerances are unverified.

The reviewer recommends refine and recognizes a substantially improved
crown, integrated broad fringe, head/bow proportion, and removed lateral
protrusions. Regression-free retention is withheld for the neck/collar gap.
Remaining major defects: drum-like side volume, strong shell transition,
blade-like locks with weak roots, and insufficient cloth layering.
Absolute scores: likeness 5/10; silhouette/proportions 5; construction 3;
identity features 5; contact/occlusion 4; intended-medium read 3;
presentation 7. Absolute stage rejected.

Coordinator verdict: reset the volume subsystem. The reviewer's refine
recommendation supports preserving the traced contour as evidence, but the
frozen plan requires reset when a major helmet or attachment defect remains.
Do not iterate material detail or disguise this with lighting. The direct
control-mesh route worked; it is the uniform annular side-volume distribution,
not window transport or inability to edit, that needs a causal change.

## Next authorized work

Goal remains open/active with no accepted attempt. Before another candidate,
replace the approximately constant-depth annular wall with a physically
motivated, spatially varying stuffed-panel profile and resolve lower-head
contact against the unchanged collar. Preserve the useful measured front
contour as an input, not accepted geometry. Freeze a new attempt and test
side/three-quarter construction and collar contact first. If those fail,
stop that state before detail work. Recheck writer bindings before reuse.

Window screenshots now work; see screenshot_evidence.md. No popup is visible.
Both runtimes were tried historically and in this continuation; the working
split remains Flatpak authoring and pinned saved-byte rendering. Desktop Use
has no user preference and is not an acceptance criterion.
