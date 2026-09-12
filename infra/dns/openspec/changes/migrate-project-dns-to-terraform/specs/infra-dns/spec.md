## MODIFIED Requirements

### Requirement: Aggregate records from their owning components

DNS declarations SHALL remain owned by their components, including project
landing CNAMEs in `dnsconfig.json` and nested projects exposed through their
declared Bazel dependencies. Aggregation SHALL preserve a complete declared
inventory while selecting only unmigrated records for central deployment.
Landing CNAMEs SHALL retain their direct `alwaldend.github.io.` targets.

#### Scenario: A registered project contributes landing records

- **WHEN** a registered project's declaration is loaded
- **THEN** its landing records retain their owning source and direct targets
- **AND** migrated records remain available to inventory and documentation
  without being declared in the central operational configuration

### Requirement: Preserve distinct global and site-local views

DNS management SHALL preserve Cloudflare global and MikroTik dc1 views,
including common apex/mail records. It SHALL route `all` to both views,
preserve every type member in a logical record entry, and reject unknown
types or destinations. Migration SHALL preserve effective record values,
names, multiplicity, TTLs, and provider-specific settings.

#### Scenario: A record is assigned only to dc1

- **WHEN** a record declares `dsp: ["dc1"]`
- **THEN** only the dc1 provider receives that declaration

#### Scenario: A common entry contains multiple types

- **WHEN** an entry declares A and AAAA records with `dsp: ["all"]`
- **THEN** both types appear once in each destination view

#### Scenario: A record names an unsupported destination

- **WHEN** an input contains an unknown destination or record type
- **THEN** offline validation fails with an actionable diagnostic

### Requirement: Validate configuration without provider access

Configuration validation SHALL check input structure, normalization, view
routing, ownership conflicts, and remaining DNSControl configuration without
provider credentials or network access. Existing JSON consumers SHALL remain
compatible throughout the migration.

#### Scenario: Validate a proposed record change locally

- **WHEN** package validation runs on a changed declaration
- **THEN** it checks the affected Terraform inputs and legacy configuration
  without contacting providers or deploying records
- **AND** VM setup consumers can still read declared names and addresses

### Requirement: Inject credentials into operational wrappers

Project Terraform and remaining DNSControl operational wrappers SHALL use
the existing AL/Vault injection workflow. DNSControl wrappers SHALL pass
`providers.json` explicitly while they remain operational. Credentials SHALL
remain outside checked-in declarations and documentation output. Once every
scope is transferred, central live DNSControl writers SHALL be retired.

#### Scenario: An authorized operator previews a DNS change

- **WHEN** the central preview wrapper is invoked during migration
- **THEN** it receives the declared configuration and `providers.json`
- **AND** it reads credentials from injected environment references

#### Scenario: A project applies its DNS records

- **WHEN** an authorized operator invokes the owning Terraform root
- **THEN** it obtains the configured provider credentials through injection
- **AND** no DNS provider credential is copied into the JSON declaration

#### Scenario: The final scope has migrated

- **WHEN** all project and common records have verified Terraform ownership
- **THEN** central DNSControl deployment entry points are no longer available
- **AND** any retained BIND renderer has no live-provider credentials

## ADDED Requirements

### Requirement: Reconcile records through their owning project state

The migration SHALL require project integrations to use the shared Terraform
DNS module with canonical JSON in `tf_setup` when present, otherwise in `tf`.
Before implementation, each integration SHALL link to its owning OpenSpec
change for project-specific bootstrap, state, and deployment requirements;
this specification owns the common DNS migration contract. Each provider
record or indivisible RRset SHALL have exactly one owning state. The module
SHALL preserve stable resource identities across value-only edits and SHALL
not reconcile undeclared records owned by another project.

#### Scenario: Both Terraform stages exist

- **WHEN** a project has both `tf_setup` and `tf`
- **THEN** only `tf_setup` instantiates its DNS module

#### Scenario: A project migration batch is prepared

- **WHEN** implementation begins for a project or the reusable module
- **THEN** the migration map links an affected-owner OpenSpec change
- **AND** that owner maintains its specific requirements without copying
  this coordination contract

#### Scenario: The setup stage does not exist

- **WHEN** a project has no `tf_setup`
- **THEN** its DNS module is instantiated in `tf`, creating that root if needed

#### Scenario: A project updates and removes its own records

- **WHEN** an owned record value changes or its declaration is removed
- **THEN** Terraform plans the corresponding update or deletion
- **AND** unrelated project records remain outside that reconciliation

### Requirement: Transfer DNS ownership without recreating records

Migration SHALL proceed by bounded scopes with central exclusions and exact
name/descendant ignores for transferred subdomains in each affected view.
Existing records SHALL be imported into their single designated state before
Terraform writes. Cutover SHALL prevent competing old and new writers and
SHALL NOT use global purge suppression or disabled ignore safety checks.

#### Scenario: A subtree is transferred

- **WHEN** `cloud.alwaldend.com` transfers to a project state
- **THEN** current central runs neither declare nor alter its apex or any
  descendant record in the transferred views
- **AND** adoption imports preserve provider identities and effective data
- **AND** central runs continue to manage unmigrated scopes

#### Scenario: A project owns an additional alias

- **WHEN** the project's inventory includes names outside its main subtree
- **THEN** those names have explicit ownership and migration protection too

#### Scenario: A scope owns a management endpoint's DNS name

- **WHEN** a batch includes DNS required to reach its provider or recovery
  dependencies, including `router1.dc1.alwaldend.com`
- **THEN** an independently reachable, authenticated management path is
  defined in the owning IaC and verified before the batch is transferred
- **AND** recovery remains possible when the managed record is unavailable

### Requirement: Evaluate BIND documentation without duplicating ownership

The migration SHALL evaluate retaining BIND generation as an offline derived
view of canonical declarations. If retained, it SHALL use shared translation
logic or its normalized output, include migrated and unmigrated records in
both views, require no credentials or live-provider reads/writes, and identify
the source revision and desired-state nature of the output. If these
constraints cannot be met practically, the migration SHALL document the
limitation and remove stale documentation claims without blocking completion.

#### Scenario: A migrated record appears in retained documentation

- **WHEN** offline documentation generation runs after a project migrates
- **THEN** its declared records remain represented in the correct BIND view
- **AND** identical source and declared generation inputs produce identical
  output without contacting live systems or copying record definitions

#### Scenario: Offline reuse is not feasible

- **WHEN** feasibility work finds no practical shared generation path
- **THEN** the decision and evidence are recorded explicitly
- **AND** obsolete BIND output is not presented as current documentation
