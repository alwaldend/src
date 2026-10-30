## 1. Dependency integration

- [x] 1.1 Declare the pinned package and payload targets; build both and compare
      their bytes with the publisher download and its inner tar.
- [x] 1.2 Generate the root module include and reconcile its lock; inspect the
      diff and run the module freshness and buildifier checks.
- [x] 1.3 Document the consumer runtime interface and upstream provenance;
      verify extracted clamscan, freshclam, and clamd report the pinned version.

## 2. Acceptance evidence

- [x] 2.1 Record artifact hashes and isolated execution evidence, then validate
      the completed OpenSpec change with strict checks.
