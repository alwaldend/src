---
title: Anki as code
linkTitle: Anki as code
description: Manage Anki collections through editable TOML
layout: landing
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Keep notes, decks, and settings in text files. Edit them with your usual tools,
review a plan of the changes, and reconcile your Anki collection with the result.

## How it works

1. Export an existing collection to TOML. Each note keeps its named fields and
   generated-card deck assignments, with deck titles derived from directory names.
2. Bind config paths to existing Anki objects with `col state import`. Export leaves
   the source collection unchanged and does not create or install state.
3. Edit the text. Change fields and tags, move notes between decks, or create
   new notes identified by their configuration paths. Create note types and edit their fields, templates, CSS, and configuration.
4. Inspect the plan, then apply it to an offline collection or build a new
   archive from the original.

## Features

- Named fields preserve Anki HTML and media references.
- Collection-backed state and explicit moves preserve existing identities.
- Plans show declared changes and known card assignments; Anki generates cards during apply.
- Anki validates templates, generates cards, and rebuilds field caches.
- Export writes configs only; state dump exports stored bindings separately.
- State import binds independently written configs to existing objects.
- Apply and build leave configurations unchanged; Anki owns IDs and GUIDs.
- Note types, ordered fields, templates, CSS, and browser fonts are editable.
- Existing scheduling and media remain in the collection.

## Status

The CLI is in development. It supports modern collection archives and offline Anki collections supported by the pinned
backend. Reconciliation preserves review history while
retaining existing card scheduling. Scheduling presets remain unmanaged; media stays
in the original archive.

[Read the documentation](/docs/projects/anki_as_code/)

The collection contains `note_types/`, root `notes/`, and `decks/`. Each deck
contains its marker, `notes/`, and optional child `decks/`. Anki titles derive
from this hierarchy. Edit CSS and templates in `note_types/` declarations.

Collection commands use the `col` group, also available as `collection`.
