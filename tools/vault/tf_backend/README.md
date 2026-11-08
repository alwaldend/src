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

Use `vault_tf_backend(path = ..., labels = ...)` from `:defs.bzl` to construct
an `al_config` plugin-call dictionary. The helper owns the plugin binding and
the default `secrets` mount; callers own their backend path and selection labels.
Optional `vault_auth`, `vault_conn`, and `mount` select other named inputs.
The shared `//tools/al:config` registers the executable and propagates runfiles.
