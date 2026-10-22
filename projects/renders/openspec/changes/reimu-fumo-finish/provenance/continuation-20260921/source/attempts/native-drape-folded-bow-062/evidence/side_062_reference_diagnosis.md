# Side reference diagnosis for cycle 062

Observed 2026-09-05 16:57 UTC. Independent pixel diagnosis; no builder,
geometry, topology, or model scripts inspected. The frame sheet covers all 30
frames of the GIF; measurements are manual pixel estimates with stated limits.
No reference/candidate camera match or overlay has been established.

## Sources and normalization

- Controlling front: `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png`, SHA-256 `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
- Controlling turn: `projects/renders/assets/reimu_fumo/references/canonical_turn_180.gif`, SHA-256 `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`.
- Construction support: `projects/renders/assets/reimu_fumo/references/physical_side.png`, SHA-256 `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`.
- Rejected candidate comparison: `out/reimu_fumo_finish/desktop_astra/macro_061_fast_review/side.png`, SHA-256 `2eb466d7b8504c58cebc38fc1b80f15d60258493d608201c304ede7f1e5aed66`.

Use zero-based GIF frame 12 as the most useful near-pure side head view,
facing left. Frame 11 retains a narrow eye/face view; frame 13 has moved into
rear-side. Foot overlap suggests the torso axis is nearer pure side around
frame 13, so the head and body do not establish one precise yaw angle.
Opposite-side frames 26/27 bracket the other side; neither reveals the entire
stuffed core. Frame 16 is rear three-quarter.

Define H as projected crown-to-chin/lower-head distance, excluding bow and
loose hanging hair. For frame 12, crown y≈65±8 and chin y≈296±6 give
H≈231±12 px. The crown is partially hidden, so normalized values are not
high-precision measurements. Candidate crown y≈144±2 and lower head
y≈331±8 give H≈187±10 px; physical support H≈200±14 px.

## Head envelope, frame 12

Image coordinates increase right/down. At this view, front is left. Normalize
horizontal coordinates from x=169, the foremost visible cap edge, by H=231.
This table distinguishes observed exterior hair/face from the hidden stuffed
core. It is not a claim that the stuffing reaches every listed exterior point.

| Height below crown | y px | Front envelope x px | Rear close cap/hair x px | Projected cap depth / H | Confidence and occlusion |
| --- | ---: | ---: | ---: | ---: | --- |
| 0.10 H | 89 | 191±3 | 330–358 | .60–.72 | Front cap measured; rear concealed by wing/tail and rear hair root |
| 0.25 H | 124 | 177±3 | 338–358 | .70–.78 | Front cap measured; rear is an estimate, not directly separable |
| 0.51 H | 183 | 169±3 | 343–358 | .75–.82 | Front cap measured; rear ribbon/flap boundaries merge |
| 0.77 H | 242 | 178±4 | 350–355 | .74–.77 | Front fringe and rear close-hair exterior visible; stuffing lies inside |
| 0.92 H | 277 | face ≈218±10 | 351–356 | .53–.64 | Front braid at x194 is separate hanging hair; rear includes skirted hair below core |

The front cap is steep through its middle: x177 at y124 to x169 at y183,
then x178 at y242. Its small variation is approximately .04H over .52H
vertically. That is a rounded side plane with a short crown roll, not a large
forward cheek bulge. At the lower row, do not connect the front braid x194 to
the stuffing contour; the beige face is farther back at approximately x218.

The large rear pointed piece reaches x≈444 at y≈304 and is a hanging hair
panel. The nearby head/close-hair rear lies around x350–356. Counting the
pointed piece as head depth adds roughly 90 px, or .39H, falsely. The red
vertical tail at x≈232–349, y≈76–266 hides most of the side head surface.
The exposed rear edge below it is hair, not an observable stuffing seam.

Candidate 061's middle cap spans approximately x169–327, or .84H using
H≈187. Reference middle close cap is approximately .75–.82H, with rear
occlusion. These ranges support at most a modest bulk-depth discrepancy;
they do not support a large global flattening prescription. The physical
photo also has appreciable depth, partly hidden by the tail. The stronger
head failure is construction: the candidate has a continuous smooth shell
and thick rounded lower lobes, while the references expose thin cut/folded
hair layers with distinct hanging pieces. Changing depth alone will not
establish sewn construction.

## Upper bow, tails excluded

An independent second pixel reader agreed on these approximate bounds.
Boundary uncertainty is generally ±4–8 px, plus the H uncertainty above.

| View | Upper-wing bounds in original pixels | Span / H | Total visible height / H | Rise above crown / H |
| --- | --- | ---: | ---: | ---: |
| GIF frame 11 | x190–367, y17–113; H236 | .75 | .41 | .19 |
| GIF frame 12 | x166–344, y18–114; H231 | .77 | .42 | .20 |
| Opposite frame 26 | x112–322, y27–157; H≈250±20 | .84 | .52 | .12 |
| Candidate 061 | x257–312, y89–216; H187 | .29 | .68 | .29 |
| Physical side, both wings | x129–295, y7–154; H200 | .83 | .74 | .53 |

The GIF side upper wing has a broad principal chord approximately 20–30°
from horizontal. Candidate 061's upper wing is approximately 80–90°. The
side horizontal span is about 60% smaller relative to H than the GIF.
The small candidate knot at x247–281,y135–150 is separate and should not be
used to inflate the wing's width. The reference tail is also separate from
the wing: frame12 approximately x232–349,y76–266, candidate approximately
x298–318,y190–311. Folding the tail into the wing bounding box hides the
orientation error.

The physical photo shows a different lifted drape: the near wing occupies
x129–242,y50–154 while a far wing rises almost vertically to y7. It supports
thin folded, asymmetric cloth. It does not establish the canonical bow pose.
Pixels cannot decide whether the difference is a manufactured variant,
manual adjustment, or another deformation.

## Front-versus-turn pose conflict and camera limitation

Canonical front upper wings occupy approximately x211–758,y150–378.
With crown232/chin589, H357, their width is 1.53H, total height .64H, and
above-crown rise .23H. Closest GIF front frame04 upper wings occupy
approximately x25–383,y25–160; H239 gives width1.50H, height.56H, rise.16H.
The upper widths are compatible, but the GIF wings sag and fold differently.
The tails differ more: the canonical front tails spread diagonally outward,
while GIF tails hang downward. Braids also appear longer/narrower in the GIF.

Both photographs have perspective/elevation and expose top-facing surfaces.
Lens, camera elevation, per-part yaw, and exact bow pose are unavailable.
Within the GIF itself, the visible wing height changes between front and
side. That can include projected depth at an elevated camera and occlusion
of lower wing regions by the head. A pure orthographic side/front pair would
share full assembly Z extrema, but these are visible image bounds rather
than recovered full assembly bounds.

Therefore .42H is NOT an exact physical Z-span target for the model. Do not
shrink the upper-wing Z span to force that number while changing the
controlling canonical-front silhouette. The robust change supported by both
side references is more front-to-back spread and a folded/rolled wing surface
over the crown, with the appropriate head occlusion. Preserve front X/Z
landmarks, then test the resulting visible side outline. Any exact side
overlay requires a camera match first.

## Diagnostic artifacts

- `side_062_frames.jpg`: all 30 GIF frames, zero-based labels.
- `side_062_frame_00.png` through `side_062_frame_29.png`: extracted frames.
- `side_062_grid_frame11.png`, `side_062_grid_frame12.png`, `side_062_grid_frame26.png`, `side_062_grid_physical.png`, `side_062_grid_candidate.png`: coordinate witnesses, 2× display with original-pixel labels.
- `side_062_head_landmarks.png`: observed front points and uncertain rear intervals, explicitly no camera-matched overlay.
- `side_062_analysis.py`: task-local extraction and annotation script.
