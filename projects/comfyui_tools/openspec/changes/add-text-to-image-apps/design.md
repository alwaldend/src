## Context

The workflow build emits native ComfyUI JSON for each model family. Existing
NoobAI Apps use a `.app.json` suffix and ComfyUI's `extra.linearData` selections.

## Goals / Non-Goals

Cover the catalog's single-image `text-to-image` and `base-text-to-image`
purposes. Reference, layout, sequence and video workflows have separate purposes.

## Decisions

Generate the Apps from the existing graph definitions. Use a shared control list
for native KSampler workflows and the native guidance and noise widgets for
FLUX.2 dev. The FLUX.2 dev App uses two native integer inputs to set both the
latent size and the size-dependent noise schedule. The existing graph documents
retain their serialized contents.

## Risks / Trade-offs

The selected widgets must resolve against the pinned ComfyUI frontend. Validate
coverage, node and widget bindings, the image output, and FLUX.2 dev dimension
links in the generated JSON. Keep runtime model setup under its existing owner.
