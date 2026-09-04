# Lower-face image diagnosis

Observed 2026-09-05 19:51 UTC. Diagnosis only; no acceptance decision. Reviewed only candidate `chin_073_close_review/chin_close.png`, `chin_073_fast_review/{front,side,three_quarter}.png`, references `canonical_front_25cm.png`, `physical_front.png`, `physical_side.png`, and turntable frames `side_062_frame_{06,08,10,12}.png`. No model, scripts, prior candidates, or implementation plans inspected. The user’s recession complaint was supplied, so this is not a completely unprompted review.

The strongest confirmed mismatch is the narrow lower facial panel. The candidate’s cheek edges converge under the eyes, making the mouth and rounded bottom read as a small inward-tucked wedge. The close and three-quarter views show this through the cheek-to-mouth region, not solely on the underside. Physical-front controls the stuffed construction; canonical-front controls the clean graphic proportions. Their variants are kept separate.

Manual projected measurements use maximum hair-cap width, excluding bow and hanging locks, as Wh: candidate approximately 211 px; canonical 359 px; physical-front 178 px. Pixel selection uncertainty is approximately 3–6 px; combined normalized uncertainty below includes occlusion. Cameras are not calibrated, so these are bounded diagnostic estimates, not overlay acceptance measurements.

| Visible front landmark | Candidate | Canonical | Physical-front |
| --- | ---: | ---: | ---: |
| Exposed face width at mouth height | 0.50–0.53 Wh | 0.57–0.62 Wh | 0.59–0.65 Wh |
| Mouth-to-visible-bottom distance | 0.10 ±0.02 Wh | 0.11 ±0.02 Wh | 0.11 ±0.02 Wh |
| Eye-top-to-mouth distance | 0.27–0.28 Wh | 0.28–0.30 Wh | 0.29–0.32 Wh |

Hair overlap contributes to exposed width; the images alone cannot separate hidden panel width from framing. Nevertheless, the visible taper needs correction. Merely lengthening the chin is unsupported by the already similar mouth-to-bottom distance.

In candidate side and frame 10, the visible face-edge run below the eye is roughly 25–40 degrees from image vertical, with approximately 10-degree uncertainty. This is an oblique, partly occluded contour, not a face-plane normal. Physical-side and frames 10/12 hide the actual chin with hair; none establish its true depth. Frames 06/08 support broad cheek-to-bottom continuity but do not provide a precise forward-displacement target.

The candidate sits close above the collar, with only a thin dark strip, approximately 0.02 Wh. Front eye-to-mouth orientation is broadly compatible. Neither whole-head rotation nor body seating is established as the primary cause.

Construction-level hypothesis: reform the lower face as a broad shallow stuffed front panel continuing from below the eyes through the mouth, retaining approximately 0.60 Wh visible width at mouth height, then turning underneath through one soft rounded seam. Avoid a local chin bulge: it could add a ledge while preserving the recessed-looking cheeks. Validate depth and hair overlap in matched three-quarter views before committing a displacement.
