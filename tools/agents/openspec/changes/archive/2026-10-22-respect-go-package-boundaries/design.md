## Context

The Go skill already requires doc.go but does not explicitly require checking implementation against intended package responsibilities. Package comments describe the design; they do not enforce it. The Anki server shortcut exposed that gap.

## Goals / Non-Goals

Add a concise reusable boundary rule and offline evaluation cases. Do not encode Anki-specific architecture as a repository-wide requirement or add tooling dependencies.

## Decisions

Extend the existing package-documentation section with responsibility, effect and dependency ownership checks. Keep command presentation separate from application orchestration while preserving lower-level business owners. Extend the existing Promptfoo cases and document that offline configuration validation does not establish agent behavior.

## Risks / Trade-offs

Avoid turning inaccurate comments into unquestionable authority: compare source and callers, correct drift, and obtain authorization for material scope changes when required by repository policy.
