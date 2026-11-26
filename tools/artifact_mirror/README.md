---
title: Artifact mirror
description: Stage pinned Helm charts and OCI images and publish them to Forgejo
---

`//tools/artifact_mirror` reads a JSON inventory of Helm chart archives and
OCI images. It verifies staged digests before publishing and verifies the
destination after copying. OCI copies preserve multi-platform indexes.

[![Actor pulls from GitHub Releases and GHCR, authenticates with Vault, and pushes charts and images to Forgejo](assets/transfer.svg)](assets/transfer.svg)

```sh
bazel run //tools/artifact_mirror -- pull --inventory /absolute/path/artifacts.json
bazel run //tools/artifact_mirror -- push --inventory /absolute/path/artifacts.json --owner ORGANIZATION
```

Direct invocations use `helm` and `oras` from `PATH`. Override either with
`--helm` or `--oras`. The default staging directory is `out/artifact-mirror`.
Push uses `FORGEJO_HOST` and `FORGEJO_API_TOKEN` from its caller.

The [third-party mirror targets](../../third_party/README.md#mirror-release-artifacts)
supply the pinned executables, package inventory, and Forgejo login.
