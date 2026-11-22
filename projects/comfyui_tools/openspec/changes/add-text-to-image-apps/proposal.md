## Why

Users should be able to generate images from a prompt in ComfyUI without opening
the graph editor. Text-to-image workflows need the same native App entry point
as the existing NoobAI Apps.

## What Changes

- Build an importable `.app.json` companion for every `text-to-image` workflow,
  including the MOSS-Image Base variant and archived Chroma Flash example.
- Present prompt, canvas and sampling controls alongside the saved-image gallery.
- Document how to build and import the App workflows.

## Capabilities

### New Capabilities

- `text-to-image-apps` provides native App companions for text-to-image graphs.

### Modified Capabilities

## Impact

The ComfyUI tools workflow catalog and its Bazel output list gain App documents.
The Apps use the models and native nodes already required by the catalog.
