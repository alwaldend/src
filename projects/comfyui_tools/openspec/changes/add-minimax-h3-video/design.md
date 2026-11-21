## Context

ComfyUI 0.44.0 provides native FL2VA conditioning, sampling, and video output.

## Decisions

- Share one graph for text conditioning, sampling, audio/video decoding, and
  MP4 output. The keyframe variant adds two image loaders.
- Resolve the four loader filenames through the rental profile registry.
- Use the base sampler at 864×480, 107 frames, and 24 fps for a short first render.

## Risks / Trade-offs

Offline checks can establish graph and installation compatibility. Duration,
memory use, and visual continuity require a live GPU render.
