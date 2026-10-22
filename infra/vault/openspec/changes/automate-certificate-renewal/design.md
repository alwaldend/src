## Context

Deferred; implementation has not started. Authority is planning only. The
[proposal](proposal.md) links the current month-long issuance request, manual
renewal instructions, and October 2 historical recovery. No current endpoint
health, scheduler, or credential availability has been observed for this plan.

## Goals / Non-Goals

**Goals:** Renew before expiry, retain working material on failure, surface
expiry risk, and make recovery prerequisites usable independently of normal
TLS validation of the expired endpoint.

**Non-Goals:** Deploy a scheduler, renew a live certificate now, persist a TLS
bypass, introduce a monitoring service, or redesign the PKI hierarchy.

## Decisions

- Extend the existing Vault role and certificate-only workflow with declarative
  scheduling and an explicit renewal window. A reminder alone leaves renewal
  dependent on operator attention; ad hoc cron changes would leave desired
  state incomplete. Reuse existing host scheduling facilities after inspection.
- Keep ordinary renewal strict about TLS and authentication. Before adopting
  any unattended identity, specify its minimum issue/read scope and credential
  lifetime through the existing Vault references and injection mechanisms.
  Store no key or token in source or diagnostic artifacts.
- Validate the replacement certificate/key pair before activating it and
  preserve the previous usable pair until activation succeeds. Treat reload
  and served-certificate verification as separate observable steps.
- Use the existing operational reporting mechanism for renewal failures and
  expiry thresholds. Select its concrete owner and destination in the initial
  deferred design review; external service additions need separate approval.
- Retain an explicitly scoped recovery procedure with independent bootstrap
  prerequisites. The existing invocation-only verification exception is
  historical evidence, not permission to weaken ordinary renewal defaults.

## Risks / Trade-offs

- Unattended authentication can fail before the certificate expires -> include
  credential lifetime and renewal failure in the same acceptance exercise.
- A scheduler running on Vault may share its outage domain -> assess the
  failure domain and ensure the recovery entry point can be invoked separately.
- Issuance succeeds but installation or reload fails -> retain usable material
  and report the failed stage without disclosing secret values.

## Migration Plan

Next action after implementation is requested: inspect existing schedulers,
reporting, and authentication, then review a concrete mechanism using
`decision-review` before code. Package its declarative configuration and run
isolated success, issuance failure, activation failure, and expired-endpoint
recovery exercises. Preserve a non-secret result artifact with expiry and
stage outcomes. Deployment and live recovery each need an exact request;
source acceptance must not be described as live renewal success.
