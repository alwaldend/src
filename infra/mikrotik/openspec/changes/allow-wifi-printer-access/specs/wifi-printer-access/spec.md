## Purpose

Provide a stable globally resolvable printer name and narrowly scoped access from the
wireless network to the wired printer without opening other wired destinations.

## ADDED Requirements

### Requirement: Global and local printer DNS

The canonical DNS declaration SHALL resolve `printer1.dc1.alwaldend.com` to
`192.168.1.221` in both the dc1 and global views.

#### Scenario: Resolve the printer in either view

- **WHEN** the adopted DNS root reconciles the printer declaration
- **THEN** it manages matching dc1 and global A records for the printer
- **AND** the global record remains unproxied and points to the private printer address

### Requirement: Host-specific wireless forwarding

The firewall configuration SHALL permit IPv4 traffic from `192.168.2.0/24`
arriving on `bridge2` to the printer through `bridge1`. The exception SHALL
match only the printer's canonical IPv4 address and precede the enabled
terminal forwarding drop. Missing or ambiguous drop-rule identification SHALL
prevent deployment of the exception.

#### Scenario: Connect from Wi-Fi

- **WHEN** a Wi-Fi client connects to the printer after authorized deployment
- **THEN** the exception permits its routed traffic to the printer
- **AND** existing connection tracking permits replies

#### Scenario: Preserve isolation

- **WHEN** traffic targets another wired address or arrives outside bridge2
- **THEN** it does not match the new exception

#### Scenario: Reject an uncertain insertion point

- **WHEN** the enabled terminal forwarding drop is absent or not unique
- **THEN** planning fails rather than deploying an ineffective or misplaced rule

### Requirement: Independent declarative operation

The firewall root SHALL use separate state and its own AppRole through the
repository injection workflow. It SHALL reuse the existing non-DNS ingress
RouterOS credential through an exact-secret read grant, without copying the
credential or assuming the ingress identity. It SHALL leave the DNS-only root
and router export snapshots intact. Offline checks SHALL not require live credentials.

#### Scenario: Validate without deployment

- **WHEN** offline checks evaluate the firewall configuration
- **THEN** they validate the configuration against the pinned provider schema without contacting a router
- **AND** they do not claim live Wi-Fi reachability or automatic printer discovery
