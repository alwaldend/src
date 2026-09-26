# molecule-qemu-runner Specification

## Purpose

Provide a shared Bazel interface for running Molecule scenarios on disposable
local QEMU guests with isolated inputs, bounded cleanup, and useful evidence.

## Requirements

### Requirement: Declared candidate and test result

The runner SHALL execute the declared Molecule scenario using declared Ansible
executables, collections, package mappings, guest image, and fixture inputs.
It SHALL preserve executable runtime files and SHALL return failure when a
scenario phase fails. Its documented integration-test policy SHALL execute
an explicitly requested run rather than reuse a previous successful result.

#### Scenario: Run a packaged scenario

- **WHEN** a consumer invokes its Bazel scenario target
- **THEN** the packaged candidate and declared tools are used
- **AND** convergence or verification failure fails the target

#### Scenario: Repeat an explicit integration run

- **WHEN** the consumer invokes the same integration target again
- **THEN** the scenario actually executes using fresh runtime state

### Requirement: Disposable local VM lifecycle

Every scenario SHALL use local QEMU with a pinned guest image, fresh writable
disks, and run-specific connection information. It SHALL leave the base
image unchanged and MUST NOT attach physical host disks or use a production
inventory. Supported execution requirements SHALL be documented and checked
before convergence.

#### Scenario: Independent runs

- **WHEN** two scenarios run concurrently or the same scenario is repeated
- **THEN** their writable disks, ports, control endpoints, and credentials
  are distinct
- **AND** neither run can reuse the other run's guest state

#### Scenario: Unavailable prerequisite

- **WHEN** required tooling, guest inputs, or the selected accelerator are
  unavailable
- **THEN** the test fails with an actionable prerequisite error
- **AND** it does not report a successful or silently skipped role test

### Requirement: Cleanup and retained evidence

The test SHALL clean up its owned VM processes and temporary runtime state
after success, failure, and handled cancellation. It SHALL retain sanitized
phase results, guest diagnostics, input identities, and verification outcomes
as Bazel test artifacts. Cleanup failures SHALL fail the test. Abrupt process
death SHALL have a documented recovery procedure scoped to the affected run.

#### Scenario: Verification fails

- **WHEN** a behavioral assertion fails after the guest starts
- **THEN** the failure and relevant diagnostics are retained
- **AND** the VM is terminated and its disposable runtime state is removed

#### Scenario: Interrupted run

- **WHEN** the runner receives a supported cancellation signal
- **THEN** it performs bounded cleanup and returns a non-success result
- **AND** its evidence records the interrupted phase and cleanup outcome

### Requirement: Credential and inventory isolation

The runner SHALL construct a test-specific environment and inventory without
inheriting production infrastructure credentials or inventories. Generated
test credentials SHALL remain temporary and excluded from retained artifacts.

#### Scenario: Developer has production credentials configured

- **WHEN** a developer launches a test with infrastructure credentials in
  their environment
- **THEN** only the declared test inventory and test credentials reach the
  scenario and guest
- **AND** retained artifacts contain no private keys, passwords, or tokens

### Requirement: Explicit execution prerequisites

Documentation SHALL identify supported platforms, tools, guest inputs,
accelerator selection, network dependencies, and execution/cache behavior.
Missing prerequisites SHALL produce an actionable failure without persistent
host reconfiguration. Mutable package services SHALL not be represented as
hermetic inputs.

#### Scenario: Network-dependent package setup

- **WHEN** a scenario requires guest packages from a mutable repository
- **THEN** its execution configuration declares the network requirement
- **AND** results record installed versions and identify package-service
  failures separately from role assertion failures
