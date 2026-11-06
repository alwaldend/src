---
title: Repository infrastructure evaluations
---

# Repository infrastructure evaluations

This suite describes the contract for routing and executing infrastructure work
through the consolidated skill: choosing the owning stage, grouping cohesive
Terraform resources into modules, preserving the shared packaging and injection
flow, distinguishing AppRole bootstrap from KV capabilities, preparing
component credential writes with private stdin and exact authorization,
inspecting token metadata without disclosure, and validating
without mutating live systems. It also covers Ansible variable precedence and
launcher listeners, host package ownership, and Host/priority/certificate
behavior across chained Traefik proxies, DNS operational prerequisites, and
runtime discovery with single-file ownership of canonical domain names across
record types and views.

Terraform cases cover text-mode inspection without saved plans or JSON,
reviewing the plan held by an interactive apply, reporting incompatible
saved-plan-only targets without bypassing them, and reviewing remaining work
after partial failure. Normal text in agent transcripts is intentional;
durable evidence contains only sanitized summaries.
Migration cases require removal of completed temporary import and moved blocks
and a no-change plan without them, while preserving unapplied transitions when
live mutation is outside the authorized scope.

Its required offline Bazel target validates the Promptfoo configuration,
referenced cases, and staged skill without making a model call.

A live target is omitted because representative behavior requires inspecting a
real component, invoking Bazel wrappers, and reasoning about Vault, provider,
backend, and remote state. The read-only Promptfoo subject has none of that tool
or state surface. Offline validation does not show that a plan is safe, that a
playbook is safe to run, or that a record matches the deployed zone; those
claims require an isolated fixture or authorized plan review.
