---
title: External consumer fixture
---

The `check` command materializes these templates in an absent or empty ignored
output directory, builds the public rules and AL with independent Bazel
configuration, verifies source identity and payload bytes, size and SHA-256,
and writes `result.json` with the verified manifest and digest. A failure exits
nonzero without a success result. The first-party module overrides reuse existing
repository pins. No parent Bazel configuration is imported.

Build a launcher from the repository root, then invoke it outside the Bazel
client so the independent build does not contend for the caller's server lock:

```sh
al bazel run --script_path=out/task/consumer-check @rules_al//tests/consumer:check
out/task/consumer-check \
  --source_root "$PWD" \
  --output "$PWD/out/task/consumer"
```

The check uses its source-built AL runfile. Optional `--al` overrides it with a
native executable; a Bazel run launcher changes its working directory to runfiles.
Supplied paths must be absolute. The check uses
an independent `.bazelrc` containing only its `agent` profile and copies the
repository's pinned `.bazeliskrc`. It runs as operator-side validation; it is not
a build action and may use Bazel's normal dependency fetching. Retain `result.json`
and the generated workspace as repeatable evidence. Build and test the AL module
separately for its internal behavior.
