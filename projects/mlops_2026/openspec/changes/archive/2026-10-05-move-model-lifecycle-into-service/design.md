## Decision

The service owns model loading, context-manager cleanup and an idempotent deletion fallback. Inject it into the application and router constructors. Keep implementation state private with explicit public app/router properties.

Preserve per-app model state and HTTP behavior. Cleanup must be safe before
loading, after context exit and on repeated calls; deletion is a fallback,
not a replacement for deterministic shutdown. Public model fields retain
serialization names, and existing public APIs are not renamed wholesale.
