## MODIFIED Requirements

### Requirement: Explicit ownership activation

The DNS stage SHALL declare no managed records before explicit activation for
authorized adoption and SHALL retain enabled ownership afterward. Its
public-only integration SHALL require only Cloudflare, and operational commands
SHALL select the existing stage's provider injection.

#### Scenario: Inspect the prepared source defaults

- **WHEN** the project's prepared stage has not yet been activated for adoption
- **THEN** record ownership is disabled and the zone input is optional, while
  operational provider setup and authentication require their Vault prerequisites.

#### Scenario: Adopt existing public DNS records

- **WHEN** an authorized adoption freezes the old central writers and enables
  this project's record ownership
- **THEN** the stage uses its project state and Cloudflare configuration to
  import existing records without a RouterOS dependency.

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted project stage uses its checked-in source defaults
- **THEN** record ownership is enabled and the zone input remains optional.

#### Scenario: Reconcile existing public DNS records

- **WHEN** the project plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
  and does not require a RouterOS provider.
