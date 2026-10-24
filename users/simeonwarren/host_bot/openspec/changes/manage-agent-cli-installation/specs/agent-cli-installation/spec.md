## Purpose

Provision the requested coding-agent commands for the managed Host Bot user
through the host's declarative Ansible deployment.

## ADDED Requirements

### Requirement: Install the requested agent CLIs for the managed user

The Host Bot deployment SHALL install Codex and Grok through their official
shell installers as the managed user, exposing both commands in that user's
`~/.local/bin`. Each installation SHALL be selectable independently.

#### Scenario: Install a missing command

- **WHEN** the selected agent command is absent during deployment
- **THEN** its official installer runs as the managed user
- **AND** the installed command can report its version

### Requirement: Preserve existing installations and report installer failures

The deployment SHALL skip an installer when its managed command already
exists. A failed installer download or execution SHALL fail the task.

#### Scenario: Repeat deployment

- **WHEN** an installation task runs after its command has been installed
- **THEN** the installer does not run again
- **AND** the task reports no change

#### Scenario: Download fails

- **WHEN** an installer download fails
- **THEN** the installation task fails
- **AND** a successful downstream shell exit does not hide that failure
