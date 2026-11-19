---
title: OpenHands canvas
description: OpenHands Agent Canvas frontend
tags:
  - ansible_role
  - openhands
---

Runs the OpenHands Agent Canvas browser client from the published npm package.
The service listens on loopback port `openhands_canvas_port`; TLS, the public
hostname, and routing to the agent server and automation API belong to the
deployment's Traefik role.

The service runs the package's `scripts/static-server.mjs` directly with Node,
an explicit loopback host, and the package's `build/` directory. It serves the
browser client without starting an agent server, an automation server, or the
npm launcher's additional ingress process. Agent Canvas owns no conversation
or automation state and executes no tools.

The browser starts with no configured backend. The operator adds one with the
canvas origin as the host and the agent server's session key as the credential.
The static server receives no session key and injects none into served assets.
This configuration does not advertise an embedded VSCode editor.

The npm prefix is installed inside `openhands_canvas_dir`.
`openhands_canvas_package_dir` identifies the installed package containing the
static server and browser assets.
