# Hugo Landing Specification

## Purpose

Describe the reusable Hugo landing site that transforms a project's README into
a project page using a shared visual shell.

Sources: [project README](../../../README.md),
[Bazel package](../../../BUILD.bazel),
[landing macros](../../../pkg/bzl/al_hugo_landing.bzl),
[link rendering](../../../layouts/_markup/render-link.html),
and [image rendering](../../../layouts/_markup/render-image.html).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08.

## Requirements

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

### Requirement: Select and adapt project README content

README selection SHALL prefer candidates associated with the requested project,
choose the shortest candidate path, and fail if no README exists. The selected
content SHALL become `content/_index.md` with full-documentation Hugo shortcodes
removed by the configured transformation.

#### Scenario: No README is provided by the documentation target

- **WHEN** the declared documentation files contain no `README.md`
- **THEN** landing-source analysis fails with a missing-README diagnostic.

### Requirement: Resolve repository-relative resources

Relative Markdown links SHALL resolve to their project source paths on the
configured GitHub repository's `master` branch, preserving query and fragment
parts. Relative image paths SHALL resolve to corresponding raw source URLs.
Absolute URLs SHALL retain their original destinations.

#### Scenario: Render a link to a project-local document

- **WHEN** a README links to `docs/README.md#usage`
- **THEN** the landing page links to that project's source document on GitHub
  with the `usage` fragment retained.

### Requirement: Keep project publication explicit

Rendered landing targets SHALL be owned by their projects, with root-workspace
landing targets for nested rules modules. The registered deployment workflow
SHALL preserve Pages history, skip unchanged sites, and support a dry run that
stages and compares without pushing.

#### Scenario: A registered landing site is unchanged

- **WHEN** deployment compares rendered content with the existing Pages branch
  and finds no differences
- **THEN** it creates no publication commit for that site.
