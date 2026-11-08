## ADDED Requirements

### Requirement: Built-in generic resource evaluation

AL SHALL evaluate selected named type/args/deps resource graphs with built-in
env, file and process resources, caching results and preserving dependency order.
Core SHALL NOT import provider-specific resource schemas or dispatch logic.

#### Scenario: Built-in resources without providers

- **WHEN** a selected graph contains only built-in resource types
- **THEN** AL evaluates it without launching resource providers for capability discovery

### Requirement: Advertised provider resources

Eligible plugins SHALL advertise resource types in PluginStartResponse and implement generic
resource-fetch RPCs. AL SHALL reject ambiguous/built-in capability collisions
and missing providers before running a child.

#### Scenario: Independent non-Vault provider

- **WHEN** a provider advertises a custom type and a selected graph depends on it
- **THEN** AL fetches the resource via RPC with its ordered dependency results and supplies its exports to the child

### Requirement: Shared resource results

PluginStart SHALL return PluginStartResponse with an optional ResourceResult
result field. ResourceFetch SHALL return ResourceResult directly. The startup
envelope SHALL allow startup fields independently of resource output. Startup MAY
produce no result; explicit empty results SHALL retain their presence.
Capability discovery SHALL remain independent of startup
outputs, and only eligible providers SHALL register advertised resource types.

#### Scenario: Optional startup output

- **WHEN** a plugin starts with an absent result, an explicitly empty result or environment, file and data output
- **THEN** startup succeeds, present fields reach the client unchanged and invocation-owned files remain available until consumers finish

#### Scenario: Invalid capability advertisement

- **WHEN** an eligible provider advertises missing, empty or colliding capabilities
- **THEN** AL prevents child startup and cleans up the attempted provider

### Requirement: Invocation-owned cleanup

AL SHALL keep files/processes/providers alive through child shutdown and release
consumers before providers, including failed startup and cancellation.

#### Scenario: Resource failure before child startup

- **WHEN** a graph or provider fails
- **THEN** the child does not start and all attempted invocation resources are stopped

### Requirement: Dedicated typed injector configuration

Config.plugin_configs SHALL contain executable declarations; Config.plugins SHALL
contain plugin invocations. Built-in injection SHALL use Config.inject with typed resource definitions,
separately from external plugin calls. Built-in resource settings SHALL use generated
protobuf messages. No plugin identity SHALL be reserved for injection.

#### Scenario: Ordinary plugin named injector

- **WHEN** an external plugin is named injector and built-in resources are selected
- **THEN** its startup and built-in injection operate independently

### Requirement: Explicit injection overrides startup exports

Explicitly injected environment exports SHALL override plugin startup defaults,
including conflicting startup values for the overridden key.

#### Scenario: Override conflicting startup defaults

- **WHEN** selected plugins export different defaults for a key and a selected resource exports that key
- **THEN** the child receives the injected value without relying on plugin startup order

### Requirement: Generic resource input data

Resources SHALL accept JSON-compatible data separately from type-specific args.
Environment and file templates SHALL expose this input as .Data, independently
of dependency result data. Provider fetch requests SHALL carry the configured
data without losing nested values or signed integer precision. File settings
SHALL NOT provide a separate extra field or .Extra template context.

#### Scenario: Local data and dependency results

- **WHEN** a selected resource defines data and depends on a provider result
- **THEN** its template can use .Data alongside .Last.Data and .Res without merging their values

#### Scenario: Generic provider input

- **WHEN** a selected provider resource has nested input data and large signed integers
- **THEN** the provider receives those values unchanged in ResourceFetchRequest.data

#### Scenario: Vault operation payload

- **WHEN** an operation resource defines a write body in its data field
- **THEN** the Vault provider sends that body to the configured operation endpoint

### Requirement: Named root configuration

vault_conn, vault_auth, plugin_configs, plugins and inject SHALL be maps keyed
by name. Their values SHALL NOT repeat that identity. Individual plugin startup
requests SHALL carry the selected map key. Selected calls SHALL be named maps.

#### Scenario: Root maps in an independent consumer

- **WHEN** a consumer builds and runs named configuration declarations
- **THEN** generated configuration uses object sections and runtime lookup preserves each key and its settings

#### Scenario: Group namespaces during composition

- **WHEN** different injection groups contain the same local resource key
- **THEN** merging configuration retains the separate groups and does not mix their settings

### Requirement: Explicit injection override ancestry

Selected injection resources SHALL resolve using explicit extension ancestry,
independently of declaration iteration order. Descendants SHALL override
ancestors, equal maximal definitions SHALL coalesce and unequal unrelated
maximal definitions SHALL fail before provider or child startup. Missing or
cyclic selected extension lineage SHALL fail before startup. Traversing an
unselected ancestor SHALL NOT select its resources.

#### Scenario: Descendant preceding its ancestor

- **WHEN** selected root map keys sort a descendant before its ancestor
- **THEN** the descendant resource wins through explicit ancestry

#### Scenario: Conflicting unrelated declarations

- **WHEN** selected unrelated groups define different resources under one name
- **THEN** configuration fails before running providers or the child
