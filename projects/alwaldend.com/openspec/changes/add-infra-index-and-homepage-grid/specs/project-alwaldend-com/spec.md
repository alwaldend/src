## MODIFIED Requirements

### Requirement: Homepage navigation panel

The homepage SHALL display one bordered panel below the site header with
four sections titled Site, External, Projects, and Infra. The sections SHALL
form a two-column grid when space permits and stack in one column on narrower
screens. Site SHALL link to Blog and Docs. External SHALL link to the Alwaldend
namespace pages on GitHub, GitLab, Docker Hub, and Hugging Face. Projects and
Infra SHALL link to their respective indexes and list their entries by creation
date, newest first. All four sections SHALL use the same link row styling.
Each row SHALL display its destination beside the title on desktop and below
it on small screens.

#### Scenario: Browse projects from the homepage

- **WHEN** a visitor opens the homepage
- **THEN** the four sections display every link without requiring expansion
- **AND** the link titles and destinations fit within the viewport

## ADDED Requirements

### Requirement: Infrastructure index

The site SHALL link to `/infra/` from its shared header and the homepage
Infra section. The index SHALL use the Projects index presentation and list
the published infrastructure documentation sections by creation date, newest
first, with links to those sections.

#### Scenario: Browse infrastructure

- **WHEN** a visitor selects Infra in the header
- **THEN** the Infra index opens
- **AND** selecting an entry opens its published documentation
