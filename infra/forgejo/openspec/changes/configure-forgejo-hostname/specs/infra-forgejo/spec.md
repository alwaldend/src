## MODIFIED Requirements

### Requirement: Canonical service routing with legacy hostname support

Forgejo SHALL configure `https://forgejo.alwaldend.com/` as its canonical
root URL and `forgejo.alwaldend.com` as its HTTP and SSH domain. It SHALL
bind HTTP to loopback port 3000 and enable built-in SSH on port 3005.
Traefik SHALL route the canonical and inventory hostnames to that service
with the Vault certificate resolver. Local `git.alwaldend.com` HTTP requests
SHALL redirect to the canonical HTTPS URL, preserving paths and queries.
Both local hostnames SHALL retain certificate coverage. Public Git routing
and the split-view DNS records SHALL remain unchanged.

Sources: [service variables](../../../../../ansible/group_vars/all.yaml),
[Forgejo configuration](../../../../../ansible/files/forgejo.ini), and
[Traefik routing](../../../../../ansible/files/traefik_dynamic.toml).

#### Scenario: Follow a local Git link

- **WHEN** a local client requests `https://git.alwaldend.com/alwaldend/src?tab=readme`
- **THEN** it receives a permanent redirect to `https://forgejo.alwaldend.com/alwaldend/src?tab=readme`
- **AND** the canonical request reaches Forgejo without another hostname redirect

#### Scenario: Inspect support for existing Forgejo clients

- **WHEN** the deployment renders Forgejo configuration
- **THEN** its root URL and HTTP and SSH domains use `forgejo.alwaldend.com`
- **AND** the inventory hostname still routes to the loopback HTTP service
