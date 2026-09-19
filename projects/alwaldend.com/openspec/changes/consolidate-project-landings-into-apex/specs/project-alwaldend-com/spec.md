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

### Requirement: Separate preview and publication

`site_serve` SHALL serve built content on `127.0.0.1:1313`. The deployment
archive SHALL contain the apex `CNAME`, and the explicit deployment workflow
SHALL write `.nojekyll` and push changed output to the configured Pages branch.
Building or previewing the site SHALL NOT itself publish it. The rendered site
SHALL consume the shared canvas styles provided by the reusable site shell
rather than restating them. Project landing pages SHALL be published through
this same site publication and SHALL NOT be published through a per-project
landing deployment.

#### Scenario: Preview a local candidate

- **WHEN** the local `site_serve` target runs
- **THEN** the rendered site is available on the loopback preview address without
  executing the deployment target

#### Scenario: Render the shared canvas

- **WHEN** the apex site is built
- **THEN** the shared canvas style file is a declared input of its assets package
  and the rendered pages use the shared background and foreground palette

#### Scenario: Publish project landing pages

- **WHEN** the site publication workflow runs
- **THEN** the project landing pages are included in the published output
  without a per-project landing deployment step
