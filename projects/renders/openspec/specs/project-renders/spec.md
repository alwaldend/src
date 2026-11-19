# Render assets and evidence

## Purpose

Own reusable render assets, controlling references, and evidence that binds
review claims to candidate bytes. This baseline was observed at repository
revision `550d7e79` on 2026-09-08. Sources are the
[project README](../../../README.md),
[Reimu asset documentation](../../../assets/reimu_fumo/README.md),
[review contract](../../../assets/reimu_fumo/review_contract.json),
[donor manifest](../../../assets/reimu_fumo/donors/a202/MANIFEST.md),
[failure-evidence documentation](../../../assets/reimu_fumo/failure_evidence/README.md),
and [validation targets](../../../BUILD.bazel).
The Reimu asset is unfinished; these baseline requirements preserve existing
evidence and publication constraints without claiming an accepted model.

## Requirements

### Requirement: Bind controlling references and review configuration

The Reimu review contract SHALL identify exact reference and landmark bytes
and define fixed review cameras and reviewer policy. Canonical front and turn
references SHALL control variant proportions and depth, with physical and
supporting references retaining the authority documented for each.

#### Scenario: Validate the review packet

- **WHEN** the review-evidence integrity test examines the checked-in contract
- **THEN** it checks the contract identity, reference and landmark hashes and sizes, fixed views, camera settings, and reviewer policy

### Requirement: Preserve the rejected state of historical material

The tracked A202 derivative SHALL remain labeled as a rejected technical parts
donor. Historical failure images SHALL remain labeled as rejected evidence
whose candidate bytes are unpublished and whose renders are not independently
reproducible from that packet.

#### Scenario: Reuse donor parts

- **WHEN** an author selects material from the tracked A202 donor
- **THEN** its presence in the repository confers no visual baseline, acceptance pass, or inherited criterion pass

### Requirement: Verify durable evidence bytes

The evidence integrity tests SHALL compare committed model and image bytes
with manifest hashes and byte counts. The sanitized donor's Blender audit
SHALL check the committed file for the declared privacy, dependency, rig, and
camera constraints.

#### Scenario: An evidence image changes without its manifest

- **WHEN** a committed evidence image's bytes no longer match its recorded digest or size
- **THEN** the corresponding packet integrity test fails

### Requirement: Keep unfinished candidates separate from accepted assets

Working candidates and intermediate renders SHALL remain under ignored
`out/reimu_fumo_finish/`. A reusable Reimu asset target SHALL be added only
after a candidate passes its visual, structural, animation, and exact-byte
delivery gates; historical rejected packets SHALL retain their actual state.

#### Scenario: A working render exists in scratch

- **WHEN** a candidate or render has been written under the ignored working directory
- **THEN** its existence alone does not establish durable evidence or acceptance
