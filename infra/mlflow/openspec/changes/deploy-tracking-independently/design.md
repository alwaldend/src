## Decision

Keep Vault OIDC enabled for tracking. Select the coordinator and its
credential injection only through the explicit coordinator deployment.

Use the observed Xen data disk attachments. The coordinator disk is at
`/dev/xvde` and the Traefik disk is at `/dev/xvdf`. Keep the preflight
checks on the declared device names and capacities.
