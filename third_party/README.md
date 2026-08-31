---
title: Third party
description: Third party tree
weight: 7
cascade:
  - categories:
      - third_party
---

This tree contains vendored and externally sourced code. All tracked content
follows the repository's public-source policy and retains its upstream license
and provenance requirements.

## Requirements

- Bazel targets MUST use repository-internal visibility unless an owning
  toolchain contract requires otherwise.
- Third-party content MUST NOT be represented or published as first-party
  source or artifacts.
- Third-party content MAY be used in builds.
- An explicit repackaging or mirroring target MAY publish an upstream artifact
  under its owner and license policy.

## Publish vendored helm

This is a state-changing operator example. An agent needs explicit publication
authority and must invoke the target through `bazel_agent`.

```sh
bazel_agent run //third_party:publish_helm.io_goharbor_helm_harbor
```
