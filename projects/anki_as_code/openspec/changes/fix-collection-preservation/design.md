## Context

The filesystem already resolves path components in order, but some public
callers cleaned paths first. Planning already validates native ownership. Export
writes a caller-selected hierarchy and does not own the directory's previous
contents. The transport client grew collection policy that belongs here.

## Goals / Non-Goals

Preserve correct source selection and required native identities. Restore
caller-owned export destinations and one application-owned collection scope.
Do not protect caller files by passing source identities through generic writers,
add generic path-planning APIs, or retain reusable transport collection phases.

## Decisions

- Resolve public paths without lexical cleanup before following symbolic links.
  Keep rooted namespace boundary checks and configuration-loader link rejection.
  For absent build outputs, resolve the original parent components in order and
  retain existing missing-parent and dangling-final-link behavior.
- Export directly writes declared files and leaves unrelated entries alone.
  Callers select output locations and are responsible for existing files and
  aliases there. Remove protectedSources and source/output preflight checks
  instead of moving them into another layer. Preserve ordinary output errors
  and consistency between generated configurations and optional state output.
- Reject omitted tracked deck ID 1 during planning ownership checks.
- Store creates a raw transport client, then prepares one private collection
  scope. ColOpen returns the shared ankiclient.Client interface exposing
  Connection and Stop. The scope owns native storage, optional archive publication
  and cleanup; its transport implements StartableClient, which embeds Client and
  adds Start. No separate public Collection interface is needed. Keep only
  synchronization required for concurrent native calls, terminal cleanup and
  source reservations, with no public restart or collection-phase API.
- Inject client constructors as functions. Observe concrete transport Done
  methods through a private Done-only interface in the collection package.
  Keep server command argument adaptation in CLI: it prefixes serve before
  calling ServerRun, so App does not need a separate ServerServe method.
- Audit every handwritten exported declaration against its actual consumers.
  Use ankiascodeapi.ColOpenOptions directly instead of a second collection
  options alias. Remove the test-only filesystem.FS.ExpandPath convenience
  method and keep resource.Store's package-local Search and Presets service
  fields private; Snapshot.Presets remains caller-visible data. Replace
  codec.Format, which has no methods or validation of its own, with string
  parameters and untyped JSON/TOML constants while retaining runtime format
  validation. Keep CollectionCloseTimeout because it configures real cleanup
  policy. Preserve public application contracts and declarations needed
  for cross-package use, framework integration or serialization. These visibility
  and alias reductions require no new runtime behavior specification.

## Risks / Trade-offs

Export overwrites declared destinations and can leave partial output on I/O
failure. The caller must choose an appropriate directory and state file.
Cancellation cannot establish native completion: application-owned restore files
and reservations remain retained if release cannot be confirmed. This does not
expand the generic transport Client interface.
