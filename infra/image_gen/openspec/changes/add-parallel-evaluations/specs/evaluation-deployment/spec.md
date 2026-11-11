## ADDED Requirements

### Requirement: Independent evaluation deployment

The evals/ subproject SHALL supply configuration, credential references and deployment composition for the reusable evaluator. Executable code SHALL remain in projects/comfyui_tools. The personal rental's resources and Terraform state SHALL remain independent.

#### Scenario: Evaluation invocation

- **WHEN** an operator explicitly starts an evaluation
- **THEN** only configured evaluation rentals SHALL be created and tracked in the evaluation namespace

### Requirement: Bounded operation and orphan cleanup

Deployment SHALL declare maximum parallelism, price/lifetime limits and a default-enabled final namespace-wide cleanup sweep and standalone interrupted-run reconciliation. The live command SHALL accept cleanup=false for explicit rental retention without changing the standalone cleanup target. Offline validation SHALL NOT rent capacity or deploy services.

#### Scenario: Controller failure

- **WHEN** the evaluation controller is unavailable before cleanup completes
- **THEN** the standalone cleanup target SHALL reconcile that matching owned rental without deleting personal instances
