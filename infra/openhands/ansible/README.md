---
title: Ansible
description: OpenHands deployment
tags:
  - ansible
---

Deploys the OpenHands components. The inventory has three hosts in three
groups, and the playbook applies one component role per group:

| Group                     | Host                                       | Component              |
| ------------------------- | ------------------------------------------ | ---------------------- |
| `openhands_canvas`        | `host1.canvas.openhands.alwaldend.com`     | Agent Canvas           |
| `openhands_server_secure` | `host1.server.openhands.alwaldend.com`     | Agent server (secured) |
| `openhands_automation`    | `host1.automation.openhands.alwaldend.com` | Automation server      |

The three XCP-ng VMs exist only for the secured agent server and for the
components that must not share a trust boundary with host-bot. The secured
group sets `openhands_server_session_api_key` and hands the same value to the
automation server so the automation dispatcher can authenticate its
dispatched conversations.

Every component applies `host` and `traefik`. Canvas runs the published
package's static server directly on loopback. The operator configures the
browser backend with the canvas origin and the shared session key.

The canvas host routes `/api/automation` to the automation host and the agent
server API prefixes to the secured agent server. These routes take precedence
over the frontend route, so the configured browser reaches one origin.
Downstream routers accept the original canvas Host header.

The unsecured agent server runs on host-bot and is deployed by
[the host-bot project](../../../users/simeonwarren/host_bot/README.md), which
owns that host entirely. It omits `OH_SESSION_API_KEYS_0` and binds to loopback
without a session API key. Access is local to host-bot; its Traefik and firewall
configuration add no route for this server.

The group variables read the session key, secret key, and automation KV secret
from the environment the AL Vault injector populates.
