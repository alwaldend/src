## Context

The analyzer owns capture and conversion. The site owns the shared Mermaid
palette and its source packaging; see the analyzer change for conversion design.

## Decisions

Use companion targets for the site's own graph, since embedding that graph in
its input target would form a cycle through `genquery.scope`. Publish a separate
converter example in the analyzer's reference documentation and use the existing
landing registry and layouts. Depth eight reaches the content packages and
repository documentation owners; depth three stopped at the source archive
before those inputs were visible. Keep source files and external repositories
hidden, and label this view as a bounded overview rather than a full graph.
The site's current search-index packaging adds a stage before the source
archive, so the initial depth-seven revision was extended after base integration.
Expose separate light and dark WebP targets using the shared renderer's
`color_scheme` attribute.

Offer an explicit grouped ownership mindmap using the flowchart's existing
depth-eight JSON capture. Keep the same target selection and default external
visibility, with SVG siblings and a dark WebP. Label it as an inventory because
its branches describe package ownership rather than dependency edges. These
companion outputs also remain outside the site's own inputs.

## Risks / Trade-offs

- A bounded graph omits deeper relationships → document the depth explicitly.
- Site labels may create a wide image → inspect the actual WebP before delivery.
- Native genquery loads the full unconfigured scope before bounding the query →
  existing foreign-platform dependency failures can block graph analysis. A
  bounded command-line capture demonstrates conversion independently; it does
  not establish that the companion graph targets build.
