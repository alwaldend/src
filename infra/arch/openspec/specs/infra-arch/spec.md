# Infrastructure architecture Specification

## Purpose

Describe the maintained infrastructure diagrams owned by `infra/arch`.
The baseline is checked-in source at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08.
These diagrams describe their source document; they do not establish live
inventory, deployment, or health.

## Requirements

### Requirement: Canonical multi-page architecture source

The project SHALL use `arch.drawio` as the canonical source for its maintained
SVG diagrams and SHALL expose the source through the `//infra/arch` editor
target. The named render mapping SHALL preserve all 15 source pages, including
the five pages classified as archives.

Sources: [project documentation](../../../README.md),
[target and page definitions](../../../BUILD.bazel), and
[Drawio source](../../../arch.drawio).

#### Scenario: Render the maintained page set

- **WHEN** the `//infra/arch:rendered` target processes the canonical document
- **THEN** it produces the SVG output mapped to each of the 15 named pages,
  including separate outputs for archived pages.

### Requirement: Reproducible diagram refresh

The project SHALL render diagrams through the repository's pinned Drawio,
headless Chrome, and font inputs. The `//infra/arch:update` target SHALL update
the maintained SVGs, and `//infra/arch:update_tests` SHALL compare them with
freshly rendered output. Rendering SHALL use local assets and report missing
pages or export failures.

Sources: [rendering guarantees](../../../README.md) and
[render and update targets](../../../BUILD.bazel).

#### Scenario: A maintained diagram becomes stale

- **WHEN** a source edit changes a rendered SVG without updating its maintained copy
- **THEN** the freshness test reports a mismatch, and the update target provides
  the regenerated source file.

### Requirement: Documentation preserves historical classification

Documentation SHALL display the maintained SVGs with links to their full-size
images and SHALL retain the source's `Archive` classification. The diagrams
SHALL remain documentation artifacts rather than evidence of live service
availability.

Source: [diagram documentation](../../../README.md).

#### Scenario: Consult a historical architecture page

- **WHEN** a reader opens the Proxmox architecture section
- **THEN** it is labeled `Archive/Proxmox` and links to the maintained SVG.
