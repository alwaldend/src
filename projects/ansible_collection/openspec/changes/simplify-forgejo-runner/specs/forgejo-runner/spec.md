## Purpose

Install a Forgejo Actions worker with a small configuration interface and
credential-safe, repeatable initial registration.

## ADDED Requirements

### Requirement: Consistent installation layout

The role SHALL derive its packaged executable, configuration, home, work and
registration paths from one installation directory, preserving the existing
default locations. Labels and capacity SHALL remain explicit inputs. Advanced
upstream configuration SHALL use one mapping and SHALL NOT override managed
registration, work-directory, label or capacity settings. The default host
executor prerequisites SHALL NOT include a container runtime.

#### Scenario: Alternate installation directory

- **WHEN** a deployment selects an alternate absolute installation directory
- **THEN** registration, daemon configuration and the unit reference that same layout
- **AND** the executable and configuration remain root-owned while runtime data belongs to the runner

#### Scenario: Advanced configuration

- **WHEN** a deployment supplies advanced runner configuration
- **THEN** those settings are rendered as upstream YAML
- **AND** role-managed paths, labels and capacity remain authoritative

### Requirement: Initial registration only

The role SHALL register only when its registration file is absent, require an
instance, token and name only in that case, and bound registration execution.
It SHALL pass the token through stdin with task logging suppressed. Failure
SHALL produce a fixed diagnostic without credential-bearing process output.
The registration file SHALL be restricted to its runner account. Neither
registration nor the daemon SHALL require OIDC or a Vault client on the worker.

#### Scenario: Existing registration without bootstrap credentials

- **WHEN** the registration file exists and no token or instance is supplied
- **THEN** registration is skipped and the existing file is preserved

#### Scenario: Failed initial registration

- **WHEN** the registration process rejects the supplied inputs
- **THEN** deployment fails with the fixed registration message
- **AND** neither stdout, stderr nor the token appears in task output

### Requirement: Change-triggered service restart

The role SHALL enable and start the worker and restart it through its handler
when its executable, configuration, registration or unit changes.

#### Scenario: Unchanged deployment

- **WHEN** registration and installed files are unchanged
- **THEN** deployment SHALL NOT request a service restart
