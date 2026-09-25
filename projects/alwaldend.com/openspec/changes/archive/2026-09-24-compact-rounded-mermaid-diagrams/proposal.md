## Why

Documentation diagrams have excessive space between ranks, while child
containers sit close to their parent titles. The user requests a more compact
diagram, more separation around child containers, and curved connectors.

## What Changes

- Reduce node and rank spacing and node padding in the site preset.
- Use Mermaid's basis curve and rounded connector caps and joins.
- Expose Dagre edge separation through a small pinned Mermaid dependency patch;
  this also increases the separation of compound-graph borders.
- Preserve Dagre placement/routing, untouched SVG output, and historical blogs.

## Capabilities

### New Capabilities

### Modified Capabilities

- `project-alwaldend-com`: compact diagrams with separated nested titles and curved connectors.

## Impact

Site rendering, browser checks, and the declared Mermaid dependency patch.
No diagram source changes or live deployment.
