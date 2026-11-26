## Context

The K3s stable channel reported v1.37.2+k3s1 on 2026-11-26. This is one
Kubernetes minor version above the previous repository pin. The installer
compares the binary and service-file hashes before and after its own work.
The Ansible role already copies the binary into place before invoking the
installer with downloads disabled.

## Decisions

Update the binary and source archive together so the shipped installer matches
the binary release. Pass `INSTALL_K3S_FORCE_RESTART=true` when the binary copy
changes. For unchanged binaries, the installer can still start or restart the
service when its own files change.

## Rollout

Check each cluster's running version and the Kubernetes v1.37 upgrade notes
before deploying. K3s requires sequential minor upgrades and server-first
ordering. This change is validated locally.

- [K3s v1.37 release notes](https://docs.k3s.io/release-notes/v1.37.X)
- [K3s upgrade procedure](https://docs.k3s.io/upgrades/manual)
