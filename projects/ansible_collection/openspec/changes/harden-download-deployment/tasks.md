## 1. Acceptance preparation

- [x] 1.1 Record all seven findings, alternatives, and failure cases before implementation; verify the owner-local proposal, design, and requirements exist.
- [x] 1.2 Prepare isolated role acceptance inputs before implementation and retain repeatable commands and result artifacts under task scratch.

## 2. Implementation and focused verification

- [x] 2.1 Correct reboot detection and policy; exercise Fedora detection and error handling without rebooting a live host.
- [x] 2.2 Protect Traefik configuration and implement validated, change-driven activation and recovery; verify unchanged runs, permissions, rejected configuration, and working HTTP after recovery.
- [x] 2.3 Reload renewed SSH host certificates; verify certificate replacement activates without an sshd configuration change.
- [x] 2.4 Separate content storage from OS boot and require it for Nginx; verify generated mount and service dependencies.
- [x] 2.5 Delete inventory-declared retired rules with the standard UFW module; verify explicit revocation, unlisted-rule preservation, and unchanged convergence.

## 3. Delivery

- [x] 3.1 Update owner and consumer documentation and evidence, preserving unverified live acceptance work; validate OpenSpec and packaging.
- [x] 3.2 Run formatting, affected semantic lint, repository quality, behavioral acceptance and exact-candidate review; commit, push, and publish the pull request.

## 4. Explicit deletion revision

- [x] 4.1 Remove the custom pruner, digest comments, and its import-check dependency; expose explicit rule deletions and enumerate download's inherited profile rules.
- [x] 4.2 Exercise standard-module deletion, quoted-comment preservation, routed and multicast rules, and unchanged convergence in isolation.
- [x] 4.3 Reconcile operator documentation and behavioral evidence for explicit deletion, retaining the live rollout limitations. Reapply the delivery checks in section 3 to the revised candidate.

## 5. Local rollout corrections

The authorized download deployment on 2026-10-26 exposed a missing-cache
failure: the enrollment block reevaluated its condition after updating
`traefik_eab_previous`, so its final write was skipped. The earlier isolated
fixture covered cache reuse, not first enrollment.

- [x] 5.1 Preserve the enrollment decision across all tasks in the block; verify a missing cache is created privately and a subsequent deployment reuses it without requesting another binding.
- [x] 5.2 Record the full local deployment and reboot observations through the download owner's rollout evidence; repeat package and affected-consumer validation before publication.
- [x] 5.3 Enable the UFW systemd service and restore its retained policy before updating rules; verify active filtering after a reboot and successful repeated deployment without adding or removing unrelated rules.

The first package-update reboot exposed `ufw.service` as disabled despite
`ENABLED=yes`; filtering was inactive and `ufw logging low` failed because
its chains were missing. Proceed with the owning systemd service, after
disabling the conflicting firewalld service. Preserve the existing policy
before changing rules. Moving logging alone would leave reboot persistence
unfixed; an imperative firewall reload would not establish desired state.

## 6. Firewall packages introduced during deployment

The first Yandex deployment installed UFW at 14:16:35 UTC while firewalld's
unit was absent, so the conditional disable task was skipped. A later host
dependency installed firewalld at 14:17:39 UTC. After the package-update reboot,
both conflicting services were enabled and UFW filtering was inactive.

Failure cases to cover before changing the role: a missing firewalld unit
that is installed later, an existing enabled firewalld service, retained UFW
rules after reboot, preserved SSH access, and an unchanged repeated role run.
Masking firewalld before enabling UFW prevents later package installation from
making it startable. Reordering only the disable task would leave the next
package transaction able to reintroduce the conflict.

- [x] 6.1 Mask firewalld even when its package is initially absent, retaining the existing stop/disable ordering for installed services.
- [x] 6.2 Apply the owning firewall role on Yandex, verify enabled active UFW and masked inactive firewalld across a verification reboot, and confirm unchanged service convergence and public HTTPS.
- [x] 6.3 Record the observed failure and live verification, validate shared-role consumers, and publish the final guarded deployment configuration.
## 7. Historical automatic security updates

- [x] 7.1 Add a reusable automatic_updates role to the shared host role using native distribution systemd timers and the user-approved OS packages.
- [x] 7.2 Verify Fedora/DNF5 configuration and enabled timers on both download hosts, repeat for convergence, and record other platform validation limits.

## 8. Reconcile the unified updater from trunk

- [x] 8.1 Keep trunk’s single update_all_packages owner, controls, backend detection and reboot opt-in; remove the superseded standalone updater and duplicate host import.
- [x] 8.2 Retain SSH certificate activation alongside trunk’s SSH hardening controls and reconcile current documentation and delta scope.

Section 7 records the earlier deployed implementation. Trunk now supplies the
canonical update role; its behavior and controls belong to that owner.
