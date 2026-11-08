## Context

See proposal.md. Thirty-three current configuration targets compose 33 Lua sources; another DNS Lua source is unreferenced. Shared helpers encode repeated injection flows. Two configurations read the public SSH server CA during build.

## Goals / Non-Goals

Replace every AL Lua source with explicit targets while retaining runtime selection and credential boundaries. No deployment or new external software.

## Decisions

Proceed with typed Vault, plugin and plugin-call targets supplied by AL. Compose targets with al_config. Repository-only infrastructure macros retain common resource construction in infra. Use executable labels and declared public file substitutions. Reject checked-in generated JSON because it would keep authoring separate from Bazel dependencies; reject expanded copies of helper resources because they create multiple owners.

Centralize repository Terraform authentication composition in tools/al rather
than adding AL/Vault dependencies to standalone rules_terraform. Repository
al_terraform_* maps accept configs and plugin_labels directly. Keep plain offline
Terraform checks on the standalone rules. The Vault backend plugin owns its
target macro; shared configuration registers its executable once.
Paths and stage labels stay with callers. Public infrastructure constructors use the infra_ namespace.
Proceed with this integration layer: it removes caller boilerplate while
retaining the standalone module boundary and exact operation contracts.

## Risks / Trade-offs

- Configuration ordering or resource drift: compare every built configuration with the captured baseline after canonical normalization of JSON fields and executable paths.
- Implicit runtime inputs: make certificates and executable runfiles dependencies explicit.
- Lua removal breaks direct consumers: accepted by the user; remove its loader and current instructions.

## Migration Plan

Capture baseline, implement the linked AL API, translate declarations, remove Lua, validate outputs and consumers, deliver the source PR. Source rollback is a Git revert; no live migration is performed.

Failure cases for the Terraform integration: lost config/runtime inputs,
incorrect AL argument boundaries, dropped operation arguments or stage labels,
altered targeted saved-plan guards, a test wrapper launching AL twice, and a
helper changing Vault paths, mounts or named authentication. The end-to-end
fixture selects a synthetic environment resource and executes AL, the injector
and the Terraform launcher against an isolated executable probe, which verifies
the received environment, arguments and packaged configuration and retains a
JSON receipt. A separate test executes the pinned Terraform version command.
No backend is initialized.

## Composable targets and configuration locality

Proceed with small native AL fragments and order-independent aggregate configs. Shared
infrastructure and backend helpers declare reusable configuration targets;
short single-use declarations stay in BUILD.bazel, including small dictionaries
and direct calls to common macros. Extract namespaced local defs.bzl macros
only for significant reuse or a substantial readability improvement for complex
graphs, construction logic, long templates or bulky settings. Fixed-attribute
wrappers alone do not justify extraction. Moving an entire configuration dictionary into another file would
retain its coupling; split declarations by authentication, plugin or call.
The dictionary compiler remains available for small declarative inputs.

Command wrappers accept plugin_labels dictionaries and expand them into AL's
CLI selection arguments. Keys are selector labels and empty string values mean
presence-only matching. Preserve selector semantics, call order, credential
references, native label dependencies and aggregate target names. Generic
run_args remain only for other AL command options.

Per-target extends relationships replace input-order overrides. Atomic targets
publish named definitions and ancestry; aggregation deduplicates shared nodes,
selects descendant definitions and rejects unrelated competitors. Commands
compose multiple config targets through the same provider-aware merge before
passing one file to AL. Aggregate reference order can follow the formatter.
Retain explicit extension edges only for real override relationships, including
shared authentication and intentional call/resource overrides.

Proportionality correction: inline the trivial local wrappers introduced by the
previous extraction pass. Retain the six local macros that isolate larger
resource graphs, long templates or certificate/OIDC settings. Compare all
39 captured production and fixture outputs exactly after parsing JSON; no
runtime target identity, dependency edge, argument or credential reference may
change. The Forgejo XO declaration reuses infra_xo_login with the same endpoints.

Injector authoring follows the AL-owned flat named-group API. Replace generic
plugin/args/res scaffolding and the narrow environment helper with
al_injector_call. Keep structured entries, explicit resource names/dependencies,
per-resource Vault metadata and existing call names. Default plugin-call identity
comes from the target; retain config_name only where it differs.
