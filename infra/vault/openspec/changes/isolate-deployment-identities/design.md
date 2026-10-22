## Context

Deferred; implementation has not started. Authority is planning only. See the
[proposal](proposal.md) for the source finding and desired outcome. The current
candidate is a future desired-state change; no live permission audit, certificate
issuance, or deployment has been performed.

The shared signer currently defaults to `ansible`, permits that principal for
multiple component identities, and uses month-scale lifetimes. A distinct
AppRole alone cannot constrain a certificate accepted by unrelated hosts.

## Goals / Non-Goals

**Goals:** Define issuer and host checks as one access boundary, preserving
necessary operator access and keeping environment inputs with their owners.

**Non-Goals:** Rotate the shared CA, redesign every Vault policy, change human
access by implication, or execute infrastructure operations during planning.

## Decisions

- Use component-scoped certificate principals and explicit host-side principal
  authorization as the candidate approach. Merely splitting Vault role names
  leaves the common accepted principal unchanged; separate CAs add unnecessary
  rotation and trust-distribution work unless later evidence requires them.
- Parameterize the existing shared host role rather than copy SSH configuration
  into every component. Component owners supply permitted principals; Vault
  owns the corresponding signing restrictions. Reuse the existing secret
  injection flow and declare host_bot inputs in its Ansible when affected.
- Select certificate maximum lifetimes and extensions from measured job needs
  during the first deferred task. Prefer a short deployment window, but do not
  guess a universal duration before inventorying long-running jobs.
- Review the concrete migration through `decision-review` when implementation
  is requested. This record commits to observable isolation, not unexamined
  per-host values or a premature production rollout.

## Risks / Trade-offs

- Host authorization and signer rollout in the wrong order can lock out
  deployment -> exercise the complete allow/deny matrix in isolation and retain
  a documented authorized operator route before removing broad grants.
- Existing valid broad certificates can outlive the policy update -> account
  for their lifetime and effective host acceptance in migration criteria.
- An incomplete component-to-host map can conceal legitimate shared access ->
  record shared administration explicitly rather than silently widening roles.

## Migration Plan

Next action, after implementation is requested: inventory desired consumers,
host principals, lifetimes, and effective sudo boundaries from source. Prepare
the matching Vault and Ansible candidate, then validate rendered configuration
and an isolated SSH allow/deny exercise with a repeatable non-secret report.
Any live rollout needs an exact operation-and-host request. Stage host and
signer transitions together; document the previous compatible policy as the
recovery candidate without automatically restoring broad grants.
