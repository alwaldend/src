## 1. Implementation

- [x] 1.1 Consolidate owner and maintained OpenSpec work; verify old source paths have no active consumers.
- [x] 1.2 Align DNS, routes, certificate scope and state partitions; verify offline Terraform and inventory outputs.
- [x] 1.3 Update consumers and documentation; run syntax, packaging, quality and semantic lint checks.
- [x] 1.4 Replace nested deployment orchestration with one standard Ansible invocation and a scoped deployment identity; verify packaged syntax, selected hosts and credential scope.

- [x] 1.5 Remove the legacy runner VM, DNS and inventory declarations; verify only secure1 and dev1 remain.
- [x] 1.6 Merge DNS into one root and shared identity, remove unused runner credentials and old inventory names, and document the pending state handoff; verify packaged configuration and host separation offline.

## 2. Authorized rollout

The user confirmed PR #305 merged and requested rebase and deployment. Dev1 reboot remains disabled. The later recreation request authorizes disposal of the old secure VM; other resources remain outside scope.

- [x] 2.1 Rebase onto the merged master and provision the scoped shared Vault identity through an interactively reviewed Terraform apply.
- [x] 2.2 Inspect both legacy states through temporary backend-only targets, reconcile DNS ownership and apply reviewed DNS changes without destructive recreation; remove temporary migration inputs after a no-change plan.
- [x] 2.3 Deploy secure1, dev1 and the changed ingress routes through packaged Ansible, verify relevant host and service state, and record sanitized results.
- [x] 2.4 Remove obsolete migration wiring and prepare the follow-up source and operational record for delivery.
- [x] 2.5 Scope Traefik account storage to its ACME directory, redeploy the proxy, and verify the new hostname directly with certificate validation enabled.

- [x] 2.6 Remove OpenHands deployment and secret access, deploy its service removal and verify it stays inactive and masked without a reboot.

## 3. Secure1 recreation

The user explicitly declined VM-state migration and requested recreation of the
secure VM. Its disk contents may be discarded; dev1 must not reboot.

- [x] 3.1 Define secure1 provisioning under infra/dev and provision its own scoped XO resource set.
- [x] 3.2 Interactively destroy only the old secure VM, create secure1, and deploy its host baseline with strict SSH verification.
- [x] 3.3 Confirm both old backends have no managed resources, remove both obsolete Vault identities and their remaining grants, and remove temporary destruction wiring.
- [x] 3.4 Verify final DNS and VM plans, host identity and service health; record the completed rollout for delivery.

- [x] 3.5 Recreate the canonical secure1 VM again at the user's request using 192.168.10.102, update its canonical DNS record and verify cloud-init, the guest address and SSH before host configuration.

- [x] 3.6 Remove the completed OpenHands cleanup block and obsolete tag documentation; verify the packaged playbook has no OpenHands tasks.

- [x] 3.7 Replace the secure1-specific host-key command with shared SSH-role reconciliation and verify packaged configuration.
