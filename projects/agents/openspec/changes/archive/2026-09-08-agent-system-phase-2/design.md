# Generate a trustworthy system map and bounded context capsule design

## Context

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

The original Phase 2 work joined owner-local facts into deterministic, bounded catalogs and context views, while making interrupted durable publication and optional runtime failures explicit.

## Goals / Non-Goals

Retain the accepted objective, criteria, decisions, plans, results, and evidence with exact provenance. Importing the archive neither reenacts past operations nor reinstates deprecated workflows.

## Decisions

### Bounded catalogs and descriptor-only index

Topology, policy, action, capability, workspace-check, and historical goal catalogs are derived from owner facts with input digests, conflicts, limits, and completeness. AgentSystemIndex contains catalog descriptors and query routes rather than copied catalog bodies.

### Explicit unavailable states

The offline context command returns relevant path, label, or task context with provenance and structured partial results when optional inputs are unavailable. Runtime isolation uses package deadlines, leases, namespaces, and expected-revision publication.

### Recoverable legacy publication

The historical goal store made multi-file publication interruption recoverable and preserved legacy source identities and raw bytes during imports. That compatibility implementation is historical; OpenSpec now owns maintained work.

## Accepted plan and evidence

The accepted attempt is `goal-publication-008`, bound to criteria revision 2. Its [plan](provenance/source/attempts/goal-publication-008/plan.md), [result](provenance/source/attempts/goal-publication-008/result.md), and [criterion review](provenance/source/attempts/goal-publication-008/attempt.yaml) retain their original bytes. Every earlier attempt, criteria revision, strategy change, and evidence file is listed in [the snapshot manifest](provenance/manifest.json).

## Risks / Trade-offs

Acceptance records describe their historical candidate and evidence limits. Some commands, paths, policy proposals, and implementation claims were superseded by later work; the current baseline and owning source determine present behavior. External historical links and unavailable scratch are documented in [provenance navigation](provenance/README.md).

## Migration Plan

The completed record was placed directly in the archive. No OpenSpec archive or spec-sync operation applied these historical deltas. Native tasks map each latest criterion to its recorded accepted verdict; original records remain immutable provenance.
