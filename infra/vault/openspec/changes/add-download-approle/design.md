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

- Declare `src_infra_download` in `infra/vault/tf` using the existing modules
  and infrastructure membership conventions. Keep all new source secret-free.
- Retain the AppRole's own state namespace and separate local, Yandex, and
  DNS state keys within that namespace. Do not grant another owner's state.
- Enable the existing Yandex folder policy and group linkage; expose identity
  outputs consumed by the Yandex folder and XCP-ng resource-set owners.
- Grant reads of the existing provider credentials only for required DNS
  views. Public ACME DNS-01 requires a credential able to change challenge
  records; distinguish Vault read access from the upstream token's DNS scope.
- Reuse the SSH host-signing path with the component's actual host identities.
  Interactive and automated publishers use the established SSH credential
  flow and known-host trust; possession of the component AppRole is not a
  public upload mechanism.
- Public website TLS is specified by the linked hosting design. Grant reads
  of the selected DNS challenge secret rather than adding a private PKI role
  merely because another component has one. Reuse private PKI only if a
  concrete management connection requires it during implementation.
- Component `al.lua` owns injection and backend selection. Matching Bazel
  plugin data and execution labels remain an explicit implementation check.

## Risks / Trade-offs

- Copying a broad sibling policy -> inspect the exact new paths and provider
  permissions, including denied access to unrelated component state.
- Provider-side bootstrap order -> create the identity before the Yandex and
  XCP-ng assignments, then configure the consuming component.
- DNS challenge credentials persist for renewal -> restrict their local
  permissions and inject references through the existing secret workflow.

## Migration Plan

Add source declarations and run offline checks first. A separately authorized
Vault apply creates the identity; provider owners then establish its resources.
The download deployment can authenticate only after those steps. Do not run
credential issuance or live permission probes as an offline validation step.
