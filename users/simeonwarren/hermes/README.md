---
title: Hermes
description: hermes.alwaldend.com
tags:
  - ansible
  - pve
  - hermes
---

Hermes Agent runs as a systemd service on a single Proxmox VM. Traefik on the
same VM terminates TLS and forwards gateway HTTP and WebSocket traffic to the
loopback-only Hermes API server.

## Deployment

```sh
bazel run //users/simeonwarren/hermes/tf_setup
bazel run //users/simeonwarren/hermes/ansible
```
