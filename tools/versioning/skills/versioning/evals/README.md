---
title: Versioning skill evaluations
---

The offline target validates that Promptfoo can load the skill and cases. The
cases cover cache-stable development, SemVer-compatible week formatting,
dependency boundaries, clean-checkout stamping, honest dry-run preflights, and
release integrity failures.

A live target is omitted because realistic evaluation requires a disposable
Git repository, local ref mutations, and Bazel stamping. Offline validation
does not claim that those operations work; deterministic behavior is covered
by the Go tests under `//tools/versioning/...`.
