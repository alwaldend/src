---
title: Traefik
description: Traefik
tags:
  - ansible_role
---

The administrator owns the proxy directory, configurations, TLS inputs, binary,
and unit. The service account can write only its data and log directories;
custom configuration and executable paths must have administrator-owned
ancestors and stay outside those writable directories. Static configuration
can contain credentials and is installed with mode `0640`.

Both rendered configurations pass controller-side TOML/YAML syntax validation
before installation. Startup inputs (static configuration, binary, trust CA,
and unit) restart the service only when changed. Dynamic routing uses Traefik's
file watcher. Unchanged runs keep the process running; the unit has no SIGHUP
reload action. Syntax validation does not prove route semantics or external
issuer reachability.

Installation retains temporary backups. If staging or startup fails, available
prior files are restored and a previously active service is restarted before
the play reports failure. A first installation has no previous deployment to
restore. Failed recovery preserves its backups; successful activation or
recovery removes them. This covers startup rejection, not failures after a
successful startup or invalid dynamic routing accepted by the process.

Local ACME enrollment is cached in administrator-only `eab.json` (mode `0600`).
Unchanged deployments reuse it without issuing another Vault request. The first
upgrade creates the cache once. Changing the issuer path or setting
`traefik_rotate_eab: true` requests a new binding; reset the flag afterwards.
If the ACME account state is lost, rotate its single-use enrollment binding
before re-enrolling. Keep the cache and ACME account data private.
