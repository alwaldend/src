---
title: Repository delivery evaluations
---

# Repository delivery evaluations

This suite describes the behavioral contract for safely finalizing a feature
branch and its pull request. Its required offline Bazel target validates the
Promptfoo configuration, referenced cases, and staged skill without making a
model call.

The skill entry point contains ordinary delivery and final review completion.
Conditional packaged references retain history rewrites and recovery, review
mutations, and Forgejo. Cases exercise those routes as well as older-runner
fallback and explicit validation-gap handling; duplicated review cases are not
counted as extra coverage.

A live target is omitted because representative delivery requires inspecting
and changing Git history, pushing a remote ref, interacting with a forge, and
responding to live review state. A read-only, tool-free model evaluation cannot
reproduce those operations safely or verify their postconditions. Promptfoo
validation therefore proves only that the evaluation assets load; delivery
behavior still needs isolated integration fixtures or review of a real,
authorized delivery.

The cases preserve failure-sensitive invariants for future live coverage. On a
supported GitHub adapter, publish readiness requires a validation decision
bound to the latest `prepare` and its exact returned HEAD OID. Required checks
run after preparation by default. A tree-preserving amendment may reuse prior
candidate-bound passes only with recorded unchanged relevant inputs; commit-,
history-, stamp-sensitive and uncertain checks rerun. The `deliver` wrapper
pauses for this decision before publication. Structured `validate` records the
caller-selected plan and exact passing candidate; explicit `continue --publish`
uses that evidence. Manual `publish` still requires the literal validated OID.
The structured path rejects stale inputs and never automatically repeats an
uncertain publication. The launcher case avoids nesting validation beneath an
ordinary lock-holding Bazel invocation on an older installed runner.
Review reads and mutations stay inside
`repo_delivery review`.
Every behavior-changing review fix also invalidates the prior correctness
verdict and requires a fresh, proportional, diff-focused scrutiny pass; green
tests alone do not replace that reasoning gate.
An explicitly authorized multi-commit range uses `--consolidate` with the
literal inspected local head only after ownership review; its linearity,
identity, oldest ownership marker, pull-request projection matching the
requested aggregate message, signature requirements, and remote lease remain
fail-closed.
The literal candidate and strict preparation receipt flow into publication;
an advancing base produces a new exact candidate and derived receipt that must
be validated directly. A divergent remote replacement remains refused unless
`$git-rebase-remote` first preserves the exact old tip and establishes
task-owned rewrite authority; only that literal freshly inspected OID may be
passed to `--replace-remote` and bound into the final snapshot and receipt.
When an isolated rebase stops on conflicts, a manual rebase is permitted only
for exact, evidence-backed minor hunks whose unambiguous resolution preserves
both sides; all other conflicts still require user direction.
Review cases also preserve pull-request/head/context,
last-comment, snapshot-digest, and reply-receipt guards, exact comment
disclaimers, post-mutation reinspection, and waiting for a started remote
review of the exact final head to reach a terminal state before delivery is
reported complete. When neither the product monitor nor the selected adapter
exposes authoritative review execution state, the suite requires a bounded
observation attempt and an explicit unverifiable result instead of inferred
success or an indefinite wait. Validation runs from a clean checkout of the
exact candidate whenever unrelated dirt could influence a check, unless the
caller can prove that the check cannot consume that dirt.
The compatibility cases separately require honest Forgejo detection, a
completely clean worktree and index without autostash, exact commit and
pull-request disclaimer strings, and only the capabilities exposed by
`$git-rebase-remote` and the installed `fj` v0.5.0 wrapper. A fork case records
that the GitHub adapter cannot reliably discover or map upstream fork pull
requests, so same-repository topology remains a caller-enforced precondition.
The transport case preserves the separate version 1 boundary: Git network
operations require SSH and must not regain HTTPS support by loading mutable
global credential-helper or repository transport configuration.
