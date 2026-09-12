---
title: OpenHands automation
description: OpenHands automation server
tags:
  - ansible_role
  - openhands
---

Runs the OpenHands automation server natively from a `uv`-managed virtual
environment. It stores schedules and event triggers, tracks run lifecycle, and
dispatches conversations to the agent server.

The automation server is not a sandbox. It records which agent server should
run a dispatched conversation and hands that work to
`AUTOMATION_AGENT_SERVER_URL`; setting that variable selects the agent
server's local mode, which uses a persistent local agent server instead of
managed OpenHands Cloud sandboxes.

`openhands_automation_agent_server_api_key` must equal the agent server's
session API key because both services authenticate the same
`X-Session-API-Key` value. `openhands_automation_local_api_key` is the key
browsers present to this service.

The service stores state in SQLite via `AUTOMATION_DB_URL` by default, which
local mode supports without a separate database role. Set
`openhands_automation_database_url` to a PostgreSQL URL to use one instead.
`openhands_automation_kv_secret` enables the automation key-value store and is
required for automations that persist state between runs.

`openhands_automation_base_url` is the externally reachable origin the service
appends `/api/automation` to when it builds callback URLs, so it must be the
origin the browser and the dispatched agent server side both use.
