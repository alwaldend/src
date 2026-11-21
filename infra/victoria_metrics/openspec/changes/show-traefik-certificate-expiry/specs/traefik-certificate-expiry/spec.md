## Purpose

Give operators a single view of certificate lifetimes across monitored Traefik instances.

## ADDED Requirements

### Requirement: Traefik certificate expiry

Monitoring overview SHALL display days until expiration for certificates
reported by all scraped Traefik instances and identify each certificate and
reporting instance.

#### Scenario: Multiple instances report certificates

- **WHEN** Traefik instances from different scrape jobs report certificate expiration
- **THEN** the panel shows a separate series for each certificate on each reporting instance
- **AND** an already expired certificate still being reported appears with a negative value
