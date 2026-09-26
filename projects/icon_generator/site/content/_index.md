---
title: Icon generator
linkTitle: Icon generator
description: Make random pixel patterns with your own colors and shapes
layout: landing
statuses:
  - active
languages:
  - go
tags:
  - cli
  - images
---

Create PNG icons from random arrangements of colored squares or circles.
Choose a palette, background, canvas size, and shape size, then adjust the
density to leave more or less of the background visible.

Choose compact groups, winding strands, or replication. In replication mode,
density sets the starting population, and each cell can produce two neighboring
children, which can reproduce in turn.

Use a transparent background for images you want to place over other
content. Save a seed to reproduce a pattern with the same settings.

```sh
icon_generator \
  --width 256 --height 256 \
  --background transparent \
  --colors '#223843,#d77a61,#e3b23c' \
  --pixel-size 16 --pixel-shape circle \
  --density 0.4 --seed 42 \
  --output icon.png
```

The command runs locally and generates the PNG directly. It needs no image
editor or external graphics programs. Existing output files are preserved.

[Usage and options](/docs/projects/icon_generator/)
