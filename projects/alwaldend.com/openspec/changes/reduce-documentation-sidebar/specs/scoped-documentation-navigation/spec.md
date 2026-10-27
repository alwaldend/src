## ADDED Requirements

### Requirement: Scope documentation navigation by owner

Documentation owner sections SHALL display their own subtree instead of the
complete repository tree. Their sidebar SHALL provide navigation to the parent
section and preserve current-page selection and nested branch navigation.

#### Scenario: Visit a project or tool document

- **WHEN** a reader visits a section beneath /docs/<category>/<owner>
- **THEN** the sidebar contains that owner's descendants and a parent link
- **AND** unrelated owner subtrees are absent from that sidebar's HTML

#### Scenario: Navigate on a narrow screen

- **WHEN** a reader opens documentation on mobile
- **THEN** the navigation toggle reveals the scoped tree and its parent link

### Requirement: Search does not block navigation

The site SHALL precompute its complete local search index during the build and
load and query it outside the browser UI thread. It SHALL preserve searchable
content, ranking weights, result limits, page exclusions, and result metadata.

#### Scenario: Navigate after a reload

- **WHEN** a visitor hovers or taps a homepage or sidebar link after reloading
- **THEN** search indexing does not delay its feedback or navigation

#### Scenario: Search while the index loads

- **WHEN** a visitor submits a query before the index is ready
- **THEN** the site shows loading state and returns the requested results when ready
- **AND** navigation and typing remain responsive

#### Scenario: Retry an unavailable index

- **WHEN** loading fails and the visitor retries after connectivity recovers
- **THEN** a clear failure is replaced by results without reloading the page

#### Scenario: Preserve the latest query

- **WHEN** multiple searches complete out of order or a query is cleared
- **THEN** stale responses do not replace current results or reopen a cleared search
