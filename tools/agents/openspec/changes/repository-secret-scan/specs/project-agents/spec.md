## ADDED Requirements

### Requirement: First-party nested modules are not new external dependencies

Repository policy and dependency skills SHALL permit creating and registering
first-party nested modules within the authorized task without separate dependency
approval. Redeclaring software already used at existing repository pins in the
module's manifests and generated locks SHALL NOT count as new external software.
Genuinely new external software SHALL retain the existing approval requirement.

#### Scenario: Create a nested rules project

- **WHEN** an authorized project needs its own module using existing pinned tools
- **THEN** the agent creates the module and its manifests and locks without
  requesting redundant dependency approval
- **AND** an explicitly selected project path takes precedence over placement defaults

### Requirement: Repository secret-scan procedure

The repository SHALL provide a discoverable `repo-secret-scan` skill for running
the pinned scanner with an explicit scope, interpreting findings separately
from scan failures, and retaining original local evidence for user inspection.
It SHALL distinguish public source metadata and demonstrated synthetic fixtures
from usable credentials. Redacting a report SHALL NOT be described as remediation
of a credential already committed to public Git history.
The procedure SHALL use scan-only execution without credential verification and
describe exact detector/path/commit/line matching and suppressed/unmatched counts.

#### Scenario: Scan repository history

- **WHEN** the user requests a normal repository secret scan
- **THEN** the procedure selects history reachable from `HEAD`, identifies
  exclusions and working-tree coverage limits, and reports the revision and
  verified and unverified finding counts
- **AND** an all-ref audit is selected explicitly rather than implied

#### Scenario: Inspect a finding in public Git history

- **WHEN** the operator needs to distinguish a real credential from a false positive
- **THEN** the procedure retains the original values in a restricted local report
  and provides its location together with detector, source path, commit, and line
- **AND** credentials are not copied into a pull request or shared report
- **AND** a public test fixture is not declared harmless solely from its file name

#### Scenario: Scanner terminates with errors

- **WHEN** a scan fails or has not demonstrably completed
- **THEN** the procedure reports that limitation rather than claiming the repository
  has no secrets
