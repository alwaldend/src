# Concrete A02 label-clearance plan

## Target

Preserve the A01 camera-fit and mechanics while moving the two-line warning
completely outside the projected proxy and impact-marker trajectory.

## Frozen input

`out/fumo_concrete_drop_scaffold/attempt_03_camera_fit/fumo_concrete_drop_camera_fit.blend`
at SHA-256
`04dda70fc072c44e6eb0820eb0d159c6b9e6f08b552f049b0efa0916e9cf7b1c`.
The A01 input and tracked Reimu/Sisyphus assets remain protected.

## Bounded work

Project the proxy and impact-marker bounds at frames 1, 12, 20, 22, 28, 40,
56, and 72. Compute one evidence-based translation for the object
`Placeholder Warning Label` that places both orange lines outside the union of
those bounds while keeping the complete text inside frame. Change only that
object's location, reset to frame 1, and save a new candidate. Do not change
camera, physics, animation, proxy, marker, label text/material/scale/rotation,
lighting, or interface.

## Outputs and gates

Clean-open without resaving. Emit the exact one-property semantic delta,
mechanics audit, protected hashes, the same eight renders, a labeled contact
sheet, pixel-overlap measurements, and result. Stop on any non-location delta,
mechanical regression, protected-hash change, image-boundary contact, text
occlusion, label/proxy or label/marker overlap, or weak complete proxy
silhouette at frames 22, 28, 40, 56, and 72. No second translation in this
attempt.
