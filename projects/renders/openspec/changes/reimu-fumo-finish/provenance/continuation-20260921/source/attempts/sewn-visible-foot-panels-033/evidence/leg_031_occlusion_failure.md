# Leg 031 occlusion failure

## Verdict

Reject 031 and retain 030b. The root warp enlarged hidden geometry rather
than creating a readable pale foot section. Its new projected extrema are
first-hit by the black pods, Hem026, the red skirt, or the internal seat.
The small increase in visible cream is not a coherent reference cue.

Decision-review verdict on a material-region replacement: **revise**. Moving
cream identity onto the retained pod shell targets the correct first-hit
surface, but no smooth existing pod loop produces the required upper-inner
cream patch while keeping the distal and ground-contact surface black. A
plain existing-face recolor is therefore not the next candidate.

The exact future causal choice is one two-material sewn foot shell derived
from each 030b pod's unchanged outer envelope and support contacts. It needs
one purpose-built closed seam boundary, with a black distal/bottom toe cap
and a cream upper-inner/proximal panel. Reuse the existing black and cream
materials without node edits. Do not move the feet, change Hem026, warp the
separate hidden roots again, or construct another tube.

This is a future representation decision only. No candidate or material
variant was made; head work can continue from retained 030b.

## Frozen evidence

- Retained source: `bow_030b_candidate.blend`, SHA-256
  `d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6`.
- Rejected candidate: `leg_031_candidate.blend`, SHA-256
  `ece85247dc07e9ac59388c20321b992c8638d4f3294ac4d4ef6436e975489b71`.
