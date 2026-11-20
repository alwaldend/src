## ADDED Requirements

### Requirement: Homepage infrastructure links

The homepage SHALL list Infra immediately after the Projects group. The Infra
row SHALL link to the infrastructure documentation index. Its indented child
links SHALL match the Projects rows and list the published top-level
infrastructure documentation sections in title order.

#### Scenario: Browse infrastructure from the homepage

- **WHEN** a visitor opens the homepage
- **THEN** the Infra group follows the complete Projects list
- **AND** the Infra heading and each displayed child link resolve to the
  corresponding documentation pages

## MODIFIED Requirements

### Requirement: Browser tab icon contrast

The site SHALL publish the S8 favicon with a graphite foreground on a
transparent background. Its fingerprinted URL SHALL reflect the icon
contents. Both browser color preferences SHALL use the same variant. The header
logo SHALL retain its transparent background and adaptive foreground.

#### Scenario: Display the favicon in a browser tab

- **WHEN** a browser renders the favicon with either color preference
- **THEN** the S8 monogram displays in graphite on a transparent background
