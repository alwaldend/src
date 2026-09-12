# Source validation and operational completion

Source validation prepared the implementation for staged ownership transfer.
The user subsequently authorized deployment one owner at a time. Shared
operational checks and completed rollout are recorded in [operational evidence](operations.md);
owner-local adoption changes record verified transfers. All adopted roots retain
enabled ownership.

## Baseline source validation

- The pinned DNSControl baseline contains 168 records: 83 global and 85 dc1.
- The canonical HCL output matches names, types, contents, effective TTLs, MX
  priorities, and provider metadata for every baseline record.
- The optional exporter and BIND implementation were removed as requested.
  Runtime tables use direct file discovery without a central registry.
- The linter loads 46 canonical files, including the empty user declaration,
  and checks the 45 record owners for conflicting DNS names across files.
  Split record types and views do not permit separate files to share a name.
- The pinned Terraform lock workflow completed for all 45 owning roots.

The source check batch passed all 53 Terraform, skill and packaging tests.
All 45 owner wrappers, the Vault wrapper and the infrastructure skill built
successfully. The runtime table contains 46 files and 97 rows. A command-line
fixture splitting one name across files, types and views failed with both file
paths and record keys. All 46 owner/module/AppRole source changes passed strict
validation and were archived. These checks describe implementation preparation.
Final rollout candidate formatting, quality and semantic-lint gates remain
pending repository delivery; Git delivery receipts own their exact candidate
and publication outcome.

## Decisions supported by failures

The initial combined module required RouterOS configuration in global-only
roots. A dedicated global entry point now reuses the canonical normalization
and Cloudflare implementation, avoiding RouterOS for those owners. Operational
root commands retain AL injection and require actual DNS credential prerequisites.
The pinned RouterOS provider probes its version even for disabled resource
sets; source tests use mocks and no live providers.

Runtime discovery includes empty canonical files without inventing an owner
state for them. The exporter and its tests are absent from the final candidate.

## Completed operations and retained boundaries

The [operational completion evidence](operations.md#verified-completion) records
all owner receipts, retained provider records, final no-op plans, and resolver
checks. [Cutover](cutover.md) retains the scoped workflow and recovery procedure.
Writer controls apply within their documented observation coverage; removal of
current entrypoints does not prove every historical revision disabled. DNS
acceptance does not establish service health or replace delivery validation.
