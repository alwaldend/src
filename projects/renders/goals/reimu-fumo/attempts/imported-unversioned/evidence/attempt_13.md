# Reimu Fumo attempt 13

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 13 — shallow-domed continuous cushion drum

**Candidate:** planned
`out/reimu_fumo_attempt_013_head_cushion_drum/checkpoint_setup.blend`, SHA-256
pending, untouched neutral seed stage, review packet
`attempt-013-continuous-drum-seed-fixed-views`.

**Failure targeted:** Attempt 12 encoded an eight-sector front outline and a
separate bevel belt. Voxel remeshing preserved the resulting chords, wedges,
hard cap-to-rim shoulder, and lower-center ridge, so the object failed before
sculpting.

**Hypothesis:** A dense superelliptic panel boundary, two broad shallow-domed
caps, and a half-ellipse rim that is tangent-continuous at both caps will
provide one smooth stuffed-cushion blank without the planar mass of a cube,
the chords of a low-sided cylinder, or the radial curvature of an ellipsoid.
One subdivision and one fine voxel remesh will erase the disposable generator's
topology while preserving its continuous occupied volume for later native
Grab sculpting.

**Representation-scope correction:** Attempt 11 prohibited scripted
superellipsoids and profile formulas because its accepted surface was supposed
to be direct clay and prior generators were being presented as sculpt. That
restriction remains for final render geometry and native candidate edits. It
is intentionally superseded for this byte-new Attempt 13 **seed only** after
three direct sharp-cube families and one faceted seed failed. The generated
mesh may only establish smooth disposable clay, must receive one topology-
destroying voxel remesh, must pass an untouched seed pixel gate, and can never
be retopologized, exported, or called the accepted head without later native
sculpt and all full gates.

**Plan written before implementation:**

1. Start factory-empty under a new temporary directory. Import no rejected
   object, mesh, coordinate, modifier, image plane, or scene datablock. Reuse
   only the exact audited render, viewport, checkpoint, and review functions.
2. Sample one 64-sector Lamé boundary with exponent `2.6`, half-width
   `0.455 Wh`, and half-height `0.510 Wh`. Use no eight- or sixteen-sector
   control surface, straight chord longer than one sample, or explicit seam.
3. Build front and rear cap boundaries at `0.82` of the outer perimeter. Fill
   each with eight concentric rings plus a center. Give the front a maximum
   shallow dome of `0.055 Wh` and the rear `0.040 Wh`, using a broad falloff
   that keeps the central half one quiet low-curvature field rather than a
   spherical bulge or flat face plate.
4. Join the cap boundaries with 13 rim rings over a half ellipse. Vary depth as
   `y = -0.235 cos(alpha)` and perimeter scale as
   `0.82 + 0.18 sin(alpha)` for `alpha = 0..pi`. Require the cap normals and
   rim flow to meet without a shoulder, seam, or constant-width gusset band.
5. Build one closed manifold mesh, apply exactly one Catmull–Clark level,
   voxel-remesh once at `0.018 Wh`, normalize to approximately
   `0.91 × 0.56 × 1.02 Wh`, shade smooth, remove every modifier, and save the
   setup checkpoint. No Sculpt Mode stroke is authorized in this attempt
   stage.
6. Render untouched neutral, silhouette, and readable grazing front, side,
   rear, canonical three-quarter, and mirrored-three-quarter views at fixed
   640-pixel cameras. Preserve exact hashes and topology metrics.
7. Hard-reject any resolvable chord, tangent kink, repeating 8/16-sided
   highlight pulse, cap/rim shoulder, seam, flat gusset strip, cleft, notch,
   face plate, straight-wall mattress, rounded box, egg, balloon, helmet, or
   molded-foam read. Require the central `0.50 W × 0.55 H` cap to retain one
   broad quiet tonal field and the three-quarter highlight to flow continuously
   from cap through rim.
8. A context-light reviewer must describe the exact pixels as one shallow
   stuffed cushion blank with two broad panels and one continuously rolled
   edge. Require continuous volume, quiet caps, rolled perimeter, facet-free
   construction, and non-egg cushion read each at least `6/10`, with no major
   defect, before writing or executing a later native-stroke family. Otherwise
   reject the seed and change representation again.

