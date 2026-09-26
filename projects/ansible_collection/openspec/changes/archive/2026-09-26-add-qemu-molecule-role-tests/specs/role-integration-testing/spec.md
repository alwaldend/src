## Purpose

Provide repeatable behavioral tests for the collection's deployment-used
`host`, `traefik`, and `forgejo` roles on disposable local QEMU machines.

The scenarios consume the `molecule-qemu-runner` contract owned by
`tools/molecule`. Its supported definition is in the linked
[runner specification](../../../../../../../../tools/molecule/openspec/specs/molecule-qemu-runner/spec.md).
Generic VM lifecycle, environment, cleanup, and evidence guarantees are owned
there and are not independently specified by this consumer.

## ADDED Requirements

### Requirement: Independently runnable role scenarios

The collection SHALL expose a Bazel test for each of `host`, `traefik`, and
`forgejo`, and a suite containing those tests. Each test SHALL apply the
packaged role implementation with documented fixture inputs. Prerequisite
setup SHALL remain separate from the behavior under test.

#### Scenario: Execute one role

- **WHEN** a developer selects one role's test
- **THEN** that scenario executes without another test having run first
- **AND** failure in convergence or verification produces a failing result

#### Scenario: Execute the initial suite

- **WHEN** a developer selects the role integration suite
- **THEN** it selects the three role tests
- **AND** it does not execute a production deployment playbook

### Requirement: Meaningful convergence and idempotence

Each scenario SHALL verify fresh convergence and an unchanged second
application. Within its documented coverage, the second application SHALL
report no changes and SHALL cause no unnecessary service restart or reload.
Tests MUST NOT conceal role changes by overriding change reporting or
disabling idempotence. Any external integration excluded from the fixture
SHALL be explicitly named in its coverage report.

#### Scenario: Unchanged service configuration

- **WHEN** a service role is applied twice with identical fixture inputs
- **THEN** the second application reports zero changes
- **AND** service lifecycle observations show no unnecessary restart or reload

#### Scenario: Intentional configuration update

- **WHEN** the Traefik or Forgejo scenario changes a supported configuration
  input and reapplies the role
- **THEN** the running service exhibits the updated behavior
- **AND** a further unchanged application is idempotent

### Requirement: Host role behavioral coverage

The `host` scenario SHALL exercise its composition of host-configuration
roles with nonempty representative user, firewall, and storage inputs.
Verification SHALL cover administrator access, configured SSH and OS
hardening, firewall behavior, installed test trust material, timezone, and a
persistent filesystem mount backed only by a disposable virtual data disk.

#### Scenario: Fresh host remains manageable

- **WHEN** `host` converges on a fresh guest
- **THEN** the configured test administrator can reconnect over SSH and use
  the declared privilege escalation
- **AND** allowed test traffic succeeds while explicitly denied traffic fails
- **AND** the expected trust, timezone, user, and hardening state is present

#### Scenario: Host storage survives a reboot

- **WHEN** the test writes a marker to the managed mount and reboots the guest
- **THEN** SSH becomes available again within a bounded deadline
- **AND** the configured mount and marker are present after reboot

### Requirement: Traefik role behavioral coverage

The `traefik` scenario SHALL verify installation, service enablement,
HTTP routing, TLS using a test CA, and supported configuration updates.
It SHALL use a disposable backend and SHALL disable production certificate
issuance through documented fixture inputs.

#### Scenario: HTTPS route works

- **WHEN** the role converges with the test listener, certificate, and route
- **THEN** a client trusting the test CA reaches the expected backend over TLS
- **AND** a request for an unconfigured route does not reach that backend

#### Scenario: Route changes become visible

- **WHEN** the fixture changes a configured route and reapplies the role
- **THEN** requests observe the new routing behavior
- **AND** the previous route no longer exposes the backend

### Requirement: Forgejo role behavioral coverage

The `forgejo` scenario SHALL run a real Forgejo service with a disposable
local database and test account. It SHALL verify service enablement, HTTP
access, Git repository writes and reads, configuration updates, and data
preservation across a service restart.

#### Scenario: Repository round trip

- **WHEN** a test account creates a repository and pushes a commit
- **THEN** a fresh clone contains the expected commit and file content
- **AND** the service's database and repository storage reside in the guest

#### Scenario: Restart preserves application state

- **WHEN** Forgejo restarts after a successful repository round trip
- **THEN** the test account, repository, and commit remain accessible
- **AND** verification does not recreate the application data to pass

### Requirement: Explicit role coverage exclusions

Each role scenario SHALL document external integrations excluded from its
fixture and SHALL report its tested inputs and assertions in the runner's
retained results. Role tests SHALL not be presented as deployment validation.

#### Scenario: Review external integration coverage

- **WHEN** a reviewer inspects scenario documentation and results
- **THEN** external SSH certificate issuance, Vault ACME integration, OIDC,
  external database integration, and full deployment behavior are explicitly
  identified as outside the initial coverage
