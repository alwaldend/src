# Ceph infrastructure specification

## Purpose

Describe the single-host Ceph CRUSH map maintained by `infra/ceph` and its
documented operator workflow. This owner contains a map and documentation;
its BUILD file exposes documentation, not a Ceph deployment target. The
baseline describes checked-in source, not an observed running cluster.

Baseline revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observed: 2026-09-08. Sources: [owner README](../../../README.md),
[BUILD](../../../BUILD.bazel), and
[CRUSH map](../../../crush_map.txt).

## Requirements

### Requirement: Single-host CRUSH topology

The maintained CRUSH map SHALL place four HDD devices, `osd.0` through
`osd.3`, in the `host1` bucket beneath the `default` root, using the declared
device weights and `straw2` bucket algorithm.

#### Scenario: Inspect the recorded storage topology

- **WHEN** a maintainer reads the checked-in CRUSH map
- **THEN** all four OSDs belong to `host1`, and the `default` root contains
  that host
- **AND** this topology establishes a single-host placement baseline without
  asserting that the running cluster matches it

### Requirement: OSD-level replicated and erasure placement

The maintained map SHALL define a replicated rule using
`chooseleaf firstn 0 type osd` and an erasure rule using
`chooseleaf indep 0 type osd` beneath the `default` root. The erasure rule
SHALL set chooseleaf attempts to 5 and choose attempts to 100.

#### Scenario: Inspect the placement failure domain

- **WHEN** a maintainer evaluates `replicated_rule` or `ceph-ec-data`
- **THEN** placement selects OSDs within the recorded topology
- **AND** the map provides no requirement for replica separation across hosts

### Requirement: Operator workflow documentation

The owner documentation SHALL describe exporting, decompiling, editing,
compiling, and installing a CRUSH map followed by a Ceph status check.
Executing that workflow MUST remain subject to the repository's explicit
infrastructure-operation authorization requirement.

#### Scenario: Use the documented map update procedure

- **WHEN** an authorized operator follows the map update instructions
- **THEN** the procedure obtains the current map before creating its compiled
  replacement and checks `ceph -s` after installation
- **AND** the presence of these instructions alone does not authorize a
  live map update
