## ADDED Requirements

### Requirement: Renew server certificates before expiry

Vault's desired host configuration SHALL declare automatic renewal with an
explicit interval and renewal margin shorter than the certificate lifetime.
Successful renewal SHALL install a valid matching certificate/key pair and
reload the service only after installation succeeds. Failed issuance or
installation MUST preserve the previously usable material and report failure.

#### Scenario: A certificate enters its renewal window

- **WHEN** the declared renewal job observes a certificate inside its configured
  renewal window and the issuer is available
- **THEN** it installs a replacement with a later expiry and reloads the service
- **AND** a normal TLS check confirms the served replacement

#### Scenario: Renewal fails

- **WHEN** issuance or installation fails
- **THEN** the existing usable certificate/key pair remains available
- **AND** the operation reports failure without claiming successful renewal

### Requirement: Report certificate expiry risk without disclosing secrets

The desired operational configuration SHALL report renewal failures and a
certificate crossing an explicit expiry-warning threshold through a declared
existing reporting destination. Diagnostic output MUST NOT contain private
keys or authentication credentials.

#### Scenario: Repeated renewal failure approaches expiry

- **WHEN** renewal remains unsuccessful and the expiry threshold is crossed
- **THEN** the declared report identifies the affected endpoint and expiry risk
- **AND** it carries no key or credential material

### Requirement: Rehearse recovery when normal TLS authentication is unavailable

The operator procedure SHALL name independently obtainable recovery
authentication prerequisites and a scoped renewal path that does not depend
on successful normal TLS verification of the expired endpoint. Any temporary
verification exception MUST remain invocation-scoped and recovery acceptance
MUST include a subsequent normal TLS check.

#### Scenario: The served certificate has expired

- **WHEN** an authorized recovery exercise starts with ordinary credential
  injection blocked by the expired endpoint
- **THEN** the documented independent prerequisites allow scoped renewal
- **AND** an unavailable prerequisite stops the procedure with an actionable error
- **AND** successful recovery is verified with normal certificate validation
