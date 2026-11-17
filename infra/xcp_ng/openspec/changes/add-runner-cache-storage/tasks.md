## 1. Storage integration

- [x] 1.1 Write failure cases for the command boundary before implementing it: missing input, missing or duplicate host, conflicting or duplicate SR, API failure or lost response, incomplete postcondition, credential exclusion.
- [x] 1.2 Implement the bootstrap; verify with an isolated TLS JSON-RPC fixture and retain a repeatable receipt.
- [x] 1.3 Declare the repository bootstrap and development resource-set grant; verify provider validation and package inputs.
- [x] 1.4 Match the NAS AppRole-prefixed cache share name in the bootstrap UNC path; retain creation-only reconciliation behavior and validate packaged Terraform.

## 2. Delivery

- [x] 2.1 Run affected Go, Terraform, lint and repository quality checks; record limitations of source-only acceptance.
- [x] 2.2 Record the coordinated verification and deployment prerequisites.

- [x] 2.3 Remove the custom cache-device command and all Ansible XO coupling; run retained bootstrap checks.
