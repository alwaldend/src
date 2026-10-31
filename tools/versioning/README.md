---
title: Versioning
description: Repository versioning is owned by the release tool
---

Version calculation, source-ref publication, stamping, and project release
publication are owned by [tools/rules_al](../rules_al/README.md).
Use `al bazel run @rules_al//:al -- version show` and the shared
`tools/rules_al/cmd/al/al.sh version bazel -- ...` bootstrap.
