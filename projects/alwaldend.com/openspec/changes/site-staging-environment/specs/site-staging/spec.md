## ADDED Requirements

### Requirement: Publish the apex site to a staging hostname

The apex project SHALL provide a staging build that renders the same packaged
source as the production site with the staging base URL, and a deployment
command that publishes that output to the staging Pages repository. The
staging build SHALL NOT alter the production base URL or the production
publication branch.

#### Scenario: Render the staging candidate

- **WHEN** the staging site target is built
- **THEN** it renders the same content as the production site with
  `https://www-staging.alwaldend.com` as its base URL
- **AND** the production target continues to render `https://alwaldend.com`

#### Scenario: Publish to staging

- **WHEN** the staging deployment command runs
- **THEN** it publishes the rendered output and the staging `CNAME` to the
  staging repository's `pages` branch

### Requirement: Derive the staging hostname from checked-in declarations

The DNS record and the repository that serve the staging hostname SHALL be
declared in their owning checked-in sources, and the repository SHALL keep a
protected default branch distinct from the branch that serves site content.

#### Scenario: Declare the staging hostname

- **WHEN** the DNS declarations are linted
- **THEN** `www-staging.alwaldend.com` is declared by exactly one source file

#### Scenario: Protect the staging default branch

- **WHEN** the staging repository is created
- **THEN** it has a default branch that the managed ruleset protects
- **AND** site content is published to a separate unprotected branch

### Requirement: Present the project index on the documentation axis

The `/projects/` section index SHALL resolve a layout that renders the same
centered documentation column as `/docs/` and `/blog/` without either sidebar,
so its heading starts at the same offset and its content shares the same
column width. It SHALL list the landings alphabetically by title with each
entry's description, status, and tags linked to their taxonomy term pages.

#### Scenario: Match the documentation page geometry

- **WHEN** `/projects/`, `/docs/`, and `/blog/` are measured at the same
  viewport width
- **THEN** each page places its heading at the same top offset and its content
  at the same left edge and width
- **AND** `/projects/` renders neither the section sidebar nor the table of
  contents sidebar

#### Scenario: List the landings

- **WHEN** the section index renders
- **THEN** it lists the registered landings in title order
- **AND** each entry carries its status and tags as links to the corresponding
  taxonomy term pages

### Requirement: Verify the replacement before retiring per-project publication

The landing, documentation, and taxonomy routes served by the apex site SHALL
be verified on the staging hostname before the production deployment and
before any per-project hostname, DNS record, or Pages repository is retired.

#### Scenario: Verify the staged replacement

- **WHEN** the staging publication completes
- **THEN** `/projects/`, `/projects/<name>/`, and `/docs/projects/<name>/` are
  served from the staging hostname
- **AND** landing pages populate the `statuses`, `languages`, and `tags`
  taxonomy pages

#### Scenario: Retire per-project publication

- **WHEN** the apex site serves the corresponding routes
- **THEN** the former per-project hostnames, their DNS records, and their
  landing Pages repositories are retired
