---
title: Xen Orchestra OIDC plugin
description: Pinned upstream authentication plugin and runtime dependencies
---

This package supplies Xen Orchestra's AGPL-3.0-or-later `auth-oidc` 0.4.0
from upstream commit `4bc614ffd1417d3e0563522aeb4a195fa494c57c`, including
OIDC group synchronization and the single-group claim fix. It can be installed
independently of the XO server; the Ansible deployment verifies its runtime
before restarting XO.

The plugin is plain CommonJS. Its complete runtime dependency set is
`passport-openidconnect` 0.1.2 (MIT), `passport-strategy` 1.0.0 (MIT), and
`oauth` 0.10.2 (MIT). The npm archives retain upstream license files and use
integrity values from the appliance revision's
[upstream lockfile](https://github.com/vatesfr/xen-orchestra/blob/961b505cfb74cbf24aaeb8d61c5613f97e78de2c/yarn.lock).
Bazel fetches and verifies all inputs before deployment; Ansible extracts them
without a package-manager resolution or lifecycle script.

[Plugin source](https://github.com/vatesfr/xen-orchestra/tree/4bc614ffd1417d3e0563522aeb4a195fa494c57c/packages/xo-server-auth-oidc).
