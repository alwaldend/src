---
title: Mermaid theme gallery
description: Six diagram types rendered with the current flat site theme
---

These fictitious document-service examples share the site's current
[flat theme](https://github.com/alwaldend/src/blob/master/data/alwaldend/assets/mermaid/config.json.tmpl) and
[palette](https://github.com/alwaldend/src/blob/master/data/alwaldend/assets/theme.json). The site selects light or dark
images with its theme menu. Open an image to inspect its native size.
The examples carry structure and labels only; appearance belongs to the theme.
Labels use embedded Liberation Mono, including when an image is opened alone.
Normal graph paths use 2.25px strokes, captions have muted frames, sequence
guides recede behind messages, and mindmap branches keep visible contrast in
both modes. Dense diagrams still benefit from opening the full-size image.

Each example also has an interactive 3D view with solid nodes, spatial paths,
and volume outlines for groups. These views use the same palette and embedded
font. Drag to orbit, right-drag to pan, scroll to zoom, or select an element to
trace its connections. Each HTML file works offline.

## Flowchart

Branching decisions, a container, two data stores, and dashed revision/retry
paths test arrow direction, caption weight, and route separation.
[Source](flowchart.mmd) · [Interactive 3D](assets/flowchart.3d.html) · [3D snapshot](assets/flowchart.3d.png) · [Dark snapshot](assets/flowchart.3d.dark.png).

[![Document processing with validation, retry, search, and audit paths](assets/flowchart.svg)](assets/flowchart.svg)

## Sequence

Calls, dashed replies, an activation, a self-call, an alternative, a loop,
and a note test the distinction between messages and structural guides.
[Source](sequence.mmd) · [Interactive 3D](assets/sequence.3d.html) · [3D snapshot](assets/sequence.3d.png) · [Dark snapshot](assets/sequence.3d.dark.png).

[![Document submission followed by repeated queue delivery and validation](assets/sequence.svg)](assets/sequence.svg)

## State

A composite processing state, retry cycle, cancellation, and terminal markers
test transition labels and nested boundaries. [Source](state.mmd) · [Interactive 3D](assets/state.3d.html) · [3D snapshot](assets/state.3d.png) · [Dark snapshot](assets/state.3d.dark.png).

[![A queued document progresses through validation and indexing or retries after failure](assets/state.svg)](assets/state.svg)

## Class

An interface, implementation, dependency, association, and composition test
open and filled markers, dashed lines, and cardinality labels.
[Source](class.mmd) · [Interactive 3D](assets/class.3d.html) · [3D snapshot](assets/class.3d.png) · [Dark snapshot](assets/class.3d.dark.png).

[![Repository interface, search implementation, indexer, documents, and revisions](assets/class.svg)](assets/class.svg)

## Entity relationships

Identifying and non-identifying relationships, optional cardinality, and
attribute tables test crow's-foot markers and fine table separators.
[Source](entity-relationship.mmd) · [Interactive 3D](assets/entity-relationship.3d.html) · [3D snapshot](assets/entity-relationship.3d.png) · [Dark snapshot](assets/entity-relationship.3d.dark.png).

[![Authors submit documents with revisions and optional search index entries](assets/entity-relationship.svg)](assets/entity-relationship.svg)

## Mindmap

Three ownership branches and two leaf levels test branch visibility and the
site's explicit surfaces and text colors. Branches represent grouping, not data flow.
[Source](mindmap.mmd) · [Interactive 3D](assets/mindmap.3d.html) · [3D snapshot](assets/mindmap.3d.png) · [Dark snapshot](assets/mindmap.3d.dark.png).

[![Document service grouped into intake, processing, and operations](assets/mindmap.svg)](assets/mindmap.svg)

## Rebuild

From the repository root:

```sh
bazel_agent bazel build //projects/alwaldend.com/docs/mermaid-theme:samples
bazel_agent bazel build //projects/mermaid_3d/samples:html //projects/mermaid_3d/samples:snapshots
bazel_agent bazel build //projects/alwaldend.com:site
```

The first target produces six SVGs and their `.dark.svg` siblings. The second
produces six self-contained HTML files and twelve light/dark PNG snapshots. The site target packages this page and
all formats. Generated files remain build outputs; the six `.mmd` files are
authoritative. [Mermaid 3D](../../../mermaid_3d/README.md) owns the spatial renderer.
