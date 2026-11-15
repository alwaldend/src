## Purpose

Complete partial SSH host-key sets before shared signing and hardening while preserving existing host identities.

## ADDED Requirements

### Requirement: Complete host identities before signing

The shared SSH role SHALL create missing RSA, ECDSA and Ed25519 key pairs before discovering and signing public keys. Key creation SHALL preserve existing private keys. New RSA keys SHALL match the configured hardening key size before signing. Repeated execution on a complete key set SHALL report no changes.

#### Scenario: Partial key set

- **WHEN** only one supported host key exists
- **THEN** the role creates the missing key pairs before discovery
- **AND** the existing private key is unchanged

#### Scenario: Repeated configuration

- **WHEN** every host key pair already exists
- **THEN** key reconciliation reports no changes and preserves every private key

#### Scenario: Configured RSA strength

- **WHEN** a missing RSA key is created with a non-default hardening key size
- **THEN** the key has the requested size before its public key is signed
