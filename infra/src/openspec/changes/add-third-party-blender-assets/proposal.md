## Why

Blender projects need reusable characters, dance motion, furniture, and scenes
that can be retrieved reproducibly, including the Pippa mesh and diffuse
texture used in animation work.

## What Changes

- Add checksum-pinned Bazel downloads for Pippa and other Blender assets.
- Pin the Mixamo dance clips used in the final Pippa album from the same
  Hugging Face dataset used for production.
- Document asset provenance, use terms, and Blender import procedures.

## Impact

The third-party catalog and root Bazel module gain opt-in asset targets. The
change introduces no new runtime or Blender add-on requirement.
