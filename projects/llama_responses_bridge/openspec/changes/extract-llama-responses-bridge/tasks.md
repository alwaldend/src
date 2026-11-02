## 1. Extraction

- [x] 1.1 Move the command, adapter, and existing protocol suite into the project;
      verify imports, loopback defaults, and command-level acceptance.
- [x] 1.2 Document the project and register its owner; regenerate inventories and
      verify owner validation and repository quality.
- [x] 1.3 Prepare the separate master-based candidate; verify semantic lint, package
      builds, protocol acceptance, and the reviewed source diff.

## 2. Project landing page

- [x] 2.1 Add and register the bridge landing page with accurate status and
      documentation/source links; inspect rendered landing, index, docs, and taxonomy HTML.
- [x] 2.2 Verify the shared Hugo site, DNS checks, and owner quality gates;
      document that a source change does not deploy the site.

## 3. Compression review finding

- [x] 3.1 Reproduce negotiated compressed upstream replies with the existing
      command fixture before fixing the request path; verify JSON/SSE routing.
- [x] 3.2 Prevent client compression negotiation from reaching the translation
      upstream; verify JSON/SSE acceptance and record the review-fix evidence.
