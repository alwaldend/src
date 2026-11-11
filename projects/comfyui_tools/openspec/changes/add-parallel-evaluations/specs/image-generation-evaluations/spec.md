## ADDED Requirements

### Requirement: Independent bounded evaluation workers

The evaluation command SHALL expand compatible hardware/model/workflow/case matrix groups before allocation and bound rental concurrency. It SHALL reuse one rental for all workloads sharing the same hardware and model profile, including multiple workflows, cases and seeds. Each ready rental SHALL immediately queue its assigned work, retrieve completed results and be destroyed independently after all its assigned results are downloaded when automatic cleanup is enabled. It SHALL preserve primary and cleanup failures and SHALL NOT mutate the personal rental.

#### Scenario: Different readiness times

- **WHEN** one rental becomes ready while another is booting
- **THEN** its workload, download and destruction SHALL proceed without a barrier

### Requirement: Recoverable rental ownership

The runner SHALL journal intended labels before creation and record confirmed instance identities. Cleanup SHALL verify exact ownership and confirmed absence. The final sweep and standalone cleanup command SHALL destroy every instance in the exact configured evaluation label namespace, regardless of age or run identifier. They SHALL preserve ordinary rental labels and verify that no matching instances remain. Every normal or cancelled run SHALL execute that sweep by default. An explicit cleanup=false run option SHALL skip both per-rental destruction and the final sweep on every outcome, report cleanup_skipped=true and cleanup_confirmed=false, and retain instance identities for the standalone cleanup command.

#### Scenario: Uncertain creation or lost controller

- **WHEN** creation acknowledgment is lost or the controller terminates
- **THEN** reconciliation SHALL discover the exact labeled rental and cleanup SHALL remain possible without touching unrelated instances

### Requirement: Repeatable cases and honest quality evidence

The two operator-approved cases SHALL include their exact prompts, references and review criteria. Results SHALL include effective workflows, seeds, model/source identities, timings, media and cleanup status. Artistic and semantic scores SHALL remain unscored until actually assessed.

#### Scenario: Valid media with poor subject adherence

- **WHEN** inference returns a valid image or GIF that fails its artistic or motion criteria
- **THEN** the report SHALL distinguish execution success from quality acceptance

### Requirement: Controlled frame comparisons

GIF sampling SHALL support original-reference or previous-frame latents and fixed or incrementing frame seeds. Saved nodes without these controls SHALL retain their previous behavior. Masked workflows SHALL take a bounded, feathered edit region from the selected case. Qwen frame workflows SHALL preserve native reference-image and prompt conditioning. The selected experiment SHALL compare multiple settings while reusing each model rental.

#### Scenario: Repeated reference pose

- **WHEN** the operator selects original-reference latents, fixed noise and identical first/last pose prompts
- **THEN** those frames SHALL use the same initial reference and seed, without recursively accumulated changes

#### Scenario: Explicit retention for workflow debugging

- **WHEN** the operator runs with cleanup=false and work succeeds, fails or is cancelled
- **THEN** the run SHALL retain its allocated rentals and the standalone cleanup command SHALL still remove all matching rentals when explicitly invoked

### Requirement: Native decoded image compatibility

The GIF node SHALL use ComfyUI's native decoded-image normalization, including the single-frame temporal layout used by Qwen's VAE. It SHALL convert frames to a NumPy-compatible CPU representation before GIF serialization.

#### Scenario: Temporal VAE output

- **WHEN** a frame decoder returns a single-frame five-dimensional tensor
- **THEN** GIF output SHALL contain one image per prompt with ordinary image dimensions and no extra temporal axis
