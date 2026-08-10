---
title: T3 Code
description: t3code.simeonwarren.users.alwaldend.com
tags:
  - ansible
  - ai
  - t3code
---

T3 Code runs as a systemd service on a single Yandex Cloud VM. Traefik on the
same VM terminates TLS and forwards HTTP and WebSocket traffic to the
loopback-only T3 Code server.

## Deployment

```sh
bazel run //users/simeonwarren/t3code/tf_setup
bazel run //users/simeonwarren/t3code/ansible
```
