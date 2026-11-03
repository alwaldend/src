## ADDED Requirements

### Requirement: Independent DNS ownership

The component SHALL manage its canonical DNS document from a separate Terraform
root and Vault state namespace using its own AppRole and DNS credentials.
DNS SHALL require neither Yandex credentials nor GPU provisioning variables.
The provisioning root MUST NOT manage these Cloudflare records.

#### Scenario: Plan DNS before host provisioning

- **WHEN** an operator invokes the DNS plan
- **THEN** only the canonical Cloudflare records are planned from independent state

### Requirement: Active configuration only

The DNS project SHALL expose only its current backend and record configuration.
It MUST NOT retain migration roots, legacy backend wiring, or import-only inputs.
The existing adopted resource IDs and DNS values MUST remain unchanged.

#### Scenario: Run the independent DNS project after adoption

- **WHEN** an operator plans the DNS project
- **THEN** it uses only the DNS backend and manages the existing records without
  compatibility or migration configuration
