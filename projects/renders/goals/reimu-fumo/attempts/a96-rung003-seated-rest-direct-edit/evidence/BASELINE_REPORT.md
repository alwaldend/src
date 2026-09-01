# A96 exact-parent baseline and ownership verification

## Verdict

The A96 baseline is valid and reviewable, and both protected Blend inputs
remain byte-identical. Repository-pinned Blender 5.2.1 LTS rendered three
fresh 512 px views directly from the exact rung003 parent. All three images
are opaque, non-black, non-blank, and visibly contain the complete subject.

The baseline itself remains an absolute visual reject for A96's lower-body
goal: the skirt is a rigid cone/trapezoid rather than a pooled seated fabric
form; the exact side view has a tall narrow support and only one long,
forward-projecting foot silhouette; the three-quarter view exposes the rigid
hem/ruffle stack as a dark floor rail. These are baseline defects, not render
failures.

## Bound source and tool

- Exact parent:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/`
  `reimu_fumo_working_rung_003.blend`
- Parent SHA-256 before and after all operations:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected tracked Blend:
  `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`
- Protected tracked SHA-256 after all operations:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`
- Tool: repository-pinned Blender `5.2.1 LTS`, background mode, automatic
  embedded-script execution disabled.
- Loaded scene: `Attempt41_Manual_Head_Maquette`.
- Read-only inventory ended with `bpy.data.is_dirty == false` and an unchanged
  parent hash. No Blender save operator was called.

## Fresh baseline views

| View | Camera | Artifact | SHA-256 | Pixel sanity |
|---|---|---|---|---|
| Front | `Review_front_Camera` | `packet/front.png` | `12323e56e611ccdb0bf826f443a2c6f89018dd263ee38d028c6de4666bdfe892` | 512×512, opaque, min 0, max 65535, mean 47106.8, 17,941 colors |
| Exact side | `Review_side_Camera` | `packet/side.png` | `22f360a13e013b66a43051e3ab697dcfaa53be0de904ebcdaf9cfaab2d548c60` | 512×512, opaque, min 0, max 65535, mean 50895.1, 13,975 colors |
| Worst 3Q | `Review_three_quarter_Camera` | `packet/three_quarter.png` | `f19a35e6ce632f3ed01545c52d5aa5d6ba7aa3b2ffc34bb91c7868b499adecf4` | 512×512, opaque, min 0, max 65535, mean 47835.1, 18,318 colors |

Pixel statistics use ImageMagick's 16-bit quantum range. I also inspected the
actual pixels of every image and the combined contact sheet; statistics alone
were not used as a visual-validity claim.

- Render manifest: `packet/manifest.json`, SHA-256
  `c0b8d18f71c865383642a07ed20eac0373adf59a2e6c1b739ab0cb5a71010812`
- Three-view contact sheet: `baseline_contact_sheet.png`, SHA-256
  `9bc294ce5541d22d60bef732d1fdcaba8043450cf9a2e63ee58f0a7e98222324`
- Render specification: `render_packet.json`, SHA-256
  `760be90ca6356090f7077909babef50cd746e1ae06e9967ea09ef77ea320b736`

`Review_three_quarter_Camera` is the controlling worst-three-quarter baseline
because it exposes both feet, the underside/hem depth order, and the long
floor-rail silhouette together. The mirror camera is less discriminating for
this lower-stack coupon.

## Reusable ownership manifest

The task-local checker is `check_inventory.py`, SHA-256
`ba909b2b7b049a0b810ff13111e7d0fb895fb1096e7c45ec1d0d3b1044097ec9`.
Its complete baseline output is `baseline_inventory.json`, SHA-256
`9994eece6d6a762b8d47106390c1784e42286cfb5111b55571891360b6336976`.

It found the expected 177 objects and partitions them into exactly 31 editable
lower-stack objects and 146 protected objects. It fingerprints:

- the complete object-name set;
- every protected object's identity, collections, transform, base geometry
  coordinates and topology, materials, modifier parameters/order, visibility,
  parenting, constraints, and vertex-group names;
- editable object identities, collections, transforms, topology, modifier
  parameters/order, and material assignments, while intentionally excluding
  editable coordinates and shape keys;
