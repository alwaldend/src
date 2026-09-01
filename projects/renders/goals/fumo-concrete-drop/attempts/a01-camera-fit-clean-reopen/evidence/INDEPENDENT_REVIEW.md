# Concrete A01 independent review

## Verdict

`REFINE`. Candidate
`04dda70fc072c44e6eb0820eb0d159c6b9e6f08b552f049b0efa0916e9cf7b1c`
passes the interface, neutral-placeholder, scale, timing, mechanics, semantic-
delta, and preservation checks. It does not pass the required pixel framing
criterion.

At frames 22 and 28 the proxy intersects or obscures the orange warning text.
At frames 40, 56, and 72 the warning text crosses the settled proxy, whose
gray-on-gray silhouette is already weak. The cyan impact marker is also partly
covered around impact and settling. In-frustum bounding boxes do not establish
mutual non-occlusion, label legibility, or silhouette separation.

## Criterion verdicts

- 001 pass: pinned Blender 5.2.1 clean-opened the exact interface and custom
  properties.
- 002 pass: all renders and the inventory show a neutral, claim-limited proxy.
- 003 pass: metric scale, 0.250 m collider, 0.245 m proxy, and floor bounds
  were independently reproduced.
- 004 pass: the 24 fps hold, release, contact, and settle timeline reproduced.
- 005 fail: required objects are inside frame but overlap in the actual pixels.
- 006 pass: 0.845343 m descent, frame-22 contact, -0.000687 m minimum bottom,
  and zero late span reproduced.
- 007 pass: every manifest hash recomputed; the direct semantic diff contains
  only the 45 to 35 mm lens change.

## Next work unit

From this exact candidate, change only the location of
`Placeholder Warning Label` so both lines remain completely outside the
sampled proxy and impact-marker trajectory. Keep camera, physics, proxy, marker,
materials, and all other properties unchanged. Rerender the same eight frames
and require unobstructed text plus a clearly separable complete proxy at frames
22, 28, 40, 56, and 72.
