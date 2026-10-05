## Implementation

- [x] 1. Replace the decorated lifespan with context-manager dunder methods.
- [x] 2. Consolidate runfile resolution while preserving failure behavior.
- [x] 3. Regenerate build files, review and validate existing HTTP E2E.

Verification: real HTTP E2E passed for recommendations, schema, query bounds,
missing CSV/model runfiles, corrupt model startup, shutdown and route isolation.
Gazelle produced no declaration changes. Delivery receipts record aggregate
checks against the final candidate.
