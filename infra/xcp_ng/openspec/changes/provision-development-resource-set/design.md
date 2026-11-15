## Context

Secure1 is being recreated under infra/dev at the user's explicit request.
The existing XO Terraform root owns per-component resource sets and bindings.

## Goals / Non-Goals

Provide the development owner's allocation through the existing OIDC workflow.
Host configuration and replacement remain owned by infra/dev.

## Decisions

Proceed with one resource_set_inventory entry. This follows existing owner
provisioning without sharing the infrastructure administrator token. A VM-state
migration was considered and explicitly declined by the user.

## Risks / Trade-offs

The first OIDC login must precede subject binding. Creation requires the
assigned template, storage and network to resolve unambiguously.

XO currently serves a self-signed certificate with no common name or SAN
(observed 2026-11-15). Scoped OIDC bootstrap fails before contacting Vault.
Use the existing certificate playbook, skipping the separately tagged OIDC
installation. The controller uses Ansible's native SSH askpass instead of an
unavailable sshpass executable. Inspect the appliance first; certificate
reconciliation requires separate exact-scope authorization under AGENTS.md.

The user explicitly approved the XO certificate repair after inspecting the
prepared scope. Keep running guest VMs up and restart only XO management.

The first certificate reconciliation was a no-op because the managed fullchain
file itself contained the same hostname-less self-signed certificate. Replace
the file-existence guard with chain and hostname verification, then reissue only
missing or invalid certificates and notify the management restart on issuance.

The installed XO code defaults autoCert to true and the inspected HTTPS
listener left it unset. The previous playbook only asserted a boolean claiming
it was disabled. Explicitly write autoCert = false after the inspected TLS key
setting, validate the TOML before installation, and notify the XO restart.
Remove the unused assertion variable. Reissuance alone was insufficient.
