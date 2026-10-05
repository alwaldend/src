## Why

Address review feedback on the application context manager and duplicate
runfiles resolution without changing the HTTP or startup contracts.

## What Changes

Implement callable async context-manager dunder methods on the application
owner. Consolidate runfiles resolution and file validation in one helper.

## Capabilities

Internal refactoring with unchanged public behavior.

## Impact

Application lifespan, packaged data resolver and project documentation.
