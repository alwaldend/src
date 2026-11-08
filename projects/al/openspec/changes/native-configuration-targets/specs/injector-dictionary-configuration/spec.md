## ADDED Requirements

### Requirement: Named injector resources

Injector arguments SHALL expose res as a dictionary keyed by resource names. Resource configuration values, including env, SHALL remain structured objects.

#### Scenario: Named environment injection

- **WHEN** AL selects an injector call containing named KV and environment resources
- **THEN** the injector evaluates the selected resources and supplies the expected environment to the child without fetching unselected calls

### Requirement: Ordered dependency references

The injector SHALL evaluate dependencies in the declared dependency list order, independently of resource dictionary insertion order.

#### Scenario: Last dependency selection

- **WHEN** a resource lists dependencies in reverse lexical order and uses Last in its template
- **THEN** Last selects the final declared dependency

### Requirement: Flat named resource-type maps

AL SHALL expose al_injector_call with named `res_<type>` maps such as res_env, res_file and res_kv. Structured entries MAY contain deps, vault_auth and
vault_conn metadata. The macro SHALL translate them into the existing args.res
map, reject duplicate names across groups and reject non-dictionary settings.
The call name SHALL default to its Bazel target name, with config_name as an
explicit override; the selected plugin SHALL be injector.

#### Scenario: Resource metadata and typed names

- **WHEN** a consumer builds flat groups containing ordered dependencies, Vault metadata and typed process/OIDC names
- **THEN** metadata remains on each generated resource, dependency order is preserved and typed name fields remain inside their type settings

#### Scenario: Duplicate or scalar entries

- **WHEN** two type maps use the same resource key or an entry is a scalar
- **THEN** package loading rejects the invalid declaration with a specific diagnostic

#### Scenario: Resource namespace and call attributes

- **WHEN** a caller supplies res_<type> maps alongside labels, extends or substitutions
- **THEN** the macro derives each resource type from its attribute suffix and forwards the other call attributes without consuming them
