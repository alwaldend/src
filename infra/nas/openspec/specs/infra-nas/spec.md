# NAS infrastructure Specification

## Purpose

Describe the TrueNAS-related configuration surface currently owned by
`infra/nas`. The baseline is checked-in source at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08.
The directory contains an AL configuration, DNS declarations, Vault wrappers,
and brief documentation. It does not contain a TrueNAS provisioning playbook
or Terraform configuration, and this baseline makes no claim about live NAS
configuration or availability.

Source baseline limitation: the owning BUILD file restricts the AL target's
visibility to `//infra/harvester:__subpackages__`, although `infra/harvester`
does not exist in this revision. This declaration does not establish a working
consumer integration.

## Requirements

### Requirement: Scoped NAS configuration packaging

The project SHALL declare its AL configuration through `//infra/nas:al`.
That configuration SHALL depend on the shared `//infra:al` configuration.
The project's Vault command map SHALL use its own AL target.

Sources: [project description](../../../README.md) and
[configuration and command targets](../../../BUILD.bazel).

#### Scenario: Inspect the NAS configuration and Vault packaging

- **WHEN** a maintainer inspects the NAS target declarations
- **THEN** the AL target includes `al.lua` and the shared infrastructure AL dependency
- **AND** the Vault command map references that NAS AL target

### Requirement: Distinct local and global cloud DNS destinations

The NAS DNS declaration SHALL expose the NAS host and root names as A records,
and SHALL distinguish the DC1 and global `cloud` aliases. The DC1 `cloud` alias
SHALL target `nas.alwaldend.com.`, while the global alias SHALL target
`ingress.alwaldend.com.`. The DNS filegroup SHALL remain visible only within
the infrastructure tree.

Sources: [DNS declarations](../../../dnsconfig.json) and
[DNS target visibility](../../../BUILD.bazel).

#### Scenario: Consume DNS declarations for different views

- **WHEN** a DNS consumer selects the entries tagged `dc1` and `global`
- **THEN** it receives the NAS destination for the DC1 `cloud` alias and the
  ingress destination for the global `cloud` alias.
