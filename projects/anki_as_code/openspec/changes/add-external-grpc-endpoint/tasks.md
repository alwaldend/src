## 1. Behavior coverage

- [x] 1.1 Add native CLI E2E scenarios before production edits for embedded defaults, TCP/Unix selection, state operations, export/apply/plan/build, invalid addresses, cancellation and server reuse; observe the unsupported flag failure. Extend endpoint syntax checks during E2E refinement.

## 2. Implementation

- [x] 2.1 Add optional `col --address` and injected application construction; configure network or embedded transports through the existing collection owner and verify CLI E2E artifacts. Preserve listener reuse by ignoring SIGPIPE in the executable's native adapter.
- [x] 2.2 Document address syntax, defaults, ownership and shared filesystem requirements; inspect help and documentation against the candidate.

## 3. Validation and delivery

- [x] 3.1 Run affected native/CLI workflows, semantic lint, repository quality and strict OpenSpec validation; inspect repeatable collection and state outputs.
- [x] 3.2 Inspect the final diff and prepare scoped delivery inputs and the validation plan; delivery receipts own the commit, push and pull-request state.

## 4. Rebase integration and application ownership

- [x] 4.1 Rebase onto the merged API CLI and inspect the complete merge diff and affected consumers.
- [x] 4.2 Exercise per-call endpoint selection on one App and native API `--address` before production edits.
- [x] 4.3 Remove the CLI application factory; pass address and namespace options into App workflows and preserve merged archive publication.
- [x] 4.4 Integrate `api --address`, refresh documentation/specification, verify native artifacts and CLI documentation, and select final lint and quality gates.
- [x] 4.5 Review the corrected candidate and prepare scoped exact-lease delivery inputs; delivery receipts own validation and publication state.

## 5. Unified collection options

- [x] 5.1 Adapt native managed/borrowed workflows before implementation to use one options value, and cover conflicting access modes, typed nil connections and archive build forwarding.
- [x] 5.2 Replace CollectionSource and source constructors with Path/Address/Connection/StateKey fields on CollectionOptions; migrate App, APIRequest, CLI and consumers, preserving native ownership and publication.
- [x] 5.3 Refresh current documentation and specifications, inspect native artifacts and prepare scoped validation/delivery inputs; receipts own final gate and publication state.
