## Decision

Expose typed public symbols through package namespaces, with `__all__` declaring
the supported API. Implementation modules within the same package may import
siblings directly to avoid circular package initialization. Package exports
must not instantiate services or introduce runtime side effects.
