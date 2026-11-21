## Purpose

Provide local MiniMax H3 Base generation with editable settings and stereo
audio using the model files installed by the rental profile.

## ADDED Requirements

### Requirement: Prepared native H3 workflows

The workflow build SHALL generate text-to-video and first/last-frame workflows
with the installed model files preselected and a native MP4 output with audio.

#### Scenario: First and last frames

- **WHEN** the operator selects both keyframes and provides a prompt
- **THEN** the keyframe workflow SHALL condition one video on the prompt and both ordered keyframes

#### Scenario: Prompt only

- **WHEN** the operator queues the text-to-video workflow
- **THEN** the workflow SHALL use the editable prompt without requiring an input image
