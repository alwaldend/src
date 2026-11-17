## ADDED Requirements

### Requirement: Runner external proxy

Secure1 jobs SHALL receive HTTP and HTTPS proxy settings using the existing
Vault-owned proxy credential. The remote proxy connection MUST use TLS with
certificate and hostname verification. Credentials MUST NOT appear in source,
Ansible output or process arguments.

#### Scenario: Download an external dependency

- **WHEN** a runner job fetches an external HTTP or HTTPS dependency
- **THEN** its proxy environment selects the existing VPN-backed HTTP proxy
- **AND** a loopback TLS forwarder supports clients using plain HTTP proxies

### Requirement: Local proxy bypass

Runner proxy settings SHALL bypass localhost, the apex `alwaldend.com` and
every subdomain of `alwaldend.com`, in uppercase and lowercase environment
forms. Proxy settings SHALL remain scoped to runner jobs.

#### Scenario: Contact the local forge

- **WHEN** a runner job requests `https://forgejo.alwaldend.com`
- **THEN** the destination is contacted directly
- **AND** `https://releases.hashicorp.com` remains proxied
