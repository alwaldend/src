## 1. Protocol and acceptance fixtures

- [x] 1.1 Extend real-process fixtures for a non-Vault provider, advertisements, dependencies and startup/fetch/graph errors before implementation.
- [x] 1.2 Add generic resource-fetch RPC and provider advertisement/config fields.

## 2. Built-in engine and providers

- [x] 2.1 Move generic resource evaluation and lifetime handling into AL without Vault schema imports.
- [x] 2.2 Integrate lazy provider startup, capability registration, dispatch and ordered rollback/shutdown.
- [x] 2.3 Move Vault resource implementations into a standalone provider and migrate registrations/fixtures.

## 3. Verification and delivery

- [x] 3.1 Verify generic provider, builtin-only, Vault and cleanup end-to-end receipts.
- [x] 3.2 Update API/owner documentation and validate root and standalone builds/tests/format/lint.

Exact-candidate review, publication and remote feedback state are recorded in
task-local al-delivery receipts rather than duplicated in source checkboxes.

## 4. Direct startup results

- [x] 4.1 Extend actual-process fixtures for optional startup resource results and separate capability discovery before implementation.
- [x] 4.2 Return ResourceResult directly from startup and fetch RPCs; move advertisements to ResourceTypes and migrate all first-party implementations.
- [x] 4.3 Validate result transport, cleanup, root/standalone builds, tests, formatting and semantic lint; update protocol documentation.

## 5. Extensible startup response

- [x] 5.1 Extend process fixtures before implementation for absent, explicitly empty and full startup results inside a response envelope.
- [x] 5.2 Add PluginStartResponse.result, retain the envelope in the client and migrate all first-party startup implementations.
- [x] 5.3 Update protocol documentation and verify root/standalone builds, tests, formatting, semantic lint and process receipts.

## 6. Human review corrections

- [x] 6.1 Extend process fixtures for startup advertisements, dedicated injector configuration, protobuf settings failures, legal plugin names and explicit environment overrides before implementation.
- [x] 6.2 Define generated protobuf settings and typed injector call/resource declarations; migrate authoring, engine and configuration consumers.
- [x] 6.3 Advertise resource types in PluginStartResponse and remove the separate discovery RPC.
- [x] 6.4 Let explicitly injected environment values override startup exports without depending on plugin startup order.
- [x] 6.5 Verify process/consumer outputs, root and standalone build/test/lint gates and complete local review before publication.

- [x] 6.6 Use plugin_configs for declarations, plugins for invocations and inject for builtin resources; align the public Bazel constructors and all consumers.
- [x] 6.7 Retain one resource-service client per plugin connection and verify the existing process fixtures.

## 7. Generic resource data

- [x] 7.1 Extend real-process and independent-consumer fixtures for local template data, provider transport, Vault operation payloads and removed extra settings before implementation.
- [x] 7.2 Add generic resource data across configuration, templates and provider RPCs; migrate the kubeconfig template and remove file-specific extra.
- [x] 7.3 Verify root/standalone tests, builds, formatting, semantic lint, generated configs and the external consumer; review the candidate for publication.

## 8. Root configuration maps

- [x] 8.1 Extend process and public-consumer fixtures for named root maps, explicit injection extensions, ambiguous collisions and invalid lineage before implementation.
- [x] 8.2 Convert all five root sections and selected plugin calls to maps; remove redundant declaration names and migrate every consumer.
- [x] 8.3 Resolve injection resource overrides through explicit ancestry, independently of map iteration or formatter order.
- [x] 8.4 Verify builds, tests, formatting, lint, generated configs and consumer artifacts; review the candidate for publication.
