---
title: Tf backend
description: Http terraform backend backed by Vault
languages:
  - go
tags:
  - terraform
  - vault
---

The plugin tracks each running backend for shutdown, including backends created
before a later call fails. Shutdown drains HTTP requests before revoking the
backend's invocation-owned Vault credentials. If the shutdown deadline expires,
connections are forcibly closed and the timeout is reported. Normal HTTP server
closure is not an error. Vault requests use their incoming request context.
