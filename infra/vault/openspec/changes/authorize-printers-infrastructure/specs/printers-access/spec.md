## Purpose

Define printer credential access and the ownership and credential guarantees required by the printer infrastructure migration.

## ADDED Requirements

### Requirement: Printer credential access

Vault configuration SHALL define a printers AppRole with both required DNS credential reads and read access to the exact MikroTik-owned printers credential. MikroTik provisioning SHALL use its own administrator secret prefix, without an ingress-secret read policy. Printers SHALL NOT receive ingress credentials or general infrastructure policies.

#### Scenario: Inspect declared ownership

- **WHEN** the declared policies are inspected
- **THEN** printers can read its dedicated credential and selected DNS credentials
- **AND** printers cannot read the ingress RouterOS credential through these grants
