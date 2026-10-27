## ADDED Requirements

### Requirement: Site tools use execution configuration

`al_hugo_site` SHALL configure its tools for execution while retaining site
source archives in the target configuration.

#### Scenario: Query build inputs without tools

- **WHEN** a site's dependency graph is queried with `--notool_deps`
- **THEN** the site archive remains reachable
- **AND** dependencies reached only through the tools attribute are excluded
