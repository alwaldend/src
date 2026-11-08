## 1. Native configuration migration

- [x] 1.1 Replace Lua configuration evaluation with native target composition and explicit label inputs.
- [x] 1.2 Verify generated configuration, nested plugin arguments, substitutions and runfiles through the independent consumer.
- [x] 1.3 Document the public API and validate the AL module and repository integration.

## 2. Runtime TLS correction

- [x] 2.1 Resolve generated TLS identifiers at runtime and verify local TLS requests from unrelated directories in directory and manifest modes.

## 3. Dictionary authoring and plugin ownership

- [x] 3.1 Implement whole-configuration dictionaries keyed by name, args payloads and argv executable argument vectors.
- [x] 3.2 Implement named args.res maps with structured env objects, explicit dependency order, selected-map overrides and missing/cycle diagnostics.
- [x] 3.3 Migrate repository configurations and shared dictionary constructors; move the external injector and schema to plugins/injector within AL.
- [x] 3.4 Verify all 33 generated configurations against the baseline, the real AL/plugin/child fixture, standalone module tests/builds and the independent consumer.
- [x] 3.5 Record final source acceptance and run the prepared-candidate quality, lint and local-review delivery gates. Git and AL receipts own publication state.

The built-in generic injection engine and resource-provider protocol are explicitly deferred.

## 4. Configuration reuse and command cleanup review

- [x] 4.1 Reuse the injector registration through a public configuration target and remove redundant plugin runtime declarations from consumers.
- [x] 4.2 Make the command cleanup fixture consume generated dictionary configuration targets, including its startup-failure override, with plugin runfiles supplied transitively.
- [x] 4.3 Verify success, startup failure and cancellation cleanup, inspect inherited runfiles and record source acceptance for delivery. Prepared-candidate checks and local review remain mandatory under al-delivery.

## 5. Target composition and plugin selection

- [x] 5.1 Provide structured plugin_labels selection in command wrappers and exercise value and presence-only selection through existing process fixtures.
- [x] 5.2 Provide reusable environment-export targets over named structured resource objects and verify file/environment injection through the existing cleanup fixture.
- [x] 5.3 Keep complex fixture graphs in local test definitions, document small target composition and validate standalone and repository consumers.
- [x] 5.4 Resolve explicit per-target extensions independently of aggregate order and reject ambiguous sibling definitions through the independent consumer.

## 6. Proportional fixture declarations

- [x] 6.1 Inline trivial fixture wrappers, retain the complex selected-resource fixture, and document a benefit-based extraction threshold.
- [x] 6.2 Verify unchanged generated configuration and real-process fixtures.

## 7. Flat named injector call groups

- [x] 7.1 Add al_injector_call with named `res_<type>` maps for all resource kinds; flatten per-resource dependency and Vault metadata and derive the default call identity from the target.
- [x] 7.2 Replace the narrow environment helper and generic injector payload boilerplate in production and fixtures.
- [x] 7.3 Verify all kinds, explicit metadata, duplicate-name rejection, structured settings and call identity through the independent consumer; preserve generated config and actual-process behavior.
- [x] 7.4 Update current API documentation and complete module/integration delivery checks.
