## ADDED Requirements

### Requirement: Personal KV browser access

The existing simeonwarren human identity SHALL receive the existing user_simeonwarren namespace policy through its personal-role member group. This change SHALL grant no secret-value access to other users' namespaces and SHALL NOT make mounts anonymously visible. List-only access to exact metadata ancestors SHALL support navigation from the mount root.

#### Scenario: Existing browser login

- **WHEN** simeonwarren accesses Vault using an existing human authentication alias
- **THEN** its KV namespace SHALL be accessible and the existing KV mount SHALL be discoverable by that authenticated identity

#### Scenario: Browse from the secrets mount

- **WHEN** the human login lists each parent from the secrets mount root to its imagegen folder
- **THEN** each listing SHALL succeed and contain the next folder or token entry
- **AND** capability checks SHALL permit navigation with and without a trailing slash
- **AND** other namespaces' secret values SHALL remain denied