- all materials and node/link contracts; and
- scenes, world, render settings, camera selection, unit scale, and frame
  state.

Baseline invariant digests are:

| Contract | SHA-256 |
|---|---|
| Complete object-name set | `d8ab0f8e4f7ef3c77a978dbed9007d1d5926c4f004465bfb1e2f687df424fde0` |
| 146 protected object records | `2259233e62922c06d5c1aa7fb1bced4eca939b58875f29f7fc76a000e96221c7` |
| 31 editable structural records | `b0c33a5605f3496fa010aee9e3fa8fc523fcc7be5cf2473c16f7dbc84e3df2f3` |
| All material contracts | `37f6d23e824c73fe317f099c575b54371fb3494fd8a18392ed7213c57ea507ba` |
| Scene/world/render contract | `c437b9777565fe7cfe02b5a1c57682acda532c1340376a2642ada9b31af326a4` |

The checker was reopened against the same parent with
`--baseline-manifest baseline_inventory.json`. All five digest gates passed in
`baseline_self_check.json`, SHA-256
`03b8e601cd9bb42720b1a3107c7f9357fab91e43be88cb827ba3e318387362b4`.
This proves the task-local comparison path is stable across clean opens of the
baseline on the pinned Blender.

### Exact editable allowlist

The only editable object identities are:

```text
Garment42 compact internal seat pad
Garment42 front gathered ruffle
Garment42 front shallow lap panel
Garment42 hem stitch front 00 through 08
Garment42 hem stitch rear 00 through 08
Garment42 left side gathered ruffle
Garment42 rear pooled dress panel
Garment42 rear pooled ruffle
Garment42 right side gathered ruffle
Garment42 side gusset left
Garment42 side gusset right
Left black stuffed foot pod
Left short hidden leg root
Right black stuffed foot pod
Right short hidden leg root
```

The two stitch ranges are inclusive and expand to 18 distinct curve objects.
The 146-object protected complement includes all head, hair, face, eye, lock,
bow, bodice, collars, cravat, current sleeves, hidden superseded geometry,
review floor, cameras, and lights. Materials, world, and render settings are
protected by their separate global digests.

## Candidate check invocation

Run the checker through the repository-pinned Blender with the candidate path
and a new task-local output path:

```sh
bazel_agent run //tools/blender:blender -- \
  --background --factory-startup --disable-autoexec \
  --python-exit-code 2 \
  --python out/reimu_fumo_attempt_096_rung003_seated_rest/verification/check_inventory.py \
  -- \
  --blend-file <candidate.blend> \
  --output <candidate-inventory.json> \
  --baseline-manifest out/reimu_fumo_attempt_096_rung003_seated_rest/verification/baseline_inventory.json
```

An admissible candidate must report `baseline_comparison.all_pass == true`.
That is a structural veto only. It does not establish pinning of attachment
rows, evaluated-surface clearance, correct foot occlusion, or visual likeness;
those remain fixed-view pixel gates.

## Concrete boundary concerns for the author and coordinator

1. `Dress red cloth` is shared with the protected bodice, and
   `Dress warm white cloth` is shared with protected collar/legacy geometry.
   Any material mutation is out of scope and will fail the global material
   digest even if performed through an editable object.
2. The editable boundary contains thin curve stitches plus mesh ruffles and
   panels. All must move as coupled surfaces; leaving a stitch or ruffle behind
   can pass topology checks but will visibly detach or cross in the render.
3. The baseline has no existing shape keys or hierarchy. Candidate shape keys
   are permitted only on the 31 editable objects. The checker intentionally
   allows their coordinate deltas but rejects topology, object-transform,
   modifier, material, scene, or protected-coordinate drift.
4. Exact front alone is insufficient: it hides the baseline's one-foot side
   silhouette and the excessive narrow/tall depth profile. Exact side and the
   selected worst 3Q must remain mandatory before any keep decision.
5. The current lower stack already creates a dark horizontal floor rail and
   rigid cone/tent silhouette. A candidate that merely moves vertices while
   preserving those pixels is a categorical reset, regardless of clean
   fingerprints.

## Result

Baseline provenance, pixels, and the ownership-check interface pass. No model
criterion is passed by this verification; it supplies the immutable comparison
and regression vetoes for the author's A96 candidate.
