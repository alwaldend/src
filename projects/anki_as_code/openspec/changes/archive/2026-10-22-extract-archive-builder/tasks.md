## 1. Extract

- [x] 1.1 Write native archive E2E coverage for success, no-op, callback failure,
      cancellation, source aliases, invalid paths and staging cleanup before copying the implementation.
- [x] 1.2 Extract the standard-library-only Builder and wire its test visibility.
- [x] 1.3 Document the API and its callback, lifetime and publication guarantees.

## 2. Verify and deliver

- [x] 2.1 Run the native E2E, inspect archive/state artifacts, and review the complete diff.
- [x] 2.2 Pass repository quality and affected-target semantic lint for the independent
      master-based candidate. Delivery receipts own publication state.

## Evidence

The native E2E passes with the master-based standalone Builder. It emits base,
changed and unchanged collection packages plus the state read from the built
package. Native reads confirm source isolation and successful state publication;
the unchanged package preserves exact bytes. Failure, cancellation and panic
paths remove staging and preserve the existing output. Logs and artifact
summaries are under ignored `out/extract-archive-builder/`.

The first fixture attempt lacked Anki's empty media directory; creating that
directory before native export corrected the fixture without production changes.

## 3. Refine package naming

- [x] 3.1 Rename the implementation package and Bazel targets to rebuild, update callers and documentation, and remove the source ignore exception; verify the native E2E through the renamed target.

## 4. Preserve cancellation causes

- [x] 4.1 Extend the native E2E to check custom cancellation causes before startup
      and after successful native work; reproduce the missing cause before fixing it.
- [x] 4.2 Preserve both the context error and custom cause at each cancellation
      check, rerun the native E2E and select exact-candidate delivery gates.
