---
title: PVE login
description: Get a login ticket using OIDC
languages:
  - go
tags:
  - pve
  - vault
---

The plugin requests an API token with a one-hour expiry and deletes that token
on shutdown using the retained login ticket and CSRF token. The token name is
registered for cleanup before creation, so cleanup is attempted even when a
creation response is lost or malformed. Failures to delete are reported; expiry
is the fallback when shutdown cannot complete. Pre-existing credentials are not
revoked. Vault tokens created by the invocation are revoked after token cleanup.

Token deletion uses the Proxmox
[user token API](https://github.com/proxmox/pve-access-control/blob/master/src/PVE/API2/User.pm).
