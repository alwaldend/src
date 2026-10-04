## Why

The user requested that repository-specific Python and OpenSpec skills reside
with their owning tools. The twelve upstream OpenSpec skills remain managed
through the existing pinned archive and discovery workflow, as clarified.

## What Changes

- Move repo-python to tools/py/skills/repo-python.
- Move the repository-specific openspec skill to tools/openspec/skills/openspec.
- Update discovery labels, the maintained documentation link, and owner READMEs.
- Preserve skill instructions, evaluation assets, and invocation names.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This ownership relocation changes no supported review or tool behavior.

## Impact

Two canonical Bazel skill packages and their generated discovery links move.
No dependency, upstream skill content, or runtime instruction changes are needed.
