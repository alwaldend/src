## Why

The heat investigation found the manually selected `throughput-performance`
profile with minimum CPU performance at 100%. The user requested `balanced`
and explicitly required setting it without Ansible.

## What Changes

- Select `balanced` directly on this host using `sudo -n tuned-adm profile balanced`.
- Verify the active profile and resulting CPU power policy.
- Declare `balanced` in the owning Ansible variables and idempotently select it during deployment.
- Document the managed setting and the original direct application.

## Capabilities

### New Capabilities

- `power-profile`: Declare and reconcile the host's selected Tuned profile through Ansible.

## Impact

The exception is limited to selecting `balanced` on this host with
`tuned-adm`. The command persists the selection through Tuned. The owning
Ansible playbook and variables carry the desired state for replacement hosts
and future deployments, using the existing Tuned installation. This change
does not execute an Ansible deployment.
