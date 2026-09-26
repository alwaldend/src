## MODIFIED Requirements

### Requirement: Credential and inventory isolation

The runner SHALL construct a test-specific environment and inventory without
inheriting production infrastructure credentials or inventories. Generated
test credentials SHALL remain temporary and excluded from retained artifacts.
Caller-supplied Ansible arguments SHALL be restricted to supported diagnostic
and tag-skipping options. Arguments that could override inventory, connection,
host, user, authentication, or extra variables SHALL fail preflight before
any guest or lifecycle playbook starts.

#### Scenario: Developer has production credentials configured

- **WHEN** a developer launches a test with infrastructure credentials in
  their environment
- **THEN** only the declared test inventory and test credentials reach the
  scenario and guest
- **AND** retained artifacts contain no private keys, passwords, or tokens

#### Scenario: Reject connection or inventory overrides

- **WHEN** caller arguments attempt to override connection, inventory, user,
  authentication, or extra variables
- **THEN** preflight rejects them with an actionable diagnostic
- **AND** no guest or Ansible lifecycle playbook starts and temporary state
  is removed
