## Decision

Proceed with a narrow configuration passthrough in the pinned Mermaid 11.17.2
ESM Dagre renderer consumed by the CLI. Mermaid passes node/rank separation
but omits Dagre's `edgesep`, whose default 20 also sets nested border separation.
The patch passes `flowchart.edgeSpacing` with fallback 20, including recursively
rendered subgraphs. It changes graph input before layout; it does not change
the Dagre algorithm or modify rendered SVGs.

The prototype reduces the architecture viewBox from 2120 × 790 to 1528 × 739
while increasing the parent/child top inset from 20 to 44 diagram units.
Node spacing is 32, rank spacing 40, and node padding 12. Basis curves soften
connector bends. The theme keeps its existing rounded rectangles and captions.

The strongest maintenance cost is a version-specific dependency patch. It
applies only to the non-minified ESM distribution the CLI imports; upgrades
must reconcile it. A config-only experiment changed overall spacing but left
the native nested inset at 20. A custom router was previously rejected and is
unnecessary. Remove the patch when upstream exposes equivalent spacing.

## Failure cases before implementation

- Compact placement clips labels, overlaps parent/child captions, or puts a
  child outside its parent.
- Connectors retain sharp bends or lose arrow endpoints.
- Edge separation is lost in recursively rendered subgraphs.
- Light/dark geometry differs; mobile/print or standalone fonts regress.
- Applying the dependency patch changes blog outputs when the option is absent.

## Verification

Extend the existing page E2E with measured nesting clearance, compact overall
bounds, and curved connector paths. Reuse its font, rounded-shape, mode,
mobile/print, and blog preservation checks, plus the legacy rendering and blog
freshness suites. Inspect page screenshots in both modes.

## Observed acceptance

The working-tree browser run on 2026-09-24 measured both light and dark outputs
at 1528.0133 × 739, with a 44-unit parent/child inset, 9-unit title clearance,
and 10 curved edge paths. Both page screenshots were inspected. All seven
selected browser, renderer, and blog tests passed; blog publication assets
remain unchanged. Final candidate gates and publication belong to delivery
receipts, not this implementation record.
