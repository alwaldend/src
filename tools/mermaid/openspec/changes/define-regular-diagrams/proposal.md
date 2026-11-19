## Why

Diagram requests need a consistent meaning for regular technical diagrams.

## What changes

Define a conventional flowchart style in the Mermaid skill and add an
evaluation case that checks the distinction from the sketch style.
Provide a `mermaid_diagram` Bazel macro for regular diagrams and migrate
maintained non-blog consumers to it. Remove superseded diagram aliases and
the single-use site WebP wrapper.

## Capabilities

### New Capabilities

- `regular-diagrams`: Conventional Mermaid diagram authoring guidance and
  a shared declaration for SVG and WebP renders.
