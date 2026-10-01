## ADDED Requirements

### Requirement: Recover an expired server certificate with scoped Ansible tasks

The packaged Vault Ansible workflow SHALL expose a `vault_tls` tag that issues
and installs the configured server certificate and private key, then reloads
Vault. Selecting only this tag SHALL NOT update packages, the service binary,
host configuration, or Raft peers. Certificate material SHALL be suppressed
from task output. Explicit certificate renewal SHALL stop on issuance or
installation failure before reloading the service.

#### Scenario: Renew one authorized host's certificate

- **WHEN** the operator selects `vault_tls` and limits the play to one Vault host
- **THEN** only certificate issuance, certificate and key installation, and
  service reload are selected, in addition to read-only fact gathering
- **AND** a failed certificate task prevents the reload

### Requirement: Scope certificate-verification bypass to recovery invocations

The recovery procedure SHALL use an invocation-scoped verification override
only when an expired server certificate prevents the authorized renewal. It
SHALL retain the existing Vault authentication and injection flow, require a
valid bootstrap login, and verify the recovered service with certificate
validation restored. The procedure SHALL NOT persist insecure TLS defaults.

#### Scenario: Certificate expiry prevents authenticated renewal

- **WHEN** an authorized operator performs the documented recovery
- **THEN** the temporary override applies to credential injection and certificate
  issuance within that command
- **AND** a failed bootstrap login prevents proceeding to renewal
- **AND** success is established separately with normal TLS verification
