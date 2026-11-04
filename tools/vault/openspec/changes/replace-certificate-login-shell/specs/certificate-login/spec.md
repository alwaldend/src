## Purpose

Provide certificate-based Vault login using the existing hardware certificate
and save a usable host token without exposing it or damaging an earlier token.

## ADDED Requirements

### Requirement: Authenticate and safely save a host token

The login command SHALL accept a Vault certificate role name and use the
configured Vault address and CA. It SHALL support PKCS#11 provider mode by
default and explicit legacy engine mode. Automatic selection SHALL ignore trust
modules and choose a currently valid hardware certificate suitable for TLS client
authentication. Role-matching identities SHALL take priority, with the newest
valid matching certificate preferred. A sole usable certificate SHALL remain
selectable when the role name differs. Selection SHALL NOT depend on token order.
An explicit certificate URI SHALL override automatic selection.
A successful login SHALL atomically replace the token file
with owner-only permissions. Input, discovery, transport, response, and storage
failures SHALL leave an existing token intact. Tokens and raw server responses
SHALL NOT appear in diagnostics. Requests SHALL have a finite timeout and
SHALL NOT follow redirects or retry authentication.

#### Scenario: Complete a certificate login

- **WHEN** a selected certificate authenticates and the response contains a token
- **THEN** the command saves that token with owner-only permissions
- **AND** its output confirms success without printing the token

#### Scenario: Reject an unusable response

- **WHEN** transport fails or the response is malformed, oversized, or tokenless
- **THEN** the command fails without replacing the existing token
- **AND** diagnostics omit the response and any token

#### Scenario: Select the hardware authentication certificate

- **WHEN** discovery exposes trust modules and hardware authentication certificates
- **THEN** login ignores the trust modules and selects a suitable hardware client
  certificate, preferring a matching role identity
- **AND** expired or server-only certificates are not selected
