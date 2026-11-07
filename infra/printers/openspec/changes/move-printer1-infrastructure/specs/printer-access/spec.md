## Purpose

Define printer network configuration and the ownership and credential guarantees required by the printer infrastructure migration.

## ADDED Requirements

### Requirement: Printer network configuration

The printers project SHALL own the canonical printer1 DNS declaration, both DNS views, and its guarded IPv4 Wi-Fi forwarding exception. It SHALL authenticate through its own AppRole and dedicated RouterOS credential. Source relocation SHALL include a non-destructive state migration procedure.

Adoption SHALL retain the existing rule identity without replacement caused by
missing creation-only insertion metadata. New rules SHALL use the guarded
insertion point; live ordering SHALL be checked during adoption.

#### Scenario: Inspect declared ownership

- **WHEN** the printer roots are inspected
- **THEN** printer1 retains its name, address, DNS views, forwarding match and exact-one-drop guard
- **AND** migration instructions preserve live provider objects and prohibit simultaneous state ownership
