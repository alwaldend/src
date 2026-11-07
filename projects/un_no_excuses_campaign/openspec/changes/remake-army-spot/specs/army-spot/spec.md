## Purpose

Provide an editable Blender remake of the supplied army comedy clip and a
repeatable way to produce its animated render in the campaign project.

## ADDED Requirements

### Requirement: Editable animated scene

The asset SHALL contain four blue-skinned stylized characters with an officer,
three helmeted recruits, a sandy base, barracks, tanks, and animated shot cameras.

#### Scenario: Open the scene

- **WHEN** the delivered Blender file is opened with automatic scripts disabled
- **THEN** its modeled subjects, materials, cameras, and timeline animation remain editable without the source video

### Requirement: Complete playable remake

The render SHALL cover the supplied clip's approximately 30-second duration,
follow its wide and close-up shot sequence, and retain its complete soundtrack.

#### Scenario: Play the render

- **WHEN** the MP4 is decoded at its beginning, middle, and end
- **THEN** animated Blender frames appear with the supplied soundtrack and the expected shot changes

### Requirement: Reproduction and honest evidence

The project SHALL provide a pinned Blender entry point, source identity,
representative comparison evidence, and an explicit account of remaining
visible deviations. A source hash mismatch SHALL fail before processing.

#### Scenario: Render the maintained scene

- **WHEN** the documented command opens the maintained Blender scene
- **THEN** it creates animation frames under the selected output directory

### Requirement: Reference-matched replacement

A replacement for the rejected draft SHALL be compared against the source's
defining facial forms, expressions, garment construction, staging, and timing.
The project SHALL record material remaining differences and the review result.

#### Scenario: The candidate retains a likeness defect

- **WHEN** a fixed-view comparison reveals a major discrepancy in an identity-defining form
- **THEN** the discrepancy is reported and the result is not described as an exact replica
