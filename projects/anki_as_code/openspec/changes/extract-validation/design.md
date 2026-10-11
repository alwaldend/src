## Context

The validation package in PR #128 depends only on merged models, upstream note types, and ankiformat. The format package provides deck codecs, preservation helpers, case folding, and field-cache text extraction; it has no storage dependency.

## Goals / Non-Goals

Goals: publish independently buildable, validation and supporting format packages. Non-goals: command implementation, storage adapters, reconciliation, and new dependencies.

## Decisions

Start with exact package copies from PR #128, then correct the two confirmed review findings. Track note IDs independently of GUIDs and card IDs, and inspect optional filtered-field presence instead of its value. Keep one Go module and both PRs against master. Retain the shared source in PR #128 until this PR merges so its checks remain valid; rebase afterward to remove overlap.

## Risks / Trade-offs

Temporary overlap is intentional. The original byte-for-byte comparison binds the extraction to PR #128. The two focused review corrections are intentional differences; package compilation, workflow tests, and lint verify the corrected candidate. Collection E2E coverage stays with PR #128, whose workflow already exercises these exact packages.

## Review verification

Before implementation, exercise valid resources, duplicate IDs with distinct GUIDs/cards, zero and negative note IDs, conflicting normal/empty-filtered declarations, and valid empty filtered configuration. A public-package workflow produces a JSON verification artifact. Storage E2E remains in PR #128.

## Deck conversion boundary

DecodeDeck takes the upstream payload and returns an editable collection.Deck
protobuf declaration. EncodeDeck accepts that declaration and returns a payload.
Neither mutates a persistence model: identity, paths, and payload assignment
belong to the caller. The existing protobuf message already represents editable
data, so no parallel DeckEditable structure or schema is introduced. Local
normal-deck values use config terminology; upstream oneof names remain Normal.

## Anki-owned derived caches

Remove custom HTML normalization. Existing sfld/csum values stay untouched; new
notes receive empty/zero caches. Users must run Anki Check Database after import
when fields or sort-field settings changed. Cache rebuilding affects sorting and
duplicate detection, not review history or scheduling. Check Database also has
broader repair/card-generation behavior, which the README explicitly discloses.

## Readable filtered configuration

Deck.filtered references the pinned upstream anki.decks.Deck.Filtered message
rather than duplicating its schema or exporting an opaque base64 payload.
TOML exposes search terms, enum ordering, rescheduling, delays, and preview
settings through the existing protobuf JSON bridge. Exactly one deck kind is
required. Unknown fields are preserved from the base archive during
reconciliation; resident cards and temporary scheduling stay unmanaged.

## Reuse regular configuration and audit other messages

Deck.normal references upstream Deck.Normal directly, removing custom DeckConfig
and opaque extra_fields_base64. Both kinds use the same direct conversion.
Unknown regular configuration and nested daily-limit data are retained from
existing messages. Preset configuration is still unmanaged. Other editable
messages were compared with pinned upstream config, notes, and note-type schemas:
the CLI manifest, named-field notes, scheduling-free cards, appearance subset,
and CLI plans have distinct contracts and remain purpose-specific.

## Editable note-type definitions

NoteTypeSettings directly owns upstream Notetype.Config, ordered upstream Field
messages, and TemplateSettings. Remove the runtime NoteType wrapper and its
parallel FieldNames projection. Snapshots index the generated settings directly.
New identities resolve by existing name or allocate after all reserved IDs.
Complete desired schemas validate named note fields and explicit card ordinals;
HTML/cloze content is not parsed. Existing note identities and card ownership
stay fixed, while declared card ordinal changes retain scheduling.

The repository stores new/changed types, field layouts, and templates in the
same transaction as notes/cards. It preserves matched original config bytes and
unchanged template synchronization metadata, releases names for atomic swaps,
and marks structural changes in col.scm. Filtered cards remain unmanaged and
protect their template positions and note-type kind. Unknown field/template
configuration matches upstream IDs, then names, with same-size position fallback
for otherwise unmatched renames. Deleted entities remove their owned config.
Nested settings preservation rejects removal of unknown-bearing messages.
Anki retains template parsing and derived requirement compilation.

## Review simplification and compact exports

Split resource validation into identity, field, card, deck hierarchy/transition,
and note-type configuration helpers. Keep validation-pass indexes private to
validation. Notes retain named fields plus derived Bodies for SQLite. Protobuf
preservation traverses descriptors instead of hard-coding daily limit and search
term fields. Exports omit implicit defaults while retaining optional zero values,
empty configuration presence, and empty named fields. Plans use note_type actions
and separate addition/update/deletion counters.

## Audited package visibility

Keep handwritten declarations private unless a production consumer crosses the
package boundary or a defined framework or serialization contract requires the
name. The independent validation package does not export local body checks. Remove
the unused note source-path field. PR #128 also moves test-only conveniences
into local fixtures and removes unused identity wrappers. Preserve
`ankiformat.PreserveUnknown`: reconciliation in PR #128 consumes it even while
PR #138 provides the independent format package. Generated protobuf declarations
remain governed by their schemas.
