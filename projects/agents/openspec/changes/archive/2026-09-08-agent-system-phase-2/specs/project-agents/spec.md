# Historical project-agents delta

These requirements describe the archived change at its recorded acceptance point. They were imported without applying them to the current `project-agents` baseline. Current source and baseline requirements own present behavior.

## ADDED Requirements

### Requirement: Catalog provenance and bounds

Derived catalogs MUST bind their owner inputs, versions, digests, conflicts, and completeness, and MUST expose truncation instead of claiming omitted entries are absent.

#### Scenario: Catalog inputs conflict

- **WHEN** the compiler encounters conflicting owner facts or a configured output bound
- **THEN** the output preserves the conflict or truncation and remains deterministic for identical inputs

### Requirement: Descriptor-only system index

AgentSystemIndex MUST contain catalog identities, versions, input digests, conflicts, and query routes without embedding catalog bodies.

#### Scenario: A catalog is indexed

- **WHEN** the system index is generated from checked catalogs
- **THEN** consumers receive bounded descriptors and routes to the owning catalogs

### Requirement: Bounded offline context

The offline context command MUST return relevant path, label, or task context with source provenance and explicit unavailable states from the same data used for human rendering.

#### Scenario: An optional input is unavailable

- **WHEN** a fresh-root context request cannot load an optional provider
- **THEN** the command returns a bounded partial result and safe discovery actions rather than inventing provider facts

### Requirement: Recoverable publication and runtime isolation

Interrupted historical goal publication MUST recover the prior valid record or idempotently complete the intended record; optional package faults MUST remain isolated from the fixed runtime control kernel.

#### Scenario: Publication or a package fails

- **WHEN** a publication boundary is interrupted or an optional package times out
- **THEN** recovery reports deterministic state and unrelated task namespaces remain usable
