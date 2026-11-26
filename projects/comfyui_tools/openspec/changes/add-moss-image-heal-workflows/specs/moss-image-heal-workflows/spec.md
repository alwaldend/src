## Purpose

Provide editable ComfyUI workflows and native Apps for generation with the MOSS-Image-heal checkpoint.

## ADDED Requirements

### Requirement: MOSS-Image-heal workflow catalog

The workflow build SHALL generate text-to-image, layout-guided and prompt-sequence graphs that select the MOSS-Image-heal checkpoint.

#### Scenario: Generate images from a native App

- **WHEN** a user imports the Heal text-to-image or prompt-sequence App in ComfyUI
- **THEN** App Mode exposes the relevant prompt, dimensions and sampling controls and displays the saved images.
