## Purpose

Own reusable Vast.ai API operations independently of Terraform resource state.

## ADDED Requirements

### Requirement: Shared generated client

The shared internal SDK SHALL generate operations from the immutable official
OpenAPI schema, keep generated source in Bazel outputs, and provide bounded
credential-safe API adaptation without importing Terraform Framework packages.

#### Scenario: Provider consumer

- **WHEN** the Terraform provider builds and runs its CLI acceptance fixtures
- **THEN** it consumes the shared client and preserves lifecycle, response-validation, and credential-protection behavior.

#### Scenario: Explicit constructor configuration

- **WHEN** a consumer constructs a client
- **THEN** it supplies the endpoint and API key through an Options value.

#### Scenario: Escaped role names

- **WHEN** a role name contains a slash
- **THEN** generated reads and deletion retain that name as one escaped URL segment, including when adapting the configured API root.

#### Scenario: Already-selected team context

- **WHEN** a generated management credential already selects the requested team
- **THEN** context activation verifies the account and sends no redundant context mutation.

#### Scenario: Canonical current-account path

- **WHEN** a consumer verifies its current account
- **THEN** the adapter uses the canonical trailing-slash endpoint without following a redirect.

#### Scenario: API-key revocation confirmation

- **WHEN** key deletion returns an object or JSON string without a success flag
- **THEN** the client confirms absence through the generated key collection before reporting success, retains errors when the key remains, and skips deletion when that collection already confirms absence.

#### Scenario: Unavailable key collection

- **WHEN** the key-collection endpoint returns HTTP 404
- **THEN** the client reports a collection failure and preserves resource ownership instead of treating a key or team as absent.

#### Scenario: Duplicate team name

- **WHEN** team creation returns the documented name-conflict message with HTTP 400
- **THEN** the client reports a fixed name-conflict diagnostic without disclosing arbitrary response text.

#### Scenario: Masked environment value

- **WHEN** an existing variable's value is returned as the API's mask
- **THEN** the client distinguishes redaction from both a retrievable value and remote absence.
