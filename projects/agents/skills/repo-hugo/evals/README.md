---
title: Repo Hugo skill evaluations
---

The offline target validates the Promptfoo configuration, case files, and
skill staging without credentials or model calls. Cases cover repository Hugo
and theme usage, implementation without deployment authority,
new-repository bootstrap alongside existing sites, and nested-module assembly
with the apex-site exclusion.

No live target is declared: representative behavior requires repository edits,
Bazel output inspection, Terraform plan/state observations, GitHub publication,
and DNS provider access. Those effects cannot be exercised safely in this
read-only staged fixture. Configuration validation does not establish correct
agent behavior or a successful deployment; live rollout evidence belongs to
the actual authorized task.
