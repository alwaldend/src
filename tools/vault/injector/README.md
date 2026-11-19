---
title: Injector
description: Secret injector
languages:
  - go
tags:
  - vault
---

Shutdown first drains plugin requests, then stops and waits for resource
processes, revokes invocation-owned Vault credentials, and deletes temporary
files and SSH key directories. Failed cleanup is reported. Registration after
shutdown is rejected; a fetcher removes any unregistered temporary material.
Temporary files use mode `0600` and directories use `0700`. Deletion is filesystem
unlinking, not secure erasure, and cannot run after `SIGKILL` or host failure.

`no_auth` explicitly sets the injected `VAULT_TOKEN` to an empty value so an
inherited token is overridden. This does not remove the user's token-helper
file; commands that independently consult that helper may still authenticate.
Template errors and OIDC status errors omit input and response contents.
OIDC requests honor cancellation and do not follow redirects.
