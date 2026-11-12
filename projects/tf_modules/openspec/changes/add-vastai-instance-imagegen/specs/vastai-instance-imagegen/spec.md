## Purpose

Provide reusable template-based image-generation rentals with SSH forwarding
and explicit separation between configuration validation and paid deployment.

## ADDED Requirements

### Requirement: Select a template and a bounded datacenter offer

The module SHALL accept an AI-Dock Forge template hash and select the cheapest
matching verified datacenter offer with one GPU, at least 24 GiB VRAM, 100 GB
storage, 99% reported reliability, direct SSH connectivity, and at most USD 1
per hour by default. GPU memory, disk allocation, hourly ceiling, and GPU model
SHALL be configurable. An optional explicit offer ID SHALL retain these filters.
The module SHALL create one instance only when its caller applies it, and
changing marketplace results SHALL NOT replace an existing rental.

#### Scenario: Caller composes the module

- **WHEN** a caller supplies a valid template and the marketplace has a matching datacenter offer
- **THEN** the module declares one SSH instance with the template's image and provisioning environment

### Requirement: Restrict access to SSH

The module SHALL reject template Docker options other than environment settings,
disable AI-Dock public tunnel services, and bind Forge to loopback.

#### Scenario: Template publishes a web port

- **WHEN** the selected template contains a port publication option
- **THEN** planning fails before rental creation

#### Scenario: SSH metadata is available

- **WHEN** the instance has been refreshed with SSH host and port
- **THEN** the output command forwards 127.0.0.1 on the local machine to the remote loopback Forge port

#### Scenario: SSH metadata is pending

- **WHEN** the API has not returned SSH host and port
- **THEN** the forwarding command is null

### Requirement: Validate without renting

The module SHALL provide offline validation without live credentials or paid
capacity, and SHALL document the remaining live acceptance checks.

#### Scenario: Offline acceptance

- **WHEN** the packaged acceptance check runs
- **THEN** it exercises the real module through Terraform against a loopback API fixture and retains a repeatable result artifact

#### Scenario: No eligible offer

- **WHEN** no offer meets the datacenter, capacity, connectivity, and hourly ceiling constraints
- **THEN** planning fails without creating a rental

### Requirement: Shared bounded transfer prices

The shared offer declaration SHALL cap incoming and outgoing transfer prices
at USD 0.04/GB each. The regular rental and ephemeral evaluator SHALL consume
that declaration while preserving their hourly-price, resource and mandatory
verified/datacenter constraints. This rate ceiling SHALL NOT be presented as
a total spending cap.

#### Scenario: Transfer price exceeds the ceiling

- **WHEN** either transfer direction costs more than USD 0.04/GB
- **THEN** the shared search predicates SHALL exclude the offer
