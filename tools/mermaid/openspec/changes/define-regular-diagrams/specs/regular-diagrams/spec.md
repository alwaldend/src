## ADDED Requirements

### Requirement: Regular diagram authoring

The Mermaid skill SHALL define a regular diagram as a conventional
technical diagram and state the repository's conventions directly.

#### Scenario: Use regular flowchart styling

- **WHEN** the user asks for a regular Mermaid flowchart
- **THEN** infrastructure diagrams use the shared site renderer and theme, curved connectors, and light and dark variants
- **AND** the definition specifies left-to-right layout, semantic shapes, and concise labels

### Requirement: Shared diagram declaration

The repository SHALL provide a `mermaid_diagram` macro as the standard
declaration for maintained non-blog Mermaid diagrams. It SHALL produce
light and dark SVG and WebP variants from one source and theme.

#### Scenario: Build the complete diagram

- **WHEN** a consumer builds the aggregate target of a `mermaid_diagram` declaration
- **THEN** it produces `.svg`, `.dark.svg`, `.webp`, and `.dark.webp` files with the same output stem
- **AND** consumers can select the SVG pair or an individual WebP variant

#### Scenario: Show a diagram in conversation

- **WHEN** an author shows a maintained diagram in conversation without a requested format or color scheme
- **THEN** the skill uses the macro's dark WebP output for the inline preview.

#### Scenario: Publish Host Bot documentation diagrams

- **WHEN** the documentation build packages the T3 Code and ComfyUI diagrams
- **THEN** it consumes the SVG pair targets emitted by `mermaid_diagram` directly

#### Scenario: Preserve existing publication workflows

- **WHEN** a published blog image uses an existing render declaration
- **THEN** it retains its render target and generated publication asset
