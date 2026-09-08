# Historical project-agents delta

These requirements describe the archived change at its recorded acceptance point. They were imported without applying them to the current `project-agents` baseline. Current source and baseline requirements own present behavior.

## ADDED Requirements

### Requirement: Bounded routing coverage

CoverageMatrix MUST bind exact catalog and normalized case inputs, distinguish emitted cases from the complete skill universe, and expose truncation and provenance.

#### Scenario: A routing baseline is generated

- **WHEN** only a subset of registered skill cases is emitted
- **THEN** the matrix reports that subset and its exact inputs without implying complete behavioral coverage

### Requirement: Writable fixture inventory and gaps

Fixture evidence MUST identify deterministic writable or hermetic coverage for each representative trajectory and explicitly retain missing Terraform, Ansible, secret-injection, and end-to-end runner coverage.

#### Scenario: A trajectory has no writable fixture

- **WHEN** the fixture inventory evaluates a missing trajectory class
- **THEN** the inventory records the gap instead of classifying a static or unrelated test as full coverage

### Requirement: Isolated comparisons and measured adoption

Live comparisons MUST remain outside ordinary tests and bind model, skill, catalog, fixture, and judge identities; adopted optimizations MUST retain measured baselines, thresholds, review, regression, fallback, and retirement rules.

#### Scenario: An optimization is proposed for adoption

- **WHEN** a live comparison or learning proposal informs a maintained change
- **THEN** its identities and evidence chain remain inspectable and ordinary wildcard tests perform no live model calls
