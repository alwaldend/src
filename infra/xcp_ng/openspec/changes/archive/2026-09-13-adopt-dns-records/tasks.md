## 1. Prepare adoption

- [x] 1.1 Verify scoped DNS access for the existing `src_infra_xcp_ng` AppRole through the authenticated owner workflow; confirm successful backend/provider access and unchanged owner authentication source.
- [x] 1.2 Record the shared writer-audit coverage and recovery prerequisites, exact provider matches, and enabled source hashes; reuse the coordinator-authorized unchanged-declaration lint evidence and record the remaining delivery source gates.

## 2. Adopt and verify

- [x] 2.1 Review the declarative imports through `//infra/xcp_ng/tf:dns.plan` and `dns.show`; verify four exact provider IDs, all DNS actions no-op, and no non-DNS managed actions.
- [x] 2.2 Apply the reviewed saved plan through `//infra/xcp_ng/tf:dns.apply`; verify four imports with zero resource changes, no remaining imports or DNS changes, complete inventory and non-DNS saved-state equality, and four matching DNS answers.
- [x] 2.3 Record sanitized adoption receipts and enabled source defaults, synchronize both modified requirements through archive, preserve the original scenario names and link targets, and pass strict change and synchronized-spec validation.
