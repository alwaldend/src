## Purpose

Make ordinary Bazel commands produce bounded, non-interactive diagnostics for
the managed Host Bot user without depending on a repository command wrapper.

## ADDED Requirements

### Requirement: Apply agent output defaults through user configuration

The managed user Bazel configuration SHALL disable color and curses, report
timestamped progress with a five-second rate limit, limit displayed action and test
output to 1 MiB each, and use a terse test summary with failing test output.
These defaults SHALL apply without selecting a named workspace configuration.

#### Scenario: Run an ordinary build or test

- **WHEN** the managed user runs a build or test that loads the managed Bazel
  configuration without overriding the output settings
- **THEN** the command receives the agent output defaults
- **AND** it does not require an AL Bazel wrapper or a workspace `agent` profile

#### Scenario: Override output for a particular invocation

- **WHEN** an operator supplies an explicit Bazel output option
- **THEN** Bazel's command-line precedence applies over the managed default

### Requirement: Preserve host-specific execution settings

The output-default migration SHALL preserve the existing managed cache paths,
SDK locations, and per-invocation resource limits. Setup documentation SHALL
identify the portable output settings separately from these host-specific values.

#### Scenario: Configure another installation

- **WHEN** an operator follows the repository setup instructions
- **THEN** they can select the output settings for their own user configuration
- **AND** the instructions do not require Host Bot's cache or SDK paths
