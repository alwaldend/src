## Context

See [proposal.md](proposal.md). Existing component identities compose the
shared AppRole and SSH modules, join the infrastructure identity group, and
receive provider access through AL. The OpenHands identity disables Yandex
folder access because that component is local-only; the download identity
must support both environments and must not copy that limitation.

## Goals / Non-Goals

**Goals:** Give the download component only the state, provider, SSH, DNS, and
certificate-secret access its deployment uses.

**Non-Goals:** A new authentication backend, credential storage outside Vault,
host-wide credential migration, or administrative provider credentials.

## Decisions

- Declare `src_infra_download` in its own module at
  `infra/vault/approles/src_infra_download`, with its own HTTP backend and
  administrative AL authentication. Resolve shared identities by name with
  data sources; `infra/vault/tf` retains sole ownership of shared group
  membership lists and looks up the new entity/group after creation. Keep all new source secret-free.
- Retain the AppRole's own state namespace and separate local, Yandex, and
  DNS state keys within that namespace. Do not grant another owner's state.
- Enable the existing Yandex folder policy and group linkage; expose identity
  outputs consumed by the Yandex folder and XCP-ng resource-set owners.
- Grant reads of existing provider credentials for required DNS provisioning
  views. Certificate validation uses HTTP-01, not DNS provider credentials.
- Reuse the SSH host-signing path with the component's actual host identities.
  Interactive and automated publishers use the established SSH credential
  flow and known-host trust; possession of the component AppRole is not a
  public upload mechanism.
- Use Let's Encrypt for external traffic and Vault ACME/EAB for internal
  traffic, following the user's PR review correction. Declare a component PKI
  role with only the three serving names, no descendant names, and no client
  certificate flag. Grant EAB creation to the component AppRole group through
  the existing module. Existing consumers retain their PKI defaults.
- Component `al.lua` owns injection and backend selection. Matching Bazel
  plugin data and execution labels remain an explicit implementation check.

## Risks / Trade-offs

- Copying a broad sibling policy -> inspect the exact new paths and provider
  permissions, including denied access to unrelated component state.
- Provider-side bootstrap order -> create the identity before the Yandex and
  XCP-ng assignments, then configure the consuming component.
- EAB and ACME account state persist for renewal -> keep them in the existing
  private Traefik configuration and state directories; use the host CA trust.

## Migration Plan

Add source declarations and run offline checks first. A separately authorized
standalone identity-root apply creates the identity. Apply the core Vault
root afterward to resolve the new entity/group by name and add its shared
memberships; provider owners then establish their resources.
The download deployment can authenticate only after those steps. Do not run
credential issuance or live permission probes as an offline validation step.

The hosting owner updates its split-trust requirements and acceptance flow in
[PR #102](https://github.com/alwaldend/src/pull/102), following the explicit
[HTTP-01 issuer correction](https://github.com/alwaldend/src/pull/102#discussion_r4097007197).
That change replaces the older public-DNS-01 local-certificate plan: clients
using XCP-ng must trust the repository CA. This Vault PR is independently
based on trunk and does not restore the superseded DNS-01 design.
