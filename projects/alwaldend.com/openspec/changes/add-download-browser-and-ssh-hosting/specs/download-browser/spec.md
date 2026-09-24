## Purpose

Let visitors browse live public release files in the main site's shared
presentation and reuse the same listing behavior on individual release pages.

## ADDED Requirements

### Requirement: Downloads page in the shared website

The site SHALL publish `/downloads/` and link it from its header. The page
SHALL use the existing navigation, theme, and responsive presentation while
rendering live JSON listings from the configured download endpoint.

#### Scenario: Open Downloads in either theme

- **WHEN** a visitor opens Downloads from the header in light or dark mode
- **THEN** the page uses the same presentation as the rest of the website
- **AND** it displays the current public project listing

### Requirement: Navigable directories and direct file links

The browser SHALL display entry names, types, dates, and available file sizes.
Directory navigation SHALL retain a shareable page URL and support refresh
and browser history. File links SHALL point to the corresponding Nginx file
endpoint. Navigation SHALL remain within the configured browse root.

#### Scenario: Share a release directory

- **WHEN** a visitor navigates to a release and shares the resulting page URL
- **THEN** another visitor opens the same directory on `/downloads/`
- **AND** selecting an artifact requests that file directly from the download host

### Requirement: Reusable independently scoped instances

One shared listing component SHALL serve the Downloads page and release
pages. An instance SHALL accept an endpoint, starting directory, and browse
boundary. Release pages SHALL derive their project/version from canonical
release metadata and SHALL NOT navigate above that release's root.

#### Scenario: Render two embedded releases

- **WHEN** two release sections contain listing instances
- **THEN** each shows files for its own project/version using the same component
- **AND** interaction with one does not change the other's state

### Requirement: Safe and resilient client rendering

The component SHALL render filenames as text and encode URL path segments.
It SHALL make public requests without credentials and show distinct loading,
empty, missing-directory, and failed-request states. An unavailable download
service SHALL NOT prevent the rest of the page from rendering.

#### Scenario: Render an unusual filename

- **WHEN** an entry contains spaces, Unicode, or HTML-like characters
- **THEN** its name is displayed literally and its link requests the correct file

#### Scenario: Download service is unavailable

- **WHEN** a listing request fails
- **THEN** the instance displays a retryable error while surrounding page content remains usable

### Requirement: Offline website build

Building the website SHALL NOT fetch live directory listings or depend on
infrastructure deployment artifacts. Listing data SHALL be requested only by
the browser. The page SHALL provide a direct endpoint link when JavaScript is
unavailable.

#### Scenario: Build with no download server

- **WHEN** the site is built with the download service unreachable
- **THEN** its pages and reusable listing assets are produced successfully
