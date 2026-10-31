## Purpose

Provide repository video guidance that respects the requested media, selects
usable tools, and verifies the resulting artifact against its inputs.

## ADDED Requirements

### Requirement: Discoverable video procedure

The repository SHALL provide a discoverable `repo-video` skill for video
creation, editing, conversion, muxing, and validation. It SHALL route diagram
authoring, scene rendering, and frame interpolation to their existing owners.

#### Scenario: Create a still-image video

- **WHEN** a user supplies a still image and an audio track
- **THEN** the procedure guides creation and verification of a standalone video
  while preserving the requested image and complete audio

### Requirement: Respect available tools and requested formats

The procedure SHALL inspect actual media and encoder capabilities, preserve
explicit output constraints, and distinguish successful encoding from a codec
listing. It SHALL reuse available tools without treating a media request as
authorization to install software or change host configuration.

#### Scenario: Listed encoder fails to initialize

- **WHEN** an encoder is listed but cannot create a video stream
- **THEN** the procedure diagnoses that failure and selects an available option
  compatible with the user's constraints or reports the unresolved constraint

### Requirement: Verify complete media

The procedure SHALL verify required streams, representative frames, and audio
coverage. It SHALL NOT infer complete audio preservation from container duration
alone. It SHALL distinguish stream-copy verification from lossy conversion.

#### Scenario: Export drops a final audio packet

- **WHEN** a successful export has lost audio samples despite a plausible duration
- **THEN** verification detects the loss and the procedure restores complete audio
  before claiming preservation

#### Scenario: Deliver a generated artifact

- **WHEN** media creation completes without a requested source change or upload
- **THEN** generated media and verification evidence remain in ignored task output
  and the user receives a usable artifact path
