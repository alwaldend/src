---
title: Bazel configs
description: Bazel config targets
---

## Harbor

Targets that pull from the private Harbor registry are disabled by default so
repository-wide builds work while Harbor is unavailable. Enable them with
`--config=harbor` when the registry is reachable.
