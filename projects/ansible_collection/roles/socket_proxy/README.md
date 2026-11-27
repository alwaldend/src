---
title: Socket proxy
description: Activate a TCP proxy on demand with systemd
---

The socket listens on `socket_proxy_listen` and starts an unprivileged
`systemd-socket-proxyd` service when a client connects. The service forwards
the byte stream to `socket_proxy_target`. Supply both endpoints as addresses
with ports and set `socket_proxy_name` to the systemd unit basename.

The caller owns firewall policy and backend readiness.

Set `socket_proxy_selinux_allow_http_ports` on hosts that need the socket
proxy to bind or connect to TCP ports labeled `http_port_t`. The role
installs the policy before starting the socket on SELinux-enabled hosts.
The policy also lets systemd transition the unprivileged service into the
socket proxy domain with `NoNewPrivileges` enabled.
