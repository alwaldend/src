---
title: Mermaid theme gallery
description: Six diagram types rendered with the current flat site theme
---

These fictitious document-service examples share the site's current
[flat theme](https://github.com/alwaldend/src/blob/master/projects/alwaldend.com/pkg/mermaid/theme.json) and
[palette](https://github.com/alwaldend/src/blob/master/projects/alwaldend.com/assets/scss/mermaid.scss). The site selects light or dark
images with its theme menu. Open an image to inspect its native size.
The examples carry structure and labels only; appearance belongs to the theme.
Labels use embedded Liberation Mono, including when an image is opened alone.
Normal graph paths use 2.25px strokes, captions have muted frames, sequence
guides recede behind messages, and mindmap branches keep visible contrast in
both modes. Dense diagrams still benefit from opening the full-size image.

## Flowchart

Branching decisions, a container, two data stores, and dashed revision/retry
paths test arrow direction, caption weight, and route separation.
[Source](flowchart.mmd).

[![Document processing with validation, retry, search, and audit paths](assets/flowchart.svg)](assets/flowchart.svg)

## Sequence

Calls, dashed replies, an activation, a self-call, an alternative, a loop,
and a note test the distinction between messages and structural guides.
[Source](sequence.mmd).

[![Document submission followed by repeated queue delivery and validation](assets/sequence.svg)](assets/sequence.svg)

## State

A composite processing state, retry cycle, cancellation, and terminal markers
test transition labels and nested boundaries. [Source](state.mmd).

[![A queued document progresses through validation and indexing or retries after failure](assets/state.svg)](assets/state.svg)

## Class

An interface, implementation, dependency, association, and composition test
open and filled markers, dashed lines, and cardinality labels.
[Source](class.mmd).

[![Repository interface, search implementation, indexer, documents, and revisions](assets/class.svg)](assets/class.svg)

## Entity relationships

Identifying and non-identifying relationships, optional cardinality, and
attribute tables test crow's-foot markers and fine table separators.
[Source](entity-relationship.mmd).

[![Authors submit documents with revisions and optional search index entries](assets/entity-relationship.svg)](assets/entity-relationship.svg)

## Mindmap

Three ownership branches and two leaf levels test branch visibility and the
site's explicit surfaces and text colors. Branches represent grouping, not data flow.
[Source](mindmap.mmd).

[![Document service grouped into intake, processing, and operations](assets/mindmap.svg)](assets/mindmap.svg)

## Rebuild

From the repository root:

```sh
bazel_agent bazel build //projects/alwaldend.com/docs/mermaid-theme:samples
bazel_agent bazel build //projects/alwaldend.com:site
```

The first target produces six SVGs and their `.dark.svg` siblings. The second
packages this page and its images. Generated images remain build outputs;
the six `.mmd` files are authoritative.
