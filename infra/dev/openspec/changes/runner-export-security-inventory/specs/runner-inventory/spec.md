## ADDED Requirements

### Requirement: Authoritative secure runner inventory

The owner inventory SHALL include the verified
`secure.forgejo-runner.alwaldend.com` host in the `forgejo_runner` group and
remain the source consumed by the shared security workflow.

#### Scenario: Shared security workflow loads runner hosts

- **WHEN** the workflow consumes `//infra/forgejo_runner/ansible:inventory`
- **THEN** it receives the existing `host1.forgejo-runner.alwaldend.com` entry
  and the confirmed `secure.forgejo-runner.alwaldend.com` entry from the same
  owner source without a separately maintained host list
