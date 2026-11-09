## ADDED Requirements

### Requirement: Explicit member KV access

The module SHALL keep namespace access restricted to the role by default. When explicitly enabled, the existing member group SHALL also receive the role's existing namespace policy.

#### Scenario: Default component role

- **WHEN** allow_member_secret_access is omitted
- **THEN** member-group policy grants SHALL remain unchanged

#### Scenario: Personal role UI access

- **WHEN** the personal user role explicitly enables member secret access
- **THEN** its existing human members SHALL have secret-value access only to that role's KV namespace under this policy

#### Scenario: Metadata navigation from the mount root

- **WHEN** member secret access is enabled
- **THEN** the policy SHALL permit list-only access to exact metadata ancestors and the namespace root without a trailing slash
- **AND** it SHALL NOT grant value access or metadata reads outside the role's namespace
