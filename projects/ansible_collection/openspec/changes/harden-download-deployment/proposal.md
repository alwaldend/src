## Why

Review of the download deployment found inactive Fedora updates, disruptive
Traefik reruns, writable proxy configuration, incomplete SSH certificate
activation, a data-disk boot dependency, stale firewall permissions, and
unvalidated proxy configuration. Fix the owning roles before public rollout.

## What Changes

- Activate Traefik changes deliberately, preserve unchanged service processes,
  validate configuration, and restore prior configuration after failed activation.
- Protect proxy configuration and its ancestors from the service identity.
- Reload SSH when signed host certificates change.
- Allow OS boot without the download disk while requiring it for content services.
- Explicitly delete retired firewall rules and remove download's
  inherited application-profile allowances.
- Reuse trunk’s canonical package-update role for reboot detection and native
  automatic updates, retaining its explicit reboot opt-in.

## Capabilities

### New Capabilities

- `host-deployment-safety`: Reliable activation and privilege boundaries for
  shared host roles used by download deployments.

### Modified Capabilities

None.

## Impact

Changes belong to shared roles and collection plugins. Download selects its
firewall adoption policy and consumes the revised storage behavior. The user
approved native dnf5-plugin-automatic, dnf-automatic, and unattended-upgrades
packages and deployment to both download hosts. The separate planned
Molecule runner remains independently owned; focused acceptance uses existing
cached fixture software and records its limitations.
