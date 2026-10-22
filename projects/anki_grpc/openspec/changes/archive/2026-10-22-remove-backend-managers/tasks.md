## 1. Remove manager layer

- [x] 1.1 Confirm all manager call sites and migrate their required failure coverage to the paired collection-store change before deleting manager tests.
- [x] 1.2 Remove manager implementations and build declarations while retaining direct backend tests; verify no Go manager call sites remain.
- [x] 1.3 Update the library's current README/package documentation to describe explicit ownership and verify examples against the remaining API.

## 2. Verify integration

- [x] 2.1 Run direct backend/collection E2E plus migrated application scopes and inspect native/process artifacts.
- [x] 2.2 Validate both owner changes and deliver together in the independent manager-removal PR.
