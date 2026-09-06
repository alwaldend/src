---
title: Harbor login
description: Create a harbor session using OIDC
languages:
  - go
tags:
  - harbor
  - vault
  - oidc
---

The plugin destroys its Harbor session on normal shutdown and checks that the
original session ID is rejected by the current-user API. Cleanup uses Harbor's
`/c/oidc/logout` endpoint (available in newer Harbor releases) and does not follow
the optional identity-provider logout redirect. Unsupported endpoints, network
failures, or sessions that remain valid are reported as cleanup errors. Session
expiry remains server-controlled; forced termination cannot guarantee logout.

The session destruction behavior is defined by Harbor's
[OIDC controller](https://github.com/goharbor/harbor/blob/main/src/core/controllers/oidc.go).
Vault tokens created by the invocation are revoked after session cleanup.
