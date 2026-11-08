## Why

Repository AL configuration is already packaged by Bazel but still evaluated from Lua. The user requests native configuration targets and removal of all AL Lua files, without backward compatibility.

## What Changes

- Replace Lua declarations and shared helpers with composable Bazel rules and repository Starlark helpers.
- **BREAKING**: remove Lua configuration sources and Lua configuration loading.
- Preserve component AppRoles, plugin selection labels, injection resources, backend paths, and existing command targets.
- Make executable and public certificate inputs explicit Bazel dependencies.

## Capabilities

### New Capabilities

None. This is a configuration representation migration; runtime command behavior is preserved. This change uses `skip_specs: true`.

### Modified Capabilities

None. The reusable rule API is tracked by the linked AL owner change.

## Impact

Repository configuration BUILD files, infrastructure helpers, AL rules and loader, current documentation and procedural references. No live authentication, provisioning, deployment, host configuration changes, or new dependencies. Related AL change: `projects/al/openspec/changes/native-configuration-targets`.
