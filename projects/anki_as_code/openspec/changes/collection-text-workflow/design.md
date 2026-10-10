## Context

The attachment has a schema-18 SQLite database compressed as `collection.anki21b`,
7,218 notes, 7,485 cards, and 126,668 input review records. Grammar lessons use the Chinese
note type's `article-title` field and A1/A2/B1/B2/C1 tags. All cards are currently
outside filtered decks. See proposal.md for motivation.

## Goals / Non-Goals

Goals: lossless existing-collection editing, reproducible text exports, and a
verified local archive. Non-goals: automatic template rendering/card generation, editing media,
legacy schema migration, or modifying a database while Anki is running.

## Decisions

Proceed with a copy-and-patch database design. Regenerating a collection from
rendered cards would lose note identities, templates, and scheduling. AnkiConnect
would require a running live profile. Direct SQLite patching is feasible for the
observed schema, with strict validation and preservation checks before handoff. Review history is unmanaged: do not export it into text, delete its rows, or
compact the database as part of reconciliation. Reconciliation only synchronizes
managed entities with their text declarations.
Anki's upstream decks.proto and colpkg exporter define deck blobs and ZIP structure.

Each TOML file represents a note and lists all of its cards as declarations;
this avoids conflicting edits to shared fields. Named TOML strings preserve every original field, including HTML. Text is independent
of any particular archive checksum. collection.anki.toml declares decks_path,
note_types_path, and optional notes_path for cardless notes. Deck discovery searches
recursively for deck.anki.toml markers, each declaring title and notes_path.
Paths resolve relative to their declaring config and may be relative or absolute.
Only \*.note.anki.toml files are recursively discovered under declared note roots;
unrelated files are ignored. Deck titles and identities do not depend on paths.
The exporter uses a conventional hierarchy as a starting layout, not a reader
constraint. Missing notes directories represent empty decks; missing configured
deck/note-type roots fail instead of silently deleting resources. Duplicate note
ownership and symlink traversal fail before mutation.

Use already pinned modernc SQLite, klauspost Zstandard, and TOML.
Register Anki's Unicode case-insensitive SQLite collation before opening databases.
Retain ZIP entries verbatim except the collection database. New normal decks declare their configuration in text; IDs are allocated deterministically.

## Risks / Trade-offs

- Unsupported schema → reject explicitly; initial support is schema 18 modern colpkg.
- Invalid/conflicting text → require unique explicit note/card identities and matching field layouts before committing a private transaction.
- Cloze edits altering generated cards → reject changes to cloze ordinal sets.
- Personal scheduling → never export it into tracked TOML; output archive is local.
- Existing HTML → preserve exact field strings; no rendering or lossy conversion occurs.
- Archive corruption → atomic output plus SQLite integrity and independent full-table
  comparisons of the attachment and rebuilt archive.

## Migration Plan

Write behavioral failure cases before implementation. Export the attachment,
organize lesson deck names in text, build from the untouched base, re-export and
compare outputs. Keep the original attachment unchanged. Deliver the CLI project in PR #128. The independent collection export and its
verification belong to users/simeonwarren/anki/openspec/changes/export-organized-collection
and a separate PR targeting master. Keep the archive in local output.

## Desired-state reconciliation

The user clarified Terraform-like reconciliation. Text owns notes, explicitly
declared cards, decks, and note type appearance. Collection settings stay unmanaged. Plan is read-only; apply
updates an offline schema-18 database or archive using a private transaction. Omitted notes/cards
are deleted; new notes/cards require unique IDs, GUIDs, existing note types, and
explicit card ordinals. Existing scheduling is preserved; new cards start new.
The base archive retains templates, scheduling presets, and media. Deck names
are declared in deck.anki.toml files; users declare required parents there.
Build, plan, and apply accept the selected compatible collection without archive-hash binding.
SQLite upgrades from existing v1.20.3 to v1.60.1 for Anki collation support.
The user approved the current release and its transitive dependencies; the
registry reports a 2026-09-29 release requiring Go 1.26.0, compatible with
the repository Go 1.26.5 toolchain.

