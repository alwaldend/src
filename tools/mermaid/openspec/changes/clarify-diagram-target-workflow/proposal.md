## Why

Diagram previews took unnecessary time because the workflow encouraged lower-level
CLI rendering and mixed system overviews with implementation detail. The user
requested clear guidance for ordinary repository diagrams.

## What Changes

- Start from the requested topology, level of detail and the owner's render target.
- Use normal SVG/WebP targets and the owning preset for local previews as well as maintained output.
- Explain dark WebP selection, publication boundaries and inspecting the actual built image.
- Extend the existing offline skill-evaluation cases for these decisions.

## Capabilities

### New Capabilities

None. This is procedural documentation with `skip_specs: true`.

### Modified Capabilities

None.

## Impact

Only the Mermaid skill, its existing evaluation cases and this change record.
Renderer code, shared themes, dependencies and discovery registration are unchanged.
