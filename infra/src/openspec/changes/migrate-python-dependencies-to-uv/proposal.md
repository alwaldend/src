## Why

Use uv for all repository-owned Python dependency management, including its
native project manifests and locks. Update Python through rules_python's bundled
runtime catalog.

## What Changes

- Upgrade rules_python and its Gazelle plugin to 2.4.2 and select bundled Python
  3.14.4, without a custom runtime archive or host installation.
- Pin Bazel-managed uv 0.13.0 and use the public upstream native lock rule directly.
- Replace maintained requirements inputs and locks with component-owned
  pyproject.toml and uv.lock files; feed the native locks to pip.parse and Gazelle.
- Own firmware tooling's manifest and native lock locally while retaining its
  pinned source, existing dependency hubs and lazy source loading.
- Remove unnecessary lock wrappers, comparison tests and cross-component internal
  loads; update workflow documentation and Python skill guidance.
- Use uv for the K3s role's existing Python environment and Autoscroll's install
  hint. Deployment is excluded.

## Capabilities

No spec-level behavior changes: this is a dependency tooling and runtime update.
`skip_specs: true` applies; application APIs and component ownership stay intact.

## Impact

tools/py, tools/ansible, tools/ykman, tools/gazelle, projects/autoscroll,
projects/infinitime, projects/ansible_collection/roles/k3s_bootstrap, generated
Bzlmod metadata and Gazelle's import manifest. Acceptance requires native lock
regeneration, package-version inspection, representative Python consumers,
offline role validation and repository quality checks. No live host change is
included.