Generate-id assigns timestamp-based unused note/card IDs and a random GUID to
existing TOML note files selected by file or recursive directory arguments. It
reserves identities in the containing export, accepts duplicate old identities
in copied notes, and preserves every non-identity declaration and filename.
Lesson organization is a one-off data adoption task outside the reusable CLI.
Notes use human-readable \*.note.anki.toml basenames. Export initially maps
Parent::child to parent/child; moving files or directories leaves Anki titles
unchanged unless a deck title or card declaration changes. Field names use bare TOML keys when valid. Deck serialization shares the field-string writer, including when persisting new deck IDs. Single-line values prefer literal strings; values containing newlines use
triple-double-quoted basic strings without blank lines between fields. Ordinary
double quotes remain unescaped in multiline values; delimiter conflicts and
backslashes are escaped. Single-line values with apostrophes or unsupported
controls use basic strings. Multiline basic strings use line continuation to
preserve the exact original trailing newline.

Apply also supports atomic replacement of an offline .colpkg; build preserves
the original and produces a separate archive. Allocated deck IDs are recorded
back into deck.anki.toml after successful reconciliation so renames preserve identity.

The nextPos allocation cursor belongs to database runtime state rather than
declarative settings. Reconciliation preserves it and advances it atomically
for newly allocated note positions; siblings share a position, including an
existing new sibling position when available. Archive apply retains file modes.

## Card appearance

note_types_path references a directory. Recursively discovered \*.notetype.anki.toml
files each contain a note type at the document root, with stable ID/name, CSS,
and [[templates]] entries for front/back/browser formats and browser fonts.
Export names these files after the note type; identities are read from content.
Collection settings are neither exported nor reconciled. The database config
records remain untouched except the computed nextPos cursor when allocating
new-card positions. Multiline content follows the note-string policy.
Changes patch only the corresponding protobuf fields in existing note types
and templates, preserving unrelated configuration and exact bytes when
unchanged. Structural identities and template ordinals cannot change; adding
or removing note types/templates is outside this editing surface. Template
changes do not generate cards; explicit note/card declarations remain authoritative.
Scheduling presets remain in the base archive at the user's clarified request.

## Decision review: references and appearance patching

Proceed with config references and explicit deck titles. Derived titles would
rename decks during ordinary file organization; fixed folders would contradict
the requested layout independence. Relative paths are portable; absolute paths
support external note storage, with --collection providing explicit identity
reservation there. Missing roots, symlinks, and overlapping note ownership are
rejected before mutation. Generated upstream messages retain unknown fields on
unmarshal, and only declared managed members are assigned before marshal.
Untouched blobs remain byte-identical; scheduling presets remain untouched.
E2E cases challenge rearranged paths,
external identities, exact CSS/template changes, and unchanged card rows.
Message definitions come directly from the pinned upstream notetypes.proto and
decks.proto; generated code is owned by Bazel rather than checked-in copies.

Path references use flat decks_path, note_types_path, and optional notes_path strings. Deck markers use notes_path. Readers and writers use only these flat keys. Repository-only protobuf bytes and scheduling cursors stay in a private snapshot rather than shared models.

## Reconciliation review corrections

Apply verifies the input hash immediately before publishing the completed
archive or SQLite copy. Archive compression completes before this final check.
The check detects observed concurrent changes; it is not an atomic filesystem
compare-and-swap against a writer that ignores the offline-use requirement.

Template edits also update the parent note type's modification time and pending
sync marker, because Anki synchronizes complete note types. CSS bytes stay
unchanged when only templates change.

Deck reconciliation temporarily renames changed or removed decks within its
private transaction, using names distinct from every current and desired title.
This permits title swaps and reuse of deleted titles with Anki's unique name
index. Temporary names never reach the published database. Conversion of an
occupied filtered deck to a normal deck fails validation rather than changing
unmanaged filtered-card scheduling state.

Card snapshots retain the effective home-deck ID (odid for filtered cards,
otherwise did) as internal state omitted from TOML. Plan and reconciliation
share a comparison against the desired title's resolved deck ID. This allows
unchanged title overrides to move to another deck after a title swap, while
renaming the same deck does not mark its cards modified.

Export omits archive checksums from the manifest. Legacy base_sha256 metadata
is accepted without binding, preserving compatibility with existing text.
Runtime hashing remains only as an apply concurrency check. Scratch creation
uses Go's standard temporary-directory selection, without reading TMPDIR
manually. Database protobuf data uses Go types generated from Anki's upstream
schemas, pinned to release 26.09.3 at commit
29bb700b951e3f0c0cb69b77c0180fc1fe33e6ba. Bazel downloads the unmodified schemas
with archive integrity verification and generates types using the existing
protobuf toolchain. No copied schemas or generated Go files are checked in.
Typed messages replace manual field-number decoding and encoding. Unknown fields
are retained by protobuf unmarshaling; untouched blobs are reused verbatim.
Deck settings outside the editable subset remain in the existing opaque extra
payload, including fields known to the upstream schema. Scheduling presets are
not decoded or modified by this change.

