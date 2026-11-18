## ADDED Requirements

### Requirement: Independent OpenCanary AppRole state

The OpenCanary AppRole and its SSH signing policy SHALL have a dedicated Terraform root and backend.

#### Scenario: Resolve shared membership

- **WHEN** the standalone root has provisioned the AppRole
- **THEN** the shared Vault root discovers the entity and group by name for Xen Orchestra and Ansible membership

### Requirement: OpenCanary server-certificate signing scope

The standalone root SHALL grant the OpenCanary AppRole access to sign server
certificate requests for `opencanary.alwaldend.com`.

#### Scenario: Request a certificate for a different name

- **WHEN** the OpenCanary AppRole submits a request for an unrelated name or a subdomain
- **THEN** the server-signing role rejects the request
