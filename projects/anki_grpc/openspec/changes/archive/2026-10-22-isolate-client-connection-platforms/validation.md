## Scope

This refactor targets master independently of the in-process feature in PR #184. Client owns lifecycle and exposes grpc.ClientConnInterface; NewClient takes only an address. Optional native dependencies require a cgo subpackage import.

## Observed verification

- All 15 tests in projects/anki_grpc and projects/anki_as_code pass, including native persistence, TCP, Unix, collection publication, application workflows and the existing WASI SQLite probe.
- Core, generic facade and application libraries build with pure Go while the cgo build setting is enabled.
- The application command builds with cgo disabled and pure Go enabled. Its help works, two sequential operations against a native TCP fixture succeed, and the external server survives both connections. Native-qualified addresses fail explicitly.
- Native constructor E2E uses an executable path containing spaces and reserved characters; repeated URL argument parameters preserve its argv. Test outputs retain database and state artifacts.

Detailed logs, executable copies and the engine-free verification database/report are under ignored out/anki-client-refactor. Exact final candidate gates and publication are recorded by repo-delivery.

An explicitly empty subprocess argument vector remains distinct from the default serve command. Omitting arg parameters encodes zero arguments; only the executable-configured factory selects the explicit default command for an empty operation address. The native E2E also starts the real worker with no arguments and confirms that it exits with CLI help instead of silently entering serve mode.

Constructor fields use the pure ankigrpc.ClientFactory interface and ClientFactoryFunc function adapter. Deprecated Backend aliases and the application BackendBinary/BackendArguments options are removed. Existing CLI, application and archive-build fixtures configure defaults through their injected client factory.

Review follow-up uses NewSTD for direct stdio construction, executable-owned subprocess registration and direct-adapter registration for optional imports and ankigrpc/native for native hosting. The existing endpoint E2E passes for registered subprocess+cgo, TCP URLs, bare TCP addresses, Unix URLs, absolute Unix paths containing spaces and the default worker. It retains the empty default-state artifact to verify no prior endpoint or namespace leaks into later calls.

## Session observations

- ANKI_RENAME_SCOPE: two compile checks found residual field names after type renaming. Review also caught a protocol-prefix string replacement before acceptance. All were corrected; upstream Backend service names remain unchanged. Future identifier migrations should distinguish syntax from literal protocol text.
- BAZEL_REPO_CACHE_ABSOLUTE_SYMLINK: a cached Rust dependency linked to a missing temporary directory in another output base. A targeted forced fetch with the repository-contents cache disabled repaired this worktree's repository without changing shared host configuration or dependency pins.
- MANUAL_PROTO_JSON_INTEGER: the engine-free smoke check initially expected a JSON number for an int64 field. The actual protobuf JSON string was valid; the check was corrected and both reconnections passed.
- Scope and continuation are recorded in the two owner-local changes. No shared instruction or tool configuration changes were needed.

Package ownership correction consolidates the native C ABI bridge and hosting in ankigrpc/native and keeps optional build composition in command sources. The shared Go registry resolves qualified scheme availability, and subprocess argv remains exact. All 15 project suites, pure generic/application builds and the applicable executable compositions pass; final candidate gates and publication remain owned by repo-delivery receipts.

Fresh review found scoped IPv6 zone escaping and, in the in-process candidate, duplicated local-strategy classification. The existing endpoint/native API E2Es reproduce the defects before correction and pass afterward. Registry ownership metadata now shares the constructor resolver; final follow-up receipts supersede earlier candidates.
