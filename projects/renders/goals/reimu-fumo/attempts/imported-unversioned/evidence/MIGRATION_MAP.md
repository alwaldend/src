# Reimu Fumo legacy migration map

The authoritative, byte-preserved legacy history remains at
`projects/renders/docs/reimu_fumo/`.

This flat staging directory exists only because the v1alpha1 importer accepts
regular Markdown files at one level. The mapping is deterministic:

- Legacy root Markdown files retain their basenames.
- `attempts/README.md` becomes `attempts_README.md`.
- Every `attempts/attempt_*.md` file retains its basename.

Every copied payload is byte-identical to its archived source. The imported
result is one closed historical snapshot and does not claim that the prose
ledger was reconstructed as individual structured attempts. New structured
work starts after that snapshot.
