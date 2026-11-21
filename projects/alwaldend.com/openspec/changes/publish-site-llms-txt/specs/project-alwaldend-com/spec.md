## ADDED Requirements

### Requirement: Site content index for agents

The site SHALL publish a plain-text Markdown index at `/llms.txt` using the
current home page, main menu, and top-level documentation sections. HTML pages
SHALL advertise the index through Docsy's discovery link.

#### Scenario: Discover the published site index

- **WHEN** a reader requests `/llms.txt` from the site root
- **THEN** the index lists the homepage, main navigation, and documentation links
- **AND** each link uses the current build's site URL