Schema review fixes in #131 use the integrity-verified 26.09.3 tag archive,
upstream `github.com/ankitects/anki/proto/anki/` Go imports, and external target
visibility restricted to the local dependency wrapper. The CLI uses the two
consumer aliases and matching Gazelle resolutions.

The schema prerequisite merged in #131. The project branch now targets master
with that dependency already present; the dependency package and root module
include are absent from the project PR diff.

## API and structure review

Export refreshes declarations in an existing output directory using a fully
staged snapshot. It preserves unrelated files, replaces local declaration
edits with the selected collection, and removes stale declarations beneath
that directory. Export continues to generate the documented default layout;
external paths are not rewritten. Archive failures leave previous output intact.

All public path inputs and configured resource paths expand current-user `~`
and `~/...`. Directly supplied root files may have any filename; directory
discovery still selects collection.anki.toml. Standard filepath.WalkDir remains
the discovery mechanism. Database setup registers the collation explicitly
once at open time. A repository type owns snapshot reads, transactions, SQL
execution, integrity verification and convergence; generated Anki note-type
messages carry structural metadata. Generated messages define the editable contract; runtime wrappers retain
filesystem information absent from protobuf messages. Deck codecs live in ankiformat. Go skill changes ship separately against master.

CLI status messages use stderr; plan JSON and help remain stdout. Shared flag definitions own names, bindings, and help text while commands retain their local values and required flag sets. Verify through the built binary, including changed note identities and required-flag failures.

## Review decisions: native encoding and package ownership

The latest review supersedes the earlier exported quoting preference: serialize complete TOML objects with the existing encoder, using its normal quoting and layout. Exact field, template, and CSS values remain the acceptance contract.

Separate internal/model (data types only), internal/validation (desired-state rules), internal/planning (state comparison), internal/ankiformat (configuration codecs and text normalization), internal/textstore (config discovery, TOML, identities, staged text output), internal/archive (ZIP/zstd and offline-file resources), internal/repository (SQLite opening, snapshots, transactions and all appearance SQL), and internal/collection (operation orchestration). Filesystem storage and archive handling do not depend on SQLite. Domain validation and planning do not depend on filesystem or database operations. Existing E2E fixtures verify raw rows through test-only handles, rather than exposing database handles in the application facade.

## Boundary decision verification

Proceed with separate Go packages under the existing repository module. Separate
module manifests would add dependency and release coordination without improving
these internal boundaries. Explicit Bazel dependencies and restricted adapter
visibility enforce the direction: collection coordinates model, textstore,
archive, and repository; repository and textstore use the pure validation, planning, and format packages
as needed; those packages depend on model and never import storage. Archive
has no project dependency. SQL handles remain private, while E2E fixtures own
independent raw handles for row comparisons. Native TOML encoding retains the
parsed-value contract and removes custom serializers. The retained workflows
challenge malformed text, appearance edits, archive preservation, and convergence.

Cards resident in filtered decks are unmanaged. Repository snapshots retain only
private references for protection and exclude those cards from the shared
managed state. Repository validation combines pure desired-state validation
with protection against identity/ordinal conflicts and orphaning removals. Both
plan and apply use that validation boundary; reconciliation never updates or
deletes filtered card rows.

The first-party collection protobuf schema owns editable TOML and plan JSON
field names, scalar types, and optional presence. TOML is loaded into generic
values, marshaled as JSON, then parsed with ProtoJSON into
generated messages used directly throughout the application. Runtime wrappers
keep metadata outside those messages. Export uses ProtoJSON with original field names and default values, then converts
integer fields to TOML integers using descriptors without floating-point decoding.
These small format helpers belong to textstore; there is no separate serialization
package or custom schema validator. CLI plans use ProtoJSON directly.
Runtime models do not duplicate serialized fields. Database protobuf values continue
to use the generated upstream Anki schemas without copied definitions.

Generated collection protobuf messages are the sole definitions of serialized data. Notes and decks wrap message pointers with runtime-only metadata; current card deck identities belong to the snapshot. Serialization consumes messages directly, without maintaining parallel field definitions or conversion tables.

The API and runtime models are provided by merged [PR #136](https://github.com/alwaldend/src/pull/136). This implementation uses those packages from master; its diff contains serialization and reconciliation behavior without duplicating the prerequisite additions.
