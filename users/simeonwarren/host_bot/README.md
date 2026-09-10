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
access** or disable approvals. The requirements file only allowlists the
`host-bot` profile; it must not define the profile because T3 Code supplies a
per-thread config layer with the same name, and Codex rejects profiles defined
by both requirements and config. The normal and isolated provider configs
define `host-bot` with workspace writes, public command-network access, and
write access to the Bazel cache and Bazelisk directories while keeping local
network binding and loopback connections available for Bazel servers. Other
private network targets remain blocked. T3 Code must select that profile for
the thread. Hosted web search is not restricted by this role.

The deployment installs pinned Bazelisk as `~/.local/bin/bazel`, provisions
the repository's `bazel_agent` runner in the same directory, and pre-seeds its
content-addressed `mcp_cordis` and `repo_delivery` runtimes when the canonical
repository checkout is present. New Codex worktrees can then start Cordis and
run delivery operations without first loading their own Bazel graph.

Two SATA SSDs, selected by their stable ATA IDs, back the `sata_ssd` volume
group. Its striped `bazel_cache` logical volume is mounted persistently at
`/var/cache/bazel`. Every Bazel workspace uses its `disk_cache` directory and
the managed user bazelrc points its output-user root at the same filesystem.
Automatic garbage collection caps the action disk cache at 700 GiB, leaving
room within the 1-TiB filesystem for output bases and install state; the
volume group retains about 900 GiB for later growth or another logical volume.
Because the volume is striped, either SSD failing invalidates the disposable
cache.

The managed user bazelrc limits each Bazel invocation to eight concurrent jobs
and budgets half of the host CPUs for local actions. These scheduling limits
apply to builds and tests; they do not impose a combined CPU cap across
independent Bazel processes or Codex sessions.

Executable repository tools use `/var/cache/bazel/tool_cache`. Unlike the
disposable action and output caches, this directory is private to the host-bot
account. Cache entries are keyed by declared source, build configuration,
platform, and dependency pins, and are installed atomically after a Bazel
build on the first miss.

The arXiv MCP server is pinned to version 0.7.2 and runs locally over stdio.
Its downloaded papers and search data live under
`/srv/misc/arxiv-mcp-server/papers` on a reusable 100-GiB `misc` logical
volume in the existing `sata_ssd` volume group. The volume is mounted for
reconstructible non-executable application data. Data with durability or
different reliability requirements belongs on a separate volume.

The managed Codex configuration allows up to 20 concurrent subagents per
session. The root agent is accounted for separately.

Node.js and npm are provisioned through the development VM role's mutable or
OSTree package backend. Open Computer Use 0.3.3 is then installed without npm
lifecycle scripts under the managed user's `~/.local` prefix. Its skill is
then installed for the whole host under `/etc/codex/skills` from the matching
immutable upstream revision. The deployment configures its stdio MCP server in
the normal Codex config and the isolated OpenRouter and Abliteration homes.
The MCP can inspect and control the logged-in desktop session, so agents must treat its actions as real user input and obtain approval before
externally visible or destructive actions.

OpenRouter is available through its isolated Codex home. Its
configuration reads its API key from the desktop keyring entry selected by
`service=openrouter` and `application=codex`; the key is not stored in the
configuration. Repository work uses the repo-owned `$codex-migration` skill, while
the user-level Codex instructions require isolated migration testing and
intentional host changes to be mirrored into this role. Existing sessions are
not redirected because provider-bound encrypted reasoning history cannot
safely migrate in place.

T3 Code can use OpenRouter as a separate Codex provider by setting its
`CODEX_HOME path` to `/var/home/simeonwarrenbot/.codex-openrouter`. That
isolated home defaults to `~deepseek/deepseek-v4-flash-latest` and shares no
conversation state with the ChatGPT-backed Codex provider. Start a new thread
when selecting it.

The OpenRouter no-tools home at `~/.codex-openrouter-no-tools` mirrors the
OpenRouter provider with every tool surface disabled. It reuses the same
keyring lookup and model, sets `web_search = "disabled"`, turns off the shell,
unified-exec, image, multi-agent, computer-use, browser, app, plugin, and goal
features, and disables the `request_user_input` tool and the MCP servers.
Codex then advertises no tools at all, so the model can only reply with text.
This matters beyond tidiness: OpenRouter rejects a request that carries any
tools when its routed endpoint supports none, so a single leftover tool fails
the turn.

T3 Code launches the Codex app server with
`-c mcp_servers.t3-code.url=...`, so its own MCP server would otherwise
reintroduce the MCP resource tools. The no-tools config therefore declares
`[mcp_servers.t3-code]` with `enabled = false`, which still wins over that
launch-time override.

Select it in T3 Code with `CODEX_HOME path` set to
`/var/home/simeonwarrenbot/.codex-openrouter-no-tools` and start a new thread.
The Abliteration and standard OpenRouter homes keep their tools.

Abliteration uses a separate home with
`CODEX_HOME="$HOME/.codex-abliteration" codex`, using the endpoint and model from its
[Codex integration guide](https://docs.abliteration.ai/integrations/codex).
It retrieves its API token at runtime with `secret-tool`, selecting
`service=abliteration` and `application=codex`. Save the token from an
interactive terminal in the same user's desktop keyring:

```sh
secret-tool store --label='Abliteration API token for Codex' service abliteration application codex
```

Enter the token at the hidden prompt. The configuration contains only the
lookup command. Provider response and tool-call validation require a saved
token. Both providers keep their configuration and conversation state outside
the main Codex home; provider profiles in the main home are removed. For
T3 Code, set its `CODEX_HOME path` to
`/var/home/simeonwarrenbot/.codex-abliteration` and start a new thread.

`T3_MCP_BEARER_TOKEN` authenticates the Codex app-server's loopback MCP
connection to T3 Code at `/mcp`. Codex excludes it from model-spawned commands
so those subprocesses do not inherit the credential.

## Deployment

```sh
bazel run //users/simeonwarren/host_bot/ansible
```
