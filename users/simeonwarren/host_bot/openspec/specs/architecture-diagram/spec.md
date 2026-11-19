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
`mermaid_svg` rule at `//users/simeonwarren/host_bot:rendered`. The project
SHALL NOT commit a rendered SVG: `//users/simeonwarren/host_bot:docs_assets`
SHALL take its `srcs` from `:rendered` so the documentation site publishes the
build output, and the README SHALL reference that published path. The diagram
SHALL use the shared `//tools/mermaid:theme.json` appearance rather than a
package-local copy.

#### Scenario: Publish the diagram

- **WHEN** the documentation site is built
- **THEN** `:docs_assets` packages the SVG produced by `:rendered`
- **AND** the README image reference resolves to that packaged file
- **AND** no rendered SVG is tracked in the source tree

#### Scenario: Render without host or network inputs

- **WHEN** the render action produces the SVG
- **THEN** it uses the repository's pinned Mermaid, headless Chrome, and font
  inputs
- **AND** it loads no stylesheet, script, or font from a CDN
- **AND** the document embeds the pinned handwriting face, so a viewer that
  loads no webfont still lays the labels out at their measured widths

### Requirement: Diagram content reflects the deployed topology

The diagram SHALL group the remote actor and the shared `Ingress` in an
`Internet` cluster, and SHALL show the remote actor reaching `Ingress` over
HTTPS with mTLS. It SHALL nest the host-bot cluster inside the local
infrastructure cluster, so the local actor, `Vault`, and the host's services
sit inside `dc1.alwaldend.com` rather than beside it, and the two nested cluster
titles SHALL carry the `dc1.alwaldend.com` and
`host-bot.simeonwarren.users.alwaldend.com` names respectively. It SHALL group
`Traefik`, `T3 Code`, and `Codex` in the host-bot cluster, SHALL show the local
actor reaching `Traefik` directly over HTTPS, `Ingress` forwarding to `Traefik`
over Wireguard, and `Traefik` obtaining ACME certificates from `Vault`. It
SHALL show `T3 Code` reaching GitHub for source and `Codex` reaching OpenAI and
OpenRouter for inference, with those providers in a second `Internet` cluster
because they are off-host. Both actors SHALL be labelled `Actor`. External
service nodes SHALL carry their own hostnames.

Sources: [diagram source](../../../t3code-architecture.mmd) and
[project documentation](../../../README.md).

#### Scenario: Consult the host-bot architecture

- **WHEN** a reader opens the architecture section of the project README
- **THEN** it displays the maintained diagram showing the ingress path, the
  in-host services, and the two external service dependencies

### Requirement: Legible diagram layout

The diagram SHALL keep node labels short by naming each host once, in the
cluster title that owns it, and SHALL NOT restate a hostname in both a cluster
title and the nodes it contains. Node labels SHALL NOT wrap mid-token. The
diagram SHALL keep subgraph titles clear of node outlines and edges: Mermaid
centers a cluster title on the cluster's top border, so a title that wraps
extends into the cluster's node area. The diagram source SHALL describe topology
only and SHALL let the shared theme own appearance and spacing; a per-diagram
directive is warranted only for a specific geometry need the theme cannot
express.

Actors SHALL render as an ordinary rounded node carrying the shared theme's
hand-drawn stick figure beside the label, rather than as a bare glyph, a pill
outline, or a Mermaid `person` shape. Mermaid's `person` shape draws the head
over the body box, so it is not used.

Sources: [diagram source](../../../t3code-architecture.mmd) and
[theme](../../../../../../tools/mermaid/theme.json).

#### Scenario: Render the maintained layout

- **WHEN** the renderer draws the diagram
- **THEN** every node outline is clear of headings, edges, and other nodes
- **AND** no label is truncated or broken inside a hostname

### Requirement: Appearance is owned by the shared theme

Diagram appearance SHALL come from `//tools/mermaid:theme.json` rather than
package-local copies or per-diagram appearance overrides, so node rounding and
edge-label plates are consistent across every rendered diagram. The theme SHALL
keep shapes semantic, so a rounded bracket node renders rounded while a square
bracket node stays square; rounding SHALL NOT be forced on every node. The
theme SHALL box edge labels by styling the label text, because Mermaid emits an
empty label plate for unlabelled edges and styling that plate directly would
draw a stray box.

Source: [theme](../../../../../../tools/mermaid/theme.json).

#### Scenario: Render another diagram with the theme

- **WHEN** any diagram renders through `mermaid_svg` or the `mmdc` target
- **THEN** a rounded node renders rounded and a square node stays square
- **AND** its edge labels use bordered plates

#### Scenario: Read an edge annotation

- **WHEN** the renderer draws an edge that carries a label
- **THEN** the label renders on a bordered plate rather than directly on the
  edge line
- **AND** an edge without a label draws no plate