- Baseline renders: `bow_030b_eevee_review/`.
- Candidate renders: `leg_031_eevee_review/`.
- Review-contract SHA-256:
  `4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
- Both packets use the contract's 512 by 512 orthographic cameras at
  `0.292 m` scale and the same Eevee, lighting, and color settings.

All `.blend` hashes were rechecked after inspection and remained unchanged.

## Measurement method

For each fixed view, an isolated evaluated-root BVH defined every pixel-center
ray that would hit either cream root if nothing occluded it. A second BVH
contained all 86 visible evaluated geometry objects and identified the actual
first-hit object for those rays. Counts therefore distinguish projected bounds
from renderable surface. They do not count antialias edge coverage; boundary
uncertainty is approximately one pixel. The relevant materials are opaque.

Decoded RGB comparisons independently counted pixels whose largest channel
changed by more than 4/255. This includes shadow changes and is an upper bound
on newly readable cream, not an object-ID mask.

## Visible-root coverage

| View | 030b projected | 030b cream | 031 projected | 031 cream |
| --- | ---: | ---: | ---: | ---: |
| Front | 1,210 | 81 (6.69%) | 2,412 | 199 (8.25%) |
| Side | 805 | 265 (32.92%) | 1,085 | 238 (21.94%) |
| Three-quarter | 1,422 | 370 (26.02%) | 2,122 | 435 (20.50%) |
| Rear | 1,210 | 0 | 2,412 | 0 |
| Mirror three-quarter | 1,422 | 316 (22.22%) | 2,122 | 360 (16.97%) |

In front, 031 doubles the root footprint but adds only 118 visible cream
pixel-center rays. The final 199 rays occupy 0.076% of the complete frame;
the increase is 0.045% of the frame. In side view, the visible count falls by
27 even though the projected footprint grows by 280. The visible fraction
also falls in both three-quarter views.

The actual PNG changes are similarly small:

| View | RGB pixels changed by more than 4/255 | Frame fraction |
| --- | ---: | ---: |
| Front | 174 | 0.066% |
| Side | 175 | 0.067% |
| Three-quarter | 380 | 0.145% |
| Rear | 0 | 0% |
| Mirror three-quarter | 259 | 0.099% |

The decoded rear RGB images are identical even though their PNG byte hashes
differ. In front, all changed pixels lie within `x=175..323`, `y=436..453`,
an 18-pixel-high foot/trim junction. This agrees with the blind finding that
the tiny light areas remain ambiguous with white trim.

## What first-hits the 031 root footprint

| View | Cream | Black pod | Hem026 | Red skirt | Internal seat |
| --- | ---: | ---: | ---: | ---: | ---: |
| Front | 199 | 1,103 | 513 | 400 | 197 |
| Side | 238 | 251 | 194 | 391 | 11 |
| Three-quarter | 435 | 749 | 496 | 227 | 215 |
| Rear | 0 | 0 | 0 | 2,412 | 0 |
| Mirror three-quarter | 360 | 674 | 742 | 240 | 106 |

The black pods are the largest front and three-quarter occluder. The red skirt
is the largest side and rear occluder, while Hem026 dominates the mirrored
three-quarter. This is not a single-object obstruction that another scalar
root enlargement can clear safely.

The projected extrema make the failed mechanism more explicit:

- Front: both 32-pixel top extrema are 100% red-skirt first hits. The left
  inward edge is seven seat-pad and two hem pixels; the right inward edge is
  eight hem and one skirt pixel. Both 28-pixel bottom extrema are black pod.
- Side: both 21-pixel top extrema and both 18-pixel seatward extrema are 100%
  red skirt. Both 27-pixel bottoms and 16-pixel distal edges are black pod.
- Three-quarter: the top extrema are red skirt. The new left inward edge is
  split six hem / six seat-pad pixels; the right inward edge is black pod.
- Rear: every one of the 2,412 projected root pixels is red-skirt first-hit.

Thus the successful 031 bounds measurements described geometry behind other
surfaces. They did not measure the final image landmark that mattered.

## Pod material-boundary feasibility

The retained pods are identical in 030b and 031. Each is a closed manifold
UV-sphere-family mesh with 1,962 vertices, 3,976 edges, and 2,016 faces:
1,904 quads, 112 pole triangles, 56 longitude segments, 35 intermediate
latitude rows, and two 56-valence poles. Each currently has only
`Feet black velour.002`; `Dress warm white cloth.002` already exists and can
be reused. Distal/front is world `-Y`, proximal/seatward is `+Y`, and the pod
midpoint is approximately `Y=-0.060 m`.

Three topology-respecting material splits were tested with 120 by 120
pixel-center rays over the pod footprint:

1. The clean `qY=0.5` proximal split is one planar 72-edge, degree-2 cycle
   through both poles and divides the faces 1,008 / 1,008. It yields zero
   cream front first hits out of 11,686 and makes 50% of ground first hits
   cream. The front hemisphere completely hides the cream rear hemisphere.
2. A mirrored 45-degree proximal/inward meridian loop exposes 14.66% cream in
   front, but still makes about 50% of the ground cream and creates severe
   `+X` side asymmetry: 85.3% left versus 14.7% right.
3. Clean 56-edge upper-latitude loops at evaluated `qZ=.695` and `.736` keep
   the ground black and expose 24.95% and 19.50% front cream, respectively.
   They color the whole upper toe rather than an upper-inner/proximal panel.

No existing smooth edge loop satisfies all three needs: front visibility,
black distal/ground coverage, and comparable bilateral side behavior. A
thresholded arbitrary face set would be a jagged color boundary, not a sewn
panel.

The original cream roots begin at pod-normalized `qY=.584/.585` and continue
beyond the pod's proximal end. That supports the intended proximal cream
identity, but 031 proves that a separate rearward root does not put that
identity on the visible surface. The physical side likewise shows cream
continuing behind and above a black distal toe. It is compatible with one
continuous stuffed shell and a real sewn cap boundary, but not evidence that
a hidden rear half or an arbitrary painted patch is sufficient.

## Exact next causal unit and gates

When feet are revisited, start from 030b and change only the two pod shells'
panel topology and material regions. Use each evaluated pod as the immutable
outer-envelope and contact donor. Add one purpose-built, mirrored closed seam
on that same surface; keep the black panel over the complete distal center,
outer/lower toe, and every ground-contact first hit. Put the existing cream
material on the upper-inner/proximal first-hit surface. The seam must be a
single manufactured boundary, not a shader gradient or face-threshold stair.

The unit passes only if:

1. Pod placement, overall rounded silhouette, `0.3995-0.4073 Wh` depth,
   ground contact, and the zero-crossing Hem026 support relationship remain
   unchanged within the existing `0.003 Wh` contact uncertainty.
2. Black and cream use the existing persistent material datablocks with no
   node changes. The panel boundary is one closed degree-2 loop, with no
   disconnected islands.
3. Every ground and distal-center first hit remains black. Both sides produce
   comparable material coverage rather than the tested 85.3% / 14.7% split.
4. The front material-ID mask contains one contiguous cream region per foot,
   visible about `0.06-0.09 Wh` inward and `0.04-0.07 Wh` above the black cap.
   Bounds behind hem, skirt, seat, or toe do not count.
5. The black cap's projected width is approximately 85% of the preserved pod
   width. This follows from the reference black widths (`0.209/0.215 Wh`)
   versus current pod widths (`0.247/0.250 Wh`) and remains subject to the
   recorded soft-edge uncertainty.
6. Fixed front, side, both three-quarter, and rear pixels show a sewn cream
   upper section without a painted halo, tube, detached bulb, new clipping,
   or loss of the retained toe/hem contact. Independent blind review, not the
   material-ID mask alone, decides retention.

If a conformal seam cannot satisfy these gates without changing the pod
envelope, stop and reconsider the whole two-panel stuffed-foot representation;
do not resume root-warp or placement sweeps.

