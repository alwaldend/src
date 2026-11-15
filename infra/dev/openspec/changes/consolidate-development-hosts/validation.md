# Validation

- Build and syntax-check the packaged Ansible deployment; inspect merged inventory to confirm host-specific storage and services remain separated.
- Render the ingress and backend routes with non-secret fixtures, retaining mTLS, matching TLS names and loopback listeners.
- Run Terraform formatting, DNS, OpenSpec, repository quality, documentation E2E and semantic lint checks.
- Inspect the shared configuration and deployment policy: one DNS root, shared identity and backend, one required OpenHands read, DNS provider reads, host-limited SSH, service-limited PKI and no legacy state access.
- Verify both renamed inventory groups and that limiting either host selects only its own play. Validate the packaged Terraform source with backend initialization disabled.

No live plan, Vault apply, DNS change, state transfer or host deployment is part of source validation. The old backends remain intact pending separately authorized adoption into the shared backend. The user confirmed the old Proxmox runner is absent and restricted retirement to IaC removal. Detailed local command receipts stay in ignored task scratch.
