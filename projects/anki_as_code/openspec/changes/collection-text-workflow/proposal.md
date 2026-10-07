## Why

Anki collections are difficult to inspect and edit in version control. A Go CLI
should expose editable TOML notes while retaining the original archive as
the authoritative source for card state, templates, and media.

## What Changes

- Add `projects/anki_as_code` with export, plan, apply, build, and generate-id commands.
- Discover decks and notes through collection.anki.toml path references and explicit deck.anki.toml titles. Expose per-deck settings and card appearance (CSS and templates) through referenced TOML; preserve collection settings without exporting or reconciling them.
- Publish accurate project documentation and a landing page through the main Hugo site.
- Represent each note's shared fields in TOML listing
  every generated card and its deck; preserve embedded HTML losslessly.
- Assign fresh identities to copied note files and directories. Publish the one-off organized user collection in a separate export PR.
- Rebuild a collection archive from the original archive and edited text, preserving
  review history and all unmanaged state without exporting it into text.

## Capabilities

### New Capabilities

- `collection-text`: Export, validate, edit, assign identities, and rebuild existing collections.

### Modified Capabilities

None.

- Treat text as desired state: plan additions, updates, and deletions, then
  reconcile an offline collection database or rebuild an archive.

## Impact

New Go project and Bazel dependency wiring. Reuse TOML and Zstandard and upgrade
the existing SQLite driver with its user-approved transitive dependencies. No
live Anki profile or host configuration change. Collection data and its
large-inventory delivery support are published separately.
