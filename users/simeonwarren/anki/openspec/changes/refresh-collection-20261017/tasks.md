## 1. Export and organize

- [x] 1.1 Export the archive with the current version-1 CLI and save its generated state; verify complete note/deck/type counts and current config loading.
- [x] 1.2 Number grammar lessons and organize Guzhenren using explicit moves; verify 133 consecutive sentence indices, 714 words and complete B2/C1 membership.
- [x] 1.3 Import state on an archive copy, build the organized collection, and verify unchanged identities, fields, due dates, intervals, review history, presets and media, with native FSRS recalculation recorded plus an empty post-build plan.
- [x] 1.4 Document current-format state import/build commands using a copy and record source-identifying acceptance evidence.

## 2. Delivery

- [x] 2.1 Validate the owner change and representative output, and select exact-candidate quality, lint and packaging gates. The repo-delivery receipt owns commit, push and pull-request completion.

## Source and organization

Source: `collection-20261017235526.colpkg`, SHA-256
`1fc03bb25a00418a19b927adbb3c18f85a93e58724932cd836152df60f83d3a5`.
Expected source membership is 7,218 notes, 7,485 cards, 485 decks and five note
types. B2 contains 1,342 notes in 142 lesson decks; C1 contains 96 notes in 12
lesson decks. Guzhenren contains 133 sentences and 714 words (847 notes).

Grammar has 449 lesson decks: A1 42, A2 95, B1 158, B2 142, C1 12. Four are
empty source decks, so 445 are populated. Use consecutive `001. ` prefixes
alphabetically within each level. Guzhenren uses 14 aligned passage groups of
ten sentences, except the final three-sentence group. Number sentence files
consecutively from 001 through 133 in chapter reading order.

Vocabulary association exceptions: note 5023 (`跪`) lacks an example and occurs
in sentences 051/054, both in group 006. Note 10003 (`敢`) has the unmatched
fragment `不敢`; neighboring vocabulary 10000–10002 and 10004–10005 places it
between sentences 105/106 in group 011. Combined examples spanning 120/121 are
assigned using each vocabulary headword's occurrence. Keep all note fields
unchanged.

The current exporter rejects two slash-bearing titles. A native plan on an
isolated normalization namespace verified exactly two deck updates and no
other operations. The scratch copy replaces `/` with `or` in those titles;
the original archive and article-title fields remain unchanged.

## Acceptance evidence

The current CLI produced `version = "1"` configs and the committed state sidecar.
No inline note/card resource identities or legacy collection/deck schema remain.
The raw state snapshot covers 485 source decks, 7,218 notes, five note types and
one referenced preset (`Default`). Both native archive presets are preserved. Its 1,298 explicit moves cover 449 grammar directories, two
filtered-deck directories and 847 Guzhenren notes; new grouping decks are unbound
until reconciliation. Numeric note filenames retain recognizable identities for
Git comparison without encoding those identities in config fields.

A fresh copy of the untouched source archive accepted the exported state.
The current CLI built `out/anki-refresh/collection-organized.colpkg`; planning
that archive against the organized configuration returned `{}`. Every source
note/card ID and GUID is retained, all 485 source deck identities remain, and
exactly 28 grouping decks are added. All note rows, note-type/field/template
rows, presets, tags, graves, and complete review rows match the original.
Every card's due date, interval, queue, ease, repetition/lapse counts, flags,
original deck/due metadata and non-FSRS custom data are unchanged. Exactly 847
Guzhenren cards move to the new groups. Native Anki SetDeck recalculates FSRS
memory estimates: 336 FSRS payloads differ, with no existing stability payload
removed. The CLI uses `reschedule: false`; the due/interval checks confirm that
no card was rescheduled. Collection configuration values other than the managed
state key are unchanged.

All 254 media filenames and decoded payloads match by cryptographic digest.
The supplied archive SHA-256 remains unchanged. Source deck kinds and common
settings match byte-for-byte, including both filtered search payloads. The
filtered display titles are `Due cards` and `New or due cards`, preserving their
queries while avoiding colons in Bazel source labels.

Current CLI config dump and native post-build plan also succeed against the
actual repository directory. Explicit membership checks confirm Guzhenren
847 notes (133 sentences and 714 words), B2 1,342 notes and C1 96 notes. Sentence
filenames are consecutive `001. ` through `133. ` in chapter reading order.
The state file is byte-identical to the actual native export's state output.

Task-local transformation inputs, full-package artifacts, independent SQLite,
media and config verification, and source-identifying JSON reports remain in
ignored `out/anki-refresh/`. Delivery receipts own exact-candidate gates and
publication; an Anki desktop import was not performed.

Existing note-field quoting was retained to keep Git comparisons readable. The
current CLI decodes the repository text to exactly the same complete config set
as the native build input, including card ordinal defaults; a canonical config
digest records this equivalence. The repository post-build plan remains empty.
