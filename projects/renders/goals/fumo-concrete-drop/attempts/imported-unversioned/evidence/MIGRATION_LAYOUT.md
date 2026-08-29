# Migration layout

The durable byte-exact legacy archive remains at
`projects/renders/docs/fumo_concrete_drop/`.

The goal importer accepts root Markdown files only. The two nested attempt
files were therefore copied byte-for-byte into this ignored staging directory
as `legacy_attempt_01.md` and `legacy_attempt_02.md`. Their source paths are:

- `attempts/attempt_01.md`
- `attempts/attempt_02.md`

No source archive file was changed by this staging operation.
