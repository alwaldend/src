---
title: Graphviz DOT rules
description: Bazel rules for exporting DOT and rendering themed SVG diagrams
---

`dot_svg` renders a DOT file on Linux using Bazel's pinned Graphviz build.
The SVG embeds the common diagram font and has light and dark palettes.
Pass a checked-in `.dot` file as `src`, or generate one with `dot_export`
from a command that writes DOT to standard output.

```starlark
load("//tools/dot:defs.bzl", "dot_export", "dot_svg")

dot_export(
    name = "state_machine_dot",
    arguments = ["state-machine", "graph"],
    command = ":mlflow-vastai",
    out = "state-machine.dot",
)

dot_svg(
    name = "state_machine_svg",
    src = ":state_machine_dot",
    out = "state-machine.svg",
    target_compatible_with = ["@platforms//os:linux"],
)
```

The MLflow Vast.ai CLI uses these rules to build its command graph.
