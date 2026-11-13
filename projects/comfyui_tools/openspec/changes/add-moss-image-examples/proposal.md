## Why

The operator requested MOSS-Image workflows as examples for the personal rental.
The existing generated workflow catalog does not include this model family.

## What Changes

- Add six native MOSS examples: SFT/Base text-to-image and layout guidance,
  image editing, and color/depth control using prepared maps.
- Generate native JSON from the existing Libsonnet graph infrastructure.
- Document inputs, sampling defaults and the limits of example acceptance.

## Capabilities

### New Capabilities

- `moss-image-examples`: Native example graphs for MOSS generation and adapters.

### Modified Capabilities

None.

## Impact

Changes workflow generation, port contracts and documentation. Uses only native
ComfyUI nodes from the existing v0.44.0 runtime. Model downloads and deployment
remain owned by `infra/image_gen` and its matching change.
