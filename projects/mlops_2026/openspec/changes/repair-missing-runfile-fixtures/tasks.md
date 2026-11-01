## 1. Runfiles fixture repair

- [x] 1.1 Derive each missing-CSV manifest key with the public directory-based runfiles resolver and retain the existing startup assertions and deadlines.
- [x] 1.2 Run the complete recommender E2E target and inspect its JSON artifact for both missing-CSV cases, missing/corrupt model checks, and normal HTTP behavior.

## 2. Validation and delivery

- [x] 2.1 Format, run scoped semantic checks, validate this change, and include the repair in the full repository build/test matrix; select guarded PR delivery checks, with receipts owning publication results.
