## Implementation

- [x] 1. Implement reviewer requests and regenerate affected BUILD files.
- [x] 2. Review the PR and validate the exact candidate with existing E2E checks.

Verification: real HTTP E2E passed with unchanged query validation, OpenAPI,
startup failures, shutdown and factory route isolation. Tracked-file import
checks and offline skill configuration passed. Aggregate publication gates
are recorded by the delivery receipt against the final candidate.
