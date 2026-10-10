---
title: Anki schemas
description: Pinned upstream protobuf schemas for Anki collection storage
---

This package fetches Anki's unmodified protobuf schemas from
[release 26.09.3](https://github.com/ankitects/anki/releases/tag/26.09.3), commit
`29bb700b951e3f0c0cb69b77c0180fc1fe33e6ba`, downloaded through its release tag with archive integrity verification.
The Anki schema headers declare GNU AGPL version 3 or later; they remain intact
in the downloaded source. No upstream schemas or generated Go sources are
copied into this repository.

Bazel generates Go messages for decks and note types with the existing protobuf
toolchain. Their import closure includes generic, sync, and collection schemas.
Generated packages use upstream `github.com/ankitects/anki/proto/anki/` import
paths. External targets are visible only to this wrapper package, which grants
its consumer access through the local aliases; no new Go runtime module is
required.
