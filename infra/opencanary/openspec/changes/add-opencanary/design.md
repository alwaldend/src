## Context

The Xen Orchestra deployment pattern uses a Fedora cloud template and per-component resource sets. Thinkst's upstream Ansible role uses Debian-specific Samba ownership and its DNF dependency tasks target releases 8 and 9.

## Goals / Non-Goals

Provision one guest and deploy a broad protocol profile with locally retained events. External alert destinations and dashboards are outside this change.

## Decisions

Use a first-party OpenCanary role with a versioned Python package and a supervised non-root service. Keep host SSH on TCP 22 and the decoy on TCP 2222. Use a dedicated operator-selected DHCP network and require explicit client and decoy source ranges. Port-scan detection requires legacy iptables and is excluded from this Fedora profile.

The OpenCanary AppRole has a dedicated Terraform root and state within the Vault project.

## Risks / Trade-offs

Syslog or a dashboard can consume the resulting events if configured separately. Decoy data can include credentials supplied by clients, so access to the journal remains restricted.
