## Context

The supplied gist already contains a skill and templates.

## Goals / Non-Goals

Preserve its outside-in documentation method, verification, and triage. Do not install its shell installer or change host configuration.

## Decisions

Use the generic tools/agents owner. Adapt the entrypoint to repository authorization and documentation ownership; retain the Markdown references. Use available link validation rather than bundling a Python utility. Attribute Steve Ruiz and record the source revision in SKILL.md.

## Risks / Trade-offs

Offline eval validation checks configuration only. Full behavior requires a real source checkout and running product.
