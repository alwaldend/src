## Acceptance evidence

The file-only `//projects/anki_as_code/pkg/ankiascode:filesystem_test` workflow
passes with absolute, relative and literal-home paths in memory storage. It
verifies nested deck discovery, default references, exact Unicode content,
JSON/TOML output, independent applications, fresh reads, missing/malformed
inputs, read-only output and absence of host output. Its undeclared outputs
contain repeatable JSON/TOML configuration artifacts.

The native `TestInjectedCollectionConfiguration` workflow passes with every
configuration and state-output file in memory. It exports, imports bindings,
plans and applies an edit, then verifies convergence and absence of disk output.
The artifact `memory-applied.note.anki.toml` preserves the edited declaration.
This workflow found and verified the fix for `AFERO-RENAMED-HANDLE-CLEANUP`:
memory-backed file handles change their reported name after rename, so cleanup
must retain the original staging pathname.

The config, export, collection, rebuild, plan, reconcile, state and CLI suites
passed during focused validation. The memory build covers success, apply
failure, cancellation, original/output preservation and staging cleanup.
Existing OS coverage includes symlink/hardlink conflicts and publication failures.
Strict OpenSpec validation and the WASI command build pass.

A source scan of App, config/layout, collection, export, rebuild and filesystem
found no direct OS file or directory operations. Host home/current-directory
expansion is confined to the OS adapter. CLI input adapters and backend/server
storage are independently owned and retain their existing contracts.

Final exact-candidate tests, quality/lint results, commit identity and publication
state are owned by the local `repo-delivery` receipts under
`out/ankiascode-afero/`; those receipts are intentionally not checked in.

## Review follow-up

PR #189's completed automated review identified `AFERO-SOURCE-STATE-PATH`:
OS metadata identity does not recognize identical paths in memory storage.
The added absolute/relative source-conflict workflow first failed because the
conflict reached backend startup. Comparing resolved source and state-output
paths directly makes those cases pass, preserves source bytes and creates no
configuration output. OS metadata comparison remains for hardlink protection.
The updated candidate uses `prepare-review.json` and its delivery validation
record under the same ignored task directory.

The next completed review identified `AFERO-BASEPATH-LINK-NAMESPACE`:
BasePathFs exposes symlink targets in backing-store coordinates. Build and
export workflows first reproduced doubled roots and misplaced output. The
resolver now uses the public Afero base-path utility to translate targets into
virtual names, preserving relative caller names and directly nested bases.
Relative link contents are interpreted against the link's backing directory;
absolute and relative targets outside the base are rejected.

Focused build/export workflows pass for absolute and relative caller names,
nested bases, ordinary relative symlinks, dangling state outputs and rejected
external links. Config input symlink rejection remains enforced. Repeatable
artifacts include `basepath-build-*.json` and `basepath-export-state.toml`.
The updated candidate uses `prepare-basepath.json` and its validation record.

The complete PR review identified `AFERO-BASEPATH-CONFIG-ESCAPE` and
`AFERO-NESTED-TEMP-NAME`, confirming the next automated review findings.
Configuration roots and declared resource paths escaped through parent traversal
or existing symlinks; nested state publication reused a backing handle's pathname.
Each case first failed its end-to-end workflow. All owners now share a
`filesystem.FS` with a private Afero backend and guarded I/O. Temporary files
return a stable virtual pathname separately from their handle.

The complete review also reproduced `AFERO-ROOTED-SOURCE-IDENTITY` and
`AFERO-SYMLINK-PARENT-ORDER`: mixed absolute/relative rooted names overwrote a
memory source, and parent components following links selected the wrong source.
Canonical rooted names protect source identity; component-by-component resolution
preserves native link semantics for both link contents and direct builder inputs.

Review of the shared wrapper caught and reproduced `AFERO-OVERLAY-OPEN` and
`AFERO-EXPORT-PATH-COORDINATES` before publication. File reads now preserve the
backend's `Open` semantics, including merged overlay directories, and the
independent exporter resolves both destinations before comparing them.
Focused checks pass for configuration, rebuild and independent export, including
native workflows, with no application I/O helpers bypassing the shared wrapper.
Artifacts include `overlay-configuration.json`, rooted/nested state TOML files,
and normalized configuration JSON. Strict OpenSpec validation also passes.
The final candidate uses `prepare-centralization.json` and its validation record.
The final independent read-only review of the complete candidate reported no
remaining actionable findings after checking the fixes and their behavioral cases.

## Session ergonomics

- `BAZEL-CACHED-CRATE-ABSOLUTE-SYMLINK` (live): native checks initially failed
  because one cached Rust crate's manifest linked into another unavailable
  output base. A bounded `bazel_agent bazel fetch --force --repo=...` refresh
  restored the manifest; no dependency pin or infrastructure change was needed.
- `DISCOVERY-OUTPUT-BOUND` (live): the first worktree inventory produced 1,671
  lines and exceeded the requested context. Subsequent discovery used bounded
  owner-local reads. Query the supplied checkout before requesting inventories.
- `CONSTRUCTOR-REWRITE-SCOPE` (live): broad textual constructor edits introduced
  invalid call arguments; compilation and diff inspection caught and corrected
  them before acceptance. Prefer scoped patches or structure-aware rewriting
  when call expressions contain nested calls or callbacks.
- `REVIEW-CAPABILITY-SEMANTICS` (live): whole-PR review found that delegating
  `Open` to `OpenFile` changed overlay directory behavior. Preserve each backend
  operation's semantics when centralizing calls; a memory-only implementation is
  insufficient evidence for layered filesystem compatibility.

These observations remain task-local evidence. No shared policy or tool behavior
was changed as part of this filesystem migration.
