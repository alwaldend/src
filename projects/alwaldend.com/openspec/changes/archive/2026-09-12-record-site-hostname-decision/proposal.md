## Why

The project has twice raised splitting the apex site's landing, documentation,
and blog sections onto separate hostnames, and the question has been answered
from conversation rather than from a durable record. The repository already
contains a deliberate earlier consolidation, so an unrecorded repeat of the
same discussion risks reversing a considered choice without weighing the
evidence that produced it. Recording the evaluation, the verdict, and the
conditions that would reopen it keeps the decision reviewable and prevents
relitigating it from memory.

## What Changes

- Record the evaluated options for the apex site's structure: one hostname with
  one build (current), separate hostnames served from separate GitHub Pages
  repositories, and one hostname with independently built and merged section
  outputs.
- Record the verified evidence bearing on the choice, including the earlier
  consolidation, the scale of the generated documentation tree, the shared
  build's failure coupling, and the existing subdomain conventions.
- Record the decision to keep the current single-hostname, single-site build,
  the reasoning, the accepted trade-offs, and the conditions that would reopen
  the decision.
- State the fallback that addresses the build-coupling concern without changing
  hostnames, so a future implementer has a supported next step.

This change adds no site behavior, build target, configuration, or DNS record.
It is a decision record only.

## Capabilities

### New Capabilities

None. This change records a decision and alters no supported behavior.

### Modified Capabilities

None. `project-alwaldend-com` and `blog-section` already describe the retained
behavior; their requirements are unchanged and deliberately not restated here.

## Impact

- `projects/alwaldend.com/openspec/changes/record-site-hostname-decision/`
  holds the proposal, design, and task record.
- No source, build, content, DNS, Terraform, or deployment artifact changes.
- No new hosting repository, certificate, or hostname is introduced.
