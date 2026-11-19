## MODIFIED Requirements

### Requirement: Separate preview and publication

`site_serve` SHALL serve built content on `127.0.0.1:1313`. The deployment
archive SHALL contain the apex `CNAME`, and the explicit deployment workflow
SHALL write `.nojekyll` and push changed output to the configured Pages branch.
Building or previewing the site SHALL NOT itself publish it. The rendered site
SHALL consume the shared canvas styles provided by the reusable site shell
rather than restating them.

#### Scenario: Preview a local candidate

- **WHEN** the local `site_serve` target runs
- **THEN** the rendered site is available on the loopback preview address without
  executing the deployment target

#### Scenario: Render the shared canvas

- **WHEN** the apex site is built
- **THEN** the shared canvas style file is a declared input of its assets package
  and the rendered pages use the shared background and foreground palette
