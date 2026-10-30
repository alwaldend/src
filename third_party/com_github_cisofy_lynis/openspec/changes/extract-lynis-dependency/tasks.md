## 1. Dependency extraction

- [x] 1.1 Preserve the Lynis package and generated module include from source
      commit `5e4bf7ccb7e7948e12c12ae12b92d048bbba06e3`; compare Git blobs to
      verify the BUILD interface and dependency pin are unchanged.
- [x] 1.2 Keep the package documentation independent of the consumer role;
      verify its local links resolve within the extracted candidate.
- [x] 1.3 Build and inspect the upstream archive, verify its SHA-256 against
      the unchanged SRI pin, and run the extracted executable's version check.

## 2. Integration validation

- [x] 2.1 Regenerate root module includes, reconcile the lock, and verify that
      only the two requested dependency includes differ from master.
- [x] 2.2 Validate the owner record with strict OpenSpec checks and retain
      fresh dependency validation evidence for the split candidate.
