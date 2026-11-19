# Skirt 035 preflight: exposed seat dominates; no global resize

No broad skirt or body resize is supported by this measurement. The red
skirt's width already agrees with the canonical landmark and its front
height overlaps the photographic uncertainty. The repeated cone complaint
is a construction/readability hypothesis, not proof that these bounds fail.
The new first-hit evidence identifies exposed internal seat geometry as the
main dark underside band. An interior front-lap/rear-drape revision with
the waist and red-to-white hem seam fixed remains a separable wall-shape
hypothesis, but cannot directly fix that dominant exposure. Defer or revise
it as the next construction choice. Root owns the decision after 033/034;
this preflight makes no model change.

## Frozen evidence and reference limits

Retained 032 input: `head_032_candidate.blend`, SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
Viewed all five existing actual-material images in
`head_032_eevee_review/`, the canonical front and physical side. Fixed
candidate views are orthographic, scale 0.292 m, 512², at frame 1; existing
Eevee 16-sample materials/lights/color settings are unchanged. No new render,
model save, construction helper, goal or canonical file was changed.

Canonical front controls width and frontal garment placement. Its existing
datum is Wh=368±4 px; candidate hood Wh=117.439255 mm, or 205.920886 pixels
in the fixed camera. The physical side is oblique and partially occluded;
it controls seated cloth behavior and depth order, not a precise overlay.
The rear model profile is measured below, but no new exact-variant rear
measurement was inferred from that side photograph.

[Annotated reference endpoint observations](skirt_035_reference_landmarks.svg)
mark the waist pick at y=685, x=399–560 and central visible red extent
through y=763. Waist-row ambiguity, sleeve/hem occlusion and perspective
give a conservative waist-width band 0.44±0.05 Wh and center drop
0.21±0.025 Wh.
These are preflight observations, not newly approved acceptance thresholds.
The image's true red sewn hem is hidden by white folds; use the existing
canonical red-width target 0.916±0.05 Wh rather than inventing its hidden edge.

A thresholded connected red component spans only 265 pixels, 0.720 Wh,
because hem folds and sleeves hide/disconnect red areas. It is explicitly
not a replacement for the authoritative full red-skirt width. No aligned
overlay or camera-match pass is claimed by this preflight.

## Narrow measurement scorecard

| Quantity | Reference evidence / uncertainty | Retained 032 |
| --- | --- | ---: |
| Full red skirt width | Canonical 0.916±0.05 Wh | 0.910005 Wh |
| Sewn waist width | Visible pick about 0.44±0.05 Wh | 0.462937 Wh |
| Front-center waist-to-red-hem drop | About 0.21±0.025 Wh | 0.231584 Wh |
| Rear-center waist-to-red-hem drop | Not precisely recoverable from supplied side | 0.404553 Wh |
| Full red depth | No new precise physical-side target | 1.005654 Wh |
| Full red vertical extent | Rear includes low seated drape | 0.408600 Wh |
| White Hem026 width / depth | Separate component; not a red resize target | 1.050792 / 1.185612 Wh |

The canonical photograph supports a short front lap, much lower rear skirt
and localized seam compression. Current front/rear red hems are respectively
Z=26.736/6.423 mm; their 20.313 mm difference is 0.173 Wh. Thus 032 has
seated front/rear height asymmetry and is not literally a symmetric cone.
The white hem depth is separately larger than the existing 0.99–1.15 Wh
combined skirt/hem construction band; that does not justify shrinking the
red surface or moving toes in this proposed unit.

Representative evaluated control-row profiles are shown in millimeters.
Rows retain the original waist-to-hem parameter t; they are not equal-height
cross sections. The full Subsurf surface supplies the bounds above.

| t | Width / Wh | Front Y,Z mm | Rear Y,Z mm |
| --- | ---: | ---: | ---: |
| 0 | 0.463 | -34.500,53.933 | 25.500,53.933 |
| 0.1875 | 0.560 | -40.888,45.054 | 33.137,42.341 |
| 0.3542 | 0.732 | -45.420,38.075 | 38.546,30.529 |
| 0.5208 | 0.805 | -49.645,33.189 | 43.574,20.147 |
| 0.7292 | 0.853 | -54.685,28.866 | 49.551,10.345 |
| 1 | 0.909 | -61.000,26.736 | 57.000,6.423 |

## Actual representation and dominant separable issue

The visible red object is `Skirt022 joined gathered panels`: 7,840 base
vertices / 7,680 faces, 49 radial rows × 160 angular columns. Despite its name,
it is one periodic loft from a rounded rectangular waist to a deep hem,
not separately shaped front and rear fabric panels. `body_skirt_022b.py`
uses radial interpolation t^0.84, monotone cubic front/rear drape heights,
five 1–1.8 mm Gaussian gathers and baked seat/leg support. That version
already removed 022's stationary-row terrace bug; do not rediscover it or
retry that repair.

Current modifier order is Subsurf 1 → Solidify 0.65 mm, centered → Armature →
`022 body proportion cage` → `023 narrow waist field`. The 023 field scales
upper X by 0.69, restoring full width below Z=28 mm; it changes the apparent
flare without giving the red surface new panel or seam structure. Most
width growth occurs between the first two-thirds of the drape; the visible
middle wall remains a very regular bilateral flare with weak localized
compression. That is a plausible contributor to the smooth lampshade read.

`Hem026 curled cotton strip` is a different mesh, generated from the 025
ruffle sampler and sewn to the evaluated red boundary. Its white pointed,
repetitive edge is visible in the retained images and cannot be corrected
by claiming a red-panel edit fixes it. Black foot-pod visibility and any 033
change are also separate. Do not combine these three complaints into one
global cone/toe diagnosis.

