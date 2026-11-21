## Purpose

Make the MiniMax H3 Base video model available on the image-generation rental
for use with native ComfyUI workflows.

## ADDED Requirements

### Requirement: Selectable H3 Base weights

The image-generation registry SHALL offer a `minimax-h3-video` profile with
revision-pinned public weights for FL2VA inference with stereo audio.

#### Scenario: Install the H3 Video profile

- **WHEN** the operator selects `minimax-h3-video` and runs the application installation
- **THEN** the installer places the denoiser, text encoder, video VAE, and audio VAE in the matching ComfyUI model directories
- **AND** the installer saves the prepared text-to-video and first/last-frame workflows in the ComfyUI workflow library
