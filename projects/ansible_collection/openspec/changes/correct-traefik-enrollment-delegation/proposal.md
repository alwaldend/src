## Why

The authorized Forgejo deployment stopped before its upgrade because Traefik's
controller-side ACME request inherited the managed host's connection variables
and privilege escalation. A local connection does not establish controller
delegation.

## What Changes

- Delegate the existing Vault enrollment command to localhost without privilege
  escalation, using the controller's runtime and existing injected credentials.
- Preserve private enrollment output and administrator-owned host cache writes.

## Capabilities

No specification delta. This repairs the existing controller-side enrollment
contract; `skip_specs: true` avoids adding a second contract.

## Impact

The reusable Traefik task and its documentation. A placeholder-only fixture,
written before the correction, reproduces failure with a remote-only Python
path and elevated managed-host context. Acceptance requires the actual task
to return its placeholder enrollment as the unprivileged controller user.
The user separately authorized the live Forgejo deployment; this fix adds
no deployment scope or dependency.
