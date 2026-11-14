## Why

Interactive ComfyUI experimentation needs a queue of deliberately prepared directions while results are retrieved and reviewed. Each finished direction needs evaluation before it is scrapped, refined or forked; automatic endless refill loses that feedback loop. Ad hoc scripts lose ownership and provenance across interruptions, and the existing evaluator couples its bounded workloads to rental provisioning rather than an already-running ComfyUI endpoint.

## What Changes

- Add a generic experiment command that uses an existing ComfyUI CLI command, maintains a configurable queue buffer, and accepts a separate experiment-producer command.
- Snapshot submitted workflows, reject duplicate workflow identities after declared non-semantic exclusions, and retain owned prompt IDs and media artifacts in a resumable private run journal.
- Refill before slow downloads and reporting, support finite or explicitly continuous runs, and stop or cancel only the run's recorded jobs without managing instances.
- Add an owner-local skill covering direction development, finite scheduling, explicit scrap/refine/fork decisions, lineage, provenance, and result review without prescribing a model or subject.
- Validate the command end to end with isolated CLI fixtures and a live existing-endpoint acceptance run.

## Capabilities

### New Capabilities

- `experiment-queue`: Buffered, resumable, ownership-scoped experiments on an existing ComfyUI endpoint, with producer-driven novelty and durable outputs.

### Modified Capabilities

None.

## Impact

The change belongs to `projects/comfyui_tools`, including its command, internal implementation, E2E fixtures, documentation, skill, and discovery registration. Authentication remains in the caller's existing CLI wrapper. The tool does not provision, stop, restart, or destroy servers; contains no model-specific narrative; and introduces no external dependency.
