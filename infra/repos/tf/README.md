---
title: Repository catalog module
---

This provider-free Terraform module reads the [catalog files](../README.md), merges
forge defaults, and projects organization-owned repository names and settings
for each consumer. It creates no resources and uses no credentials.

Fork and mirror names follow the [catalog naming contract](../README.md).
Output keys retain the organization and catalog key, independently of the
derived destination name. The consumers retain their existing state addresses
where resources were already managed.

The catalog test runs without network access or provider initialization:

```sh
bazel_agent bazel test //infra/repos/tf:tf_tests.catalog_test
```
