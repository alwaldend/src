---
title: Mermaid 3D
linkTitle: Mermaid 3D
description: Interactive spatial views of Mermaid diagrams
layout: landing
date: 2026-10-29
statuses:
  - active
languages:
  - ts
tags:
  - diagrams
---

Explore Mermaid diagrams as solid nodes, spatial connections, and group volumes.
Rotate the camera, trace relationships, and switch between light and dark views
using the same theme as the static diagrams. Generated HTML files contain their
fonts and renderer, so they work offline. PNG and WebP snapshots use the same
scene and can be embedded as ordinary documentation images.

## How it works

The renderer CLI obtains graph data through Mermaid's public layout extension. The adapter arranges nodes, groups, and paths in
three dimensions, using text measurements and colors from the shared theme.
The CLI writes self-contained HTML and can export the scene and metadata as JSON.

[![Mermaid 3D flowchart: mermaid_3d calls Mermaid and writes HTML and JSON](/docs/projects/mermaid_3d/docs/assets/architecture.3d.png)](/docs/projects/mermaid_3d/docs/assets/architecture.3d.png)

[Explore the flowchart in 3D](/docs/projects/mermaid_3d/docs/assets/architecture.3d.html).

The [sample gallery](/docs/projects/alwaldend.com/docs/mermaid-theme/) includes
3D flowcharts, state diagrams, classes, entity relationships, and
mindmaps. See the [project documentation](/docs/projects/mermaid_3d/) for build
commands and current limitations.
