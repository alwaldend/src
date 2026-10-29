## Purpose

Let readers inspect existing Mermaid diagrams as themed three-dimensional
scenes with spatial connections and interactive navigation.

## ADDED Requirements

### Requirement: Spatial diagrams provide static snapshots

Generated spatial views SHALL support hermetic PNG and WebP snapshots for
documentation embeds. Snapshots SHALL render the same scene, theme, and font
as the interactive view, without viewer controls, and SHALL include all content
in the selected camera view. Light and dark appearances SHALL be selectable.

#### Scenario: Embed a spatial diagram

- **WHEN** a generated spatial HTML view is built as a static image
- **THEN** its scene is rendered without remote requests or interactive controls
- **AND** all scene geometry and labels fit within the image
- **AND** the chosen light/dark appearance and camera view are preserved
- **AND** the original interactive HTML remains available

### Requirement: Spatial representations preserve source meaning

The renderer SHALL produce interactive 3D views of the existing flowchart,
sequence, state, class, entity-relationship, and mindmap examples from their
authoritative Mermaid definitions. Nodes SHALL have volume, available groups
SHALL have three-dimensional bounds, and connections SHALL follow spatial
paths. Labels and relationship direction or cardinality SHALL remain available.

#### Scenario: Inspect the sample collection

- **WHEN** the sample collection is built and opened
- **THEN** each of the six existing examples has an interactive 3D view
- **AND** its source labels and connections remain represented
- **AND** groups contain their members and edge paths clear unrelated nodes

#### Scenario: Inspect a sequence

- **WHEN** the sequence example is opened
- **THEN** participants and messages retain their source order
- **AND** replies, activation, self-calls, loop/alternative blocks, and notes remain represented in three dimensions

#### Scenario: Follow curved connections

- **WHEN** graph connections or routed self-calls change direction
- **THEN** their paths bend smoothly while clearing unrelated nodes
- **AND** endpoint markers follow the path direction and dashed connections retain visible gaps through bends
- **AND** straight sequence messages retain their time alignment

#### Scenario: Preserve sequence structure

- **WHEN** a sequence uses numeric participant IDs, arrowless messages, or notes inside interaction blocks
- **THEN** participant declaration order and the absence of message arrowheads are retained
- **AND** each block contains its notes and nested blocks, including blocks containing only notes

#### Scenario: Multiline sequence labels

- **WHEN** participants or notes contain line breaks or long labels
- **THEN** their card dimensions and separation accommodate the rendered text
- **AND** subsequent messages and containing blocks retain adequate space

#### Scenario: Preserve adapter display and layout semantics

- **WHEN** supported graphs use encoded text, generic class headings, multiline decisions, invisible links, or local subgraph directions
- **THEN** displayed text is decoded before measurement and generic type parameters remain visible
- **AND** decision labels fit inside their diamond faces
- **AND** invisible links influence placement without appearing as selectable relationships
- **AND** groups retain their local direction, with nested groups inheriting the enclosing direction when none is declared

#### Scenario: Preserve activation intervals

- **WHEN** sequence messages start or end explicit or shorthand activations, including nested self-calls
- **THEN** each activation interval aligns with its initiating and terminating message events
- **AND** nested activations remain visually distinct

#### Scenario: Read a self-loop

- **WHEN** a supported graph connects a node to itself
- **THEN** the curved route remains outside that node's face in the front view
- **AND** its two ports remain distinct

#### Scenario: Read a long group title

- **WHEN** a graph group has a long or multiline title
- **THEN** its measured title contributes to the group dimensions and retains line breaks
- **AND** desktop and phone framing include the title

### Requirement: Spatial views retain the maintained theme

The views SHALL derive colors and typography from the same theme, palette,
and pinned font inputs as their existing site diagrams, including light and
dark appearances. Rendering and viewing SHALL require no remote assets.

#### Scenario: Change appearance offline

- **WHEN** a generated view is opened without network access and its appearance changes
- **THEN** the scene uses the selected maintained palette and embedded font
- **AND** its geometry and labels remain available

### Requirement: Interactive inspection

Readers SHALL be able to rotate, pan, zoom, reset the camera, and select a
diagram element to inspect its identity and connected paths. Selection SHALL
remain available through keyboard-operable controls.

#### Scenario: Trace a connection

- **WHEN** a reader selects an element and changes the camera
- **THEN** the element and its connections can be distinguished
- **AND** resetting the view restores a useful framing

#### Scenario: Limit selection to the requested relationships

- **WHEN** a reader selects a connection
- **THEN** only that connection and its endpoints remain highlighted
- **WHEN** a reader selects a node or group
- **THEN** its direct connections and group descendants remain highlighted without following further graph connections

#### Scenario: Select elements with shared identifiers

- **WHEN** a node and a connection use the same Mermaid identifier
- **THEN** each remains independently selectable with its own details and relationship highlighting

#### Scenario: Read and resize a diagram

- **WHEN** a sample is opened or the browser changes between desktop and phone dimensions
- **THEN** the default view preserves readable text and exposes any larger canvas through scrolling
- **AND** a fitted overview frames the complete scene at the current dimensions
- **AND** relationship labels avoid covering node text in the maintained default views
- **AND** fitted overview panels avoid unnecessary vertical space for wide scenes

#### Scenario: Inspect a large flowchart

- **WHEN** a supported flowchart has fifty sequential nodes
- **THEN** its geometry and labels remain visible within the camera's depth range
- **AND** opening, orbiting, resetting, changing view or appearance, and resizing preserve visible canvas content

### Requirement: Reproducible verification

The build SHALL produce self-contained HTML views. Browser verification SHALL
exercise the built collection and preserve screenshots and a machine-readable
verification result as repeatable artifacts. Existing 2D rendering SHALL
continue to satisfy its maintained tests after the dependency upgrade.

#### Scenario: Verify a candidate

- **WHEN** the browser verification runs against the built collection
- **THEN** it verifies scene geometry, source coverage, interaction, and theme
- **AND** it emits screenshots and a result identifying the checked artifacts
- **AND** the maintained views stay within measured offline payload and draw-call budgets
