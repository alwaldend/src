---
title: Anki collection
description: Editable collection notes, decks, and settings
---

Editable TOML exported from `collection-20261007190123.colpkg`. The base archive
SHA-256 is recorded in `collection.anki.toml`. The export contains 7,218 notes and
7,485 cards. All 445 titled Chinese grammar lessons live under their proficiency
levels in `decks/main/chinese/chinese-grammar-wiki-782551504/`.

Use the Anki as code CLI from [PR #128](https://github.com/alwaldend/src/pull/128) to export, plan,
apply, assign fresh note identities, or rebuild an archive. The root config
references deck, note, and note-type locations through tables with a `path` key. Each `deck.anki.toml` declares its exact
Anki `title` and `notes.path`; note files use `<name>.note.anki.toml`. Deck strings follow the same quoting policy as note fields. The initial
hierarchy is convenient for browsing, but directory names do not define deck
titles. Notes use single-quoted literal values when safe and triple-double-quoted multiline values,
avoiding escapes in HTML. They retain compact named fields, stable identities, tags, and explicit
card declarations. Each card defaults to the title of its declaring deck marker,
unless it has an explicit override.

The root `note_types.path` reference points to `note_types/`. Each
`*.notetype.anki.toml` contains one note type's CSS and `[[templates]]` entries.
There are five note types and six card templates. Styles, front/back HTML,
browser formats, and browser fonts are editable. Collection settings and
scheduling presets are neither exported nor reconciled; they remain preserved
in the archive. The archive and review history are not stored here. Rebuilt
archives retain media, card styles/templates, deck settings, existing card
scheduling, and review records unchanged.
Keep the original archive locally as the base for rebuilding:

```sh
bazel_agent bazel run //projects/anki_as_code -- build \
  --base /absolute/path/collection-20261007190123.colpkg \
  --text "$PWD/users/simeonwarren/anki/collection.anki.toml" \
  --output "$PWD/out/anki-as-code/collection-organized.colpkg"
```

The initial adoption and its verification are tracked in
[export-organized-collection](openspec/changes/export-organized-collection/).
