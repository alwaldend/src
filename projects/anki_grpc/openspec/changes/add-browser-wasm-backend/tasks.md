## 1. Scoped extraction and browser acceptance

- [x] 1.1 Record ownership, prerequisite boundaries and behavioral failure cases before adding implementation.
- [x] 1.2 Write independent Chromium E2E source and target using production WASI hosting and generated clients before restoring browser implementation.
- [x] 1.3 Extract browser-only production sources and owned shared-file hunks from the immutable source; omit both prerequisite BUILD files.
- [x] 1.4 Update browser documentation and validate OpenSpec and formatting.

## 2. Historical composition and delivery

- [x] 2.1 Validate the earlier startup contract in Node plus Chromium E2Es in composition with PR195 core `5e839d9611d293bbde51469234c7e1e539161cc7`; inspect JSON output and verify the matching PR2 test/runtime blobs against tested tree `c1e893c248d5cb3af30dea51f80b4f3cfc08b5df`.
- [x] 2.2 After PR195 merges, pass both E2Es, repository quality and affected semantic lint/build gates on the exact master-based PR2 candidate; composed behavior results do not complete this task.
- [x] 2.3 Publish the verified PR against master after PR1 merges; do not create a stack or mutate PR191.

## 3. Startup completion contract

- [x] 3.1 Record the joined-factory cancellation contract and adapt the existing delayed-Promise startup case before the core production change; preserve storage and non-startup cases.
- [x] 3.2 Update browser comments, README and specifications: canceled Start waits for factory completion, while commands retain cancellation and concurrent Stop/Done retain their cleanup behavior.
- [x] 3.3 Validate both browser E2Es with core `cd110f71f7a88187979a45037f94f3ac25cae0ff` in composed tree `de7883261b2bcf9478295217ff510b48eda0d107`; retain matching PR2 runtime/test blobs and passing reports. SDK browser E2E and native package verification also passed in that composition.
- [x] 3.4 Drain the delayed factory, pending Start and Stop/Done on every startup-case exit so a failed assertion retains its original error before the Go runtime exits.
- [x] 3.5 Retain the expected startup assertion in JSON with old core `5e839d9611d293bbde51469234c7e1e539161cc7` and cleanup fixture blob `c0e56eb31f6473be082c5c7742b719a4d2cb2b71`, before applying the simplified core for task 3.3.

## 4. Independent delivery review of the published Afero host

- [x] 4.1 Verify PR195 merged and rebase onto fetched master while preserving both native and browser startup documentation.
- [x] 4.2 Obtain a fresh independent review and add a real-browser concurrent media case before changing command execution.
- [x] 4.3 Reproduce and resolve reactor reentry during filesystem callbacks, then verify concurrent writes and existing browser cases.
- [x] 4.4 Complete fresh review of the correction and exact-candidate delivery gates before publication.

## 5. Caller-owned JavaScript filesystem

- [x] 5.1 Write failure cases before implementation for the production WASI host over Node fs: root mapping, flags/positions/errors, borrowed storage survival and descriptor cleanup on initialization failure, close and trap.
- [x] 5.2 Replace Go filesystem callbacks with Node-compatible synchronous filesystem calls; keep the command mutex and extract the separate Go timezone host.
- [x] 5.3 Add shared installGoFilesystem(fs) using the pinned Go runtime contract and adapt the generated-client harness to NewOsFs/BasePathFs over a caller-supplied root.
- [x] 5.4 Pass dependency-independent Node production-host and existing lifecycle E2Es; preserve JSON artifacts and exact source identities.
- [x] 5.5 Complete provider-backed Chromium and composed SDK acceptance after integrating the separately owned Emscripten provider; ZenFS and the npm `memfs` package are not selected; historical Afero-host results do not validate this revision.
- [x] 5.6 Record fresh independent review and the exact-candidate delivery plan. Repository quality, affected lint and consumer checks remain mandatory before updating PR196; the delivery receipt records their execution and publication.

## 6. Production-host I/O correctness

- [x] 6.1 Extend the existing real-reactor filesystem harness before changing production code; retain separate old-source failures for partial read/write progress, sync-after-write errors, final symlink lookup and recognized filesystem errors.
- [x] 6.2 Preserve completed transfer counts, report deferred sync failures, honor final-component no-follow before mutation and map concrete filesystem error codes without swallowing programming exceptions.
- [x] 6.3 Pass the production-host and existing Node suites, affected lint, configured formatting and strict OpenSpec; retain exact source/report identities and prepare one local candidate for parent-coordinated composition and fresh review.
- [x] 6.4 Before fixing close error precedence, retain failures for recognized sync plus recognized close errors before and after provider release, with unknown-close controls; preserve one-attempt ownership and rerun affected Node checks.

## 7. Standalone Emscripten MEMFS composition

- [x] 7.1 Preserve the build owner's failure-first generated-artifact E2E, verified SDK pin and hermetic Python/Node toolchain support.
- [x] 7.2 Selectively compose build candidate 15927fd and the scoped adapter manifest without inheriting b3 source, replacing shared BUILD files or adding SDK code.
- [x] 7.3 Write the explicit application-owned fs/root Chromium fixture before importing provider implementation, and retain the passing initial-source report separately.
- [x] 7.4 Integrate the independently reviewed adapter correction and rerun affected Chromium/provider acceptance with retained exact source and artifact hashes.
- [x] 7.5 Verify matching corrected-source composed SDK/native evidence.
- [x] 7.6 Reconcile the source inventory and evidence for delivery. Parent coordination released the hold; exact-candidate gates and publication remain controlled by the delivery receipt and Git state.
