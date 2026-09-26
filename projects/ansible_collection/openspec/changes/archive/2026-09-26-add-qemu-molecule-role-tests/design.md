## Context

See [proposal.md](proposal.md) for scope. Source inspection at
`b84d302a7a5b15b48a6a79d0795271a394e60d51` on 2026-09-26 found that `host`
composes host configuration and invokes Vault for SSH certificate signing.
Traefik and Forgejo request reloads on unchanged applications. Those are
source-level reasons to expect failures, not observed VM-test results.

The shared runner is defined by the linked
[tools-owned change](../../../../../../tools/molecule/openspec/changes/archive/2026-09-26-add-qemu-molecule-runner/proposal.md).
Its design owns VM provisioning, dependencies, package materialization,
credential isolation, cleanup, execution policy, and retained evidence.
This design owns collection inputs, assertions, and any narrow role fixes.

## Goals / Non-Goals

**Goals:** Exercise real packaged roles with representative inputs and
observable behavior, documenting fixture exclusions.

**Non-Goals:** Full deployment topologies, production identity integrations,
other independent role scenarios, or a second runner contract.

## Decisions

### 1. Collection-owned scenarios

Place scenarios under `extensions/molecule/{host,traefik,forgejo}/` and shared
role fixtures under `extensions/molecule/shared/`. Expose `:molecule_test`
single-test suites from those role packages and a
`//projects/ansible_collection:molecule_test` suite. Each scenario consumes
the shared runner's declared package and guest inputs. Keep fixtures outside
production role source globs and collection package mappings.

Service scenarios prepare only their prerequisites; they need not execute
`host` first. This makes failures attributable while the separate `host`
scenario verifies the host composition. The alternative of a single combined
scenario would make role regressions harder to isolate.

### 2. Representative fixtures and honest exclusions

| Scenario  | Fixture and acceptance focus                                                                                                                                                                                                                                                      |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `host`    | One administrator and regular user, generated SSH keys, explicit allow/deny ports, test trust material, enabled SSH/OS hardening, and LVM/filesystem/mount inputs on a blank virtual disk. Reconnect after convergence and reboot; verify access controls and a persisted marker. |
| `traefik` | Packaged binary, real systemd unit, local HTTP backend, generated CA/server certificate, static and dynamic configuration, and `traefik_disable_eab: true`. Verify HTTPS, unmatched routes, route updates, and unchanged lifecycle behavior.                                      |
| `forgejo` | Packaged binary, real systemd unit, SQLite database, test account, and test Git repository. Verify HTTP access, push/clone content, a supported application configuration update, and preservation after restart.                                                                 |

For `host`, prepare locally signed SSH host certificates at the paths the
existing hardening tasks consume, and skip only the `ssh_sign_key` tag.
Record external certificate issuance as excluded. Continue exercising SSH
configuration, host certificate loading, access, and hardening; do not use
`ssh_sign_ignore_errors` to turn a failure into a pass. Apply test-only trust
file mappings where the role uses fixed source filenames, preserving the
production mappings. The fixture owns its keys and certificates.

Use nonempty inputs for host storage and firewall behavior so a passing test
cannot simply skip them. Test CA installation must not contact production
Vault. Use only temporary generated credentials for Forgejo and TLS; suppress
secret-bearing command output and omit keys, tokens, and passwords from
artifacts. SQLite keeps an additional database role out of scope. External
databases, OIDC, public ACME, and the deployment's Xen disk assertions remain
outside coverage.

The alternative of running a Vault fixture would cover certificate issuance
but adds a separate service and renewal lifecycle. That can be a later
scenario; it is not required to establish these three role tests.

### 3. Idempotence failures lead to narrow role corrections

Write the behavioral scenarios before changing role behavior. Expect the
first unchanged run to expose the unconditional Traefik and Forgejo reloads.
Inspect the pinned binaries' supported reload semantics, then notify handlers
only for actual relevant changes. Account for binary, unit, static config,
and dynamic config updates without triggering duplicate restarts.

Validate both unchanged and changed configurations through the running
service. Use Ansible change results together with systemd lifecycle timestamps
or journal observations; process identity alone cannot detect an in-process
reload. Do not add blanket idempotence exclusions or fake change reporting.
Any other fix must be a demonstrated blocker in these scenarios and preserve
existing deployment defaults; broader refactoring is separate work.

## Risks / Trade-offs

- Fedora or pinned hardening incompatibility -> validate the runner's guest
  against the role prerequisites before completing the host fixture.
- Firewall or SSH hardening can sever management access -> verify reconnect
  using the declared management route and retain the runner's diagnostics.
- Fixture exclusions reduce coverage -> publish them with scenario results;
  never describe role acceptance as production deployment validation.
- Existing defects may emerge -> make only demonstrated in-scope corrections
  and verify unchanged deployment defaults.

## Migration Plan

Complete the linked runner smoke acceptance, then the host, Traefik, and
Forgejo scenarios in that order. Validate the three role entry points and suite,
inspect retained outputs, and verify test fixtures do not enter production
payloads. No live deployment or host configuration is part of this rollout.
Rollback removes the scenarios and separately reverts any incompatible role
correction. Update cross-change links when either change is archived.

## Implementation evidence (2026-09-26)

