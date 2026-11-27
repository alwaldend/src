## ADDED Requirements

### Requirement: Invocation-scoped SSH local forwarding

The SSH resource provider SHALL authenticate using a supplied SSH client
certificate and verify the destination host certificate against its configured
CA. It SHALL publish an OS-assigned loopback address after confirming that the
remote destination accepts a TCP connection.

#### Scenario: Run a command through a private service forward

- **WHEN** an AL command selects an SSH forwarding resource
- **THEN** the provider opens a verified SSH forward and returns its loopback address
- **AND** it releases the listener and SSH connection when the invocation ends

#### Scenario: The destination or host trust fails

- **WHEN** the host certificate is rejected or the remote service cannot be reached
- **THEN** resource evaluation fails before the command runs
