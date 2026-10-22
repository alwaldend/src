## Purpose

Provide a reviewable Anki CLI command surface independently of collection processing and storage implementation.

## ADDED Requirements

### Requirement: Discoverable commands

The CLI SHALL expose export, plan, apply, build, and generate-id with help and argument validation.

#### Scenario: Help

- **WHEN** the user requests root or subcommand help
- **THEN** the CLI prints the command names and supported arguments to stdout and exits successfully

### Requirement: Explicit operation placeholders

The scaffold SHALL return a nonzero status and an explicit not-implemented error on stderr for valid operation invocations, without reading or modifying collection files.

#### Scenario: Valid export invocation

- **WHEN** export receives its required input and output flags
- **THEN** the CLI reports that export is not implemented and creates no output

#### Scenario: Missing required arguments

- **WHEN** an operation lacks required flags or note paths
- **THEN** argument validation fails before the operation is invoked
