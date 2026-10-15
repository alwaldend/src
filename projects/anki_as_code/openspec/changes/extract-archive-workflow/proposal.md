## Why

Extract collection archive resources from #128 and delegate package import and
export to Anki. Manual ZIP and Zstandard handling duplicates Anki's format rules
and can lose media as those rules evolve.

## What Changes

- Use the pinned Anki backend's collection-package methods, including media.
- Accept an existing gRPC connection; the caller owns worker lifetime.
- Share guarded, destination-side publication for packages and offline databases.
- Remove direct archive encoding and the public database-path accessor.
- Exercise real native packages, identities, review history, media, and failures.
- Preserve #128 and its unpublished diagram edits until this prerequisite merges.

## Capabilities

### New Capabilities

- `archive-resources`: Native package resources and guarded file publication.

### Modified Capabilities

None.

## Impact

The archive package and README change along with narrow build visibility and
bindings for already pinned upstream schemas. No dependency is added. Anki owns
archive compatibility and database semantics; application code contains no SQL.
