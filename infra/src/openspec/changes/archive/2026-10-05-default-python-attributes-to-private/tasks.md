## Implementation

- [x] 1. Implement the reviewer requests and update the owning guidance.
- [x] 2. Review cleanup failure cases and verify existing HTTP E2E and consumers.

Verification: HTTP E2E passed and service-lifecycle.json records successful
context exit, repeated close and deletion cleanup using the real model.
The Python skill offline configuration check passed. Public response fields
and exported app/router properties retain their contracts. Aggregate delivery
gates are bound to the prepared candidate by the trusted receipt.
