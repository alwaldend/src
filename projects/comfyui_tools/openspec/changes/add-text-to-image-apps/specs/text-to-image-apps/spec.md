## Purpose

Provide a focused ComfyUI App interface for each catalog workflow that generates
a single image directly from a text prompt.

## ADDED Requirements

### Requirement: Native text-to-image Apps

The workflow catalog SHALL produce a native ComfyUI App companion for each
text-to-image workflow, preserving that workflow's model and generation defaults.

#### Scenario: Import a text-to-image App

- **WHEN** a user imports a generated `.app.json` companion into ComfyUI
- **THEN** it opens in App mode with editable prompt, canvas dimensions, seed,
  steps, guidance and a saved-image gallery
- **AND** it exposes a negative prompt when the graph uses negative conditioning

#### Scenario: Change canvas dimensions in the FLUX.2 dev App

- **WHEN** a user changes the width or height
- **THEN** the selected value controls both the starting latent and the scheduler
