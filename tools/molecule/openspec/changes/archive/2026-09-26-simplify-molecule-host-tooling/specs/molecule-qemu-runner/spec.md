## MODIFIED Requirements

### Requirement: Explicit execution prerequisites

Documentation SHALL identify supported platforms, tools, guest inputs,
accelerator selection, network dependencies, and execution/cache behavior.
QEMU and qemu-img SHALL be provided by the host installation together with
its compatible firmware and runtime modules. The runner SHALL check those
executables before convergence and record their versions and binary hashes.
Missing prerequisites SHALL produce an actionable failure without persistent
host reconfiguration. Host QEMU runtime dependencies and mutable guest package
services SHALL not be represented as hermetic inputs.

#### Scenario: Use host QEMU

- **WHEN** the documented host QEMU prerequisites are present
- **THEN** the guest runs using those executables and their installed firmware
- **AND** retained results identify the host executable versions and hashes

#### Scenario: Missing host QEMU

- **WHEN** a required host QEMU executable is unavailable
- **THEN** preflight fails with an actionable host-prerequisite diagnostic
- **AND** no guest starts and temporary runtime state is removed

#### Scenario: Network-dependent package setup

- **WHEN** a scenario requires guest packages from a mutable repository
- **THEN** its execution configuration declares the network requirement
- **AND** results record installed versions and identify package-service
  failures separately from role assertion failures
