## ADDED Requirements

### Requirement: Shared development identity

The src_infra_dev AppRole SHALL read the existing dc1 and global DNS credentials through dns_access and use its own KV namespace for DNS state. It SHALL NOT read the retired OpenHands service secret. It SHALL NOT gain legacy Terraform-state access. Standard provisioning-group membership SHALL support its own scoped XO resource set.

#### Scenario: Inspect development permissions

- **WHEN** the Vault source for src_infra_dev is inspected
- **THEN** DNS credential reads are present and the OpenHands read is absent
- **AND** the unused runner-config read and legacy state grants are absent
- **AND** obsolete identities and their remaining grants are removed after the authorized handoff and empty-backend verification
