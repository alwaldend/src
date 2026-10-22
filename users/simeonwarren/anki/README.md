---
title: Anki collection
description: Editable collection notes, decks, and settings
---

Current CLI version-1 TOML exported from `collection-20261017235526.colpkg`.
The export contains 7,218 notes, 7,485 existing cards, five note types and 513
decks: 485 source decks plus 28 Guzhenren grouping decks. Native identities are
saved in `state.anki.toml`; note files contain editable fields, tags, note-type
references and card assignments. Deck titles derive from literal directory
names, with a `decks/` directory between hierarchy levels.

The source archive SHA-256 is
`1fc03bb25a00418a19b927adbb3c18f85a93e58724932cd836152df60f83d3a5`.
The original archive, media and review history remain local and are not
committed. The original archive remains byte-identical. Note fields and tags,
native note/card identities, templates, presets, due dates, intervals and review
history were verified against a rebuilt archive. Anki recalculates FSRS memory
estimates when moving cards between decks; this build changes 336 FSRS payloads
while leaving due dates, intervals, review counts, queues and other card data
unchanged.

Grammar lessons live under
`decks/main/decks/chinese/decks/chinese-grammar-wiki-782551504/decks/`.
All 449 lesson decks use `001. `, `002. ` and subsequent three-digit prefixes,
by the earliest contained note creation ID within each proficiency level.
Imported deck timestamps can reflect alphabetical recreation, so only the four
empty source lessons fall back to their deck creation IDs. Existing note files
within each lesson sort by their original note creation IDs. There are 445
populated lessons. B2 contains 142 lessons and 1,342 notes; C1 contains
12 lessons and 96 notes. Their directories are `B2/decks/` and `C1/decks/` under
the grammar level directory. Two lesson titles replace `/` with `or` for valid
literal filenames; article-title fields are unchanged. The filtered decks use
the readable titles `Due cards` and `New or due cards`; their search queries
are unchanged. The source's removed,
empty `B1_unsorted` deck is absent.

Guzhenren lives under `decks/main/decks/chinese/decks/guzhenren/`.
Chapter 1 has 133 sentence notes and 714 vocabulary notes. Its `sentences` and
`words` parents each have 14 matching numbered groups, from
`001. Sentences 001-010` through `014. Sentences 131-133`. Sentence files have
global reading-order prefixes, such as `001. 1692893034664.note.anki.toml`.
The sequence includes the chapter heading, narrative and author's closing note.
Vocabulary follows its example passages and uses global `001. ` through `714. `
filename prefixes, preserving numeric note-ID order within each passage group.
Most vocabulary IDs are imported sequence numbers rather than recoverable
wall-clock timestamps. Shared import/edit modification times do not determine
the order. Two incomplete examples use the
associations recorded in the [verification](openspec/changes/archive/2026-10-22-refresh-collection-20261017/tasks.md).

The state file is the current CLI's unmodified export snapshot: it binds all
485 source decks, 7,218 notes, five note types with their fields/templates, and
one referenced scheduling preset (`Default`). Both native presets remain in the
archive. `[[moved]]` entries in `collection.anki.toml` connect
those exported addresses to the organized paths. Grouping decks are created by
native reconciliation. Keep the state and moves together to retain identities.

State update writes its collection argument in place. Import state into a copy
of the original archive, inspect the plan, then build a separate organized
archive with the current CLI:

```sh
mkdir -p "$PWD/out/anki-refresh"
cp /absolute/path/collection-20261017235526.colpkg \
  "$PWD/out/anki-refresh/collection-state-bound.colpkg"
bazel_agent bazel run //projects/anki_as_code -- col state update \
  --collection "$PWD/out/anki-refresh/collection-state-bound.colpkg" \
  "$PWD/users/simeonwarren/anki/state.anki.toml"
bazel_agent bazel run //projects/anki_as_code -- col plan \
  --collection "$PWD/out/anki-refresh/collection-state-bound.colpkg" \
  --text "$PWD/users/simeonwarren/anki"
bazel_agent bazel run //projects/anki_as_code -- col build \
  --base "$PWD/out/anki-refresh/collection-state-bound.colpkg" \
  --text "$PWD/users/simeonwarren/anki" \
  --output "$PWD/out/anki-refresh/collection-organized.colpkg"
```

The native build succeeds, and planning the rebuilt package against this tree
returns an empty plan. All 254 media filenames and decoded contents are
preserved. The text uses the current CLI format; no legacy CLI is required.
Verification is tracked in
[refresh-collection-20261017](openspec/changes/archive/2026-10-22-refresh-collection-20261017/).
