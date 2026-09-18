# architecture-diagram Specification

## Purpose

Describe the maintained architecture diagram owned by
`users/simeonwarren/host_bot`. The diagram replaced the `T3code` page of the
`infra/arch` Drawio document, which no longer renders or publishes that page.
The diagram describes checked-in configuration; it is not a live inventory or
health check.

## Requirements

### Requirement: Owner-local Mermaid diagram source

The project SHALL keep its architecture diagram as Mermaid source at
`t3code-architecture.mmd` and SHALL render it through the repository's
`mermaid_svg` rule at `//users/simeonwarren/host_bot:rendered`. The maintained
SVG SHALL live at `assets/t3code-architecture.svg` and SHALL be refreshed by
`//users/simeonwarren/host_bot:update`. The diagram SHALL use the shared
`//tools/mermaid:theme.json` appearance rather than a package-local copy, and
documentation SHALL display the maintained SVG with a link to its full-size
image.

#### Scenario: Refresh the maintained diagram

- **WHEN** the diagram source changes without the maintained SVG being refreshed
- **THEN** `//users/simeonwarren/host_bot:update_tests` reports a mismatch
- **AND** the update target regenerates the maintained SVG

#### Scenario: Render without host or network inputs

- **WHEN** the render action produces the SVG
- **THEN** it uses the repository's pinned Mermaid, headless Chrome, and font
  inputs
- **AND** it loads no stylesheet, script, or font from a CDN

### Requirement: Diagram content reflects the deployed topology

The diagram SHALL show the remote actor reaching the shared `Ingress` over
HTTPS with mTLS, the intranet actor reaching `Traefik` directly over HTTPS, and
`Ingress` forwarding to `Traefik` over Wireguard. It SHALL group `Traefik`,
`T3 Code`, and `Codex` inside the host-bot host in the `dc1.alwaldend.com` local
infrastructure, SHALL show `Traefik` obtaining ACME certificates from `Vault`,
and SHALL show `T3 Code` reaching GitHub for source and `Codex` reaching OpenAI
for inference.

Sources: [diagram source](../../../t3code-architecture.mmd) and
[project documentation](../../../README.md).

#### Scenario: Consult the host-bot architecture

- **WHEN** a reader opens the architecture section of the project README
- **THEN** it displays the maintained diagram showing the ingress path, the
  in-host services, and the two external service dependencies
