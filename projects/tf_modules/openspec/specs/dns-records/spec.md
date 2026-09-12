# dns-records Specification

## Purpose

Provide one canonical DNS transformation for project-local Terraform resources
and offline inspection, preserving declared names, values, and view ownership.

## Requirements

### Requirement: Normalize canonical declarations without providers

The module SHALL accept decoded owner DNS documents and a zone, flatten every
supported type member, expand and deduplicate destinations, and expose one
normalized map without provider configuration. It SHALL reject unsupported
types, destinations, malformed members, and absolute names outside the zone.

#### Scenario: One logical declaration has several types and destinations

- **WHEN** an entry contains A and AAAA with destinations all and global
- **THEN** normalization produces exactly four records, one per type and view

### Requirement: Preserve record multiplicity and stable identities

Normalized and provider resource keys SHALL use logical declaration key, type,
and view without mutable record values. Separate logical keys SHALL preserve
multiple records at the same name and type, including MX priority and TXT values.
Explicit TTLs SHALL override the type-specific compatibility defaults.

#### Scenario: A record changes address

- **WHEN** an A member changes its address without changing its logical key
- **THEN** its normalized and provider resource keys remain unchanged

### Requirement: Keep provider ownership explicit

The module SHALL accept provider instances from its caller and create individual
Cloudflare global and RouterOS dc1 records only when explicitly enabled. It SHALL
default to disabled provider ownership while preserving normalized outputs.
Cloudflare records SHALL remain unproxied.
Disabled ownership SHALL NOT imply that operational provider initialization
or credential prerequisites are suppressed.

#### Scenario: A root prepares a future migration

- **WHEN** a root calls the module without enabling ownership
- **THEN** it declares no provider record resources and still exposes every normalized declaration
