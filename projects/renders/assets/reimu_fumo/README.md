# Reimu Fumo asset

The target is a reference-faithful, reusable, animated 25 cm Reimu Fumo made
locally in Blender. The standalone asset must remain separate from downstream
scenes.

## Resume now

The generated [goal projection](../../goals/reimu-fumo-finish/README.md) is
the durable state authority and records the current execution and acceptance
state. In the originating worktree, `out/reimu_fumo_finish/CURRENT.md`, when
present, is an ignored convenience projection that points directly to local
models and renders. Its absence in GitHub or a fresh worktree is expected, and
it cannot override the goal record.

The
[`flatten-dose-response-018` result](../../goals/reimu-fumo-finish/attempts/flatten-dose-response-018/result.md)
records the capability stop that established blocked execution at goal
resource version 58. Confirm the goal projection before acting because later
goal state or evidence may supersede that historical result.

The
[ergonomics provenance correction](ERGONOMICS_PROVENANCE.md)
anchors the closed session review to a reachable commit and identifies the
session-local measurements that a fresh checkout cannot reproduce.

This model is unofficial fan work based on Touhou Project. Touhou Project ©
Team Shanghai Alice. The official
[Touhou fan-content guidelines](https://touhou-project.news/guidelines_en/)
are the recorded franchise-use basis.

## Source authority

- `references/canonical_front_25cm.png` controls exact-variant front
  proportions, graphics, silhouette, and the 25 cm overall scale.
- `references/canonical_turn_180.gif` controls exact-variant depth, layer
  order, and side/rear silhouettes.
- `references/physical_front.png` and `references/physical_side.png` control
  sewn construction, stuffing, compression, and contacts when they do not
  conflict with the canonical variant.
- `references/clean_front.png`, `references/sofa.gif`, and
  `references/turn.gif` are supporting evidence only.

[review_contract.json](review_contract.json) binds the exact reference and
[LANDMARKS.md](LANDMARKS.md) bytes, fixed cameras, and reviewer-role policy.
[PROCESS.md](PROCESS.md) defines the stage gates and the meaning of acceptance.
The reference files were recovered from PR #24; their acquisition and rights
status are recorded conservatively in the reference dossier and contract. At
the user's explicit direction, those public reference bytes remain in the
repository. This retention decision does not assert ownership, relicensing, or
an asset-specific sublicense.

## Authoring and verification

One explicitly named local Blender session or native artist is the sole model
writer. Every promoted checkpoint is a new file under
`out/reimu_fumo_finish/`, clean-reopened and rendered with the
repository-pinned Blender 5.2.1 toolchain. References, review cameras, lights,
and diagnostics remain outside the reusable export collection.

No external service or third-party generated asset may supply geometry,
textures, materials, rigging, or animation. This work does not modify or
integrate into a downstream scene.

No historical candidate is accepted. The exact A157 and original A202 bytes
remain only in the ignored recovery workspace because their files are not a
privacy-clean publication or resume contract. The terminal pixels for the five
rejected representation families are visible under
[failure_evidence/](failure_evidence/) and explicitly marked as historical and
non-reproducible. A privacy-sanitized A202 derivative under `donors/a202/` is
the sole tracked parts donor and resume input; starting a clean model is also
allowed. Neither route inherits a criterion pass.

The per-image creator and license fields remain `unverified` until better
provenance is found. That status is transparent metadata, not a model-quality
acceptance criterion and not a claim that the files may be relicensed.
