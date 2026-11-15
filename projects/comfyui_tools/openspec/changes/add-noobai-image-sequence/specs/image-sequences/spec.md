## Purpose

Provide a simple way to generate ordered still images with persistent references,
previous-image continuity within a scene and deliberate composition resets.
Also provide a separate shared-prompt composer for independent scene images.

## ADDED Requirements

### Requirement: Scene-aware sequential generation

The reference workflow SHALL accept shared appearance/setting text, one action per nonempty
line, and a line containing `---` to start a new scene. The first action in each
scene SHALL start from fresh noise; subsequent actions SHALL use the previous
decoded result as the starting image. The original references SHALL remain active
across scene changes. A change-strength control SHALL govern carry-forward edits.

#### Scenario: Two scenes with two actions each

- **WHEN** four actions are separated into two scenes with `---`
- **THEN** four images are generated in order, with fresh compositions for actions
  one and three and previous-image edits for actions two and four

### Requirement: Repeatable partial generation

The workflow SHALL support selecting an action range, stable per-action seeds and
an optional continuation image. A selection beginning inside a scene SHALL require
that continuation image rather than silently discarding continuity.

#### Scenario: Continue an accepted sequence

- **WHEN** a user selects action three inside a scene and supplies accepted image two
- **THEN** generation starts by editing that image with action three and uses the
  same action-three seed as a complete run with the same base seed

### Requirement: Immediate ordered outputs

Each completed image SHALL be saved as an individual PNG with grouped run identity,
scene/action numbering and actual generation metadata before the next action starts.
Interruption SHALL preserve completed files and prevent further sampling.

#### Scenario: Interrupt after the second image

- **WHEN** generation is interrupted after two images have been saved
- **THEN** both files remain available and no further action is sampled

### Requirement: Native ComfyUI App

The workflow SHALL include a native ComfyUI App definition exposing references,
shared text, actions, change strength and routine generation controls, with ordered
image output. It SHALL retain an editable native graph and use existing model pins.

#### Scenario: Import the sequence app

- **WHEN** the generated workflow is imported into a compatible ComfyUI frontend
- **THEN** the app exposes the controls and output without requiring node rewiring

### Requirement: Sequence GIF playback

The workflow SHALL retain individual PNGs and produce exactly one looping GIF
of the selected sequence in action order. It SHALL NOT publish per-scene GIFs.
The App SHALL expose frame duration in milliseconds. The GIF SHALL preserve the
canvas and each generated image's playback time. When interpolation is disabled,
it SHALL NOT synthesize intermediate motion. Identical adjacent images MAY
share a GIF frame with combined duration. Interruption SHALL retain completed
PNGs without publishing an incomplete sequence GIF.

#### Scenario: Generate two scenes

- **WHEN** a run produces two images in each of two scenes at 500 milliseconds per image
- **THEN** it saves four PNGs and one two-second combined GIF, exposing five outputs
  through the native App and no per-scene GIF files

#### Scenario: Interrupted or partial generation

- **WHEN** generation is interrupted before the selected run completes
- **THEN** completed PNGs remain available and no GIF is published
- **AND** a successful partial run's GIF contains only its selected images

### Requirement: Optional previous-image structure guidance

The sequence SHALL optionally apply NoobAI Canny ControlNet conditioning from
the preceding generated image during same-scene continuation. Every fresh scene
SHALL start without a guide from the preceding scene. A zero control strength
SHALL disable guidance. The App SHALL expose strength and end percentage.

#### Scenario: Continue across a scene break

- **WHEN** Canny guidance is enabled for two scenes with two actions each
- **THEN** only the second image of each scene uses the preceding image's edges,
  and partial continuation uses its supplied preceding image as the guide

### Requirement: Optional scene-local frame interpolation

The App SHALL expose optional 2x, 4x and 8x frame interpolation, with 1x disabling
it. It SHALL preserve keyframe PNGs and total playback time and SHALL combine
interpolated scenes with a hard cut, without interpolating across scene breaks.

#### Scenario: Smooth a two-scene sequence

- **WHEN** two two-image scenes are interpolated 2x at 500 milliseconds per keyframe
- **THEN** the combined GIF retains two seconds of playback and one intermediate frame
  within each scene, and cuts directly from scene one's endpoint to scene two's start

### Requirement: Distinct reference and composer Apps

The project SHALL provide separate native reference-sequence and scene-composer
Apps with editable graph equivalents. The reference App SHALL retain uploaded
references and previous-image continuation. The composer SHALL accept shared
character, clothing and render text, per-scene setting text and desired-image
actions. It SHALL generate every image independently, using one fixed sampling
seed by default and a second latent refinement pass with the same conditioning
and seed. Older sequence inputs SHALL retain their existing defaults.

#### Scenario: Compose two scenes

- **WHEN** four actions are grouped into two scenes with two setting descriptions
- **THEN** every image uses the same character/clothing/render text, its scene's
  setting and its action; every first pass starts from fresh noise and its second
  pass refines only that image's first-pass latent
- **AND** four PNGs and one combined GIF are saved with actual
  prompts, modes, seeds and refinement settings

#### Scenario: Regenerate an independent action

- **WHEN** the composer selects one action inside a scene without continuation
- **THEN** it generates that action with the same fixed seed and effective prompt
  as its complete run, without using any preceding generated image

### Requirement: Optional installed LoRA and preflight validation

The composer SHALL offer an optional native LoRA selection, defaulting to none.
It SHALL reject unknown installed filenames and nonfinite strengths. Prompt
composition SHALL accept either one shared setting or exactly one per scene,
rejecting mismatches. Invalid generation/refinement settings SHALL fail before
sampling. Native interruption between passes SHALL prevent further sampling and
preserve every previously completed PNG.

#### Scenario: Reject an invalid composition

- **WHEN** settings do not match the scene count or the selected LoRA is absent
- **THEN** generation fails without sampling or publishing images

#### Scenario: Cancel during two-pass generation

- **WHEN** interruption is requested after an image's first pass
- **THEN** its second pass does not start and earlier completed files remain
