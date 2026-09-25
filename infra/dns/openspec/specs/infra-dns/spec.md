# Infrastructure DNS

## Purpose

Describe owner-local Terraform DNS declarations and runtime source-ownership
validation for global and dc1 records. This contract describes checked-in
source; it does not establish live adoption or authorize infrastructure writes.

Sources: [component documentation](../../../README.md),
[target definitions](../../../BUILD.bazel),
[runtime linter](../../../cmd/lint/main.go), and
[Terraform root](../../../tf/dns.tf).

## Requirements

### Requirement: Aggregate records from their owning components

DNS declarations SHALL remain in their owners' `dnsconfig.json` files. Runtime
inventory SHALL discover all canonical files, including nested modules and
empty declarations, without a checked-in ownership registry. The inventory
SHALL print a deterministic table of files, DNS names, types and views.

Project landing declarations retain their declared direct targets.

#### Scenario: A new project adds a canonical declaration

- **WHEN** the runtime linter scans a workspace containing the new file
- **THEN** it includes the file without requiring a registry update

### Requirement: Preserve distinct global and site-local views

The Terraform module SHALL preserve Cloudflare global and MikroTik dc1 views.
It SHALL retain every type member, expand and deduplicate destinations, reject
unknown inputs and preserve effective values, TTLs, priorities and multiplicity.

#### Scenario: An entry declares A and AAAA in all views

- **WHEN** the shared module normalizes the entry
- **THEN** both types appear exactly once in each destination view

#### Scenario: An input names an unsupported destination

- **WHEN** a record names an unsupported destination
- **THEN** normalization fails with an actionable diagnostic

### Requirement: Validate configuration without provider access

The linter and module tests SHALL validate source ownership, names, inputs,
normalization and provider mappings without credentials or live provider access.
Existing JSON-based VM consumers SHALL remain compatible.

#### Scenario: Validate a changed declaration

- **WHEN** offline checks run
- **THEN** they validate the declaration and fixtures without deploying records

### Requirement: Inject credentials into operational wrappers

Owner Terraform roots SHALL use the existing AL/Vault authentication, backend
and credential injection workflow. They SHALL obtain only their required view
credentials. Operational provider initialization may require live credentials
even with resource creation disabled. Central DNSControl write entrypoints
SHALL be removed from the candidate.

#### Scenario: A global-only owner operates its root

- **WHEN** an authorized operator invokes the root after bootstrap
- **THEN** it uses that owner's AppRole and Cloudflare credentials
- **AND** RouterOS credentials and resources are not required

### Requirement: Reconcile DNS through its owner's Terraform state

Each owner SHALL instantiate the shared module in `tf_setup` when present,
otherwise in `tf`. Missing roots and AppRoles SHALL be added. Related
AppRole resources SHALL be grouped in modules under
`infra/vault/tf/approles/<name>`. Provider instance keys SHALL be stable across
value-only changes; unrelated owners SHALL remain outside the state.

#### Scenario: An owner has both stages

- **WHEN** the integration is prepared
- **THEN** only its setup root declares the DNS module

### Requirement: One declaration file owns each DNS name

The runtime linter SHALL reject different files managing the same canonical
fully qualified DNS name, across types and views. Multiple records within one
file SHALL be allowed. Diagnostics SHALL identify both files.

#### Scenario: Different files split a name across types or views

- **WHEN** one file declares an A record and another an AAAA record or a
  different view for the same name
- **THEN** lint fails and identifies both source files

#### Scenario: One file manages multiple records for a name

- **WHEN** one file declares several values, types or views for a name
- **THEN** ownership lint accepts the sole source

### Requirement: Preserve live records during transfer

Operational cutover SHALL audit active and scheduled central writers, stop
identified competing writers, and coordinate one owner at a time. Evidence
SHALL state the audit's coverage and unavailable observations. Each existing
record SHALL be imported using its actual provider ID into exactly one owner
state, followed by a no-change adoption plan before writes. Fresh complete
inventories MAY establish missing declarations for an exact additions-only
plan that preserves every existing record. Endpoint-owning batches SHALL
verify an independent authenticated recovery path first.

The shared apex and mail records owned by `infra/dns` SHALL remain disabled
by default in its `tf` root before the authorized adoption revision and SHALL
retain enabled ownership after import. Reconciliation against the adopted state
and unchanged declarations SHALL propose no record additions, changes,
replacements, or deletions and SHALL preserve unrelated provider records.

#### Scenario: An owner is adopted

- **WHEN** its authorized cutover runs
- **THEN** exact-ID imports preserve existing records
- **AND** the operator coordinates exclusive active writers and verifies
  unrelated records remain unchanged

#### Scenario: Deploy verified missing declarations

- **WHEN** complete provider inventories establish that declared names are absent
  and conflict-free
- **THEN** the reviewed DNS plan creates only those missing declarations
- **AND** every existing provider identity and attribute remains unchanged
- **AND** a follow-up DNS plan contains no changes

#### Scenario: Inspect the prepared shared-record defaults

- **WHEN** the shared `tf` root uses its checked-in defaults before its
  authorized adoption revision
- **THEN** shared-record ownership is disabled and the root retains the
  canonical shared declaration and module inputs

#### Scenario: Inspect adopted shared-record defaults

- **WHEN** the adopted shared `tf` root uses its checked-in source defaults
- **THEN** ownership is enabled for the shared apex and mail records in
  their declared views

#### Scenario: Reconcile existing shared records

- **WHEN** the shared `tf` root plans against adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
- **AND** unrelated provider records remain unchanged

### Requirement: Avoid a separate DNS exporter

The implementation SHALL derive its declaration views from the owning
`dnsconfig.json` files rather than maintaining a normalized inventory or BIND
exporter. Historical provider snapshots SHALL NOT be presented as current
desired state, and no provider snapshot SHALL be committed as a declaration
source.

#### Scenario: Inspect declared ownership

- **WHEN** the runtime linter succeeds
- **THEN** its table reports the discovered declarations without an export step

### Requirement: Publish declaration pages per destination view

The implementation SHALL render one documentation page per destination view
from the declared records, with the owning declaration for every record.
Regeneration SHALL require an explicit manual command. Ordinary validation
SHALL check declarations and ownership without requiring checked-in pages to
match the declarations or rewriting those pages. Documentation SHALL identify
the declarations as authoritative and explain that snapshots may lag.

#### Scenario: Regenerate a declaration page

- **WHEN** an operator explicitly runs the documentation generation command
- **THEN** each destination page lists the current declarations with their owner
- **AND** an explicit freshness check succeeds against the regenerated pages

#### Scenario: A declaration page is stale

- **WHEN** a declaration changes while its checked-in page remains unchanged
- **THEN** ordinary declaration validation does not fail because of that stale page
- **AND** the page is not rewritten unless regeneration is explicitly invoked

#### Scenario: Check freshness on demand

- **WHEN** an operator explicitly checks a stale page
- **THEN** the command reports that regeneration is needed
