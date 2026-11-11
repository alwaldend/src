## Why

The personal rental is inconvenient for parallel automation. Evaluations need an independent deployment composition with bounded resource use and recoverable cleanup.

## What Changes

- Instantiate the reusable projects/comfyui_tools evaluator under evals/.
- Declare the run matrix, hardware/price/lifetime limits, credentials and final/standalone cleanup here.
- Reuse this component's owned Vault authentication and approved model registry without accessing personal rental state.

## Capabilities

### New Capabilities

- `evaluation-deployment`: bounded ephemeral evaluation composition and cleanup wiring.

### Modified Capabilities

None.

## Impact

All executable code remains in projects/comfyui_tools. This component owns configuration and deployment definitions only. Do not rent, write Vault or deploy a persistent controller/reaper during offline validation.
