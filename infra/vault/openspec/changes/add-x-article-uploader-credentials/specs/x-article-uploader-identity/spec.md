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

The identity's policy SHALL grant read access to the uploader's own credential
path and SHALL NOT grant access to another component's secrets. The credential
value SHALL live in Vault, and the checked-in source SHALL reference it by path.

#### Scenario: Read the publication credential

- **WHEN** the uploader reads its X API credential through the injection flow
- **THEN** the read succeeds against the uploader's own path
- **AND** the credential value is not present in checked-in source

#### Scenario: Attempt to read another component's secret

- **WHEN** the uploader's effective policy is evaluated against another
  component's credential path
- **THEN** the policy grants no read access to that path

### Requirement: Injected credential reference

The injected reference SHALL name exactly one environment variable for the
publication credential, and this identity's configuration SHALL declare that
variable by name. The reference SHALL be what checked-in configuration carries;
no credential value enters source.

#### Scenario: Inject the credential at run time

- **WHEN** the uploader's injection runs
- **THEN** the credential is placed in the named environment variable
- **AND** the reference, not the value, is what checked-in configuration names
