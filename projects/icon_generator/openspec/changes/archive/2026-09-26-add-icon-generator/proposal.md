## Why

Provide a small Go CLI for generating random pixel-pattern icons with
user-selected canvas dimensions, colors, density, size, and shape. Go's
standard image and PNG packages support the entire rendering pipeline, so
the tool can run without external image libraries or FFmpeg.

## What Changes

- Add `projects/icon_generator` with an `icon_generator` CLI that writes PNG
  images.
- Configure canvas width and height, an opaque or transparent background,
  and a nonempty palette of allowed pixel colors.
- Configure logical pixel size and shape (`square` or `circle`). Place
  shapes in grid cells and leave uncovered areas as the background.
- Configure per-cell occupancy probability with density from zero to one.
- Support an optional seed for reproducible patterns and an output path.
- Validate inputs and report file or encoding failures through a nonzero
  exit status.
- Provide Bazel integration, usage documentation, and end-to-end checks
  that retain generated images as inspectable artifacts.

## Capabilities

### New Capabilities

- `random-icon-generation`: Generate configurable, optionally reproducible
  random arrangements of square or circular pixels as PNG images.

### Modified Capabilities

None.

## Impact

The new project owns the command, renderer, tests, and OpenSpec artifacts.
Implementation will use the existing Go and Bazel toolchains with only Go
standard-library imports; no new external dependencies are needed.
Repository integration includes the project catalog, normal project
documentation and landing page, and registration of its OpenSpec sources
with the existing validation workflow.
