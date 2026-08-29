# Reimu Fumo attempt 1

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 1 — repository migration

**Candidate:**
`projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`, SHA-256
`489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`,
migrated rejected-baseline stage, review packet `goal-baseline-five-view`.

**Failure targeted:** Blender sources were stored as undifferentiated private
data, and the reusable Fumo had no owning project, durable goal record, or
stable Bazel label.

**Hypothesis:** Moving each LFS object without rewriting it, exposing explicit
targets from one project package, and recording acceptance evidence here would
create a stable base without changing any scene content.

**Plan:** Create the `projects/renders` project; classify Fumo and
miscellaneous sources under separate category directories; give each source
its own named directory; remove `data/blender`; preserve binary checksums; add
narrow Bazel filegroups; and validate the package and LFS attributes.

**Work performed:** Moved `reimu_fumo.blend` and `fumo_sisyphus.blend` beneath
`blender/fumo/`, moved `scenes.blend` beneath `blender/misc/`, added project
documentation and three explicit filegroups, excluded the non-deployable
project from the repository release aggregator, and removed the old package.

**Evidence:**

- Reimu checksum remained
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Sisyphus checksum remained
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.
- Miscellaneous scenes checksum remained
  `dee6abe31054cd882fcefa0c54e722080c092779eadcc097bae08ba6d4b867e1`.
- `git check-attr` reported the `lfs` filter, diff, and merge drivers for all
  three destination files.
- `bazel query //projects/renders:all` returned all three filegroups and both
  documentation targets.
- `bazel build //projects/renders:all` completed successfully.
- `bazel test //projects/renders:all` correctly reported that the asset-only
  package contains no test targets.
- Blender inspection found no live external libraries in the Sisyphus file;
  its appended Fumo data and packed images remain self-contained.

**Criterion results:** Repository delivery passes for the migrated baseline.
Reference likeness remains failed. Every modeling, presentation, reuse, and
animation criterion remains failed or unverified.

**Decision:** Accept the repository migration only. Continue the goal; do not
accept or refine the rejected baseline geometry.

### Progress and approach audit after attempt 1

- **Improved:** ownership, discoverability, LFS verification, and preservation
  of protected scene content.
- **Unchanged:** every visual and animation criterion, as expected for a pure
  migration.
- **Absolute result:** the tracked model is still the 3/10 rejected baseline.
- **Approach evidence:** a stable project and goal record are useful, but they
  do not support retaining the existing geometry.
- **Highest-leverage problem:** quantify the reference silhouette before
  constructing any replacement mesh.
- **Next approach:** measure the raw references, establish fixed camera and
  landmark contracts, then create an untextured panel-based blockout from a
  clean file.
