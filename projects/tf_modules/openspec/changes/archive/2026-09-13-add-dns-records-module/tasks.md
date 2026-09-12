## 1. Canonical declarations

- [x] 1.1 Implement normalization and verify all supported types, destination deduplication, stable identities, multiplicity, and invalid-input rejection with offline Terraform tests.
- [x] 1.2 Verify compatibility TTLs and normalized names against migration parity fixtures.

## 2. Provider resources and packaging

- [x] 2.1 Map Cloudflare and RouterOS scalar resources, verify default-disabled behavior and owned-resource changes with provider mocks.
- [x] 2.2 Package source and pinned test providers, then pass Terraform formatting, Bazel package checks, and strict OpenSpec validation.
