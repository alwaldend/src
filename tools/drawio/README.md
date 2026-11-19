---
title: Drawio
description: Bazel rules for drawio
languages:
  - bzl
tags:
  - bzl_rules
  - drawio
---

## Links

- https://github.com/rlespinasse/docker-drawio-desktop-headless

## Reproducible SVGs

`drawio_svg` in `defs.bzl` renders named pages from a `.drawio` file with the
pinned Drawio web exporter and the existing pinned Chrome-for-Testing browser.
The extraction action unpacks the web application from the Drawio AppImage;
rendering uses a small Puppeteer bridge for its export messages. Neither action
uses a host display server or accesses the network. JavaScript keeps this
adapter close to Drawio's browser API and the existing Puppeteer dependency.

Use the rule from an explicit source-update target with `write_source_files`,
as in `//infra/arch:update`. Normal documentation consumes the maintained SVGs,
not the rendering tool. The generated freshness tests catch stale output.

Exports use a light theme with an opaque white canvas, keeping labels and
connectors readable in both light and dark documentation themes and standalone
image viewers. The renderer adds this canvas on every export.

Chrome's inner sandbox is disabled because it cannot nest within the Bazel
Linux sandbox used for these actions. The Bazel sandbox remains enabled;
requests outside local files and embedded data are rejected. The renderer is
intended for checked-in, reviewed diagram sources and propagates export failures.

The repeat-render test checks byte equality, including stable SVG identifiers.
Rendering uses only the declared [Liberation Fonts 2.1.5](https://github.com/liberationfonts/liberation-fonts/releases/tag/2.1.5)
through an isolated Fontconfig configuration. Helvetica/Arial map to Liberation
Sans, Times to Liberation Serif, and Courier to Liberation Mono. The fonts are
licensed under the SIL Open Font License 1.1; their archive retains `LICENSE`.
