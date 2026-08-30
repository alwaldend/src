---
title: Mermaid
description: Bazel-managed Mermaid diagram renderer
---

This package exposes the Mermaid CLI as a repository-wide Bazel tool. Bazel
provisions the pinned Node toolchain, JavaScript dependency graph, Mermaid CLI,
and Chrome-for-Testing browser used by Puppeteer. Rendering runs as a Bazel
target or action; it does not use a host browser or download one during an npm
lifecycle hook.

The JavaScript launcher runs from Bazel's output tree. Pass absolute paths when
rendering source-tree files directly:

```sh
repo_root="$PWD"
bazel_agent run //tools/mermaid:mmdc -- \
  -i "${repo_root}/path/to/diagram.mmd" \
  -o "${repo_root}/path/to/diagram.svg"
```

Prefer a `mermaid_svg` action plus `write_source_file` for maintained diagrams;
those targets use declared Bazel paths and need no absolute-path handling.
