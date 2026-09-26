## Context

See [proposal.md](proposal.md) for scope. Source inspection at
`b84d302a7a5b15b48a6a79d0795271a394e60d51` on 2026-09-26 found that `host`
composes host configuration and invokes Vault for SSH certificate signing.
Traefik and Forgejo request reloads on unchanged applications. Those are
source-level reasons to expect failures, not observed VM-test results.

The shared runner is defined by the linked
[tools-owned change](../../../../../tools/molecule/openspec/changes/add-qemu-molecule-runner/proposal.md).
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
aliases from those role packages and a
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
Forgejo scenarios in that order. Validate the three role aliases and suite,
inspect retained outputs, and verify test fixtures do not enter production
payloads. No live deployment or host configuration is part of this rollout.
Rollback removes the scenarios and separately reverts any incompatible role
correction. Update cross-change links when either change is archived.
