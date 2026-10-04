## Why

The reviewer requests idiomatic direct symbol imports in package export files.

## What Changes

Allow and encourage direct imports in `__init__.py` for explicit `__all__`
exports. Simplify affected package exports and update repo-python and lint
fixtures while preserving namespaced consumer imports.

## Capabilities

Modified `python-import-validation`: allow explicit package symbol re-exports.

## Impact

Python skill, import quality checker, and recommender and Autoscroll exports.
