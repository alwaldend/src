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

The arXiv MCP server is pinned to version 0.7.2 and runs locally over stdio.
Its downloaded papers and search data live under
`/srv/misc/arxiv-mcp-server/papers` on a reusable 100-GiB `misc` logical
volume in the existing `sata_ssd` volume group. The volume is mounted for
reconstructible non-executable application data. Data with durability or
different reliability requirements belongs on a separate volume.

The managed Codex configuration allows up to 20 concurrent subagents per
session. The root agent is accounted for separately.

OpenRouter is available to new sessions through `codex -p openrouter`. The
profile reads its API key from the desktop keyring entry selected by
`service=openrouter` and `application=codex`; the key is not stored in the
profile. Repository work uses the repo-owned `$codex-migration` skill, while
the user-level Codex instructions require isolated migration testing and
intentional host changes to be mirrored into this role. Existing sessions are
not redirected because provider-bound encrypted reasoning history cannot
safely migrate in place.

T3 Code can use OpenRouter as a separate Codex provider by setting its
`CODEX_HOME path` to `/var/home/simeonwarrenbot/.codex-openrouter`. That
isolated home defaults to `~deepseek/deepseek-v4-flash-latest` and shares no
conversation state with the ChatGPT-backed Codex provider. Start a new thread
when selecting it.

`T3_MCP_BEARER_TOKEN` authenticates the Codex app-server's loopback MCP
connection to T3 Code at `/mcp`. Codex excludes it from model-spawned commands
so those subprocesses do not inherit the credential.

## Deployment

```sh
bazel run //users/simeonwarren/host_bot/ansible
```
