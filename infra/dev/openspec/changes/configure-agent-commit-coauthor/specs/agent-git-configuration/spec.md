## ADDED Requirements

### Requirement: Managed agent co-author

The dev1 Ansible configuration SHALL provide the operator's co-author identity
through `al.deliveryCoAuthor` in the agent account's global Git configuration.

#### Scenario: Configure the dev1 agent account

- **WHEN** the dev1 Git configuration is deployed
- **THEN** the global Git config contains the operator's co-author identity for agent delivery
