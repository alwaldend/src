# Seat 035 successor test: pad exclusion reveals the far red skirt

Every tested ray whose first surface is the internal seat becomes a hit on
`Skirt022 joined gathered panels` when the complete seat is excluded from
the ray query. None becomes background, hem, leg or foot. The replacement
wall is typically 81–85 mm behind the seat. This falsifies the simple
expectation that concealing/removing the pad will resolve the opening:
the same opening instead reveals the distant red skirt wall.

Recommend rejecting pad-only hiding/removal as the next construction
strategy under the current fixed outer-cloth/hem/foot guards. Root owns
that decision. This test does not prove that every conceivable pad shape
would fail; it evaluates the specific visibility consequence of excluding
the existing pad, without authoring a deformation or changing coverage.

## Bounded method and frozen source

- Input: `head_032_candidate.blend`, SHA256
  `6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
- Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`, through `bazel_agent`.
- Exact front and both three-quarter cameras from the existing 032 render
  receipt; orthographic scale 0.292 m, resolution 512×512.
- Same inclusive pixel ROI x=214–297, y=438–474 in each view: 3,108
  pixel-center rays per view. The ROI is the same image coordinates,
  not a claim of identical surface sampling across views.
- Build per-object BVHs from evaluated visible/renderable geometry. For
  each ray, identify the first seat hit and nearest non-seat hit from the
  original camera-ray origin. Retain only rays where the seat is the
  overall nearest hit. Excluding the entire seat BVH avoids incorrectly
  treating the pad's own back face as its successor.
- Ray limit: 2 m. Opaque geometry only; no shader transparency,
  antialiasing, lighting change or successor render was simulated.
- No geometry, object visibility or scene settings changed; no model save,
  support-patch enumeration, shape-parameter trial or render. Source hash
  remained exact after the process, which exited normally.

The front seat-first count reproduces the prior attribution's 1,824 rays.
This is a new successor/occlusion question, not a repeated shape trial.

## Results

Depth is the extra camera-ray distance from the original seat surface to
the successor surface, in millimeters.

| View | Seat-first rays | Successor: red skirt | Other / background | Depth min / median / max |
| --- | ---: | ---: | ---: | ---: |
| Front | 1,824 | 1,824 | 0 / 0 | 51.404 / 81.247 / 86.123 |
| Three-quarter | 300 | 300 | 0 / 0 | 54.777 / 84.725 / 94.169 |
| Mirrored three-quarter | 438 | 438 | 0 / 0 | 54.613 / 84.387 / 97.347 |

Both the seat and all successor hits use `Dress red cloth.004`. These are
geometric depths, not a prediction of rendered brightness after removal.
The successor surface lies across the interior opening, not immediately
beneath the white hem. No missing-ray uncertainty affects these selected
seat-first counts: every one has a recorded non-seat successor.

Representative world points (X,Y,Z in mm):

| View / pixel | Original seat | Successor red skirt | Depth |
| --- | --- | --- | ---: |
| Front (250,454) | −3.137, −38.128, 16.793 | −3.137, +44.994, 16.793 | 83.123 |
| Three-quarter (273,459) | +31.599, −17.485, 14.174 | −29.466, +43.580, 13.587 | 86.361 |
| Mirror (214,458) | −38.804, −5.332, 14.721 | +12.667, +46.139, 14.226 | 72.792 |

The fixed front example travels from the pad's forward underside to the
rear skirt at nearly the same image height. The slight Z difference in
three-quarter rays comes from the recorded camera pitch, not a shape edit.

## Decision limit

This is sufficient to stop a pad-only removal construction trial based on
the premise that the undesirable band merely covers a harmless background.
It instead covers a deep red interior wall. Changing that outer-cloth
coverage/interface would be a materially different modeling hypothesis,
not automatic permission to move hem, feet, pad supports or the whole body.

No new acceptance threshold, reference-normal inference, whole-character
audit or visual approval is implied. The pre-existing red-wall drape,
pointed white hem and rigid collar complaints remain separate. No further
probe or deformation is proposed in this evidence packet.

## Artifacts

Compact summary and representative witnesses:
`seat_035_successor_evidence.json`, SHA256
`64db6bbbd472be5a7c06c7fb333ca5f9419bb48a7d6f5dcf0ecdd35ab55d88ff`.

All 2,562 seat-first ray witnesses, separately stored:
`seat_035_successor_witnesses.json`, SHA256
`2483f012c7ae11e7a443a938492550444f93bae64193045ea2ecc0f3437d3380`.

Read-only script: `probe_seat_035_successor.py`. All outputs remain ignored
task-local evidence; root owns any canonical evidence links or promotion.
