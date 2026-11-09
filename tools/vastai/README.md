---
title: Vast.ai CLI
description: Vast.ai command-line client and public Forge image discovery
---

Runs the unmodified upstream `vastai` CLI with exact versions and hashes recorded
in [`uv.lock`](uv.lock). The manifest in [`pyproject.toml`](pyproject.toml) leaves
the Vast.ai version unconstrained so lock updates can select newer releases. The existing
[`Python dependency workflow`](../py/README.md) defines the update mechanism: `bazel run //tools/vastai:uv_lock.update`.
The component-local environment requires Python 3.11 or later and overrides
upstream's vulnerable Pillow and cryptography pins with secure release minimums.
The lock isolates these dependencies from the shared Python environment.

```sh
bazel run //tools/vastai -- --help
bazel run //tools/vastai -- search offers --help
```

This target does not install software on the host or supply credentials.
Commands that contact Vast.ai require network access; resource operations can
incur charges. Configure authentication through upstream's supported credential
mechanism before making authenticated requests. Keep API keys out of source,
command arguments, and captured output.

## Links

- [Upstream source and usage](https://github.com/vast-ai/vast-cli)
- [Upstream releases](https://github.com/vast-ai/vast-cli/releases)

## Forge image discovery

`//tools/vastai/cmd/latest_image` is the first-party adapter for Terraform's
external-data JSON protocol. It reads public Docker Hub metadata and returns the
newest published Vast Forge Neo/CUDA 12.9 image tag. It does not supply credentials,
create a template, or rent an instance. The bounded lookup fails for registry
errors or a missing matching image family; it does not substitute another runtime.

The personal [image-generation root](../../users/simeonwarren/imagegen/README.md)
packages this executable as runtime data. The executable-protocol fixture runs
without external network access and emits `latest-image.json`:

```sh
bazel test //tools/vastai/test/latest_image:lookup_test
```
