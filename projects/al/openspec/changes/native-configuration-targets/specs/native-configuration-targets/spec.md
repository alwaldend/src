## ADDED Requirements

### Requirement: Native configuration targets

AL SHALL provide composable Bazel targets for Vault authentication, Vault connections, plugins and plugin calls, preserving explicit override relationships, labels and payloads. Plugin executables and runtime inputs MUST use declared label dependencies.

#### Scenario: Component configuration composition

- **WHEN** a component builds al_config with shared and component fragment dependencies
- **THEN** the output is valid AL JSON containing definitions selected by explicit extension ancestry and the runtime files required by its plugins

#### Scenario: Reusable plugin configuration

- **WHEN** a command consumes an inherited configuration target registering an external plugin
- **THEN** the plugin executable and its runfiles are available without a direct plugin dependency in the command's data

#### Scenario: Explicit public file input

- **WHEN** a plugin call substitutes a declared public file into a payload
- **THEN** its generated JSON contains the escaped file contents and Bazel tracks that file as an action input

#### Scenario: TLS runfile resolution

- **WHEN** a generated Vault connection selects CA and client certificate/key labels and a command runs outside its runfiles directory
- **THEN** runtime TLS initialization resolves those identifiers through the available directory or manifest runfiles context before opening the files

### Requirement: Declarative configuration loading

AL SHALL load JSON and YAML configuration without evaluating Lua.

#### Scenario: Removed Lua format

- **WHEN** a Lua file is passed as configuration
- **THEN** AL rejects the unsupported extension without evaluating its contents

### Requirement: Named configuration dictionaries

AL SHALL accept a whole configuration dictionary with vault_auth, vault_conn, plugins and plugin_calls sections keyed by name. Plugin payload dictionaries SHALL use args and executable argument vectors SHALL use argv.

#### Scenario: Dictionary compilation

- **WHEN** a consumer builds an inherited configuration with local named declarations
- **THEN** the generated configuration preserves names, labels, arguments, executable/TLS label dependencies and explicit local override relationships

### Requirement: Explicit configuration extensions

Atomic configuration targets SHALL accept an extends dependency of the same
kind. Aggregates and command wrappers SHALL resolve definitions by identity
and extension ancestry independently of input list order, include shared
ancestors once, and reject unrelated definitions of the same kind and name.
A descendant with the same runtime identity SHALL replace the complete parent
definition. Differently named call descendants SHALL follow their ancestors.

#### Scenario: Auth override listed before its parent

- **WHEN** component auth extends a shared default auth target and an aggregate or command lists the component first
- **THEN** the generated file contains only the component's default auth definition

#### Scenario: Unrelated competing definitions

- **WHEN** two unrelated targets define the same auth name
- **THEN** analysis fails with a conflict diagnostic requiring explicit extensions

### Requirement: Dictionary plugin selectors

Command wrappers SHALL accept plugin_labels dictionaries and encode exact-value
or presence-only selectors into AL arguments.

#### Scenario: Empty selector value

- **WHEN** a wrapper selects a label with an empty dictionary value
- **THEN** the real plugin receives the call through presence-only selection

### Requirement: Target-derived plugin call defaults

Generic al_plugin_call targets SHALL default their call identity and selected
plugin name to the Bazel target name. config_name and plugin SHALL provide
independent explicit overrides. Dictionary compiler entries SHALL default plugin
to the entry key, keeping generated private fragment names out of runtime identities.

#### Scenario: Different runtime call identity

- **WHEN** a generic call overrides config_name but omits plugin
- **THEN** its runtime call identity uses the override and its selected plugin still uses the target name
