---
title: Openhands server
description: OpenHands agent server
tags:
  - ansible_role
  - openhands
---

Runs the OpenHands agent server natively from a `uv`-managed virtual
environment. It hosts conversations, executes tools, and streams events for
Agent Canvas. The service binds loopback only; TLS and external reachability
belong to the calling deployment.

Set `openhands_server_secure` to `true` to require a session API key. The
unsecured mode omits the session key and is only appropriate where the loopback
boundary is trusted.

The caller can set `openhands_server_install_packages` to `false` and install
`openhands_server_packages` through its host-specific package workflow.
OSTree hosts use the shared `dev_vm` package entry point and must boot into a
deployment containing those dependencies before installing the service. They
also set `openhands_server_service_path` under `/etc/systemd/system`.
