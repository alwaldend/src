## MODIFIED Requirements

### Requirement: Reconcile DNS through its owner's Terraform state

Terraform-managed DNS SHALL use the shared DNS packaging with an owner-local
`dns/` root and state. Terraform DNS commands SHALL use the owner's Vault
AppRole and DNS credentials without selecting service provisioning settings.
Service commands SHALL retain their independent state and authentication.
The root SHALL consume the owner's canonical
`dns/dnsconfig.json` and use the provider views that owner requires. Provider
instance keys SHALL remain stable across value-only changes. Empty inventories
do not require an operational root.

#### Scenario: DNS currently shares provisioning state

- **WHEN** the DNS configuration moves to its dedicated root
- **THEN** the old root SHALL relinquish DNS bindings without deleting records
- **AND** the operator SHALL complete the state transfer before ordinary plans or applies from the new revision
- **AND** an authorized transfer SHALL import the same provider IDs and verify a no-change plan
- **AND** provisioning SHALL retain any canonical declaration inputs it uses for host configuration

#### Scenario: An owner already has independent DNS state

- **WHEN** the owner adopts shared packaging
- **THEN** its existing backend and Terraform resource addresses SHALL be preserved
