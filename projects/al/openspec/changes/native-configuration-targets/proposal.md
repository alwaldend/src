## Why

Repository configuration was expressed through Lua evaluated during Bazel builds. The user requests Bazel-native declarations, convenient dictionary authoring throughout, and AL ownership of its injector plugin.

## What Changes

- Replace Lua source and evaluation with an al_config dictionary compiler over native Bazel fragments.
- **BREAKING**: use dictionaries keyed by name for configuration sections and injector resources; remove duplicated name fields.
- **BREAKING**: rename plugin payload parameters to args; use argv for executable command-line arguments. Keep args.res and structured env objects.
- Move the external injector and schema to projects/al/plugins/injector and make its builds/tests work in the standalone AL module.
- Preserve executable/TLS label dependencies and runtime runfile resolution.

## Capabilities

### New Capabilities

- `native-configuration-targets`: named configuration dictionaries, target composition, runtime input resolution and declarative loading.
- `injector-dictionary-configuration`: named resource graphs and explicit ordered dependencies.

### Modified Capabilities

None.

## Impact

AL rules, loader, plugin sources, existing-pinned module dependency declarations, tests and documentation. Repository integration is tracked by infra/src/openspec/changes/migrate-al-configuration-to-bazel. No live infrastructure/authentication operations or new external software. The proposed built-in injection engine and resource-provider registration protocol are explicitly deferred.
