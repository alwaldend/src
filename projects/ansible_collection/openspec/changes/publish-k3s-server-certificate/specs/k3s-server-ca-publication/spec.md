## ADDED Requirements

### Requirement: Optional public K3s CA publication

The role SHALL publish only the public K3s server CA certificate to the
configured Vault KV v2 path after the service starts.

#### Scenario: A caller enables publication

- **WHEN** a caller configures a Vault path and runs the role
- **THEN** the role publishes `server-ca.crt` as `certificate`
- **AND** an unchanged certificate does not create another Vault version

#### Scenario: A caller leaves publication disabled

- **WHEN** the destination is empty or the role runs in check mode
- **THEN** the publication tasks make no Vault requests
