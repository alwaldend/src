# Infrastructure MikroTik

## Purpose

Describe the checked-in RouterOS export snapshots and site DNS declarations for
the dc1 routers. This package is a documentation and record-input owner: its
BUILD file declares no router deployment executable. Export contents are
historical source evidence, not a verification of current device configuration.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [export workflow](../../../README.md),
[package targets](../../../BUILD.bazel),
[router1 snapshot](../../../router1.rsc),
[router2 snapshot](../../../router2.rsc), and
[DNS declarations](../../../dnsconfig.json).

## Requirements

### Requirement: Preserve export snapshots as the package's documented interface

The package SHALL retain the router export files as documentation inputs and
describe the RouterOS `/export` collection workflow. Documentation SHALL
distinguish these snapshots from an automatically applied desired-state
configuration.

#### Scenario: Inspect the router package

- **WHEN** a reader opens the component documentation
- **THEN** the reader can inspect both router export snapshots
- **AND** the available Bazel targets expose documentation and DNS inputs rather
  than an automated router deployment

### Requirement: Describe the two routers according to their recorded topology

The baseline description SHALL identify router1's separate wired and wireless
bridges with their own DHCP networks, and router2's single bridge containing its
eight Ethernet interfaces with a DHCP client. It SHALL NOT infer additional
live devices or enabled topology from comments elsewhere in the repository.

#### Scenario: Compare the recorded network roles

- **WHEN** the two checked-in exports are compared
- **THEN** router1 records `bridge1` for wired ports and `bridge2` for Wi-Fi ports
- **AND** router2 records `bridge01` and an IP DHCP client on that bridge

### Requirement: Retain the recorded router1 traffic-policy structure

The router1 snapshot SHALL describe named interface lists for input and
forwarding permissions, explicit terminal IPv4 input and forwarding drops,
and source NAT masquerading for WAN egress. This contract describes the export
and SHALL NOT be represented as a live firewall audit.

#### Scenario: Review router1's recorded IPv4 policy

- **WHEN** a reader follows the IPv4 firewall section of the snapshot
- **THEN** the listed service permissions precede terminal input and forwarding
  drop rules
- **AND** the NAT section declares masquerading for traffic leaving the WAN list

### Requirement: Supply site-local DNS records to the DNS owner

The package SHALL expose its `dnsconfig.json` as a Bazel filegroup. Its router,
switch, and bare-metal records SHALL be assigned to the `dc1` destination view.

#### Scenario: Assemble the dc1 DNS configuration

- **WHEN** the DNS owner consumes `//infra/mikrotik:dnsconfig`
- **THEN** router1 IPv4 and IPv6 records and the declared switch and bare-metal
  IPv4 records are available for the dc1 view
