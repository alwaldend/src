# Validation

- Build and syntax-check the packaged Ansible deployment; inspect merged inventory to confirm host-specific storage and services remain separated.
- Render the ingress and backend routes with non-secret fixtures, retaining mTLS, matching TLS names and loopback listeners.
- Run Terraform formatting, DNS, OpenSpec, repository quality, documentation E2E and semantic lint checks.
- Inspect the shared configuration and deployment policy: one DNS root, shared identity and backend, no OpenHands secret read, DNS provider reads, host-limited SSH, service-limited PKI and no legacy state access.
- Verify both renamed inventory groups and that limiting either host selects only its own play. Validate the packaged Terraform source with backend initialization disabled.

Live deployment was separately authorized after merge. See the [deployment record](deployment.md)
for DNS adoption, host and ingress results, and the OpenHands removal. Source validation and live verification are distinct; detailed
local command receipts remain in ignored task scratch.
