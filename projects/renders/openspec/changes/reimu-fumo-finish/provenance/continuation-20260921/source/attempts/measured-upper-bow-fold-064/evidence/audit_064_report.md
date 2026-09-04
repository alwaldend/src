# Read-only 061-to-064 regression audit

No covered non-target surface, transform, or material-node regression found.
This is structural evidence only; it grants no visual or whole-asset acceptance.

Candidate `macro_064_candidate.blend` SHA256:
`502a52cbd825460f20f49d15d9b6443edfc68a8a9883c6d139aead7741db0211`.
Protected comparison `macro_061_candidate.blend` SHA256:
`b00c2459b1d8df746cd9dd0895ecefe3a9e92672a2798682bdaf41a0fcf0436b`.
Both hashes remained exact during the checks and in a fresh shell hash check
after both Blender processes exited. Pinned Blender 5.2.1 LTS,
build `9e2066aef7ef`, background, factory startup, automatic scripts disabled,
four threads; both audit runs exited 0. Main audit observed at
2026-09-05T17:28:53.871953+00:00. No save or render was performed.

## Results

| Check | Bounded result |
| --- | --- |
| Scene inventory | 405 objects in each file; no additions/removals after Macro061/Macro064 normalization. |
| Head and facial graphics | Mesh coordinates/topology, per-face material assignments, and transforms match exactly, including visible clipped graphics and hidden legacy eye meshes. Closed head has zero nonmanifold edges. |
| Hair | Four sheets, two locks and four hair-tie objects match exactly in covered geometry, materials, modifiers and transforms. |
| Lower-body controls and bow tails | Protected visible mesh/curve controls match; all 405 object transforms match. |
| Arms | Equivalent indexed vertex positions, edge connections and face loops; only edge/face storage order differs. See diagnosis below. |
| Sleeves and red dashes | Covered mesh data, assignments, modifiers and transforms match exactly. |
| Upper bow | Four meshes and ten curves changed in Y only; all 4,734 checked world positions retain X/Z exactly. Maximum absolute Y displacement is 38.850004 mm. |
| Materials | All 48 material property/graph fingerprints match; no added or removed material. |
| Saved simulation residue | No live CLOTH or COLLISION modifier and no object named temporary. |
| Finite geometry | All scene base mesh vertex coordinates are finite. |
| Legacy visibility | No old cuffs, eye overlays, or detached small hands render-enabled. |
| Scene basics | Both saved at frame 1 with no linked libraries. |

The first broad report flags the two arms because its fingerprint includes
edge/face array order. The focused diagnosis checks each arm's 1,106 indexed
vertex coordinates and world positions, with zero differences; all 2,256
edge connections and seam/sharp flags match after sorting; all 1,152 face
loops, preserved winding, material indices and smooth flags match after
sorting and rotation of each loop's starting vertex. No shape keys are present.
The broad raw finding is retained in `audit_064_regression.json` and resolved
by `audit_064_arm_detail.json`; neither report was overwritten to hide it.

## Hair contact coverage

The four base center surfaces contain 7,252 vertices and 13,824 triangles.
In both 061 and 064, the bounded checks found:

- Zero nonadjacent same-panel triangle overlap pairs and zero cross-panel
  overlap pairs across all six panel pairs.
- Zero degenerate triangles at the stated squared-area threshold.
- 581 vertices unanimously inside the closed head across three ray directions,
  all at or above Z135 mm; zero such vertices below Z135 mm. Maximum nearest
  surface depth is 0.8862213 mm. Six vertices have conflicting ray votes and
  two lie within 1 micrometer of the surface.
- 333 hair/head triangle contact pairs. These contacts and every reported
  per-panel metric are identical in 061 and 064.

These are inherited upper-root contacts, not a new 064 hair regression.
The measurements do not decide whether the contact or drape is visually
acceptable.

## Limits and evidence

No pixels were reviewed. Checks omit animation, rig correctness, full asset
acceptance, arbitrary mesh attributes/UVs, external image bytes, and arbitrary
custom ID properties. Material checks cover editable scalar/enum/array
properties, socket defaults, ramps/mappings, links and nested node groups;
they are not an exhaustive serialization of every Blender datablock field.

Hair overlap uses unthickened base surfaces, excludes triangles sharing a
vertex within one panel, and can count coplanar/tangent contact. Three-ray
penetration classifies vertices only, omits depths between vertices, and leaves
disagreement indeterminate. It does not test Solidify shells, torso, sleeves,
bow, locks, ties, or general scene collisions. Z135 mm is a saved-coordinate
diagnostic partition, not an independently reconstructed solver pin group.

Evidence beside this report: `audit_064_regression.py`,
`audit_064_regression.json`, `audit_064_regression.log`,
`audit_064_arm_detail.py`, `audit_064_arm_detail.json`,
`audit_064_arm_detail.log`. The 062 and 063 prepared audit scripts were never
executed because those attempts produced no saved candidate.
