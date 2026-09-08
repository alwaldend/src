# Historical project-agents delta

These requirements describe the archived change at its recorded acceptance point. They were imported without applying them to the current `project-agents` baseline. Current source and baseline requirements own present behavior.

## ADDED Requirements

### Requirement: Deterministic advisory impact plans

Identical candidate and causal inputs MUST produce the same ImpactPlan digest, with selected capabilities, effects, narrow targets, required checks, gaps, cost, and reusable evidence.

#### Scenario: The same plan is requested twice

- **WHEN** the intent, candidate, and causal catalog inputs are unchanged
- **THEN** both plans have the same digest and disclose the same minimum checks and gaps

### Requirement: Exact effect admission

Provider action admission MUST reject unknown actions, stale pre-state, wrong environment, widened authority, weakened safety flags, and missing required cancellation or resource budgets before execution.

#### Scenario: A remote write has stale pre-state

- **WHEN** an action is prepared against an earlier remote generation
- **THEN** admission refuses execution until its exact authority and state are reconciled

### Requirement: Candidate-bound reusable evidence

ValidationSets MUST be immutable and candidate-bound; EvidenceAssertions MUST apply those sets to exact criterion revisions without mutating the validation evidence.

#### Scenario: Evidence is reused for a second criterion

- **WHEN** two assertions reference one validation set
- **THEN** the immutable set retains its digest and each assertion records its own criterion revision and semantic verdict

### Requirement: Bounded resumable full checks

Structured full-repository reports MUST record exact inputs, target-universe counts, coverage, bounds, and interrupted work and MUST resume only when relevant inputs are identical.

#### Scenario: An audit resumes after inputs change

- **WHEN** the candidate or relevant execution inputs differ from the saved audit
- **THEN** the prior progress cannot be represented as a completed check of the changed inputs
