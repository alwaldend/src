## MODIFIED Requirements

### Requirement: Produce reusable source and rendered sites

`al_hugo_landing` SHALL combine a selected project README, generated Hugo
configuration, shared layouts and assets, and declared theme dependencies into
a Hugo-ready archive. `al_hugo_landing_site` SHALL additionally render that
archive into a site output. The shared assets SHALL include the canvas styles
that the repository's Hugo sites render, and the archive SHALL expose those
styles for the apex site to package into its own assets tree.

#### Scenario: Define a landing site for a project

- **WHEN** the caller supplies a project, title, documentation target,
  documentation URL, and repository URL
- **THEN** the macros package those inputs with the reusable shell and can
  produce rendered site files

#### Scenario: Render the shared canvas

- **WHEN** a landing site is rendered
- **THEN** it uses the shared canvas styles, and the apex site obtains the same
  styles from this shell rather than from a copy in its own tree

#### Scenario: Consume the shared canvas from the apex site

- **WHEN** the apex site packages the shared canvas styles into its assets tree
- **THEN** the shared style file is a declared input of that site's build and no
  second copy of its rules exists in the apex site tree
