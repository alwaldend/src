## Context

The seven findings and read-only local observations were recorded at
`1c8cf9412d9e9c0f634306237d78664ba2450ad9`. Shared roles serve other consumers,
including host_bot; changes must retain their ordinary defaults and packaging.

## Goals / Non-Goals

**Goals:** Fix the reviewed behavior in its existing owner, using already
available Ansible, Python, systemd, UFW, and pinned Traefik inputs.

**Non-Goals:** Live provisioning, reboot, firewall changes, certificate issuance,
new dependencies, backup infrastructure, or implementation of the separate
Molecule/QEMU runner proposal.

## Decisions

Verdict: proceed with source fixes and isolated behavioral checks.

- Use DNF's reboot hint on Red Hat-family hosts, accepting only exit codes 0
  and 1. Preserve the Debian marker. Expose the reboot policy and result so
  operators can separate detection from an approved maintenance reboot.
- Validate rendered TOML/YAML on the controller with its existing parsers.
  Keep content out of parse-error messages. Use a change-driven Traefik task transaction
  with startup verification and restoration of available backups on failure.
  A syntax check alone does not establish routing or issuer reachability.
  The activation task block owns rescue: the initial handler-block experiment
  did not execute rescue after a startup failure.
- Keep configuration, TLS inputs, binary ancestors, and the main proxy
  directory administrator-owned. Limit service writes to data and logs.
  Cache local ACME enrollment material privately to avoid generating a new
  binding and restarting the proxy on every unchanged run.
- Notify a role-owned SSH reload handler on certificate replacement. The
  upstream hardening handler alone does not observe renewed certificate bytes.
- Add `nofail` and a bounded device wait to the content mount. Give Nginx an
  explicit mount dependency; deduplication already has one. Dropping the boot
  dependency without constraining content services would be unsafe.
- Use `community.general.ufw` with `delete: true` for inventory-declared
  `firewall_delete_rules`, after adding desired allows. Download enumerates
  inherited SSH and mDNS profile rules, including multicast destinations.
  Removing an allow entry alone leaves the existing rule in place; retirement
  requires an explicit deletion entry retained for hosts that missed the change.
  Do not parse command output, infer ownership from comments, or reset UFW.
- Stop and disable an installed firewalld service before restoring UFW, and
  mask firewalld even when absent. Later dependency installation must not
  introduce a competing enabled service that prevents UFW startup at reboot.

## Failure cases and acceptance

These cases are recorded before implementation and guide isolated acceptance:

- DNF returns reboot-needed, no-reboot, or an unexpected error; Debian marker
  detection remains intact and disabled reboot policy does not hide the result.
- Unchanged Traefik rerun, executable/unit/static updates, dynamic changes,
  malformed configuration, startup-rejected configuration, and initial install.
- Enrollment material is reused without printing it; required service data
  remains writable while configuration and all of its writable ancestors do not.
- Signed SSH certificate replacement triggers a daemon reload with unchanged
  sshd configuration; invalid sshd configuration prevents activation.
- Generated mount dependencies allow OS boot while Nginx and maintenance
  require the data filesystem.
- Explicitly deleting a port and a routed rule, preserving unlisted rules
  (including a comment with an apostrophe), deleting generic and multicast
  profile rules, and unchanged convergence.

## Risks / Trade-offs

- Shared-role consumers may use customized paths -> inspect current consumers,
  package their inputs, and verify a representative custom binary location.
- Firewall rule revocation can disrupt access -> add desired rules first and
  limit deletion to the explicit rule specifications.
- Configuration syntax is not a full routing proof -> exercise live fixture
  HTTP behavior and verify activation/restoration separately.
- Containers cannot prove VM disk-loss boot recovery -> inspect generated
  systemd dependencies and retain live reboot/disk-failure testing as a rollout
  limitation, without claiming it passed.

## Migration Plan

Publish validated source only. A later explicitly authorized deployment adopts
protected directories, explicit firewall retirements, and the mount dependency.
Review its maintenance reboot and initial enrollment-cache migration separately.
Rollback uses the prior source and retained service configuration; no live
rollback operation is performed by this implementation task.

## Native automatic security updates

The later authorized rollout adds automatic_updates to host. Use the existing
package, ini_file, copy, and systemd_service modules. Select DNF5, DNF4, or APT
from gathered package-manager facts. Install the specifically approved native
package and use its existing systemd timer. DNF selects security advisories;
APT retains the distribution's packaged allowed origins. Disable automatic
reboots explicitly. Skip immutable OSTree and unsupported hosts visibly.

Keep one public enabled setting. Disabling stops the managed installation
timer and disables its update configuration, without removing packages.
Do not start a package-upgrade transaction merely to test the role.
Before implementation, acceptance cases are: first package installation,
minimal security settings with preserved vendor options, active enabled timer,
unchanged rerun, disable with and without a preinstalled package, and unsupported
or OSTree host selection. Live E2E covers the two Fedora download hosts; report
DNF4 and APT runtime coverage separately if those hosts are unavailable.

The 2026-10-26 scoped deployment passed on both Fedora download hosts. Each
first run reported five tasks, three changes, and zero failures; repeated runs
reported zero changes. Both dnf5-automatic.timer units are enabled and active,
with security-only updates enabled and automatic reboot disabled. Verification
did not trigger an update transaction. DNF4, APT, disable, and OSTree branches
remain configuration-reviewed rather than live-tested.
