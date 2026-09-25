## MODIFIED Requirements

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