**Attempt 13 work actually performed:** A byte-new dormant driver created the
predeclared 64-sector, eight-cap-ring, 13-rim-ring closed shell. The control
mesh had 1,730 vertices, 3,520 edges, 1,792 faces, Euler characteristic two,
zero boundary or non-manifold edges, and consistent opposite face use on every
edge. One Catmull-Clark level and one `0.018 Wh` voxel remesh were applied,
every modifier was removed, and the untouched seed was normalized to
`0.910 × 0.560 × 1.020 Wh`. No Sculpt Mode stroke ran.

**Attempt 13 raw evidence:** The exact builder, setup checkpoint, seed metrics,
reviewed state, and contact sheet hashes are
`752e3c9725a4fd3f9bb9a6f36b2e64eeb9c14c31386d8b5cf4ad3155014b8ac0`,
`2ff65a7136d67b28b0cbbb1f352d5f7ab0da31658516cd6ca6c40ec60c64e25c`,
`4d5d85a91f4875c27e07700c2836a9e10e6f24d9b36b7485ceb125dae8255cab`,
`652f617cc5c4e3f823527990e1a2d8aed9f7296538d783cd7b6e90741b447ad5`,
and
`80d850831407e3508869a2b033aeb224893b67f81a942968028943f01010b0b0`.
The implementation-blind pixel review scored continuous volume `8/10`, quiet
caps `5/10`, rolled perimeter `5/10`, facet-free `8/10`, and non-egg cushion
read `3/10`. It saw no chord, isolated facet, seam, cleft, or notch, but called
the object a smooth rounded-square slab with broad planar caps, a visible
cap-to-rim shoulder, a narrow constant-depth sidewall, and a uniform bezel.

**Attempt 13 criterion result and decision:** Factory-empty intent, size,
closure, modifier removal, continuity, and facet freedom pass. Quiet caps,
rolled perimeter, non-egg cushion read, reference fidelity, constructed-plush
read, and presentation readiness fail. Reject before native strokes. No
coordinate, mesh, modifier, image datablock, world, collection, or scene
datablock may enter Attempt 14.

### Progress and approach audit after Attempt 13

- **Improved:** the dense outline and fine remesh eliminated Attempt 12's
  visible octagon, wedge facets, lower V ridge, and split. Continuous volume
  and facet-free construction rose from `4/1` to `8/8`.
- **Regressed or unchanged:** the object remained a regular manufactured solid,
  now expressed as a smooth rounded-square foam slab. The full visual gate
  still failed before identity work, and non-egg cushion read was only `3/10`.
- **Absolute result:** this is cleaner geometry but not a viable plush head.
  It does not satisfy any user-visible likeness or medium criterion.
- **Evidence against continuing:** every rim ring is the same outline under one
  uniform scale. That homothetic construction necessarily repeats one edge
  radius around the entire object and cannot independently express crown,
  temple, cheek, underside, rear fullness, compression, or stuffing drift.
  The analytic surfaces are tangent-continuous, but their sampled cap/rim
  secants still differ by approximately `8–9°`, matching the highlight
  shoulder. More sectors, remesh, or parameter tuning cannot remove the
  structural bezel.
- **Workflow defects found by audit:** the verdict was not cryptographically
  bound to every reviewed byte; non-finite scores were not rejected; idempotent
  setup checked file existence rather than saved hashes; the review packet did
  not place controlling references beside the candidate; cleanup did not
  assert every permitted datablock; connectivity, face area, orientation, and
  signed volume were not all checked; and the untouched generated seed was
  inaccurately labeled direct-sculpt-only. Attempt 14 must correct all seven.
- **Highest-leverage unresolved problem:** derive front, rear, crown, temple,
  lower-cheek, underside, and depth behavior independently from the physical
  and turn references, then let fabric behavior and native sculpt compression
  replace the uniform mathematical bezel.
- **Approach decision:** permanently discard homothetic cap/rim generators and
  the rule that an untouched analytic primitive must already satisfy the full
  plush-medium gate. Attempt 14 will use distinct reference-derived panels and
  a nonuniform gusset, preferably pressure-settled with Blender cloth, then
  require one bounded native sculpt family to establish stuffing compression
  before the full plush gate. Its setup gate may prove only viable panel
  construction, topology, provenance, and absence of slab/egg failure.
