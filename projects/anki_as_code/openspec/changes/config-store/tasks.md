## 1. Store API

- [x] 1.1 Adapt existing configuration/export/state E2E consumers to Store before implementation, including reuse for export and load; verify the existing behavior assertions remain intact.
- [x] 1.2 Implement the root-bound Store and migrate production callers; verify all consumers build and existing config/native E2E workflows pass.

## 2. Documentation and delivery

- [x] 2.1 Update package and README API documentation and diagram labels; verify maintained diagram renders and include the integrated site build in the delivery checks.
- [x] 2.2 Inspect the complete scoped diff and select required exact-candidate quality/lint, affected E2E and site-build gates. Delivery receipts own validation and publication facts.

Verification: the existing config, CLI, app, state and native exporter E2E targets pass with Store-based callers. Rendered operation/export diagrams show the Store boundary; no serializer or layout implementation changed. The existing export/load workflow reuses one store and preserves exact values and source objects. Evidence lives under `out/config-store/`.
