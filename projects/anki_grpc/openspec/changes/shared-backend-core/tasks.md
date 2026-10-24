## 1. Shared implementation

- [x] 1.1 Extract directclient, shared C ABI and native adaptation; verify the
      source manifest and existing generated-client/native compatibility checks.
- [x] 1.2 Extract timezone calculations and WASI delegation; verify worker
      timezone and host fixtures, and exclude browser implementation and feature docs.

## 2. Integration and delivery

- [x] 2.1 Validate managed collection and application workflow consumers;
      inspect their repeatable output artifacts and record results.
- [x] 2.2 Resolve independent bounded review and validate the specification;
      define exact-candidate compatibility, root quality and affected semantic lint
      publication gates. Guarded delivery receipts record final validation and
      publication against master under task scratch.

## 3. Synchronous startup refinement

- [x] 3.1 Update the existing native E2E before production edits and demonstrate
      that cancellation must wait for initialization while Stop remains bounded.
- [x] 3.2 Inline synchronous initialization into Start, preserve lifecycle and
      error ownership, and pass native and consumer compatibility checks.
- [x] 3.3 Refresh documentation and independent review, and select the guarded
      exact-candidate validation plan. Receipts own final gate results;
      publication remains held for parent-owned browser validation.
