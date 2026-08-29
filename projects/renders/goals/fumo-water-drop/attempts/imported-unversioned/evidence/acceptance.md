# Fumo water-drop acceptance and evidence plan

[Back to goal](README.md)

## Attempt 00 preflight gate

The preflight passes only if one inspectable packet proves all of the
following without Blender geometry:

1. The exact collection, object, rig, action, axis, unit, root, collider, and
   ownership contract is both human-readable and machine-readable.
2. The tank has at least `0.20 m` horizontal clearance around the conservative
   `0.30 x 0.22 x 0.25 m` Fumo envelope and at least `0.24 m` water depth.
3. The liquid domain encloses the tank interior, rest surface, contact pose,
   and a splash headroom band, while allowing the held asset to begin partly
   above the domain.
4. At `48 fps`, release frame `25` and first-contact frame `36` agree with a
   `0.26 m` ballistic fall to within half a frame.
5. Camera framing contains the tank base, water line, held placeholder, and
   impact pose with at least a five-percent projected margin in the analytic
   frustum estimate.
6. The placeholder is named and visibly labeled as neutral, has no Reimu
   colors, face, hair, bow, or garment imitation, and cannot be mistaken for
   accepted character work.
7. The packet explicitly records that Blender MCP was unavailable, no
   `.blend` was created, no fluid bake ran, and no tracked blend changed.

Any missing or contradictory field is a `NO-GO` for scene construction.

### Attempt 00 result

`NO-GO`. The machine report passes its `15/15` implemented checks, but the
packet fails this broader gate. The sheet places the bottom-center root at the
water line while drawing the contact envelope below it, contradicting its own
contact convention, and neither contract nor report records the protected
tracked blend hashes unchanged. Stage B is not authorized.

## Stage-B cheap Blender gate

When Blender MCP is connected, a temporary candidate must contain the exact
interface and render only frames `1`, `24`, and `36` at `640 x 360` or less.
It must pass:

- unique names and expected collection ownership;
- meters, `Z` up, front `-Y`, unit scale, and root-on-ground checks;
- no placeholder object inside `FUMO`;
- no Fumo source-library mutation or tracked-blend write;
- conservative envelope fully inside tank X/Y at contact;
- drop control at `z = 0.62 m` through frame `24` and `z = 0.36 m` at frame
  `36`;
- water line and tank base visible in all three diagnostic renders; and
- no liquid modifier, cache directory, mesh bake, foam, spray, or bubbles.

A failure returns to the contract or scaffold. It does not authorize a bake.

## Later gates, not yet attempted

1. Append the exact approved `FUMO` collection and validate its root, rig,
   neutral action, envelope, and collider without touching source bytes.
2. Run a resolution-`48` liquid preview with the diagnostic collider. Review
   every frame for tunneling, exploding liquid, domain clipping, collider lag,
   and tank leakage.
3. Re-estimate the domain and resolution from the approved asset's actual
   dimensions. A higher-cost bake is authorized only after the proxy survives.
4. Review animation, water behavior, materials, and character likeness as
   separate criteria. Passing the scaffold never inherits those passes.

## Fixed regressions

- placeholder mistaken for character work;
- scene animation inserted into the reusable asset;
- asset root or rig renamed to satisfy the scene;
- tank or domain resized after simulation merely to hide clipping;
- fluid resolution too low to resolve the conservative collider;
- camera clipping the held pose, splash headroom, tank base, or water line;
- post-impact motion presented as physically solved without evidence; and
- any final-simulation, final-animation, or plush-integration claim from proxy
  evidence.
