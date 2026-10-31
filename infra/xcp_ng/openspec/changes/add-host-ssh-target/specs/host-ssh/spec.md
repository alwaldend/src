## Purpose

Provide component-owned authenticated host access for authorized XCP-ng operations.

## ADDED Requirements

### Requirement: Packaged host SSH access

The component SHALL expose an SSH target that obtains its password from the
existing Vault host credential reference, verifies the host key, and forwards
the caller's remote command without printing the password or passing it as an argument.

#### Scenario: Authorized hardware inspection

- **WHEN** the caller runs the target with a read-only hardware command
- **THEN** the command runs on the XCP-ng host using the packaged authentication flow.

#### Scenario: Untrusted host key

- **WHEN** the server host key is unknown or differs from the trusted key
- **THEN** the target refuses the connection without accepting the key automatically.
