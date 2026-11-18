## ADDED Requirements

### Requirement: Sign a server certificate using an existing key

The role SHALL create a certificate signing request using an existing managed-host
key and submit it to the configured Vault PKI signing endpoint. It SHALL verify
the certificate trust and hostname before installing the certificate and SHALL
notify the configured handler only when installation changes the certificate.

#### Scenario: Renew a missing, untrusted, or expiring certificate

- **WHEN** the current certificate fails trust or hostname verification, is absent, or expires within seven days
- **THEN** the role submits the signing request and installs the verified certificate without transferring or replacing the private key

#### Scenario: Retain a certificate that is still valid

- **WHEN** the current certificate remains trusted for the configured name and has at least seven days of validity
- **THEN** the role skips signing and certificate installation
