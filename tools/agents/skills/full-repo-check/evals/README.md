---
title: Full repository check evaluations
---

# Full repository check evaluations

This suite describes the required orchestration and reporting behavior of the
full-repository check. Its offline target validates the Promptfoo configuration,
referenced case, and skill staging without making a model call.

A live target is omitted because representative behavior requires invoking
Bazel across every workspace and inspecting generated diagnostic logs. The Go
runner tests cover all eighteen child commands, continuation after failures,
restricted artifact permissions, and report generation deterministically.
Promptfoo validation alone is not evidence that an agent completes the audit
correctly.
