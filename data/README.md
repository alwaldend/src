---
title: Data
description: Data tree
weight: 6
cascade:
  - categories:
      - data
---

This tree contains repository-owned data and documentation assets. All tracked
content follows the repository's public-source policy. Actual credentials,
personal information, and secret-bearing generated values remain prohibited.

## Requirements

- Bazel targets MUST use repository-internal visibility.
- The tree MUST NOT be independently published as a product artifact.
- Data and documentation assets MAY be used in builds and MAY be embedded in
  an explicitly owned published artifact when that owner permits it.
