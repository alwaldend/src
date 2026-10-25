## ADDED Requirements

### Requirement: Keep the public root site on GitHub Pages

The public apex SHALL declare the four GitHub Pages A records and four AAAA
records with a 300-second TTL while the public download deployment is deferred.
It SHALL NOT declare an apex CNAME to the download service. Local apex routing,
mail, and unrelated names SHALL retain their existing declarations.

#### Scenario: Reconcile the deferred public deployment

- **WHEN** the DNS owner applies its maintained declaration
- **THEN** the public root site resolves through the eight GitHub Pages addresses
- **AND** the local apex continues using the download VM's IPv4 endpoint

#### Scenario: Restore the public apex after an early cutover

- **WHEN** an authorized rollback replaces the download apex CNAME
- **THEN** its removal completes before the address records are created
- **AND** a subsequent plan contains no changes and unrelated provider records remain unchanged
