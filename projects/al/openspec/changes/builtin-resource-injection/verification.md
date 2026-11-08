## Source acceptance

The source implements all requirements and their scenarios. The core graph
and generic resources live in `pkg/injector`; `pkg/al_plugin` discovers advertised
capabilities and dispatches the separate resource RPC. Vault schema and fetchers
live in `plugins/vault`. The former external injector package is removed.

| Requirement                 | Verification                                                                                                                                                                                                                                                                                                          |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Generic built-in evaluation | `tests/resources:resources_test` runs AL with built-in-only calls and an invalid, unselected provider executable. It also verifies actual helper-process termination.                                                                                                                                                 |
| Advertised providers        | The same actual-process fixture serves a non-Vault type over RPC, validates dependency file transport and preserves nested data and integers above 2^53. It rejects missing/cyclic graphs, missing types, duplicate/built-in collisions and empty advertisements, and rejects advertisements from ineligible plugins. |
| Shared resource results     | `TestStartupResourceResults` verifies absent and full env/files/data startup results over an actual plugin process. Types arrive in the startup envelope independently of optional output; fetches return the generic result directly and invalid advertisements roll back the provider.                              |
| Invocation-owned cleanup    | The fixture covers startup/fetch failures, provider death and cancellation, retains files/providers until child cleanup finishes and verifies final file/process removal. Root command fixtures additionally cover success, startup failure and cancellation.                                                         |
| Vault separation            | `tests/injector:dictionary_e2e_test`, `plugins/vault:vault_test` and the root `tools/al:vault_dev_server_test` verify selected resources, overrides, dependency order, lazy authentication, OIDC boundaries and the pinned local Vault server.                                                                        |

The startup process fixture distinguishes absent responses, responses without a
result, explicitly empty results and full env/files/data results. Startup metadata
has its own envelope and carries the provider type advertisements alongside optional output.
The same fixture verifies generated protobuf settings validation, a legal external
plugin named injector, and explicit injection overriding conflicting startup defaults.
Dedicated InjectorCall protobuf entries and the dictionary API are checked by the
independent consumer.

Generic resource input is verified by real-process scenarios for local file and
environment template data, typed/nested provider input, separate dependency
results and rejection of removed extra settings. The dictionary process fixture
also checks the Vault operation write body through its local HTTP server. The
independent consumer checks flat and dictionary data declarations, nested values
and integer precision in generated configuration.

Root-map process scenarios verify an override whose key sorts before its parent,
unselected ancestors, equal definitions, conflicting siblings, missing/cyclic
lineage, empty identities and rejection of legacy lists. Startup verifies plugin
and call keys over RPC. The independent consumer checks all five root maps and
merges groups with overlapping local resource keys without mixing their settings.

The provider fixture writes `resource-provider-results.json` to Bazel undeclared
test outputs. It also writes `startup-resource-results.json` for the optional
startup-result scenarios. The root command fixture writes success, startup-failure and
cancellation receipts. These are repeatable acceptance artifacts using synthetic
data and local processes.

## Validation scope

- From `projects/al`, `bazel test //...` passes all 24 normal tests and
  `bazel build //...` builds all 124 normal targets.
- From the root, test `//:repo_quality_test`, `@rules_al//...`, `//tools/al/...`,
  `//tools/vault/...`, `//infra/xcp_ng/cmd/xo_login:all`,
  `//infra/src/openspec/validation:all`, the changed infrastructure/secrets/Bazel skill
  `eval_config_test` targets and `//.agents:write_skills_test`.
- Run `bazel build --config=lint` for `@rules_al//...`, `//tools/al/...`, the
  changed skill libraries, `//infra/github:all`, `//tools/vault/...` and
  `//infra/xcp_ng/cmd/xo_login:all`.
- Build all 33 repository AL config targets. Generated output inspection confirms
  generic resource definitions and Vault provider registration; no external
  injector registration remains. Built-in-only fixture configs carry no provider.
- Run `@rules_al//tests/consumer:check` through its documented external launcher.
  Its independent workspace verifies configuration, extensions, flat resource
  authoring, runtime files, directory/manifest TLS and payload bytes/digest.
- Apply the configured Go, Starlark, protobuf, Markdown and YAML formatters.
  Strict OpenSpec validation passes for this AL change.

Prepared-candidate gate logs, local review, transport verification and publication
receipts are retained under ignored `out/builtin-injector/`. Normal generated
standalone Bazel convenience links are disposable outputs and are removed before
the parent workspace recursively scans the module. No live infrastructure or
credential operations form part of this validation.

Offline skill eval validation proves staging/configuration only; no model-backed
behavioral eval was run. Capability discovery starts every eligible provider when
an external type is needed, and providers must keep startup free of secret reads.
