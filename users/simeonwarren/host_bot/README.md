---
title: Host Bot
description: host-bot.simeonwarren.users.alwaldend.com
tags:
  - ansible
  - t3code
  - traefik
---

Host Bot runs Traefik with mTLS in front of T3 Code. The host firewall blocks
direct external access to T3 Code and accepts new connections only on SSH,
HTTP, and HTTPS ports.

T3 Code controls permissions per thread. Its **Auto** mode maps to Codex
Auto-review. The managed Codex requirements allow read-only and `host-bot`
permissions with on-request approvals, so threads cannot select **Full
access** or disable approvals. The `host-bot` profile grants workspace writes,
public command-network access, and write access to the Bazel cache and Bazelisk
directories while keeping local and private network targets blocked. Hosted
web search is not restricted by this role.

The deployment installs pinned Bazelisk as `~/.local/bin/bazel` and provisions
the repository's `bazel_agent` runner in the same directory.

Two SATA SSDs, selected by their stable ATA IDs, back the `sata_ssd` volume
group. Its striped `bazel_cache` logical volume is mounted persistently at
`/var/cache/bazel`. Every Bazel workspace uses its `disk_cache` directory and
the managed user bazelrc points its output-user root at the same filesystem.
Automatic garbage collection caps the action disk cache at 700 GiB, leaving
room within the 1-TiB filesystem for output bases and install state; the
volume group retains about 900 GiB for later growth or another logical volume.
Because the volume is striped, either SSD failing invalidates the disposable
cache.

`T3_MCP_BEARER_TOKEN` authenticates the Codex app-server's loopback MCP
connection to T3 Code at `/mcp`. Codex excludes it from model-spawned commands
so those subprocesses do not inherit the credential.

## Deployment

```sh
bazel run //users/simeonwarren/host_bot/ansible
```
