## Purpose

Run reproducible experiments on an existing ComfyUI endpoint while maintaining an authorized queue buffer and preserving job ownership, meaningful variation and original outputs across interruptions.

## ADDED Requirements

### Requirement: Buffered execution on an existing endpoint

The runner SHALL accept an existing CLI command and maintain one executing job plus a configurable pending buffer while its authorized producer supplies experiments. Retrieval and reporting SHALL NOT block replenishment. Continuous execution MUST be explicitly enabled.

#### Scenario: Slow retrieval

- **WHEN** a job completes while an earlier artifact download is still running
- **THEN** the runner submits replacement work before waiting for retrieval, preserving the configured outstanding buffer after warm-up

### Requirement: Generic producer and workflow identity

The runner SHALL accept producer-generated native UI or API workflows without model or subject assumptions. It SHALL snapshot the workflow, compute a canonical identity after validated declared exclusions, and reject duplicate identities with bounded producer attempts.

#### Scenario: Seed-only repetition

- **WHEN** a producer repeats a workflow while changing only fields explicitly excluded from identity
- **THEN** the runner rejects that candidate and requests a different experiment without submitting the duplicate

### Requirement: Durable ownership and resumption

The runner SHALL persist owned prompt IDs and private experiment state, prevent concurrent controllers using one run directory, and resume known jobs without submitting them again. It MUST NOT blindly retry an ambiguously acknowledged submission.

#### Scenario: Interrupted controller

- **WHEN** a controller restarts with recorded queued prompt IDs
- **THEN** it observes or retrieves those same jobs and does not create replacements for their already-submitted experiments

#### Scenario: Lost acknowledgment

- **WHEN** submission may have reached the endpoint but no valid prompt ID was retained
- **THEN** the runner records ambiguity and stops automatic submission rather than retrying that experiment

### Requirement: Ownership-scoped control

The runner SHALL provide status and a stop request that drains owned work. Explicit cancellation SHALL target only prompt IDs recorded in that run. It MUST NOT perform instance lifecycle operations or clear another user's queue.

#### Scenario: Foreign work shares the endpoint

- **WHEN** the user requests cancellation of a run's owned jobs
- **THEN** only its recorded prompt IDs are passed to the CLI and unrelated jobs remain untouched

### Requirement: Original artifacts and bounded failures

The runner SHALL retain original downloaded media, workflow snapshots, producer notes, timings and hashes in private run artifacts. It SHALL report failures and producer exhaustion without inventing successful results, and SHALL confine reported media to the experiment's output directory.

#### Scenario: Failed or escaping output

- **WHEN** a workflow fails or a download receipt names a file outside its experiment directory
- **THEN** the runner records the failure and does not publish that file as a successful result

### Requirement: Reusable agent procedure

The project SHALL supply a discoverable skill for authorized ComfyUI experimentation that covers stable directions, finite ready jobs, evaluation before explicit scrap/refine/fork decisions, lineage, provenance and result review without embedding a particular model, artwork or rental.

#### Scenario: A different experiment domain

- **WHEN** an agent uses the skill with an unrelated workflow and producer
- **THEN** it can operate the same generic runner without editing its source or inheriting this task's subject choices

#### Scenario: Review governs successor jobs

- **WHEN** a direction's job finishes
- **THEN** the agent evaluates its original outputs before queuing another revision or independent fork

#### Scenario: A direction is scrapped

- **WHEN** evaluation scraps a direction
- **THEN** its history is retained and it is not generated again unless the user explicitly reopens it

#### Scenario: No reviewed next step exists

- **WHEN** the ready set is exhausted or the user pauses experimentation
- **THEN** new submissions stop without inventing variations to keep the endpoint occupied
