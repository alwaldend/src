## Why

Upstream fixed POSIX-shell compatibility and restored Flatcar detection in the
K3s installer after the release packaged by the Ansible collection.

## What Changes

Pin the newer installer by upstream commit and checksum. Keep the Ansible role
on the existing K3s binary release.

## Capabilities

The existing Ansible role contract applies.

## Impact

The shared K3s role packages the updated script at its existing install path.
