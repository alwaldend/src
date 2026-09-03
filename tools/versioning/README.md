---
title: Versioning
description: Global repository and project versioning
---

`versioning` owns versions for this repository and its first-party projects.
It deliberately does not manage versions of third-party dependencies.

The repository uses SemVer-compatible calendar versions:

- ordinary development: `0.0.0-dev`;
- nightly trunk tag: `vYYYY.W.0-nightly.YYYYMMDD`;
- weekly release branch: `releases/YYYY.W`;
- release tag: `vYYYY.W.PATCH`.

`YYYY` and `W` are the ISO week-year and week. The week is not zero-padded,
because SemVer forbids leading zeroes in numeric identifiers. Patch zero is
the release branch point. Every first-parent commit after that point advances
the calculated patch number by one.

Calculated versions and Bazel status omit the Git tag's leading `v`.

A commit may carry one nightly tag and one release tag when a nightly is
promoted. Branch context selects the channel automatically. On detached HEAD,
pass `--channel release` or `--channel nightly` for that exact co-tagged
commit. For an untagged detached commit from a release branch, pass
`--release YYYY.W` so the tool can calculate its patch from the correct
branch-point tag.

Build and inspect the tool with:

```sh
bazel_agent bazel run //tools/versioning/cmd/versioning -- show
```

For a stamped build, use the bootstrap entry point. It generates a
source-current Bazel launcher under `out/versioning/`, then the Go tool runs
Bazel with itself as workspace status:

```sh
tools/versioning/cmd/versioning/versioning.sh bazel -- \
  build --config=release //path/to:artifact
```

Read `$versioning` for the guarded nightly, release, and Bazel stamping
workflows.

## Guarded release-ref publication

`versioning` also owns a typed, reviewed release-ref plan and a
provider-neutral guarded publisher for the generated nightly and release
tags. It never merges versioning, delivery, or goal authority; the tool
retains its own `release-refs` authority.

Generate a deterministic plan from the resolved version state:

```sh
bazel_agent bazel run //tools/versioning/cmd/versioning -- release-plan \
  --plan out/delivery/release-plan.json
```

The plan records the exact version, channel, commit, tree state, target refs
(tag-only for nightly, branch-plus-tag for release), the expected remote
preconditions, and whether atomic multi-ref publication is required. It is
reviewed before consumption and is not an authorization.

Publish the reviewed plan only after explicit release scope:

```sh
bazel_agent bazel run //tools/versioning/cmd/versioning -- release-publish \
  --plan out/delivery/release-plan.json \
  --receipt out/delivery/release-ref-receipt.json
```

The guarded publisher fetches expected remote state, acquires a distinct
`release-refs` lease, publishes the refs (atomically when required and
supported), and verifies the remote before emitting a `ReleaseRefReceipt`.
An existing immutable release tag never moves; a remote that cannot guarantee
atomic multi-ref publication is an explicit refusal, never a generic success.
