## Why

Deferred; implementation has not started. Review finding 2 identifies manual
renewal of a roughly month-long Vault server certificate; the October 2
recovery record demonstrates that expiry can block the authentication flow
needed to renew it. Future work should prevent routine expiry and make recovery
usable when normal certificate validation fails.

## What Changes

- Declare automatic renewal before expiry through the owning IaC workflow.
- Expose renewal failure and approaching expiry through existing operational
  reporting mechanisms, without logging keys or credentials.
- Document and rehearse an independently usable recovery path with explicit
  authentication prerequisites and invocation-scoped exceptions only.
- Accept with isolated renewal, failure, and recovery exercises plus packaged
  desired-state inspection; live success requires separate authorization and
  observed evidence.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-vault`: maintain server certificate validity and recovery readiness.

## Impact

This request authorizes planning only. Future code belongs to the reusable
Vault role and the Vault owner's scheduling, identity, and alert configuration.
No timer, alert destination, renewal, authentication, or deployment is activated
now; adding an external dependency would require specific approval.

Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[operator renewal procedure](../../../README.md#vault-certificates),
[issuance tasks](../../../../../projects/ansible_collection/roles/vault/tasks/main.yaml),
and the historical
[expired-certificate recovery](../archive/2026-10-02-recover-expired-server-certificate/design.md).
That historical incident does not establish current live certificate health.
