---
title: Users
description: User tree
weight: 8
cascade:
  - categories:
      - user
---

This tree contains user-specific code and infrastructure. Tracked source
follows the repository's public-source policy. Personal information and
secrets must not be tracked or disclosed. Non-secret, non-personal operational
facts remain public even when they describe generated or live state.

## Requirements

- Bazel targets MUST use repository-internal visibility.
- User artifacts MUST NOT be published.
- User targets MUST NOT be dependencies of production build targets.
