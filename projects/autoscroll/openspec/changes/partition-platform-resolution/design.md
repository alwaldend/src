## Context

The upstream importer reads package resolution markers but does not follow
incoming platform-qualified edges. The shared repair owns this requirement.

## Decisions

Use `tool.uv.environments` to partition macOS and its complement. Their union
preserves all platforms. uv then records PyObjC packages with Darwin resolution
markers understood by the unmodified importer. Keep native lock imports and
the existing `uv_lock.update` target. The user approved tomli 2.5.0, required
by the existing PySide6 package on Python 3.10. Existing distribution hashes
and versions remain unchanged.

## Risks / Trade-offs

Inspect generated package markers and actual import behavior on Linux, macOS,
and Windows. Verify repeated locking is stable and full-site genquery loads.
Native compilation and rendering acceptance run on the available Linux host;
macOS and Windows checks establish import metadata, not execution there.
