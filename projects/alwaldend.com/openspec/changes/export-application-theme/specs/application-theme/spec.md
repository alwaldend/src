## ADDED Requirements

### Requirement: Shared application theme

The site SHALL expose an offline application stylesheet derived from the same
maintained palette and typography used by its public pages. Application consumers
SHALL compile those sources with the repository-pinned Sass and Bootstrap.
Shared selectors SHALL have one owner; extracting them SHALL preserve public-site
and Mermaid palette behavior. Applications MAY add workflow-specific layout rules
without duplicating shared theme values.

#### Scenario: Cleaner browser interface

- **WHEN** the cleaner serves its interface without internet access
- **THEN** its controls and content use the site's light and dark palette
- **AND** no runtime CDN request or copied palette is required

#### Scenario: Selection controls in either palette

- **WHEN** an application uses primary buttons or checkbox inputs
- **THEN** their filled, checked, and keyboard-focus states follow the active palette
- **AND** their shared presentation is defined by the site rather than its consumer
