## Why

PR feedback clarifies that namespaced imports allow importing modules from
their parent package, and require that form for long module paths.

## What Changes

Correct the Python skill, syntax enforcement and fixtures, and affected imports.
Preserve the prohibition on importing individual symbols or using wildcards.

## Capabilities

Correction to existing import policy enforcement; no new runtime capability.

## Impact

Python quality checks, the Python skill, and existing Python consumers.
