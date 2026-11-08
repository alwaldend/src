## Decisions

Proceed with a built-in injector and runtime resource provider registry. Core owns
Resource definitions (name/type/args/deps), results, graph validation, deterministic
evaluation/cache, templates, files, environment and processes. Core imports no
Vault plugin code or resource schema. The existing public res_<type> maps compile
into generic type/args entries, with deps on the resource and other settings in
args. Resource keys and dependency order remain explicit.

PluginStart returns PluginStartResponse with an optional ResourceResult in its
result field. The envelope leaves room for startup metadata independently of
resource outputs; no speculative fields are added. An absent result means
startup produced no resource output, while an explicitly empty result retains
its presence. PluginStartResponse.resource_types advertises capabilities during
startup, separately from resource fetching. ResourceFetch takes name/type/args and
dependency results and returns ResourceResult directly. Startup results retain
env/files/data, and selected startup environment exports still reach the child.
Explicit resource_provider declarations make
plugins eligible for discovery only when selected graphs need external types.
Providers must initialize capabilities at startup without fetching resources.
Reject duplicate/empty capabilities, collisions with built-ins, missing providers
and invalid graphs before starting the child. Never fetch unselected resources.

Register provider processes with the plugin lifecycle. After providers start,
build the registry and evaluate resources. Stop built-in consumer processes and
files before provider processes/tokens. On failure/cancellation, roll back all
attempted resources and wait for completion. The CLI owns the child first, then
injection shutdown, then provider shutdown. Normal plugin env conflicts retain
the existing deterministic error policy without disclosing values.

Move Vault fetchers/schema into plugins/vault; retain lazy Vault authentication,
error redaction, token revocation, TLS/runfile resolution and private temporary
SSH material. Shared Vault config owns the provider executable/runfiles.
No global host registration, new dependency, engine-specific Vault dispatch or
live infrastructure action. Injector configuration has its own field and does
not reserve a plugin name.

## Provider discovery decision

Verdict: proceed with an eligibility marker and runtime advertisements. A static
type list in configuration would duplicate the plugin's capability ownership.
Starting every configured plugin for discovery would also start unselected legacy
authentication plugins. The boolean marker limits discovery without duplicating
types, while the separate resource service leaves environment-only plugins usable.
All eligible providers start when an external type is needed; capability startup
must therefore avoid secret reads. Actual-process fixtures demonstrate this
boundary using a non-Vault provider, builtin-only calls and rejected advertisements.

Publication and review status belongs to task-local al-delivery receipts, keeping
source acceptance independent of the eventual published commit identity.

## Failure cases before implementation

Missing/cyclic dependencies; unknown/duplicate advertised types; built-in type
collision; unsupported type requested from a provider; provider startup/fetch
failure; provider process death/cancellation; leaked files/tokens/processes;
child running before all resources succeed; lost metadata or dependency results;
discovery requests sent to ineligible plugins; duplicate environment exports; builtin-only
calls accidentally starting configured providers. Extend actual-process fixtures
first and retain repeatable receipts. Reuse Vault loopback and command cleanup
coverage, including startup failure/cancellation. Test a synthetic non-Vault
provider to prove the engine is generic.

## Startup response contract

Startup has a dedicated response envelope with an optional generic result,
rather than placing environment values directly on the startup response or
conflating all startup metadata with resource output. Fetch RPCs return the
generic result directly. Register capabilities advertised in startup responses
only for eligible providers. Do not
add implicit startup-resource dependencies or change the resource authoring API.
Before implementation, extend process fixtures for absent/present startup
results, env/files/data transport and cleanup, startup capability advertisements, invalid
advertisement rollback, and skipped discovery on ineligible plugins. Preserve all
existing graph, environment-conflict and child-lifetime checks. Extend the process
fixture before implementation to distinguish nil/empty startup responses,
explicitly empty resource results and full results with env/files/data.

## Human review decisions

Proceed with all four requested corrections. Built-in env/file/process settings
use generated protobuf messages rather than anonymous JSON structs. Dedicated
Config.inject entries contain named typed resource definitions and labels;
the public res_<type> authoring maps remain unchanged. The injector no longer uses
a plugin name, plugin call payload or reserved-name check. Explicit extensions
still order deliberate graph overrides.

PluginStartResponse advertises resource_types alongside its optional result;
remove ResourceTypes RPC/messages. The provider marker still limits lazy startup
to eligible plugins when external resources are selected. Startup initializes the
provider registry without fetching unselected resources.

Explicit injected environment values take precedence over startup exports,
including keys with conflicting startup defaults. Independent startup-plugin
conflicts without an explicit injected override retain deterministic validation.
Failure cases before implementation: invalid protobuf field/value types, lost
resource settings/dependency order during schema migration, missed selectors,
startup advertisements lost when result is absent, accidental plugin-name
reservation, missing/colliding capabilities, and incorrect override precedence or
cleanup. Extend real-process fixtures and the independent authoring consumer.

The final configuration names are plugin_configs, plugins and inject. Public
constructors are al_plugin_config, al_plugin and al_inject respectively.

## Generic resource data

Proceed with ResourceConfig.data for JSON-compatible per-resource input, exposed
as .Data to environment and file templates and forwarded in ResourceFetchRequest.
Type-specific args remain separate. Remove InjectorFile.extra and .Extra;
migrate the kubeconfig CA to generic data. Dependency result data keeps its
existing .Res[name].Data and .Last.Data meaning; configured input is not merged
into provider output.

A file-only data rename would retain the special case, while a new constant-data
resource would add graph nodes for simple local variables. The shared field uses
the existing JSON transport and also serves external providers. Vault operations
consume this input as their write payload, preserving res_op.data authoring.
Failure cases before implementation: flattening data into args, losing nested or
int64 values, dropping provider input, confusing local and dependency data,
accepting removed extra settings, and losing Vault operation payloads.

The separate question about making inject a map is investigated: the current
list also encodes cross-declaration override order. A map should follow explicit
override resolution rather than introducing unordered last-writer behavior.

## Root configuration maps

Proceed with the requested complete migration. vault_conn, vault_auth,
plugin_configs, plugins and inject are maps keyed by declaration name. Remove
their redundant value names; PluginStartRequest carries the selected plugin's
name because an individual RPC no longer has its enclosing map. Selected plugin
calls also use named maps. Arguments, dependencies and result file sequences keep
their ordered meaning. No list-format compatibility layer remains.

Injection declarations retain selectors and explicit parent names derived from
Bazel extends. Resolve resource collisions by ancestry after selection and before
provider startup: descendants override ancestors; identical unrelated definitions
coalesce; different unrelated definitions fail. Inspect every maximal declaration
so coalescing equal siblings cannot hide a conflicting descendant. Extension
lookup can traverse unselected ancestors without fetching their resources.
Same-name Bazel extensions replace their atomic definition as before and retain
the inherited parent reference rather than introducing a self-reference.

Keeping lists internally would retain the API mismatch; sorting keys and applying
last-writer overrides would merely replace one hidden ordering with another.
Named maps and explicit lineage preserve selection and make conflicts observable.
Keep plugin and injection groups as namespaces, as requested: merge root entries
by group identity rather than flattening their contents. The consumer fixture
merges two groups with the same local resource key and verifies both settings
remain separate. Selected graph collisions still require explicit ancestry.
Failure cases before implementation: missing map-key identities, lost RPC plugin
names, duplicate selection, random sibling override, extension cycles/missing
parents, lost selected resources, accidental fetching of unselected ancestors,
and dropping data or dependency order while changing the schema.
