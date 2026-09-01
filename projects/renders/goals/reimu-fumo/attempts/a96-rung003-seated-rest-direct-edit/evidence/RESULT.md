# A96 author result — seated-rest lower-stack coupon

## Verdict

**STOP / RESET. Do not promote or correct this candidate.**

The candidate is structurally admissible and the front view gives both foot
pods useful additional vertical fullness. It nevertheless fails the first
absolute pixel gate categorically. In exact side, the red lower assembly still
reads as a long rearward cape or ramp that terminates in a thin floor rail,
with only one forward ball-like foot visible. The worst three-quarter view
repeats that long trailing silhouette and rail. Front still has a nearly
horizontal, dark hem/underside band rather than a soft U-shaped skirt that
partly occludes two stuffed feet.

This is not a one-local-contact-defect case, so the plan's one-correction
allowance does not apply. No further geometry edit was made after this pixel
judgment, and no implementation-blind preference review was requested.

## Bound input and outputs

- Exact parent:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/`
  `reimu_fumo_working_rung_003.blend`
- Parent SHA-256 before and after authoring:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected tracked canonical Blend SHA-256 after authoring:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`
- Candidate:
  `work/reimu_fumo_a96_seated_rest.blend`
- Candidate SHA-256:
  `73099d33e19e7cc73be6b1184df3b3a41a437bcad02d481ae02bf8988f6699a7`
- Author script: `work/build_candidate.py`, SHA-256
  `af43de56225766e912fd8e160cc6127410557b3b17dd08b9128cebae76169b1f`
- Edit manifest: `work/EDIT_MANIFEST.json`, SHA-256
  `877bdc83e5c8541fff2d7d0e60c3614c56c1e385a1eabfa826dfe4ab49d314f0`

All paths in the sections below are relative to this `author/` directory
unless shown otherwise.

## Edit performed

One coupled direct-edit coupon was made on the exact rung003 parent:

- reversible `A96 Seated Rest` shape-key deltas were added to the changed
  lower-stack objects;
- the front and rear skirt panels, side gussets, ruffles, and all eighteen hem
  stitches were moved as coupled surfaces;
- existing foot-pod width and depth were preserved while their vertical
  fullness was increased by a factor of 1.40 from the ground plane; and
- the two existing hidden leg roots were lengthened slightly.

Thirty of the thirty-one allowlisted lower-stack objects changed; the compact
internal seat pad remained unchanged. No object, topology, modifier, material,
camera, light, world, scene setting, rig, or visible owner was added, removed,
or reassigned.

## Structural verification

The candidate was reopened cleanly with repository-pinned Blender 5.2.1 LTS,
automatic embedded-script execution disabled, and checked against the
independent verifier baseline.

- Candidate inventory: `work/CANDIDATE_INVENTORY.json`, SHA-256
  `1330936cde85926b40436a1ddd35d3c2696e92395bd27cada23689bf4e137bda`
- `baseline_comparison.all_pass`: `true`
- Object-name digest: pass
- 146 protected-object digest: pass
- 31 editable-object structural digest: pass
- Global material-contract digest: pass
- Scene/world/render-contract digest: pass
- Candidate Blender version: `5.2.1 LTS`

This proves boundary and structure preservation only; it does not qualify the
failed pixels.

## Fixed-view artifacts and pixel judgment

| View | Artifact | SHA-256 | Judgment |
|---|---|---|---|
| Front | `renders/packet/front.png` | `530623567a886536e4d5b30a673729713afd9182d0db4338aa42fc2efa593b7f` | Foot pods are taller and separately readable, but the hem remains a dark near-horizontal rail and the skirt is not a compact soft U-shaped pool. |
| Exact side | `renders/packet/side.png` | `a4e7dd72bfb389cfacdb8a0d457a2dce45ed056365d001e06d6560c282eb1441` | Categorical fail: a single forward ball-like foot plus a long rearward cape/ramp ending as a floor rail. |
| Worst 3Q | `renders/packet/three_quarter.png` | `073cfc22f0061a3cdea3725ff58728d297101beb96c32a817c48d4e4ac93a5df` | Categorical fail: the lower panel trails far behind the seat and lies along the floor; it does not read as compact pooled plush fabric. |

All three artifacts are 512 by 512, opaque, non-black, and non-blank. Their
render manifest is `renders/packet/manifest.json`, SHA-256
`e6436e644f8cb81135dd87c547c6e20b0431dd8abb12ac4d46e66986a373fa08`.
The combined three-view artifact is
`renders/candidate_contact_sheet.png`, SHA-256
`61e1339dbb9c01806792f2a5a397a568455b3dcb0118fc71801d1af1767436d7`.

## Reference evidence used

- `projects/renders/blender/fumo/reimu_fumo/references/`
  `canonical_front_25cm.png`
- `projects/renders/blender/fumo/reimu_fumo/references/physical_front.png`
- `projects/renders/blender/fumo/reimu_fumo/references/physical_side.png`
- `out/reference_metrics/turn_contact.png`
- `out/reference_metrics/sofa_contact.png`
- `out/reference_metrics/turn_side_grid.png`
- `out/reference_metrics/turn_3q_grid.png`
- exact rung003 baseline fixed views under
  `out/reimu_fumo_attempt_096_rung003_seated_rest/verification/packet/`
- rung003 five-view contact sheet under
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/five_view_gate/`

The controlling references show a short compact seat, red and white fabric
dropping around and partly onto the feet, and stuffed foot pods rising into
the hem. They do not support the candidate's rearward sheet-like extension.

## Honest result and reset signal

A96 demonstrates that sparse reversible coordinate deltas can safely create
large lower-stack changes while preserving the proven upper model. It also
demonstrates that this particular panel-row remap is the wrong geometry
strategy: it converts a rigid trapezoid into a longer trailing surface rather
than producing a compact seated volume. The candidate must remain isolated
evidence only. Any future lower-stack attempt should first prove, in exact
side, a short self-supporting pooled silhouette with no trailing rail before
carrying the full panel/ruffle/stitch system or spending a second correction.
