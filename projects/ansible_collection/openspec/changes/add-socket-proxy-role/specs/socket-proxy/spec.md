## ADDED Requirements

### Requirement: Generic socket activation

The role SHALL install and activate a systemd socket at the listener
specified by the caller.

#### Scenario: Connect to the listener

- **WHEN** a client connects to the configured listener
- **THEN** systemd starts the unprivileged proxy for the configured backend

#### Scenario: Use HTTP ports on an SELinux host

- **WHEN** the caller opts in to HTTP ports on an SELinux-enabled host
- **THEN** the role allows the systemd socket proxy domain to bind and connect to HTTP-labeled TCP ports
- **AND** systemd can enter that domain while keeping `NoNewPrivileges` enabled
