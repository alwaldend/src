## 1. Dependency repair

- [x] 1.1 Define platform, extras, groups, cycles, version/source, and rootless
      failure cases before implementation and verify they reject the old importer.
- [x] 1.2 Initially repair platform reachability at the existing pin; verify
      isolated cases and consumers. Section 5 replaces that implementation.
- [x] 1.3 Verify the locked Rust source acquisition; retire the initial archive
      workaround through the Bazel 9 upgrade recorded below.

## 2. Acceptance

- [x] 2.1 Build the full site and graph WebP, inspect representative HTML and
      rendered diagrams, and complete the two consuming project records.
- [x] 2.2 Run formatting, semantic lint, and affected consumer checks, then
      hand the candidate to al-delivery for publication verification.

## 3. Reassess the Anki archive override

- [x] 3.1 Reproduce the original locked Git acquisition with clean isolated
      download and repository-content caches, including reuse across output bases.
- [x] 3.2 Remove the archive override and its generated projections if ordinary
      acquisition is sufficient; otherwise identify the smallest justified repair.
- [x] 3.3 Verify the Anki metadata, site graph, and site build after the decision,
      record the evidence, and hand the authorized revision to al-delivery.

## 4. Upgrade Bazel and retire the archive workaround

- [x] 4.1 Pin Bazel 9.3.1 with the publisher's SHA-256 and remove the Anki archive
      override, projections, generators, and formatting exception.
- [x] 4.2 Reconcile metadata and compatibility, keeping dependency versions
      unchanged unless separately required and authorized.
- [x] 4.3 Validate native Git crate acquisition, full-site loading/build, and
      both diagram formats; prepare the revision for the authorized existing PR
      update. Delivery receipts own exact commit, push, and PR state.

## 5. Replace the local Python importer patch

- [x] 5.1 Compare marker-qualified overrides, native environment partitions,
      and supported platform-specific compilation before choosing the replacement.
      Acceptance requires Linux to omit PyObjC, macOS to retain it, preserved
      locked versions and hashes, and reproducible generated artifacts.
- [x] 5.2 Implement the smallest supported patch-free workflow, remove obsolete
      importer code and fixtures, and document the owning update commands.
- [x] 5.3 Validate the actual autoscroll build, full-site capture and site test,
      representative diagram outputs, formatting, and semantic lint.
- [x] 5.4 Hand the completed candidate to al-delivery for the authorized
      update to PR 206; its receipts own publication state.
