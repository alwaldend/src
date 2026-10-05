## Context

The HTTP harness already overrides runfiles with an exact manifest entry to
verify missing-model startup behavior. The model harness consumes the real
trained artifact. See proposal.md for the verification gap.

## Goals / Non-Goals

Exercise missing and corrupt runtime inputs and real upstream training failure
without changing production code. Keep the diagram source authoritative.

## Decisions

Reuse manifest overrides for missing CSVs and a deliberately invalid pickle.
Invoke the actual build wrapper and upstream trainer with malformed CSV data,
asserting failure and absence of the output. Use no substitute model or trainer.
Add a native Mermaid person actor connected to Uvicorn by an HTTP arrow; use
the existing SVG and WebP targets instead of new render tooling. The WebP uses
the requested dark palette and 5:2 aspect ratio; documentation keeps SVG output.

## Risks / Trade-offs

Failure checks must time out and retain bounded diagnostics. Malformed training
data must fail before the expensive similarity calculation. Rendering must be
inspected for actor and label readability.
