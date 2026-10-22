## Purpose

Give every registered project a visitor-facing home inside the main site, with
landing content owned by the project and published at `/projects/<name>/`, so
project landings no longer require a dedicated Pages repository, custom domain,
or DNS record per project.

## ADDED Requirements

### Requirement: Publish a project section in the main site

The main site SHALL publish a `/projects` section. The section index SHALL list
the projects registered for landing publication, and each registered project
SHALL have a landing page at `/projects/<name>/`. A project without landing
content SHALL NOT fail the site build and SHALL NOT produce a landing page.

#### Scenario: Render a registered project landing page

- **WHEN** the site is built while a registered project declares landing content
- **THEN** the rendered site contains the landing page at that project's
  `/projects/<name>/` path

#### Scenario: Build the section index

- **WHEN** the site is built and the project section contains no landing pages
- **THEN** the section still renders its index page without failing the build

#### Scenario: A registered project has no landing content yet

- **WHEN** a project is registered for landing publication but owns no landing
  content
- **THEN** the build succeeds and the project has no landing page in the output

### Requirement: Package project-owned landing content

Landing content SHALL be owned by each project and packaged into the main site
from that project's content package, using one registry as the single owner of
landing membership. Project landing packages SHALL contain content only and
SHALL NOT declare layouts, styles, site configuration, or build rules.

#### Scenario: Add a project to the registry

- **WHEN** a project is added to the landing registry and declares landing
  content
- **THEN** the site's packaged sources include that content without a second
  membership list being edited

#### Scenario: A project landing package declares presentation

- **WHEN** a project landing content package declares a layout, style, or site
  configuration of its own
- **THEN** that declaration is a defect because presentation is owned by the
  shared shell

### Requirement: Serve user-facing project documentation beside the landing page

A project that owns user-facing documentation SHALL publish it under its
landing path at `/projects/<name>/docs/`, including subpages. Visitor-facing
pages SHALL be classified by declared page type rather than by directory
position, so that a landing page and documentation pages under one project path
render with their respective layouts.

#### Scenario: Render project documentation with subpages

- **WHEN** a project declares user-facing documentation containing subpages
- **THEN** the rendered site serves those pages beneath the project's
  `/projects/<name>/docs/` path

#### Scenario: Classify landing and documentation pages

- **WHEN** a landing page and a documentation page are rendered under the same
  project path
- **THEN** each renders with the layout for its declared page type

### Requirement: Keep repository documentation separate from landing pages

Project landing content SHALL be distinct from the project's repository
documentation. A project's `README.md` SHALL remain packaged by the repository
documentation tree at its existing `/docs/projects/<name>/` path, and project
landing content SHALL NOT be packaged into the repository documentation tree.

#### Scenario: Resolve the project's documentation and landing URLs

- **WHEN** a visitor requests a project's repository documentation path and its
  landing path
- **THEN** each serves content owned by the corresponding source

#### Scenario: Package a landing content directory

- **WHEN** the site is built while a project owns landing content
- **THEN** that content is not published under the repository documentation path

### Requirement: Publish landing pages without dedicated hosting

Registered project landings SHALL be published as part of the main site, and
SHALL NOT require a per-project Pages repository, custom domain, `CNAME`
record, or per-project DNS stage. Landing pages SHALL NOT be published at
per-project hostnames.

#### Scenario: Publish a landing page

- **WHEN** the site is published
- **THEN** the landing page is served from the main site's hostname and requires
  no per-project repository or DNS record

#### Scenario: Inspect a retired project hostname

- **WHEN** a visitor requests a project's former dedicated hostname after
  retirement
- **THEN** no dedicated landing site is published there

### Requirement: Participate in site taxonomies

Landing pages SHALL participate in the site's configured taxonomies, so
projects are discoverable by their declared statuses, languages, and tags.

#### Scenario: List projects by taxonomy term

- **WHEN** a landing page declares a taxonomy term that the site configures
- **THEN** the corresponding taxonomy term page lists that project's landing page

### Requirement: Own the shared landing presentation in the main site

The main site SHALL own the presentation the project landings use, including
the shared canvas and accent styles and the landing page rules. A separate
reusable landing project SHALL NOT be required for the main site to render
landings, and each shared style SHALL have exactly one owner. The site SHALL
use the theme's taxonomy and term layouts rather than a site-local copy.

#### Scenario: Render a landing page after the reusable project is removed

- **WHEN** the site is built after the reusable landing project is removed
- **THEN** the landing pages still render with the shared canvas, accent, and
  landing page rules

#### Scenario: Locate the owner of a shared style

- **WHEN** a shared canvas or accent style is changed
- **THEN** it is changed once in the main site's own tree rather than in a
  second location

### Requirement: Retire per-project DNS ownership

Project landings SHALL NOT own DNS records, dedicated hostnames, or
per-project Terraform DNS stages. A project whose Terraform root declares only
DNS SHALL have that root removed entirely, and a project root that declares
other resources SHALL retain them. Shared apex and mail declarations SHALL
retain their existing owner, and Vault configuration SHALL NOT change.

#### Scenario: Inspect DNS ownership after retirement

- **WHEN** the declarations are inspected after retirement
- **THEN** no project landing owns a hostname, DNS record, or DNS stage, and the
  shared apex declarations are unchanged

#### Scenario: Lint declarations after retirement

- **WHEN** the DNS declaration lint runs after retirement
- **THEN** it passes with no landing hostname declared by any source
