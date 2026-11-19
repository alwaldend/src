## 1. Prepare project DNS ownership

- [x] 1.1 Create this linked owner change and verify strict OpenSpec validation and local links.
- [x] 1.2 Add the canonical `tf` root and source exports; verify it consumes the existing declaration with DNS record creation disabled by default.
- [x] 1.3 Verify the root-workspace operational wrapper, project AppRole/backend references, and normal `tf=main` credential injection in the packaged output.

## 2. Validate integration

- [x] 2.1 Run the owning offline package checks and inspect a representative packaged output against the candidate; record the result in the linked migration evidence.

Live adoption and the freeze of all older central deployment candidates follow
the [cutover runbook](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
and [shared migration tasks](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/tasks.md).
AppRole and credential deployment is a prerequisite for new root use; the
disabled DNS resource default does not remove that prerequisite.
