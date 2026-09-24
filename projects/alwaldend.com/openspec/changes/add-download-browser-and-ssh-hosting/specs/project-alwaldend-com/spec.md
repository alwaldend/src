## MODIFIED Requirements

### Requirement: Separate preview and publication

`site_serve` SHALL serve built content on `127.0.0.1:1313`. Production site
publication SHALL package the release-mode output as a website archive and
explicitly deploy it through the release tool's SSH deployment to the selected
environment. The archive SHALL contain `index.html` at its root and the assets
needed by the complete rendered site, including its packaged project landing
pages. Building or previewing the site SHALL NOT itself publish it. The
rendered site SHALL consume the shared canvas styles provided
by the reusable site shell rather than restating them.

#### Scenario: Preview a local candidate

- **WHEN** the local `site_serve` target runs
- **THEN** the rendered site is available on the loopback preview address without
  executing the deployment target

#### Scenario: Render the shared canvas

- **WHEN** the apex site is built
- **THEN** the shared canvas style file is a declared input of its assets package
  and the rendered pages use the shared background and foreground palette

#### Scenario: Explicitly publish the main website

- **WHEN** an authorized operator selects a concrete release and SSH target
- **THEN** the archive is published and its extracted release becomes selected
- **AND** the public site URL remains `https://alwaldend.com`
- **AND** the production workflow does not require publishing a GitHub Pages branch

#### Scenario: Publish project landing pages

- **WHEN** the SSH site publication workflow runs
- **THEN** the project landing pages are included in the selected website release
- **AND** no per-project landing deployment step is needed
