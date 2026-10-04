---
name: repo-python
description: >-
  Implement, refactor, and review Python code in this Bazel monorepo with
  mandatory type hints, repository formatting, dependency locks, and Gazelle.
  Use alongside repo-bazel and bazel-agent for builds and validation.
---

# Work with Python in this repository

Read the owning README and BUILD.bazel before editing. Use project-layout for
new source locations and preserve unrelated legacy layouts.

## Require type hints everywhere

Type hints are mandatory everywhere in Python code. Annotate every function
and method parameter and return value, including constructors, async functions,
callbacks, and test helpers. Use `-> None` for functions without a return value;
`self` and `cls` use their implicit receiver types. Annotate module and class
attributes and local variables as well, even when inference is possible.
Use precise collection and optional types; do not silence missing annotations
with `Any`, blanket ignores, or disabled type checks. Model external boundaries
with appropriate typed interfaces and narrow validated values before use.

## Respect Python support and formatting

The root pyproject.toml owns shared formatter and type-checker settings. Read
it rather than duplicating settings. The project support range is declared in
tools/py/pyproject.toml. That range, root formatter syntax targets, and the mypy
target are distinct from the Bazel runtime selected in
tools/py/include.MODULE.bazel; do not assume the runtime permits newer syntax
in every component. Check narrower component settings when present.

Use existing repository Black, isort, mypy, Flake8, or Ruff targets as applicable
and the configured semantic lint profile. Discover available labels before
running checks; do not replace the pinned toolchain with host Python or install
packages into the host environment for validation.

## Dependencies and generated BUILD files

Use repo-external-dependency for dependency approval and pinning. Root Python
dependencies are declared in tools/py/pyproject.toml and locked in
tools/py/requirements.txt.
Component-specific manifests retain their own ownership; inspect their wiring
in tools/py/include.MODULE.bazel before changing them.

For root dependency changes, use the owning workflows in order:

```sh
bazel_agent bazel run //tools/py:requirements.update
bazel_agent bazel run //tools/gazelle:gazelle_python_manifest.update
bazel_agent bazel run //:gazelle
```

Do not hand-edit requirements.txt or gazelle_python.yaml. Review generated
diffs, including transitive packages and import mappings. Preserve requested
extras. Gazelle wires Python imports to Bazel dependencies; declare runtime
resources in data and inspect generated entry points and library boundaries.

## Verify the affected behavior

Use bazel-agent and repo-bazel for pinned execution, package checks, and
semantic lint. Validate root dependency changes with
//tools/py:requirements.test and //tools/gazelle:gazelle_python_manifest.test,
plus affected consumers. Prefer repeatable end-to-end artifacts for complex
behavior and follow the root test policy before implementing isolated tests.
Complete repo-delivery's publication gates. Offline skill eval validation proves
configuration and packaging, not that an agent follows these instructions.
