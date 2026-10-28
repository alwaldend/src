## Context

The confirmed subject is `pkg/mermaid/theme.json` with the shared site palette,
not the historical hand-drawn renderer. See `proposal.md` for motivation.

## Goals / Non-Goals

Use one authoritative source for each diagram and render light/dark siblings
through `mermaid_site_svg`. Preserve the active theme and renderer contract.
Experimental appearance changes are comparison inputs, not a new supported
preset or changes to the site's deployed appearance.

## Decisions

- Keep maintained examples in `docs/mermaid-theme/`, owned by the site. A
  separate tool-level gallery would misrepresent ownership of the flat theme.
- Package generated SVGs directly through `docs_filegroup`; do not check in
  generated images or copy shared theme configuration into the examples.
- Cover six diagram grammars with fictitious document-processing examples.
  Include branching, return paths, dashed edges, and relationship markers.
- Inspect native SVG output at natural and reduced widths in both modes.
  Derive comparison configurations from the current theme under ignored
  `out/mermaid-readability/`; retain exact changes and observations in the
  proposal report. Do not hand-edit SVGs to simulate a rendering change.

## Risks / Trade-offs

- Diagram families have different CSS and variables: a flowchart improvement
  is not evidence that sequence or relationship markers also improve.
- Small samples cannot prove arbitrary graph readability: include one
  moderately dense flowchart and document downscaling limits.
- Proposed variants are temporary evidence: record their input changes so
  later implementation can reproduce and reassess them.
