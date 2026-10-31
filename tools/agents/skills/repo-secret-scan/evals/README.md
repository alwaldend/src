---
title: Repository secret-scan evaluations
---

The offline `eval_config_test` validates the configuration, case files, and
staged skill without model calls. Cases cover scope selection, original local
evidence, public-source triage, and incomplete scans.

A live model target is omitted because representative execution requires a
real Git history, the pinned Bazel scanner, and potentially credential-bearing
findings. The offline harness does not execute those tools or prove correct
answers. The actual repository scan supplies execution evidence separately.
