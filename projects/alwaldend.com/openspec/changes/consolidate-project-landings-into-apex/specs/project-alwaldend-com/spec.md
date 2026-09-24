## MODIFIED Requirements

### Requirement: Build packaged repository documentation

The `site` target SHALL render the packaged repository documentation tree with
the configured Hugo, Docsy, Bootstrap, and local layouts and assets. Packaged
`README.md` documentation SHALL be renamed to Hugo `_index.md` pages. The
rendered site SHALL also package project landing content declared by registered
projects, so the published output contains both the repository documentation
tree and the `/projects` landing section.

#### Scenario: Build the main documentation site

- **WHEN** the `//projects/alwaldend.com:site` target is built
- **THEN** it produces rendered site output from declared documentation and
  asset inputs, with the release configuration using `https://alwaldend.com`.

#### Scenario: Build the site with project landing content

- **WHEN** the site is built while a registered project declares landing content
- **THEN** the rendered output contains both the repository documentation tree
  and that project's landing page

#### Scenario: Render the site without project landing content

- **WHEN** the site is built while no registered project declares landing
  content
- **THEN** the build succeeds and the rendered output contains the repository
  documentation tree
