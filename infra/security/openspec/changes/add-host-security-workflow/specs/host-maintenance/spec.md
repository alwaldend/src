## ADDED Requirements

### Requirement: Reusable security maintenance

The system SHALL compose existing owner inventories and SHALL run the pinned Lynis role and canonical package updater without copying host identities.

#### Scenario: Normal maintenance

- **WHEN** an accessible supported host is selected
- **THEN** the playbook SHALL preserve a private audit report, update installed packages, enable automatic updates, disable automatic and immediate reboots, and verify the boot identity is unchanged

#### Scenario: Inaccessible host

- **WHEN** authenticated access fails
- **THEN** the playbook SHALL report the host as uninspected and SHALL NOT update it or label it inactive

#### Scenario: Audit or updater failure

- **WHEN** an audit or package-update task fails while the host remains reachable
- **THEN** the playbook SHALL compare the original and final boot identities,
  retain the task failure, and reject a changed boot

#### Scenario: Tagged maintenance

- **WHEN** the operator selects only `security_audit` or `update_all_packages`
- **THEN** only the selected role SHALL execute, with boot checks before and
  after it, including after an ordinary task failure
