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
   card declarations, with deck titles and locations declared in config files.
2. Edit the text. Change fields and tags, move notes between decks, or create
   new notes with generated identities. Edit card CSS and template HTML.
3. Inspect the plan, then apply it to an offline collection or build a new
   archive from the original.

## Features

- Named fields preserve Anki HTML and media references.
- Stable identities keep file renames separate from note creation.
- Plans show additions, updates, and deletions before reconciliation.
- Fresh identities can be assigned to copied notes or directory trees.
- Card CSS, template formats, and browser fonts are editable.
- Existing scheduling and media remain in the collection.

## Status

The CLI is in development. It supports schema-18 modern collection archives
and offline SQLite collections. Reconciliation preserves review history while
retaining existing card scheduling. Note type structure and scheduling presets remain unchanged; media stays
in the original archive.

[Read the documentation](/docs/projects/anki_as_code/)

Deck and note locations are configured through `collection.anki.toml` and
`deck.anki.toml` markers. Anki titles are explicit; directories are flexible.
Edit card CSS and front/back templates in the files under the configured `note_types_path` directory.
