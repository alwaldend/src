## Implementation

- [x] Document SSH target discovery and missing-target handling.
- [x] Add and validate an evaluation case.
- [x] Complete repository quality gates.

## Session evidence

Two direct SSH attempts preceded discovering the repository SSH abstraction:
one failed on controller configuration permissions, the other on key authentication.
The new routing guidance and evaluation address this discovery failure. A fresh
agent can use the XCP-ng target, which succeeded after Vault login refresh and
an sshpass compatibility correction. Offline fixture validation does not
establish live host access or model behavior.
