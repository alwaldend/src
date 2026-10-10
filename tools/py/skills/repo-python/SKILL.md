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

## Document modules and declarations

Every first-party Python module, including package `__init__.py` files, must
have a module docstring describing its purpose, responsibilities, and important
boundaries or guarantees. Document every class, function, and method, including
private declarations, constructors, async functions, and test helpers, with a
docstring immediately inside the definition. Describe purpose and behavior;
record meaningful contracts without merely restating the name.

Document every module-level constant, variable, and type alias, and every class
or instance attribute, including private and serialized fields, with a variable
or attribute docstring immediately after its declaration or assignment. Use a
standalone string literal describing its purpose and contract, including in
nested assignment contexts and methods other than `__init__`. Do not substitute
comments for docstrings or move assignments merely to satisfy extraction tools.
Source documentation tools may not extract these strings in every context;
that limitation does not relax the documentation requirement. These docstrings
do not become the assigned value's runtime `__doc__`.
Apply these rules to declarations and fields in tests too. Preserve accurate
existing documentation and update it with the declaration. Do not edit generated
code or upstream dependencies; document generated contracts at their source.
Project architecture belongs in the owning README.

## Require type hints everywhere

Type hints are mandatory everywhere in Python code. Annotate every function
and method parameter and return value, including constructors, async functions,
callbacks, and test helpers. Use `-> None` for functions without a return value;
`self` and `cls` use their implicit receiver types. Annotate module and class
attributes and local variables as well, even when inference is possible.
Use precise collection and optional types; do not silence missing annotations
with `Any`, blanket ignores, or disabled type checks. Model external boundaries
with appropriate typed interfaces and narrow validated values before use.

## Keep implementation attributes private

Instance and class attributes are private by default: prefix implementation
state with `_`, for example `self._model` and `self._recommendations`.
Expose public attributes only as deliberate APIs or required framework contracts.
Prefer typed read-only properties for public access to owned implementation
objects. Keep Pydantic response fields and other serialization contracts public;
do not change their names or established public APIs merely to add `_`.

## Require namespaced imports

In implementation and consumer files, import modules and reference their members
through the module namespace.
Use `import pydantic` with `pydantic.BaseModel` and `pydantic.Field`, and
`import fastapi` with `fastapi.Query`. Apply this to standard-library and
repository modules too: use `import pathlib` and `pathlib.Path`, or
`from projects.example.internal import models` and `models.Response`.
For long module paths, use `from package.path import module`, with a clear
module alias when needed. This imports a module, not one of its members.
Combine imports from the same package in one `from package import a, b`
statement, letting the formatter wrap long lines. This includes aliased imports.
Use a common exported package namespace when combining imports with different
parent paths; do not duplicate exports solely to shorten an import. Ruff import
lint (`I`) and isort enforce combining imports; Ruff formatting alone does not
sort them. The repository sorter settings enable combining aliased imports.
Strongly prefer package APIs with explicit exports in `__init__.py` over
cross-package imports of individual implementation files. Import the package
and use its exported members through that namespace, for example
`from projects.example.internal import models` and `models.Response`.
Keep exports typed and declare `__all__`; do not instantiate services in exports.
Implementation modules within one package may import siblings directly to
avoid circular initialization. Keep package initialization free of side effects.
Do not import individual classes, functions, or constants with `from ... import`
outside `__init__.py`. In `__init__.py`, direct symbol imports are encouraged
for public re-exports listed in an explicit `__all__`, for example
`from .models import Response` with `__all__: tuple[str, ...] = ("Response",)`.
These imports preserve the original type annotations and signatures without
redundant alias assignments. Wildcard imports remain forbidden.
Prefer the imported module's own name. Do not append redundant suffixes such as
`_module`, and do not write `module as module`. Use a descriptive alias only
to resolve a name collision or clarify an otherwise ambiguous module name.

The `//tools/repo_quality:python_import_quality_test` target rejects symbol imports
outside explicit `__init__.py` re-exports in tracked first-party Python and is included in `//:repo_quality_test`.

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
