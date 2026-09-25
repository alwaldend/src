---
title: Mermaid dependency patch
description: Dagre edge-spacing configuration for the pinned Mermaid renderer
---

The JavaScript manifest and lock in `tools/` pin the Mermaid dependency graph.
`tools/js/include.MODULE.bazel` applies `edge-spacing.patch` to Mermaid 11.17.2
through `npm_translate_lock`. No version or downloaded package is duplicated here.

The patch forwards `flowchart.edgeSpacing` to Dagre's `edgesep` option, including
recursively rendered subgraphs. Dagre uses that separation for edge placeholders
and nested container borders. Unconfigured consumers retain the default 20;
historical blog rendering must remain byte-identical.

It patches the non-minified ESM Dagre module loaded by Mermaid CLI's
`mermaid.esm.mjs` entry point. Other distribution formats are not extended.
Reconcile the patch and rerun the site browser and blog freshness tests on a
Mermaid upgrade. Remove it when upstream exposes this spacing option.

Upstream source: [Mermaid Dagre renderer](https://github.com/mermaid-js/mermaid/blob/mermaid%4011.17.2/packages/mermaid/src/rendering-util/layout-algorithms/dagre/index.js).
The patch changes configuration before layout, without changing Dagre routing
or processing generated SVGs.
