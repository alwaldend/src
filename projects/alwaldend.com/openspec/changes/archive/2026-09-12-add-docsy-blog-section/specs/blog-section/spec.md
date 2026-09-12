## Purpose

Give alwaldend.com a dated article section that reuses the site's Docsy shell,
packaging, and outputs instead of a separate page type.

## ADDED Requirements

### Requirement: Package an article section

The site SHALL package a `blog` content section whose index page and posts are
declared as Bazel packages and included in the site source archive, so posts
are rendered from declared inputs rather than copied into the build.

#### Scenario: Build the site with a blog post

- **WHEN** the site source archive is built while the blog section contains a
  post package
- **THEN** the rendered site contains the post page under the blog section's
  published path

#### Scenario: Build the blog section index

- **WHEN** the site is built and the blog section contains no posts
- **THEN** the section still renders its index page without failing the build

### Requirement: Publish section outputs

The blog section SHALL publish an RSS feed and a print edition in addition to
its HTML pages, and SHALL be reachable from the site's main navigation.

#### Scenario: Request the blog feed

- **WHEN** a client requests the blog section's feed
- **THEN** the site serves a feed of the section's published posts

#### Scenario: Navigate to the blog section

- **WHEN** a visitor opens any site page that renders the main navigation
- **THEN** the navigation links to the blog section index

### Requirement: Render posts with author and date metadata

Rendered posts SHALL present their publication date and, when the post declares
one, its author, and the section index SHALL order posts by date.

#### Scenario: A post declares a date and author

- **WHEN** a post declares a publication date and an author
- **THEN** the rendered post shows both, and the section index lists the post in
  date order

### Requirement: Withhold unpublished posts

A post whose front matter declares it a draft SHALL be excluded from the
release build while remaining available to local previews, so the draft state
alone keeps a post off the deployed site.

#### Scenario: Release a site containing a draft post

- **WHEN** the site is built for release while the blog section contains a post
  that declares itself a draft
- **THEN** the release output contains no page, feed entry, or sitemap entry for
  that post

#### Scenario: Preview a draft post locally

- **WHEN** the site is built for local preview while the blog section contains a
  post that declares itself a draft
- **THEN** the preview output contains the post so its author can review it
