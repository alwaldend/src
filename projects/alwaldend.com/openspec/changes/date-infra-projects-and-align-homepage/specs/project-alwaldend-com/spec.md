## MODIFIED Requirements

### Requirement: Homepage infrastructure links

The homepage SHALL link to the infrastructure documentation index and list
the published top-level Infra projects in descending creation-date order.
Each Infra project SHALL declare its creation date in its README front matter.

#### Scenario: Browse infrastructure from the homepage

- **WHEN** a visitor opens the homepage
- **THEN** Infra shows its project links from newest to oldest
- **AND** the Infra index and each project link resolve to their documentation

## ADDED Requirements

### Requirement: Responsive homepage project groups

The homepage SHALL display Projects above Infra. Each group's index link
SHALL span its list. Each list SHALL display two columns of entries when
space permits and one column otherwise. Entries SHALL read across rows in
date order. General navigation links SHALL appear above both lists.

#### Scenario: Resize the homepage

- **WHEN** the available width fits two columns of entries within a list
- **THEN** each list displays its entries in two columns
- **AND** each list retains its reading order on a narrow screen
