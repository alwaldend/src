## ADDED Requirements

### Requirement: Forgejo HTTPS authentication uses the native login

The dev1 Git configuration SHALL use the managed user's fj login for the
Forgejo alwaldend/src HTTPS remote.

#### Scenario: Git requests credentials for the configured repository

- **WHEN** Git requests credentials for HTTPS on forgejo.alwaldend.com and alwaldend/src
- **THEN** the helper returns the login credential for that Forgejo host

#### Scenario: Git requests credentials outside the configured scope

- **WHEN** Git requests credentials for another host, protocol, or repository
- **THEN** the helper returns no credential
