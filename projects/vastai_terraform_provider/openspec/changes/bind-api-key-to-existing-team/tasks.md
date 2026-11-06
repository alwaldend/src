## 1. Existing-team credentials

- [x] 1.1 Write Terraform CLI failure fixtures before implementation and verify unsupported binding fails.
- [x] 1.2 Implement independent key binding with recoverable partial state and pass success, rejection, and missing-team fixtures.
- [x] 1.3 Run provider acceptance, owner validation, and record the actual live credential outcome.

## 2. Provider authentication

- [x] 2.1 Write E2E cases for provider-only authentication and removal of legacy issuer metadata without key rotation before implementation.
- [x] 2.2 Make issuer override optional and preserve key identity on issuer-only updates; pass lifecycle and migration fixtures.
- [x] 2.3 Remove redundant authentication from the ordinary test root and validate packaging and existing bootstrap behavior.

## 3. Imported module composition

- [x] 3.1 Exercise imported-team module reuse and managed bootstrap before implementation; defer permission validation only for unresolved computed team IDs.
- [x] 3.2 Pass the full Terraform CLI suite, retaining scoped-owner rejection once the ID resolves and preserving failed-key cleanup.
