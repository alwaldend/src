## Why

Make the supplied Anki study collection editable and reviewable as TOML, with
lessons grouped into proficiency levels and lesson subdecks. Keep this data
publication separate from the CLI project in PR #128.

## What Changes

- Add notes with \*.note.anki.toml names, explicit deck.anki.toml markers, referenced per-note-type CSS/template files, and a collection.anki.toml root under users/simeonwarren/anki.
- Organize all titled grammar lessons under their proficiency level and lesson.
- Preserve original fields and identities; retain scheduling, templates, and media in the base and rebuilt archives.
- Produce and independently verify a local rebuilt archive with original review history preserved.

## Capabilities

No product requirements change. This is user-owned collection data adopting the
CLI proposed in https://github.com/alwaldend/src/pull/128; skip_specs is true.

## Impact

Only user study data and its documentation/build exposure change here. The
large-status delivery support has its own maintained change under tools/repo_delivery.
No dependency upgrade, live Anki profile mutation, or binary archive publication.
