## Context

Owner inventories are authoritative; live reachability does not prove activity or security.

## Decisions

Package owner inventory files under distinct paths, without deployment group variables. Import canonical Lynis and updater roles. Explicitly enforce automatic updates true and reboots false. Keep evidence private. Host Bot may use a task-specific local connection override only on its verified host.

Wrap audit and update imports in one block and compare boot identities in its
`always` section. Keep role tags on the imports and `always` tags on the boot
checks; do not tag the enclosing block or rescue an original failure. Ansible
transport loss remains unreachable and does not establish unchanged boot state.

The standalone SSH target uses the existing `al_binary_run` wrapper and a
quoted shell `exec` to read `HOST_MAINTENANCE_SSH_CONFIG` after injection.
That generated file owns trust, identity, and connection policy for both
maintenance entry points. Quote the whole shell command for Bazel argument
tokenization and quote the config path and forwarded arguments within it.
Require a nonempty config variable before executing SSH. Preserve host-PATH
SSH selection from `tools/ssh`, the AL config, injector dependency, and
maintenance plugin selector.

## Risks

Legacy appliances need separate transport and vendor maintenance. Scanners cannot establish absence of compromise. Package scriptlets can restart services.
