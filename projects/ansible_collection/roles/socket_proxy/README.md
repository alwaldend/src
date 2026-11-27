---
title: Socket proxy
description: Activate a TCP proxy on demand with systemd
---

The socket listens on `socket_proxy_listen` and starts an unprivileged
`systemd-socket-proxyd` service when a client connects. The service forwards
the byte stream to `socket_proxy_target`. Supply both endpoints as addresses
with ports and set `socket_proxy_name` to the systemd unit basename.

The caller owns firewall policy and backend readiness.
