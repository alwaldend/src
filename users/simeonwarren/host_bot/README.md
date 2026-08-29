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
public command-network access, and write access to the `~/.cache/bazel` and
`~/.cache/bazelisk` caches while keeping local and private network targets
blocked. Hosted web search is not restricted by this role.

The deployment installs pinned Bazelisk as `~/.local/bin/bazel` and provisions
the repository's `bazel_agent` runner in the same directory.

`T3_MCP_BEARER_TOKEN` authenticates the Codex app-server's loopback MCP
connection to T3 Code at `/mcp`. Codex excludes it from model-spawned commands
so those subprocesses do not inherit the credential.

## Deployment

```sh
bazel run //users/simeonwarren/host_bot/ansible
```
