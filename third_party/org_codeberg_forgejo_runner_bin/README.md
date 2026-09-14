---
title: org_codeberg_forgejo_runner_bin
description: Forgejo Actions runner binary
tags:
  - forgejo
---

## Links

- Site: https://forgejo.org/
- Repo: https://code.forgejo.org/forgejo/runner

The [binary lock](binary_toolchain.json) pins the runner version and publisher
checksum. Native Forgejo Actions OIDC support is required by the
[runner deployment](../../infra/forgejo_runner/README.md).

Validate a workflow without executing jobs:

```sh
bazel_agent bazel run //third_party/org_codeberg_forgejo_runner_bin -- validate --workflow --path "$PWD/.forgejo/workflows/secure.yaml"
```
