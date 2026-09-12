# Alwaldend.com Specification

## Purpose

Describe the repository's main Hugo and Docsy documentation website and its
packaged build, local preview, and explicit GitHub Pages deployment behavior.

Sources: [project README](../../../README.md),
[site targets](../../../BUILD.bazel), and
[Hugo dependencies](../../../include.MODULE.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08. This specification describes source behavior, not a live
deployment observation.

## Requirements

### Requirement: Build packaged repository documentation

The `site` target SHALL render the packaged repository documentation tree with
the configured Hugo, Docsy, Bootstrap, and local layouts and assets. Packaged
`README.md` documentation SHALL be renamed to Hugo `_index.md` pages.

#### Scenario: Build the main documentation site

- **WHEN** the `//projects/alwaldend.com:site` target is built
- **THEN** it produces rendered site output from declared documentation and
  asset inputs, with the release configuration using `https://alwaldend.com`.

### Requirement: Resolve links relative to documentation sources

The site SHALL resolve Markdown links and images relative to their source
directory, map packaged documentation pages and resources to published URLs,
and preserve unknown internal destinations rather than silently redirecting
them to GitHub.

#### Scenario: Documentation contains an unresolved internal link

- **WHEN** a Markdown destination has no matching packaged page or resource
- **THEN** it remains unresolved in the rendered output rather than being
  silently redirected to GitHub.

### Requirement: Preserve distinct anchors in combined print pages

Print rendering SHALL scope document IDs and their fragment and control
references to the source document so that combined documents retain distinct
anchors.

#### Scenario: Two printed documents contain the same heading anchor

- **WHEN** the documents are combined into a print page
- **THEN** their scoped IDs and references remain distinguishable.

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
