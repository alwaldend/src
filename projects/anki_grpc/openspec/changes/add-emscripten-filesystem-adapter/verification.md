## Corrected candidate

The corrected adapter and E2E are frozen in
`out/emscripten-adapter/source-ready-2/manifest.json`. The incremental patch
starts at the original adapter blob
`417aec826d72e599b2b4faea6e660d4b9a41fba3`. Borrowed build inputs remain the
dependency owner's `15927fd` projection and are excluded from owned patches.

Before production fixes, the four new independent scenario groups all failed
against that original blob while the original eight groups passed. The
pre-fix test SHA-256 was
`ddf1a91b6fd4b11a0c718e5d8a62c7ad7dd32bba7c2f55452b6fdcc72fddfc9a`
at 2026-10-24T19:21:10Z. Retained failure evidence is under
`out/emscripten-adapter/corrections/red/test.outputs/`.

The final formatted provider E2E and Buildifier pass. The report under
`out/emscripten-adapter/corrections/green-2/test.outputs/` records eight original
and four correction groups: exclusive creation on existing/dangling symlinks
with string and numeric flags, concurrent shared options and preserved hooks,
unrepresentable write extents, and complete asynchronous callback error tuples
with caller-exception propagation.

- Corrected adapter SHA-256: `88db5dc4fe456421d5bf7e3c7344038e7680659789998ba5eacc8c85c5212554`.
- Corrected E2E SHA-256: `63e7734f3056beabd818980629b7d2f458a0e6c8ff2cc18383e53afcf7d3a8a6`.
- Generated module remains `97a0bd83047f70c80ed30ca04d61206e6fcff9c8825201131a22464e0239c41d`.
- Retained `memfs-corrections.bin` is four bytes, SHA-256
  `132d8a24110c564e939eb0af7f177d2a45fa6fefdf8dfed38a0e22cf43195c09`.
- The original 196,613-byte binary artifact retains its hash below.

The first corrected run exposed an overly specific E2E expectation: safe
integer overflow is `EFBIG`, while the engine can reject a safely represented
extent with an allocation failure, `ENOMEM`. The test now accepts either
appropriate native-allocation classification for the latter, and verifies
unchanged bytes and stream position. That intermediate evidence is retained
under `corrections/green-1/`; production error handling was unchanged afterward.
No large backing store was allocated by these cases. The adapter adds no file
quota: arithmetic bounds and actual engine allocation limits control support.

Configured target lint and strict OpenSpec validation pass. Parent-reported
19-group Chromium and six-group Node SDK acceptance on composition
`f3bfa93c699cb1d38ade5b28f0980c7b5062d8e2` used the original adapter. Those results
are historical, not corrected-candidate proof. The parent owns fresh correction
review, composition, browser acceptance and publication.

## Historical initial candidate and scope

The adapter was verified with the generated Emscripten 6.0.13 classic MEMFS
module supplied by the dependency owner. Only the adapter, its provider E2E,
minimal consumer/test BUILD additions and this change belong to this candidate.
The borrowed build projection is identified separately in
`out/emscripten-adapter/source-ready-1/manifest.json`.

The failure-oriented E2E existed before production code on
2026-10-24T19:01:26Z; its pre-implementation SHA-256 was
`e942705187d0f39b465cba7905f56cca74241f4a6e6e23bc8b6f1c2f583e98fa`.
The test and adapter were subsequently formatted by pinned Prettier.

## Historical initial evidence

`bazel_agent bazel test
//projects/anki_grpc/test/browser_backend:emscripten_filesystem_test
//:buildifier_test` passed. The retained provider report and bytes are under
`out/emscripten-adapter/provider-test-3/test.outputs/`.

The provider ran on pinned Node v24.13.0 and passed eight scenario groups:
initialization and isolation, callback/Sync binary transfers, append/truncate,
open inode identity after rename/unlink, numeric/bigint metadata, symlinks and
no-follow, directory/error behavior, and initialization failure recovery.
The test observes upstream FS through an initialization hook to verify that
reported inode identities and file bytes are the actual MEMFS values.

- Adapter SHA-256: `8e2c9022a635ef1ae3be3bf6e1b6e041ff283375369c869c597b381af25754ac`.
- Generated module SHA-256: `97a0bd83047f70c80ed30ca04d61206e6fcff9c8825201131a22464e0239c41d`.
- Retained `memfs-binary.bin`: 196,613 bytes, SHA-256
  `5d00a4a45adb20b272e76a7c821fb6eae5f3dc2e63044f7a13aabe83417a61dc`.

Strict OpenSpec validation passed. Semantic lint of the adapter and provider
E2E targets passed through the repository's lint configuration. The first test
analysis identified rules_js requiring a package-local `js_library`; packaging
was corrected before the successful runs.

## Upstream differences and acceptance boundary

MEMFS owns metadata, including its fixed ownership, block accounting and nlink
on open unlinked nodes. Birthtime uses Node's ctime fallback. Ownership changes
return `ENOTSUP`; the upstream absent hard-link primitive returns `EMLINK`.
Standard descriptors cannot be closed through the adapter, preserving the
consumer's reserved zero through two range. `fsync` only confirms in-memory
completion. IndexedDB durability stays with caller-owned snapshots.

The adapter translates explicit positioned append writes to append at the
upstream file size while retaining the stream's previous position. Recursive
mkdir delegates each creation and lookup to FS; it owns no directory data.

This evidence does not claim Go/Rust Chromium or application acceptance. Those
tests and aggregate delivery belong to the parent. No push, pull request,
global installation, backend-host change or SDK change was performed here.
