## 1. Configuration migration

- [x] 1.1 Inventory configuration targets and capture baseline JSON before editing.
- [x] 1.2 Replace all AL Lua sources with typed targets and shared Starlark helpers.
- [x] 1.3 Update current documentation and procedures to the target-based flow.
- [x] 1.4 Compare all generated configurations and build representative consumers offline.
- [x] 1.5 Record acceptance evidence and the required source-delivery gates. Git and AL receipts own commit, push and PR status.

## 2. Review correction

- [x] 2.1 Resolve generated TLS runfile identifiers at runtime and verify local TLS requests in directory and manifest modes from an unrelated working directory.

## 3. Dictionary authoring follow-up

- [x] 3.1 Implement and migrate named configuration dictionaries, args payloads and dictionary resource graphs; verify output and runtime behavior. The injector contract is owned by projects/al/openspec/changes/native-configuration-targets.
- [x] 3.2 Remove redundant wrapper plugin dependencies in favor of inherited configuration runfiles; preserve all generated configurations and verify affected command consumers. AL's change record owns the cleanup fixture and acceptance evidence.

## 4. Reusable Terraform integration and namespaced APIs

- [x] 4.1 Centralize AL Terraform map integration and Vault backend call construction without changing standalone rules_terraform dependencies.
- [x] 4.2 Migrate callers to direct configs/run_args and namespaced infrastructure helpers, and record namespace requirements in repo-bazel.
- [x] 4.3 Verify configuration outputs and operation contracts, execute isolated Terraform wrapper fixtures, and record source acceptance for final delivery checks.

## 5. Composable configuration declarations

- [x] 5.1 Replace component-wide configuration dictionaries with small authentication, connection, plugin and call targets composed through configs.
- [x] 5.2 Keep complex component graphs and templates in local defs.bzl macros; declare common patterns through reusable namespaced target macros.
- [x] 5.3 Replace plugin-selection CLI lists with plugin_labels dictionaries expanded by the owning AL wrapper.
- [x] 5.4 Update the skill and documentation, preserve generated configuration and runtime dependencies, and verify source acceptance for delivery.

- [x] 5.5 Migrate the remaining XO bootstrap Lua configuration to native targets, retaining runtime-only approved-role selection and its existing validation coverage.

## 6. Proportional configuration extraction

- [x] 6.1 Inline single-use local wrappers whose small declarations are clearer in BUILD; preserve genuinely complex definitions and reusable common helpers.
- [x] 6.2 Refine locality guidance and evaluation cases to require a significant readability or reuse benefit.
- [x] 6.3 Compare generated configurations and complete source-delivery validation and review.

## 7. Flat injector declarations

- [x] 7.1 Migrate injector calls to named `res_<type>` maps while retaining small-target composition and proportional extraction.
- [x] 7.2 Preserve generated definitions, runtime identities, labels, dependencies and credential references; update guidance and deliver the validated aggregate.
