## 1. Provider contract and implementation

- [x] 1.1 Write failure-oriented provider E2E cases before adapter implementation.
- [x] 1.2 Translate Node callback and Sync APIs onto the generated classic MEMFS,
      preserving upstream ownership and reporting unsupported behavior.
- [x] 1.3 Package adapter and generated assets in the consumer files target.

## 2. Verification and handoff

- [x] 2.1 Run the standalone E2E against the real generated module and immediately
      retain its JSON/binary artifacts with exact input hashes.
- [x] 2.2 Validate formatting and OpenSpec, document observed limitations and hand
      a scoped local candidate to the parent for composition and browser acceptance.

## 3. Independent review corrections

- [x] 3.1 Extend the real generated-module E2E before production corrections and
      retain failing evidence against adapter blob `417aec826d72e599b2b4faea6e660d4b9a41fba3`.
- [x] 3.2 Correct exclusive symlink creation, shared factory options, write extent
      errors and read/write callback failure tuples; verify all added scenarios.
- [x] 3.3 Retain corrected artifacts and immutable source projections for parent
      composition, preserving previous integration evidence as historical.