Before role corrections, real VM runs of Traefik 3.7.1 and Forgejo 15.0.3
completed initial convergence but failed their unchanged-application lifecycle
assertion. In both cases the `Reload the service` task reported `changed: true`
on the second application. Structured task results and package inventories are
retained locally under `out/molecule-role-tests/baseline/{traefik,forgejo}`.

The correction replaces unconditional lifecycle tasks with one notified restart
handler per service, covering binary, unit, and configuration changes. Traefik
also notifies for its dynamic configuration and client CA. A just-started
service skips a duplicate restart. This supports declared configurations even
when file watching is disabled and coalesces multiple changed inputs. The
trade-off is a restart for dynamic changes that a watching file provider could
otherwise reload automatically. Unchanged applications perform neither action.
[Forgejo's v15 configuration reference](https://forgejo.org/docs/v15.0/admin/config-cheat-sheet/)
requires a full restart for configuration changes. Traefik separates
[startup and routing configuration](https://doc.traefik.io/traefik/v3.7/getting-started/configuration-overview/);
the tests must verify actual updated routes after application. The subsequent behavioral runs below establish acceptance of these corrections.

The host baseline reached its second application with changes in firewalld,
SSH configuration, and `su` permissions after the initial role's final package
upgrade. The installed versions included OpenSSH 10.2p1-14.fc44, util-linux
2.41.5-1.fc44, and firewalld 2.4.4-1.fc44. The scoped correction moves the same
package-upgrade role before host configuration, so later package replacement
does not undo the already-applied policy. No update or hardening task is
excluded. The initial fixture also reused Fedora's existing `operator` account
(UID 11); it now uses a new `molecule_admin` account to exercise real user
creation without colliding with a system account. Subsequent host idempotence and reboot acceptance is recorded below.

The corrected host rerun narrowed idempotence changes to firewalld, the SSH
configuration, and sysctls. UFW's pinned module runs `ufw enable` even when
already enabled; devsec's UFW sysctl protection applies only to Debian. The
repository OS role now disables UFW's competing sysctl reload on Red Hat
systems with UFW installed. Host composition configures the firewall after SSH
package dependencies, so Fedora's fail2ban dependency installation cannot
re-enable firewalld after its final configuration. SSH certificate paths are
sorted before rendering; filesystem enumeration order is not configuration.
The fixture now verifies an actual SSH certificate handshake against its test
CA, in addition to administrator access and effective configuration.

Forgejo startup generates and saves a missing OAuth2 JWT secret even when the
fixture uses basic authentication. The fixture supplies a generated per-run
secret, following [the v15 configuration contract](https://forgejo.org/docs/v15.0/admin/config-cheat-sheet/),
so the managed file remains stable without suppressing role change reporting.

Candidate `3ee3787109a858b26f6cbfe2dc2d7a49396ba6a1` passed the complete
Traefik scenario, including configuration update with boot enablement drift,
repeat idempotence, and restart persistence. Host passed idempotence,
administrator privilege access, and a real test-CA SSH certificate handshake.
Its trust assertion used a removed Fedora 44 bundle path; the fixture now uses
OpenSSL's default system trust lookup, matching
[Fedora's trust-store change](https://fedoraproject.org/wiki/Changes/droppingOfCertPemFile).
Forgejo still changed its managed configuration at that candidate; private diff
diagnosis followed, and no secret values are retained in its diagnostic summary.

Private diff inspection isolated the remaining Forgejo change to an added
`WORK_PATH` setting. Live comparisons confirmed every generated secret matched
the fixture, including the valid 32-byte JWT secret. The fixture now declares
`WORK_PATH` from the role's `forgejo_work_dir`, so startup does not rewrite the
managed file. The diagnostic summaries retain only field names, lengths, and
match booleans; raw configuration and secrets were not retained.


## Accepted result

Candidate `1bfe561c394e52d514b2bc756154ff24cc07a320` passed all ten phases for
each role, including both unchanged applications, the intentional update, and
repeat verification. The host passed fresh administrator SSH/sudo access,
CA-authenticated SSH, test trust, effective hardening, allowed/denied traffic,
UTC, the additional user, and two reboots preserving its mounted filesystem
and original marker. Traefik passed trusted HTTPS, unmatched and removed route
rejection, updated routing, and restart persistence. Forgejo passed the real
account/API/Git round trip and fresh clones of the original commit before and
after restarts, including after its application configuration update.

All three reported zero changed tasks in both idempotence phases, successful
cleanup, and removed private runtime directories. Their 5,299 packaged inputs
contained no scenario fixture files and included the real service binaries.
The convenience targets and aggregate resolve to exactly these three tests.
The aggregate and individual entry points were selected in the accepted Bazel
invocation; their underlying scenarios also ran independently during diagnosis.
All used the restored pinned Fedora image digest recorded by the runner owner.

Affected semantic lint, collection packaging, and 27 repository quality and
OpenSpec checks passed. Retained outputs live under
`out/molecule-role-tests/candidate-1bfe561`; the validation receipt is
`out/molecule-role-tests/fixture-corrections-prepare.json.validation.json`.
[Implementation PR #111](https://github.com/alwaldend/src/pull/111) was published
and verified for that candidate. Later archive/documentation updates preserve
all tested runtime inputs and rerun the mandatory publication gates.
