## Context

The supplied collection must be exported using the current CLI's version-1 configuration schema. IDs belong in native state; titles derive from the fixed directory hierarchy. Two existing lesson titles contain `/`, which cannot be represented as literal path components (`ANKI-EXPORT-LITERAL-PATH`). See proposal.md for the requested data refresh.

## Goals / Non-Goals

Preserve exact source fields, tags, note/card identities, note-type presentation, presets, due dates, intervals and review history while updating the export and its organization. Native Anki deck moves recalculate FSRS memory estimates; that observed behavior is disclosed and checked separately. No CLI implementation change, live profile mutation, or publication of collection binaries, media, or review history is included.

## Decisions

On a task-local archive copy, use the native CLI with an isolated state namespace to change only the two slash-bearing deck titles, replacing `/` with `or`. Preserve native identities and article-title fields. Reuse existing numeric note filenames through filename bindings, then run the current CLI export with a state output. Replace only collection declarations; retain the owner README, BUILD file and OpenSpec history.

Number grammar lessons alphabetically within each proficiency level with three digits and `. ` in both Anki titles and filesystem paths. Keep their article-title fields unchanged.

Rename the two filtered deck display titles to `Due cards` and `New or due cards` so their literal filenames are valid Bazel labels. Preserve their native filtered search queries and settings.

For Guzhenren chapter 1, creation IDs provide the existing reading sequence: the title, opening dialogue, narrative, and author closing note are ordered consecutively. Divide its 133 sentence notes into consecutive groups of ten (the last contains three). Keep words and sentences as separate parents, with matching numbered sentence-range children. Match vocabulary example sentences to the sentence notes after removing HTML and punctuation; inspect ambiguous and missing examples separately. Prefix sentence filenames with their global three-digit reading positions. The names document reading order without modifying review scheduling.

Save the actual export's state as `state.anki.toml`. Keep its original exported addresses and declare explicit directory/note moves to the organized addresses in `collection.anki.toml`. The 28 grouping decks are new declarations without native identity bindings. When state is installed on a copy of the original archive, native plan/build applies the renames and card moves while preserving existing identities. Verify a rebuilt package, media, scheduling, complete review history and an empty post-build plan.

## Risks / Trade-offs

The old checked-in format is incompatible with the current CLI, so retaining it is not an acceptable final export. Native package reconciliation and an empty plan establish compatibility of the final candidate. State update mutates its collection argument; documented commands operate on a copy so the original archive checksum remains valid. Sentence creation order is supported by the text sequence but there is no explicit upstream ordinal field. Uncertain vocabulary associations remain explicit in verification evidence.
