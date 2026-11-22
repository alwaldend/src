---
title: Graphviz
description: Graphviz DOT renderer
---

The Bazel target builds Graphviz's `dot_static` executable on Linux with
built-in plugins for rendering DOT, including SVG output.

```sh
bazel run //third_party/org_graphviz_graphviz:dot -- -Tsvg < graph.dot > graph.svg
```

The source archive is pinned to the checksum published with the
[Graphviz release](https://graphviz.org/download/source/).
