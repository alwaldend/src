## ADDED Requirements

### Requirement: Dedicated runner scheduling budgets

The dedicated Forgejo runner's managed Bazel configuration SHALL derive its job
limit and CPU scheduling budget from 80% of the detected host CPUs and its memory
scheduling budget from 80% of the detected host RAM. These SHALL be per-invocation
scheduling budgets, not hard process or aggregate host-wide resource limits.
Forgejo runner capacity SHALL remain one job at a time.

#### Scenario: Runner builds or tests the checked-out revision

- **WHEN** a Bazel build or test loads the runner's managed defaults without overrides
- **THEN** the job limit and CPU budget are calculated from `HOST_CPUS*0.8`
- **AND** the memory budget is calculated from `HOST_RAM*0.8`
- **AND** Forgejo schedules at most one job on the runner

#### Scenario: Workstation uses its independent configuration

- **WHEN** a workstation build loads its managed defaults
- **THEN** the runner-specific resource budgets do not change the workstation defaults
