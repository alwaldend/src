## Purpose

Give visitors a discoverable browser workspace for editing Anki configurations,
applying them to collections and saving completed work within the main site's
established presentation.

## ADDED Requirements

### Requirement: Discoverable application page

The project landing SHALL provide a Try it now link to `/projects/anki_as_code/app/`. The app page SHALL reuse the site's navigation, typography, light/dark modes and responsive content layout.

#### Scenario: Open the app from the landing

- **WHEN** a visitor activates Try it now
- **THEN** the application page opens within the same site and retains its shared presentation

#### Scenario: Small screens and theme preferences

- **WHEN** a visitor changes the site's color mode or uses a narrow viewport
- **THEN** the application controls follow that mode and remain usable without horizontal page overflow

### Requirement: Configuration editing workflow

The app page SHALL restore its last saved workspace on startup, or load an
editable example when no saved workspace exists, and present one TOML editor
with a file selector. Visitors SHALL be able to open configuration JSON, check
their edits and download the complete configuration as JSON. Opening JSON SHALL
load its files directly. Switching files, checking and downloading SHALL retain
the current editor contents without a separate action to update in-memory
files. Operations SHALL use the existing browser App API and display recoverable
errors.

#### Scenario: Begin editing immediately

- **WHEN** a visitor opens the application without a saved workspace and Go becomes ready
- **THEN** the example's note is selected and editable without a generation step

#### Scenario: Edit and download

- **WHEN** a visitor edits a note's TOML, switches files and downloads the configuration
- **THEN** the edit is retained and included in the JSON download

#### Scenario: Reopen downloaded work

- **WHEN** a visitor opens a previously downloaded configuration JSON
- **THEN** its files are immediately available in the editor

#### Scenario: Import without browser storage

- **WHEN** IndexedDB is unavailable or a saved snapshot cannot be read
- **THEN** visitors can still open configuration JSON and collection packages, edit them and download completed work
- **AND** local save failures remain visible without discarding those edits

#### Scenario: Import after another tab saves

- **WHEN** another tab saves a newer revision while this page opens a local file
- **THEN** the local import succeeds but a subsequent save rejects the stale revision until the visitor explicitly restores the newer snapshot

#### Scenario: Recover from malformed input

- **WHEN** an operation receives invalid JSON or TOML
- **THEN** the page reports the error, retains the existing workspace and editor text, and remains usable after the visitor corrects the input

#### Scenario: Reject configuration paths outside the editor workspace

- **WHEN** imported JSON writes files outside the configuration directory, including reserved `/tmp` storage
- **THEN** the page rejects the candidate without changing the current editor or collection
- **AND** the import cannot create a collection in a configuration-only session or replace an open collection, including after local save and reload

### Requirement: Clear storage lifetime

The app page SHALL distinguish active in-memory work from completed local saves
and downloaded files. It SHALL explain that unsaved work lasts for the tab
session and that browser-local storage can be cleared. Portable configuration
and collection downloads SHALL remain available.

#### Scenario: Open a fresh session

- **WHEN** a visitor opens the application
- **THEN** the page describes how to retain work and exposes configuration editing without requiring collection backend startup

### Requirement: Collection workflow controls

The page SHALL provide Open collection, Plan, Apply and Download collection
actions through the existing browser App API. Collection imports SHALL populate
the existing TOML editor. Planning SHALL display the application's change report;
application SHALL use its established reconciliation semantics. Current editor
text SHALL be retained before these operations. Check SHALL continue to describe
configuration format validity independently of reconciliation.

#### Scenario: Edit and apply an imported collection

- **WHEN** a visitor opens a collection package, edits its configuration, inspects a plan and applies it
- **THEN** the page displays the result and allows downloading the completed collection with its media retained

#### Scenario: Failed collection import

- **WHEN** a visitor opens an invalid collection package
- **THEN** the page reports the failure and retains the previous workspace and editor contents

### Requirement: Explicit local save and startup restore

The page SHALL provide Save locally and restore the last saved workspace on
startup using the SDK's atomic snapshot API. Saving SHALL include current
editable config bytes and, when present, the last completed collection package
without implicitly applying edits. Configuration-only work SHALL also be
saveable without loading the server. It SHALL report
success only after persistence completes. Restore SHALL replace the current
workspace only after successfully preparing the complete saved snapshot.
Errors and conflicting saves SHALL remain visible and preserve current edits.

#### Scenario: Continue after reload

- **WHEN** a visitor saves locally and reloads the page
- **THEN** the editor and collection recover the completed saved snapshot

#### Scenario: Failed local save

- **WHEN** browser storage rejects or conflicts with a save
- **THEN** the page reports the failure, keeps current editor contents and leaves portable downloads available

#### Scenario: Failed startup restore

- **WHEN** browser storage cannot restore the last saved workspace
- **THEN** the page reports the error, preserves the saved snapshot and leaves configuration editing available
