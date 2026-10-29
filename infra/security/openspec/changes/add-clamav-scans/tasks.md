Sections 1–5 record the original implementation and historical acceptance.
They do not establish validation of the simplified candidate in section 6.

## 1. Original controller implementation

- [x] 1.1 Preserve and remove the superseded custom helper and synthetic fixtures; implement native tasks/templates with private signature verification and a bounded transient unit, validated by package/syntax checks and inspected command/configuration artifacts.

## 2. Packaged workflow

- [x] 2.1 Add separate fleet and XOA-default appliance targets sharing controller tasks; build both packages and run local syntax checks without authentication.
- [x] 2.2 Document variables, limit handling, scan policy, cleanup and evidence; inspect packaged commands, templates, and file mappings against documentation.

## 3. Handoff

- [x] 3.1 Record validation evidence and remaining live prerequisites; run OpenSpec validation and diff checks, then hand source to the parent without committing or pushing.

## 4. Parent runtime acceptance

- [x] 4.1 Verify the historical native candidate's clean, EICAR, and oversized-file behavior with retained receipts and confirmed controller cleanup; see validation.md.
- [ ] 4.2 Complete authorized production host scans and final delivery validation; retain per-host outcomes without claiming complete coverage or host cleanliness. Earlier remote canary results remain historical evidence.

## 5. Historical approved SSH maintenance

- [x] 5.1 Package a limited serial SSH hardening play through the existing shared role, retaining original identity, exact policy checks, reload-only activation, rollback, and private boot/operation receipts.
- [x] 5.2 Build and syntax-check the package and inspect its source and required native-role facts.
- [x] 5.3 Complete parent-owned fleet policy apply and acceptance: all nine validated configurations, unchanged root policy/CA/boot identities, fresh signed logins, native TCP denials, and owned-master cleanup are verified by rollout-result.json; the secure-runner Unix-forwarding canary passed. Preserve the ingress pre-play failure/retry and preview token-revocation cleanup failure. Full scans and final delivery remain pending under 4.2.

## 6. Simplification requested after rollout

- [x] 6.1 Retire the completed one-time SSH maintenance playbook, targets, and helper; keep the forwarding default in the existing SSH hardening role and document normal owner deployment controls with upstream restart-handler semantics.
- [x] 6.2 Consolidate duplicated ClamAV lifecycle playbooks and remove redundant controller checks using native Ansible and ClamAV behavior, preserving signature checks, private evidence, RAM/disk minimums, bounded execution, and cleanup.
- [x] 6.3 Verify the simplified packaged workflows and representative failure paths, preserving historical live evidence and the running scans' captured runtime. Delivery receipts own publication state; production outcomes remain under 4.2.

## 7. Explicit controller resource overrides

- [x] 7.1 Add validated thread, memory, and background-scheduling overrides with unchanged defaults; document native mappings and retain selected settings in lifecycle evidence.
- [x] 7.2 Validate packaged defaults and explicit overrides with native E2E execution; verify controller headroom and the separate eight-worker, 16 GiB normal-priority ingress run without altering captured running runtimes.
- [x] 7.3 Retain the priority ingress scan outcome and cleanup evidence; native termination with findings and errors does not establish complete coverage.

## Native simplification

- [x] Use systemd RuntimeDirectory with mode 0700 for the daemon socket; remove PID/cwd aliases and retain evidence under ignored out.
- [x] Verify successful and failed execution, private retention, owned cleanup, and the packaged candidate.
