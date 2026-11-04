## ADDED Requirements

### Requirement: Order projects by explicit creation date

Every project landing page MUST declare its project creation date in the `date`
front matter field. The `/projects/` index SHALL list projects from newest to
oldest by that date. Rebuilding or editing a page SHALL NOT change its creation
date unless an author explicitly corrects the field.

#### Scenario: Browse projects of different ages

- **WHEN** the project index contains landings with different creation dates
- **THEN** each newer project appears above every older project

#### Scenario: Author a project landing page

- **WHEN** an author follows the repository Hugo skill to add a landing
- **THEN** the instructions require a manually maintained creation `date`
- **AND** no Git history or filesystem timestamp is needed during rendering
