# Normalize agent semantics and declare safety contracts design

## Context

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

The original Phase 1 audit found that shared types alone did not establish a closed registered operation universe, isolated scratch, or explicit mutation ownership. The accepted work normalized agent semantics and made those safety boundaries inspectable.

## Goals / Non-Goals

Retain the accepted objective, criteria, decisions, plans, results, and evidence with exact provenance. Importing the archive neither reenacts past operations nor reinstates deprecated workflows.

## Decisions

### Owner-local declarations

Operation, lifecycle, and generated-artifact facts stay with component owners. A report-only checker derives a closed registered universe and names legacy or unavailable entries instead of becoming a hand-maintained central registry.

### Task and run isolation

Cordis scratch is scoped to task and run identities with bounded manifests. Bazel action and test temporary directories retain Bazel ownership.

### Explicit mutation and measured resource limits

The unnamed Terraform mutation alias was replaced by its explicit operation. The runner doctor reports source, pin, profile, platform, scratch, and staleness without changing host state. A criteria-bound baseline records measurements and unavailable fields.

## Accepted plan and evidence

The accepted attempt is `phase1-completion-001`, bound to criteria revision 3. Its [plan](provenance/source/attempts/phase1-completion-001/plan.md), [result](provenance/source/attempts/phase1-completion-001/result.md), and [criterion review](provenance/source/attempts/phase1-completion-001/attempt.yaml) retain their original bytes. Every earlier attempt, criteria revision, strategy change, and evidence file is listed in [the snapshot manifest](provenance/manifest.json).

## Risks / Trade-offs

Acceptance records describe their historical candidate and evidence limits. Some commands, paths, policy proposals, and implementation claims were superseded by later work; the current baseline and owning source determine present behavior. External historical links and unavailable scratch are documented in [provenance navigation](provenance/README.md).

## Migration Plan

The completed record was placed directly in the archive. No OpenSpec archive or spec-sync operation applied these historical deltas. Native tasks map each latest criterion to its recorded accepted verdict; original records remain immutable provenance.
