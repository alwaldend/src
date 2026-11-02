## 1. Conversion and build rules

- [x] 1.1 Define behavioral failure cases before implementation and add an E2E
      command fixture that preserves Mermaid and validation results as artifacts.
- [x] 1.2 Implement the JSON converter and CLI; verify shared edges, generated
      files, filters, cycles, escaping, invalid input, and deterministic ordering.
- [x] 1.3 Add target and JSON-input rules; verify a real fixture graph builds
      through genquery, conversion, and WebP rendering.

## 2. Documentation and acceptance

- [x] 2.1 Document rule and command usage, defaults, graph semantics, and site
      companion targets; build the documentation assets and inspect their SVG.
- [x] 2.2 Verify site integration through its owner-local change, inspect the
      candidate site graph WebP, and record exact validation evidence.
- [x] 2.3 Validate formatting and affected package checks, then hand the
      candidate to al-delivery; its receipts own publication verification.

## 3. Optional external boundary view

- [x] 3.1 Add E2E failure cases for an opt-in external-leaf boundary, shared
      children, external roots, depth, and unchanged existing modes.
- [x] 3.2 Add the command and rule option without changing defaults; document
      its graph traversal and edge semantics.
- [x] 3.3 Validate and show the site's dark diagram locally. Publication is
      withheld by the user's instruction.

## 4. Optional package grouping

- [x] 4.1 Define E2E cases for full package captions, short node captions,
      nested external repositories, shared target names, escaping, stable
      ordering, and preservation of the selected dependency edges.
- [x] 4.2 Add opt-in package grouping to the command and rules, retaining the
      existing flat output by default.
- [x] 4.3 Render and inspect the grouped dark graph, validate locally, and
      retain source changes without publication.

## 5. Compact grouped layout

- [x] 5.1 Compare native layouts at the same font size and with identical
      nodes, edges, and package groups; measure canvas dimensions and inspect
      caption clearance.
- [x] 5.2 Apply the selected grouped layout, verify the rule and actual site
      preview, and show a local comparison without publication.
- [x] 5.3 Nest all visible external targets under `External` in flat and grouped
      modes, with behavioral E2E coverage defined before implementation.
- [x] 5.4 Render mindmap, treemap, and block alternatives; show the attempts and
      document their relationship and readability limitations.

## 6. Selectable mindmap overview

- [x] 6.1 Define command and native-render failure cases for an opt-in ownership
      mindmap, unchanged defaults, unique targets, filters, grouping, and escaping.
- [x] 6.2 Add `diagram_type` / `--diagram-type` with flowchart as the default and
      mindmap as an ownership overview; document the relationship distinction.
- [x] 6.3 Add site companion targets, validate both native pipelines, and show
      the dark outputs locally. Canvas size is secondary to readable structure.
      The user subsequently authorized publication; hand the completed revision
      to al-delivery for the existing pull request.
