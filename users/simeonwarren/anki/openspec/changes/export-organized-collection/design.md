## Context

The source is the supplied schema-18 modern collection archive. The CLI project
is a separate change in PR #128. This PR targets master independently; its
filegroup and formatting checks do not require the CLI to be merged.

## Goals / Non-Goals

Goals: preserve every note/card identity and named field; group grammar lessons;
verify a rebuilt local archive and expose editable TOML in Git.
Non-goals: publishing scheduling, review history, media, or binary archives in Git,
changing a live Anki profile, or implementing CLI behavior here.

## Decisions

Keep each deck at decks/<hierarchy>/deck.anki.toml with its notes under notes/.
Anki's :: maps to directory separators with unsafe components escaped. Named
TOML fields use compact one-line strings unless the value has newlines. The base
archive hash binds the export; the archive remains local and supplies immutable
assets and scheduling. Runtime nextPos allocation state remains in the database.

Copy the verified collection projection from the task-owned project candidate,
then verify the independent data candidate against the original archive and the
rebuilt output. Record only source-identifying verification summaries here;
raw logs and archives stay under ignored task output.

## Risks / Trade-offs

The text does not contain media or templates, so rebuilding requires the exact
original archive and the CLI from the companion project PR. Independent SQLite,
ZIP, and TOML comparisons verify preservation; an Anki desktop import has not
been performed.

The format revision adds explicit deck titles and path references, names notes
\*.note.anki.toml, and exports existing CSS/templates in one file per note type under the referenced note_types.path directory. Collection settings are omitted from text. Collection configuration and scheduling presets are not modified.
