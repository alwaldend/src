## Context

The data plane uses explicit Terraform provider aliases and a shared node module.

## Goals / Non-Goals

Add work1 while preserving the existing node identities and credential injection.
Architecture rendering belongs to a separate pull request.

## Decisions

Host work1 at njal.la using the supplied IPv4 and IPv6 addresses. Give it its own
provider and module, using njalla1's Mullvad relay selection with the retired
Düsseldorf relay replaced by de-fra-wg-001 for work1.
Keep its credentials in the existing Vault maps.

## Risks / Trade-offs

The initial deployment must configure the new node before registering it with
the control plane. The user authorized applying Terraform only to
module.node_work1 while the master is unavailable. Its API token is obtained
after installation and is needed for later master registration.
