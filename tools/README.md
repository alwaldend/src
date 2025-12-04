---
title: Tools
description: Tools tree
weight: 4
cascade:
  - categories:
      - tool
---

This tree contains tools that are used inside the repo

## Requirements

- MUST NOT be public
- MUST NOT be published
- MUST NOT be used in builds
- MUST be available to the whole repo (`visibility = ["//:__subpackages__"]`)
- MAY be used in tests
