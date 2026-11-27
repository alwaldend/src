## ADDED Requirements

### Requirement: Generic socket activation

The role SHALL install and activate a systemd socket at the listener
specified by the caller.

#### Scenario: Connect to the listener

- **WHEN** a client connects to the configured listener
- **THEN** systemd starts the unprivileged proxy for the configured backend
