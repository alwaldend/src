---
title: Anki as code
linkTitle: Anki as code
description: Export Anki collections and work with schema-backed configurations
layout: landing
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code represents notes, decks, and card appearance in schema-backed TOML.
Deck directory names define Anki titles, with each deck containing `notes/` and
child `decks/` directories. Note types live in collection-root `note_types/`.
Use `config dump` (alias `conf dump`) to inspect parsed configuration as JSON by default or select
TOML with `--format toml`. Configuration paths identify resources; references use
readable paths and names. These commands work directly on text files without Anki running.

Use `col export --collection collection.colpkg --output configs` to write an editable
configuration tree from an offline archive or `.anki2` database. Export preserves
the source and writes configs only; explicit state commands manage collection-backed
identity bindings separately. Keep Anki closed during collection operations.

Planning, apply, and archive rebuilding remain in the separate reconciliation
change.

[Read the documentation](/docs/projects/anki_as_code/)

Collection commands use the `col` group, also available as `collection`.
