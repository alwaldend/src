---
title: Third party
description: Third party tree
weight: 7
cascade:
  - categories:
      - third_party
---

This tree contains third party code

## Requirements

- MUST NOT be public
- MUST NOT be published
- MAY be used in builds

## Publish vendored helm

```sh
bazel run //third_party:publish_helm.io_goharbor_helm_harbor
```
