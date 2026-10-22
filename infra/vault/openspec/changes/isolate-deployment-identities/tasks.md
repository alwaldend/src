Deferred; implementation has not started. These unchecked tasks are future
work; this planning request authorizes no implementation or deployment.

## 1. Confirm the access boundary

- [ ] 1.1 Inventory component identities, authorized hosts, SSH principals,
      certificate lifetimes, and sudo scope; deliver a source-linked allow/deny
      matrix and review the concrete migration through decision-review.
- [ ] 1.2 Establish an isolated SSH acceptance exercise before configuration
      changes; verify it exposes the current cross-component acceptance and
      records the expected permitted and denied connections without secrets.

## 2. Implement and verify the desired state

- [ ] 2.1 Restrict Vault signing and parameterize host principal authorization
      with owner-specific inputs; verify formatted Terraform and packaged
      Ansible agree with the matrix and preserve an authorized operator route.
- [ ] 2.2 Run the isolated matrix against the exact candidate, including
      certificate lifetime/extensions and retirement of broad grants; save a
      repeatable non-secret report showing permitted access and denied lateral
      access, then document the separately authorized rollout/recovery sequence.
