## ADDED Requirements

### Requirement: Dedicated component identity

Vault configuration SHALL declare a `src_projects_x_article_uploader` identity
through the existing AppRole backend and reusable identity module, so the
uploader authenticates as itself rather than sharing another component's role.
Checked-in declarations SHALL contain secret references only.

#### Scenario: Inspect component authentication

- **WHEN** the uploader's identity is evaluated
- **THEN** it authenticates through the shared AppRole backend as
  `src_projects_x_article_uploader`
- **AND** it does not reuse another component's role or entity

### Requirement: Least-privilege access to the publication credential

The identity's policy SHALL grant access to the uploader's own credential path
and SHALL NOT grant access to another component's secrets. The credential value
SHALL live in Vault under the identity's own AppRole subtree, read through the
shared module's own-subtree policy, and the checked-in source SHALL reference it
by mount and path only.

#### Scenario: Reference the credential by path

- **WHEN** the credential mount and its logical path are evaluated for the
  injected reference
- **THEN** the checked-in source names the mount and the path and no credential
  value, so no secret enters source
- **AND** the path lies inside the identity's own AppRole subtree, so the shared
  module's own-subtree policy grants the read

#### Scenario: Read the publication credential

- **WHEN** the uploader reads its X API credential through the injection flow
- **THEN** the read succeeds against the uploader's own path
- **AND** the credential value is not present in checked-in source

#### Scenario: Attempt to read another component's secret

- **WHEN** the uploader's effective policy is evaluated against another
  component's credential path
- **THEN** the policy grants no read access to that path

### Requirement: Injected credential reference

The publication credential is an OAuth 1.0a user-context credential, which
requires four fields: the API key, the API secret, the access token, and the
access token secret. The injected reference SHALL name one environment variable
per field, and this identity's configuration SHALL declare those variables by
name. The references, not the values, are what checked-in configuration
carries; no credential value enters source.

#### Scenario: Inject the credential at run time

- **WHEN** the uploader's injection runs
- **THEN** each of the four credential fields is placed in its named environment variable
- **AND** the reference, not the value, is what checked-in configuration names
