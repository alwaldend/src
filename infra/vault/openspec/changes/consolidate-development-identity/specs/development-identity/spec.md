## ADDED Requirements

### Requirement: Shared development identity

The src_infra_dev AppRole SHALL read the existing dc1 and global DNS credentials through dns_access and use its own KV namespace for DNS state. Its legacy service-secret read SHALL be limited to the required OpenHands key. It SHALL NOT gain legacy Terraform-state access or generic infrastructure-group membership.

#### Scenario: Inspect development permissions

- **WHEN** the Vault source for src_infra_dev is inspected
- **THEN** DNS credential reads and the OpenHands read are present
- **AND** the unused runner-config read and legacy state grants are absent
- **AND** old identities remain available until a separately authorized DNS-state handoff is complete
