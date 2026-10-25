## Purpose

Provide repository SSH commands that use the host client and AL's configurable
credential injection lifecycle.

## ADDED Requirements

### Requirement: Run the host SSH client

The SSH target SHALL resolve `ssh` from the host's `PATH` and pass fixed
arguments followed by invocation arguments without shell interpretation.

#### Scenario: Query the host client version

- **WHEN** the default SSH target is run with `-V`
- **THEN** it runs the host SSH client's version command

#### Scenario: Forward SSH options

- **WHEN** a consumer specifies fixed SSH arguments and invocation arguments
- **THEN** the host client receives them in that order with argument boundaries
  preserved

### Requirement: Preserve AL injection configuration

The macro SHALL accept AL configurations, plugin selection arguments, and
runtime data so consumers can inject keys before SSH starts and retain AL's
cleanup behavior. The default target SHALL select no authentication.

#### Scenario: Consumer selects key injection

- **WHEN** a consumer provides an AL configuration, plugin data, and matching
  plugin selection arguments
- **THEN** AL prepares the configured environment before starting SSH and owns
  cleanup after SSH exits
