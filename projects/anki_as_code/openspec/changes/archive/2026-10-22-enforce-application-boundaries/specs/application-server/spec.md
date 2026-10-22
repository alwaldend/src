## Purpose

Expose native server hosting through the same application boundary as collection and configuration operations.

## ADDED Requirements

### Requirement: Application-owned server operations

App SHALL expose ServerRun and ServerServe using an explicitly injected native runner. ServerServe SHALL delegate to ServerRun with serve prepended. App SHALL reject missing runner configuration and cancelled startup before invoking the runner, preserve runner error causes, and leave caller arguments unchanged. The executable SHALL supply native hosting and process signal adaptation without assembling Cobra server commands. CLI handlers SHALL forward command context, display name and native arguments to the application without running the native server or implementing server delegation.

#### Scenario: Serve through the application

- **WHEN** a CLI or Go caller requests ServerServe with native transport arguments
- **THEN** App invokes the configured runner with the display name, serve and unchanged transport arguments

#### Scenario: Rejected startup

- **WHEN** startup is cancelled or no runner is configured
- **THEN** App returns an error before native execution

### Requirement: Package responsibility documentation

Each first-party Go package SHALL describe its responsibilities, side effects, dependency and lifecycle boundaries in doc.go. Comments SHALL remain descriptive; implementation and review instructions belong in the Go skill.

#### Scenario: Command boundary

- **WHEN** the server command is executed
- **THEN** CLI syntax and output concerns remain in cli while App owns the server operation and native packages retain transport and storage behavior
