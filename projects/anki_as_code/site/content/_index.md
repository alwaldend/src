---
title: Anki as code
linkTitle: Anki as code
description: Inspect and manipulate schema-backed Anki configuration files
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
Use `config dump` to inspect parsed configuration as JSON by default or select
TOML with `--format toml`. Configuration paths identify resources; references use
readable paths and names. These commands work directly on text files without Anki running.

Collection export, planning, apply, and archive rebuilding are developed in the
separate reconciliation change.

[Read the documentation](/docs/projects/anki_as_code/)
