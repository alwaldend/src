## Purpose

Provide verified private controller scanning without loading signatures on
small hosts or changing the behavior of host maintenance workflows.

## ADDED Requirements

### Requirement: Separate packaged scan entry points

The workflow SHALL expose independent fleet and appliance ClamAV targets using
the existing owner inventory exports and operator authentication boundaries.
They SHALL invoke the shared streaming scan role without selecting maintenance.
They SHALL preserve server SSH policy and report denied forwarding as an
incomplete scan, without automatically relaxing authentication or forwarding
restrictions.

#### Scenario: Operator selects a limited host scan

- **WHEN** the operator selects a ClamAV target and includes localhost in a limit
- **THEN** controller preparation, only selected endpoint scans, and controller
  shutdown run, without package updates, reboots, deletion, or quarantine

### Requirement: Private verified bounded controller

The workflow SHALL use the pinned unprivileged engine, one private persistent
database cache, explicit runtime library paths, a private Unix socket, and
bounded engine resources. It SHALL inspect actual signature versions and
freshness after freshclam and require successful database loading.

#### Scenario: Freshclam reports success with stale signatures

- **WHEN** freshclam returns zero but the daily database exceeds the allowed age
- **THEN** preparation fails with private signature and lifecycle receipts and
  the remote scan role cannot use the controller

#### Scenario: A daemon load or startup operation fails

- **WHEN** signatures or the engine cannot initialize successfully
- **THEN** incomplete evidence is retained and owned processes are stopped

#### Scenario: Operator selects controller resources and scheduling

- **WHEN** the operator supplies controller thread and memory overrides or
  disables the shared background setting
- **THEN** preparation accepts integer threads from 1 through 32 (default 2),
  integer memory ceilings from 1 through 64 GiB (default 6), and a boolean
  background setting (default true), rejecting invalid values before execution
- **AND** native MaxQueue is four times MaxThreads, MemoryMax uses the selected
  GiB ceiling, and scheduling is Nice 19 with idle I/O for true or Nice 0 with
  best-effort I/O for false, with selected settings retained in lifecycle evidence
- **AND** signatures, socket privacy, runtime bounds, scan limits, and captured
  running invocations remain unchanged

### Requirement: Cleanup preserves unrelated processes

The controller SHALL record lifecycle evidence before engine launch and use a
unique transient systemd user unit with a maximum runtime to bound daemon
lifetime after failure or interruption.

#### Scenario: The host play fails or Ansible exits

- **WHEN** the stop play runs or the transient unit reaches its maximum runtime
- **THEN** systemd stops only that unit's engine without signaling unrelated
  processes or deleting the reusable database cache

#### Scenario: Native private socket lifecycle

- **WHEN** the transient daemon starts
- **THEN** systemd creates its mode 0700 RuntimeDirectory under the operator's private runtime root
- **AND** clients use its short absolute mode 0600 socket without PID/cwd aliases
- **AND** stopping the unit removes its runtime directory while preserving workspace evidence and cache
