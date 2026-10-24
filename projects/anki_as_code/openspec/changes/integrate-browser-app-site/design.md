## Context

The `add-browser-app` change provides a reusable worker API and standalone SDK example. Project landing content already flows through `site/content:docs` into the apex Hugo site. Its existing `landing` layout and Bootstrap components supply all required presentation.

## Goals / Non-Goals

**Goals:** Provide a visitor-facing configuration and collection workspace,
discoverable from the project landing, using existing application and theme
owners. Let visitors explicitly save completed work locally and continue from
their last saved workspace after reload.

**Non-Goals:** VS Code or another editor dependency, offline installation,
automatic background saving, a parallel design system or live deployment.

## Decisions

- Add project-owned `site/content/app/index.md` with `layout: landing` and accessible form markup. Reuse Bootstrap classes and the existing theme; introduce no stylesheet, shortcode, site layout or iframe.
- Add the page controller at `assets/browser/app.js`, importing `./index.js`. It owns DOM interaction only and calls the existing worker API. Automatically load the example and select its note so visitors can edit immediately. Show one file selector and TOML editor with Open JSON, Check and Download JSON actions. Import exports JSON to the virtual workspace directly; checks and downloads read the current configuration through Go. This removes the separate JSON editor, generation step, save step, output format selector and dump panel.
- Keep TOML as the editing surface because it is the project's source format. A JSON-only editor would expose the serialized envelope; a form editor would require another representation of the configuration schema. Explicit Check gives feedback when requested without per-keystroke scheduling or errors during incomplete edits. It checks syntax, schema types and filesystem layout through `configDump`; success describes format validity, not collection reconciliation.
- Split the SDK's existing package mapping into runtime files and its standalone example. A project-owned filegroup maps runtime files plus the page controller under `content/projects/anki_as_code/app/` as resources of the page's leaf bundle, excluding standalone `index.html` so Hugo owns the page. Add that mapping to the landing's `docs_filegroup` dependencies; the registry requires no change. The leaf bundle uses the existing content mount; the site's explicit webfont static mount suppresses Hugo's implicit default static mount.
- Serialize page actions, show readiness/errors in a live status region and disable conflicting controls while an operation runs. Write the current editor text before switching files, checking or downloading, so visitors need no separate action to update in-memory files. Retain malformed text for correction and validate the whole workspace before download. Replacing an existing workspace with JSON uses a new app instance, swapped only after successful export/load, to preserve existing work on invalid input and avoid stale undeclared files.
- Keep the SDK's existing offline E2E and add a separate test of the generated site. Reuse the pinned site browser workflow, allowing only the site's existing SRI-pinned CDN scripts. The app runtime and file operations require no external service.
- Add Open collection, Plan, Apply and Download collection controls using the
  SDK's collection methods. Open collection accepts an Anki package and makes
  its exported configuration available to the existing TOML editor. Stage a
  replacement workspace before swapping it into the page, preserving current
  work on import failure. Plan shows the existing application's change report;
  Apply uses the existing reconciliation workflow. Download collection returns
  the completed package from the workspace. Check retains its configuration
  format meaning and remains independent of collection reconciliation.
- Add explicit Save locally using the SDK's atomic IndexedDB snapshots. The
  page invokes the SDK's restore method at startup to continue the last saved
  workspace, or loads the existing example when no snapshot exists. Flush
  current editor text into the workspace before saving and
  wait for any active collection work to complete. Show success only after the
  snapshot transaction commits. Restore stages a complete replacement before
  changing the visible workspace. A restore failure reports the error without
  overwriting the saved snapshot and leaves configuration editing available.
  Saving does not apply configuration edits to
  the collection implicitly; the saved package reflects the last completed
  collection operation and the saved config files retain current editable text.
  Save locally also supports the configuration-only editor before a collection
  is opened, without loading the server.
- Keep Open JSON and Download JSON available for portable configuration work.
  Before retaining a collection, export and validate the candidate configuration
  and reject any files outside `/workspace/`, including reserved `/tmp` storage. Inspect exported paths
  rather than duplicating the protobuf schema in JavaScript. This prevents JSON
  paths from creating or replacing a collection package; failed candidates leave
  the active worker and editor untouched. The SDK retains its full virtual
  namespace for callers that explicitly manage their own files.
  Imports prepare a fresh memory workspace without reading IndexedDB. Retain
  the page's last observed saved revision across imports and pass it explicitly
  when saving, so replacement never adopts another tab's unseen revision.
  Use one staged-workspace helper for import, restore and initial setup, with
  one success swap and one failure-disposal path.
  Distinguish unsaved in-memory edits, completed local snapshots and downloaded
  files in concise status/help text. Report unavailable browser capabilities,
  storage errors and conflicting saves without losing current edits. Browser
  storage is specific to the site origin and can be cleared; downloads remain
  the portable way to retain work.

## Risks / Trade-offs

- Static demo HTML could override Hugo output → exclude it from site assets and verify generated app HTML includes the shared site shell.
- Unsaved work could be lost on reload → explain the explicit Save locally
  boundary and retain configuration JSON and collection package downloads.
- Switching files could discard text or hide validation errors → commit editor contents to memory before switching, retain the current selection on read failures and validate the whole workspace on Check and Download JSON.
- Collection or storage failures could replace good work → stage imports and
  restores, retain the editor on failures and delegate atomic save/conflict
  behavior to the SDK.

## Acceptance Cases Before Implementation

Navigate via Try it now and edit the automatically opened Unicode note. Switch files and return without losing the edit; check and download JSON without a separate save step. Reopen the download in one action. Invalid JSON, schema/layout errors and malformed TOML must preserve existing work and allow correction. Verify the actual shared header, stylesheet, body typography/colors, theme preference and mobile overflow; save compact desktop light/dark/mobile screenshots and a structured report. Build the site and run its existing build/DNS checks. No deployment is authorized.

Before collection UI implementation, extend the rendered-site E2E to open a
collection with binary media, edit TOML, inspect a plan, apply, download the
completed collection, save locally and restore after reload. Verify the restored
editor contents and collection result, failed import/restore/save feedback and
retention of current edits. Recheck the shared theme and narrow-screen layout
with the added actions. The SDK's E2E owns detailed transaction conflict, quota
and native Anki package verification; the site test proves users can reach and
complete those workflows through the rendered page.

Reject JSON exports into `/tmp` with and without an open collection. Preserve
the current editor and exact package bytes through a subsequent save and reload.

Restore a saved workspace with more files than JavaScript permits as function
arguments. Every file must remain selectable, edits must reach the restored
workspace, and the previous worker must be disposed. Prepare the selector in a
detached document fragment before switching the active workspace so option
construction failures preserve the current editor, worker and saved revision.
