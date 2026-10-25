# Deployment hardening acceptance

Source baseline: `1c8cf9412d9e9c0f634306237d78664ba2450ad9`.
The isolated checks below were observed on 2026-10-25 without changing live
hosts. The subsequent authorized
[local rollout](../../../../../infra/download/openspec/changes/add-static-hosting/evidence.md#authorized-local-deployment-2026-10-26)
records the 2026-10-26 deployment, two reboots, and the corrections it exposed.
Publication receipts bind the final candidate.

## Isolated behavior

A disposable rootless Fedora container ran the actual collection roles, pinned
external collections from the download Ansible runfiles, and the repository's
Traefik 3.7.1 binary. Runtime networking was disabled. Container-local systemd
and network-administration capabilities did not grant host privileges.

| Case | Observation |
| --- | --- |
| Unchanged proxy | Zero changed tasks; the process ID remained identical across the unchanged run and dynamic-route update. |
| Dynamic configuration | HTTP redirect changed to the requested destination without replacing the process. |
| Configuration permissions | Running as the service account could not write the configuration directory. Static/TLS inputs and their parent are administrator-owned. |
| Invalid TOML | Play failed before installing the candidate, with zero changed tasks and censored diagnostics. |
| Rejected startup | An invalid listener address failed startup; rescue restored the old configuration and HTTP service, then reported failure. |
| Binary-only failure | At the host_bot-style custom executable path, replacing the binary with an executable that exits immediately triggered restoration and successful service recovery. |
| ACME enrollment | A private public-data fixture cache was reused with an intentionally nonexistent Vault executable; an unchanged repeat reported zero changes. No issuer request was made. |
| SSH certificate renewal | Replacing only the certificate activated serial 2, observed over a new SSH connection. The daemon configuration remained unchanged. |
| Fedora reboot detection | Actual DNF detection succeeded. Injected DNF exits 0 and 1 produced the corresponding fact; exit 2 failed rather than reporting no reboot. No reboot task ran. |
| Explicit firewall retirement | The standard `community.general.ufw` module removed inventory-listed incoming/routed rules and inherited SSH/mDNS profile allowances, including multicast destinations, from UFW's actual rule files. Desired TCP 22/80/443, an unlisted rule with an apostrophe in its comment, and an unlisted rule with a legacy `# Ansible` comment remained. A repeat reported zero changes. |
| Missing content device | The systemd generator produced a `local-fs.target.wants` mount and a ten-second device timeout. Nginx required the mount and stayed inactive when it was absent; SSH and Traefik stayed active. |

## Evidence and reproduction

Task-local inputs, command logs, source digests, result summaries, and the
pinned binary digest are retained in `out/download-fixes/acceptance.json` and
its adjacent fixtures. The disposable container is named
`download-role-acceptance-20261025`; its image is
`localhost/download-role-acceptance:20261025`. Example repeat:

```sh
podman exec --env ANSIBLE_COLLECTIONS_PATH=/fixture/collections \
  download-role-acceptance-20261025 ansible-playbook -i localhost, \
  /fixture/traefik-acceptance.yaml -e @/fixture/custom-vars.json
```

Expected failure logs are `traefik-syntax.log`, `traefik-recovery.log`,
`traefik-binary-recovery.log`, `reboot-error.log`, and
`nginx-missing-mount.log`. Their nonzero result is required by the scenario.
The final valid proxy run serves the configured route and reports zero changes.

The explicit-deletion revision supersedes the earlier custom-pruner fixture.
Its inputs and logs are in `out/download-explicit-firewall/`, and its disposable
container is `download-explicit-firewall-acceptance`, using the same image.
The source collection and download inventory are mounted read-only. `seed.py`
installs an offline UFW wrapper and seeds the real rule files; `acceptance.yaml`
runs the production `firewall` role's `rules.yaml` tasks and asserts the resulting
policy. `retirement.json` extends the download retirement list with an incoming
port and a routed rule. Repeat after starting the existing container:

```sh
podman exec --env ANSIBLE_COLLECTIONS_PATH=/fixture/collections \
  download-explicit-firewall-acceptance ansible-playbook -i localhost, \
  /acceptance/acceptance.yaml -e @/acceptance/retirement.json
```

`first.log` records successful retirement; `second.log` records zero changes.
Removing an allow declaration without a corresponding retirement intentionally
leaves its installed rule in place, regardless of its comment.

## Limits and rollout

These limitations describe the original isolated checks; the linked local
rollout supplies later live evidence for firewall activation and reboot
persistence, initial enrollment-cache creation, and unchanged service reruns.
It does not test disk-loss recovery or fresh TLS certificate issuance.

The host kernel lacks legacy `ip_tables`, so full UFW activation could not run.
The offline check disabled UFW's kernel capability checks, substituted its
status query with an inactive result, and kept its firewall disabled. It
exercised the real rule-file operations through `community.general.ufw` and the
production rule tasks. It does not establish packet filtering or complete
firewall-role convergence on a live host. Confirm those during an authorized
rollout.

Container device-dependency checks are not a VM reboot or disk-loss recovery
exercise. Certificate issuance and ACME re-enrollment against live Vault remain
unverified. Syntax checks do not prove dynamic route semantics or external
issuer reachability. Existing download live-acceptance tasks remain open.

## Session corrections

- `acceptance-runtime-input`: an initial copy selected Bazel's launcher instead
  of its ELF payload; the runfiles manifest supplied the correct pinned binary.
- `acceptance-container-capabilities`: systemd sandboxing needed SYS_ADMIN
  inside the rootless container's user namespace. No host configuration changed.
- `acceptance-firewall-kernel`: missing legacy netfilter support prevented full
  firewall execution; the result was narrowed explicitly to rule-file behavior.
- `activation-handler-rescue`: the first handler-based recovery failed to run
  rescue. A regular task transaction passed the startup and binary failure cases.

These are fixture-tested observations, not changes to shared agent policy or
claims that the separate Molecule/QEMU integration is implemented.

## Yandex firewall persistence correction

The [authorized Yandex rollout](../../../../../infra/download/openspec/changes/add-static-hosting/evidence.md#yandex-zone-change-2026-10-26)
found that firewalld's unit was absent when the firewall role first ran. Its
conditional disable task was skipped. A later host dependency installed and
enabled firewalld, which conflicts with UFW. After the package-update reboot,
UFW was enabled but inactive and filtering was disabled.

The role now masks firewalld regardless of whether its package is installed,
while retaining the existing stop/disable ordering for an installed service.
The pinned `ansible.builtin.systemd_service` explicitly supports masking absent
units. No new dependency or imperative host repair was introduced.

The live correction run changed only firewalld disablement, its mask, and UFW
activation. An owning-playbook verification reboot passed. Before another
configuration run, UFW was active/enabled and firewalld was inactive/masked.
A subsequent firewall/Nginx/Traefik run reported 33 successful tasks and zero
changes, with unchanged Nginx/Traefik process IDs. Ten public HTTPS checks
passed before and after reboot and after convergence. The download owner's
record contains commands, counters, retained artifacts, and remaining limits.
The corrected missing-package path was checked against the pinned module's
supported behavior; a second fresh VM was not provisioned solely for that case.

## Unified updater reconciliation, 2026-11-01

Trunk commit b8b60ae7 consolidates immediate upgrades, reboot detection and
native automatic updates in update_all_packages. The rebase preserves that
role unchanged, including OSTree support and explicit immediate reboot opt-in.
The prior standalone automatic_updates role and its second host import are
removed. Earlier dated acceptance above remains historical; final delivery
checks validate the integrated source and packaged consumers. The certificate
reload handler is retained alongside upstream SSH hardening controls.
