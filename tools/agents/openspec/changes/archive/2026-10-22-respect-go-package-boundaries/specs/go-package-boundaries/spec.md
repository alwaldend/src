## Purpose

Keep Go package descriptions accurate while implementation and review instructions direct agents to respect intended boundaries.

## ADDED Requirements

### Requirement: Respect documented Go package boundaries

The repo-go skill SHALL require agents to read affected packages' doc.go before implementation or review, identify intended responsibility and effect owners, and check implementation and call paths against the approved design. Package comments SHALL remain descriptive; instructions to respect the design belong in the skill. The skill SHALL require accurate package descriptions and preserve a single owner for each responsibility. A comment SHALL NOT be broadened solely to legitimize a boundary violation. CLI packages SHALL retain input/output adaptation while applications coordinate operations and lower-level packages retain their documented domain and lifecycle duties.

#### Scenario: Implementation bypasses its owner

- **WHEN** a command directly starts native work assigned to the application
- **THEN** the workflow identifies the misplaced call and restores application delegation rather than merely rewriting the CLI documentation

#### Scenario: Inaccurate package documentation

- **WHEN** doc.go contradicts actual ownership or documented effects
- **THEN** the workflow checks source and consumers and corrects the contract or misplaced behavior within authorized scope
