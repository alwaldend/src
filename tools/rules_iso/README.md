---
title: rules_iso
description: ISO image download and flash rules
statuses:
  - active
---

Module extension that downloads ISO images with `http_file` and adds a
runnable `iso_flash` target next to each image.

Develop and test the rules from the monorepo root:

```sh
bazel test //tools/rules_iso/...
bazel build //tools/rules_iso:release_archive
```

The release archive is `bazel-bin/tools/rules_iso/rules_iso.tar.gz`. It contains
the rules under `pkg/bzl`, their module declaration, and the repository license.
For a registered release, use [release.MODULE.bazel](release.MODULE.bazel) as
the basis for the version's registry `MODULE.bazel` and set its release
version. The registry's `source.json` identifies the archive URL and integrity
hash. The archive uses the module root directly and needs no strip prefix.

Registry consumers load the extension from
`@rules_iso//pkg/bzl:iso_extension.bzl`.
