## Purpose

Let readers inspect existing Mermaid diagrams as themed three-dimensional
scenes with spatial connections and interactive navigation.

## ADDED Requirements

### Requirement: Spatial representations preserve source meaning

The renderer SHALL produce interactive 3D views of the existing flowchart,
state, class, entity-relationship, and mindmap examples from their
authoritative Mermaid definitions. Nodes SHALL have volume, available groups
SHALL have three-dimensional bounds, and connections SHALL follow spatial
paths. Labels and relationship direction or cardinality SHALL remain available.

#### Scenario: Inspect the sample collection

- **WHEN** the sample collection is built and opened
- **THEN** each of the five supported examples has an interactive 3D view
- **AND** its source labels and connections remain represented
- **AND** groups contain their members and edge paths clear unrelated nodes

#### Scenario: Follow curved connections

- **WHEN** graph connections or self-loops change direction
- **THEN** their paths bend smoothly while clearing unrelated nodes
- **AND** endpoint markers follow the path direction and dashed connections retain visible gaps through bends

#### Scenario: Preserve adapter display and layout semantics

- **WHEN** supported graphs use encoded text, generic class headings, multiline decisions, invisible links, or local subgraph directions
- **THEN** displayed text is decoded before measurement and generic type parameters remain visible
- **AND** decision labels fit inside their diamond faces
- **AND** invisible links influence placement without appearing as selectable relationships
- **AND** groups retain their local direction, with nested groups inheriting the enclosing direction when none is declared

#### Scenario: Read a self-loop

- **WHEN** a supported graph connects a node to itself
- **THEN** the curved route remains outside that node's face in the front view
- **AND** its two ports remain distinct

#### Scenario: Read a long group title

- **WHEN** a graph group has a long or multiline title
- **THEN** its measured title contributes to the group dimensions and retains line breaks
- **AND** desktop and phone framing include the title

#### Scenario: Connect a decision node

- **WHEN** multiple connections share a side of a diamond decision node
- **THEN** each port follows the diamond surface with only its endpoint clearance
- **AND** horizontal, vertical, incoming, outgoing, and self-loop connections remain attached

### Requirement: Supported Mermaid integration

The renderer SHALL consume graph data through Mermaid's public layout loader
and render APIs. Flowchart, state, class, entity-relationship, and mindmap
are supported; sequence diagrams are outside the 3D renderer's scope.

#### Scenario: Obtain graph data

- **WHEN** a supported source is rendered
- **THEN** the registered layout receives nodes, groups, and edges from Mermaid
- **AND** no internal Mermaid database API is required

#### Scenario: Reject an unsupported diagram

- **WHEN** a sequence diagram is supplied to the 3D renderer
- **THEN** rendering fails with an explicit unsupported-diagram error
- **AND** no HTML viewer is produced

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

#### Scenario: Render a large diagram on a high-density display

- **WHEN** the fifty-node flowchart opens with device pixel ratio 2
- **THEN** the drawing buffer fits the WebGL renderbuffer and viewport limits
- **AND** the complete geometry remains visible within the scrollable CSS canvas
- **AND** switching to a small overview restores the native pixel density

### Requirement: Reproducible verification

The build SHALL produce self-contained HTML views. Browser verification SHALL
exercise the built collection and preserve screenshots and a machine-readable
verification result as repeatable artifacts. Existing 2D rendering SHALL
continue to satisfy its maintained tests alongside the additional renderer.

#### Scenario: Verify a candidate

- **WHEN** the browser verification runs against the built collection
- **THEN** it verifies scene geometry, source coverage, interaction, and theme
- **AND** it emits screenshots and a result identifying the checked artifacts
- **AND** the maintained views stay within measured offline payload and draw-call budgets
