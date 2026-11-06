## Why

Saved Terraform plans can retain secrets after an operation. Agents can review
Terraform's normal text output and approve the plan held by an interactive
apply without creating a plan artifact or a separate filtering layer.

## What Changes

- Require text-mode plan review and interactive apply for agents, without saved
  plans or JSON output; allow normal terminal output in agent transcripts.
- Replace saved-plan instructions and evaluation expectations with the
  interactive workflow, including reporting incompatible saved-plan-only targets.
- Check repository workspaces and temporary directories for leftover plans,
  reporting metadata without displaying values or deleting unrelated files.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-agents`: Require the infrastructure skill to guide agents through
  artifact-free Terraform review while preserving operation authorization.

## Impact

Root agent policy, the canonical repo-infra skill and its evaluation cases.
Reusable Terraform APIs, component target contracts, dependencies, live
infrastructure, existing transcripts, and unrelated task artifacts are unchanged.
