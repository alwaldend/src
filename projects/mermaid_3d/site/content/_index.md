---
title: Mermaid 3D
linkTitle: Mermaid 3D
description: Interactive spatial views of Mermaid diagrams
layout: landing
statuses:
  - active
languages:
  - js
tags:
  - diagrams
---

Explore Mermaid diagrams as solid nodes, spatial connections, and group volumes.
Rotate the camera, trace relationships, and switch between light and dark views
using the same theme as the static diagrams. Generated HTML files contain their
fonts and renderer, so they work offline. PNG and WebP snapshots use the same
scene and can be embedded as ordinary documentation images.

## How it works

Mermaid parses the source. Diagram adapters arrange nodes, groups, and paths in
three dimensions, using text measurements and colors from the shared theme.
The resulting HTML contains everything Three.js needs to display the scene.

[![Mermaid 3D architecture: from source and theme to an offline interactive scene](/docs/projects/mermaid_3d/docs/assets/architecture.3d.png)](/docs/projects/mermaid_3d/docs/assets/architecture.3d.png)

[Explore the architecture in 3D](/docs/projects/mermaid_3d/docs/assets/architecture.3d.html).

The [sample gallery](/docs/projects/alwaldend.com/docs/mermaid-theme/) includes
flowcharts, sequence and state diagrams, classes, entity relationships, and
mindmaps. See the [project documentation](/docs/projects/mermaid_3d/) for build
commands and current limitations.
