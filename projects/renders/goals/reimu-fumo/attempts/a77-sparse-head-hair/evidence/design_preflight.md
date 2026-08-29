# A77 sparse-form design preflight

## Verdict

**REVISE before building P0.** The head-rest and rear-leaf contracts are
implementable, but the hair-mantle face selection contradicts its required
connected topology. Building it as written would spend the first candidate on
a known floating two-face island.

## Blocking contradiction and exact correction

The written mantle selection has the advertised `89 V / 165 E / 76 Q`, but it
has two face-connected components:

- 74-face crown/side/rear shell; and
- isolated two-face central-front island in Z band 1.

The two outer faces in front Z band 2 do not touch that central island. A
Subsurf then Solidify stack would preserve two disconnected fabric pieces,
contradicting “connected open result” and risking an obvious floating tongue.

Minimal purpose-preserving correction: include the two missing central faces
in front Z band 2 as well. Front band 2 then contains all four faces, and the
two central band-1 faces form the intended lower central tongue. Update the
mantle contract and any object suffix/assertions to:

```text
89 V / 166 E / 78 quads
one face-connected component
20 boundary edges
```

This correction preserves the proposed crown, side, rear, and low off-center
front boundary. Do not delete arbitrary faces merely to retain the old `76`
face count; connectivity and the visible opening are the design requirements.

## Verified contracts

| Contract | Preflight result |
| --- | --- |
| Head boundary grid | Exact `114 V / 224 E / 112 Q`; one closed component; zero boundary edges |
| Mantle as written | Exact `89 V / 165 E / 76 Q`; **two components (`74 + 2`)** |
| Corrected mantle | Exact `89 V / 166 E / 78 Q`; one component |
| Rear-leaf grid | Exact `20 V / 31 E / 12 Q`; one open component |
| Seven source hide names | All exist exactly; five meshes and two curves |
| Blender modifier API | Blender 5.2.1 accepts Catmull-Clark `levels=1`, `render_levels=2`, followed by Solidify with outward `offset=1.0` and 3.3--4.0 mm thickness |
| Protected rung | Still `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b` after both read-only probes |

The reproducible in-memory check and machine result are
`preflight_api.py` and `preflight_api_probe.json`. No `.blend` was saved or
mutated.

## Interface implementation guards

The face contract is numerically plausible, but Catmull-Clark does not
interpolate its control vertices. Fit and measure the **evaluated** surface at
the nine anchors, then separately sample every eye and mouth witness. The old
eye clearances (`0.076/0.080 mm`) and mouth clearance (`0.320 mm`) already sit
inside the proposed ranges, so there is no preflight contradiction if the
local evaluated fit succeeds.

The `0.6--1.0 mm` mantle rest offset is comparable to the old cap/receiver
clearance (`0.54--1.02 mm`), but it cannot be assumed to satisfy the stricter
`0.25 mm` old-cap root match after the macro head changes. Solve that offset
against the old cap only under retained front/temple roots and run the stated
early interface veto before deriving the leaf.

“Snap” the rear-leaf root must not mean zero clearance: project each root
control to the evaluated mantle, then displace it outward along the mantle
normal by a declared target such as `0.4 mm`. Verify the final evaluated
Subsurf/Solidify surfaces against the `0.2--0.8 mm` gate. Otherwise the snap
instruction and minimum-clearance gate conflict operationally.

For the clay discriminator, assign the existing face clay to the head-rest
and existing hair clay to mantle/leaf. Hiding the old receiver without this
explicit assignment would erase the face carrier even though the topology
probe passes.
