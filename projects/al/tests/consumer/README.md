---
title: External consumer fixture
---

The `check` command materializes these templates in an absent or empty ignored
output directory, builds the public rules and AL with independent Bazel
configuration, verifies source identity and payload bytes, size and SHA-256,
verifies order-independent configuration extensions, rejects ambiguous sibling
overrides, checks wrapper composition, flat injector groups and metadata,
target-derived call/plugin defaults, duplicate names and invalid settings,
nested plugin arguments, public file
substitution and propagated plugin/TLS runfiles, performs local TLS requests
with CA and client certificate/key runfile inputs in directory and manifest
modes from an unrelated working directory, and writes `result.json` with the verified manifest and digest. A failure exits
nonzero without a success result. The first-party module overrides reuse existing
repository pins. No parent Bazel configuration is imported.

Build a launcher from the repository root, then invoke it outside the Bazel
client so the independent build does not contend for the caller's server lock:

```sh
bazel run --script_path=out/task/consumer-check @rules_al//tests/consumer:check
out/task/consumer-check \
  --source_root "$PWD" \
  --output "$PWD/out/task/consumer"
```

The check invokes `bazel` from `PATH` in the generated workspace and copies the
repository's pinned `.bazeliskrc`. Supplied paths must be absolute. It runs as
operator-side validation; it is not a build action and may use Bazel's normal
dependency fetching and user configuration. Retain `result.json` and the generated
workspace as repeatable evidence. Build and test the AL module separately for its
internal behavior.