## First-hit underside attribution changes the priority

The broad black/red band is not simply a skirt-width issue. Pixel-center
rays from the exact saved front-camera contract identify these surfaces:

| Front-image location | First visible geometry |
| --- | --- |
| Row y=454, x=215–296 | `Garment42 compact internal seat pad` |
| Row y=460, x=217–294 | Same internal seat pad |
| Row y=466, x=217–294 | Far/rear surface of `Skirt022 joined gathered panels` |
| Row y=471, x=216–295 | Same far/rear red skirt surface |

The seat hit at pixel (250,454) is world Y=−38.128 mm, Z=16.793 mm,
material `Dress red cloth.004`, with actual decoded image RGB
(0.0235,0,0). Thus the near-black band is red seat fabric rendered very dark,
not a separate black foot or trim material. At pixels (250,466)/(250,471),
the lower red stripe hits the rear skirt at Y=+49.464/+52.646 mm and
Z=9.949/7.097 mm. The camera sees through the front opening to that far wall.
Foot pods bound this band laterally; white Hem026 folds cross its top and
appear again beneath the lower red edge.

In the explicit x=214–297, y=438–474 ROI, 3,108 first-hit rays attribute
1,824 pixels to the seat, 656 to red skirt, 470 to white hem, 92 to foot
pods, and 66 to no surface. These ROI counts are not whole-image scores.
Rays identify opaque geometric surfaces; they do not simulate lighting,
transparency, antialiasing or subpixel coverage. Actual PNG pixels corroborate
the sampled dark seat and lower red stripe. No shading cause was repaired.

An interior-only red drape test that preserves the seat and present front
seam cannot directly remove the dominant exposed-seat surface. Concealing
it would require a newly authorized coverage/interface change, not merely
adding red-wall folds. The rigid white collar is also a separate target and
was not inspected or altered here. These visible differences take priority
over an argument based only on acceptable red width. Root should not select
red interior drape as a claimed fix for this underside band.

## One bounded hypothesis for root's decision

If the red wall still reads rigid after the pending units, test one
panel-defined seated compression field on the red interior only: distinguish
a broad front lap from the lower rear drape, with a few seam-rooted broad
folds and unequal side compression. Use the existing cloth topology or a
replacement interior joined to the same boundary rings. The mechanism is
panel tension and contact-shaped drape, not more radial oscillation, random
noise, global smoothing, uniform scaling or surface stitches. This is a
construction hypothesis; the photograph does not determine a unique 3D fold
amplitude, so this preflight does not prescribe one as measured truth.

Exact proposed direct target: `Skirt022 joined gathered panels` only.
Keep both evaluated midsurface boundary loops and enough adjacent rows to
preserve their tangents. Waist: 320 vertices, perimeter 207.814983 mm,
Z=53.933022 mm; lower red seam: 320 vertices, perimeter 378.648855 mm,
Z=6.423603–28.268075 mm. Their hashes are in the probe. The lower loop is
the receiver sampled by Hem026's 768-point sewn root. Do not change that
receiver, Hem026, bodice, internal seat, legs, black pods, sleeves, rig or
shared 022/023 cages. New foot work remains protected input, not an excuse
to regenerate historical support geometry.

Expected view effect: front and both three-quarter views gain localized
seated panel compression without width/height expansion; side/rear lose
the uniform middle-wall sweep while retaining existing low rear and raised
front endpoints. The pointed white edge is expected to remain unchanged.
Rebind these guards to whichever exact model root retains after 033/034;
032 measurements cannot silently authorize modifying newer geometry.

## Strongest rejection and falsifier

The strongest objection is that width and front height are already within
their evidence bands, and the rigid impression may be dominated by the
exposed seat, white edge and rectangular bodice. A red-only fold revision could add noise
without improving likeness, or trade a smooth skirt for exaggerated pleats.
Holding both boundaries also limits how much its outer silhouette can
change. This evidence argues against a wholesale resize/rebuild by
impression and against asserting that red geometry is the sole cause.

Provisional recommendation to root: do not resize globally, and defer a
red-interior test as the next fix for the dominant underside complaint.
Keep that bounded hypothesis only for a separately persistent rigid-wall
read after coverage has been addressed. Keeping the present red surface is
a credible alternative. Root owns the final proceed/revise verdict and
any new interface scope; this preflight authorizes no construction trial.

Reject the test if fixed side/three-quarter pixels retain the same stiff
red-wall read, gain an artificial pleat/balloon/shelf, or require shifting a
protected seam to appear improved. Reject new hem detachment or visible
clipping, silhouette movement beyond the existing uncertainty, or a claimed
improvement attributable only to altered white hem/toes/materials. Technical
seam preservation cannot pass this visual question; an implementation-blind
review must judge the red wall separately before any module retention.

Measurement evidence: `skirt_035_probe.json`, SHA256
`5057ae16a37e7c5c7af584c7e4e6ee764e08f4eac974cfb401252dfed51726f1`.
First-hit evidence: `skirt_035_underside_first_hits.json`, SHA256
`8d9437f11c5160b5256df53d9eaa7db3ecdcc38700890664f97e47427c376442`.
Two pinned read-only processes measured profiles, image-row aids and the
newly requested underside attribution; temporary in-memory diagnostic
copies were discarded. Input 032 hash remained exact in both processes.
This is a scoped, implementation-aware preflight, not a blind review or
whole-character scorecard/acceptance gate. Canonical links are root-owned.
