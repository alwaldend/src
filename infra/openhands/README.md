---
title: OpenHands
description: openhands.alwaldend.com
tags:
  - openhands
  - traefik
  - terraform
  - ansible
---

OpenHands deploys as separate components because each owns a distinct
responsibility and trust boundary:

| Component         | Responsibility                                                           |
| ----------------- | ------------------------------------------------------------------------ |
| Agent Canvas      | Browser client for conversations, files, settings, backends, automations |
| Agent Server      | Runs conversations, agents, tools, and workspace operations              |
| Automation Server | Stores schedules and triggers, tracks runs, dispatches conversations     |

The agent server runs directly on the host with the service account's
permissions. Its workspace is a working directory; operating-system account
permissions define its access to the host. Treat the agent server host as
trusted infrastructure.

## Layout

- `tf_setup` creates the Agent Canvas, secured agent server, and automation
  server VMs in the `src_infra_openhands` Xen Orchestra resource set.
- `ansible` deploys the three VM components. The
  [host-bot project](../../users/simeonwarren/host_bot/README.md) owns and
  deploys its unsecured agent server.
- `dnsconfig.json` owns the component hostnames and static addresses.
- `al.lua` authenticates with the `src_infra_openhands` Vault AppRole, signs
  host keys for Ansible, and injects the component secrets.

There is no `tf` package: these components keep their state in files and
SQLite on their own hosts, so no API-side resources exist to declare.

Canvas serves the packaged browser assets through its
[static-server service](../../projects/ansible_collection/roles/openhands_canvas/README.md)
on loopback. Configure the browser's backend with the canvas origin as its
host and the agent server session key as its credential.

## Traefik and TLS

Each component host runs Traefik from the shared role with certificates issued
by `src_infra_openhands_pki_server`. The canvas origin proxies `/api/*`,
`/api/automation`, `/sockets`, and `/server_info` to the agent server and
automation server, so the browser talks to one origin and issues no
cross-origin request.

The public `canvas.openhands.alwaldend.com` record is a CNAME to
`ingress.alwaldend.com`, and `infra/ingress` forwards it to
`dc1.canvas.openhands.alwaldend.com`. That ingress keeps its
`RequireAndVerifyClientCert` client-auth policy, so reaching canvas from
outside still requires a client certificate.

## Secrets

`agent_server` holds `session_api_key` and `secret_key`; `automation` holds
`kv_secret`. The session API key is shared by the agent server, the automation
server, and the browser because all three authenticate the same
`X-Session-API-Key` value.

The unsecured host-bot server has a separate `secret_key` reference owned by
[its host's AL configuration](../../users/simeonwarren/host_bot/al.lua).

Prepare one private JSON input for each path through the
[secret-handling workflow](../../projects/agents/skills/repo-secrets/SKILL.md).
Set `OPENHANDS_AGENT_SERVER_JSON`, `OPENHANDS_AUTOMATION_JSON`, and
`OPENHANDS_HOST_BOT_JSON` to those input file paths. Each JSON object contains
the fields listed above for its path. Initialize the paths through standard
input; `-cas=0` prevents overwriting existing data:

```sh
bazel_agent bazel run //infra/openhands:vault.kv_put -- -cas=0 \
  alwaldend.com/vault1/approles/src_infra_openhands/agent_server \
  - < "$OPENHANDS_AGENT_SERVER_JSON"
bazel_agent bazel run //infra/openhands:vault.kv_put -- -cas=0 \
  alwaldend.com/vault1/approles/src_infra_openhands/automation \
  - < "$OPENHANDS_AUTOMATION_JSON"
bazel_agent bazel run //users/simeonwarren:vault.kv_put -- -cas=0 \
  alwaldend.com/vault1/approles/user_simeonwarren/openhands \
  - < "$OPENHANDS_HOST_BOT_JSON"
```

## Deployment

```sh
bazel run //infra/openhands/tf_setup:tf_setup.plan
bazel run //infra/openhands/tf_setup:tf_setup.apply
bazel run //infra/openhands/ansible
```

The component AppRole must first exist in Vault. Follow the
[XO authentication and resource-set workflow](../xcp_ng/tf/README.md) to
synchronize its OIDC identity and apply its resource-set assignment before
`tf_setup` runs. The assignment in
[`resource_set_inventory`](../xcp_ng/tf/inventory.tf) supplies the template,
storage, and network by name. Component secrets must be written to Vault
before Ansible runs. State and locking use the repository's `tf_backend`
plugin, which stores them as Vault KV entries and creates the entries on first
use, so no separate state bucket is required.
