## ADDED Requirements

### Requirement: Shared forwarding default

The SSH wrapper SHALL supply its `ssh_hardening_match_user` default to the
existing devsec role as `ssh_server_match_user`. The default SHALL match only
`ansible` with `AllowTcpForwarding remote`, `PermitListen none`, and
`AllowStreamLocalForwarding remote`. The scanner SHALL NOT change SSH server
policy.

#### Scenario: Use the shared default

- **WHEN** a caller invokes SSH hardening without overriding the wrapper list
- **THEN** the existing devsec template renders the three-rule `ansible` match
- **AND** no owner inventory needs to repeat that default

#### Scenario: Override the wrapper list

- **WHEN** a caller supplies an empty or custom `ssh_hardening_match_user` list
- **THEN** the supplied list replaces the default without recursive forwarding

### Requirement: Normal SSH setup

The existing SSH role SHALL retain its client CA installation, host-key
discovery and signing, signed-certificate mapping, and native devsec handler.
Its existing devsec include SHALL supply the shared forwarding default.

#### Scenario: Select normal SSH setup

- **WHEN** an owner deployment selects its SSH role using the `ssh` tag
- **THEN** normal discovery and signing precede the upstream hardening tasks
- **AND** the normal role applies the shared policy without a separate
  maintenance playbook or certificate-input interface
