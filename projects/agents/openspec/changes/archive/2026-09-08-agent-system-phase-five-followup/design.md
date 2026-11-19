# Complete Phase 5 fixture and optimization follow-up design

## Context

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

The recorded Phase 5 follow-up accepted a bounded routing-coverage baseline, an inventory of deterministic writable fixtures and their gaps, a policy for isolated live comparisons, and an evidence chain for optimization adoption.

## Goals / Non-Goals

Retain the accepted objective, criteria, decisions, plans, results, and evidence with exact provenance. Importing the archive neither reenacts past operations nor reinstates deprecated workflows.

## Decisions

### Bounded coverage is an observed subset

CoverageMatrix binds normalized routing cases and the exact capability catalog digest. Emitted cases are distinguished from the complete skill universe and truncation remains explicit.

### Fixture inventory preserves missing coverage

Accepted evidence inventories Git, forge, runner, and Cordis fixtures. It explicitly records missing Terraform, Ansible, Vault/injector, and Bazel end-to-end writable fixtures; those gaps were not implemented by the acceptance review.

### Live evaluation stays separate

Live comparisons are manual or scheduled, outside wildcard tests, and bind subject model, skill, catalog, fixture, judge, and isolated authentication identities. Offline tests validate harness configuration without model calls.

### Optimization adoption follows evidence

The accepted example traces measured friction, predeclared threshold, validated proposal, owner-local change, regression, delivered revision, fallback, and retirement. The source result still names remaining work although the final criterion review accepted the documented inventory and policy; both facts are retained.

## Accepted plan and evidence

The accepted attempt is `followup-baseline-001`, bound to criteria revision 1. Its [plan](provenance/source/attempts/followup-baseline-001/plan.md), [result](provenance/source/attempts/followup-baseline-001/result.md), and [criterion review](provenance/source/attempts/followup-baseline-001/attempt.yaml) retain their original bytes. Every earlier attempt, criteria revision, strategy change, and evidence file is listed in [the snapshot manifest](provenance/manifest.json).

## Risks / Trade-offs

Acceptance records describe their historical candidate and evidence limits. Some commands, paths, policy proposals, and implementation claims were superseded by later work; the current baseline and owning source determine present behavior. External historical links and unavailable scratch are documented in [provenance navigation](provenance/README.md).

## Migration Plan

The completed record was placed directly in the archive. No OpenSpec archive or spec-sync operation applied these historical deltas. Native tasks map each latest criterion to its recorded accepted verdict; original records remain immutable provenance.

## Recorded limitations

- The accepted attempt result still names remaining work. Its final review marks all four criteria pass; evidence explicitly records Terraform, Ansible, Vault, and Bazel end-to-end fixture gaps. Import preserves both without treating gaps as implemented.
