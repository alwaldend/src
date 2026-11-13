## Purpose

Supply editable native MOSS-Image examples for text generation, layout guidance
and image adapters using the matching deployment profile's local models.

## ADDED Requirements

### Requirement: Six native example workflows

The workflow catalog SHALL generate SFT/Base text-to-image and layout-guided
examples, one edit example and one color/depth example using native nodes and
filenames selected from the model configuration.

#### Scenario: Load generation examples

- **WHEN** a user loads an SFT or Base text-to-image or layout example
- **THEN** the graph exposes editable prompt and sampling controls and saves an image

#### Scenario: Supply an edit reference

- **WHEN** a user selects a reference image and edit instruction
- **THEN** the edit graph uses the edit LoRA and MOSS edit conditioning on both branches

#### Scenario: Supply prepared control maps

- **WHEN** a user selects a region map and depth map at the output size
- **THEN** the control graph uses the control adapter without installing a depth estimator
