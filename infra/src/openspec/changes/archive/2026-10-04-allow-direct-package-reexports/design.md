## Decision

In `__init__.py`, prefer direct symbol re-exports listed in explicit `__all__`.
The checker accepts those bindings only in export files, including relative
imports, while rejecting wildcard imports and non-exported direct symbols.
Other source files continue to import module or package namespaces.
