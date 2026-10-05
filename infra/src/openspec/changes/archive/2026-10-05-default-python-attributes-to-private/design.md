## Decision

The Python skill requires implementation attributes to start with \_, exposing public members only as deliberate APIs or framework/serialization contracts.

Preserve per-app model state and HTTP behavior. Cleanup must be safe before
loading, after context exit and on repeated calls; deletion is a fallback,
not a replacement for deterministic shutdown. Public model fields retain
serialization names, and existing public APIs are not renamed wholesale.
