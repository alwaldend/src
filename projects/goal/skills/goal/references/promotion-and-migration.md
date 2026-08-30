# Promotion and migration

Use this reference only when promoting a workspace goal or importing an
unversioned record.

## Promote workspace to project

Promotion changes storage and retention, not goal identity. It must:

1. require an explicit maintained-project decision and an owning project root;
2. lock the source and destination paths in canonical order, then validate the
   source at an expected resource version;
3. preserve the goal ID and closed attempt history;
4. record source scope and promotion provenance;
5. require repository-relative or stable external links;
6. retain, copy, or make reproducible every acceptance-critical artifact;
7. scan for absolute local paths, credentials, private environment details,
   and ignored evidence that would become broken; and
8. validate and render the destination before making it canonical.

Do not promote caches or routine coordination logs. Refuse a destination that
would shadow another goal with the same ID or that cannot retain critical
evidence. Leave the workspace source intact until the destination is verified.

## Migrate an unversioned record

A README-only goal has no claimed schema version. Migration is intentionally
conservative and is always a non-destructive import. Supply a distinct legacy
source directory and destination goals root; the command never converts the
source in place.

- keep every source byte and source directory entry unchanged;
- derive the stable goal ID from the source directory name and publish only to
  `<destination-goals-root>/<goal-id>`;
- reject canonical source/target equality or ancestry overlap, then lock both
  paths in canonical sorted order;
- preserve the original prose as an immutable legacy snapshot;
- import only fields that are unambiguous;
- mark imported criteria unverified in the closed migration attempt and begin
  a new structured lifecycle at the migration checkpoint;
- build a complete `<goal-id>` directory in hidden staging under the
  destination root and run the normal record validator there;
- re-read and digest the source immediately before publishing the staged
  directory with one rename;
- make rerunning the same provenance and options idempotent; and
- refuse an existing target with different provenance, options, ownership, or
  invalid content.

Do not pretend a prose history can be losslessly normalized. The snapshot is
evidence of what was recorded; new structured attempts begin from the migration
checkpoint.
