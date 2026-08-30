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
bazel_agent run //tools/versioning/cmd/versioning -- show
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
